package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/sentinel"
)

func main() {
	once := flag.Bool("once", false, "executa uma única rodada de manutenção e sai")
	interval := flag.Duration("interval", 15*time.Minute, "intervalo de execução")
	flag.Parse()

	observability.Configure(observability.Options{})
	logger := observability.Logger("cmd.sentinel")

	swarm := sentinel.NewSwarm()
	ctx := context.Background()

	if *once {
		if err := swarm.RunOnce(ctx); err != nil {
			logger.Error("maintenance round failed", slog.Any("err", err))
			os.Exit(1)
		}
		return
	}

	logger.Info("starting continuous maintenance", slog.Duration("interval", *interval))
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	// Initial run
	_ = swarm.RunOnce(ctx)

	for {
		select {
		case <-ticker.C:
			_ = swarm.RunOnce(ctx)
		case <-ctx.Done():
			logger.Info("shutting down sentinel")
			return
		}
	}
}
