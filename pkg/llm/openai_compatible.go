package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

const openAICompatibleTimeout = 8 * time.Minute

type OpenAICompatibleConfig struct {
	Provider   string
	APIKey     string
	BaseURL    string
	Model      string
	UserAgent  string
	Headers    map[string]string
	HTTPClient *http.Client
	Request    OpenAICompatibleRequestOptions
}

type OpenAICompatibleRequestOptions struct {
	MaxTokens   int
	Temperature *float64
	ExtraFields map[string]any
}

// OpenAICompatibleProvider implements agent.LLMProvider for chat-completions APIs
// that follow the OpenAI-compatible request and response shape.
type OpenAICompatibleProvider struct {
	client    *http.Client
	provider  string
	apiKey    string
	baseURL   string
	model     string
	userAgent string
	headers   map[string]string
	request   OpenAICompatibleRequestOptions
}

func NewOpenAICompatibleProvider(cfg OpenAICompatibleConfig) *OpenAICompatibleProvider {
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: openAICompatibleTimeout}
	}

	return &OpenAICompatibleProvider{
		client:    client,
		provider:  cfg.Provider,
		apiKey:    cfg.APIKey,
		baseURL:   cfg.BaseURL,
		model:     cfg.Model,
		userAgent: cfg.UserAgent,
		headers:   cloneStringMap(cfg.Headers),
		request: OpenAICompatibleRequestOptions{
			MaxTokens:   cfg.Request.MaxTokens,
			Temperature: cfg.Request.Temperature,
			ExtraFields: cloneAnyMap(cfg.Request.ExtraFields),
		},
	}
}

func (p *OpenAICompatibleProvider) Close() {}

func (p *OpenAICompatibleProvider) GenerateContent(
	ctx context.Context,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (*agent.ModelResponse, error) {
	if err := ensureVisionSupport(p.provider, p.model, history); err != nil {
		return nil, err
	}
	reqBody, err := buildOpenAICompatibleRequest(p.model, systemPrompt, history, tools, p.request)

	if err != nil {
		return nil, err
	}

	respBody, headers, err := p.doChatCompletionRequest(ctx, reqBody)
	if err != nil {
		return nil, err
	}

	return parseChatCompletionResponse(respBody, headers)
}

func buildOpenAICompatibleRequest(
	model string,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
	opts OpenAICompatibleRequestOptions,
) (map[string]any, error) {
	messages := []map[string]any{{
		"role":    "system",
		"content": systemPrompt,
	}}
	messages = append(messages, buildOpenAICompatibleHistory(history)...)

	reqBody := map[string]any{
		"model":    model,
		"messages": messages,
	}
	if opts.MaxTokens > 0 {
		reqBody["max_tokens"] = opts.MaxTokens
	}
	if opts.Temperature != nil {
		reqBody["temperature"] = *opts.Temperature
	}
	for key, value := range opts.ExtraFields {
		reqBody[key] = value
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

func buildOpenAICompatibleHistory(history []agent.Message) []map[string]any {
	messages := make([]map[string]any, 0, len(history))
	for _, msg := range history {
		messages = append(messages, mapOpenAICompatibleMessage(msg))
	}
	return messages
}

func mapOpenAICompatibleMessage(msg agent.Message) map[string]any {
	cMsg := map[string]any{
		"role": mapChatRole(msg.Role),
	}

	if msg.Role == "tool" {
		cMsg["content"] = msg.Content
		cMsg["tool_call_id"] = msg.ToolCallID
		cMsg["name"] = msg.ToolCallID
		return cMsg
	}

	if len(msg.Parts) != 0 {
		cMsg["content"] = buildOpenAICompatibleContent(msg.Parts)
	} else {
		cMsg["content"] = msg.Content
	}

	if msg.Role == "assistant" && len(msg.ToolCalls) > 0 {
		cMsg["reasoning_content"] = msg.ReasoningContent
		cMsg["tool_calls"] = mapToolCalls(msg.ToolCalls)
	}

	return cMsg
}

func buildOpenAICompatibleContent(parts []agent.ContentPart) []map[string]any {
	content := make([]map[string]any, 0, len(parts))
	for _, part := range parts {
		switch part.Type {
		case agent.ContentPartImage:
			content = append(content, map[string]any{
				"type": "image_url",
				"image_url": map[string]any{
					"url": openAIImageURL(part),
				},
			})
		default:
			content = append(content, map[string]any{
				"type": "text",
				"text": part.Text,
			})
		}
	}
	return content
}

func (p *OpenAICompatibleProvider) doChatCompletionRequest(ctx context.Context, reqBody map[string]any) ([]byte, http.Header, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create http request: %w", err)
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
		return nil, nil, fmt.Errorf("http request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if isVisionUnsupportedAPIResponse(resp.StatusCode, respBody) {
			return nil, nil, VisionUnsupportedError{provider: p.provider, model: p.model}
		}
		return nil, nil, fmt.Errorf("openai-compatible API error: %d, response: %s", resp.StatusCode, string(respBody))
	}
	return respBody, resp.Header, nil
}

func isVisionUnsupportedAPIResponse(statusCode int, body []byte) bool {
	if statusCode != http.StatusNotFound && statusCode != http.StatusBadRequest && statusCode != http.StatusUnprocessableEntity {
		return false
	}
	lower := strings.ToLower(string(body))
 	patterns := []string{
 		"no endpoints found that support image input",
 		"does not support image input",
		"vision is not supported",
		"invalid model with image",
 		"image input",
 		"vision",
 	}
	for _, pattern := range patterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
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

func cloneAnyMap(input map[string]any) map[string]any {
	if len(input) == 0 {
		return nil
	}
	cloned := make(map[string]any, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
