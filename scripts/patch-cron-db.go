package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "/home/will/.aurelia/data/aurelia.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Update Sentinel Watchdog (5 min)
	watchdogPrompt := "[sys:sentinel-watchdog] Você é o Sentinel. Monitore a saúde crítica." +
		" **Ações**: 1. Execute 'nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader'. 2. Verifique 'docker ps'." +
		" **Análise**: Use o motor local qwen3.5. Se temp > 75°C, execute 'docker stop $(docker ps --format \"{{.ID}}\")' e ALERTE IMEDIATAMENTE." +
		" Se tudo estiver OK, NÃO envie mensagem (STAY SILENT)."

	_, err = db.ExecContext(ctx, "UPDATE cron_jobs SET prompt = ?, cron_expr = ? WHERE prompt LIKE '[sys:sentinel-watchdog]%'", watchdogPrompt, "*/5 * * * *")
	if err != nil {
		log.Fatalf("failed to update watchdog: %v", err)
	}
	fmt.Println("✅ [sys:sentinel-watchdog] atualizado.")

	// Update GPU Report (5 hours)
	reportPrompt := "[sys:gpu-report] Você é o Sentinel, monitor industrial de GPU." +
		" **Ação**: Execute 'node scripts/grafana-capture.mjs'." +
		" **Análise**: Use o motor local qwen3.5 para analisar o dashboard de métricas." +
		" Envie o resumo executivo de 5 horas. Se tudo estiver estável, apenas confirme: '✓ Homelab Industrial Estável'."

	_, err = db.ExecContext(ctx, "UPDATE cron_jobs SET prompt = ?, cron_expr = ? WHERE prompt LIKE '[sys:gpu-report]%'", reportPrompt, "* * * * *")
	if err != nil {
		log.Fatalf("failed to update gpu-report: %v", err)
	}
	fmt.Println("✅ [sys:gpu-report] set to 1 minute for verification.")
}
