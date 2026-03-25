package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/dashboard"
	"github.com/kocar/aurelia/internal/health"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/pkg/llm"
)

const (
	breakerFailureThreshold = 3
	breakerCooldown         = 2 * time.Minute
)

type laneBudget struct {
	Soft        int     `json:"soft"`
	Hard        int     `json:"hard"`
	CostHardUSD float64 `json:"cost_hard_usd,omitempty"`
}

type routeState struct {
	Day               string    `json:"day,omitempty"`
	Requests          int       `json:"requests"`
	Failures          int       `json:"failures"`
	ConsecutiveFails  int       `json:"consecutive_failures"`
	BreakerState      string    `json:"breaker_state"`
	BreakerOpenUntil  time.Time `json:"breaker_open_until,omitempty"`
	LastError         string    `json:"last_error,omitempty"`
	LastDecisionModel string    `json:"last_decision_model,omitempty"`
	TotalInputTokens  int       `json:"total_input_tokens"`
	TotalOutputTokens int       `json:"total_output_tokens"`
	TotalCostUSD      float64   `json:"total_cost_usd"`
	RateLimitRequests int       `json:"ratelimit_requests"`
	RateLimitRemaining int      `json:"ratelimit_remaining"`
	RateLimitReset     time.Time `json:"ratelimit_reset"`
}

type DegradedState struct {
	Active         bool      `json:"active"`
	Reason         string    `json:"reason"`
	CheckedAt      time.Time `json:"checked_at"`
	RemoteDisabled bool      `json:"remote_disabled"`
}

type StatusSnapshot struct {
	PrimaryLane string                `json:"primary_lane"`
	PrimaryMode string                `json:"primary_mode"`
	Budgets     map[string]laneBudget `json:"budgets"`
	Routes      map[string]routeState `json:"routes"`
	Degraded    DegradedState         `json:"degraded"`
}

type closableProvider interface {
	agent.LLMProvider
	Close()
}

type Provider struct {
	planner *Planner
	metrics *gatewayMetrics
	store   StateStore

	localFast         closableProvider
	localBalanced     closableProvider
	remoteCheapLong   closableProvider
	remoteCheapVision closableProvider
	remotePremium     closableProvider
	miniMaxDirect     closableProvider
	localVision       closableProvider
	judge             Judge

	mu       sync.Mutex
	budgets  map[string]laneBudget
	states   map[string]*routeState
	degraded DegradedState
}

func NewProvider(cfg *config.AppConfig) (*Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	lowTemp := 0.1
	localFast := llm.NewOllamaProviderWithOptions(cfg.OllamaURL, modelGemma3, llm.OpenAICompatibleRequestOptions{
		MaxTokens:   192,
		Temperature: &lowTemp,
		ExtraFields: map[string]any{"think": false},
	})
	localBalanced := llm.NewOllamaProviderWithOptions(cfg.OllamaURL, modelGemma3, llm.OpenAICompatibleRequestOptions{
		MaxTokens:   1024,
		Temperature: &lowTemp,
		ExtraFields: map[string]any{"think": false},
	})
	localVision := llm.NewOllamaProviderWithOptions(cfg.OllamaURL, modelGemma3, llm.OpenAICompatibleRequestOptions{
		MaxTokens:   512,
		Temperature: &lowTemp,
	})

	var remoteCheapLong closableProvider
	var remoteCheapVision closableProvider
	var remotePremium closableProvider
	if cfg.OpenRouterAPIKey != "" {
		remoteCheapLong = llm.NewOpenRouterProviderWithOptions(cfg.OpenRouterAPIKey, modelDeepSeekChat, llm.OpenAICompatibleRequestOptions{
			MaxTokens:   1024,
			Temperature: &lowTemp,
		})
		remoteCheapVision = llm.NewOpenRouterProviderWithOptions(cfg.OpenRouterAPIKey, modelKimiK25, llm.OpenAICompatibleRequestOptions{
			MaxTokens:   512,
			Temperature: &lowTemp,
		})
		remotePremium = llm.NewOpenRouterProviderWithOptions(cfg.OpenRouterAPIKey, modelMiniMaxM27, llm.OpenAICompatibleRequestOptions{
			MaxTokens:   1024,
			Temperature: &lowTemp,
		})
	}

	var miniMaxDirect closableProvider
	if cfg.MiniMaxAPIKey != "" {
		miniMaxDirect = llm.NewMiniMaxProviderWithOptions(cfg.MiniMaxAPIKey, modelMiniMaxDirect, llm.OpenAICompatibleRequestOptions{
			MaxTokens:   1024,
			Temperature: &lowTemp,
		})
	}

	judge := NewGemmaJudge(cfg.OllamaURL, "gemma3:12b")

	provider := &Provider{
		planner:           NewPlanner(),
		metrics:           defaultMetrics(),
		store:             newSQLiteStateStore(cfg.DBPath),
		localFast:         localFast,
		localBalanced:     localBalanced,
		localVision:       localVision,
		remoteCheapLong:   remoteCheapLong,
		remoteCheapVision: remoteCheapVision,
		remotePremium:     remotePremium,
		miniMaxDirect:     miniMaxDirect,
		judge:             judge,
		budgets: map[string]laneBudget{
			"local":             {Soft: 1000000, Hard: 1000000},
			"remote_cheap":      {Soft: 400, Hard: 800, CostHardUSD: 0.50},
			"remote_vision":     {Soft: 120, Hard: 240, CostHardUSD: 0.25},
			"remote_premium":    {Soft: 80, Hard: 160, CostHardUSD: 2.00},
			"audio":             {Soft: 800, Hard: 1200},
			"research":          {Soft: 40, Hard: 80},
		},
		states: make(map[string]*routeState),
		degraded: DegradedState{
			Active: false,
		},
	}
	if err := provider.loadState(); err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *Provider) EnableDegradedMode(reason string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.degraded = DegradedState{
		Active:         true,
		Reason:         reason,
		CheckedAt:      time.Now(),
		RemoteDisabled: true,
	}
}

func (p *Provider) IsDegraded() DegradedState {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.degraded
}

func (p *Provider) Close() {
	if p.store != nil {
		_ = p.store.Close()
	}
	for _, provider := range []closableProvider{
		p.localFast,
		p.localBalanced,
		p.localVision,
		p.remoteCheapLong,
		p.remoteCheapVision,
		p.remotePremium,
		p.miniMaxDirect,
	} {
		if provider != nil {
			provider.Close()
		}
	}
}

func (p *Provider) PrimaryLLMDescription() string {
	return "gateway/" + modelGemma3
}

// modelSupportsTools returns whether the model supports OpenAI-style tool calling.
// Gemma3 does not support tools via Ollama's OpenAI-compatible API.
func modelSupportsTools(model string) bool {
	switch {
	case strings.HasPrefix(model, "gemma"):
		return false
	default:
		return true
	}
}

func (p *Provider) GenerateContent(ctx context.Context, systemPrompt string, history []agent.Message, tools []agent.Tool) (*agent.ModelResponse, error) {
	opts, _ := agent.RunOptionsFromContext(ctx)

	// Calculate context size roughly
	contextSize := len(systemPrompt)
	for _, m := range history {
		contextSize += len(m.Content)
	}

	// 1. Call Judge for classification
	var latestUser string
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			latestUser = history[i].Content
			break
		}
	}

	judgeRes, err := p.judge.Judge(ctx, latestUser, history)
	if err != nil {
		// Fallback to default classification if judge fails
		observability.Logger("gateway.provider").Warn("judge failed, using default coding_main", slog.String("error", err.Error()))
		judgeRes = &JudgeResult{Class: "coding_main", Confidence: 0.5, Reason: "judge error fallback"}
	}

	// Inject "Be concise" only for simple/curation tasks.
	// Professional and other classes get full token budget without truncation pressure.
	const conciseInstruction = "Be concise. Avoid unnecessary explanations. Output minimal sufficient answer."
	switch judgeRes.Class {
	case "simple_short", "curation":
		if !strings.Contains(systemPrompt, conciseInstruction) {
			systemPrompt = conciseInstruction + "\n" + systemPrompt
		}
	}

	req := DryRunRequest{
		Task:            latestUser,
		JudgeClass:      judgeRes.Class,
		JudgeConfidence: judgeRes.Confidence,
		ContextSize:     contextSize,
		RequiresVision:  hasVisionParts(history),
		// Add other flags if needed
	}
	
	if opts.LocalOnly {
		req.LocalOnly = true
	}
	req.OutputMode = opts.OutputMode

	degraded := p.IsDegraded()
	if degraded.Active && degraded.RemoteDisabled {
		req.LocalOnly = true
	}

	candidates := p.planner.Plan(req)

	logger := observability.Logger("gateway.provider")
	logger.Info("gateway routing attempt",
		slog.String("class", judgeRes.Class),
		slog.Float64("confidence", judgeRes.Confidence),
		slog.Int("candidates", len(candidates)),
	)

	var lastErr error
	var previousCandidate *RouteCandidate
	var failureReasons []string

	for _, candidate := range candidates {
		// S-31: Strict Sovereignty — Skip remote models if LocalOnly is requested
		if req.LocalOnly && candidate.UseRemote {
			continue
		}

		if degraded.Active && degraded.RemoteDisabled && candidate.UseRemote {
			continue
		}

		if previousCandidate != nil {
			p.recordFallback(*previousCandidate, candidate)
		}

		resp, err := p.generateWithDecision(ctx, candidate, systemPrompt, history, tools)

		if err != nil {
			failureReasons = append(failureReasons,
				fmt.Sprintf("%s:%s=%s", candidate.Provider, candidate.Model, err.Error()))
			previousCandidate = &candidate
			lastErr = err
			// O gateway continua e tenta o proximo candidato da fila silenciando o erro atual.
			continue
		}

		// Se bateu aqui, retornou sucesso.
		return resp, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all gateway routes failed or were filtered by degraded mode")
	}
	logger.Error("all gateway routes exhausted", slog.String("failures", strings.Join(failureReasons, " | ")))
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
		PrimaryLane: "local",
		PrimaryMode: modelGemma3,
		Budgets:     budgets,
		Routes:      routes,
		Degraded:    p.degraded,
	}
}

func (p *Provider) GatewayStatusJSON() ([]byte, error) {
	return json.Marshal(p.StatusSnapshot())
}

func (p *Provider) HealthCheck() health.CheckResult {
	snapshot := p.StatusSnapshot()
	var warnings []string

	if snapshot.Degraded.Active {
		warnings = append(warnings, "DEGRADED_MODE: "+snapshot.Degraded.Reason)
	}

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
	if p.costExceeded(decision.BudgetLane) {
		p.observeResult(decision, "cost_exceeded", time.Since(startedAt))
		return nil, fmt.Errorf("cost limit exceeded for %s", decision.BudgetLane)
	}

	provider := p.providerFor(decision.Lane)
	if provider == nil {
		p.observeResult(decision, "provider_unavailable", time.Since(startedAt))
		return nil, fmt.Errorf("provider unavailable for %s", providerKey)
	}

	systemPrompt = applyGuards(systemPrompt, decision)
	p.markRequest(decision.BudgetLane, providerKey, decision.Model)

	// Publicar decisão de rota ao dashboard via SSE
	reason := "cost_lane:" + decision.BudgetLane
	if !decision.UseRemote {
		reason = "local"
	}
	dashboard.Publish(dashboard.Event{
		Type:      "route_decision",
		Agent:     "Gateway",
		Action:    fmt.Sprintf("%s → %s:%s", reason, decision.Provider, decision.Model),
		Timestamp: time.Now().Format("15:04:05"),
	})

	// Strip tools for local models that don't support tool calling (e.g. gemma3).
	effectiveTools := tools
	toolsStripped := false
	if !decision.UseRemote && !modelSupportsTools(decision.Model) {
		effectiveTools = nil
		toolsStripped = true
	}

	resp, err := provider.GenerateContent(ctx, systemPrompt, history, effectiveTools)
	if err != nil {
		p.markFailure(decision.BudgetLane, providerKey, err.Error())
		p.observeResult(decision, "error", time.Since(startedAt))
		return nil, err
	}

	// Clean up thought tags for local models that might mix them in.
	if resp != nil {
		resp.Content = stripThoughtTags(resp.Content)
	}
	// Se tools foram stripped para modelo local e a resposta tem conteudo, aceitar sem guard.
	if toolsStripped && resp != nil && strings.TrimSpace(resp.Content) != "" {
		p.markSuccess(providerKey)
		p.observeResult(decision, "success_no_tools", time.Since(startedAt))
		return resp, nil
	}
	if !responseSatisfies(decision, resp) {
		p.markFailure(decision.BudgetLane, providerKey, "empty guarded content")
		p.observeResult(decision, "guard_reject", time.Since(startedAt))
		return resp, fmt.Errorf("empty guarded content for %s", providerKey)
	}
	p.markSuccess(providerKey)
	if resp.InputTokens > 0 || resp.OutputTokens > 0 || len(resp.Metadata) > 0 {
		p.recordTokens(decision.BudgetLane, providerKey, decision.Model, resp.InputTokens, resp.OutputTokens, resp.Metadata)
	}
	p.observeResult(decision, "success", time.Since(startedAt))
	return resp, nil
}

func (p *Provider) providerFor(lane string) agent.LLMProvider {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch lane {
	case "local-fast":
		return p.localFast
	case "local-vision":
		return p.localVision
	case "local-balanced", "local":
		return p.localBalanced
	case "remote-cheap":
		return p.remoteCheapLong
	case "remote-long-context", "remote-cheap-vision":
		return p.remoteCheapVision
	case "remote-premium":
		if p.miniMaxDirect != nil {
			return p.miniMaxDirect
		}
		return p.remotePremium
	default:
		return p.localBalanced
	}
}

func hasVisionParts(history []agent.Message) bool {
	for _, m := range history {
		// This depends on how agent.Message is structured. 
		// Assuming for now it's simple string content check or metadata.
		// If agent.Message has a structured content type, check it here.
		if strings.Contains(m.Content, "[image]") || strings.Contains(m.Content, "base64,") {
			return true
		}
	}
	return false
}

func buildDryRunRequest(systemPrompt string, history []agent.Message, tools []agent.Tool) DryRunRequest {
	var latestUser string
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			latestUser = history[i].Content
			break
		}
	}

	outputMode := detectOutputMode(systemPrompt, latestUser)

	return DryRunRequest{
		Task:           latestUser,
		OutputMode:     outputMode,
		RequiresTools:  len(tools) > 0,
		RequiresVision: containsVisionKeywords(latestUser),
		LocalOnly:      containsLocalOnly(systemPrompt, latestUser),
	}
}

func detectOutputMode(systemPrompt, userMessage string) string {
	systemLower := strings.ToLower(systemPrompt)
	userLower := strings.ToLower(userMessage)
	switch {
	case strings.Contains(systemLower, "json valido") || strings.Contains(systemLower, "only valid json"):
		return "structured_json"
	case strings.Contains(userLower, "3 facts") || strings.Contains(userLower, "3 tags") || strings.Contains(userLower, "facts curtos") || strings.Contains(userLower, "fatos curtos"):
		return "curation"
	default:
		return "text"
	}
}

func containsVisionKeywords(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "screenshot") || strings.Contains(lower, "imagem")
}

func containsLocalOnly(systemPrompt, userMessage string) bool {
	systemLower := strings.ToLower(systemPrompt)
	userLower := strings.ToLower(userMessage)
	return strings.Contains(systemLower, "local only") || strings.Contains(userLower, "local only")
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
	// Se chamou ferramentas, a resposta é válida idependente de texto.
	if len(resp.ToolCalls) > 0 {
		return true
	}
	
	content := strings.TrimSpace(resp.Content)
	reasoning := strings.TrimSpace(resp.ReasoningContent)

	// Se tem conteúdo textual final, é válida.
	if content != "" {
		return true
	}

	// Se o modo é 'minimize', exigimos conteúdo textual final (Content).
	// No entanto, se o modelo for local (Gemma 3), aceitamos o ReasoningContent 
	// como resposta se ele não estiver vazio (às vezes o provider falha em separar).
	if decision.Guards.ReasoningMode == "minimize" && decision.UseRemote {
		return false
	}
	
	// Permitimos que a resposta contenha apenas raciocínio se o modelo for local
	// ou se o modo não for estrito.
	return reasoning != ""
}

// stripThoughtTags removes <thought>...</thought> blocks from the text.
func stripThoughtTags(s string) string {
	for {
		start := strings.Index(s, "<thought>")
		if start == -1 {
			break
		}
		end := strings.Index(s, "</thought>")
		if end == -1 {
			// If tag is not closed, strip until the end of string? 
			// Or just the tag itself. Let's strip the tag and keep the rest.
			return strings.Replace(s, "<thought>", "", 1)
		}
		s = s[:start] + s[end+10:]
	}
	// Also handle the case where it might use markdown-like variations or just the tags
	s = strings.ReplaceAll(s, "<thought>", "")
	s = strings.ReplaceAll(s, "</thought>", "")
	return strings.TrimSpace(s)
}

func (p *Provider) markRequest(budgetLane, providerKey, model string) {
	p.mu.Lock()
	var toPersist []persistableRouteState
	defer p.mu.Unlock()
	state := p.ensureStateLocked(providerKey)
	state.Requests++
	state.LastDecisionModel = model
	toPersist = append(toPersist, persistableRouteState{Key: providerKey, State: *state})
	if budgetLane != "" {
		budgetState := p.ensureStateLocked(budgetLane)
		budgetState.Requests++
		p.updateBudgetMetricLocked(budgetLane, budgetState.Requests)
		toPersist = append(toPersist, persistableRouteState{Key: budgetLane, State: *budgetState})
	}
	p.persistLocked(toPersist...)
}

func (p *Provider) markSuccess(providerKey string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.ensureStateLocked(providerKey)
	state.ConsecutiveFails = 0
	state.BreakerState = "closed"
	state.LastError = ""
	p.setBreakerMetricLocked(providerKey, state.BreakerState)
	p.persistLocked(persistableRouteState{Key: providerKey, State: *state})
}

func (p *Provider) markFailure(budgetLane, providerKey, reason string) {
	p.mu.Lock()
	var toPersist []persistableRouteState
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
	toPersist = append(toPersist, persistableRouteState{Key: providerKey, State: *state})
	if budgetLane != "" {
		budgetState := p.ensureStateLocked(budgetLane)
		budgetState.Failures++
		toPersist = append(toPersist, persistableRouteState{Key: budgetLane, State: *budgetState})
	}
	p.persistLocked(toPersist...)
}

func (p *Provider) recordTokens(budgetLane, providerKey, model string, inputTokens, outputTokens int, metadata map[string]string) {
	cost := CalculateCostUSD(model, inputTokens, outputTokens)
	
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.ensureStateLocked(providerKey)
	state.TotalInputTokens += inputTokens
	state.TotalOutputTokens += outputTokens
	state.TotalCostUSD += cost

	// Parse Rate Limits
	if val, ok := metadata["X-RateLimit-Limit-Requests"]; ok {
		fmt.Sscanf(val, "%d", &state.RateLimitRequests)
	}
	if val, ok := metadata["X-RateLimit-Remaining-Requests"]; ok {
		fmt.Sscanf(val, "%d", &state.RateLimitRemaining)
	}
	if val, ok := metadata["X-RateLimit-Reset-Requests"]; ok {
		// Reset is usually seconds or timestamp. OpenAI uses seconds.
		var seconds float64
		if _, err := fmt.Sscanf(val, "%fs", &seconds); err == nil {
			state.RateLimitReset = time.Now().Add(time.Duration(seconds * float64(time.Second)))
		}
	} else if val, ok := metadata["Retry-After"]; ok {
		var seconds int
		if _, err := fmt.Sscanf(val, "%d", &seconds); err == nil {
			state.RateLimitReset = time.Now().Add(time.Duration(seconds) * time.Second)
		}
	}
	var toPersist []persistableRouteState
	toPersist = append(toPersist, persistableRouteState{Key: providerKey, State: *state})
	if budgetLane != "" {
		budgetState := p.ensureStateLocked(budgetLane)
		budgetState.TotalInputTokens += inputTokens
		budgetState.TotalOutputTokens += outputTokens
		budgetState.TotalCostUSD += cost
		toPersist = append(toPersist, persistableRouteState{Key: budgetLane, State: *budgetState})
	}
	p.persistLocked(toPersist...)
	if p.metrics != nil {
		p.metrics.tokens.WithLabelValues(budgetLane, "input").Add(float64(inputTokens))
		p.metrics.tokens.WithLabelValues(budgetLane, "output").Add(float64(outputTokens))
		p.metrics.costUSD.WithLabelValues(budgetLane).Add(cost)
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
		p.persistLocked(persistableRouteState{Key: providerKey, State: *state})
		return false
	}
	if state.RateLimitRemaining == 0 && !state.RateLimitReset.IsZero() && time.Now().Before(state.RateLimitReset) {
		return true
	}
	p.setBreakerMetricLocked(providerKey, state.BreakerState)
	return true
}

func (p *Provider) costExceeded(lane string) bool {
	if lane == "" {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.ensureStateLocked(lane)
	budget, ok := p.budgets[lane]
	if !ok || budget.CostHardUSD <= 0 {
		return false
	}
	return state.TotalCostUSD >= budget.CostHardUSD
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
	p.rolloverStateLocked(state)
	return state
}

func (p *Provider) rolloverStateLocked(state *routeState) {
	if state == nil {
		return
	}
	today := time.Now().UTC().Format("2006-01-02")
	if state.Day == today {
		return
	}
	state.Day = today
	state.Requests = 0
	state.Failures = 0
	state.ConsecutiveFails = 0
	state.BreakerState = "closed"
	state.BreakerOpenUntil = time.Time{}
	state.LastError = ""
	state.TotalInputTokens = 0
	state.TotalOutputTokens = 0
	state.TotalCostUSD = 0
}

func (p *Provider) loadState() error {
	if p == nil || p.store == nil {
		return nil
	}
	states, err := p.store.Load()
	if err != nil {
		return err
	}
	if len(states) == 0 {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for key, value := range states {
		stateCopy := value
		p.rolloverStateLocked(&stateCopy)
		if stateCopy.BreakerState == "" {
			stateCopy.BreakerState = "closed"
		}
		p.states[key] = &stateCopy
		p.setBreakerMetricLocked(key, stateCopy.BreakerState)
		p.updateBudgetMetricLocked(key, stateCopy.Requests)
	}
	return nil
}

type persistableRouteState struct {
	Key   string
	State routeState
}

func (p *Provider) persistLocked(states ...persistableRouteState) {
	if p == nil || p.store == nil {
		return
	}
	for _, entry := range states {
		if entry.Key == "" {
			continue
		}
		_ = p.store.Save(entry.Key, entry.State)
	}
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

	// Publicar fallback ao dashboard
	dashboard.Publish(dashboard.Event{
		Type:      "route_fallback",
		Agent:     "Gateway",
		Action:    fmt.Sprintf("Fallback: %s:%s → %s:%s", from.Provider, from.Model, to.Provider, to.Model),
		Timestamp: time.Now().Format("15:04:05"),
	})
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
