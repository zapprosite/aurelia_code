package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/kocar/aurelia/internal/observability"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	observability.Configure(observability.Options{})
	logger := observability.Logger("cmd.main")

	if len(args) > 1 {
		switch args[1] {
		case "live":
			if err := runLiveCommand(args[2:], os.Stdout); err != nil {
				logger.Error("failed to run live command", "err", err)
				return 1
			}
			return 0
		default:
			logger.Error("unknown command", "command", args[1])
			return 1
		}
	}

	app, err := bootstrapApp(args)
	if err != nil {
		logger.Error("failed to bootstrap Aurelia", "err", err)
		return 1
	}
	defer app.close()

	app.start()
	waitForShutdownSignal()

	logger.Info("shutting down Aurelia")
	app.shutdown(context.Background())
	return 0
}

func waitForShutdownSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func runLiveCommand(args []string, out io.Writer) error {
	// SOTA 2026: Live command starts a terminal-based coordination session
	// with the sovereign team. Placeholder for now to satisfy build.
	_, err := io.WriteString(out, "Aurelia SOTA 2026.2 - Live Coordination Mode (Headless)\n")
	return err
}
