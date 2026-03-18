package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

const openAICompatibleTimeout = 8 * time.Minute

type OpenAICompatibleConfig struct {
	APIKey     string
	BaseURL    string
	Model      string
	UserAgent  string
	Headers    map[string]string
	HTTPClient *http.Client
}

// OpenAICompatibleProvider implements agent.LLMProvider for chat-completions APIs
// that follow the OpenAI-compatible request and response shape.
type OpenAICompatibleProvider struct {
	client    *http.Client
	apiKey    string
	baseURL   string
	model     string
	userAgent string
	headers   map[string]string
}

func NewOpenAICompatibleProvider(cfg OpenAICompatibleConfig) *OpenAICompatibleProvider {
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: openAICompatibleTimeout}
	}

	return &OpenAICompatibleProvider{
		client:    client,
		apiKey:    cfg.APIKey,
		baseURL:   cfg.BaseURL,
		model:     cfg.Model,
		userAgent: cfg.UserAgent,
		headers:   cloneStringMap(cfg.Headers),
	}
}

func (p *OpenAICompatibleProvider) Close() {}

func (p *OpenAICompatibleProvider) GenerateContent(
	ctx context.Context,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (*agent.ModelResponse, error) {
	reqBody, err := buildOpenAICompatibleRequest(p.model, systemPrompt, history, tools)
	if err != nil {
		return nil, err
	}

	respBody, err := p.doChatCompletionRequest(ctx, reqBody)
	if err != nil {
		return nil, err
	}

	return parseChatCompletionResponse(respBody)
}

func buildOpenAICompatibleRequest(
	model string,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (map[string]any, error) {
	messages := append([]chatMessage{{
		Role:    "system",
		Content: systemPrompt,
	}}, buildChatHistory(history)...)

	reqBody := map[string]any{
		"model":    model,
		"messages": messages,
	}

	if len(tools) > 0 {
		chatTools, err := buildChatTools(tools)
		if err != nil {
			return nil, err
		}
		reqBody["tools"] = chatTools
	}

	return reqBody, nil
}

func (p *OpenAICompatibleProvider) doChatCompletionRequest(ctx context.Context, reqBody map[string]any) ([]byte, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}
	if p.userAgent != "" {
		req.Header.Set("User-Agent", p.userAgent)
	}
	for key, value := range p.headers {
		req.Header.Set(key, value)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai-compatible API error: %d, response: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
