package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
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

		ollamaURL := cfg.OllamaURL
		if ollamaURL == "" {
			ollamaURL = "http://127.0.0.1:11434"
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(ollamaURL, "/")+"/v1/models", nil)
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

// startSentinelHealthLoop runs periodic health probes and updates squad status for gemma and sentinel.
// S-24: Sentinel→Squad Health Probes.
func startSentinelHealthLoop(cfg *config.AppConfig) {
	ollamaURL := cfg.OllamaURL
	if ollamaURL == "" {
		ollamaURL = "http://127.0.0.1:11434"
	}
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			runSentinelProbe(ollamaURL)
		}
	}()
}

func runSentinelProbe(ollamaURL string) {
	totalChecks := 2
	passed := 0

	// Probe 1: Ollama
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(ollamaURL, "/")+"/api/tags", nil)
	if err == nil {
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				passed++
				agent.UpdateSquadAgentStatus("gemma", "online", 10)
			} else {
				agent.UpdateSquadAgentStatus("gemma", "offline", 0)
			}
		} else {
			agent.UpdateSquadAgentStatus("gemma", "offline", 0)
		}
	}

	// Probe 2: self health endpoint
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	req2, err2 := http.NewRequestWithContext(ctx2, http.MethodGet, "http://localhost:8080/ready", nil)
	if err2 == nil {
		resp2, err2 := http.DefaultClient.Do(req2)
		if err2 == nil {
			resp2.Body.Close()
			if resp2.StatusCode == http.StatusOK {
				passed++
			}
		}
	}

	sentinelLoad := (passed * 100) / totalChecks
	agent.UpdateSquadAgentStatus("sentinel", "online", 100-sentinelLoad) // more checks passing = lower load
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

// ── Instance guard ────────────────────────────────────────────────────────────

// shouldSuppressDuplicateLaunch returns true when a duplicate-launch error
// was caused by an orphan process (PID 1 = systemd adopted it) rather than
// a legitimate service conflict.
func shouldSuppressDuplicateLaunch(parentPID int, err error) bool {
	if err == nil {
		return false
	}
	if !strings.Contains(err.Error(), "another Aurelia instance is already running") {
		return false
	}
	return parentPID == 1
}

func recordDuplicateLaunch(args []string, err error) {
	home := os.Getenv("AURELIA_HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOME"), ".aurelia")
	}
	logDir := filepath.Join(home, "logs")
	_ = os.MkdirAll(logDir, 0o755)
	logPath := filepath.Join(logDir, "duplicate-launch.log")
	entry := fmt.Sprintf("[%s] args=%v error=%v\n", time.Now().Format(time.RFC3339), args, err)
	_ = os.WriteFile(logPath, []byte(entry), 0o644)
}

func exitCodeForBootstrapError(logger *slog.Logger, args []string, err error) int {
	if shouldSuppressDuplicateLaunch(os.Getppid(), err) {
		recordDuplicateLaunch(args, err)
		logger.Warn("suppressed orphan duplicate launch", slog.Any("err", err))
		return 0
	}
	logger.Error("bootstrap failed", slog.Any("err", err))
	return 1
}
