package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/config"
)

func TestBuildGeminiHealthCheck_ConfiguredButNoSmokeIsWarningForNonGoogleProvider(t *testing.T) {
	t.Parallel()

	check := buildGeminiHealthCheck(&config.AppConfig{
		LLMProvider:  "openrouter",
		GoogleAPIKey: "secret",
	}, filepath.Join(t.TempDir(), "missing.json"))

	result := check()
	if result.Status != "warning" {
		t.Fatalf("status = %q", result.Status)
	}
}

func TestBuildGeminiHealthCheck_MissingKeyIsErrorForGoogleProvider(t *testing.T) {
	t.Parallel()

	check := buildGeminiHealthCheck(&config.AppConfig{
		LLMProvider: "google",
	}, "")

	result := check()
	if result.Status != "error" {
		t.Fatalf("status = %q", result.Status)
	}
}

func TestBuildGeminiHealthCheck_RecentSmokeIsOk(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "gemini_smoke.json")
	payload := geminiSmokeStatus{
		Status:    "ok",
		Model:     "gemini-2.5-flash",
		CheckedAt: time.Now().Add(-5 * time.Minute),
	}
	if err := writeGeminiSmokeFixture(path, payload); err != nil {
		t.Fatalf("writeGeminiSmokeFixture() error = %v", err)
	}

	check := buildGeminiHealthCheck(&config.AppConfig{
		LLMProvider:  "openrouter",
		GoogleAPIKey: "secret",
	}, path)

	result := check()
	if result.Status != "ok" {
		t.Fatalf("status = %q", result.Status)
	}
}

func TestBuildPrimaryLLMHealthCheck_OllamaModelInstalledIsOk(t *testing.T) {
	t.Parallel()

	originalTransport := http.DefaultTransport
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"gemma3:27b"}]}`))
	}))
	defer server.Close()
	http.DefaultTransport = rewriteLocalhostTransport(t, server)
	defer func() { http.DefaultTransport = originalTransport }()

	check := buildPrimaryLLMHealthCheck(&config.AppConfig{
		LLMProvider: "ollama",
		LLMModel:    "gemma3:27b",
	}, nil)

	result := check()
	if result.Status != "ok" {
		t.Fatalf("status = %q message=%q", result.Status, result.Message)
	}
}

type fakePrimaryDescriber struct{}

func (fakePrimaryDescriber) PrimaryLLMDescription() string { return "gateway/gemma3:27b" }

func TestBuildPrimaryLLMHealthCheck_UsesProviderDescriptionWhenAvailable(t *testing.T) {
	t.Parallel()

	check := buildPrimaryLLMHealthCheck(&config.AppConfig{
		LLMProvider: "openrouter",
		LLMModel:    "minimax/minimax-m2.7",
	}, fakePrimaryDescriber{})

	result := check()
	if result.Status != "ok" {
		t.Fatalf("status = %q", result.Status)
	}
	if result.Message != "gateway/gemma3:27b" {
		t.Fatalf("message = %q", result.Message)
	}
}

func rewriteLocalhostTransport(t *testing.T, server *httptest.Server) http.RoundTripper {
	t.Helper()
	base := server.Client().Transport
	return roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "127.0.0.1:11434" {
			cloned := req.Clone(req.Context())
			cloned.URL.Scheme = "http"
			cloned.URL.Host = strings.TrimPrefix(server.URL, "http://")
			return base.RoundTrip(cloned)
		}
		return base.RoundTrip(req)
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func writeGeminiSmokeFixture(path string, payload geminiSmokeStatus) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
