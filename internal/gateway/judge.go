package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/pkg/llm"
)

type JudgeResult struct {
	Class      string  `json:"class"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

type Judge interface {
	Judge(ctx context.Context, task string, history []agent.Message) (*JudgeResult, error)
}

type GemmaJudge struct {
	provider agent.LLMProvider
}

func NewGemmaJudge(baseURL, model string) *GemmaJudge {
	lowTemp := 0.1
	// Usamos o provedor Ollama diretamente para o Gemma 3
	provider := llm.NewOllamaProviderWithOptions(baseURL, model, llm.OpenAICompatibleRequestOptions{
		Temperature: &lowTemp,
		MaxTokens:   256,
		ExtraFields: map[string]any{
			"format": "json", // Força o Ollama a retornar JSON estruturado se o modelo suportar
		},
	})
	return &GemmaJudge{provider: provider}
}

const judgeSystemPrompt = `You are a specialized task classifier for an LLM Gateway.
Your goal is to categorize the user task into one of the following classes:

- simple_short: Greetings, simple questions, direct factual queries, or small modifications.
- professional: Business responses — sales proposals, project management, CRM updates, client briefings, budgets, or domain-specific content.
- coding_main: Programming tasks, debugging, architecture, or logic.
- computer_use_jarvis: Tasks requiring web navigation, GUI interaction, or browser-based actions (clicking, typing, extracting from sites) via Stagehand/MCP.
- voice_multimodal: Requests specifically asking for voice response, processing audio, or multimodal analysis (images/PDFs/Live Voice).
- critical: High-stakes decisions, security audits, or complex multi-step reasoning.

Output ONLY a valid JSON object with this structure:
{
  "class": "simple_short | professional | coding_main | computer_use_jarvis | voice_multimodal | critical",
  "confidence": float (0-1),
  "reason": "short explanation"
}`

func (g *GemmaJudge) Judge(ctx context.Context, task string, history []agent.Message) (*JudgeResult, error) {
	prompt := fmt.Sprintf("Task to classify: %s", task)
	
	resp, err := g.provider.GenerateContent(ctx, judgeSystemPrompt, []agent.Message{
		{Role: "user", Content: prompt},
	}, nil)
	
	if err != nil {
		return nil, fmt.Errorf("judge generation failed: %w", err)
	}

	var result JudgeResult
	content := strings.TrimSpace(resp.Content)
	jsonContent := extractJSON(content)

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		// Try to extract class from free-form text as last resort.
		if cls := extractClassFromText(content); cls != "" {
			return &JudgeResult{Class: cls, Confidence: 0.5, Reason: "text extraction fallback"}, nil
		}
		return nil, fmt.Errorf("failed to parse judge output: %w | content: %s", err, content)
	}

	return &result, nil
}

// extractJSON finds the first '{' and last '}' to isolate a JSON block.
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || start >= end {
		return s // Return original as fallback
	}
	return s[start : end+1]
}

// extractClassFromText tries to find a valid class name in free-form text from the judge.
func extractClassFromText(s string) string {
	lower := strings.ToLower(s)
	classes := []string{"professional", "simple_short", "coding_main", "long_context_or_multimodal", "critical", "curation", "maintenance"}
	for _, cls := range classes {
		if strings.Contains(lower, cls) {
			return cls
		}
	}
	return ""
}
