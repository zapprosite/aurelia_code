package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const (
	ollamaChatCompletionsURL = "http://127.0.0.1:11434/v1/chat/completions"
	ollamaModelsURL          = "http://127.0.0.1:11434/v1/models"
)

func NewOllamaProvider(model string) *OpenAICompatibleProvider {
	return NewOllamaProviderWithOptions(model, OpenAICompatibleRequestOptions{})
}

func NewOllamaProviderWithOptions(model string, request OpenAICompatibleRequestOptions) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		BaseURL:   ollamaChatCompletionsURL,
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
