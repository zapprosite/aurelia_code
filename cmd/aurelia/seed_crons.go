package main

import (
	"context"
	"log/slog"
	"strings"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/observability"
)

// systemCronDefs define os cron jobs sistêmicos gerenciados pela Aurélia.
// Cada job é identificado por um marcador [sys:nome] no início do prompt.
// Idempotente: jobs com o mesmo marcador não são recriados.
// homelabMonitorCron is the cron definition for homelab health monitoring.
// Registered separately from systemCronDefs to control the marker.
var homelabMonitorCron = cron.CronJob{
	ScheduleType: "cron",
	CronExpr:     "*/15 * * * *",
	Prompt: "[sys:homelab-monitor] Use a ferramenta homelab_status para verificar o estado do homelab." +
		" Reporte apenas se houver componente degradado ou offline." +
		" Se todos estiverem saudáveis, NÃO envie mensagem (STAY SILENT).",
}

var repoGuardianCron = cron.CronJob{
	ScheduleType: "cron",
	CronExpr:     "0 */6 * * *",
	Prompt: "[sys:repo-guardian] Você é o guardião do repositório Aurelia." +
		" Execute a skill repo-guardian: audite todos os arquivos .md, verifique links quebrados, ADRs órfãs e arquivos fora do lugar." +
		" Use mcp__ai-context__explore para descoberta semântica e mcp__ai-context__sync para regenerar o codebase-map." +
		" Reporte: contagem de .md, problemas encontrados, ações tomadas." +
		" Se tudo estiver em ordem, responda: '✓ Repositório organizado (sem anomalias detectadas)'.",
}

var systemCronDefs = []cron.CronJob{
	{
		ScheduleType: "cron",
		CronExpr:     "*/5 * * * *",
		Prompt: "[sys:sentinel-watchdog] Você é o Sentinel. Monitore a saúde crítica." +
			" **Ações**: 1. Execute 'nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader'. 2. Verifique 'docker ps'." +
			" **Análise**: Use o motor local gemma3:12b. Se temp > 75°C, execute 'docker stop $(docker ps --format \"{{.ID}}\")' e ALERTE IMEDIATAMENTE." +
			" Se tudo estiver OK, NÃO envie mensagem (STAY SILENT).",
	},
	{
		ScheduleType: "cron",
		CronExpr:     "0 7 * * *",
		Prompt: "[sys:sentinel-smoke] Você é o Sentinel, monitor do homelab." +
			" Smoke test diário completo: bash /home/will/aurelia/scripts/smoke-test-homelab.sh" +
			" Gere relatório executivo com status de containers, VRAM livre, ZFS e tunnel Cloudflared.",
	},
	{
		ScheduleType: "cron",
		CronExpr:     "*/30 * * * *",
		Prompt: "[sys:memory-sync] Você é o Sentinel, monitor de memória vetorial." +
			" Verifique Qdrant: curl -s http://localhost:6333/collections/conversation_memory" +
			" Reporte status e contagem de vetores. Alerte se offline ou count=0.",
	},
	{
		ScheduleType: "cron",
		CronExpr:     "0 */5 * * *",
		Prompt: "[sys:gpu-report] Você é o Sentinel, monitor industrial de GPU." +
			" **Ação**: Execute 'node scripts/grafana-capture.mjs'." +
			" **Análise**: Use o motor local gemma3:12b para analisar o dashboard de métricas." +
			" Envie o resumo executivo de 5 horas. Se tudo estiver estável, apenas confirme: '✓ Homelab Industrial Estável'.",
	},
	{
		ScheduleType: "cron",
		CronExpr:     "*/15 * * * *",
		Prompt: "[sys:aurelia-ping] Você é o Sentinel." +
			" Health check do daemon Aurelia: curl -s http://localhost:8484/health" +
			" Responda em 1 linha com: status, LLM provider ativo e voz habilitada." +
			" Só notifique se status != ok.",
	},
}

// seedHomelabMonitorCron registra o cron de monitoramento do homelab (idempotente).
func seedHomelabMonitorCron(ctx context.Context, store *cron.SQLiteCronStore, adminChatID int64) {
	logger := observability.Logger("cmd.seed_crons")
	svc := cron.NewService(store, nil)

	existing, err := store.ListJobsByChat(ctx, adminChatID)
	if err != nil {
		logger.Warn("seed_homelab_monitor: failed to list jobs", slog.Any("err", err))
		return
	}
	for _, j := range existing {
		if extractSysMarker(j.Prompt) == "[sys:homelab-monitor]" {
			return
		}
	}

	def := homelabMonitorCron
	def.OwnerUserID = "system"
	def.TargetChatID = adminChatID
	if _, err := svc.CreateJob(ctx, def); err != nil {
		logger.Warn("seed_homelab_monitor: failed to create cron", slog.Any("err", err))
		return
	}
	logger.Info("seed_homelab_monitor: cron criado", slog.String("expr", def.CronExpr))
}

// seedObsidianCron registra o cron de sync do Obsidian se habilitado.
func seedObsidianCron(ctx context.Context, store *cron.SQLiteCronStore, cfg *config.AppConfig, adminChatID int64) {
	if !cfg.ObsidianSyncEnabled || cfg.ObsidianVaultPath == "" {
		return
	}
	logger := observability.Logger("cmd.seed_crons")
	svc := cron.NewService(store, nil)

	existing, err := store.ListJobsByChat(ctx, adminChatID)
	if err != nil {
		logger.Warn("seed_obsidian: failed to list jobs", slog.Any("err", err))
		return
	}
	for _, j := range existing {
		if extractSysMarker(j.Prompt) == "[sys:obsidian-sync]" {
			return
		}
	}

	job := cron.CronJob{
		ScheduleType: "cron",
		CronExpr:     "*/30 * * * *",
		OwnerUserID:  "system",
		TargetChatID: adminChatID,
		Prompt: "[sys:obsidian-sync] Use a ferramenta obsidian_sync para sincronizar o vault do Obsidian com o Qdrant." +
			" Reporte quantas notas foram indexadas. Se nenhuma mudou, responda: '✓ Obsidian vault sincronizado (sem mudanças)'.",
	}
	if _, err := svc.CreateJob(ctx, job); err != nil {
		logger.Warn("seed_obsidian: failed to create cron", slog.Any("err", err))
		return
	}
	logger.Info("seed_obsidian: cron criado", slog.String("vault", cfg.ObsidianVaultPath))
}

// seedRepoGuardianCron registra o cron de governança do repositório (idempotente).
func seedRepoGuardianCron(ctx context.Context, store *cron.SQLiteCronStore, adminChatID int64) {
	logger := observability.Logger("cmd.seed_crons")
	svc := cron.NewService(store, nil)

	existing, err := store.ListJobsByChat(ctx, adminChatID)
	if err != nil {
		logger.Warn("seed_repo_guardian: failed to list jobs", slog.Any("err", err))
		return
	}
	for _, j := range existing {
		if extractSysMarker(j.Prompt) == "[sys:repo-guardian]" {
			return
		}
	}

	def := repoGuardianCron
	def.OwnerUserID = "system"
	def.TargetChatID = adminChatID
	if _, err := svc.CreateJob(ctx, def); err != nil {
		logger.Warn("seed_repo_guardian: failed to create cron", slog.Any("err", err))
		return
	}
	logger.Info("seed_repo_guardian: cron criado", slog.String("expr", def.CronExpr))
}

// seedSystemCrons cria os cron jobs sistêmicos caso ainda não existam.
// É idempotente: identifica jobs pelo marcador [sys:nome] no prompt.
func seedSystemCrons(ctx context.Context, store *cron.SQLiteCronStore, adminChatID int64) {
	logger := observability.Logger("cmd.seed_crons")
	svc := cron.NewService(store, nil)

	existing, err := store.ListJobsByChat(ctx, adminChatID)
	if err != nil {
		logger.Warn("seed_crons: failed to list existing jobs", slog.Any("err", err))
		return
	}

	presentMarkers := map[string]bool{}
	for _, j := range existing {
		if m := extractSysMarker(j.Prompt); m != "" {
			presentMarkers[m] = true
		}
	}

	created := 0
	for _, def := range systemCronDefs {
		marker := extractSysMarker(def.Prompt)
		if presentMarkers[marker] {
			continue
		}
		def.OwnerUserID = "system"
		def.TargetChatID = adminChatID
		if _, err := svc.CreateJob(ctx, def); err != nil {
			logger.Warn("seed_crons: failed to create job", slog.String("marker", marker), slog.Any("err", err))
			continue
		}
		logger.Info("seed_crons: created system cron", slog.String("marker", marker), slog.String("expr", def.CronExpr))
		created++
	}

	if created > 0 {
		logger.Info("seed_crons: done", slog.Int("created", created))
	}
}

// extractSysMarker retorna o marcador [sys:nome] do início do prompt, ou "".
func extractSysMarker(prompt string) string {
	if !strings.HasPrefix(prompt, "[sys:") {
		return ""
	}
	end := strings.Index(prompt, "]")
	if end < 0 {
		return ""
	}
	return prompt[:end+1]
}
