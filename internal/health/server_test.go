package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealth_WarningDoesNotDegrade(t *testing.T) {
	t.Parallel()

	srv := NewServer(0)
	srv.RegisterCheck("gemini", func() CheckResult {
		return CheckResult{Status: "warning", Message: "not configured"}
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.handleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload HealthStatus
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if payload.Status != "ok" {
		t.Fatalf("status = %q", payload.Status)
	}
	if payload.Checks["gemini"].Status != "warning" {
		t.Fatalf("gemini status = %q", payload.Checks["gemini"].Status)
	}
}

func TestHandleHealth_ErrorDegrades(t *testing.T) {
	t.Parallel()

	srv := NewServer(0)
	srv.RegisterCheck("gemini", func() CheckResult {
		return CheckResult{Status: "error", Message: "unreachable"}
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.handleHealth(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload HealthStatus
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if payload.Status != "degraded" {
		t.Fatalf("status = %q", payload.Status)
	}
}
