package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

type app struct {
	resolver       *runtime.PathResolver
	instanceLock   *runtime.InstanceLock
	mem            *memory.MemoryManager
	cronStore      *cron.SQLiteCronStore
	mcpManager     *mcp.Manager
	taskStore      *agent.SQLiteTaskStore
	llmProvider    closableLLMProvider
	bot            *telegram.BotController
	cronScheduler  *cron.Scheduler
	cronCtx        context.Context
	cronCancel     context.CancelFunc
	heartbeat      *heartbeat.HeartbeatService
	healthServer   *health.Server
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

func bootstrapApp(args []string) (*app, error) {
	logger := observability.Logger("cmd.app")

	resolver, err := runtime.New()
	if err != nil {
		return nil, fmt.Errorf("resolve instance root: %w", err)
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		return nil, fmt.Errorf("bootstrap instance directory: %w", err)
	}
	instanceLock, err := runtime.AcquireInstanceLock(resolver, args)
	if err != nil {
		return nil, fmt.Errorf("acquire instance lock: %w", err)
	}

	cfg, err := config.Load(resolver)
	if err != nil {
		_ = instanceLock.Release()
		return nil, fmt.Errorf("load config: %w", err)
	}

	mem, err := memory.NewMemoryManager(cfg.DBPath, cfg.MemoryWindowSize)
	if err != nil {
		_ = instanceLock.Release()
		return nil, fmt.Errorf("initialize memory manager: %w", err)
	}

	personasDir := resolver.MemoryPersonas()
	memoryDir := resolver.Memory()
	ownerPlaybookPath := filepath.Join(memoryDir, "OWNER_PLAYBOOK.md")
	lessonsLearnedPath := filepath.Join(memoryDir, "LESSONS_LEARNED.md")
	cwd, err := os.Getwd()
	if err != nil {
		logger.Warn("failed to resolve working directory for project playbook", slog.Any("err", err))
		cwd = ""
	}
	if err := runtime.BootstrapProject(cwd); err != nil {
		logger.Warn("failed to bootstrap project-local Aurelia directory", slog.Any("err", err))
	}
	var projectPlaybookPath string
	var projectSkillsDir string
	if cwd != "" {
		projectPlaybookPath = filepath.Join(cwd, "docs", "PROJECT_PLAYBOOK.md")
		projectSkillsDir = runtime.ProjectSkills(cwd)
	}
	llmProvider, err := buildLLMProvider(cfg, resolver)
	if err != nil {
		_ = instanceLock.Release()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize llm provider: %w", err)
	}
	canonicalService := persona.NewCanonicalIdentityService(
		mem,
		filepath.Join(personasDir, "IDENTITY.md"),
		filepath.Join(personasDir, "SOUL.md"),
		filepath.Join(personasDir, "USER.md"),
		ownerPlaybookPath,
		lessonsLearnedPath,
		projectPlaybookPath,
	)
	registry := buildToolRegistry()

	cronStore, err := cron.NewSQLiteCronStore(cfg.DBPath)
	if err != nil {
		_ = instanceLock.Release()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize cron store: %w", err)
	}

	registerScheduleTools(registry, cronStore)
	mcpManager, err := maybeRegisterMCPTools(cfg, registry)
	if err != nil {
		logger.Warn("failed to initialize MCP tools", slog.Any("err", err))
	}

	loop := agent.NewLoop(llmProvider, registry, cfg.MaxIterations)
	skillLoader := skill.NewLoader(resolver.Skills(), projectSkillsDir)
	skillRouter := skill.NewRouter(llmProvider)
	skillExecutor := skill.NewExecutor(loop)
	skillInstaller := skill.NewInstaller(resolver.Skills(), projectSkillsDir)
	transcriber, err := buildTranscriber(cfg)
	if err != nil {
		_ = instanceLock.Release()
		_ = cronStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize transcriber: %w", err)
	}

	installSkillTool := tools.NewInstallSkillTool(skillInstaller)
	registry.Register(installSkillTool.Definition(), installSkillTool.Execute)

	bot, err := telegram.NewBotController(
		cfg,
		mem,
		skillRouter,
		skillExecutor,
		skillLoader,
		transcriber,
		canonicalService,
		personasDir,
	)
	if err != nil {
		_ = instanceLock.Release()
		_ = cronStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize telegram block: %w", err)
	}

	var voiceProcessor *voice.Processor
	var voiceCapture *voice.CaptureWorker
	var sqliteMirror closableVoiceMirror
	if cfg.VoiceCaptureEnabled && !cfg.VoiceEnabled {
		_ = instanceLock.Release()
		_ = cronStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("voice capture requires voice_enabled=true")
	}
	if cfg.VoiceEnabled {
		voiceSpool, err := voice.NewSpool(cfg.VoiceSpoolPath)
		if err != nil {
			_ = instanceLock.Release()
			_ = cronStore.Close()
			_ = mem.Close()
			return nil, fmt.Errorf("initialize voice spool: %w", err)
		}
		var fallback stt.Transcriber
		if cfg.STTFallbackCommand != "" {
			fallback = stt.NewCommandTranscriber(cfg.STTFallbackCommand)
		}
		sqliteMirror = voice.NewSQLiteMirror(cfg.DBPath)
		voiceMirror := voice.NewMultiMirror(
			sqliteMirror,
			voice.NewSupabaseMirror(cfg.SupabaseURL, cfg.SupabaseServiceRoleKey, cfg.SupabaseEventsTable),
			voice.NewQdrantMirror(cfg.QdrantURL, cfg.QdrantAPIKey, cfg.QdrantCollection, cfg.QdrantEmbeddingModel),
		)
		voiceProcessor = voice.NewProcessor(voiceSpool, transcriber, fallback, voiceBotDispatcher{bot: bot}, voice.Config{
			PollInterval:       time.Duration(cfg.VoicePollIntervalMS) * time.Millisecond,
			HeartbeatPath:      cfg.VoiceHeartbeatPath,
			HeartbeatFreshness: time.Duration(cfg.VoiceHeartbeatFreshSec) * time.Second,
			WakePhrase:         cfg.VoiceWakePhrase,
			DefaultUserID:      cfg.VoiceReplyUserID,
			DefaultChatID:      cfg.VoiceReplyChatID,
			SoftCapDaily:       cfg.GroqSoftCapDaily,
			HardCapDaily:       cfg.GroqHardCapDaily,
			PrimaryLabel:       cfg.STTProvider,
			Mirror:             voiceMirror,
		})
		if cfg.VoiceCaptureEnabled {
			if cfg.VoiceCaptureCommand == "" {
				_ = instanceLock.Release()
				_ = cronStore.Close()
				_ = mem.Close()
				return nil, fmt.Errorf("voice capture is enabled but voice_capture_command is missing")
			}
			voiceCapture = voice.NewCaptureWorker(
				voiceSpool,
				voice.NewCommandCaptureSource(cfg.VoiceCaptureCommand, map[string]string{
					"AURELIA_VOICE_WAKE_PHRASE": cfg.VoiceWakePhrase,
					"AURELIA_VOICE_DROP_PATH":   cfg.VoiceDropPath,
					"AURELIA_VOICE_USER_ID":     strconv.FormatInt(cfg.VoiceReplyUserID, 10),
					"AURELIA_VOICE_CHAT_ID":     strconv.FormatInt(cfg.VoiceReplyChatID, 10),
				}),
				voice.CaptureConfig{
					PollInterval:       time.Duration(cfg.VoiceCapturePollMS) * time.Millisecond,
					HeartbeatPath:      cfg.VoiceCaptureHeartbeat,
					HeartbeatFreshness: time.Duration(cfg.VoiceCaptureFreshSec) * time.Second,
					DefaultUserID:      cfg.VoiceReplyUserID,
					DefaultChatID:      cfg.VoiceReplyChatID,
					DefaultSource:      "capture",
				},
			)
		}
	}

	taskStore, err := agent.NewSQLiteTaskStore(cfg.DBPath + ".teams")
	if err != nil {
		_ = instanceLock.Release()
		_ = cronStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize team task store: %w", err)
	}

	if err := registerSpawnAgentTool(cfg, registry, llmProvider, bot, taskStore); err != nil {
		_ = instanceLock.Release()
		_ = taskStore.Close()
		_ = cronStore.Close()
		_ = mem.Close()
		return nil, err
	}

	cronScheduler, cronCtx, cronCancel, err := buildCronScheduler(cronStore, loop, canonicalService, bot)
	if err != nil {
		_ = instanceLock.Release()
		_ = taskStore.Close()
		_ = cronStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize cron scheduler: %w", err)
	}

	// Initialize heartbeat service
	hbService := heartbeat.NewHeartbeatService(
		resolver.Root(),
		cfg.HeartbeatIntervalMinutes,
		cfg.HeartbeatEnabled,
		loop,
	)

	// Initialize health server
	healthSrv := health.NewServer(8484) // Default port, could be configurable
	registerAuxiliaryHealthChecks(healthSrv, cfg, resolver, llmProvider)
	healthSrv.RegisterRoute("/metrics", promhttp.Handler())
	healthSrv.RegisterRoute("/v1/router/dry-run", gateway.NewDryRunHandler(gateway.NewPlanner()))
	if gwProvider, ok := llmProvider.(*gateway.Provider); ok {
		healthSrv.RegisterRoute("/v1/router/status", gwProvider.StatusHandler())
	}
	if voiceProcessor != nil {
		healthSrv.RegisterCheck("voice_processor", voiceProcessor.HealthCheck)
		healthSrv.RegisterRoute("/v1/voice/status", voiceProcessor.StatusHandler())
	}
	if voiceCapture != nil {
		healthSrv.RegisterCheck("voice_capture", voiceCapture.HealthCheck)
		healthSrv.RegisterRoute("/v1/voice/capture/status", voiceCapture.StatusHandler())
	}

	return &app{
		resolver:       resolver,
		instanceLock:   instanceLock,
		mem:            mem,
		cronStore:      cronStore,
		mcpManager:     mcpManager,
		taskStore:      taskStore,
		llmProvider:    llmProvider,
		bot:            bot,
		cronScheduler:  cronScheduler,
		cronCtx:        cronCtx,
		cronCancel:     cronCancel,
		heartbeat:      hbService,
		healthServer:   healthSrv,
		voiceProcessor: voiceProcessor,
		voiceCapture:   voiceCapture,
		voiceMirrorDB:  sqliteMirror,
	}, nil
}

func buildLLMProvider(cfg *config.AppConfig, resolver *runtime.PathResolver) (closableLLMProvider, error) {
	_ = resolver
	if (cfg.LLMProvider == "openrouter" || cfg.LLMProvider == "ollama") && cfg.OpenRouterAPIKey != "" {
		return gateway.NewProvider(cfg)
	}
	switch cfg.LLMProvider {
	case "anthropic":
		return llm.NewAnthropicProvider(cfg.AnthropicAPIKey, cfg.LLMModel), nil
	case "google":
		return llm.NewGeminiProvider(context.Background(), cfg.GoogleAPIKey, cfg.LLMModel)
	case "kilo":
		return llm.NewKiloProvider(cfg.KiloAPIKey, cfg.LLMModel), nil
	case "ollama":
		return llm.NewOllamaProvider(cfg.LLMModel), nil
	case "openrouter":
		return llm.NewOpenRouterProvider(cfg.OpenRouterAPIKey, cfg.LLMModel), nil
	case "zai":
		return llm.NewZAIProvider(cfg.ZAIAPIKey, cfg.LLMModel), nil
	case "alibaba":
		return llm.NewAlibabaProvider(cfg.AlibabaAPIKey, cfg.LLMModel), nil
	case "openai":
		if cfg.OpenAIAuthMode == "codex" {
			if err := llm.EnsureCodexCLIAvailable(); err != nil {
				return nil, err
			}
			return llm.NewCodexCLIProvider(cfg.LLMModel)
		}
		return llm.NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.LLMModel), nil
	case "", "kimi":
		return llm.NewKimiProvider(cfg.KimiAPIKey, cfg.LLMModel), nil
	default:
		return nil, fmt.Errorf("unsupported llm provider %q", cfg.LLMProvider)
	}
}

func buildTranscriber(cfg *config.AppConfig) (stt.Transcriber, error) {
	switch cfg.STTProvider {
	case "", "groq":
		return stt.NewGroqTranscriber(cfg.GroqAPIKey), nil
	default:
		return nil, fmt.Errorf("unsupported stt provider %q", cfg.STTProvider)
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
		if err := a.heartbeat.Start(); err != nil {
			logger.Warn("failed to start heartbeat service", slog.Any("err", err))
		}
	}

	if a.healthServer != nil {
		if err := a.healthServer.Start(); err != nil {
			logger.Warn("failed to start health server", slog.Any("err", err))
		}
	}
	if a.voiceProcessor != nil {
		if err := a.voiceProcessor.Start(); err != nil {
			logger.Warn("failed to start voice processor", slog.Any("err", err))
		}
	}
	if a.voiceCapture != nil {
		if err := a.voiceCapture.Start(); err != nil {
			logger.Warn("failed to start voice capture worker", slog.Any("err", err))
		}
	}

	dashboard.RegisterRoute("/api/squad", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := json.NewEncoder(w).Encode(agent.GetFixedSquad()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := dashboard.StartServer(logger); err != nil {
		logger.Warn("failed to start ultratrink dashboard", slog.Any("err", err))
	}

	go a.bot.Start()
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
	if a.bot != nil {
		a.bot.Stop()
	}
	_ = ctx
}

type voiceBotDispatcher struct {
	bot *telegram.BotController
}

func (d voiceBotDispatcher) DispatchVoice(ctx context.Context, userID, chatID int64, text string, requiresAudio bool) error {
	if d.bot == nil {
		return fmt.Errorf("telegram bot dispatcher unavailable")
	}
	return d.bot.ProcessExternalInput(ctx, userID, chatID, text, requiresAudio)
}

func (a *app) close() {
	logger := observability.Logger("cmd.app")

	if a.taskStore != nil {
		if err := a.taskStore.Close(); err != nil {
			logger.Warn("failed to close team task store", slog.Any("err", err))
		}
	}
	if a.mcpManager != nil {
		if err := a.mcpManager.Close(); err != nil {
			logger.Warn("failed to close MCP manager", slog.Any("err", err))
		}
	}
	if a.cronStore != nil {
		if err := a.cronStore.Close(); err != nil {
			logger.Warn("failed to close cron store", slog.Any("err", err))
		}
	}
	if a.mem != nil {
		if err := a.mem.Close(); err != nil {
			logger.Warn("failed to close memory manager", slog.Any("err", err))
		}
	}
	if a.voiceMirrorDB != nil {
		if err := a.voiceMirrorDB.Close(); err != nil {
			logger.Warn("failed to close voice sqlite mirror", slog.Any("err", err))
		}
	}
	if a.llmProvider != nil {
		a.llmProvider.Close()
	}
	if a.instanceLock != nil {
		if err := a.instanceLock.Release(); err != nil {
			logger.Warn("failed to release instance lock", slog.Any("err", err))
		}
	}
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

	notifyingRuntime := cron.NewNotifyingRuntime(cronRuntime, func(ctx context.Context, job cron.CronJob, output string, execErr error) error {
		chat := &telebot.Chat{ID: job.TargetChatID}
		if execErr != nil {
			return telegram.SendText(bot.GetBot(), chat, "Falha na rotina agendada:\n"+execErr.Error())
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
