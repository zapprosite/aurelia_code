package memory

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContextAssembler_IncludesMarkdownBrainForSovereignBot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/embed":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"embeddings": [][]float32{{0.1, 0.2, 0.3}},
			})
		case "/collections/conversation_memory/points/search":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": []map[string]any{{
					"id":    "mem-1",
					"score": 0.88,
					"payload": map[string]any{
						"text":             "Memória de chat validada",
						"canonical_bot_id": "aurelia_code",
						"domain":           "conversation",
					},
				}},
			})
		case "/collections/aurelia_markdown_brain/points/search":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": []map[string]any{{
					"id":    "md-1",
					"score": 0.91,
					"payload": map[string]any{
						"text":             "Title: Brain\nPath: docs/brain.md\nSection: Overview\n\nCérebro markdown ativo",
						"canonical_bot_id": "aurelia_code",
						"repo_path":        "docs/brain.md",
						"section":          "Overview",
						"source_id":        "repo_markdown:docs/brain.md",
					},
				}},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	assembler := NewContextAssembler(server.URL, "", "conversation_memory", "aurelia_markdown_brain", "nomic-embed-text", server.URL, nil)
	got := assembler.AssembleContextForBot(context.Background(), "aurelia", "brain")

	if !strings.Contains(got, "Arquivos Históricos") {
		t.Fatalf("expected conversation memory section, got %q", got)
	}
	if !strings.Contains(got, "Markdown Brain") {
		t.Fatalf("expected markdown brain section, got %q", got)
	}
}

func TestContextAssembler_SkipsMarkdownBrainForSpecialistBots(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/embed":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"embeddings": [][]float32{{0.1, 0.2, 0.3}},
			})
		case "/collections/conversation_memory/points/search":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": []map[string]any{{
					"id":    "mem-1",
					"score": 0.88,
					"payload": map[string]any{
						"text":             "Memória de vendas",
						"canonical_bot_id": "ac-vendas",
					},
				}},
			})
		case "/collections/aurelia_markdown_brain/points/search":
			t.Fatal("markdown brain should not be queried for specialist bots")
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	assembler := NewContextAssembler(server.URL, "", "conversation_memory", "aurelia_markdown_brain", "nomic-embed-text", server.URL, nil)
	got := assembler.AssembleContextForBot(context.Background(), "ac-vendas", "proposta")

	if strings.Contains(got, "Markdown Brain") {
		t.Fatalf("did not expect markdown brain section, got %q", got)
	}
}
