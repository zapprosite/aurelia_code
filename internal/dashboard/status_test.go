package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuildStatusSnapshot_DegradesOnDownComponent(t *testing.T) {
	t.Parallel()

	snapshot := BuildStatusSnapshot(time.Unix(10, 0).UTC(), []ComponentStatus{
		{Name: "dashboard_api", Status: StatusUp},
		{Name: "health_server", Status: StatusDown},
	})

	if snapshot.Status != StatusDegraded {
		t.Fatalf("status = %q", snapshot.Status)
	}
}

func TestBuildStatusSnapshot_DegradesOnDegradedComponent(t *testing.T) {
	t.Parallel()

	snapshot := BuildStatusSnapshot(time.Unix(10, 0).UTC(), []ComponentStatus{
		{Name: "dashboard_api", Status: StatusUp},
		{Name: "voice_capture", Status: StatusDegraded},
	})

	if snapshot.Status != StatusDegraded {
		t.Fatalf("status = %q", snapshot.Status)
	}
}

func TestNewStatusHandler_ReturnsSnapshotJSON(t *testing.T) {
	t.Parallel()

	handler := NewStatusHandler(func(r *http.Request) StatusSnapshot {
		return BuildStatusSnapshot(time.Unix(20, 0).UTC(), []ComponentStatus{
			{Name: "dashboard_api", Status: StatusUp, Summary: "ok", Source: "test"},
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload StatusSnapshot
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if payload.Status != StatusUp {
		t.Fatalf("payload status = %q", payload.Status)
	}
	if len(payload.Components) != 1 || payload.Components[0].Name != "dashboard_api" {
		t.Fatalf("components = %#v", payload.Components)
	}
}
