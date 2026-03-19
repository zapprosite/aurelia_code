package observability

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Options controls the default logger configuration for the process.
type Options struct {
	Format string
	Level  slog.Level
}

// Configure installs the process-wide slog logger and bridges the stdlib log package into it.
func Configure(opts Options) {
	level := opts.Level
	if envLevel := strings.TrimSpace(os.Getenv("AURELIA_LOG_LEVEL")); envLevel != "" {
		switch strings.ToLower(envLevel) {
		case "debug":
			level = slog.LevelDebug
		case "warn", "warning":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
	} else if level == 0 {
		level = slog.LevelInfo
	}

	format := strings.ToLower(strings.TrimSpace(opts.Format))
	if format == "" {
		format = strings.ToLower(strings.TrimSpace(os.Getenv("AURELIA_LOG_FORMAT")))
	}
	if format == "" {
		format = "text"
	}

	handlerOpts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	default:
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	log.SetFlags(0)
	log.SetOutput(stdLogBridge{logger: logger.With(slog.String("component", "stdlog"))})
}

// Logger returns a component-scoped logger.
func Logger(component string) *slog.Logger {
	return slog.Default().With(slog.String("component", component))
}

// Redact returns a generic placeholder for sensitive values.
func Redact(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return "[redacted]"
}

// RedactToolArgs returns a list of keys from a tool argument map, hiding the values.
func RedactToolArgs(args map[string]any) []string {
	if args == nil {
		return nil
	}
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Basename removes host-specific directory details from a path before logging it.
func Basename(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	return filepath.Base(path)
}

// MapKeys returns a stable, sorted list of map keys for safe diagnostics.
func MapKeys(values map[string]any) []string {
	if len(values) == 0 {
		return nil
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

type stdLogBridge struct {
	logger *slog.Logger
}

func (b stdLogBridge) Write(p []byte) (int, error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		b.logger.Info(msg)
	}
	return len(p), nil
}
