package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListOpenAIModels(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"id":"gpt-5.2"},
				{"id":"o4-mini"},
				{"id":"text-embedding-3-large"}
			]
		}`))
	}))
	defer server.Close()

	models, err := listOpenAIModels(context.Background(), "secret", server.URL, server.Client())
	if err != nil {
		t.Fatalf("listOpenAIModels() error = %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
}

func TestNewOpenAIProvider(t *testing.T) {
	t.Parallel()

	provider := NewOpenAIProvider("secret", "gpt-5.2")
	if provider.baseURL != openAIChatCompletionsURL {
		t.Fatalf("baseURL = %q", provider.baseURL)
	}
}
