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

type SemanticPoint struct {
	ID      any
	Score   float64
	Payload map[string]any
}

func NewSemanticHTTPClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &http.Client{Timeout: timeout}
}

func EmbedText(ctx context.Context, client *http.Client, ollamaURL, model, text string) ([]float32, error) {
	ollamaURL = strings.TrimRight(strings.TrimSpace(ollamaURL), "/")
	model = strings.TrimSpace(model)
	text = strings.TrimSpace(text)
	if ollamaURL == "" {
		return nil, fmt.Errorf("ollama url is required for semantic search")
	}
	if model == "" {
		return nil, fmt.Errorf("embedding model is required for semantic search")
	}
	if text == "" {
		return nil, fmt.Errorf("query text is required for semantic search")
	}
	if client == nil {
		client = NewSemanticHTTPClient(10 * time.Second)
	}

	body, err := json.Marshal(map[string]any{
		"model": model,
		"input": text,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaURL+"/api/embed", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
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

func SearchSemantic(ctx context.Context, client *http.Client, qdrantURL, collection, apiKey string, vector []float32, limit int) ([]SemanticPoint, error) {
	qdrantURL = strings.TrimRight(strings.TrimSpace(qdrantURL), "/")
	collection = strings.TrimSpace(collection)
	if qdrantURL == "" {
		return nil, fmt.Errorf("qdrant url is required for semantic search")
	}
	if collection == "" {
		return nil, fmt.Errorf("qdrant collection is required for semantic search")
	}
	if len(vector) == 0 {
		return nil, fmt.Errorf("query vector is required for semantic search")
	}
	if limit <= 0 {
		limit = 10
	}
	if client == nil {
		client = NewSemanticHTTPClient(10 * time.Second)
	}

	body, err := json.Marshal(map[string]any{
		"vector":       vector,
		"limit":        limit,
		"with_payload": true,
		"with_vector":  false,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, qdrantURL+"/collections/"+collection+"/points/search", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("api-key", apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("qdrant search returned %s", resp.Status)
	}

	var raw struct {
		Result []struct {
			ID      any            `json:"id"`
			Score   float64        `json:"score"`
			Payload map[string]any `json:"payload"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	points := make([]SemanticPoint, 0, len(raw.Result))
	for _, item := range raw.Result {
		points = append(points, SemanticPoint{
			ID:      item.ID,
			Score:   item.Score,
			Payload: NormalizeSemanticPayload(item.Payload),
		})
	}
	return points, nil
}

func ScrollPoints(ctx context.Context, client *http.Client, qdrantURL, collection, apiKey string, limit int) ([]SemanticPoint, error) {
	qdrantURL = strings.TrimRight(strings.TrimSpace(qdrantURL), "/")
	collection = strings.TrimSpace(collection)
	if qdrantURL == "" {
		return nil, fmt.Errorf("qdrant url is required for scroll")
	}
	if collection == "" {
		return nil, fmt.Errorf("qdrant collection is required for scroll")
	}
	if limit <= 0 {
		limit = 20
	}
	if client == nil {
		client = NewSemanticHTTPClient(10 * time.Second)
	}

	body, err := json.Marshal(map[string]any{
		"limit":        limit,
		"with_payload": true,
		"with_vector":  false,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, qdrantURL+"/collections/"+collection+"/points/scroll", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("api-key", apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("qdrant scroll returned %s", resp.Status)
	}

	var raw struct {
		Result struct {
			Points []struct {
				ID      any            `json:"id"`
				Payload map[string]any `json:"payload"`
			} `json:"points"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	points := make([]SemanticPoint, 0, len(raw.Result.Points))
	for _, item := range raw.Result.Points {
		points = append(points, SemanticPoint{
			ID:      item.ID,
			Payload: NormalizeSemanticPayload(item.Payload),
		})
	}
	return points, nil
}

func FilterPointsLexical(points []SemanticPoint, query string) []SemanticPoint {
	query = strings.TrimSpace(query)
	if query == "" {
		return points
	}

	lowered := strings.ToLower(query)
	filtered := make([]SemanticPoint, 0, len(points))
	for _, point := range points {
		if strings.Contains(strings.ToLower(ExtractSearchableText(point.Payload)), lowered) {
			filtered = append(filtered, point)
		}
	}
	return filtered
}

func FilterPointsByCanonicalBotID(points []SemanticPoint, botID string) []SemanticPoint {
	botID = strings.TrimSpace(botID)
	if botID == "" {
		return points
	}

	aliases := canonicalBotAliases(botID)
	filtered := make([]SemanticPoint, 0, len(points))
	for _, point := range points {
		candidate := firstString(NormalizeSemanticPayload(point.Payload), "canonical_bot_id")
		for _, alias := range aliases {
			if candidate == alias {
				filtered = append(filtered, point)
				break
			}
		}
	}
	return filtered
}

func ExtractSearchableText(payload map[string]any) string {
	if len(payload) == 0 {
		return ""
	}

	title := firstString(payload, "title")
	main := firstString(payload, "text", "content", "summary", "transcript", "message", "body")
	switch {
	case title != "" && main != "" && title != main:
		return title + " — " + main
	case main != "":
		return main
	default:
		return title
	}
}

func NormalizeSemanticPayload(payload map[string]any) map[string]any {
	if payload == nil {
		payload = map[string]any{}
	}

	normalized := make(map[string]any, len(payload)+3)
	for key, value := range payload {
		normalized[key] = value
	}

	if _, ok := normalized["canonical_bot_id"]; !ok {
		if botID := firstString(normalized, "bot_id"); botID != "" {
			normalized["canonical_bot_id"] = botID
		}
	}
	if _, ok := normalized["text"]; !ok {
		if text := ExtractSearchableText(normalized); text != "" {
			normalized["text"] = text
		}
	}
	if _, ok := normalized["ts"]; !ok {
		if ts := ExtractPayloadTimestamp(normalized); !ts.IsZero() {
			normalized["ts"] = ts.Unix()
		}
	}

	return normalized
}

func ExtractPayloadTimestamp(payload map[string]any) time.Time {
	if payload == nil {
		return time.Time{}
	}

	if raw, ok := payload["ts"]; ok {
		switch value := raw.(type) {
		case float64:
			return time.Unix(int64(value), 0).UTC()
		case int64:
			return time.Unix(value, 0).UTC()
		case int:
			return time.Unix(int64(value), 0).UTC()
		case json.Number:
			if unix, err := value.Int64(); err == nil {
				return time.Unix(unix, 0).UTC()
			}
		case string:
			if value != "" {
				if parsed, err := time.Parse(time.RFC3339, value); err == nil {
					return parsed.UTC()
				}
			}
		}
	}

	for _, key := range []string{"updated_at", "created_at", "mirrored_at"} {
		if raw, ok := payload[key]; ok {
			if value, ok := raw.(string); ok && strings.TrimSpace(value) != "" {
				if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
					return parsed.UTC()
				}
				if parsed, err := time.Parse(time.RFC3339, value); err == nil {
					return parsed.UTC()
				}
			}
		}
	}

	return time.Time{}
}

func firstString(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		if raw, ok := payload[key]; ok {
			if value, ok := raw.(string); ok {
				value = strings.TrimSpace(value)
				if value != "" {
					return value
				}
			}
		}
	}
	return ""
}

func canonicalBotAliases(botID string) []string {
	botID = strings.TrimSpace(botID)
	if botID == "" {
		return nil
	}

	aliases := []string{botID}
	switch botID {
	case "aurelia":
		aliases = append(aliases, "aurelia_code")
	case "aurelia_code":
		aliases = append(aliases, "aurelia")
	}
	return aliases
}
