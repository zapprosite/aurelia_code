package llm

const zAIChatCompletionsURL = "https://api.z.ai/api/coding/paas/v4/chat/completions"

func NewZAIProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		APIKey:    apiKey,
		BaseURL:   zAIChatCompletionsURL,
		Model:     model,
		UserAgent: "Aurelia/1.0",
	})
}
