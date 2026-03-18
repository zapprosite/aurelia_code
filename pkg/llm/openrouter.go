package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	openRouterChatCompletionsURL = "https://openrouter.ai/api/v1/chat/completions"
	openRouterModelsURL          = "https://openrouter.ai/api/v1/models"
	openRouterReferer            = "https://github.com/Lordymine/aurelia"
	openRouterTitle              = "Aurelia"
)

func NewOpenRouterProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		APIKey:    apiKey,
		BaseURL:   openRouterChatCompletionsURL,
		Model:     model,
		UserAgent: "Aurelia/1.0",
		Headers: map[string]string{
			"HTTP-Referer": openRouterReferer,
			"X-Title":      openRouterTitle,
		},
	})
}

func listOpenRouterModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelOption, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("HTTP-Referer", openRouterReferer)
	req.Header.Set("X-Title", openRouterTitle)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter models list returned %s", resp.Status)
	}

	var payload struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	models := []ModelOption{
		{ID: "openrouter/auto", Name: "OpenRouter Auto"},
		{ID: "openrouter/free", Name: "OpenRouter Free Router"},
	}
	for _, model := range payload.Data {
		models = append(models, ModelOption{
			ID:   model.ID,
			Name: model.Name,
		})
	}
	return models, nil
}
