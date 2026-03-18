package main

import (
	"reflect"
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/pkg/llm"
)

func TestBuildLLMProvider_UsesConfiguredModel(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider: "kimi",
		LLMModel:    "moonshot-v1-32k",
		KimiAPIKey:  "secret",
	}

	provider, err := buildLLMProvider(cfg, nil)
	if err != nil {
		t.Fatalf("buildLLMProvider() error = %v", err)
	}
	defer provider.Close()

	kimiProvider, ok := provider.(*llm.KimiProvider)
	if !ok {
		t.Fatalf("provider type = %T, want *llm.KimiProvider", provider)
	}
	model := reflect.ValueOf(kimiProvider).Elem().FieldByName("model").String()
	if model != "moonshot-v1-32k" {
		t.Fatalf("model = %q", model)
	}
}

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

func TestBuildLLMProvider_Google(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider:  "google",
		LLMModel:     "gemini-2.5-pro",
		GoogleAPIKey: "secret",
	}

	provider, err := buildLLMProvider(cfg, nil)
	if err != nil {
		t.Fatalf("buildLLMProvider() error = %v", err)
	}
	defer provider.Close()

	if _, ok := provider.(*llm.GeminiProvider); !ok {
		t.Fatalf("provider type = %T, want *llm.GeminiProvider", provider)
	}
}

func TestBuildLLMProvider_Kilo(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider: "kilo",
		LLMModel:    "gpt-5.4",
		KiloAPIKey:  "secret",
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

	if _, ok := provider.(*llm.OpenAICompatibleProvider); !ok {
		t.Fatalf("provider type = %T, want *llm.OpenAICompatibleProvider", provider)
	}
}

func TestBuildLLMProvider_ZAI(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider: "zai",
		LLMModel:    "glm-5",
		ZAIAPIKey:   "secret",
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

func TestBuildLLMProvider_Alibaba(t *testing.T) {
	cfg := &config.AppConfig{
		LLMProvider:   "alibaba",
		LLMModel:      "qwen3-coder-plus",
		AlibabaAPIKey: "secret",
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

func TestBuildLLMProvider_OpenAICodex(t *testing.T) {
	originalLookPath := llm.CodexLookPathForTest(func(string) (string, error) {
		return "codex", nil
	})
	defer originalLookPath()
	restoreFactory := llm.UseNoopCodexCallerForTest()
	defer restoreFactory()

	cfg := &config.AppConfig{
		LLMProvider:    "openai",
		LLMModel:       "gpt-5.2-codex",
		OpenAIAuthMode: "codex",
	}

	provider, err := buildLLMProvider(cfg, nil)
	if err != nil {
		t.Fatalf("buildLLMProvider() error = %v", err)
	}
	defer provider.Close()

	if _, ok := provider.(*llm.CodexCLIProvider); !ok {
		t.Fatalf("provider type = %T, want *llm.CodexCLIProvider", provider)
	}
}
