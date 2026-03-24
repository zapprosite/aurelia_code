package main

import (
	"context"
	"log/slog"
	"strings"

	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/observability"
)

// systemCronDefs define os cron jobs sistêmicos gerenciados pela Aurélia.
// Cada job é identificado por um marcador [sys:nome] no início do prompt.
// Idempotente: jobs com o mesmo marcador não são recriados.
var systemCronDefs = []cron.CronJob{
	{
		ScheduleType: "cron",
		CronExpr:     "*/5 * * * *",
		Prompt: "[sys:sentinel-watchdog] Você é o Sentinel, monitor do homelab." +
			" Execute o watchdog: bash /home/will/aurelia/scripts/homelab-watchdog.sh" +
			" Reporte em até 3 linhas: containers unhealthy, status GPU, Ollama e ZFS." +
			" Só notifique se houver anomalia.",
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
		CronExpr:     "0 9 * * *",
		Prompt: "[sys:gpu-report] Você é o Sentinel, monitor de GPU." +
			" Relatório diário: execute nvidia-smi --query-gpu=name,memory.used,memory.free,utilization.gpu --format=csv,noheader" +
			" e curl -s localhost:11434/api/tags para modelos Ollama carregados." +
			" Alerte se VRAM usada > 80%.",
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
