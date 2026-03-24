package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const defaultOllamaBaseURL = "http://127.0.0.1:11434"

func ollamaChatURL(baseURL string) string {
	if baseURL == "" {
		baseURL = defaultOllamaBaseURL
	}
	return strings.TrimRight(baseURL, "/") + "/v1/chat/completions"
}

func ollamaModelsURL(baseURL string) string {
	if baseURL == "" {
		baseURL = defaultOllamaBaseURL
	}
	return strings.TrimRight(baseURL, "/") + "/v1/models"
}

func NewOllamaProvider(baseURL, model string) *OpenAICompatibleProvider {
	return NewOllamaProviderWithOptions(baseURL, model, OpenAICompatibleRequestOptions{})
}

func NewOllamaProviderWithOptions(baseURL, model string, request OpenAICompatibleRequestOptions) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		BaseURL:   ollamaChatURL(baseURL),
		Model:     model,
		UserAgent: "Aurelia/1.0",
		Request:   request,
	})
}

func listOllamaModels(ctx context.Context, baseURL string, client *http.Client) ([]ModelOption, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama models list returned %s", resp.Status)
	}

	var payload struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	models := make([]ModelOption, 0, len(payload.Data))
	for _, model := range payload.Data {
		if !looksLikeOllamaChatModel(model.ID) {
			continue
		}
		models = append(models, ModelOption{
			ID:   model.ID,
			Name: model.ID,
		})
	}

	sort.SliceStable(models, func(i, j int) bool {
		return rankOllamaModel(models[i].ID) < rankOllamaModel(models[j].ID)
	})
	return models, nil
}

func looksLikeOllamaChatModel(modelID string) bool {
	lower := strings.ToLower(modelID)
	if strings.Contains(lower, "embed") {
		return false
	}
	if strings.HasPrefix(lower, "bge-") || strings.HasPrefix(lower, "mxbai-") {
		return false
	}
	return true
}

func rankOllamaModel(modelID string) int {
	switch {
	case strings.HasPrefix(modelID, "gemma3:27b"):
		return 0
	case strings.HasPrefix(modelID, "gemma3:12b"):
		return 1
	default:
		return 10
	}
}
