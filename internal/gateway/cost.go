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
	"gemma3:12b":          {InputPerMToken: 0, OutputPerMToken: 0},
	"gemma3:27b-it-q4_K_M": {InputPerMToken: 0, OutputPerMToken: 0},

	// Tier 1 — Remote Cheap (OpenRouter)
	"deepseek/deepseek-chat-v3.1": {InputPerMToken: 0.014, OutputPerMToken: 0.14},

	// Tier 2 — Remote Premium
	"minimax/minimax-m2.7": {InputPerMToken: 0.30, OutputPerMToken: 0.60},

	// Tier 2 — Long Context / Multimodal (OpenRouter)
	"moonshotai/kimi-k2.5": {InputPerMToken: 0.14, OutputPerMToken: 0.60},

	// Tier 4 — Emergency
	"anthropic/claude-haiku-4.5": {InputPerMToken: 0.80, OutputPerMToken: 4.00},

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
