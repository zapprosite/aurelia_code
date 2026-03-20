package telegram

import "strings"

const (
	genericExecutionFailureMessage = "Nao consegui concluir isso agora por uma falha temporaria do runtime. Tente novamente em alguns segundos."
	genericResponseGuardMessage    = "Nao consegui produzir uma resposta confiavel agora. Tente novamente com um pedido mais direto."
)

var (
	technicalErrorMarkers = []string{
		"provider error",
		"openai-compatible api error",
		"route breaker open",
		"budget exceeded",
		"empty guarded content",
		"invalid_request_error",
		"registry.ollama.ai",
		"does not support tools",
		"all gateway routes failed",
		"deepseek/deepseek-v3.2",
		"openrouter:",
		"minimax/minimax",
		"qwen3.5:",
	}
	internalOutputMarkers = []string{
		"router parse error",
		"active skill context",
		"executing tool:",
		`"skillname":`,
		"tool_name=",
		"arg_keys=",
		"mcp_",
		"provider error:",
		"openai-compatible api error",
		"route breaker open",
		"budget exceeded",
		"empty guarded content",
		"deepseek/deepseek-v3.2",
	}
)

func sanitizeUserVisibleErrorMessage(errMsg string) string {
	trimmed := strings.TrimSpace(errMsg)
	if trimmed == "" {
		return genericExecutionFailureMessage
	}

	lower := strings.ToLower(trimmed)
	switch {
	case strings.Contains(lower, "max iterations reached"),
		strings.Contains(lower, "context cancelled by timer"),
		strings.Contains(lower, "timeout"):
		return "Nao consegui concluir essa tarefa a tempo. Tente novamente com um pedido mais curto ou em alguns segundos."
	case containsAny(lower, technicalErrorMarkers...):
		return genericExecutionFailureMessage
	default:
		return trimmed
	}
}

func sanitizeAssistantOutputForUser(text string) string {
	trimmed := strings.TrimSpace(strings.ReplaceAll(text, "\x00", ""))
	if trimmed == "" {
		return genericResponseGuardMessage
	}

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "{\"skillname\"") || strings.HasPrefix(lower, "```json\n{\"skillname\"") {
		return genericResponseGuardMessage
	}
	if containsAny(lower, internalOutputMarkers...) {
		return genericResponseGuardMessage
	}
	return trimmed
}

func containsAny(text string, markers ...string) bool {
	for _, marker := range markers {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return false
}
