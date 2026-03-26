package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// ModelOption is a selectable model entry for onboarding and config UIs.
type ModelOption struct {
	ID                 string
	Name               string
	SupportsImageInput bool
	SupportsTools      bool
	IsFree             bool
}

func (m ModelOption) Label() string {
	badges := make([]string, 0, 3)
	if m.SupportsImageInput {
		badges = append(badges, "vision")
	}
	if m.SupportsTools {
		badges = append(badges, "tools")
	}
	if m.IsFree {
		badges = append(badges, "free")
	}
	suffix := ""
	if len(badges) != 0 {
		suffix = " [" + strings.Join(badges, ", ") + "]"
	}
	if m.Name == "" || m.Name == m.ID {
		return m.ID + suffix
	}
	return fmt.Sprintf("%s (%s)%s", m.Name, m.ID, suffix)
}

// ModelCatalogCredentials carries provider-specific credentials used by model catalogs.
type ModelCatalogCredentials struct {
	AnthropicAPIKey  string
	GoogleAPIKey     string
	OpenRouterAPIKey string
	OpenAIAPIKey     string
	OpenAIAuthMode   string
}

const googleModelsURL = "https://generativelanguage.googleapis.com/v1beta/models"

// ListModels returns the available model options for a provider.
// Providers may use remote discovery or curated fallbacks internally.
func ListModels(ctx context.Context, provider string, creds ModelCatalogCredentials) ([]ModelOption, error) {
	_ = ctx
	_ = creds

	switch provider {
	case "google":
		if creds.GoogleAPIKey != "" {
			return listGoogleModels(ctx, creds.GoogleAPIKey, googleModelsURL, http.DefaultClient)
		}
		return fallbackModels("google"), nil
	case "ollama":
		models, err := listOllamaModels(ctx, ollamaModelsURL(""), http.DefaultClient)
		if err == nil && len(models) != 0 {
			return models, nil
		}
		return fallbackModels("ollama"), nil
	case "openrouter":
		return listOpenRouterModels(ctx, creds.OpenRouterAPIKey, openRouterModelsURL, http.DefaultClient)
	case "openai":
		if creds.OpenAIAPIKey != "" {
			return listOpenAIModels(ctx, creds.OpenAIAPIKey, openAIModelsURL, http.DefaultClient)
		}
		return fallbackModels("openai"), nil
	case "groq":
		// Groq uses OpenAI-compatible models list
		return listOpenAIModels(ctx, creds.AnthropicAPIKey, "https://api.groq.com/openai/v1/models", http.DefaultClient)
	default:
		return nil, fmt.Errorf("unsupported llm provider %q", provider)
	}
}

// FallbackModels returns curated default models when discovery is unavailable.
func FallbackModels(provider string) []ModelOption {
	return fallbackModels(provider)
}

func fallbackModels(provider string) []ModelOption {
	switch provider {
	case "google":
		return []ModelOption{
			{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", SupportsImageInput: true},
			{ID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", SupportsImageInput: true},
			{ID: "gemini-2.5-flash-lite", Name: "Gemini 2.5 Flash-Lite", SupportsImageInput: true},
		}
	case "ollama":
		return []ModelOption{
			{ID: "gemma3:12b", Name: "Gemma 3 12B", SupportsTools: true},
		}
	case "openrouter":
		return []ModelOption{
			{ID: "openrouter/auto", Name: "OpenRouter Auto"},
			{ID: "openrouter/free", Name: "OpenRouter Free Router", IsFree: true},
		}
	case "openai":
		return []ModelOption{
			{ID: "gpt-5.4", Name: "GPT-5.4", SupportsImageInput: true, SupportsTools: true},
			{ID: "gpt-5-mini", Name: "GPT-5 mini", SupportsImageInput: true, SupportsTools: true},
			{ID: "o4-mini", Name: "o4-mini", SupportsImageInput: true, SupportsTools: true},
		}
default:
		return nil
	}
}

func listGoogleModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelOption, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"?key="+apiKey, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google models list returned %s", resp.Status)
	}

	var payload struct {
		Models []struct {
			Name                      string   `json:"name"`
			DisplayName               string   `json:"displayName"`
			SupportedGenerationMethod []string `json:"supportedGenerationMethods"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	var models []ModelOption
	for _, model := range payload.Models {
		id := strings.TrimPrefix(model.Name, "models/")
		if !supportsGenerateContent(model.SupportedGenerationMethod) || !looksLikeGoogleChatModel(id) {
			continue
		}
		models = append(models, ModelOption{
			ID:                 id,
			Name:               model.DisplayName,
			SupportsImageInput: true,
		})
	}
	sort.SliceStable(models, func(i, j int) bool {
		return rankGoogleModel(models[i].ID) < rankGoogleModel(models[j].ID)
	})
	return models, nil
}

func supportsGenerateContent(methods []string) bool {
	for _, method := range methods {
		if method == "generateContent" {
			return true
		}
	}
	return false
}

func looksLikeGoogleChatModel(modelID string) bool {
	if !strings.HasPrefix(modelID, "gemini-") {
		return false
	}
	blocked := []string{"preview", "exp", "tts", "live", "image", "embedding", "aqa"}
	for _, token := range blocked {
		if strings.Contains(modelID, token) {
			return false
		}
	}
	return true
}

func rankGoogleModel(modelID string) int {
	switch modelID {
	case "gemini-2.5-pro":
		return 0
	case "gemini-2.5-flash":
		return 1
	case "gemini-2.5-flash-lite":
		return 2
	default:
		return 10
	}
}
