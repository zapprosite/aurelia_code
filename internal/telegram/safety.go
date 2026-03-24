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
		"openrouter:",
		"groq:",
		"minimax/minimax",
		"gemma3:",
		"google/gemini-2.5-flash",
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
		"google/gemini-2.5-flash",
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
	case strings.Contains(lower, "all gateway routes failed"):
		return "Nenhum provedor de IA respondeu. Verifique: Groq API, Ollama local, e OpenRouter. Use /status para diagnostico."
	case strings.Contains(lower, "budget exceeded"):
		return "Limite de custo atingido para este provedor. Tente mais tarde ou use /status."
	case strings.Contains(lower, "route breaker open"):
		return "Provedor temporariamente indisponivel (circuit breaker aberto). Tente em 1-2 minutos."
	case strings.Contains(lower, "rate_limit_exceeded"),
		strings.Contains(lower, "request too large"),
		strings.Contains(lower, "tokens per minute"):
		return "Limite de tokens atingido. Tente uma mensagem mais curta ou aguarde alguns segundos."
	case containsAny(lower, technicalErrorMarkers...):
		return genericExecutionFailureMessage
	default:
		return trimmed
	}
}

func sanitizeAssistantOutputForUser(text string) string {
	text = strings.ReplaceAll(text, "\x00", "")
	
	// Remove tags de pensamento/raciocínio interno (padrão DeepSeek/Claude/Gemma)
	text = removeTag(text, "thought")
	
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return genericResponseGuardMessage
	}

	lower := strings.ToLower(trimmed)
	
	// Se a resposta for puramente JSON de roteamento, é um erro de escape do LLM
	if strings.HasPrefix(lower, "{\"skillname\"") || 
	   strings.HasPrefix(lower, "```json\n{\"skillname\"") ||
	   strings.HasPrefix(lower, "{\"tool_code\"") {
		return genericResponseGuardMessage
	}

	if containsAny(lower, internalOutputMarkers...) {
		return genericResponseGuardMessage
	}

	return trimmed
}

func removeTag(text, tag string) string {
	startTag := "<" + tag + ">"
	endTag := "</" + tag + ">"
	
	for {
		startIdx := strings.Index(text, startTag)
		if startIdx == -1 {
			break
		}
		endIdx := strings.Index(text, endTag)
		if endIdx == -1 {
			// Tag não fechada, removemos do início até o fim
			return strings.TrimSpace(text[:startIdx])
		}
		// Remove a tag e o conteúdo interno
		text = text[:startIdx] + text[endIdx+len(endTag):]
	}
	return text
}

func containsAny(text string, markers ...string) bool {
	for _, marker := range markers {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return false
}
