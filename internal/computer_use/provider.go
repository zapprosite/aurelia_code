// Package computer_use provides autonomous computer use agent.
package computer_use

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Provider defines an interface for LLM providers.
type Provider interface {
	Message(ctx context.Context, req *MessageRequest) (*MessageResponse, error)
	DefineTool(name, description string, schema map[string]interface{})
}

// MessageRequest represents a chat message request.
type MessageRequest struct {
	Role    string
	Content string
}

// MessageResponse represents a chat message response.
type MessageResponse struct {
	Content   string
	ToolCalls []ToolCall
}

// ToolCall represents a tool call from the LLM.
type ToolCall struct {
	Name      string
	Arguments map[string]interface{}
}

// AnthropicProvider uses the official Anthropic Go SDK.
type AnthropicProvider struct {
	apiKey  string
	model   string
	tools   []Tool
}

// Tool represents an LLM tool definition.
type Tool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
}

// NewAnthropicProvider creates a new Anthropic provider.
func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey: apiKey,
		model:  model,
		tools:  make([]Tool, 0),
	}
}

// DefineTool adds a tool definition.
func (p *AnthropicProvider) DefineTool(name, description string, schema map[string]interface{}) {
	p.tools = append(p.tools, Tool{
		Name:        name,
		Description: description,
		InputSchema: schema,
	})
}

// Message sends a message to the LLM.
func (p *AnthropicProvider) Message(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	// Build request for Anthropic API
	body := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]interface{}{
			{"role": req.Role, "content": req.Content},
		},
		"max_tokens": 4096,
	}

	if len(p.tools) > 0 {
		body["tools"] = p.tools
	}

	// Make request
	jsonBody, _ := json.Marshal(body)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST",
		"https://api.anthropic.com/v1/messages",
		strings.NewReader(string(jsonBody)))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic request: %w", err)
	}
	defer resp.Body.Close()

	var anthropicResp struct {
		Content []struct {
			Type   string `json:"type"`
			Text   string `json:"text"`
			ID     string `json:"id"`
			Name   string `json:"name"`
			Input  map[string]interface{} `json:"input"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	result := &MessageResponse{}

	for _, c := range anthropicResp.Content {
		if c.Type == "text" {
			result.Content += c.Text
		} else if c.Type == "tool_use" {
			result.ToolCalls = append(result.ToolCalls, ToolCall{
				Name:      c.Name,
				Arguments: c.Input,
			})
		}
	}

	return result, nil
}

// LiteLLMProvider uses LiteLLM for local models.
type LiteLLMProvider struct {
	baseURL string
	apiKey  string
	model   string
	tools   []Tool
}

// NewLiteLLMProvider creates a new LiteLLM provider.
func NewLiteLLMProvider(baseURL, apiKey, model string) *LiteLLMProvider {
	return &LiteLLMProvider{
		baseURL: baseURL,
		apiKey: apiKey,
		model:  model,
		tools:  make([]Tool, 0),
	}
}

// DefineTool adds a tool definition.
func (p *LiteLLMProvider) DefineTool(name, description string, schema map[string]interface{}) {
	p.tools = append(p.tools, Tool{
		Name:        name,
		Description: description,
		InputSchema: schema,
	})
}

// Message sends a message to the LLM.
func (p *LiteLLMProvider) Message(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	messages := []map[string]interface{}{
		{"role": req.Role, "content": req.Content},
	}

	body := map[string]interface{}{
		"model": p.model,
		"messages": messages,
	}

	if len(p.tools) > 0 {
		// Convert tools to OpenAI format
		openaiTools := make([]map[string]interface{}, len(p.tools))
		for i, t := range p.tools {
			openaiTools[i] = map[string]interface{}{
				"type": "function",
				"function": map[string]interface{}{
					"name":        t.Name,
					"description": t.Description,
					"parameters":  t.InputSchema,
				},
			}
		}
		body["tools"] = openaiTools
	}

	jsonBody, _ := json.Marshal(body)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST",
		p.baseURL+"/chat/completions",
		strings.NewReader(string(jsonBody)))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("litellm request: %w", err)
	}
	defer resp.Body.Close()

	var litellmResp struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
				ToolCalls []struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&litellmResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(litellmResp.Choices) == 0 {
		return &MessageResponse{}, nil
	}

	result := &MessageResponse{
		Content: litellmResp.Choices[0].Message.Content,
	}

	for _, tc := range litellmResp.Choices[0].Message.ToolCalls {
		var args map[string]interface{}
		json.Unmarshal([]byte(tc.Function.Arguments), &args)
		result.ToolCalls = append(result.ToolCalls, ToolCall{
			Name:      tc.Function.Name,
			Arguments: args,
		})
	}

	return result, nil
}

// ProviderChain tries providers in order until one works.
type ProviderChain struct {
	providers []Provider
}

// NewProviderChain creates a new provider chain.
func NewProviderChain(providers ...Provider) *ProviderChain {
	return &ProviderChain{providers: providers}
}

// Message tries each provider in order.
func (c *ProviderChain) Message(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	var lastErr error
	for _, p := range c.providers {
		resp, err := p.Message(ctx, req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all providers failed: %w", lastErr)
}

// DefineTool adds a tool to all providers.
func (c *ProviderChain) DefineTool(name, description string, schema map[string]interface{}) {
	for _, p := range c.providers {
		p.DefineTool(name, description, schema)
	}
}

// ToolRunner executes tool calls from LLM responses.
type ToolRunner struct {
	provider    Provider
	toolHandlers map[string]ToolHandlerFunc
	maxSteps    int
}

// ToolHandlerFunc is a function that handles a tool call.
type ToolHandlerFunc func(args map[string]interface{}) (string, error)

// NewToolRunner creates a new tool runner.
func NewToolRunner(provider Provider, maxSteps int) *ToolRunner {
	return &ToolRunner{
		provider:    provider,
		toolHandlers: make(map[string]ToolHandlerFunc),
		maxSteps:    maxSteps,
	}
}

// Register registers a tool handler.
func (tr *ToolRunner) Register(name string, handler ToolHandlerFunc) {
	tr.toolHandlers[name] = handler
}

// Run executes the tool runner loop.
func (tr *ToolRunner) Run(ctx context.Context, task string) (string, error) {
	messages := []*MessageRequest{
		{Role: "user", Content: task},
	}

	for step := 0; step < tr.maxSteps; step++ {
		resp, err := tr.provider.Message(ctx, messages[len(messages)-1])
		if err != nil {
			return "", fmt.Errorf("provider message: %w", err)
		}

		// Add assistant response
		messages = append(messages, &MessageRequest{
			Role:    "assistant",
			Content: resp.Content,
		})

		// If no tool calls, we're done
		if len(resp.ToolCalls) == 0 {
			return resp.Content, nil
		}

		// Execute tool calls
		for _, tc := range resp.ToolCalls {
			handler, ok := tr.toolHandlers[tc.Name]
			if !ok {
				messages = append(messages, &MessageRequest{
					Role:    "user",
					Content: fmt.Sprintf("Error: unknown tool %s", tc.Name),
				})
				continue
			}

			result, err := handler(tc.Arguments)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			messages = append(messages, &MessageRequest{
				Role:    "user",
				Content: fmt.Sprintf("<tool_result name=\"%s\">%s</tool_result>", tc.Name, result),
			})
		}
	}

	return "Max steps reached", nil
}
