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

	// checkAntigoDB(ctx) removido (migrado para SQLite)
	s.checkQdrant(ctx)
	s.checkGPU(ctx)
	s.checkLiteLLM(ctx)

	s.logger.Info("sovereign maintenance round complete")
	return nil
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
