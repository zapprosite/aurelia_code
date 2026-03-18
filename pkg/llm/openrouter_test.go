package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListOpenRouterModels(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("Authorization = %q", got)
		}
		if got := r.Header.Get("HTTP-Referer"); got != openRouterReferer {
			t.Fatalf("HTTP-Referer = %q", got)
		}
		if got := r.Header.Get("X-Title"); got != openRouterTitle {
			t.Fatalf("X-Title = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"id":"anthropic/claude-sonnet-4","name":"Claude Sonnet 4"},
				{"id":"google/gemini-2.5-flash","name":"Gemini 2.5 Flash"}
			]
		}`))
	}))
	defer server.Close()

	models, err := listOpenRouterModels(context.Background(), "secret", server.URL, server.Client())
	if err != nil {
		t.Fatalf("listOpenRouterModels() error = %v", err)
	}
	if len(models) != 4 {
		t.Fatalf("expected 4 models, got %d", len(models))
	}
	if models[0].ID != "openrouter/auto" || models[1].ID != "openrouter/free" {
		t.Fatalf("unexpected router defaults: %+v", models[:2])
	}
}

func TestNewOpenRouterProvider(t *testing.T) {
	t.Parallel()

	provider := NewOpenRouterProvider("secret", "openrouter/auto")
	if provider.baseURL != openRouterChatCompletionsURL {
		t.Fatalf("baseURL = %q", provider.baseURL)
	}
	if provider.headers["HTTP-Referer"] != openRouterReferer {
		t.Fatalf("HTTP-Referer = %q", provider.headers["HTTP-Referer"])
	}
}
