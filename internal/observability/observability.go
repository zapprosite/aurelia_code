package observability

import (
	"log/slog"
	"github.com/kocar/aurelia/internal/purity/alog"
)

// Options is a proxy to alog.Options to maintain compatibility.
type Options = alog.Options

// Configure is a proxy to alog.Configure.
func Configure(opts Options) {
	alog.Configure(opts)
}

// Logger is a proxy to alog.Logger.
func Logger(component string) *slog.Logger {
	return alog.Logger(component)
}

// Redact is a proxy to alog.Redact.
func Redact(value string) string {
	return alog.Redact(value)
}

// Basename is a proxy to alog.Basename.
func Basename(path string) string {
	return alog.Basename(path)
}

// MapKeys is a proxy to alog.MapKeys.
func MapKeys(values map[string]any) []string {
	return alog.MapKeys(values)
}
