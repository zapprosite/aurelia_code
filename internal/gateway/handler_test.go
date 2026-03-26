package gateway

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDryRunHandler_ReturnsDecision(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(DryRunRequest{
		TaskClass:  "routing",
		OutputMode: "structured_json",
	})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/router/dry-run", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	NewDryRunHandler(NewPlanner()).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload DryRunDecision
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if payload.Model != modelDeepSeekV31 {
		t.Fatalf("unexpected decision: %+v", payload)
	}
}

func TestDryRunHandler_RejectsInvalidPayload(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/v1/router/dry-run", bytes.NewBufferString(`{`))
	rec := httptest.NewRecorder()

	NewDryRunHandler(NewPlanner()).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code = %d", rec.Code)
	}
}
