package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/dashboard"
	"github.com/kocar/aurelia/internal/gateway"
	"github.com/kocar/aurelia/internal/health"
	"github.com/kocar/aurelia/internal/heartbeat"
	"github.com/kocar/aurelia/internal/mcp"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/metrics"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/internal/skill"
	"github.com/kocar/aurelia/internal/telegram"
	"github.com/kocar/aurelia/internal/tools"
	"github.com/kocar/aurelia/internal/voice"
	"github.com/kocar/aurelia/pkg/llm"
	"github.com/kocar/aurelia/pkg/stt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/telebot.v3"
)

// app is the central application container, managing lifecycle and dependencies.
type app struct {
	cfg            *config.AppConfig
	resolver       *runtime.PathResolver
	instanceLock   *runtime.InstanceLock

	// Core Infrastructure
	mem            *memory.MemoryManager
	cronStore      *cron.SQLiteCronStore
	mcpManager     *mcp.Manager
	taskStore      *agent.SQLiteTaskStore
	llmProvider    closableLLMProvider

	// Features
	pool           *telegram.BotPool          // S-32: multi-bot pool
	primaryBot     *telegram.BotController    // S-32: primary bot reference (aurelia)
	cronScheduler  *cron.Scheduler
	cronCtx        context.Context
	cronCancel     context.CancelFunc
	heartbeat      *heartbeat.HeartbeatService
	healthServer   *health.Server

	// Voice Stack
	voiceProcessor *voice.Processor
	voiceCapture   *voice.CaptureWorker
	voiceMirrorDB  closableVoiceMirror
}

type closableLLMProvider interface {
	agent.LLMProvider
	Close()
}

type closableVoiceMirror interface {
	voice.Mirror
	Close() error
}

// bootstrapApp initializes the entire application stack in a modular fashion.
func bootstrapApp(args []string) (*app, error) {
	logger := observability.Logger("cmd.bootstrap")
	a := &app{}

	// 1. Initial Infrastructure
	if err := a.initInfrastructure(args, logger); err != nil {
		return nil, err
	}

	// 2. Core Dependencies
	if err := a.initCore(logger); err != nil {
		a.close() // Cleanup partially initialized state
		return nil, err
	}

	// 3. Skills and Intelligence
	loop, err := a.initSkills(logger)
	if err != nil {
		a.close()
		return nil, err
	}

	// 4. Features & Interfaces
	if err := a.initFeatures(loop, logger); err != nil {
		a.close()
		return nil, err
	}

	// 5. Voice Processing (Optional)
	if err := a.initVoice(logger); err != nil {
		a.close()
		return nil, err
	}

	// 6. Servers & Connectivity
	a.initServers(logger)

	return a, nil
}

func (a *app) initInfrastructure(args []string, logger *slog.Logger) error {
	resolver, err := runtime.New()
	if err != nil {
		return fmt.Errorf("resolve instance root: %w", err)
	}
	a.resolver = resolver

	if err := runtime.Bootstrap(resolver); err != nil {
		return fmt.Errorf("bootstrap instance directory: %w", err)
	}

	instanceLock, err := runtime.AcquireInstanceLock(resolver, args)
	if err != nil {
		return fmt.Errorf("acquire instance lock: %w", err)
	}
	a.instanceLock = instanceLock

	cfg, err := config.Load(resolver)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	a.cfg = cfg

	return nil
}

func (a *app) initCore(logger *slog.Logger) error {
	mem, err := memory.NewMemoryManager(a.cfg.DBPath, a.cfg.MemoryWindowSize)
	if err != nil {
		return fmt.Errorf("initialize memory manager: %w", err)
	}
	a.mem = mem

	cronStore, err := cron.NewSQLiteCronStore(a.cfg.DBPath)
	if err != nil {
		return fmt.Errorf("initialize cron store: %w", err)
	}
	a.cronStore = cronStore

	taskStore, err := agent.NewSQLiteTaskStore(a.cfg.DBPath + ".teams")
	if err != nil {
		return fmt.Errorf("initialize team task store: %w", err)
	}
	a.taskStore = taskStore

	llmProvider, err := buildLLMProvider(a.cfg, a.resolver)
	if err != nil {
		return fmt.Errorf("initialize llm provider: %w", err)
	}
	a.llmProvider = llmProvider

	return nil
}

func (a *app) initSkills(logger *slog.Logger) (*agent.Loop, error) {
	registry := buildToolRegistry()
	registerScheduleTools(registry, a.cronStore)

	mcpManager, err := maybeRegisterMCPTools(a.cfg, registry)
	if err != nil {
		logger.Warn("failed to initialize MCP tools (degraded mode likely)", slog.Any("err", err))
	}
	a.mcpManager = mcpManager

	// Perform Preflight
	if gwProvider, ok := a.llmProvider.(*gateway.Provider); ok {
		ctxP, cancelP := context.WithTimeout(context.Background(), 3*time.Second)
		if err := performPreflightCheck(ctxP, a.cfg); err != nil {
			logger.Warn("preflight check failed, enabling degraded mode", slog.Any("err", err))
			gwProvider.EnableDegradedMode("Preflight failures: " + err.Error())
		}
		cancelP()
	}

	go performOllamaWarmup(a.cfg)

	assembler := memory.NewContextAssembler(
		a.cfg.QdrantURL, a.cfg.QdrantAPIKey,
		a.cfg.QdrantCollection, a.cfg.QdrantEmbeddingModel,
		a.cfg.OllamaURL, a.mem,
	)

	loop := agent.NewLoop(a.llmProvider, registry, a.cfg.MaxIterations).
		WithMemoryAssembler(assembler).
		WithToolCatalog(agent.NewToolCatalog(registry), 7)

	cwd, _ := os.Getwd()
	if err := runtime.BootstrapProject(cwd); err != nil {
		logger.Warn("failed to bootstrap project-local directory", slog.Any("err", err))
	}
	projectSkillsDir := runtime.ProjectSkills(cwd)

	skillInstaller := skill.NewInstaller(a.resolver.Skills(), projectSkillsDir)

	semanticRouter := skill.NewSemanticRouter(
		a.cfg.QdrantURL, a.cfg.QdrantAPIKey, "aurelia_skills",
		a.cfg.QdrantEmbeddingModel, a.cfg.OllamaURL,
	)
	loop.WithSemanticRouter(semanticRouter)

	skillLoader := skill.NewLoader(a.resolver.Skills(), projectSkillsDir)
	if loadedSkills, err := skillLoader.LoadAll(); err == nil {
		go semanticRouter.SyncSkills(context.Background(), loadedSkills)
	}

	installSkillTool := tools.NewInstallSkillTool(skillInstaller)
	registry.Register(installSkillTool.Definition(), installSkillTool.Execute)

	return loop, nil
}

func (a *app) initFeatures(loop *agent.Loop, logger *slog.Logger) error {
	cwd, _ := os.Getwd()
	projectPlaybookPath := filepath.Join(cwd, "docs", "PROJECT_PLAYBOOK.md")
	projectSkillsDir := runtime.ProjectSkills(cwd)

	canonicalService := persona.NewCanonicalIdentityService(
		a.mem,
		filepath.Join(a.resolver.MemoryPersonas(), "IDENTITY.md"),
		filepath.Join(a.resolver.MemoryPersonas(), "SOUL.md"),
		filepath.Join(a.resolver.MemoryPersonas(), "USER.md"),
		filepath.Join(a.resolver.Memory(), "OWNER_PLAYBOOK.md"),
		filepath.Join(a.resolver.Memory(), "LESSONS_LEARNED.md"),
		projectPlaybookPath,
	)

	transcriber, err := stt.NewTranscriber(a.cfg.STTProvider, a.cfg.GroqAPIKey, a.cfg.STTBaseURL, a.cfg.STTModel, a.cfg.STTLanguage)
	if err != nil {
		return fmt.Errorf("initialize transcriber: %w", err)
	}

	skillLoader := skill.NewLoader(a.resolver.Skills(), projectSkillsDir)
	skillRouter := skill.NewRouter(a.llmProvider)
	skillExecutor := skill.NewExecutor(loop)

	// S-32: Build multi-bot pool
	pool := telegram.NewBotPool(
		a.cfg, a.mem, skillRouter, skillExecutor, skillLoader,
		transcriber, canonicalService, a.resolver.MemoryPersonas(),
	)
	for _, botCfg := range a.cfg.Bots {
		if err := pool.Add(botCfg); err != nil {
			logger.Warn("failed to add bot to pool", slog.String("bot_id", botCfg.ID), slog.Any("err", err))
		}
	}
	// Fallback: if pool is empty but primary token exists, add aurelia directly
	if pool.Size() == 0 && a.cfg.TelegramBotToken != "" {
		bc, err := telegram.NewBotController(
			a.cfg, a.mem, skillRouter, skillExecutor, skillLoader,
			transcriber, canonicalService, a.resolver.MemoryPersonas(),
		)
		if err != nil {
			return fmt.Errorf("initialize telegram: %w", err)
		}
		_ = pool.AddController("aurelia", bc)
	}
	a.pool = pool
	a.primaryBot = pool.Primary()

	if a.primaryBot == nil {
		return fmt.Errorf("no telegram bot configured (set telegram_bot_token or bots in config)")
	}

	if gw, ok := a.llmProvider.(*gateway.Provider); ok {
		a.primaryBot.SetHealthReporter(gw)
	}

	// S-27: Wire squad + cron status reporters for /status command
	a.primaryBot.SetSquadReporter(squadStatusAdapter{})
	a.primaryBot.SetCronJobReporter(&cronNextJobAdapter{store: a.cronStore})

	// Wire gemma3 input guard
	ollamaURL := a.cfg.OllamaURL
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	a.primaryBot.SetInputGuard(telegram.NewInputGuard(ollamaURL))

	if err := registerSpawnAgentTool(a.cfg, loop.Registry(), a.llmProvider, a.primaryBot, a.taskStore); err != nil {
		return fmt.Errorf("register spawn agent tool: %w", err)
	}

	// S-32: Sync bots to squad
	botEntries := make([]agent.BotConfigEntry, 0, len(a.cfg.Bots))
	for _, b := range a.cfg.Bots {
		botEntries = append(botEntries, agent.BotConfigEntry{
			ID: b.ID, Name: b.Name, PersonaID: b.PersonaID, FocusArea: b.FocusArea,
		})
	}
	agent.SyncBotsToSquad(botEntries)

	cronScheduler, cronCtx, cronCancel, err := buildCronScheduler(a.cronStore, loop, canonicalService, a.primaryBot)
	if err != nil {
		return fmt.Errorf("initialize cron scheduler: %w", err)
	}
	a.cronScheduler = cronScheduler
	a.cronCtx = cronCtx
	a.cronCancel = cronCancel

	seedSystemCrons(context.Background(), a.cronStore, a.cfg.VoiceReplyChatID)

	// S-22: Squad Live Load — updates agent metrics every 10s
	agent.StartLiveLoad(a.cronScheduler, ollamaURL, a.cfg.OpenRouterAPIKey)

	// S-24: Sentinel health probes — updates gemma/sentinel squad status every 30s
	startSentinelHealthLoop(a.cfg)

	a.heartbeat = heartbeat.NewHeartbeatService(
		a.resolver.Root(),
		a.cfg.HeartbeatIntervalMinutes,
		a.cfg.HeartbeatEnabled,
		loop,
	)

	return nil
}

func (a *app) initVoice(logger *slog.Logger) error {
	if !a.cfg.VoiceEnabled {
		return nil
	}

	voiceSpool, err := voice.NewSpool(a.cfg.VoiceSpoolPath)
	if err != nil {
		return fmt.Errorf("initialize voice spool: %w", err)
	}

	transcriber, _ := stt.NewTranscriber(a.cfg.STTProvider, a.cfg.GroqAPIKey, a.cfg.STTBaseURL, a.cfg.STTModel, a.cfg.STTLanguage)
	var fallback stt.Transcriber
	if a.cfg.STTFallbackCommand != "" {
		fallback = stt.NewCommandTranscriber(a.cfg.STTFallbackCommand)
	}

	a.voiceMirrorDB = voice.NewSQLiteMirror(a.cfg.DBPath)
	voiceMirror := voice.NewMultiMirror(
		a.voiceMirrorDB,
		voice.NewQdrantMirror(
			a.cfg.QdrantURL, a.cfg.QdrantAPIKey,
			a.cfg.QdrantCollection, a.cfg.QdrantEmbeddingModel,
			a.cfg.OllamaURL,
		),
	)

	a.voiceProcessor = voice.NewProcessor(
		voiceSpool, transcriber, fallback,
		voiceBotDispatcher{bot: a.primaryBot},
		voice.Config{
			PollInterval:       time.Duration(a.cfg.VoicePollIntervalMS) * time.Millisecond,
			HeartbeatPath:      a.cfg.VoiceHeartbeatPath,
			HeartbeatFreshness: time.Duration(a.cfg.VoiceHeartbeatFreshSec) * time.Second,
			WakePhrase:         a.cfg.VoiceWakePhrase,
			DefaultUserID:      a.cfg.VoiceReplyUserID,
			DefaultChatID:      a.cfg.VoiceReplyChatID,
			SoftCapDaily:       a.cfg.GroqSoftCapDaily,
			HardCapDaily:       a.cfg.GroqHardCapDaily,
			PrimaryLabel:       a.cfg.STTProvider,
			Mirror:             voiceMirror,
		},
	)

	if a.cfg.VoiceCaptureEnabled {
		if a.cfg.VoiceCaptureCommand == "" {
			return fmt.Errorf("voice capture is enabled but voice_capture_command is missing")
		}
		a.voiceCapture = voice.NewCaptureWorker(
			voiceSpool,
			voice.NewCommandCaptureSource(a.cfg.VoiceCaptureCommand, map[string]string{
				"AURELIA_VOICE_WAKE_PHRASE": a.cfg.VoiceWakePhrase,
				"AURELIA_VOICE_DROP_PATH":   a.cfg.VoiceDropPath,
				"AURELIA_VOICE_USER_ID":     strconv.FormatInt(a.cfg.VoiceReplyUserID, 10),
				"AURELIA_VOICE_CHAT_ID":     strconv.FormatInt(a.cfg.VoiceReplyChatID, 10),
			}),
			voice.CaptureConfig{
				PollInterval:       time.Duration(a.cfg.VoiceCapturePollMS) * time.Millisecond,
				HeartbeatPath:      a.cfg.VoiceCaptureHeartbeat,
				HeartbeatFreshness: time.Duration(a.cfg.VoiceCaptureFreshSec) * time.Second,
				DefaultUserID:      a.cfg.VoiceReplyUserID,
				DefaultChatID:      a.cfg.VoiceReplyChatID,
				DefaultSource:      "capture",
			},
		)
	}

	return nil
}

func (a *app) initServers(logger *slog.Logger) {
	healthSrv := health.NewServer(a.cfg.HealthPort)
	registerAuxiliaryHealthChecks(healthSrv, a.cfg, a.resolver, a.llmProvider)

	healthSrv.RegisterRoute("/metrics", promhttp.Handler())
	healthSrv.RegisterRoute("/v1/router/dry-run", gateway.NewDryRunHandler(gateway.NewPlanner()))

	if gw, ok := a.llmProvider.(*gateway.Provider); ok {
		healthSrv.RegisterRoute("/v1/router/status", gw.StatusHandler())
	}
	if a.voiceProcessor != nil {
		healthSrv.RegisterCheck("voice_processor", a.voiceProcessor.HealthCheck)
		healthSrv.RegisterRoute("/v1/voice/status", a.voiceProcessor.StatusHandler())
	}
	if a.voiceCapture != nil {
		healthSrv.RegisterCheck("voice_capture", a.voiceCapture.HealthCheck)
		healthSrv.RegisterRoute("/v1/voice/capture/status", a.voiceCapture.StatusHandler())
	}

	// S-28: Telegram Impersonation for CLI/Automation
	healthSrv.RegisterRoute("/v1/telegram/impersonate", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			UserID int64  `json:"user_id"`
			ChatID int64  `json:"chat_id"`
			Text   string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		// Fallback para o Master se omitido
		if req.UserID == 0 {
			req.UserID = a.cfg.VoiceReplyUserID
		}
		if req.ChatID == 0 {
			req.ChatID = a.cfg.VoiceReplyChatID
		}

		ctx := r.Context()
		logger.Info("Telegram impersonation request received", slog.Int64("user_id", req.UserID), slog.String("text", req.Text))

		err := a.primaryBot.ProcessExternalInput(ctx, req.UserID, req.ChatID, req.Text, false)
		if err != nil {
			http.Error(w, "Pipeline execution failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok", "message": "Pipeline processing started"}`))
	}))

	a.healthServer = healthSrv
}

// registerBotAPIRoutes registers S-32 multi-bot REST endpoints on the dashboard server.
func (a *app) registerBotAPIRoutes(resolver *runtime.PathResolver, appCfg *config.AppConfig) {
	corsJSON := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	// GET /api/bots — list all bots with running status
	dashboard.RegisterRoute("/api/bots", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			corsJSON(w)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		corsJSON(w)
		type botStatus struct {
			config.BotConfig
			Running bool `json:"running"`
		}
		cfgs := a.pool.Configs()
		out := make([]botStatus, 0, len(cfgs))
		for _, bc := range cfgs {
			running := a.pool.Get(bc.ID) != nil
			out = append(out, botStatus{BotConfig: bc, Running: running})
		}
		_ = json.NewEncoder(w).Encode(out)
	}))

	// POST /api/bots — create a new bot
	dashboard.RegisterRoute("/api/bots/create", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			corsJSON(w)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		corsJSON(w)
		var req config.BotConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.ID == "" || req.Token == "" {
			http.Error(w, "id and token are required", http.StatusBadRequest)
			return
		}
		req.Enabled = true
		if err := a.pool.Add(req); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		// Persist
		all := a.pool.Configs()
		_ = config.SaveBots(resolver, all)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(req)
	}))

	// DELETE /api/bots/remove?id=xxx
	dashboard.RegisterRoute("/api/bots/remove", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			corsJSON(w)
			return
		}
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		corsJSON(w)
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query param required", http.StatusBadRequest)
			return
		}
		a.pool.Remove(id)
		all := a.pool.Configs()
		_ = config.SaveBots(resolver, all)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"removed"}`))
	}))

	// GET /api/personas — list persona templates
	dashboard.RegisterRoute("/api/personas", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			corsJSON(w)
			return
		}
		corsJSON(w)
		_ = json.NewEncoder(w).Encode(persona.BuiltinTemplates())
	}))
}

func buildLLMProvider(cfg *config.AppConfig, resolver *runtime.PathResolver) (closableLLMProvider, error) {
	if (cfg.LLMProvider == "openrouter" || cfg.LLMProvider == "ollama") && cfg.OpenRouterAPIKey != "" {
		return gateway.NewProvider(cfg)
	}
	switch cfg.LLMProvider {
	case "anthropic":
		return llm.NewAnthropicProvider(cfg.AnthropicAPIKey, cfg.LLMModel), nil
	case "google":
		return llm.NewGeminiProvider(context.Background(), cfg.GoogleAPIKey, cfg.LLMModel)
	case "ollama":
		return llm.NewOllamaProvider(cfg.OllamaURL, cfg.LLMModel), nil
	case "openrouter":
		return llm.NewOpenRouterProvider(cfg.OpenRouterAPIKey, cfg.LLMModel), nil
	case "openai":
		return llm.NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.LLMModel), nil
	case "groq":
		return llm.NewGroqProvider(cfg.GroqAPIKey, cfg.LLMModel), nil
	default:
		return nil, fmt.Errorf("unsupported llm provider %q", cfg.LLMProvider)
	}
}


func (a *app) start() {
	logger := observability.Logger("cmd.app")
	go func() {
		if err := a.cronScheduler.Start(a.cronCtx); err != nil && err != context.Canceled {
			logger.Warn("cron scheduler stopped with error", slog.Any("err", err))
		}
	}()
	if a.heartbeat != nil {
		_ = a.heartbeat.Start()
	}
	if a.healthServer != nil {
		_ = a.healthServer.Start()
	}
	if a.voiceProcessor != nil {
		_ = a.voiceProcessor.Start()
	}
	if a.voiceCapture != nil {
		_ = a.voiceCapture.Start()
	}

	dashboard.RegisterRoute("/api/squad", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_ = json.NewEncoder(w).Encode(agent.GetFixedSquad())
	})
	// Expor gateway status ao dashboard
	if gw, ok := a.llmProvider.(*gateway.Provider); ok {
		dashboard.RegisterRoute("/api/router/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			_ = json.NewEncoder(w).Encode(gw.StatusSnapshot())
		})
	}
	dashboard.RegisterRoute("/api/commands", dashboard.HandleCommands)
	dashboard.RegisterRoute("/api/metrics", metrics.Handler())

	// S-26: The Brain — Qdrant semantic search endpoints
	qdrantURL := a.cfg.QdrantURL
	if qdrantURL == "" {
		qdrantURL = "http://localhost:6333"
	}
	qdrantCollection := a.cfg.QdrantCollection
	if qdrantCollection == "" {
		qdrantCollection = "conversation_memory"
	}
	qdrantAPIKey := a.cfg.QdrantAPIKey
	dashboard.RegisterRoute("/api/brain/search", buildBrainSearchHandler(qdrantURL, qdrantCollection, qdrantAPIKey))
	dashboard.RegisterRoute("/api/brain/recent", buildBrainRecentHandler(qdrantURL, qdrantCollection, qdrantAPIKey))

	// Proxy VRV homelab dashboard (/api/vrv/ → http://localhost:3333/)
	dashboard.RegisterRoute("/api/vrv/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subPath := strings.TrimPrefix(r.URL.Path, "/api/vrv")
		if subPath == "" || subPath == "/" {
			subPath = "/"
		}
		targetURL := "http://localhost:3333" + subPath
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}
		resp, err := http.Get(targetURL) //nolint:noctx
		if err != nil {
			http.Error(w, "VRV backend unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "VRV read error: "+err.Error(), http.StatusBadGateway)
			return
		}
		ct := resp.Header.Get("Content-Type")
		if strings.Contains(ct, "text/html") {
			body = bytes.ReplaceAll(body, []byte(`href="/`), []byte(`href="/api/vrv/`))
			body = bytes.ReplaceAll(body, []byte(`src="/`), []byte(`src="/api/vrv/`))
		}
		for k, vs := range resp.Header {
			if k == "Content-Encoding" || k == "Content-Length" || k == "Transfer-Encoding" {
				continue
			}
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(body)
	}))

	// S-32: REST API for multi-bot management
	resolver := a.resolver
	appCfg := a.cfg
	a.registerBotAPIRoutes(resolver, appCfg)

	_ = dashboard.StartServer(logger, a.cfg.DashboardPort)
	a.pool.StartAll()
}

func (a *app) shutdown(ctx context.Context) {
	if a.cronCancel != nil {
		a.cronCancel()
	}
	if a.heartbeat != nil {
		a.heartbeat.Stop()
	}
	if a.healthServer != nil {
		_ = a.healthServer.Stop()
	}
	if a.voiceProcessor != nil {
		a.voiceProcessor.Stop()
	}
	if a.voiceCapture != nil {
		a.voiceCapture.Stop()
	}
	if a.pool != nil {
		a.pool.StopAll()
	}
}

func (a *app) close() {
	if a.taskStore != nil {
		_ = a.taskStore.Close()
	}
	if a.mcpManager != nil {
		_ = a.mcpManager.Close()
	}
	if a.cronStore != nil {
		_ = a.cronStore.Close()
	}
	if a.mem != nil {
		_ = a.mem.Close()
	}
	if a.voiceMirrorDB != nil {
		_ = a.voiceMirrorDB.Close()
	}
	if a.llmProvider != nil {
		a.llmProvider.Close()
	}
	if a.instanceLock != nil {
		_ = a.instanceLock.Release()
	}
}

type voiceBotDispatcher struct {
	bot *telegram.BotController // primary bot for voice dispatch
}

func (d voiceBotDispatcher) DispatchVoice(ctx context.Context, userID, chatID int64, text string, requiresAudio bool) error {
	if d.bot == nil {
		return fmt.Errorf("telegram bot dispatcher unavailable")
	}
	return d.bot.ProcessExternalInput(ctx, userID, chatID, text, requiresAudio)
}

func buildCronScheduler(
	cronStore *cron.SQLiteCronStore,
	loop *agent.Loop,
	canonicalService *persona.CanonicalIdentityService,
	bot *telegram.BotController,
) (*cron.Scheduler, context.Context, context.CancelFunc, error) {
	cronPrompt, cronAllowedTools := loadCronPromptConfig(canonicalService)
	cronRuntime := cron.NewAgentCronRuntimeWithPromptBuilder(
		&loopExecutorAdapter{loop: loop},
		cronPrompt,
		cronAllowedTools,
		func(ctx context.Context, job cron.CronJob) (string, []string, error) {
			if canonicalService == nil {
				return cronPrompt, cronAllowedTools, nil
			}
			return canonicalService.BuildPromptForQuery(ctx, job.OwnerUserID, strconv.FormatInt(job.TargetChatID, 10), job.Prompt)
		},
	)

	notifyingRuntime := cron.NewNotifyingRuntime(cronRuntime, func(ctx context.Context, job cron.CronJob, output string, parts []agent.ContentPart, execErr error) error {
		chat := &telebot.Chat{ID: job.TargetChatID}
		if execErr != nil {
			return telegram.SendText(bot.GetBot(), chat, "Falha na rotina agendada:\n"+execErr.Error())
		}
		if len(parts) > 0 {
			return telegram.SendMediaParts(bot.GetBot(), chat, parts)
		}
		return telegram.SendText(bot.GetBot(), chat, output)
	})

	scheduler, err := cron.NewScheduler(cronStore, notifyingRuntime, nil, cron.SchedulerConfig{PollInterval: time.Minute})
	if err != nil {
		return nil, nil, nil, err
	}

	cronCtx, cancel := context.WithCancel(context.Background())
	return scheduler, cronCtx, cancel, nil
}

func performPreflightCheck(ctx context.Context, cfg *config.AppConfig) error {
	if cfg.OpenRouterAPIKey == "" {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/v1/auth/key", nil)
	if err != nil {
		return fmt.Errorf("req build failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.OpenRouterAPIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unreachable (%w)", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth rejected HTTP %d", resp.StatusCode)
	}
	return nil
}

func performOllamaWarmup(cfg *config.AppConfig) {
	if cfg.LLMProvider != "ollama" {
		return
	}
	logger := observability.Logger("ollama.warmup")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	url := strings.TrimRight(cfg.OllamaURL, "/") + "/api/generate"
	payload := map[string]any{
		"model":  cfg.LLMModel,
		"prompt": "", // Empty prompt just triggers load
		"stream": false,
	}
	body, _ := json.Marshal(payload)
	
	logger.Info("warming up local model", "model", cfg.LLMModel)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		resp.Body.Close()
		logger.Info("model warmed up successfully", "model", cfg.LLMModel)
	} else {
		logger.Warn("model warm up failed", "model", cfg.LLMModel, "err", err)
	}

	// Also warm up embedding model
	embedURL := strings.TrimRight(cfg.OllamaURL, "/") + "/api/embed"
	embedPayload := map[string]any{
		"model": cfg.QdrantEmbeddingModel,
		"input": "warmup",
	}
	ebody, _ := json.Marshal(embedPayload)
	req2, _ := http.NewRequestWithContext(ctx, http.MethodPost, embedURL, bytes.NewReader(ebody))
	if resp2, err := http.DefaultClient.Do(req2); err == nil {
		resp2.Body.Close()
		logger.Info("embedding model warmed up successfully", "model", cfg.QdrantEmbeddingModel)
	}
}
