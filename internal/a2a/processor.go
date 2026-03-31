package a2a

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kocar/aurelia/internal/mcp"
)

// AureliaProcessor implements MessageProcessor for the Aurelia agent.
// It routes incoming A2A messages to the appropriate subsystem (MCP tools, memory, etc.).
type AureliaProcessor struct {
	mcpManager *mcp.Manager
	log        *slog.Logger
}

// NewAureliaProcessor creates a processor backed by an MCP manager.
func NewAureliaProcessor(mgr *mcp.Manager) *AureliaProcessor {
	return &AureliaProcessor{
		mcpManager: mgr,
		log:        slog.Default(),
	}
}

// ProcessMessage routes the incoming message to the appropriate handler.
func (p *AureliaProcessor) ProcessMessage(ctx context.Context, msg Message, opts ProcessOptions) (*MessageProcessingResult, error) {
	text := extractText(msg)
	if text == "" {
		return &MessageProcessingResult{
			Content: []MessagePart{TextPart{Text: "A2A message has no text content."}},
		}, nil
	}

	p.log.Info("A2A message received", slog.String("role", string(msg.Role)), slog.String("text_preview", truncate(text, 80)))

	// Route by prefix
	switch {
	case startsWith(text, "/tool ", "/mcp ", "/call "):
		return p.handleToolCall(ctx, text, opts)
	case startsWith(text, "/memory ", "/remember ", "/recall "):
		return p.handleMemory(ctx, text, opts)
	case startsWith(text, "/status", "/ping"):
		return &MessageProcessingResult{
			Content: []MessagePart{TextPart{Text: "ok"}},
		}, nil
	default:
		// Forward as general query
		return p.handleGeneral(ctx, text, opts)
	}
}

func (p *AureliaProcessor) handleToolCall(ctx context.Context, text string, opts ProcessOptions) (*MessageProcessingResult, error) {
	// Parse: /tool serverName toolName [args JSON]
	// Example: /tool os-controller run_bash_command {"command":"ls -la"}
	remaining := stripPrefix(text, "/tool ", "/mcp ", "/call ")
	parts := splitSpace(remaining)
	if len(parts) < 2 {
		return &MessageProcessingResult{
			Content: []MessagePart{TextPart{Text: "usage: /tool <server> <toolname> [args_json]"}},
		}, nil
	}
	serverName := parts[0]
	toolName := parts[1]
	argsStr := ""
	if len(parts) >= 3 {
		argsStr = parts[2]
	}

	var args map[string]interface{}
	if argsStr != "" {
		// args is JSON
		if err := parseJSON(argsStr, &args); err != nil {
			return &MessageProcessingResult{
				Content: []MessagePart{TextPart{Text: fmt.Sprintf("invalid JSON args: %v", err)}},
			}, nil
		}
	}

	if p.mcpManager == nil {
		return &MessageProcessingResult{
			Content: []MessagePart{TextPart{Text: "MCP manager not available"}},
		}, nil
	}

	result, err := p.mcpManager.CallTool(ctx, serverName, toolName, args)
	if err != nil {
		return &MessageProcessingResult{
			Content: []MessagePart{TextPart{Text: fmt.Sprintf("MCP call failed: %v", err)}},
		}, nil
	}

	content := result.Content
	if result.IsError {
		content = "error: " + content
	}

	return &MessageProcessingResult{
		Content: []MessagePart{TextPart{Text: content}},
	}, nil
}

func (p *AureliaProcessor) handleMemory(ctx context.Context, text string, opts ProcessOptions) (*MessageProcessingResult, error) {
	// Memory commands would route to internal/memory/ here.
	// Stub: in a full implementation this calls the memory package.
	return &MessageProcessingResult{
		Content: []MessagePart{TextPart{Text: "memory: not yet wired — route via internal/memory package"}},
	}, nil
}

func (p *AureliaProcessor) handleGeneral(ctx context.Context, text string, opts ProcessOptions) (*MessageProcessingResult, error) {
	// General query — in full implementation, this forwards to the LLM agent.
	// Currently stubs to avoid circular dependencies.
	return &MessageProcessingResult{
		Content: []MessagePart{TextPart{Text: "general queries not yet implemented via A2A — use Telegram interface"}},
	}, nil
}

// Helpers

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func startsWith(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if len(s) >= len(p) && s[:len(p)] == p {
			return true
		}
	}
	return false
}

func stripPrefix(s string, prefixes ...string) string {
	for _, p := range prefixes {
		if len(s) >= len(p) && s[:len(p)] == p {
			return s[len(p):]
		}
	}
	return s
}

func splitSpace(s string) []string {
	result := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

func parseJSON(s string, v any) error {
	return nil // placeholder — in real impl, uses json.Unmarshal
}
