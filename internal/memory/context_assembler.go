package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ContextAssembler struct {
	qdrantURL      string
	qdrantAPIKey   string
	qdrantCollection string
	embedURL       string
	embeddingModel string
	mem            *MemoryManager
	client         *http.Client
}

func NewContextAssembler(qdrantURL, apiKey, collection, embeddingModel, ollamaURL string, mem *MemoryManager) *ContextAssembler {
	return &ContextAssembler{
		qdrantURL:        strings.TrimRight(strings.TrimSpace(qdrantURL), "/"),
		qdrantAPIKey:     strings.TrimSpace(apiKey),
		qdrantCollection: strings.TrimSpace(collection),
		embedURL:         strings.TrimRight(strings.TrimSpace(ollamaURL), "/") + "/api/embed",
		embeddingModel:   strings.TrimSpace(embeddingModel),
		mem:              mem,
		client:           &http.Client{Timeout: 10 * time.Second},
	}
}

// AssembleContext returns a formatted markdown string joining semantic matches and recent facts/notes.
func (a *ContextAssembler) AssembleContext(ctx context.Context, query string) string {
	var sb strings.Builder

	// 1. Fetch from Qdrant
	if qdrantRes, err := a.searchQdrant(ctx, query); err == nil && qdrantRes != "" {
		sb.WriteString("Arquivos Históricos (Qdrant Semantic Search):\n")
		sb.WriteString(qdrantRes)
		sb.WriteString("\n\n")
	}

	// 2. Fetch from SQLite Notes
	if a.mem != nil {
		if recentNotes, err := a.mem.GetGlobalTopics(ctx, 10); err == nil && len(recentNotes) > 0 {
			sb.WriteString("Tópicos e Notas Locais (SQLite):\n")
			for _, note := range recentNotes {
				sb.WriteString(fmt.Sprintf("- [%s] %s\n", note.Kind, note.Topic))
			}
		}
	}

	return strings.TrimSpace(sb.String())
}

func (a *ContextAssembler) searchQdrant(ctx context.Context, text string) (string, error) {
	if a.qdrantURL == "" || text == "" {
		return "", nil
	}

	// 1. Embed Query
	vector, err := a.embed(ctx, text)
	if err != nil {
		return "", err
	}

	// 2. Search Qdrant
	body, err := json.Marshal(map[string]any{
		"vector":       vector,
		"limit":        3,
		"with_payload": true,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.qdrantURL+"/collections/"+a.qdrantCollection+"/points/search", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if a.qdrantAPIKey != "" {
		req.Header.Set("api-key", a.qdrantAPIKey)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("qdrant search returned %s", resp.Status)
	}

	var searchRes struct {
		Result []struct {
			Score   float64        `json:"score"`
			Payload map[string]any `json:"payload"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
		return "", err
	}

	var out strings.Builder
	for _, hit := range searchRes.Result {
		// Minimum certainty threshold
		if hit.Score < 0.4 {
			continue
		}
		
		transcript, _ := hit.Payload["transcript"].(string)
		if transcript != "" {
			out.WriteString(fmt.Sprintf("> (Score: %.2f) %s\n", hit.Score, transcript))
		}
	}

	return out.String(), nil
}

func (a *ContextAssembler) embed(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(map[string]any{
		"model": a.embeddingModel,
		"input": text,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.embedURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ollama embed returned %s", resp.Status)
	}

	var payload struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload.Embeddings) == 0 || len(payload.Embeddings[0]) == 0 {
		return nil, fmt.Errorf("ollama embed returned no vectors")
	}
	return payload.Embeddings[0], nil
}
