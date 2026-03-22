package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestFallbackModels_UnknownProviderReturnsEmpty(t *testing.T) {
	t.Parallel()

	models := FallbackModels("unknown")
	if len(models) != 0 {
		t.Fatalf("expected empty fallback catalog, got %v", models)
	}
}

func TestModelOptionLabel(t *testing.T) {
	t.Parallel()

	option := ModelOption{ID: "qwen3.5:9b", Name: "Qwen 3.5 9B"}
	if got := option.Label(); got != "Qwen 3.5 9B (qwen3.5:9b)" {
		t.Fatalf("Label() = %q", got)
	}

	option = ModelOption{ID: "gpt-5.4", Name: "GPT-5.4", SupportsImageInput: true, SupportsTools: true, IsFree: true}
	if got := option.Label(); got != "GPT-5.4 (gpt-5.4) [vision, tools, free]" {
		t.Fatalf("Label() with badges = %q", got)
	}
}

func TestListAnthropicModels(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("x-api-key"); got != "secret" {
			t.Fatalf("x-api-key = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"id":"claude-sonnet-4-6","display_name":"Claude Sonnet 4.6","created_at":"2026-01-01T00:00:00Z","type":"model"},
				{"id":"claude-haiku-4-5","display_name":"Claude Haiku 4.5","created_at":"2025-10-01T00:00:00Z","type":"model"}
			],
			"has_more": false,
			"first_id": "claude-sonnet-4-6",
			"last_id": "claude-haiku-4-5"
		}`))
	}))
	defer server.Close()

	models, err := listAnthropicModels(context.Background(), "secret", option.WithBaseURL(server.URL+"/"))
	if err != nil {
		t.Fatalf("listAnthropicModels() error = %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
	if models[0].ID != "claude-sonnet-4-6" {
		t.Fatalf("first model = %+v", models[0])
	}
}

func TestListGoogleModels(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("key"); got != "secret" {
			t.Fatalf("key = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"models": [
				{
					"name": "models/gemini-2.5-flash",
					"displayName": "Gemini 2.5 Flash",
					"supportedGenerationMethods": ["generateContent"]
				},
				{
					"name": "models/text-embedding-004",
					"displayName": "Text Embedding 004",
					"supportedGenerationMethods": ["embedContent"]
				}
			]
		}`))
	}))
	defer server.Close()

	models, err := listGoogleModels(context.Background(), "secret", server.URL, server.Client())
	if err != nil {
		t.Fatalf("listGoogleModels() error = %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}
	if models[0].ID != "gemini-2.5-flash" {
		t.Fatalf("first model = %+v", models[0])
	}
}

func TestListOllamaModels(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"id":"bge-m3:latest"},
				{"id":"qwen3.5:9b"},
				{"id":"qwen3.5:27b-q4_K_M"},
				{"id":"gemma3:27b-it-q4_K_M"}
			]
		}`))
	}))
	defer server.Close()

	models, err := listOllamaModels(context.Background(), server.URL+"/v1/models", server.Client())
	if err != nil {
		t.Fatalf("listOllamaModels() error = %v", err)
	}
	if len(models) != 3 {
		t.Fatalf("expected 3 chat models, got %d", len(models))
	}
	got := []string{models[0].ID, models[1].ID, models[2].ID}
	if !containsModelID(got, "qwen3.5:9b") {
		t.Fatalf("models missing qwen3.5:9b: %v", got)
	}
	if !containsModelID(got, "qwen3.5:27b-q4_K_M") {
		t.Fatalf("models missing qwen3.5:27b-q4_K_M: %v", got)
	}
	if !containsModelID(got, "gemma3:27b-it-q4_K_M") {
		t.Fatalf("models missing gemma3:27b-it-q4_K_M: %v", got)
	}
}

func TestFallbackModels_Ollama(t *testing.T) {
	t.Parallel()

	models := FallbackModels("ollama")
	if len(models) == 0 {
		t.Fatal("expected ollama fallback catalog")
	}
	if models[0].ID != "qwen3.5:9b" {
		t.Fatalf("first model = %+v", models[0])
	}
}

func containsModelID(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

