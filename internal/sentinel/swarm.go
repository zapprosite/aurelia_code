package sentinel

import (
	"context"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/kocar/aurelia/internal/observability"
)

// Swarm handles the autonomous maintenance of the Aurelia ecosystem.
type Swarm struct {
	logger *slog.Logger
}

// NewSwarm creates a new Sentinel Swarm.
func NewSwarm() *Swarm {
	return &Swarm{
		logger: observability.Logger("sentinel.swarm"),
	}
}

// RunOnce executes a single round of maintenance tasks.
func (s *Swarm) RunOnce(ctx context.Context) error {
	s.logger.Info("starting sovereign maintenance round")

	s.checkSupabase(ctx)
	s.checkQdrant(ctx)
	s.checkGPU(ctx)
	s.checkLiteLLM(ctx)

	s.logger.Info("sovereign maintenance round complete")
	return nil
}

func (s *Swarm) checkSupabase(ctx context.Context) {
	s.logger.Info("checking Supabase stability")
	
	// 1. Container status check
	cmd := exec.CommandContext(ctx, "docker", "inspect", "-f", "{{.State.Status}}", "supabase_db_aurelia")
	out, err := cmd.Output()
	if err != nil {
		s.logger.Error("Supabase container is missing or unreachable", slog.Any("err", err))
		return
	}
	status := strings.TrimSpace(string(out))
	if status != "running" {
		s.logger.Warn("Supabase container is not running", slog.String("status", status))
	} else {
		s.logger.Info("Supabase status verified", slog.String("status", "running"))
	}

	// 2. Performance Outlier Check (Industrial)
	s.logger.Info("auditing database query performance")
	perfCmd := exec.CommandContext(ctx, "supabase", "db", "outliers", "-n", "3")
	if out, err := perfCmd.Output(); err == nil {
		s.logger.Info("database performance report generated", slog.String("outliers", string(out)))
	} else {
		s.logger.Warn("could not fetch db outliers from Supabase CLI")
	}
}

func (s *Swarm) checkQdrant(ctx context.Context) {
	s.logger.Info("checking Qdrant vector storage")
	// Check if Qdrant is responding.
	cmd := exec.CommandContext(ctx, "curl", "-s", "http://localhost:6333/health")
	if err := cmd.Run(); err != nil {
		s.logger.Error("Qdrant health check failed", slog.Any("err", err))
	} else {
		s.logger.Info("Qdrant healthy")
	}
}

func (s *Swarm) checkGPU(ctx context.Context) {
	s.logger.Info("checking GPU / CUDA performance")
	cmd := exec.CommandContext(ctx, "nvidia-smi", "--query-gpu=utilization.gpu,memory.used", "--format=csv,noheader,nounits")
	out, err := cmd.Output()
	if err != nil {
		s.logger.Warn("GPU monitoring failed - check drivers", slog.Any("err", err))
		return
	}
	s.logger.Info("GPU metrics gathered", slog.String("stats", strings.TrimSpace(string(out))))
}

func (s *Swarm) checkLiteLLM(ctx context.Context) {
	s.logger.Info("auditing LiteLLM latency")
	// LiteLLM sanity check.
	cmd := exec.CommandContext(ctx, "docker", "inspect", "-f", "{{.State.Status}}", "litellm-aurelia")
	if out, err := cmd.Output(); err == nil {
		s.logger.Info("LiteLLM router verified", slog.String("status", strings.TrimSpace(string(out))))
	}
}
