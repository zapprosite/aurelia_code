package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

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
		default:
			logger.Error("unknown command", slog.String("command", args[1]))
			return 1
		}
	}

	app, err := bootstrapApp(args)
	if err != nil {
		return exitCodeForBootstrapError(logger, args, err)
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

func exitCodeForBootstrapError(logger *slog.Logger, args []string, err error) int {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	if shouldSuppressDuplicateLaunch(os.Getppid(), err) {
		recordDuplicateLaunch(args, err)
		return 0
	}

	logger.Error("failed to bootstrap Aurelia", bootstrapErrorAttrs(args, err)...)
	return 1
}

func bootstrapErrorAttrs(args []string, err error) []any {
	return []any{
		slog.Any("err", err),
		slog.Int("pid", os.Getpid()),
		slog.Int("ppid", os.Getppid()),
		slog.String("argv", stringsForLog(args)),
		slog.String("parent_cmd", parentCommandForLog(os.Getppid())),
		slog.String("cwd", cwdForLog()),
		slog.String("invocation_id", os.Getenv("INVOCATION_ID")),
	}
}

func isAlreadyRunningInstanceError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "acquire instance lock: another Aurelia instance is already running")
}

func shouldSuppressDuplicateLaunch(ppid int, err error) bool {
	return isAlreadyRunningInstanceError(err) && ppid == 1
}

func recordDuplicateLaunch(args []string, err error) {
	logDir := filepath.Join(aureliaHomeForLog(), "logs")
	if mkErr := os.MkdirAll(logDir, 0o755); mkErr != nil {
		return
	}

	file, openErr := os.OpenFile(filepath.Join(logDir, "duplicate-launch.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if openErr != nil {
		return
	}
	defer func() { _ = file.Close() }()

	_, _ = fmt.Fprintf(
		file,
		"%s pid=%d ppid=%d argv=%s cwd=%q parent_cmd=%q env=%q err=%q\n",
		time.Now().Format(time.RFC3339Nano),
		os.Getpid(),
		os.Getppid(),
		stringsForLog(args),
		cwdForLog(),
		parentCommandForLog(os.Getppid()),
		selectedEnvForLog(),
		err.Error(),
	)
}

func aureliaHomeForLog() string {
	if root := strings.TrimSpace(os.Getenv("AURELIA_HOME")); root != "" {
		return root
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return ".aurelia"
	}
	return filepath.Join(home, ".aurelia")
}

func cwdForLog() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return cwd
}

func selectedEnvForLog() string {
	keys := []string{
		"PWD",
		"SHLVL",
		"TERM",
		"TERM_PROGRAM",
		"TERM_PROGRAM_VERSION",
		"CODEX_HOME",
		"CLAUDE_CWD",
		"CLAUDE_PROJECT_DIR",
		"INVOCATION_ID",
	}
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			continue
		}
		values = append(values, key+"="+value)
	}
	return strings.Join(values, ", ")
}

func parentCommandForLog(ppid int) string {
	if ppid <= 1 {
		return ""
	}

	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(ppid), "cmdline"))
	if err != nil {
		return ""
	}

	parts := strings.Split(string(data), "\x00")
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	return stringsForLog(filtered)
}
