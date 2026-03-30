package main

import (
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/gateway"
	"github.com/kocar/aurelia/pkg/llm"
)

func TestBuildLLMProvider_Anthropic(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider:     "anthropic",
		LLMModel:        "claude-sonnet-4-6",
		AnthropicAPIKey: "secret",
	}

	provider, err := buildLLMProvider(cfg, nil)
	if err != nil {
		t.Fatalf("buildLLMProvider() error = %v", err)
	}
	defer provider.Close()

	if _, ok := provider.(*llm.AnthropicProvider); !ok {
		t.Fatalf("provider type = %T, want *llm.AnthropicProvider", provider)
	}
}

func TestBuildLLMProvider_Ollama(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider: "ollama",
		LLMModel:    "qwen3.5",
	}

	provider, err := buildLLMProvider(cfg, nil)
	if err != nil {
		t.Fatalf("buildLLMProvider() error = %v", err)
	}
	defer provider.Close()

	if _, ok := provider.(*llm.OpenAICompatibleProvider); !ok {
		t.Fatalf("provider type = %T, want *llm.OpenAICompatibleProvider", provider)
	}
}

func TestBuildLLMProvider_OllamaIgnoresOpenRouterKey(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider:      "ollama",
		LLMModel:         "qwen3.5",
		OpenRouterAPIKey: "secret",
	}

	provider, err := buildLLMProvider(cfg, nil)
	if err != nil {
		t.Fatalf("buildLLMProvider() error = %v", err)
	}
	defer provider.Close()

	if _, ok := provider.(*gateway.Provider); ok {
		t.Fatalf("provider type = %T, expected direct ollama provider without gateway", provider)
	}
	if _, ok := provider.(*llm.OpenAICompatibleProvider); !ok {
		t.Fatalf("provider type = %T, want *llm.OpenAICompatibleProvider", provider)
	}
}

func TestBuildLLMProvider_OpenRouter(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider:      "openrouter",
		LLMModel:         "openrouter/auto",
		OpenRouterAPIKey: "secret",
	}

	provider, err := buildLLMProvider(cfg, nil)
	if err != nil {
		t.Fatalf("buildLLMProvider() error = %v", err)
	}
	defer provider.Close()

	if _, ok := provider.(*gateway.Provider); !ok {
		t.Fatalf("provider type = %T, want *gateway.Provider", provider)
	}
}

func TestBuildLLMProvider_OpenAI(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider:  "openai",
		LLMModel:     "gpt-5.4",
		OpenAIAPIKey: "secret",
	}

	provider, err := buildLLMProvider(cfg, nil)
	if err != nil {
		t.Fatalf("buildLLMProvider() error = %v", err)
	}
	defer provider.Close()

	if _, ok := provider.(*llm.OpenAICompatibleProvider); !ok {
		t.Fatalf("provider type = %T, want *llm.OpenAICompatibleProvider", provider)
	}
}
