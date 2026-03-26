//go:build integration

package memory

import (
	"context"
	"math"
	"os"
	"testing"
	"time"
)

func ollamaURLForTest(t *testing.T) string {
	t.Helper()
	u := os.Getenv("OLLAMA_URL")
	if u == "" {
		u = "http://127.0.0.1:11434"
	}
	return u
}

func TestEmbedText_RealOllama(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: -short skips")
	}
	ollamaURL := ollamaURLForTest(t)
	client := NewSemanticHTTPClient(30 * time.Second)

	vec, err := EmbedText(context.Background(), client, ollamaURL, "nomic-embed-text", "aurelia homelab test")
	if err != nil {
		t.Fatalf("EmbedText() error = %v (is Ollama running at %s with nomic-embed-text?)", err, ollamaURL)
	}
	if len(vec) < 512 {
		t.Fatalf("expected vector dim >= 512, got %d", len(vec))
	}
	t.Logf("embed dim=%d OK", len(vec))
}

func TestEmbedText_CosineSimilarity_SimilarTexts(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: -short skips")
	}
	ollamaURL := ollamaURLForTest(t)
	client := NewSemanticHTTPClient(30 * time.Second)
	ctx := context.Background()

	v1, err := EmbedText(ctx, client, ollamaURL, "nomic-embed-text", "fluxo de caixa da empresa")
	if err != nil {
		t.Fatalf("EmbedText(v1) error = %v", err)
	}
	v2, err := EmbedText(ctx, client, ollamaURL, "nomic-embed-text", "caixa financeiro do negócio")
	if err != nil {
		t.Fatalf("EmbedText(v2) error = %v", err)
	}
	vDiff, err := EmbedText(ctx, client, ollamaURL, "nomic-embed-text", "receita de bolo de chocolate")
	if err != nil {
		t.Fatalf("EmbedText(vDiff) error = %v", err)
	}

	simSimilar := cosine(v1, v2)
	simDiff := cosine(v1, vDiff)

	t.Logf("cosine(similar)=%.4f cosine(diff)=%.4f", simSimilar, simDiff)

	if simSimilar < 0.6 {
		t.Fatalf("expected similar texts to have cosine >= 0.6, got %.4f", simSimilar)
	}
	if simDiff >= simSimilar {
		t.Fatalf("expected dissimilar text to score lower than similar (%.4f >= %.4f)", simDiff, simSimilar)
	}
}

func cosine(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
