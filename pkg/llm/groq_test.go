package llm

import (
	"testing"
)

func TestNewGroqProvider(t *testing.T) {
	provider := NewGroqProvider("test-key", "llama-3.3-70b-versatile")
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestNewGroqProviderWithOptions(t *testing.T) {
	temp := 0.5
	opts := OpenAICompatibleRequestOptions{
		Temperature: &temp,
		MaxTokens:   512,
	}
	provider := NewGroqProviderWithOptions("test-key", "llama-3.3-70b-versatile", opts)
	if provider == nil {
		t.Fatal("expected non-nil provider with options")
	}
}

func TestGroqURL(t *testing.T) {
	if groqChatCompletionsURL != "https://api.groq.com/openai/v1/chat/completions" {
		t.Errorf("unexpected groq URL: %s", groqChatCompletionsURL)
	}
}
