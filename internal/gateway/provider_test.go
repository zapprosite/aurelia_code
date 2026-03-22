package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

type fakeProvider struct {
	response *agent.ModelResponse
	err      error
	prompts  []string
}

func (f *fakeProvider) GenerateContent(ctx context.Context, systemPrompt string, history []agent.Message, tools []agent.Tool) (*agent.ModelResponse, error) {
	f.prompts = append(f.prompts, systemPrompt)
	return f.response, f.err
}

func (f *fakeProvider) Close() {}

func newTestGatewayProvider() *Provider {
	return &Provider{
		planner: NewPlanner(),
		budgets: map[string]laneBudget{
			"local":             {Soft: 1000000, Hard: 1000000},
			"remote_cheap":      {Soft: 400, Hard: 800},
			"remote_vision":     {Soft: 120, Hard: 240},
			"remote_structured": {Soft: 200, Hard: 400},
			"remote_premium":    {Soft: 80, Hard: 160},
			"audio":             {Soft: 800, Hard: 1200},
			"research":          {Soft: 40, Hard: 80},
		},
		states: make(map[string]*routeState),
	}
}

func TestProviderGenerateContent_MaintenancePrefersLocalBalanced(t *testing.T) {
	t.Parallel()

	localBalanced := &fakeProvider{response: &agent.ModelResponse{Content: "ok local"}}
	provider := newTestGatewayProvider()
	provider.localBalanced = localBalanced

	resp, err := provider.GenerateContent(context.Background(), "system", []agent.Message{{
		Role:    "user",
		Content: "validar homelab depois do reboot",
	}}, nil)
	if err != nil {
		t.Fatalf("GenerateContent() error = %v", err)
	}
	if resp.Content != "ok local" {
		t.Fatalf("Content = %q", resp.Content)
	}
	if len(localBalanced.prompts) != 1 {
		t.Fatalf("local balanced calls = %d", len(localBalanced.prompts))
	}
	state := provider.StatusSnapshot().Routes["ollama:"+defaultLocalBalancedModel]
	if state.Requests != 1 {
		t.Fatalf("requests = %d", state.Requests)
	}
}

func TestProviderGenerateContent_FallsBackWhenGuardedResponseIsEmpty(t *testing.T) {
	t.Parallel()

	// Na nova política, curadoria prioriza remoteStructured (Gemini).
	// Vamos fazer o Gemini falhar na guarda (só reasoning) para testar o fallback pro local.
	remoteStructured := &fakeProvider{response: &agent.ModelResponse{ReasoningContent: "gemini analyzing..."}}
	localBalanced := &fakeProvider{response: &agent.ModelResponse{Content: "{\"status\":\"ok\"}"}}

	provider := newTestGatewayProvider()
	provider.localBalanced = localBalanced
	provider.remoteStructured = remoteStructured

	resp, err := provider.GenerateContent(context.Background(), "system", []agent.Message{{
		Role:    "user",
		Content: "validar homelab e resumir fatos curtos",
	}}, nil)
	if err != nil {
		t.Fatalf("GenerateContent() error = %v", err)
	}
	if resp.Content != "{\"status\":\"ok\"}" {
		t.Fatalf("Content = %q", resp.Content)
	}
	if len(localBalanced.prompts) != 1 {
		t.Fatalf("local balanced calls = %d", len(localBalanced.prompts))
	}
	if len(remoteStructured.prompts) != 1 {
		t.Fatalf("remote structured calls = %d", len(remoteStructured.prompts))
	}
	if !strings.Contains(remoteStructured.prompts[0], "# RESPONSE GUARD") {
		t.Fatalf("expected response guard in gemini prompt, got %q", remoteStructured.prompts[0])
	}
	state := provider.StatusSnapshot().Routes["openrouter:"+defaultRemoteStructuredModel]
	if state.Failures != 1 {
		t.Fatalf("gemini failures = %d", state.Failures)
	}
}

func TestProviderHealthCheck_WarnsOnOpenBreakerAndBudget(t *testing.T) {
	t.Parallel()

	provider := newTestGatewayProvider()
	provider.budgets["remote_premium"] = laneBudget{Soft: 1, Hard: 1}
	provider.states["remote_premium"] = &routeState{Requests: 1, BreakerState: "closed"}
	provider.states["openrouter:"+defaultRemotePremiumModel] = &routeState{
		BreakerState:     "open",
		BreakerOpenUntil: time.Now().Add(time.Minute),
	}

	check := provider.HealthCheck()
	if check.Status != "warning" {
		t.Fatalf("status = %q", check.Status)
	}
	if !strings.Contains(check.Message, "remote_premium budget-hard") {
		t.Fatalf("message = %q", check.Message)
	}
	if !strings.Contains(check.Message, "openrouter:"+defaultRemotePremiumModel+" open") {
		t.Fatalf("message = %q", check.Message)
	}
}

func TestProviderStatusHandler_ReturnsSnapshot(t *testing.T) {
	t.Parallel()

	provider := newTestGatewayProvider()
	provider.states["local"] = &routeState{Requests: 2, BreakerState: "closed"}

	req := httptest.NewRequest(http.MethodGet, "/v1/router/status", nil)
	rec := httptest.NewRecorder()

	provider.StatusHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d", rec.Code)
	}

	var payload StatusSnapshot
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if payload.PrimaryMode != defaultLocalBalancedModel {
		t.Fatalf("primary_mode = %q", payload.PrimaryMode)
	}
	if payload.Routes["local"].Requests != 2 {
		t.Fatalf("routes.local.requests = %d", payload.Routes["local"].Requests)
	}
}

func TestSQLiteStateStore_PersistsAcrossReload(t *testing.T) {
	t.Parallel()

	store := newSQLiteStateStore(filepath.Join(t.TempDir(), "gateway.db"))
	if store == nil {
		t.Fatal("expected sqlite state store")
	}
	defer func() { _ = store.Close() }()

	state := routeState{
		Requests:          3,
		Failures:          1,
		ConsecutiveFails:  1,
		BreakerState:      "half-open",
		LastError:         "boom",
		LastDecisionModel: "qwen3.5:9b",
	}
	if err := store.Save("remote_structured", state); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	got := loaded["remote_structured"]
	if got.Requests != 3 || got.Failures != 1 || got.BreakerState != "half-open" {
		t.Fatalf("loaded state = %+v", got)
	}
}
