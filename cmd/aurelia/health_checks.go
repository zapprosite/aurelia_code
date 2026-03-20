package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/health"
	"github.com/kocar/aurelia/internal/runtime"
)

const geminiSmokeFreshness = 24 * time.Hour

type geminiSmokeStatus struct {
	Status    string    `json:"status"`
	CheckedAt time.Time `json:"checked_at"`
	Model     string    `json:"model,omitempty"`
	Error     string    `json:"error,omitempty"`
}

type primaryLLMDescriber interface {
	PrimaryLLMDescription() string
}

type gatewayHealthReporter interface {
	HealthCheck() health.CheckResult
}

func registerAuxiliaryHealthChecks(healthSrv *health.Server, cfg *config.AppConfig, resolver *runtime.PathResolver, provider any) {
	if healthSrv == nil || cfg == nil {
		return
	}

	healthSrv.RegisterCheck("primary_llm", buildPrimaryLLMHealthCheck(cfg, provider))
	if reporter, ok := provider.(gatewayHealthReporter); ok {
		healthSrv.RegisterCheck("gateway_routing", reporter.HealthCheck)
	}

	healthSrv.RegisterCheck("gemini_api", buildGeminiHealthCheck(cfg, geminiSmokeStatusPath(resolver)))
}

func buildPrimaryLLMHealthCheck(cfg *config.AppConfig, provider any) func() health.CheckResult {
	return func() health.CheckResult {
		if cfg == nil || cfg.LLMProvider == "" {
			return health.CheckResult{Status: "warning", Message: "llm provider not configured"}
		}
		if describer, ok := provider.(primaryLLMDescriber); ok {
			return health.CheckResult{Status: "ok", Message: describer.PrimaryLLMDescription()}
		}
		if cfg.LLMProvider != "ollama" {
			return health.CheckResult{Status: "ok", Message: cfg.LLMProvider + "/" + cfg.LLMModel}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:11434/v1/models", nil)
		if err != nil {
			return health.CheckResult{Status: "error", Message: "failed to build ollama health request"}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return health.CheckResult{Status: "error", Message: "ollama not reachable: " + err.Error()}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return health.CheckResult{Status: "error", Message: "ollama returned " + resp.Status}
		}

		var payload struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return health.CheckResult{Status: "error", Message: "failed to decode ollama models response"}
		}

		for _, model := range payload.Data {
			if model.ID == cfg.LLMModel {
				return health.CheckResult{Status: "ok", Message: cfg.LLMProvider + "/" + cfg.LLMModel}
			}
		}
		return health.CheckResult{Status: "error", Message: "configured ollama model not installed: " + cfg.LLMModel}
	}
}

func geminiSmokeStatusPath(resolver *runtime.PathResolver) string {
	if resolver == nil {
		return ""
	}
	return filepath.Join(resolver.Data(), "gemini_smoke.json")
}

func buildGeminiHealthCheck(cfg *config.AppConfig, smokePath string) func() health.CheckResult {
	return func() health.CheckResult {
		if cfg == nil {
			return health.CheckResult{Status: "warning", Message: "config unavailable"}
		}
		if cfg.GoogleAPIKey == "" {
			if cfg.LLMProvider == "google" {
				return health.CheckResult{Status: "error", Message: "provider=google but google_api_key missing"}
			}
			return health.CheckResult{Status: "warning", Message: "google_api_key not configured"}
		}

		status, err := readGeminiSmokeStatus(smokePath)
		if err != nil {
			if cfg.LLMProvider == "google" {
				return health.CheckResult{Status: "error", Message: "provider=google but gemini smoke missing"}
			}
			return health.CheckResult{Status: "warning", Message: "configured, but gemini smoke has not run yet"}
		}

		age := time.Since(status.CheckedAt)
		if status.Status != "ok" {
			message := fmt.Sprintf("last gemini smoke failed (%s)", status.Error)
			if cfg.LLMProvider == "google" {
				return health.CheckResult{Status: "error", Message: message}
			}
			return health.CheckResult{Status: "warning", Message: message}
		}
		if age > geminiSmokeFreshness {
			message := fmt.Sprintf("last gemini smoke is stale (%s ago)", age.Round(time.Minute))
			if cfg.LLMProvider == "google" {
				return health.CheckResult{Status: "error", Message: message}
			}
			return health.CheckResult{Status: "warning", Message: message}
		}
		return health.CheckResult{Status: "ok", Message: fmt.Sprintf("%s validated %s ago", status.Model, age.Round(time.Minute))}
	}
}

func readGeminiSmokeStatus(path string) (geminiSmokeStatus, error) {
	var status geminiSmokeStatus
	if path == "" {
		return status, os.ErrNotExist
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return status, err
	}
	if err := json.Unmarshal(data, &status); err != nil {
		return status, err
	}
	if status.CheckedAt.IsZero() {
		return status, fmt.Errorf("missing checked_at")
	}
	return status, nil
}
