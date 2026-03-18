package llm

import "testing"

func TestNewAlibabaProvider(t *testing.T) {
	t.Parallel()

	provider := NewAlibabaProvider("secret", "qwen-plus")
	if provider.baseURL != alibabaChatCompletionsURL {
		t.Fatalf("baseURL = %q", provider.baseURL)
	}
	if provider.model != "qwen-plus" {
		t.Fatalf("model = %q", provider.model)
	}
}
