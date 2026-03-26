package dashboard

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	StatusUp         = "up"
	StatusDown       = "down"
	StatusDegraded   = "degraded"
	StatusDisabled   = "disabled"
	StatusConfigured = "configured"
	StatusUnknown    = "unknown"
)

// ComponentStatus is an explicit per-component status contract for dashboard consumers.
type ComponentStatus struct {
	Name      string         `json:"name"`
	Status    string         `json:"status"`
	Summary   string         `json:"summary"`
	Source    string         `json:"source"`
	CheckedAt time.Time      `json:"checked_at"`
	Details   map[string]any `json:"details,omitempty"`
}

// StatusSnapshot is a coarse-grained operational snapshot for the dashboard.
type StatusSnapshot struct {
	Status     string            `json:"status"`
	Timestamp  time.Time         `json:"timestamp"`
	Components []ComponentStatus `json:"components"`
}

// BuildStatusSnapshot derives an overall status from explicit component states.
func BuildStatusSnapshot(checkedAt time.Time, components []ComponentStatus) StatusSnapshot {
	overall := StatusUp
	for _, component := range components {
		switch component.Status {
		case StatusDown:
			overall = StatusDegraded
		case StatusDegraded:
			if overall == StatusUp {
				overall = StatusDegraded
			}
		case StatusUnknown:
			if overall == StatusUp {
				overall = StatusUnknown
			}
		}
	}
	return StatusSnapshot{
		Status:     overall,
		Timestamp:  checkedAt,
		Components: components,
	}
}

// NewStatusHandler exposes a snapshot route with stable JSON and dashboard-friendly CORS.
func NewStatusHandler(provider func(r *http.Request) StatusSnapshot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		snapshot := provider(r)
		statusCode := http.StatusOK
		if snapshot.Status == StatusDegraded {
			statusCode = http.StatusServiceUnavailable
		}
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(snapshot)
	}
}
