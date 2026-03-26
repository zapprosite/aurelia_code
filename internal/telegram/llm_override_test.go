package telegram

import (
	"testing"

	"github.com/kocar/aurelia/internal/config"
)

func TestBuildBotOverrideProvider_NoneConfigured(t *testing.T) {
	provider, effectiveProvider, effectiveModel, err := buildBotOverrideProvider(&config.AppConfig{}, config.BotConfig{ID: "aurelia"})
	if err != nil {
		t.Fatalf("buildBotOverrideProvider() error = %v", err)
	}
	if provider != nil || effectiveProvider != "" || effectiveModel != "" {
		t.Fatalf("expected no override, got provider=%v effectiveProvider=%q effectiveModel=%q", provider, effectiveProvider, effectiveModel)
	}
}

func TestBuildBotOverrideProvider_ControleDBPinnedToMiniMaxM27(t *testing.T) {
	provider, effectiveProvider, effectiveModel, err := buildBotOverrideProvider(&config.AppConfig{
		OpenRouterAPIKey: "test-key",
	}, config.BotConfig{
		ID:          "controle-db",
		LLMProvider: "groq",
		LLMModel:    "llama-3.3-70b-versatile",
	})
	if err != nil {
		t.Fatalf("buildBotOverrideProvider() error = %v", err)
	}
	if provider == nil {
		t.Fatal("expected provider")
	}
	if effectiveProvider != "openrouter" {
		t.Fatalf("effectiveProvider = %q, want openrouter", effectiveProvider)
	}
	if effectiveModel != "minimax/minimax-m2.7" {
		t.Fatalf("effectiveModel = %q, want minimax/minimax-m2.7", effectiveModel)
	}
}
