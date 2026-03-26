package gateway

import (
	"testing"
	"time"
)

func TestCalculateCostUSD_GroqLlama(t *testing.T) {
	t.Parallel()
	// Groq Llama-3.3-70B is free tier (direct API, no per-token charge)
	cost := CalculateCostUSD("llama-3.3-70b-versatile", 1000, 500)
	if cost != 0 {
		t.Fatalf("cost = %f, want 0 for groq free tier", cost)
	}
}

func TestCalculateCostUSD_MiniMaxM2_7(t *testing.T) {
	t.Parallel()
	// 1000 input tokens + 500 output tokens on minimax/minimax-m2.7
	// Pricing: $0.30/$1.20 per 1M tokens
	cost := CalculateCostUSD("minimax/minimax-m2.7", 1000, 500)
	// Expected: (1000 * 0.30 / 1M) + (500 * 1.20 / 1M) = 0.0003 + 0.0006 = 0.0009
	if cost < 0.0008 || cost > 0.001 {
		t.Fatalf("cost = %f, want ~0.0009", cost)
	}
}

func TestCalculateCostUSD_LocalModelFree(t *testing.T) {
	t.Parallel()
	cost := CalculateCostUSD("gemma3:12b", 10000, 5000)
	if cost != 0 {
		t.Fatalf("cost = %f, want 0 for local model", cost)
	}
}

func TestCalculateCostUSD_UnknownModelFree(t *testing.T) {
	t.Parallel()
	cost := CalculateCostUSD("unknown-model", 1000, 500)
	if cost != 0 {
		t.Fatalf("cost = %f, want 0 for unknown model", cost)
	}
}

func TestLookupCost_KnownModel(t *testing.T) {
	t.Parallel()
	// Tier 1 — DeepSeek V3.1 cheap remote fallback
	cost := LookupCost("deepseek/deepseek-chat-v3.1")
	if cost.InputPerMToken != 0.15 || cost.OutputPerMToken != 0.75 {
		t.Fatalf("cost = %+v", cost)
	}
}

func TestCostExceeded_BlocksWhenOverLimit(t *testing.T) {
	t.Parallel()
	p := &Provider{
		budgets: map[string]laneBudget{
			"remote_premium": {Soft: 80, Hard: 160, CostHardUSD: 2.00},
		},
		states: map[string]*routeState{
			"remote_premium": {Day: time.Now().UTC().Format("2006-01-02"), TotalCostUSD: 2.50},
		},
	}
	if !p.costExceeded("remote_premium") {
		t.Fatal("expected cost exceeded")
	}
}

func TestCostExceeded_AllowsUnderLimit(t *testing.T) {
	t.Parallel()
	p := &Provider{
		budgets: map[string]laneBudget{
			"remote_premium": {Soft: 80, Hard: 160, CostHardUSD: 2.00},
		},
		states: map[string]*routeState{
			"remote_premium": {Day: time.Now().UTC().Format("2006-01-02"), TotalCostUSD: 0.50},
		},
	}
	if p.costExceeded("remote_premium") {
		t.Fatal("expected cost not exceeded")
	}
}
