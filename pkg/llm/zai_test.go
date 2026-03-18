package llm

import "testing"

func TestNewZAIProvider(t *testing.T) {
	t.Parallel()

	provider := NewZAIProvider("secret", "glm-5")
	if provider.baseURL != zAIChatCompletionsURL {
		t.Fatalf("baseURL = %q", provider.baseURL)
	}
	if provider.model != "glm-5" {
		t.Fatalf("model = %q", provider.model)
	}
}
