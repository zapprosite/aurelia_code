//go:build integration

package memory

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

func qdrantURLForTest(t *testing.T) string {
	t.Helper()
	u := os.Getenv("QDRANT_URL")
	if u == "" {
		u = "http://127.0.0.1:6333"
	}
	return u
}

func qdrantAPIKeyForTest(t *testing.T) string {
	t.Helper()
	k := os.Getenv("QDRANT_API_KEY")
	if k == "" {
		t.Skip("QDRANT_API_KEY not set — skipping live Qdrant test")
	}
	return k
}

// qdrantDo executes a raw HTTP request to Qdrant and returns the status code.
func qdrantDo(t *testing.T, client *http.Client, method, url, apiKey string, body any) int {
	t.Helper()
	var buf *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("json.Marshal: %v", err)
		}
		buf = bytes.NewReader(b)
	} else {
		buf = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, url, buf)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("api-key", apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func TestQdrant_CRUDCanonicalPayload(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: -short skips")
	}
	qdrantURL := qdrantURLForTest(t)
	apiKey := qdrantAPIKeyForTest(t)
	client := NewSemanticHTTPClient(15 * time.Second)
	collection := fmt.Sprintf("aurelia_integration_test_%d", time.Now().UnixMilli())

	// 1. Criar coleção temporária com dim=4 (sem precisar de Ollama)
	createStatus := qdrantDo(t, client, http.MethodPut,
		qdrantURL+"/collections/"+collection, apiKey,
		map[string]any{
			"vectors": map[string]any{
				"size":     4,
				"distance": "Cosine",
			},
		},
	)
	if createStatus < 200 || createStatus >= 300 {
		t.Fatalf("criar coleção retornou HTTP %d (Qdrant rodando em %s?)", createStatus, qdrantURL)
	}
	t.Cleanup(func() {
		qdrantDo(t, client, http.MethodDelete, qdrantURL+"/collections/"+collection, apiKey, nil)
	})

	// 2. Upsert de um point com payload canônico
	now := time.Now().UTC()
	payload := map[string]any{
		"app_id":           "aurelia",
		"repo_id":          "github.com/kocar/aurelia",
		"environment":      "test",
		"text":             "teste de integração Qdrant",
		"canonical_bot_id": "aurelia_code",
		"source_system":    "test",
		"source_id":        "integration-test-001",
		"domain":           "system",
		"ts":               now.Unix(),
		"version":          1,
	}
	if err := ValidateCanonicalMemoryPayload(payload); err != nil {
		t.Fatalf("payload canônico inválido: %v", err)
	}
	upsertStatus := qdrantDo(t, client, http.MethodPut,
		qdrantURL+"/collections/"+collection+"/points?wait=true", apiKey,
		map[string]any{
			"points": []map[string]any{
				{
					"id":      1,
					"vector":  []float32{0.1, 0.2, 0.3, 0.4},
					"payload": payload,
				},
			},
		},
	)
	if upsertStatus < 200 || upsertStatus >= 300 {
		t.Fatalf("upsert retornou HTTP %d", upsertStatus)
	}

	// 3. Scroll e verificar que o point foi gravado
	points, err := ScrollPoints(context.Background(), client, qdrantURL, collection, apiKey, 10)
	if err != nil {
		t.Fatalf("ScrollPoints() error = %v", err)
	}
	if len(points) == 0 {
		t.Fatalf("esperava >= 1 point, got 0")
	}

	found := points[0]
	if found.Payload["canonical_bot_id"] != "aurelia_code" {
		t.Fatalf("canonical_bot_id esperado aurelia_code, got %v", found.Payload["canonical_bot_id"])
	}
	if found.Payload["source_system"] != "test" {
		t.Fatalf("source_system esperado test, got %v", found.Payload["source_system"])
	}
	t.Logf("CRUD Qdrant OK: collection=%s point_id=%v", collection, found.ID)
}

func TestQdrant_SearchSemantic(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: -short skips")
	}
	qdrantURL := qdrantURLForTest(t)
	apiKey := qdrantAPIKeyForTest(t)
	client := NewSemanticHTTPClient(15 * time.Second)
	collection := fmt.Sprintf("aurelia_search_test_%d", time.Now().UnixMilli())

	// Setup: criar coleção + upsert 2 pontos
	qdrantDo(t, client, http.MethodPut, qdrantURL+"/collections/"+collection, apiKey,
		map[string]any{"vectors": map[string]any{"size": 4, "distance": "Cosine"}})
	t.Cleanup(func() {
		qdrantDo(t, client, http.MethodDelete, qdrantURL+"/collections/"+collection, apiKey, nil)
	})
	qdrantDo(t, client, http.MethodPut, qdrantURL+"/collections/"+collection+"/points?wait=true", apiKey,
		map[string]any{
			"points": []map[string]any{
				{"id": 1, "vector": []float32{1.0, 0.0, 0.0, 0.0}, "payload": map[string]any{"text": "caixa financeira"}},
				{"id": 2, "vector": []float32{0.0, 1.0, 0.0, 0.0}, "payload": map[string]any{"text": "pipeline de obras"}},
			},
		})

	// Buscar com vetor próximo ao ponto 1
	results, err := SearchSemantic(context.Background(), client, qdrantURL, collection, apiKey,
		[]float32{0.99, 0.01, 0.0, 0.0}, 5)
	if err != nil {
		t.Fatalf("SearchSemantic() error = %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("esperava >= 1 resultado, got 0")
	}
	if results[0].Payload["text"] != "caixa financeira" {
		t.Fatalf("esperava hit em 'caixa financeira', got %v", results[0].Payload["text"])
	}
	t.Logf("SearchSemantic OK: top=%v score=%.4f", results[0].ID, results[0].Score)
}
