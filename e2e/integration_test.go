//go:build integration

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// TestFullStack_EmbedQdrantDashboard verifica a cadeia completa:
// Ollama embed → Qdrant upsert+search → sem panic ou erro silencioso.
func TestFullStack_EmbedQdrantDashboard(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: -short skips")
	}

	ollamaURL := getEnvOrDefault("OLLAMA_URL", "http://127.0.0.1:11434")
	qdrantURL := getEnvOrDefault("QDRANT_URL", "http://127.0.0.1:6333")
	qdrantKey := os.Getenv("QDRANT_API_KEY")
	httpClient := &http.Client{Timeout: 30 * time.Second}
	ctx := context.Background()

	// 1. Embed via Ollama
	embedBody, _ := json.Marshal(map[string]any{
		"model": "bge-m3",
		"input": "aurelia fullstack integration test",
	})
	embedReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, ollamaURL+"/api/embed", bytes.NewReader(embedBody))
	embedReq.Header.Set("Content-Type", "application/json")

	embedResp, err := httpClient.Do(embedReq)
	if err != nil {
		t.Fatalf("Ollama embed: %v (OLLAMA_URL=%s, bge-m3 instalado?)", err, ollamaURL)
	}
	defer embedResp.Body.Close()
	if embedResp.StatusCode >= 300 {
		t.Fatalf("Ollama embed HTTP %d", embedResp.StatusCode)
	}

	var embedPayload struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.NewDecoder(embedResp.Body).Decode(&embedPayload); err != nil {
		t.Fatalf("decode embed response: %v", err)
	}
	if len(embedPayload.Embeddings) == 0 || len(embedPayload.Embeddings[0]) == 0 {
		t.Fatal("embed retornou vetor vazio")
	}
	vec := embedPayload.Embeddings[0]
	t.Logf("Step 1/3 Ollama embed OK: dim=%d", len(vec))

	// 2. Qdrant: criar coleção temp, upsert, search, cleanup
	collection := fmt.Sprintf("aurelia_fullstack_%d", time.Now().UnixMilli())

	qdrantReq := func(method, path string, body any) *http.Response {
		t.Helper()
		var buf *bytes.Reader
		if body != nil {
			b, _ := json.Marshal(body)
			buf = bytes.NewReader(b)
		} else {
			buf = bytes.NewReader(nil)
		}
		req, _ := http.NewRequestWithContext(ctx, method, qdrantURL+path, buf)
		req.Header.Set("Content-Type", "application/json")
		if qdrantKey != "" {
			req.Header.Set("api-key", qdrantKey)
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s: %v", method, path, err)
		}
		return resp
	}

	// Criar coleção
	r := qdrantReq(http.MethodPut, "/collections/"+collection,
		map[string]any{"vectors": map[string]any{"size": len(vec), "distance": "Cosine"}})
	r.Body.Close()
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		t.Fatalf("criar coleção Qdrant HTTP %d (QDRANT_URL=%s)", r.StatusCode, qdrantURL)
	}
	t.Cleanup(func() {
		r2 := qdrantReq(http.MethodDelete, "/collections/"+collection, nil)
		r2.Body.Close()
	})

	// Upsert com payload canônico
	r = qdrantReq(http.MethodPut, "/collections/"+collection+"/points?wait=true",
		map[string]any{
			"points": []map[string]any{
				{
					"id":     1,
					"vector": vec,
					"payload": map[string]any{
						"app_id":           "aurelia",
						"repo_id":          "github.com/kocar/aurelia",
						"environment":      "integration",
						"text":             "aurelia fullstack integration test",
						"canonical_bot_id": "aurelia_code",
						"source_system":    "e2e",
						"source_id":        "fullstack-001",
						"domain":           "system",
						"ts":               time.Now().Unix(),
						"version":          1,
					},
				},
			},
		})
	r.Body.Close()
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		t.Fatalf("upsert Qdrant HTTP %d", r.StatusCode)
	}

	// Search com o mesmo vetor
	r = qdrantReq(http.MethodPost, "/collections/"+collection+"/points/search",
		map[string]any{"vector": vec, "limit": 1, "with_payload": true})
	defer r.Body.Close()
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		t.Fatalf("search Qdrant HTTP %d", r.StatusCode)
	}
	var searchResult struct {
		Result []struct {
			Score   float64        `json:"score"`
			Payload map[string]any `json:"payload"`
		} `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&searchResult); err != nil {
		t.Fatalf("decode search: %v", err)
	}
	if len(searchResult.Result) == 0 {
		t.Fatal("search retornou 0 resultados")
	}
	if searchResult.Result[0].Score < 0.99 {
		t.Fatalf("esperava score >= 0.99 para auto-busca, got %.4f", searchResult.Result[0].Score)
	}
	t.Logf("Step 2/3 Qdrant upsert+search OK: score=%.4f source_id=%v",
		searchResult.Result[0].Score, searchResult.Result[0].Payload["source_id"])

	// 3. Dashboard health check (se rodando)
	dashURL := getEnvOrDefault("DASHBOARD_URL", "http://127.0.0.1:3334")
	healthURL := getEnvOrDefault("HEALTH_URL", "http://127.0.0.1:8484")
	dashClient := &http.Client{Timeout: 3 * time.Second}

	if resp, err := dashClient.Get(healthURL + "/health"); err == nil {
		resp.Body.Close()
		t.Logf("Step 3/3 Health OK: HTTP %d", resp.StatusCode)
	} else if resp2, err2 := dashClient.Get(dashURL + "/"); err2 == nil {
		resp2.Body.Close()
		t.Logf("Step 3/3 Dashboard OK: HTTP %d", resp2.StatusCode)
	} else {
		t.Logf("Step 3/3 Dashboard/Health não alcançável (OK se serviço não rodando): %v", err)
	}
}
