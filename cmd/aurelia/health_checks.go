package main

import (
	"encoding/json"
	"fmt"
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

func registerAuxiliaryHealthChecks(healthSrv *health.Server, cfg *config.AppConfig, resolver *runtime.PathResolver) {
	if healthSrv == nil || cfg == nil {
		return
	}

	healthSrv.RegisterCheck("primary_llm", func() health.CheckResult {
		if cfg.LLMProvider == "" {
			return health.CheckResult{Status: "warning", Message: "llm provider not configured"}
		}
		return health.CheckResult{Status: "ok", Message: cfg.LLMProvider + "/" + cfg.LLMModel}
	})

	healthSrv.RegisterCheck("gemini_api", buildGeminiHealthCheck(cfg, geminiSmokeStatusPath(resolver)))
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
