# ADR 20260328: Anthropic SDK Go Integration

## Status
🟡 Proposto (P2 - Qualidade)

## Contexto
Integrar o **anthropic-sdk-go** (SDK oficial da Anthropic em Go) para computer use agent. O SDK oficial suporta tool definitions e betatoolrunner, permitindo que o agente chame tools de browser via LLM reasoning. Alternativa open source ao SDK Python oficial.

Baseado em: [anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go) (Go SDK oficial)

## Decisões Arquiteturais

### 1. SDK Integration

```go
// internal/llm/anthropic.go
package llm

import (
    "context"

    "github.com/anthropics/anthropic-sdk-go"
)

type AnthropicProvider struct {
    client   *anthropic.Anthropic
    model    string
    tools    []anthropic.ToolParam
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
    return &AnthropicProvider{
        client: anthropic.New(
            anthropic.WithAuthToken(apiKey),
        ),
        model: model,
    }
}

// DefineTool adiciona uma tool de computer use
func (p *AnthropicProvider) DefineTool(tool ToolDef) {
    p.tools = append(p.tools, anthropic.ToolParam{
        Name:        tool.Name,
        Description: anthropic.String(tool.Description),
        InputSchema: tool.Params,
    })
}

// Message sends a message with tools and returns the response
func (p *AnthropicProvider) Message(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
    messages := make([]anthropic.MessageParam, len(req.Messages))
    for i, msg := range req.Messages {
        messages[i] = anthropic.NewUserMessage(anthropic.TextBlockParam{
            Text: msg.Content,
        })
    }

    input := anthropic.MessageInput{
        Model:      p.model,
        Messages:   messages,
        Tools:      p.tools,
        MaxTokens:  4096,
    }

    message, err := p.client.Messages.New(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("anthropic message: %w", err)
    }

    return &MessageResponse{
        Content:   message.Content[0].GetText(),
        ToolCalls: extractToolCalls(message),
    }, nil
}

func extractToolCalls(msg *anthropic.Message) []ToolCall {
    var calls []ToolCall
    for _, block := range msg.Content {
        if tc, ok := block.GetToolUse(); ok {
            calls = append(calls, ToolCall{
                Name:      string(tc.Name),
                Arguments: tc.Input,
            })
        }
    }
    return calls
}
```

### 2. Beta Tool Runner Pattern

```go
// internal/llm/tool_runner.go
package llm

import (
    "context"
    "encoding/json"
)

type ToolRunner struct {
    provider  *AnthropicProvider
    tools     map[string]ToolHandler
    maxSteps  int
}

type ToolHandler func(ctx context.Context, args map[string]any) (string, error)

func (tr *ToolRunner) Register(name string, handler ToolHandler) {
    tr.tools[name] = handler
}

func (tr *ToolRunner) Run(ctx context.Context, req RunRequest) (*RunResult, error) {
    messages := []MessageRequest{{
        Role:    "user",
        Content: req.Task,
    }}

    for step := 0; step < tr.maxSteps; step++ {
        // Envia mensagem e recebe resposta
        resp, err := tr.provider.Message(ctx, MessageRequest{
            Messages: messages,
        })
        if err != nil {
            return nil, err
        }

        // Adiciona resposta do assistant
        messages = append(messages, MessageRequest{
            Role:    "assistant",
            Content: resp.Content,
        })

        // Se não há tool calls, termina
        if len(resp.ToolCalls) == 0 {
            return &RunResult{
                Output:   resp.Content,
                Steps:    len(messages),
            }, nil
        }

        // Executa tool calls
        for _, tc := range resp.ToolCalls {
            handler, ok := tr.tools[tc.Name]
            if !ok {
                messages = append(messages, MessageRequest{
                    Role:    "user",
                    Content: fmt.Sprintf("Error: unknown tool %s", tc.Name),
                })
                continue
            }

            result, err := handler(ctx, tc.Arguments)
            if err != nil {
                result = fmt.Sprintf("Error: %v", err)
            }

            // Adiciona resultado como tool result
            messages = append(messages, MessageRequest{
                Role:    "user",
                Content: formatToolResult(tc.Name, result),
            })
        }
    }

    return &RunResult{
        Output: "Max steps reached",
        Steps:  tr.maxSteps,
    }, nil
}

func formatToolResult(name string, result string) string {
    return fmt.Sprintf(`<tool_result name="%s">%s</tool_result>`, name, result)
}
```

### 3. LiteLLM Integration (Fallback)

```go
// internal/llm/litellm.go
package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
)

// LiteLLMProvider para gemma3:27b e outros modelos locais
type LiteLLMProvider struct {
    baseURL string
    apiKey  string
    model   string
    client  *http.Client
}

type LiteLLMRequest struct {
    Model       string          `json:"model"`
    Messages    []MessagePart  `json:"messages"`
    Temperature float64        `json:"temperature,omitempty"`
    Tools       []ToolDef     `json:"tools,omitempty"`
}

type LiteLLMResponse struct {
    Choices []struct {
        Message struct {
            Content   string `json:"content"`
            ToolCalls []struct {
                ID       string         `json:"id"`
                Type     string         `json:"type"`
                Function struct {
                    Name      string          `json:"name"`
                    Arguments json.RawMessage `json:"arguments"`
                } `json:"function"`
            } `json:"tool_calls"`
        } `json:"message"`
    } `json:"choices"`
    Usage struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
    } `json:"usage"`
}

func (p *LiteLLMProvider) Message(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
    llmReq := LiteLLMRequest{
        Model:       p.model,
        Temperature: 0.7,
        Messages:    p.buildMessages(req),
    }

    body, _ := json.Marshal(llmReq)

    httpReq, _ := http.NewRequestWithContext(ctx, "POST",
        p.baseURL+"/chat/completions",
        bytes.NewReader(body))
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

    resp, err := p.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("litellm request: %w", err)
    }
    defer resp.Body.Close()

    var llmResp LiteLLMResponse
    if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
        return nil, fmt.Errorf("litellm decode: %w", err)
    }

    choice := llmResp.Choices[0].Message

    var toolCalls []ToolCall
    for _, tc := range choice.ToolCalls {
        var args map[string]any
        json.Unmarshal(tc.Function.Arguments, &args)
        toolCalls = append(toolCalls, ToolCall{
            Name:      tc.Function.Name,
            Arguments: args,
        })
    }

    return &MessageResponse{
        Content:   choice.Content,
        ToolCalls: toolCalls,
    }, nil
}

func (p *LiteLLMProvider) buildMessages(req MessageRequest) []MessagePart {
    var msgs []MessagePart
    for _, m := range req.Messages {
        msg := MessagePart{
            Role:    m.Role,
            Content: m.Content,
        }
        msgs = append(msgs, msg)
    }
    return msgs
}
```

### 4. Provider Fallback Chain

```go
// internal/llm/provider.go
package llm

type Provider interface {
    Message(ctx context.Context, req MessageRequest) (*MessageResponse, error)
    DefineTool(tool ToolDef)
}

// ProviderChain tenta provedores em ordem até um funcionar
type ProviderChain struct {
    providers []Provider
}

func NewProviderChain(providers ...Provider) *ProviderChain {
    return &ProviderChain{providers: providers}
}

func (c *ProviderChain) Message(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
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

// Common chain: Claude → gemma3 via LiteLLM → OpenRouter
func DefaultProviderChain() *ProviderChain {
    return NewProviderChain(
        NewAnthropicProvider(os.Getenv("ANTHROPIC_API_KEY"), "claude-sonnet-4-20250514"),
        NewLiteLLMProvider(
            os.Getenv("LITELLM_URL"),
            os.Getenv("LITELLM_API_KEY"),
            "aurelia-smart", // gemma3:27b
        ),
    )
}
```

## Consequências

### Positivas
- SDK oficial Anthropic = melhor suporte e updates
- betatoolrunner pattern é bem testado
- Fallback chain garante resiliência
- LiteLLM suporta gemma3:27b e outros modelos locais

### Negativas
- Claude API é paid (custo por token)
- LiteLLM pode ter latência adicional
- Múltiplos providers = complexidade

### Trade-offs
- Claude vs gemma3: Claude é mais inteligente, gemma3 é local/sovereign
- Streaming vs non-streaming: streaming é melhor UX mas mais complexo

## Dependências
- ⚠️ `github.com/anthropics/anthropic-sdk-go` - SDK oficial
- ⚠️ `internal/llm/` - existente via gateway
- ❌ `configs/litellm/config.yaml` - ADR separado

## Referências
- [anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [ADR-20260328-smart-router-litellm-gemma3-27b.md](./20260328-smart-router-litellm-gemma3-27b.md)
- [ADR-20260328-bua-browser-use-agent-go.md](./20260328-bua-browser-use-agent-go.md)

## Links Obrigatórios
- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/neon-sentinel
**Progress**: 0%
