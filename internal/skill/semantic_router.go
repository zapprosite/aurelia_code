package skill

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/observability"
)

const defaultSemanticSkillsCollection = "aurelia_skills"
const defaultSemanticSkillsEmbedding = "bge-m3"

// SemanticRouter uses local embeddings (e.g. Ollama) and Qdrant to route skills based on vector similarity.
type SemanticRouter struct {
	qdrantURL      string
	qdrantAPIKey   string
	collection     string
	embeddingModel string
	embedURL       string
	client         *http.Client
	ensureOnce     sync.Once
	ensureErr      error
}

// NewSemanticRouter creates a new vector-based skill router.
func NewSemanticRouter(qdrantURL, qdrantAPIKey, collection, embeddingModel, ollamaURL string) *SemanticRouter {
	col := strings.TrimSpace(collection)
	if col == "" {
		col = defaultSemanticSkillsCollection
	}
	emb := strings.TrimSpace(embeddingModel)
	if emb == "" {
		emb = defaultSemanticSkillsEmbedding
	}
	return &SemanticRouter{
		qdrantURL:      strings.TrimRight(strings.TrimSpace(qdrantURL), "/"),
		qdrantAPIKey:   strings.TrimSpace(qdrantAPIKey),
		collection:     col,
		embeddingModel: emb,
		embedURL:       strings.TrimRight(strings.TrimSpace(ollamaURL), "/") + "/api/embed",
		client:         &http.Client{Timeout: 30 * time.Second}, // Sync can take longer for batches
	}
}

// SyncSkills embeds the descriptions and upserts to aurelia_skills collection
func (r *SemanticRouter) SyncSkills(ctx context.Context, skills map[string]Skill) error {
	logger := observability.Logger("skill.semantic_router")
	if r.qdrantURL == "" || len(skills) == 0 {
		return nil
	}

	points := make([]map[string]any, 0, len(skills))
	for name, skill := range skills {
		text := fmt.Sprintf("%s: %s", name, skill.Metadata.Description)
		vector, err := r.embed(ctx, text)
		if err != nil {
			logger.Warn("failed to embed skill for sync", "skill", name, "err", err)
			continue
		}

		// Ensure collection on first successful vector (implies dimension size)
		r.ensureOnce.Do(func() {
			r.ensureErr = r.ensureCollection(ctx, len(vector))
		})
		if r.ensureErr != nil {
			return fmt.Errorf("ensure collection failed: %w", r.ensureErr)
		}

		// Generate consistent UUID from skill name
		id := uuid.NewMD5(uuid.NameSpaceURL, []byte("skill:"+name)).String()

		points = append(points, map[string]any{
			"id":     id,
			"vector": vector,
			"payload": map[string]any{
				"name":        name,
				"description": skill.Metadata.Description,
				"synced_at":   time.Now().UTC().Format(time.RFC3339),
			},
		})
	}

	if len(points) == 0 {
		return nil
	}

	body, err := json.Marshal(map[string]any{"points": points})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, r.qdrantURL+"/collections/"+r.collection+"/points?wait=true", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.qdrantAPIKey != "" {
		req.Header.Set("api-key", r.qdrantAPIKey)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("qdrant upsert request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("qdrant upsert returned %s", resp.Status)
	}

	logger.Info("semantic skill router synced", "count", len(points))
	return nil
}

// Search returns top K skill names for a query
func (r *SemanticRouter) Search(ctx context.Context, query string, limit int) ([]string, error) {
	if r.qdrantURL == "" || query == "" {
		return nil, nil // Silently degrade to fallback
	}

	vector, err := r.embed(ctx, query)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(map[string]any{
		"vector":       vector,
		"limit":        limit,
		"with_payload": true,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.qdrantURL+"/collections/"+r.collection+"/points/search", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.qdrantAPIKey != "" {
		req.Header.Set("api-key", r.qdrantAPIKey)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Collection probably doesn't exist yet
		return nil, nil
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("qdrant search returned %s", resp.Status)
	}

	var searchRes struct {
		Result []struct {
			Score   float64        `json:"score"`
			Payload map[string]any `json:"payload"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
		return nil, err
	}

	var matchedNames []string
	for _, hit := range searchRes.Result {
		// Allow slightly lower threshold for skill routing compared to strict memory
		if hit.Score < 0.3 {
			continue
		}
		if name, ok := hit.Payload["name"].(string); ok && name != "" {
			matchedNames = append(matchedNames, name)
		}
	}

	return matchedNames, nil
}

func (r *SemanticRouter) embed(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(map[string]any{
		"model": r.embeddingModel,
		"input": text,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.embedURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
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

func (r *SemanticRouter) ensureCollection(ctx context.Context, size int) error {
	body, err := json.Marshal(map[string]any{
		"vectors": map[string]any{
			"size":     size,
			"distance": "Cosine",
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, r.qdrantURL+"/collections/"+r.collection, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.qdrantAPIKey != "" {
		req.Header.Set("api-key", r.qdrantAPIKey)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusConflict && resp.StatusCode != http.StatusBadRequest {
		// Qdrant returns 400 Bad Request usually if "Collection already exists" (depending on version/mode)
		// But 409 Conflict is also standard. We ignore errors that sound like "already exists".
		// For safety we'll log it if it's 400, but not fail the whole app.
		return fmt.Errorf("qdrant ensure collection returned %s", resp.Status)
	}
	return nil
}
