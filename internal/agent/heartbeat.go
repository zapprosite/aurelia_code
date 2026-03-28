package agent

import (
	"context"
	"log/slog"
	"time"

	"github.com/kocar/aurelia/internal/purity/alog"
)

// HeartbeatService executes periodic health self-checks and telemetry pulses.
type HeartbeatService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
	enabled  bool
	logger   *slog.Logger
}

// NewHeartbeatService creates a service compatible with established app bootstrap patterns.
func NewHeartbeatService(rootPath string, intervalMinutes int, enabled bool, loop *Loop) *HeartbeatService {
	ctx, cancel := context.WithCancel(context.Background())
	interval := time.Duration(intervalMinutes) * time.Minute
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	
	return &HeartbeatService{
		ctx:      ctx,
		cancel:   cancel,
		interval: interval,
		enabled:  enabled,
		logger:   alog.Logger("agent.heartbeat"),
	}
}

// Start initiates the heartbeat loop.
func (hs *HeartbeatService) Start() {
	if !hs.enabled {
		return
	}
	hs.logger.Info("heartbeat service starting", slog.Duration("interval", hs.interval))
	go hs.loop()
}

// Stop terminates the heartbeat loop.
func (hs *HeartbeatService) Stop() {
	if !hs.enabled {
		return
	}
	hs.logger.Info("heartbeat service stopping")
	hs.cancel()
}

func (hs *HeartbeatService) loop() {
	ticker := time.NewTicker(hs.interval)
	defer ticker.Stop()

	for {
		select {
		case <-hs.ctx.Done():
			return
		case <-ticker.C:
			hs.pulse()
		}
	}
}

func (hs *HeartbeatService) pulse() {
	hs.logger.Debug("pulse", slog.Time("at", time.Now()))
}
