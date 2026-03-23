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
	"qwen3.5:4b": {InputPerMToken: 0, OutputPerMToken: 0},
	"qwen3.5:9b": {InputPerMToken: 0, OutputPerMToken: 0},

	// Tier 1 — Remote Cheap
	"qwen/qwen3.5-flash-02-23": {InputPerMToken: 0.05, OutputPerMToken: 0.10},

	// Tier 2 — Remote Structured
	"google/gemini-2.5-flash": {InputPerMToken: 0.10, OutputPerMToken: 0.40},

	// Tier 3 — Remote Premium
	"minimax/minimax-m2.7": {InputPerMToken: 0.30, OutputPerMToken: 0.60},

	// Tier 4 — Emergency
	"anthropic/claude-haiku-4.5": {InputPerMToken: 0.80, OutputPerMToken: 4.00},

	// Vision
	"qwen/qwen3.5-9b": {InputPerMToken: 0.05, OutputPerMToken: 0.10},

	// Audio (Groq STT)
	"whisper-large-v3-turbo": {InputPerMToken: 0, OutputPerMToken: 0},

	// Groq Text
	"llama-3.3-70b-versatile": {InputPerMToken: 0.59, OutputPerMToken: 0.79},
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
