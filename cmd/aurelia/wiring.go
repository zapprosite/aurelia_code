package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/computer_use"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/mcp"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/telegram"
	"github.com/kocar/aurelia/internal/tools"
	"gopkg.in/telebot.v3"
)

// ── Tool registry ─────────────────────────────────────────────────────────────

func buildToolRegistry() *agent.ToolRegistry {
	registry := agent.NewToolRegistry()
	tools.RegisterCoreTools(registry)
	registry.RegisterPlannerTools()
	registry.RegisterMemoryTools()
	registry.RegisterVerifierTools()

	// Enterprise Kit (SOTA 2026)
	ek := computer_use.NewEnterpriseKit("https://google.com")
	registry.Register(agent.Tool{
		Name:        "playwright_codegen",
		Description: "Abre o gravador do Playwright para gerar scripts de automação. Reclama um URL.",
		JSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{"type": "string"},
			},
		},
	}, func(ctx context.Context, args map[string]any) (string, error) {
		url, _ := args["url"].(string)
		return ek.PlaywrightCodeGen(url)
	})

	registry.Register(agent.Tool{
		Name:        "inspect_dom",
		Description: "Extrai a árvore de acessibilidade/DOM de um URL para análise profunda.",
		JSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{"type": "string"},
			},
		},
	}, func(ctx context.Context, args map[string]any) (string, error) {
		url, _ := args["url"].(string)
		return ek.InspectDOM(url)
	})

	return registry
}

func registerHomelabTool(cfg *config.AppConfig, registry *agent.ToolRegistry) {
	tools.SetSearXNGURL(cfg.SearXNGURL)
	tools.RegisterHomelabTool(registry, cfg.OllamaURL, cfg.QdrantURL)
}

func maybeRegisterMarkdownBrainTool(cfg *config.AppConfig, registry *agent.ToolRegistry, repoRoot string, db *sql.DB) *tools.MarkdownBrainSyncTool {
	if cfg == nil {
		return nil
	}
	vaultPath := ""
	if cfg.ObsidianSyncEnabled {
		vaultPath = cfg.ObsidianVaultPath
	}
	tool := tools.NewMarkdownBrainSyncTool(
		repoRoot,
		vaultPath,
		cfg.OllamaURL,
		cfg.QdrantEmbeddingModel,
		cfg.QdrantURL,
		cfg.QdrantAPIKey,
		"",
		db,
		observability.Logger("markdown_brain"),
	)
	if tool == nil {
		return nil
	}
	registry.Register(tool.Definition(), tool.Execute)
	logger := observability.Logger("cmd.wiring")
	logger.Info("markdown_brain_sync tool registered", slog.String("repo_root", repoRoot))
	if vaultPath != "" {
		logger.Info("markdown brain vault source enabled", slog.String("vault", vaultPath))
	}
	return tool
}

func registerScheduleTools(registry *agent.ToolRegistry, cronStore *cron.SQLiteCronStore) *cron.Service {
	cronService := cron.NewService(cronStore, nil)
	registry.Register(tools.NewCreateScheduleTool(cronService).Definition(), tools.NewCreateScheduleTool(cronService).Execute)
	registry.Register(tools.NewListSchedulesTool(cronService).Definition(), tools.NewListSchedulesTool(cronService).Execute)
	registry.Register(tools.NewPauseScheduleTool(cronService).Definition(), tools.NewPauseScheduleTool(cronService).Execute)
	registry.Register(tools.NewResumeScheduleTool(cronService).Definition(), tools.NewResumeScheduleTool(cronService).Execute)
	registry.Register(tools.NewDeleteScheduleTool(cronService).Definition(), tools.NewDeleteScheduleTool(cronService).Execute)
	return cronService
}

func maybeRegisterMCPTools(cfg *config.AppConfig, registry *agent.ToolRegistry) (*mcp.Manager, error) {
	mcpCfg, err := config.LoadMCPConfig(cfg.MCPConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, loggableError("failed to load MCP config from %s: %v", cfg.MCPConfigPath, err)
	}
	if mcpCfg == nil || !mcpCfg.Enabled {
		return nil, nil
	}
	cwd, _ := os.Getwd()
	mcpManager, err := mcp.NewManager(*mcpCfg, cwd)
	if err != nil {
		return nil, loggableError("failed to start MCP Manager: %v", err)
	}
	tools.RegisterMCPTools(registry, mcpManager)
	// Inject MCP manager for computer use tools (stagehand handlers)
	tools.SetGlobalMCPManager(mcpManager)
	observability.Logger("cmd.wiring").Info("MCP manager initialized", slog.Int("tool_count", len(mcpManager.ToolSpecs())))
	return mcpManager, nil
}

func registerSpawnAgentTool(
	cfg *config.AppConfig,
	registry *agent.ToolRegistry,
	llmProvider agent.LLMProvider,
	bot *telegram.BotController,
	taskStore *agent.SQLiteTaskStore,
) error {
	teamManager, err := agent.NewTeamManager(taskStore)
	if err != nil {
		return loggableError("initialize team manager: %v", err)
	}
	masterTeams, err := agent.NewMasterTeamService(
		teamManager,
		llmProvider,
		registry,
		cfg.MaxIterations,
		func(teamKey string, message string) { notifyMasterTeam(bot, teamKey, message) },
	)
	if err != nil {
		return loggableError("initialize master team service: %v", err)
	}
	if err := masterTeams.Rehydrate(context.Background()); err != nil {
		observability.Logger("cmd.wiring").Warn("failed to rehydrate master teams", slog.Any("err", err))
	}

	registry.Register(tools.NewSpawnAgentTool(masterTeams).Definition(), tools.NewSpawnAgentTool(masterTeams).Execute)
	registry.Register(tools.NewCreateSquadTool(masterTeams).Definition(), tools.NewCreateSquadTool(masterTeams).Execute)
	registry.Register(tools.NewPauseTeamTool(masterTeams).Definition(), tools.NewPauseTeamTool(masterTeams).Execute)
	registry.Register(tools.NewResumeTeamTool(masterTeams).Definition(), tools.NewResumeTeamTool(masterTeams).Execute)
	registry.Register(tools.NewCancelTeamTool(masterTeams).Definition(), tools.NewCancelTeamTool(masterTeams).Execute)
	registry.Register(tools.NewTeamStatusTool(masterTeams).Definition(), tools.NewTeamStatusTool(masterTeams).Execute)
	registry.Register(tools.NewSendTeamMessageTool(teamManager).Definition(), tools.NewSendTeamMessageTool(teamManager).Execute)
	registry.Register(tools.NewReadTeamInboxTool(teamManager).Definition(), tools.NewReadTeamInboxTool(teamManager).Execute)
	return nil
}

func notifyMasterTeam(bot *telegram.BotController, teamKey, message string) {
	logger := observability.Logger("cmd.wiring")
	chatID, err := strconv.ParseInt(teamKey, 10, 64)
	if err != nil {
		logger.Warn("invalid team key for master notification", slog.String("team_key", teamKey), slog.Any("err", err))
		return
	}
	if err := telegram.SendText(bot.GetBot(), &telebot.Chat{ID: chatID}, message); err != nil {
		logger.Warn("failed to send master team update", slog.Int64("chat_id", chatID), slog.Any("err", err))
	}
}

func loggableError(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// ── Auth (legacy stub) ────────────────────────────────────────────────────────

func runOpenAIAuthLogin(_ io.Reader, _ io.Writer) error {
	return fmt.Errorf("this CLI auth has been removed; use the openai provider with an API key instead")
}

// ── Cron support ──────────────────────────────────────────────────────────────

type loopExecutorAdapter struct {
	loop *agent.Loop
}

func (a *loopExecutorAdapter) Execute(ctx context.Context, systemPrompt string, history []agent.Message, allowedTools []string) ([]agent.Message, string, error) {
	return a.loop.Run(ctx, systemPrompt, history, allowedTools)
}

func (a *loopExecutorAdapter) RunCommand(ctx context.Context, command string) (string, error) {
	return a.loop.Registry().Execute(ctx, "run_command", map[string]any{"command": command})
}

func loadCronPromptConfig(canonicalPromptLoader interface {
	BuildPrompt(ctx context.Context, userID, conversationID string) (string, []string, error)
}) (string, []string) {
	if canonicalPromptLoader != nil {
		prompt, tools, err := canonicalPromptLoader.BuildPrompt(context.Background(), "", "")
		if err == nil {
			return prompt, tools
		}
		log.Printf("Warning: failed to load canonical prompt for cron runtime, using default prompt: %v", err)
	}
	return "Voce e o agente pessoal Aurelia. Execute a tarefa agendada com precisao e retorne um resumo objetivo do resultado.", []string{"web_search", "read_file", "write_file", "list_dir", "run_command"}
}

// ── LLM runtime snapshot ──────────────────────────────────────────────────────

type llmRuntimeSnapshot struct {
	RequestedProvider string    `json:"requested_provider"`
	RequestedModel    string    `json:"requested_model"`
	EffectiveProvider string    `json:"effective_provider"`
	EffectiveModel    string    `json:"effective_model"`
	CheckedAt         time.Time `json:"checked_at"`
}

func buildLLMRuntimeSnapshot(a *app, checkedAt time.Time) llmRuntimeSnapshot {
	snapshot := llmRuntimeSnapshot{
		RequestedProvider: "unconfigured",
		EffectiveProvider: "unconfigured",
		CheckedAt:         checkedAt,
	}
	if a == nil || a.cfg == nil {
		return snapshot
	}
	snapshot.RequestedProvider = a.cfg.LLMProvider
	snapshot.RequestedModel = a.cfg.LLMModel
	snapshot.EffectiveProvider = a.cfg.LLMProvider
	snapshot.EffectiveModel = a.cfg.LLMModel

	return snapshot
}

// ── Status adapters ───────────────────────────────────────────────────────────

type squadStatusAdapter struct{}

func (squadStatusAdapter) GetSquadStatus() []telegram.AgentStatus {
	members := agent.GetFixedSquad()
	out := make([]telegram.AgentStatus, 0, len(members))
	for _, m := range members {
		out = append(out, telegram.AgentStatus{
			Name:   m.Name,
			Icon:   m.IconName,
			Role:   m.Role,
			Status: m.Status,
			Load:   m.Load,
		})
	}
	return out
}

type cronNextJobAdapter struct {
	store *cron.SQLiteCronStore
}

func (a *cronNextJobAdapter) GetNextJobs(ctx context.Context, limit int) []telegram.NextJob {
	if a.store == nil {
		return nil
	}
	now := time.Now().UTC()
	jobs, err := a.store.ListDueJobs(ctx, now.Add(24*time.Hour), limit*5)
	if err != nil || len(jobs) == 0 {
		return nil
	}
	var result []telegram.NextJob
	seen := 0
	for _, j := range jobs {
		if !j.Active || j.NextRunAt == nil {
			continue
		}
		dur := j.NextRunAt.Sub(now)
		if dur < 0 {
			dur = 0
		}
		name := j.ID
		if len(name) > 16 {
			name = name[:16]
		}
		result = append(result, telegram.NextJob{Name: name, NextIn: formatDuration(dur)})
		seen++
		if seen >= limit {
			break
		}
	}
	return result
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dmin", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh%02dmin", int(d.Hours()), int(d.Minutes())%60)
}
