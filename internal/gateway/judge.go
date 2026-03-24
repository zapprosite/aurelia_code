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

func NewGemmaJudge(model string) *GemmaJudge {
	lowTemp := 0.1
	// Usamos o provedor Ollama diretamente para o Gemma 3
	provider := llm.NewOllamaProviderWithOptions(model, llm.OpenAICompatibleRequestOptions{
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
- coding_main: Programming tasks, debugging, architecture design, or logic implementation.
- long_context_or_multimodal: Tasks involving images, screenshots, PDFs, or requests that require analyzing very large chunks of text.
- critical: High-stakes decisions, security audits, or complex multi-step reasoning where accuracy is paramount.

Output ONLY a valid JSON object with this structure:
{
  "class": "simple_short | coding_main | long_context_or_multimodal | critical",
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
	// Limpar possíveis markdown blocks se o modelo injetar
	content := strings.TrimSpace(resp.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse judge output: %w | content: %s", err, content)
	}

	return &result, nil
}
