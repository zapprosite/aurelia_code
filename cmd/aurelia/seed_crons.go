package main

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/observability"
)

const systemMonitorCronExpr = "*/180 * * * *"

// systemCronDefs define os cron jobs sistêmicos gerenciados pela Aurélia.
// Cada job é identificado por um marcador [sys:nome] no início do prompt.
// Idempotente: jobs com o mesmo marcador não são recriados.
// homelabMonitorCron is the cron definition for homelab health monitoring.
// Registered separately from systemCronDefs to control the marker.
var homelabMonitorCron = cron.CronJob{
	ScheduleType: "cron",
	CronExpr:     systemMonitorCronExpr,
	Prompt: "[sys:homelab-monitor] Use a ferramenta homelab_status para verificar o estado do homelab." +
		" Reporte apenas se houver componente degradado ou offline." +
		" Se todos estiverem saudáveis, NÃO envie mensagem (STAY SILENT).",
}

var controleDBCron = cron.CronJob{
	ScheduleType: "cron",
	CronExpr:     "0 9 * * 1", // toda segunda-feira às 09:00
	Prompt: "[sys:controle-db-audit] Você é o CONTROLE DB. Execute a auditoria semanal da camada de dados:" +
		" 1) Liste todas as collections do Qdrant (curl http://127.0.0.1:6333/collections) e compare com a lista canônica (aurelia_skills, conversation_memory, aurelia_markdown_brain)." +
		" 2) Verifique tamanho do SQLite e integridade (PRAGMA integrity_check)." +
		" 3) Confirme que db_audit_log existe — se não existir, crie." +
		" 4) Detecte tabelas fora do schema canônico." +
		" 5) Cheque WAL size." +
		" 6) Grave o relatório em db_audit_log com action=INVENTORY." +
		" Se tudo estiver OK, NÃO envie mensagem (STAY SILENT). Se houver problema, detalhe risco e ação sugerida.",
}

var repoGuardianCron = cron.CronJob{
	ScheduleType: "cron",
	CronExpr:     "0 */6 * * *",
	Prompt: "[sys:repo-guardian] Você é o guardião do repositório Aurelia." +
		" Execute a skill repo-guardian: audite todos os arquivos .md, verifique links quebrados, ADRs órfãs e arquivos fora do lugar." +
		" Use mcp__ai-context__explore para descoberta semântica e mcp__ai-context__sync para regenerar o codebase-map." +
		" Reporte apenas se houver problemas encontrados ou ações corretivas tomadas." +
		" Se tudo estiver em ordem, NÃO envie mensagem (STAY SILENT).",
}

var systemCronDefs = []cron.CronJob{
	{
		ScheduleType: "cron",
		CronExpr:     systemMonitorCronExpr,
		Prompt: "[sys:sentinel-watchdog] Você é o Sentinel. Monitore a saúde crítica." +
			" **Ações**: 1. Execute 'nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader'. 2. Verifique 'docker ps'." +
			" **Análise**: Use o motor local gemma3:27b-it-qat. Se temp > 75°C, execute 'docker stop $(docker ps --format \"{{.ID}}\")' e ALERTE IMEDIATAMENTE." +
			" Se tudo estiver OK, NÃO envie mensagem (STAY SILENT).",
	},
	{
		ScheduleType: "cron",
		CronExpr:     "0 7 * * *",
		Prompt: "[sys:sentinel-smoke] Você é o Sentinel, monitor do homelab." +
			" Smoke test diário completo: bash /home/will/aurelia/scripts/smoke-test-homelab.sh" +
			" Gere relatório executivo apenas se houver falha, degradação ou drift operacional." +
			" Se tudo estiver saudável, NÃO envie mensagem (STAY SILENT).",
	},
	{
		ScheduleType: "cron",
		CronExpr:     systemMonitorCronExpr,
		Prompt: "[sys:memory-sync] Você é o Sentinel, monitor de memória vetorial." +
			" Verifique Qdrant: curl -s http://localhost:6333/collections/conversation_memory" +
			" Alerte apenas se estiver offline, count=0 ou houver drift relevante." +
			" Se tudo estiver saudável, NÃO envie mensagem (STAY SILENT).",
	},
	{
		ScheduleType: "cron",
		CronExpr:     systemMonitorCronExpr,
		Prompt: "[sys:gpu-report] Você é o Sentinel, monitor industrial de GPU." +
			" **Ação**: Execute 'node scripts/grafana-capture.mjs'." +
			" **Análise**: Use o motor local gemma3:27b-it-qat para analisar o dashboard de métricas." +
			" Envie resumo executivo somente se detectar instabilidade, saturação, erro ou tendência anômala." +
			" Se tudo estiver estável, NÃO envie mensagem (STAY SILENT).",
	},
	{
		ScheduleType: "cron",
		CronExpr:     systemMonitorCronExpr,
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

// seedMarkdownBrainCron registra o cron de sync do cérebro Markdown canônico.
func seedMarkdownBrainCron(ctx context.Context, store *cron.SQLiteCronStore, cfg *config.AppConfig, adminChatID int64) {
	if cfg == nil || adminChatID == 0 {
		return
	}
	logger := observability.Logger("cmd.seed_crons")
	svc := cron.NewService(store, nil)

	existing, err := store.ListJobsByChat(ctx, adminChatID)
	if err != nil {
		logger.Warn("seed_markdown_brain: failed to list jobs", slog.Any("err", err))
		return
	}
	for _, j := range existing {
		if extractSysMarker(j.Prompt) == "[sys:markdown-brain-sync]" {
			return
		}
	}

	job := cron.CronJob{
		ScheduleType: "cron",
		CronExpr:     systemMonitorCronExpr,
		OwnerUserID:  "system",
		TargetChatID: adminChatID,
		Prompt: "[sys:markdown-brain-sync] Use a ferramenta markdown_brain_sync para sincronizar o cérebro Markdown canônico da Aurelia." +
			" Sincronize os `.md` do repositório e, se houver vault configurado, as notas externas também." +
			" Reporte repo_docs, vault_docs, synced_docs, synced_chunks e removed_docs." +
			" Se nenhuma mudança ocorrer, responda: '✓ Markdown Brain sincronizado (sem mudanças)'.",
	}
	if _, err := svc.CreateJob(ctx, job); err != nil {
		logger.Warn("seed_markdown_brain: failed to create cron", slog.Any("err", err))
		return
	}
	logger.Info("seed_markdown_brain: cron criado", slog.String("expr", job.CronExpr))
}

// seedControleDBCron registra a auditoria semanal de dados (idempotente).
func seedControleDBCron(ctx context.Context, store *cron.SQLiteCronStore, adminChatID int64) {
	logger := observability.Logger("cmd.seed_crons")
	svc := cron.NewService(store, nil)

	existing, err := store.ListJobsByChat(ctx, adminChatID)
	if err != nil {
		logger.Warn("seed_controle_db: failed to list jobs", slog.Any("err", err))
		return
	}
	for _, j := range existing {
		if extractSysMarker(j.Prompt) == "[sys:controle-db-audit]" {
			return
		}
	}

	def := controleDBCron
	def.OwnerUserID = "system"
	def.TargetChatID = adminChatID
	if _, err := svc.CreateJob(ctx, def); err != nil {
		logger.Warn("seed_controle_db: failed to create cron", slog.Any("err", err))
		return
	}
	logger.Info("seed_controle_db: cron criado", slog.String("expr", def.CronExpr))
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

func reconcileSystemCronCadence(ctx context.Context, store *cron.SQLiteCronStore, cfg *config.AppConfig) {
	if store == nil || cfg == nil || cfg.VoiceReplyChatID == 0 {
		return
	}

	logger := observability.Logger("cmd.seed_crons")
	jobs, err := store.ListJobsByChat(ctx, cfg.VoiceReplyChatID)
	if err != nil {
		logger.Warn("reconcile_system_crons: failed to list jobs", slog.Any("err", err))
		return
	}

	desired := map[string]string{
		"[sys:homelab-monitor]":     systemMonitorCronExpr,
		"[sys:sentinel-watchdog]":   systemMonitorCronExpr,
		"[sys:memory-sync]":         systemMonitorCronExpr,
		"[sys:gpu-report]":          systemMonitorCronExpr,
		"[sys:aurelia-ping]":        systemMonitorCronExpr,
		"[sys:markdown-brain-sync]": systemMonitorCronExpr,
	}

	now := time.Now().UTC()
	nextRun := now.Truncate(time.Minute).Add(3 * time.Hour)
	updated := 0

	for _, job := range jobs {
		marker := extractSysMarker(job.Prompt)
		wantExpr, ok := desired[marker]
		if !ok || !strings.EqualFold(job.ScheduleType, "cron") {
			continue
		}
		needsUpdate := job.CronExpr != wantExpr || job.NextRunAt == nil || !job.NextRunAt.After(now)
		if !needsUpdate {
			continue
		}
		job.CronExpr = wantExpr
		job.NextRunAt = &nextRun
		if err := store.UpdateJob(ctx, job); err != nil {
			logger.Warn("reconcile_system_crons: failed to update job", slog.String("job_id", job.ID), slog.String("marker", marker), slog.Any("err", err))
			continue
		}
		updated++
		logger.Info("reconcile_system_crons: updated cadence", slog.String("job_id", job.ID), slog.String("marker", marker), slog.String("expr", wantExpr))
	}

	if updated > 0 {
		logger.Info("reconcile_system_crons: done", slog.Int("updated", updated))
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
