package telegram

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/observability"
)

const (
	inputGuardTimeout = 3 * time.Second
)

// Suspicious patterns for secrets extraction
var secretPatterns = []string{
	"api[_-]?key",
	"secret[_-]?key",
	"password",
	"token",
	"credential",
	"auth[_-]?token",
	"access[_-]?token",
	"private[_-]?key",
	".env",
	"config",
	"show[_-]?me[_-]?the[_-]?system[_-]?prompt",
	"what[_-]?is[_-]?your[_-]?system",
	"ignore[_-]?previous",
	"disregard[_-]?instructions",
}

// Destructive command patterns
var destructivePatterns = []string{
	"drop[_-]?database",
	"delete[_-]?all",
	"rm[_-]?rf",
	"destroy",
	"truncate",
	"shutdown",
	"halt",
	"kill[_-]?all",
	"systemctl[_-]?stop",
	"docker[_-]?stop[_-]?all",
}

// InputGuardLight uses keyword matching for fast, lightweight security filtering.
// Only blocks obvious attempts to extract secrets or run destructive commands.
type InputGuardLight struct {
	logger *slog.Logger
}

// NewInputGuardLight creates a lightweight guard that uses keyword matching.
func NewInputGuardLight() *InputGuardLight {
	return &InputGuardLight{
		logger: observability.Logger("telegram.input_guard.light"),
	}
}

// CheckWithUser runs the guard but bypasses it for trusted user IDs.
func (g *InputGuardLight) CheckWithUser(ctx context.Context, userID int64, trustedIDs []int64, text string) (blocked bool, reason string) {
	for _, id := range trustedIDs {
		if id == userID {
			return false, ""
		}
	}
	return g.Check(ctx, text)
}

// Check runs keyword-based filtering. Returns blocked=true and reason if suspicious.
func (g *InputGuardLight) Check(ctx context.Context, text string) (blocked bool, reason string) {
	if strings.TrimSpace(text) == "" {
		return false, ""
	}

	lower := strings.ToLower(text)

	// Check for secret extraction attempts
	for _, pattern := range secretPatterns {
		if contains(lower, pattern) {
			g.logger.Warn("blocked: secret extraction attempt",
				slog.String("pattern", pattern),
				slog.String("text_preview", truncate(text, 80)),
			)
			return true, "Tentativa de extração de secrets detectada"
		}
	}

	// Check for destructive commands
	for _, pattern := range destructivePatterns {
		if contains(lower, pattern) {
			g.logger.Warn("blocked: destructive command attempt",
				slog.String("pattern", pattern),
				slog.String("text_preview", truncate(text, 80)),
			)
			return true, "Comando destrutivo bloqueado"
		}
	}

	return false, ""
}

func contains(s, substr string) bool {
	if strings.Contains(s, substr) {
		return true
	}
	words := strings.Fields(s)
	for _, w := range words {
		if strings.Contains(w, substr) {
			return true
		}
	}
	return false
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
