package agent

import (
	"context"
	"log/slog"
	"os"

	"github.com/kocar/aurelia/internal/purity/alog"
)

// SetupTest centraliza a configuração do ambiente de teste SOTA 2026.1.
func SetupTest() {
	// Força o logger estruturado para diagnósticos em caso de falha.
	alog.Configure(alog.Options{
		Level: slog.LevelDebug,
	})
	os.Setenv("AURELIA_LOG_FORMAT", "text")
}

type MockLLMProvider struct {
	response *ModelResponse
	err      error
}

func (m *MockLLMProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	return m.response, m.err
}
