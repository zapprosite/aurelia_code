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

func TestRegisterRoute_ExposesCustomHandler(t *testing.T) {
	t.Parallel()

	srv := NewServer(0)
	srv.RegisterRoute("/v1/router/dry-run", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))

	mux := http.NewServeMux()
	mux.HandleFunc("/health", srv.handleHealth)
	mux.HandleFunc("/ready", srv.handleReady)
	srv.mu.RLock()
	for path, handler := range srv.routes {
		mux.Handle(path, handler)
	}
	srv.mu.RUnlock()

	req := httptest.NewRequest(http.MethodPost, "/v1/router/dry-run", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status code = %d", rec.Code)
	}
	if body := rec.Body.String(); body != `{"status":"ok"}` {
		t.Fatalf("body = %q", body)
	}
}
