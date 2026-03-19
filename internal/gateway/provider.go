package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/health"
	"github.com/kocar/aurelia/pkg/llm"
)

const (
	breakerFailureThreshold = 3
	breakerCooldown         = 2 * time.Minute
)

type laneBudget struct {
	Soft int `json:"soft"`
	Hard int `json:"hard"`
}

type routeState struct {
	Requests          int       `json:"requests"`
	Failures          int       `json:"failures"`
	ConsecutiveFails  int       `json:"consecutive_failures"`
	BreakerState      string    `json:"breaker_state"`
	BreakerOpenUntil  time.Time `json:"breaker_open_until,omitempty"`
	LastError         string    `json:"last_error,omitempty"`
	LastDecisionModel string    `json:"last_decision_model,omitempty"`
}

type StatusSnapshot struct {
	PrimaryLane string                `json:"primary_lane"`
	PrimaryMode string                `json:"primary_mode"`
	Budgets     map[string]laneBudget `json:"budgets"`
	Routes      map[string]routeState `json:"routes"`
}

type closableProvider interface {
	agent.LLMProvider
	Close()
}

type Provider struct {
	planner *Planner
	metrics *gatewayMetrics

	localFast         closableProvider
	localBalanced     closableProvider
	remoteCheapLong   closableProvider
	remoteCheapVision closableProvider
	remoteStructured  closableProvider
	remotePremium     closableProvider

	mu      sync.Mutex
	budgets map[string]laneBudget
	states  map[string]*routeState
}

func NewProvider(cfg *config.AppConfig) (*Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	lowTemp := 0.1
	localFast := llm.NewOllamaProviderWithOptions(defaultLocalFastModel, llm.OpenAICompatibleRequestOptions{
		MaxTokens:   192,
		Temperature: &lowTemp,
		ExtraFields: map[string]any{"think": false},
	})
	localBalanced := llm.NewOllamaProviderWithOptions(defaultLocalBalancedModel, llm.OpenAICompatibleRequestOptions{
		MaxTokens:   384,
		Temperature: &lowTemp,
		ExtraFields: map[string]any{"think": false},
	})

	var remoteCheapLong closableProvider
	var remoteCheapVision closableProvider
	var remoteStructured closableProvider
	var remotePremium closableProvider
	if cfg.OpenRouterAPIKey != "" {
		remoteCheapLong = llm.NewOpenRouterProviderWithOptions(cfg.OpenRouterAPIKey, defaultRemoteCheapLongModel, llm.OpenAICompatibleRequestOptions{
			MaxTokens:   384,
			Temperature: &lowTemp,
		})
		remoteCheapVision = llm.NewOpenRouterProviderWithOptions(cfg.OpenRouterAPIKey, defaultRemoteVisionModel, llm.OpenAICompatibleRequestOptions{
			MaxTokens:   384,
			Temperature: &lowTemp,
		})
		remoteStructured = llm.NewOpenRouterProviderWithOptions(cfg.OpenRouterAPIKey, defaultRemoteStructuredModel, llm.OpenAICompatibleRequestOptions{
			MaxTokens:   256,
			Temperature: &lowTemp,
		})
		remotePremium = llm.NewOpenRouterProviderWithOptions(cfg.OpenRouterAPIKey, defaultRemotePremiumModel, llm.OpenAICompatibleRequestOptions{
			MaxTokens:   640,
			Temperature: &lowTemp,
		})
	}

	return &Provider{
		planner:           NewPlanner(),
		metrics:           defaultMetrics(),
		localFast:         localFast,
		localBalanced:     localBalanced,
		remoteCheapLong:   remoteCheapLong,
		remoteCheapVision: remoteCheapVision,
		remoteStructured:  remoteStructured,
		remotePremium:     remotePremium,
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
	}, nil
}

func (p *Provider) Close() {
	for _, provider := range []closableProvider{
		p.localFast,
		p.localBalanced,
		p.remoteCheapLong,
		p.remoteCheapVision,
		p.remoteStructured,
		p.remotePremium,
	} {
		if provider != nil {
			provider.Close()
		}
	}
}

func (p *Provider) PrimaryLLMDescription() string {
	return "gateway/qwen3.5:9b"
}

func (p *Provider) GenerateContent(ctx context.Context, systemPrompt string, history []agent.Message, tools []agent.Tool) (*agent.ModelResponse, error) {
	req := buildDryRunRequest(systemPrompt, history, tools)
	decision := p.planner.Decide(req)
	resp, err := p.generateWithDecision(ctx, decision, systemPrompt, history, tools)
	if err == nil && responseSatisfies(decision, resp) {
		return resp, nil
	}

	fallbacks := p.fallbackDecisions(req, decision)
	lastErr := err
	for _, fallback := range fallbacks {
		p.recordFallback(decision, fallback)
		resp, err = p.generateWithDecision(ctx, fallback, systemPrompt, history, tools)
		if err == nil && responseSatisfies(fallback, resp) {
			return resp, nil
		}
		lastErr = err
	}
	if resp != nil && responseSatisfies(decision, resp) {
		return resp, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("all gateway routes failed or returned empty guarded content")
	}
	return nil, lastErr
}

func (p *Provider) StatusSnapshot() StatusSnapshot {
	p.mu.Lock()
	defer p.mu.Unlock()

	routes := make(map[string]routeState, len(p.states))
	for key, state := range p.states {
		routes[key] = *state
	}
	budgets := make(map[string]laneBudget, len(p.budgets))
	for key, value := range p.budgets {
		budgets[key] = value
	}

	return StatusSnapshot{
		PrimaryLane: "local-balanced",
		PrimaryMode: defaultLocalBalancedModel,
		Budgets:     budgets,
		Routes:      routes,
	}
}

func (p *Provider) HealthCheck() health.CheckResult {
	snapshot := p.StatusSnapshot()
	var warnings []string
	for key, state := range snapshot.Routes {
		if state.BreakerState == "open" {
			warnings = append(warnings, key+" open")
		}
		budget, ok := snapshot.Budgets[key]
		if ok && budget.Hard > 0 && state.Requests >= budget.Hard {
			warnings = append(warnings, key+" budget-hard")
		}
	}
	if len(warnings) == 0 {
		return health.CheckResult{Status: "ok", Message: snapshot.PrimaryMode}
	}
	return health.CheckResult{Status: "warning", Message: strings.Join(warnings, ", ")}
}

func (p *Provider) StatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(p.StatusSnapshot())
	})
}

func (p *Provider) generateWithDecision(ctx context.Context, decision DryRunDecision, systemPrompt string, history []agent.Message, tools []agent.Tool) (*agent.ModelResponse, error) {
	startedAt := time.Now()
	providerKey := decision.Provider + ":" + decision.Model
	if p.breakerOpen(providerKey) {
		p.observeResult(decision, "breaker_open", time.Since(startedAt))
		return nil, fmt.Errorf("route breaker open for %s", providerKey)
	}
	if p.budgetExceeded(decision.BudgetLane) {
		p.observeResult(decision, "budget_exceeded", time.Since(startedAt))
		return nil, fmt.Errorf("budget exceeded for %s", decision.BudgetLane)
	}

	provider := p.providerFor(decision)
	if provider == nil {
		p.observeResult(decision, "provider_unavailable", time.Since(startedAt))
		return nil, fmt.Errorf("provider unavailable for %s", providerKey)
	}

	systemPrompt = applyGuards(systemPrompt, decision)
	p.markRequest(decision.BudgetLane, providerKey, decision.Model)

	resp, err := provider.GenerateContent(ctx, systemPrompt, history, tools)
	if err != nil {
		p.markFailure(decision.BudgetLane, providerKey, err.Error())
		p.observeResult(decision, "error", time.Since(startedAt))
		return nil, err
	}
	if !responseSatisfies(decision, resp) {
		p.markFailure(decision.BudgetLane, providerKey, "empty guarded content")
		p.observeResult(decision, "guard_reject", time.Since(startedAt))
		return resp, fmt.Errorf("empty guarded content for %s", providerKey)
	}
	p.markSuccess(providerKey)
	p.observeResult(decision, "success", time.Since(startedAt))
	return resp, nil
}

func (p *Provider) providerFor(decision DryRunDecision) closableProvider {
	switch decision.Lane {
	case "local-fast":
		return p.localFast
	case "local-balanced":
		return p.localBalanced
	case "remote-cheap-long-context":
		return p.remoteCheapLong
	case "remote-cheap-vision":
		return p.remoteCheapVision
	case "remote-tool-long-output":
		return p.remoteStructured
	case "remote-premium-workflow":
		return p.remotePremium
	default:
		return p.localBalanced
	}
}

func (p *Provider) fallbackDecisions(req DryRunRequest, original DryRunDecision) []DryRunDecision {
	switch original.Lane {
	case "remote-premium-workflow":
		return []DryRunDecision{
			p.planner.Decide(DryRunRequest{Task: req.Task, TaskClass: "curation", OutputMode: req.OutputMode, RequiresTools: req.RequiresTools, LocalOnly: false}),
			p.planner.Decide(DryRunRequest{Task: req.Task, TaskClass: "maintenance", RequiresTools: req.RequiresTools, LocalOnly: true}),
		}
	case "remote-tool-long-output", "remote-cheap-long-context", "remote-cheap-vision":
		return []DryRunDecision{
			p.planner.Decide(DryRunRequest{Task: req.Task, TaskClass: "maintenance", OutputMode: req.OutputMode, RequiresTools: req.RequiresTools, LocalOnly: true, LatencyBudgetMS: req.LatencyBudgetMS}),
		}
	default:
		return []DryRunDecision{
			{
				Lane:       "remote-tool-long-output",
				Provider:   "openrouter",
				Model:      defaultRemoteStructuredModel,
				UseRemote:  true,
				UseTools:   req.RequiresTools,
				Reason:     "fallback estruturado para lane remota previsivel",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_structured",
			},
		}
	}
}

func buildDryRunRequest(systemPrompt string, history []agent.Message, tools []agent.Tool) DryRunRequest {
	var latestUser string
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			latestUser = history[i].Content
			break
		}
	}
	outputMode := "text"
	systemLower := strings.ToLower(systemPrompt)
	userLower := strings.ToLower(latestUser)
	switch {
	case strings.Contains(systemLower, "json valido") || strings.Contains(systemLower, "only valid json"):
		outputMode = "structured_json"
	case strings.Contains(userLower, "3 facts") || strings.Contains(userLower, "3 tags") || strings.Contains(userLower, "facts curtos"):
		outputMode = "curation"
	}

	taskClass := ""
	switch {
	case strings.Contains(userLower, "screenshot") || strings.Contains(userLower, "ide"):
		taskClass = "vision"
	case strings.Contains(userLower, "roteamento") || strings.Contains(systemLower, "classifier"):
		taskClass = "routing"
	case outputMode == "curation":
		taskClass = "curation"
	case strings.Contains(userLower, "reboot") || strings.Contains(userLower, "nvidia-gpu") || strings.Contains(userLower, "homelab"):
		taskClass = "maintenance"
	}

	return DryRunRequest{
		Task:           latestUser,
		TaskClass:      taskClass,
		OutputMode:     outputMode,
		RequiresTools:  len(tools) > 0,
		RequiresVision: strings.Contains(userLower, "screenshot") || strings.Contains(userLower, "imagem"),
		LocalOnly:      strings.Contains(systemLower, "local only") || strings.Contains(userLower, "local only"),
	}
}

func applyGuards(systemPrompt string, decision DryRunDecision) string {
	switch decision.Guards.ReasoningMode {
	case "minimize":
		return strings.TrimSpace(systemPrompt) + "\n\n# RESPONSE GUARD\n- Minimize reasoning.\n- Do not spend the answer budget on hidden analysis.\n- Prioritize final content over internal deliberation.\n- If structured output was requested, return the final structure directly."
	default:
		return systemPrompt
	}
}

func responseSatisfies(decision DryRunDecision, resp *agent.ModelResponse) bool {
	if resp == nil {
		return false
	}
	if len(resp.ToolCalls) > 0 {
		return true
	}
	if strings.TrimSpace(resp.Content) != "" {
		return true
	}
	if decision.Guards.ReasoningMode == "minimize" {
		return false
	}
	return strings.TrimSpace(resp.ReasoningContent) != ""
}

func (p *Provider) markRequest(budgetLane, providerKey, model string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.ensureStateLocked(providerKey)
	state.Requests++
	state.LastDecisionModel = model
	if budgetLane != "" {
		budgetState := p.ensureStateLocked(budgetLane)
		budgetState.Requests++
		p.updateBudgetMetricLocked(budgetLane, budgetState.Requests)
	}
}

func (p *Provider) markSuccess(providerKey string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.ensureStateLocked(providerKey)
	state.ConsecutiveFails = 0
	state.BreakerState = "closed"
	state.LastError = ""
	p.setBreakerMetricLocked(providerKey, state.BreakerState)
}

func (p *Provider) markFailure(budgetLane, providerKey, reason string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.ensureStateLocked(providerKey)
	state.Failures++
	state.ConsecutiveFails++
	state.LastError = reason
	if state.ConsecutiveFails >= breakerFailureThreshold {
		state.BreakerState = "open"
		state.BreakerOpenUntil = time.Now().Add(breakerCooldown)
	}
	p.setBreakerMetricLocked(providerKey, state.BreakerState)
	if budgetLane != "" {
		budgetState := p.ensureStateLocked(budgetLane)
		budgetState.Failures++
	}
}

func (p *Provider) breakerOpen(providerKey string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.ensureStateLocked(providerKey)
	if state.BreakerState != "open" {
		return false
	}
	if time.Now().After(state.BreakerOpenUntil) {
		state.BreakerState = "half-open"
		p.setBreakerMetricLocked(providerKey, state.BreakerState)
		return false
	}
	p.setBreakerMetricLocked(providerKey, state.BreakerState)
	return true
}

func (p *Provider) budgetExceeded(lane string) bool {
	if lane == "" {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.ensureStateLocked(lane)
	budget, ok := p.budgets[lane]
	if !ok || budget.Hard <= 0 {
		return false
	}
	return state.Requests >= budget.Hard
}

func (p *Provider) ensureStateLocked(key string) *routeState {
	state, ok := p.states[key]
	if !ok {
		state = &routeState{BreakerState: "closed"}
		p.states[key] = state
	}
	return state
}

func (p *Provider) observeResult(decision DryRunDecision, result string, duration time.Duration) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.requests.WithLabelValues(decision.Lane, decision.Provider, decision.Model, result).Inc()
	p.metrics.latency.WithLabelValues(decision.Lane, decision.Provider, decision.Model, result).Observe(duration.Seconds())
	if result != "success" {
		p.metrics.failures.WithLabelValues(decision.Lane, decision.Provider, decision.Model).Inc()
	}
}

func (p *Provider) recordFallback(from, to DryRunDecision) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.fallbacks.WithLabelValues(from.Lane, to.Lane).Inc()
}

func (p *Provider) updateBudgetMetricLocked(lane string, requests int) {
	if p == nil || p.metrics == nil || lane == "" {
		return
	}
	budget, ok := p.budgets[lane]
	if !ok || budget.Hard <= 0 {
		return
	}
	p.metrics.budgets.WithLabelValues(lane).Set(float64(requests) / float64(budget.Hard))
}

func (p *Provider) setBreakerMetricLocked(providerKey, state string) {
	if p == nil || p.metrics == nil || providerKey == "" {
		return
	}
	parts := strings.SplitN(providerKey, ":", 2)
	if len(parts) != 2 {
		return
	}
	value := 0.0
	switch state {
	case "open":
		value = 1
	case "half-open":
		value = 0.5
	}
	p.metrics.breakers.WithLabelValues(parts[0], parts[1]).Set(value)
}
