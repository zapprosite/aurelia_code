package gateway

// ModelCost holds pricing per million tokens for a model.
type ModelCost struct {
	InputPerMToken  float64
	OutputPerMToken float64
}

// modelCosts maps model identifiers to their pricing.
// Prices sourced from OpenRouter/provider pricing pages.
var modelCosts = map[string]ModelCost{
	// Tier 0 — Local (free)
	"qwen3.5": {InputPerMToken: 0, OutputPerMToken: 0},

	// Tier 1 — Remote Cheap (OpenRouter)
	"mistralai/devstral-2512":     {InputPerMToken: 0.05, OutputPerMToken: 0.22},  // cheapest paid Tier1
	"deepseek/deepseek-v3.2":      {InputPerMToken: 0.28, OutputPerMToken: 0.42},
	"deepseek/deepseek-chat-v3.1": {InputPerMToken: 0.15, OutputPerMToken: 0.75},

	// Tier 1 — Free (rate-limited, 1000 req/day with $10+ credit)
	"deepseek/deepseek-r1-0528:free": {InputPerMToken: 0, OutputPerMToken: 0},

	// Tier 2 — Remote Premium
	"minimax/minimax-m2.7": {InputPerMToken: 0.30, OutputPerMToken: 1.20},
	"minimax/minimax-m2.5": {InputPerMToken: 0.20, OutputPerMToken: 1.17},

	// Long Context (OpenRouter)
	"meta-llama/llama-4-scout": {InputPerMToken: 0.08, OutputPerMToken: 0.30},  // 10M ctx
	"moonshotai/kimi-k2.5":     {InputPerMToken: 0.50, OutputPerMToken: 2.00},  // vision fallback

	// Tier 4 — Emergency
	"anthropic/claude-haiku-4.5": {InputPerMToken: 0.80, OutputPerMToken: 4.00},

	// Audio (Groq STT)
	"whisper-large-v3-turbo": {InputPerMToken: 0, OutputPerMToken: 0},

	// Groq Text (direct API — free tier, $0 cost)
	"llama-3.3-70b-versatile": {InputPerMToken: 0, OutputPerMToken: 0},
}

// LookupCost returns the cost for a model. Returns zero cost for unknown models.
func LookupCost(model string) ModelCost {
	if cost, ok := modelCosts[model]; ok {
		return cost
	}
	return ModelCost{}
}

// CalculateCostUSD computes USD cost from token counts and model pricing.
func CalculateCostUSD(model string, inputTokens, outputTokens int) float64 {
	cost := LookupCost(model)
	return (float64(inputTokens) * cost.InputPerMToken / 1_000_000) +
		(float64(outputTokens) * cost.OutputPerMToken / 1_000_000)
}
