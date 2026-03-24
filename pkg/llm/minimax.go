package llm

import (
)

const (
	MiniMaxInternationalBaseURL = "https://api.minimax.io/v1"
	// MiniMaxChinaBaseURL        = "https://api.minimaxi.com/v1"
)

// NewMiniMaxProvider cria um novo provedor para MiniMax usando o adaptador OpenAI-compatible
func NewMiniMaxProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewMiniMaxProviderWithOptions(apiKey, model, OpenAICompatibleRequestOptions{})
}

// NewMiniMaxProviderWithOptions cria um novo provedor para MiniMax com opções customizadas
func NewMiniMaxProviderWithOptions(apiKey string, model string, request OpenAICompatibleRequestOptions) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Provider: "minimax",
		APIKey:   apiKey,
		BaseURL:  MiniMaxInternationalBaseURL,
		Model:    model,
		Request:  request,
	})
}
