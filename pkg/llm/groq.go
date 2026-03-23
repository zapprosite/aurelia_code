package llm

const groqChatCompletionsURL = "https://api.groq.com/openai/v1/chat/completions"

func NewGroqProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewGroqProviderWithOptions(apiKey, model, OpenAICompatibleRequestOptions{})
}

func NewGroqProviderWithOptions(apiKey, model string, request OpenAICompatibleRequestOptions) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Provider:  "groq",
		APIKey:    apiKey,
		BaseURL:   groqChatCompletionsURL,
		Model:     model,
		UserAgent: "Aurelia/1.0",
		Request:   request,
	})
}
