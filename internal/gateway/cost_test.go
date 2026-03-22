package gateway

import "testing"

func TestCalculateCostUSD_MiniMaxM2_7(t *testing.T) {
	t.Parallel()
	// 1000 input tokens + 500 output tokens on minimax/minimax-m2.7
	cost := CalculateCostUSD("minimax/minimax-m2.7", 1000, 500)
	// Expected: (1000 * 0.30 / 1M) + (500 * 0.60 / 1M) = 0.0003 + 0.0003 = 0.0006
	if cost < 0.0005 || cost > 0.0007 {
		t.Fatalf("cost = %f, want ~0.0006", cost)
	}
}

func TestCalculateCostUSD_LocalModelFree(t *testing.T) {
	t.Parallel()
	cost := CalculateCostUSD("qwen3.5:9b", 10000, 5000)
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
	cost := LookupCost("google/gemini-2.5-flash")
	if cost.InputPerMToken != 0.10 || cost.OutputPerMToken != 0.40 {
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
			"remote_premium": {Day: "2026-03-22", TotalCostUSD: 2.50},
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
			"remote_premium": {Day: "2026-03-22", TotalCostUSD: 0.50},
		},
	}
	if p.costExceeded("remote_premium") {
		t.Fatal("expected cost not exceeded")
	}
}
