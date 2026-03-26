package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kocar/aurelia/internal/dashboard"
	"github.com/kocar/aurelia/internal/gateway"
	"github.com/kocar/aurelia/internal/health"
)

func buildDashboardStatusSnapshot(a *app) dashboard.StatusSnapshot {
	healthURL := ""
	if a != nil && a.cfg != nil && a.cfg.HealthPort > 0 {
		healthURL = fmt.Sprintf("http://127.0.0.1:%d/health", a.cfg.HealthPort)
	}
	return buildDashboardStatusSnapshotWithHealthURL(a, healthURL)
}

func buildDashboardStatusSnapshotWithHealthURL(a *app, healthURL string) dashboard.StatusSnapshot {
	checkedAt := time.Now().UTC()
	components := []dashboard.ComponentStatus{
		{
			Name:      "dashboard_api",
			Status:    dashboard.StatusUp,
			Summary:   "dashboard routes registered in-process",
			Source:    "cmd/aurelia/app.go",
			CheckedAt: checkedAt,
		},
	}

	var healthSnapshot *health.HealthStatus
	if healthURL == "" {
		components = append(components, dashboard.ComponentStatus{
			Name:      "health_server",
			Status:    dashboard.StatusUnknown,
			Summary:   "health port is not configured",
			Source:    "config.health_port",
			CheckedAt: checkedAt,
		})
	} else {
		snapshot, err := fetchHealthSnapshot(healthURL)
		if err != nil {
			components = append(components, dashboard.ComponentStatus{
				Name:      "health_server",
				Status:    dashboard.StatusDown,
				Summary:   "health endpoint probe failed",
				Source:    healthURL,
				CheckedAt: checkedAt,
				Details: map[string]any{
					"error": err.Error(),
				},
			})
		} else {
			healthSnapshot = snapshot
			components = append(components, dashboard.ComponentStatus{
				Name:      "health_server",
				Status:    mapHealthStatus(snapshot.Status),
				Summary:   snapshot.Status,
				Source:    healthURL,
				CheckedAt: checkedAt,
				Details: map[string]any{
					"uptime": snapshot.Uptime,
				},
			})
		}
	}

	components = append(components, describeCoreComponents(a, healthSnapshot, checkedAt)...)
	return dashboard.BuildStatusSnapshot(checkedAt, components)
}

func describeCoreComponents(a *app, healthSnapshot *health.HealthStatus, checkedAt time.Time) []dashboard.ComponentStatus {
	components := make([]dashboard.ComponentStatus, 0, 6)

	if a == nil {
		return append(components, dashboard.ComponentStatus{
			Name:      "app_container",
			Status:    dashboard.StatusDown,
			Summary:   "application container is nil",
			Source:    "cmd/aurelia/app",
			CheckedAt: checkedAt,
		})
	}

	components = append(components, simpleRuntimeStatus("cron_scheduler", a.cronScheduler != nil, checkedAt, "cmd/aurelia/app"))
	components = append(components, simpleRuntimeStatus("telegram_primary", a.primaryBot != nil, checkedAt, "cmd/aurelia/app"))
	components = append(components, simpleRuntimeStatus("telegram_pool", a.pool != nil, checkedAt, "cmd/aurelia/app"))

	llmRuntime := buildLLMRuntimeSnapshot(a, checkedAt)
	llmStatus := dashboard.StatusConfigured
	llmSummary := llmRuntime.EffectiveProvider
	if a.llmProvider != nil {
		llmStatus = dashboard.StatusUp
	}
	if llmRuntime.ViaGateway {
		llmSummary = llmRuntime.RequestedProvider + " via gateway"
	}
	components = append(components, dashboard.ComponentStatus{
		Name:      "llm_runtime",
		Status:    llmStatus,
		Summary:   llmSummary,
		Source:    "cmd/aurelia/buildLLMProvider",
		CheckedAt: checkedAt,
		Details: map[string]any{
			"requested_provider": llmRuntime.RequestedProvider,
			"requested_model":    llmRuntime.RequestedModel,
			"effective_provider": llmRuntime.EffectiveProvider,
			"effective_model":    llmRuntime.EffectiveModel,
			"via_gateway":        llmRuntime.ViaGateway,
		},
	})

	if gw, ok := a.llmProvider.(*gateway.Provider); ok && gw != nil {
		gwSnapshot := gw.StatusSnapshot()
		status := dashboard.StatusUp
		summary := gwSnapshot.PrimaryMode
		if gwSnapshot.Degraded.Active {
			status = dashboard.StatusDegraded
			summary = gwSnapshot.Degraded.Reason
		}
		components = append(components, dashboard.ComponentStatus{
			Name:      "gateway_router",
			Status:    status,
			Summary:   summary,
			Source:    "gateway.StatusSnapshot",
			CheckedAt: checkedAt,
			Details: map[string]any{
				"primary_lane": gwSnapshot.PrimaryLane,
				"degraded":     gwSnapshot.Degraded.Active,
			},
		})
	} else if a.llmProvider != nil {
		components = append(components, dashboard.ComponentStatus{
			Name:      "gateway_router",
			Status:    dashboard.StatusConfigured,
			Summary:   "llm provider active without gateway status snapshot",
			Source:    "cmd/aurelia/buildLLMProvider",
			CheckedAt: checkedAt,
		})
	}

	components = append(components, componentFromHealthCheck("voice_processor", a.voiceProcessor != nil, healthSnapshot, checkedAt))
	components = append(components, componentFromHealthCheck("voice_capture", a.voiceCapture != nil, healthSnapshot, checkedAt))

	qdrantURL := ""
	collection := ""
	if a.cfg != nil {
		qdrantURL = a.cfg.QdrantURL
		collection = a.cfg.QdrantCollection
	}
	if qdrantURL == "" {
		qdrantURL = "http://localhost:6333"
	}
	if collection == "" {
		collection = "conversation_memory"
	}
	components = append(components, dashboard.ComponentStatus{
		Name:      "brain_search",
		Status:    dashboard.StatusConfigured,
		Summary:   "dashboard brain routes are configured; query health is not probed here",
		Source:    "cmd/aurelia/brain_handlers.go",
		CheckedAt: checkedAt,
		Details: map[string]any{
			"qdrant_url": qdrantURL,
			"collection": collection,
		},
	})

	return components
}

func simpleRuntimeStatus(name string, active bool, checkedAt time.Time, source string) dashboard.ComponentStatus {
	status := dashboard.StatusDown
	summary := "component is not initialized"
	if active {
		status = dashboard.StatusUp
		summary = "component is initialized in-process"
	}
	return dashboard.ComponentStatus{
		Name:      name,
		Status:    status,
		Summary:   summary,
		Source:    source,
		CheckedAt: checkedAt,
	}
}

func componentFromHealthCheck(name string, enabled bool, snapshot *health.HealthStatus, checkedAt time.Time) dashboard.ComponentStatus {
	if !enabled {
		return dashboard.ComponentStatus{
			Name:      name,
			Status:    dashboard.StatusDisabled,
			Summary:   "component is disabled",
			Source:    "cmd/aurelia/app",
			CheckedAt: checkedAt,
		}
	}
	if snapshot == nil {
		return dashboard.ComponentStatus{
			Name:      name,
			Status:    dashboard.StatusUnknown,
			Summary:   "health snapshot unavailable",
			Source:    "health:/health",
			CheckedAt: checkedAt,
		}
	}
	check, ok := snapshot.Checks[name]
	if !ok {
		return dashboard.ComponentStatus{
			Name:      name,
			Status:    dashboard.StatusUnknown,
			Summary:   "component enabled but no named health check is registered",
			Source:    "health:/health",
			CheckedAt: checkedAt,
		}
	}
	return dashboard.ComponentStatus{
		Name:      name,
		Status:    mapHealthStatus(check.Status),
		Summary:   check.Message,
		Source:    "health:/health",
		CheckedAt: checkedAt,
	}
}

func mapHealthStatus(status string) string {
	switch status {
	case "ok", "ready":
		return dashboard.StatusUp
	case "warning":
		return dashboard.StatusDegraded
	case "error", "degraded":
		return dashboard.StatusDown
	default:
		return dashboard.StatusUnknown
	}
}

func fetchHealthSnapshot(endpoint string) (*health.HealthStatus, error) {
	client := &http.Client{Timeout: 750 * time.Millisecond}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snapshot health.HealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}
