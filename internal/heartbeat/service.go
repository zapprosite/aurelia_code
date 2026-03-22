package heartbeat

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/observability"
)

const (
	minIntervalMinutes     = 5
	defaultIntervalMinutes = 30
)

// HeartbeatService manages periodic health checks via agent loop
type HeartbeatService struct {
	workspace string
	loop      *agent.Loop
	interval  time.Duration
	enabled   bool
	mu        sync.RWMutex
	stopChan  chan struct{}
	logger    *slog.Logger
}

// NewHeartbeatService creates a new heartbeat service
func NewHeartbeatService(workspace string, intervalMinutes int, enabled bool, loop *agent.Loop) *HeartbeatService {
	if intervalMinutes < minIntervalMinutes && intervalMinutes != 0 {
		intervalMinutes = minIntervalMinutes
	}
	if intervalMinutes == 0 {
		intervalMinutes = defaultIntervalMinutes
	}

	return &HeartbeatService{
		workspace: workspace,
		interval:  time.Duration(intervalMinutes) * time.Minute,
		enabled:   enabled,
		loop:      loop,
		logger:    observability.Logger("heartbeat"),
	}
}

// Start begins the heartbeat service
func (hs *HeartbeatService) Start() error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if hs.stopChan != nil {
		hs.logger.Info("Heartbeat service already running")
		return nil
	}

	if !hs.enabled {
		hs.logger.Info("Heartbeat service disabled")
		return nil
	}

	if hs.loop == nil {
		hs.logger.Warn("Heartbeat service requires agent loop")
		return nil
	}

	hs.stopChan = make(chan struct{})
	go hs.runLoop(hs.stopChan)

	hs.logger.Info("Heartbeat service started", slog.Float64("interval_minutes", hs.interval.Minutes()))
	return nil
}

// Stop gracefully stops the heartbeat service
func (hs *HeartbeatService) Stop() {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if hs.stopChan == nil {
		return
	}

	hs.logger.Info("Stopping heartbeat service")
	close(hs.stopChan)
	hs.stopChan = nil
}

// IsRunning returns whether the service is running
func (hs *HeartbeatService) IsRunning() bool {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	return hs.stopChan != nil
}

// runLoop runs the heartbeat ticker
func (hs *HeartbeatService) runLoop(stopChan chan struct{}) {
	ticker := time.NewTicker(hs.interval)
	defer ticker.Stop()

	// Run first heartbeat after initial delay
	time.AfterFunc(time.Second, func() {
		hs.executeHeartbeat()
	})

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			hs.executeHeartbeat()
		}
	}
}

// executeHeartbeat performs a single heartbeat check
func (hs *HeartbeatService) executeHeartbeat() {
	hs.mu.RLock()
	enabled := hs.enabled
	if !hs.enabled || hs.stopChan == nil {
		hs.mu.RUnlock()
		return
	}
	hs.mu.RUnlock()

	if !enabled {
		return
	}

	hs.logger.Debug("Executing heartbeat")

	prompt := hs.buildPrompt()
	if prompt == "" {
		hs.logger.Info("No heartbeat prompt (HEARTBEAT.md empty or missing)")
		return
	}

	// Execute heartbeat via agent loop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	history := []agent.Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	systemPrompt := `You are a proactive home lab guardian. This is a scheduled heartbeat check.
Review the following tasks and execute any necessary actions using available tools.
If there is nothing that requires attention, respond with: HEARTBEAT_OK

Be concise. Focus on essential health checks and corrective actions.`

	_, result, err := hs.loop.RunWithOptions(ctx, systemPrompt, history, []string{
		"run_command", "docker_control", "system_monitor", "service_control",
	}, agent.RunOptions{LocalOnly: true})

	if err != nil {
		hs.logger.Error("Heartbeat error", slog.Any("err", err))
		return
	}

	hs.logger.Info("Heartbeat completed", slog.String("result", result))
}

// buildPrompt builds the heartbeat prompt from HEARTBEAT.md
func (hs *HeartbeatService) buildPrompt() string {
	heartbeatPath := filepath.Join(hs.workspace, "HEARTBEAT.md")

	data, err := os.ReadFile(heartbeatPath)
	if err != nil {
		if os.IsNotExist(err) {
			hs.createDefaultHeartbeatTemplate()
			return ""
		}
		hs.logger.Error("Error reading HEARTBEAT.md", slog.Any("err", err))
		return ""
	}

	content := string(data)
	if len(content) == 0 {
		return ""
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(`# Heartbeat Check - %s

%s`, now, content)
}

// createDefaultHeartbeatTemplate creates the default HEARTBEAT.md file
func (hs *HeartbeatService) createDefaultHeartbeatTemplate() {
	heartbeatPath := filepath.Join(hs.workspace, "HEARTBEAT.md")

	defaultContent := `# Heartbeat Check List

This file contains tasks for the heartbeat service to check periodically every 30 minutes.

## Health Checks

1. **Docker containers** — List active containers and check for unhealthy status. Restart unhealthy containers if needed.
2. **GPU status** — Check NVIDIA RTX 4090 temperature and utilization. Alert if temperature > 85°C.
3. **Disk usage** — Check disk space. Alert if > 85% used.
4. **System stats** — Report CPU load, memory usage, and uptime.
5. **Ollama service** — Check if Ollama is running. Restart if down.
6. **Service health** — Check critical systemd services (Docker, SSH, any custom services). Restart if failed.

## Instructions

- Execute ALL health checks listed above.
- For simple checks, report findings. For issues, attempt automatic remediation (restart failed services/containers).
- If all checks pass, respond with: HEARTBEAT_OK
- If issues are found and corrected, report what was done.
- Only report critical errors or unusual findings.
`

	if err := os.WriteFile(heartbeatPath, []byte(defaultContent), 0o644); err != nil {
		hs.logger.Error("Failed to create default HEARTBEAT.md", slog.Any("err", err))
	} else {
		hs.logger.Info("Created default HEARTBEAT.md template")
	}
}
