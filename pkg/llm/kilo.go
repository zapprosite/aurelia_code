package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	kiloChatCompletionsURL = "https://api.kilo.ai/api/gateway/chat/completions"
	kiloModelsURL          = "https://api.kilo.ai/api/gateway/models"
)

var kiloModelsCatalogURL = kiloModelsURL

func NewKiloProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		APIKey:    apiKey,
		BaseURL:   kiloChatCompletionsURL,
		Model:     model,
		UserAgent: "Aurelia/1.0",
	})
}

func listKiloModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelOption, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kilo models list returned %s", resp.Status)
	}

	var payload struct {
		Data []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			OwnedBy string `json:"owned_by"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	models := make([]ModelOption, 0, len(payload.Data))
	for _, model := range payload.Data {
		if model.ID == "" {
			continue
		}
		name := model.Name
		if name == "" {
			name = model.ID
		}
		if model.OwnedBy != "" && !containsFold(name, model.OwnedBy) {
			name = fmt.Sprintf("%s · %s", name, model.OwnedBy)
		}
		models = append(models, ModelOption{
			ID:   model.ID,
			Name: name,
		})
	}
	return models, nil
}

func kiloModelsURLForTest(url string) func() {
	original := kiloModelsCatalogURL
	kiloModelsCatalogURL = url
	return func() {
		kiloModelsCatalogURL = original
	}
}

func containsFold(value string, needle string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}
