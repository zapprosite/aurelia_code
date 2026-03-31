package llm

import (
	"context"

	"github.com/kocar/aurelia/internal/agent"
)

// GeminiProvider implements agent.LLMProvider for Google Gemini models using its OpenAI-compatible endpoint.
// For SOTA 2026, we use the standard OpenAI-compatible base structure to ensure uniform integration.
type GeminiProvider struct {
	*OpenAICompatibleProvider
}

// NewGeminiProvider creates a new LLM provider for Gemini models.
func NewGeminiProvider(apiKey string, modelName string) *GeminiProvider {
	baseURL := "https://generativelanguage.googleapis.com/v1beta/openai/"
	
	cfg := OpenAICompatibleConfig{
		Provider: "gemini",
		APIKey:   apiKey,
		BaseURL:  baseURL,
		Model:    modelName,
		Request: OpenAICompatibleRequestOptions{
			MaxTokens: 4096,
		},
	}
	
	// Gemini OpenAI-compatible endpoint uses the API key in a specific way sometimes, 
	// but standard Bearer token works for /v1beta/openai/
	return &GeminiProvider{
		OpenAICompatibleProvider: NewOpenAICompatibleProvider(cfg),
	}
}

func (p *GeminiProvider) GenerateContent(
	ctx context.Context,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (*agent.ModelResponse, error) {
	return p.OpenAICompatibleProvider.GenerateContent(ctx, systemPrompt, history, tools)
}

func (p *GeminiProvider) GenerateStream(
	ctx context.Context,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (<-chan agent.StreamResponse, error) {
	return p.OpenAICompatibleProvider.GenerateStream(ctx, systemPrompt, history, tools)
}
