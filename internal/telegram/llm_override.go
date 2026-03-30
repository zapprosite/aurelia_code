package telegram

import (
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/pkg/llm"
)

type ownedLLMProvider interface {
	agent.LLMProvider
	Close()
}

func buildBotOverrideProvider(appCfg *config.AppConfig, botCfg config.BotConfig) (ownedLLMProvider, string, string, error) {
	if appCfg == nil {
		return nil, "", "", fmt.Errorf("app config is required")
	}

	if !botHasPinnedLLM(botCfg.ID) && strings.TrimSpace(botCfg.LLMProvider) == "" && strings.TrimSpace(botCfg.LLMModel) == "" {
		return nil, "", "", nil
	}

	provider, model := EffectiveBotLLM(botCfg, appCfg.LLMProvider, appCfg.LLMModel)
	provider = strings.ToLower(strings.TrimSpace(provider))
	model = strings.TrimSpace(model)

	if provider == "" {
		return nil, "", "", fmt.Errorf("bot %q has llm override without provider", botCfg.ID)
	}
	if model == "" {
		return nil, "", "", fmt.Errorf("bot %q requires llm_model for provider %q", botCfg.ID, provider)
	}
	if !botHasPinnedLLM(botCfg.ID) &&
		strings.EqualFold(provider, strings.TrimSpace(appCfg.LLMProvider)) &&
		strings.EqualFold(model, strings.TrimSpace(appCfg.LLMModel)) {
		return nil, "", "", nil
	}

	switch provider {
	case "anthropic":
		if strings.TrimSpace(appCfg.AnthropicAPIKey) == "" {
			return nil, "", "", fmt.Errorf("bot %q requires anthropic_api_key", botCfg.ID)
		}
		return llm.NewAnthropicProvider(appCfg.AnthropicAPIKey, model), provider, model, nil
	case "groq":
		if strings.TrimSpace(appCfg.GroqAPIKey) == "" {
			return nil, "", "", fmt.Errorf("bot %q requires groq_api_key", botCfg.ID)
		}
		return llm.NewGroqProvider(appCfg.GroqAPIKey, model), provider, model, nil
	case "minimax":
		if strings.TrimSpace(appCfg.MiniMaxAPIKey) == "" {
			return nil, "", "", fmt.Errorf("bot %q requires minimax_api_key", botCfg.ID)
		}
		return llm.NewMiniMaxProvider(appCfg.MiniMaxAPIKey, model), provider, model, nil
	case "ollama":
		return llm.NewOllamaProvider(appCfg.OllamaURL, model), provider, model, nil
	case "openai":
		if strings.TrimSpace(appCfg.OpenAIAPIKey) == "" {
			return nil, "", "", fmt.Errorf("bot %q requires openai_api_key", botCfg.ID)
		}
		return llm.NewOpenAIProvider(appCfg.OpenAIAPIKey, model), provider, model, nil
	case "openrouter":
		if strings.TrimSpace(appCfg.OpenRouterAPIKey) == "" {
			return nil, "", "", fmt.Errorf("bot %q requires openrouter_api_key", botCfg.ID)
		}
		return llm.NewOpenRouterProvider(appCfg.OpenRouterAPIKey, model), provider, model, nil
	default:
		return nil, "", "", fmt.Errorf("bot %q has unsupported llm provider %q", botCfg.ID, provider)
	}
}
