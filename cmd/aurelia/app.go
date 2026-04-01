package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/health"
	"github.com/kocar/aurelia/internal/heartbeat"
	"github.com/kocar/aurelia/internal/infra"
	"github.com/kocar/aurelia/internal/markdownbrain"
	"github.com/kocar/aurelia/internal/mcp"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/middleware"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/internal/skill"
	"github.com/kocar/aurelia/internal/telegram"
	"github.com/kocar/aurelia/internal/tools"
	"github.com/kocar/aurelia/internal/voice"
	"github.com/kocar/aurelia/pkg/llm"
	"github.com/kocar/aurelia/pkg/voice/stt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/telebot.v3"
)

// app is the central application container, managing lifecycle and dependencies.
type app struct {
	cfg          *config.AppConfig
	resolver     *runtime.PathResolver
	instanceLock *runtime.InstanceLock

	// Core Infrastructure
	mem         *memory.MemoryManager
	cronStore   *cron.SQLiteCronStore
	mcpManager  *mcp.Manager
	taskStore   *agent.SQLiteTaskStore
	llmProvider closableLLMProvider
	redis       *infra.RedisProvider

	// Features
	porteiro      *middleware.PorteiroMiddleware
	pool          *telegram.BotPool       // S-32: multi-bot pool
	primaryBot    *telegram.BotController // S-32: primary bot reference (aurelia)
	cronScheduler *cron.Scheduler
	cronCtx       context.Context
	cronCancel    context.CancelFunc
	heartbeat     *heartbeat.HeartbeatService
	healthServer  *health.Server

	// Voice Stack
	voiceProcessor *voice.Processor
	voiceCapture   *voice.CaptureWorker
	voiceMirrorDB  closableVoiceMirror

	// Agent Loop (Streaming Support 2026.1)
	agentLoop *agent.Loop
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
	a.agentLoop = loop

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
	db, err := sql.Open("sqlite", a.cfg.DBPath)
	if err != nil {
		return fmt.Errorf("open sqlite db: %w", err)
	}

	// Inicializar provedor de LLM temporário para o summarizer de memória
	p := llm.NewOllamaProvider(a.cfg.OllamaURL, a.cfg.LLMModel)
	wrapper := &memoryWrapper{llm: p}

	mem := memory.NewMemoryManager(db, wrapper)
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

	redis, err := infra.NewRedisProvider(a.cfg)
	if err != nil {
		logger.Warn("falha ao inicializar o Redis (Porteiro operará sem cache)", slog.Any("err", err))
	}
	a.redis = redis

	// [SOTA GODMODE] Build and assign the LLM provider so it is non-nil
	// when registerSpawnAgentTool is called in initFeatures.
	provider, err := buildLLMProvider(a.cfg, a.resolver)
	if err != nil {
		return fmt.Errorf("build llm provider: %w", err)
	}
	a.llmProvider = provider

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

	go performOllamaWarmup(a.cfg)

	cwd, _ := os.Getwd()
	if err := runtime.BootstrapProject(cwd); err != nil {
		logger.Warn("failed to bootstrap project-local directory", slog.Any("err", err))
	}
	projectSkillsDir := runtime.ProjectSkills(cwd)
	projectSkillOverlayDir := runtime.ProjectSkillOverlay(cwd)

	markdownBrainTool := maybeRegisterMarkdownBrainTool(a.cfg, registry, cwd, a.mem.DB())
	if markdownBrainTool != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			stats, err := markdownBrainTool.Sync(ctx)
			if err != nil {
				logger.Warn("markdown brain bootstrap sync failed", slog.Any("err", err))
				return
			}
			logger.Info("markdown brain bootstrap sync complete",
				slog.Int("repo_docs", stats.RepoDocs),
				slog.Int("vault_docs", stats.VaultDocs),
				slog.Int("synced_docs", stats.SyncedDocs),
				slog.Int("synced_chunks", stats.SyncedChunks),
				slog.Int("removed_docs", stats.RemovedDocs),
			)
		}()
	}

	assembler := memory.NewContextAssembler(
		a.cfg.QdrantURL, a.cfg.QdrantAPIKey,
		a.cfg.QdrantCollection, markdownbrain.DefaultCollection,
		a.cfg.QdrantEmbeddingModel,
		a.cfg.OllamaURL, a.mem,
	)

	loop := agent.NewLoop(a.llmProvider, registry, a.cfg.MaxIterations).
		WithMemoryAssembler(assembler).
		WithToolCatalog(agent.NewToolCatalog(registry), 7)

	skillInstaller := skill.NewInstaller(a.resolver.Skills(), projectSkillsDir)

	semanticRouter := skill.NewSemanticRouter(
		a.cfg.QdrantURL, a.cfg.QdrantAPIKey, "aurelia_skills",
		a.cfg.QdrantEmbeddingModel, a.cfg.OllamaURL,
	)

	skillLoader := skill.NewLoader(a.resolver.Skills(), projectSkillOverlayDir, projectSkillsDir)
	if auditReport, err := skill.AuditCatalog(a.resolver.Skills(), projectSkillOverlayDir, projectSkillsDir); err != nil {
		logger.Warn("skill catalog audit failed", slog.Any("err", err))
	} else {
		logger.Info("skill catalog audited",
			slog.Int("skills", auditReport.ScannedSkills),
			slog.Int("warnings", auditReport.WarningCount()),
			slog.Int("errors", auditReport.ErrorCount()),
		)
		for _, issue := range auditReport.Issues {
			logger.Warn("skill catalog issue",
				slog.String("severity", issue.Severity),
				slog.String("code", issue.Code),
				slog.String("skill", issue.SkillName),
				slog.String("path", issue.Path),
				slog.String("message", issue.Message),
			)
		}
	}
	if loadedSkills, err := skillLoader.LoadAll(); err == nil {
		go semanticRouter.SyncSkills(context.Background(), loadedSkills)
	}

	installSkillTool := tools.NewInstallSkillTool(skillInstaller)
	registry.Register(installSkillTool.Definition(), installSkillTool.Execute)

	// Porteiro
	a.porteiro = middleware.NewPorteiroMiddleware()
	logger.Info("Porteiro SOTA 2026 (In-Memory) inicializado com sucesso")

	return loop, nil
}

func (a *app) initFeatures(loop *agent.Loop, logger *slog.Logger) error {
	logger.Info("Aurelia Daemon starting",
		slog.String("mode", string(a.cfg.AureliaMode)),
		slog.String("version", "starlight-2026"))
	cwd, _ := os.Getwd()
	projectPlaybookPath := filepath.Join(cwd, "docs", "PROJECT_PLAYBOOK.md")
	projectSkillsDir := runtime.ProjectSkills(cwd)
	projectSkillOverlayDir := runtime.ProjectSkillOverlay(cwd)

	canonicalService := persona.NewCanonicalIdentityService(
		a.mem,
		filepath.Join(a.resolver.MemoryPersonas(), "IDENTITY.md"),
		filepath.Join(a.resolver.MemoryPersonas(), "SOUL.md"),
		filepath.Join(a.resolver.MemoryPersonas(), "USER.md"),
		filepath.Join(a.resolver.Memory(), "OWNER_PLAYBOOK.md"),
		filepath.Join(a.resolver.Memory(), "LESSONS_LEARNED.md"),
		projectPlaybookPath,
	)

	transcriber, err := stt.NewTranscriber(a.cfg.STTProvider, a.cfg.STTBaseURL, a.cfg.STTModel, a.cfg.STTLanguage, a.cfg.GroqAPIKey)
	if err != nil {
		return fmt.Errorf("initialize transcriber: %w", err)
	}

	skillLoader := skill.NewLoader(a.resolver.Skills(), projectSkillOverlayDir, projectSkillsDir)
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
	notificationBot := selectNotificationBot(a.cfg, a.pool, a.primaryBot, logger)

	// S-27: Wire squad + cron status reporters for /status command
	a.primaryBot.SetSquadReporter(squadStatusAdapter{})
	a.primaryBot.SetCronJobReporter(&cronNextJobAdapter{store: a.cronStore})

	ollamaURL := a.cfg.OllamaURL
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	if a.porteiro != nil {
		a.primaryBot.SetPorteiro(a.porteiro)
		logger.Info("Proteção de entrada ativa - Porteiro (In-Memory + Regex)")
	}

	// InputGuard removed: deduplication is handled by Porteiro middleware.

	if err := registerSpawnAgentTool(a.cfg, loop.Registry(), a.llmProvider, notificationBot, a.taskStore); err != nil {
		return fmt.Errorf("register spawn agent tool: %w", err)
	}

	// Squad sync removed (Module Pruned)

	cronScheduler, cronCtx, cronCancel, err := buildCronScheduler(a.cronStore, loop, canonicalService, notificationBot)
	if err != nil {
		return fmt.Errorf("initialize cron scheduler: %w", err)
	}
	a.cronScheduler = cronScheduler
	a.cronCtx = cronCtx
	a.cronCancel = cronCancel

	seedSystemCrons(context.Background(), a.cronStore, a.cfg.VoiceReplyChatID)
	seedMarkdownBrainCron(context.Background(), a.cronStore, a.cfg, a.cfg.VoiceReplyChatID)
	seedRepoGuardianCron(context.Background(), a.cronStore, a.cfg.VoiceReplyChatID)
	seedControleDBCron(context.Background(), a.cronStore, a.cfg.VoiceReplyChatID)
	reconcileSystemCronCadence(context.Background(), a.cronStore, a.cfg)

	// Sentinel health probes removed (Simplified in SOTA 2026)

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

	transcriber, _ := stt.NewTranscriber(a.cfg.STTProvider, a.cfg.STTBaseURL, a.cfg.STTModel, a.cfg.STTLanguage, a.cfg.GroqAPIKey)
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
			logger.Warn("voice capture enabled but command is missing; leaving capture runtime disabled")
			return nil
		}
		if missing := voice.MissingCommandPath(a.cfg.VoiceCaptureCommand); missing != "" {
			logger.Warn("voice capture command path is unavailable; leaving capture runtime disabled", slog.String("path", missing))
			return nil
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

func selectNotificationBot(cfg *config.AppConfig, pool *telegram.BotPool, primary *telegram.BotController, logger *slog.Logger) *telegram.BotController {
	if primary == nil || cfg == nil || pool == nil {
		return primary
	}
	botID := strings.TrimSpace(cfg.TelegramNotificationBotID)
	if botID == "" {
		return primary
	}
	if bc := pool.Get(botID); bc != nil {
		return bc
	}
	logger.Warn("configured notification bot not found; falling back to primary", slog.String("bot_id", botID))
	return primary
}

func (a *app) initServers(logger *slog.Logger) {
	healthSrv := health.NewServer(a.cfg.HealthPort)
	healthSrv.RegisterRoute("/metrics", promhttp.Handler())

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
			BotID  string `json:"bot_id"` // S-32: route to specific bot; empty = primary
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

		// Select target bot: named bot if bot_id provided, otherwise primary.
		target := a.primaryBot
		if req.BotID != "" {
			if bc := a.pool.Get(req.BotID); bc != nil {
				target = bc
			}
		}

		ctx := r.Context()
		logger.Info("Telegram impersonation request received", slog.Int64("user_id", req.UserID), slog.String("text", req.Text))

		err := target.ProcessExternalInput(ctx, req.UserID, req.ChatID, req.Text, false)
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
	// Dashboard routes removed (Module Pruned)
}

func maskTokenPreview(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	if len(token) <= 10 {
		return "***"
	}
	return token[:6] + "..." + token[len(token)-4:]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func (a *app) effectiveBotLLM(botCfg config.BotConfig, appCfg *config.AppConfig) (string, string) {
	if strings.TrimSpace(botCfg.LLMProvider) != "" || strings.TrimSpace(botCfg.LLMModel) != "" {
		return telegram.EffectiveBotLLM(botCfg, appCfg.LLMProvider, appCfg.LLMModel)
	}
	snapshot := buildLLMRuntimeSnapshot(a, time.Now().UTC())
	return telegram.EffectiveBotLLM(botCfg, snapshot.EffectiveProvider, snapshot.EffectiveModel)
}

// Chat API routes removed
func (a *app) registerChatAPIRoutes() {}

func buildLLMProvider(cfg *config.AppConfig, resolver *runtime.PathResolver) (closableLLMProvider, error) {
	var tiers []agent.LLMProvider

	// Tier 1: Primary Cloud Provider (Sovereign Choice)
	primary, err := instantiateProvider(cfg.LLMProvider, cfg)
	if err == nil && primary != nil {
		tiers = append(tiers, primary)
	}

	// Tier 2: OpenRouter Failover (if not already primary)
	if cfg.LLMProvider != "openrouter" && cfg.OpenRouterAPIKey != "" {
		or, err := instantiateProvider("openrouter", cfg)
		if err == nil && or != nil {
			tiers = append(tiers, or)
		}
	}

	// Tier 0: Local Sovereignty (Ollama)
	// Always present as final fallback to ensure zero-drift autonomy.
	if cfg.LLMProvider != "ollama" {
		local := llm.NewOllamaProvider(cfg.OllamaURL, cfg.LLMModel)
		tiers = append(tiers, local)
	}

	if len(tiers) == 0 {
		return nil, fmt.Errorf("nenhum provedor de LLM disponível (verifique as chaves de API)")
	}

	router := agent.NewTieredRouter(tiers...)
	return &closableTieredRouter{TieredRouter: router, tiers: tiers}, nil
}

type closableTieredRouter struct {
	*agent.TieredRouter
	tiers []agent.LLMProvider
}

func (r *closableTieredRouter) Close() {
	for _, t := range r.tiers {
		if c, ok := t.(interface{ Close() }); ok {
			c.Close()
		}
	}
}

func instantiateProvider(name string, cfg *config.AppConfig) (agent.LLMProvider, error) {
	switch name {
	case "anthropic":
		return llm.NewAnthropicProvider(cfg.AnthropicAPIKey, cfg.LLMModel), nil
	case "ollama":
		return llm.NewOllamaProvider(cfg.OllamaURL, cfg.LLMModel), nil
	case "openrouter":
		// Use OpenRouter with default provider or custom one if needed
		return llm.NewOpenRouterProvider(cfg.OpenRouterAPIKey, cfg.LLMModel), nil
	case "openai":
		return llm.NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.LLMModel), nil
	case "groq":
		return llm.NewGroqProvider(cfg.GroqAPIKey, cfg.LLMModel), nil
	case "gemini":
		return llm.NewGeminiProvider(cfg.GoogleAPIKey, cfg.LLMModel), nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", name)
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
		if strings.TrimSpace(output) == "" && len(parts) == 0 {
			return nil
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

// memoryWrapper satisfaz a interface memory.Summarizer usando um agent.LLMProvider.
type memoryWrapper struct {
llm agent.LLMProvider
}

func (w *memoryWrapper) Summarize(ctx context.Context, history string) (string, error) {
prompt := "Sintetize os pontos principais da conversa de forma extremamente concisa. Responda apenas o resumo Markdown."
msgs := []agent.Message{{Role: "user", Content: history}}
resp, err := w.llm.GenerateContent(ctx, prompt, msgs, nil)
	if err != nil {
		return "", err
	}
return resp.Content, nil
}
