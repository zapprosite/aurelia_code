package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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
