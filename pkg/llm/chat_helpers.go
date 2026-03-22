package llm

import (
	"encoding/json"
	"fmt"

	"github.com/kocar/aurelia/internal/agent"
)

func parseChatCompletionResponse(respBody []byte) (*agent.ModelResponse, error) {
	var apiResp chatCompletionResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chat completion response: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned by chat completion API")
	}

	choice := apiResp.Choices[0].Message
	result := &agent.ModelResponse{
		Content:          choice.Content,
		ReasoningContent: choice.ReasoningContent,
	}
	if apiResp.Usage != nil {
		result.InputTokens = apiResp.Usage.PromptTokens
		result.OutputTokens = apiResp.Usage.CompletionTokens
	}

	for _, call := range choice.ToolCalls {
		if call.Type != "function" {
			continue
		}

		var args map[string]any
		_ = json.Unmarshal([]byte(call.Function.Arguments), &args)

		result.ToolCalls = append(result.ToolCalls, agent.ToolCall{
			ID:        call.ID,
			Name:      call.Function.Name,
			Arguments: args,
		})
	}

	return result, nil
}

func mapChatRole(role string) string {
	switch role {
	case "assistant":
		return "assistant"
	case "tool":
		return "tool"
	default:
		return "user"
	}
}

func mapToolCalls(calls []agent.ToolCall) []chatToolCall {
	toolCalls := make([]chatToolCall, 0, len(calls))
	for _, call := range calls {
		argsRaw, _ := json.Marshal(call.Arguments)
		toolCalls = append(toolCalls, chatToolCall{
			ID:   call.ID,
			Type: "function",
			Function: chatFunctionCall{
				Name:      call.Name,
				Arguments: string(argsRaw),
			},
		})
	}
	return toolCalls
}

func buildChatTools(tools []agent.Tool) ([]chatTool, error) {
	chatTools := make([]chatTool, 0, len(tools))
	for _, t := range tools {
		schemaRaw, err := json.Marshal(t.JSONSchema)
		if err != nil {
			return nil, fmt.Errorf("marshal tool schema %s: %w", t.Name, err)
		}
		if string(schemaRaw) == "null" {
			schemaRaw = []byte(`{"type":"object","properties":{}}`)
		}

		chatTools = append(chatTools, chatTool{
			Type: "function",
			Function: chatToolDef{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  json.RawMessage(schemaRaw),
			},
		})
	}
	return chatTools, nil
}
