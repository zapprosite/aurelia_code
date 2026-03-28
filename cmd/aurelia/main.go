package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kocar/aurelia/internal/purity/alog"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	alog.Configure(alog.Options{})
	logger := alog.Logger("cmd.main")

	if len(args) > 1 {
		switch args[1] {
		case "onboard":
			if err := runOnboard(os.Stdin, os.Stdout); err != nil {
				logger.Error("failed to run onboarding", slog.Any("err", err))
				return 1
			}
			return 0
		case "auth":
			if len(args) > 2 && args[2] == "openai" {
				if err := runOpenAIAuthLogin(os.Stdin, os.Stdout); err != nil {
					logger.Error("failed to run OpenAI auth login", slog.Any("err", err))
					return 1
				}
				return 0
			}
			logger.Error("unknown auth command", slog.String("command", stringsForLog(args[1:])))
			return 1
		case "voice":
			if err := runVoiceCommand(args[2:], os.Stdout); err != nil {
				logger.Error("failed to run voice command", slog.Any("err", err))
				return 1
			}
			return 0
		default:
			logger.Error("unknown command", slog.String("command", args[1]))
			return 1
		}
	}

	app, err := bootstrapApp(args)
	if err != nil {
		logger.Error("failed to bootstrap Aurelia", slog.Any("err", err))
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

func stringsForLog(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return fmt.Sprint(values)
}
