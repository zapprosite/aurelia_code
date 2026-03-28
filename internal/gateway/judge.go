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
	// [SOTA 2026] Unificação de protocolo: Juiz agora fala OpenAI nativo via LiteLLM na porta 4000
	provider := llm.NewOpenAICompatibleProvider(llm.OpenAICompatibleConfig{
		BaseURL: baseURL + "/v1/chat/completions",
		Model:   model,
		Request: llm.OpenAICompatibleRequestOptions{
			Temperature: &lowTemp,
			MaxTokens:   256,
			ExtraFields: map[string]any{
				"format": "json",
			},
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
- linux_god_mode: Full OS control, bash commands, system logs, patching, or infrastructure management.

Output ONLY a valid JSON object with this structure:
{
  "class": "simple_short | professional | coding_main | computer_use_jarvis | voice_multimodal | critical | linux_god_mode",
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
		if cls := extractClassFromText(content); cls != "" {
			return &JudgeResult{Class: cls, Confidence: 0.5, Reason: "text extraction fallback"}, nil
		}
		return nil, fmt.Errorf("failed to parse judge output: %w | content: %s", err, content)
	}

	return &result, nil
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || start >= end {
		return s
	}
	return s[start : end+1]
}

func extractClassFromText(s string) string {
	lower := strings.ToLower(s)
	classes := []string{"professional", "simple_short", "coding_main", "long_context_or_multimodal", "critical", "linux_god_mode"}
	for _, cls := range classes {
		if strings.Contains(lower, cls) {
			return cls
		}
	}
	return ""
}
