package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/dashboard"
)

func TestBuildDashboardStatusSnapshotWithHealthURL_UsesHealthPayload(t *testing.T) {
	t.Parallel()

	healthSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status":"ok",
			"timestamp":"2026-03-25T10:00:00Z",
			"uptime":"2h0m0s",
			"checks":{
				"voice_processor":{"status":"warning","message":"voice backlog=2"}
			}
		}`))
	}))
	defer healthSrv.Close()

	a := &app{cfg: &config.AppConfig{QdrantURL: "http://localhost:6333", QdrantCollection: "conversation_memory"}}

	snapshot := buildDashboardStatusSnapshotWithHealthURL(a, healthSrv.URL)

	if snapshot.Status != dashboard.StatusDegraded {
		t.Fatalf("snapshot status = %q", snapshot.Status)
	}

	foundHealth := false
	foundBrain := false
	for _, component := range snapshot.Components {
		switch component.Name {
		case "health_server":
			foundHealth = true
			if component.Status != dashboard.StatusUp {
				t.Fatalf("health_server status = %q", component.Status)
			}
		case "brain_search":
			foundBrain = true
			if component.Status != dashboard.StatusConfigured {
				t.Fatalf("brain_search status = %q", component.Status)
			}
		}
	}
	if !foundHealth || !foundBrain {
		t.Fatalf("components = %#v", snapshot.Components)
	}
}

func TestBuildDashboardStatusSnapshotWithHealthURL_ReportsProbeFailure(t *testing.T) {
	t.Parallel()

	snapshot := buildDashboardStatusSnapshotWithHealthURL(&app{cfg: &config.AppConfig{}}, "http://127.0.0.1:1/health")

	if snapshot.Status != dashboard.StatusDegraded {
		t.Fatalf("snapshot status = %q", snapshot.Status)
	}
}
