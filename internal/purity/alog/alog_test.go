package alog_test

import (
	"log/slog"
	"testing"

	"github.com/kocar/aurelia/internal/purity/alog"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	// Padrão SOTA 2026.1: Testes isolados e declarativos.
	t.Run("default configuration", func(t *testing.T) {
		alog.Configure(alog.Options{
			Level: slog.LevelInfo,
		})
		// Apenas validação de boot, o slog global é mutado.
		alog.Info("test log info")
		assert.NotNil(t, alog.Logger("test"))
	})

	t.Run("redact sensitive values", func(t *testing.T) {
		assert.Equal(t, "[redacted]", alog.Redact("secret-token"))
		assert.Equal(t, "", alog.Redact(""))
	})

	t.Run("basename path utility", func(t *testing.T) {
		assert.Equal(t, "file.go", alog.Basename("/home/will/aurelia/file.go"))
	})
}
