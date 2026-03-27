package skill

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSemanticRouter_SyncSkills_WritesCanonicalPayload(t *testing.T) {
	t.Parallel()

	var pointsPayload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/embed":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"embeddings":[[0.1,0.2,0.3]]}`))
		case "/collections/test_skills":
			w.WriteHeader(http.StatusOK)
		case "/collections/test_skills/points":
			if err := json.NewDecoder(r.Body).Decode(&pointsPayload); err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	router := NewSemanticRouter(server.URL, "secret", "test_skills", "nomic-embed-text", server.URL)
	router.embedURL = server.URL + "/api/embed"

	err := router.SyncSkills(context.Background(), map[string]Skill{
		"controle-db": {
			Metadata: Metadata{
				Description: "governa supabase qdrant obsidian e sqlite",
			},
		},
	})
	if err != nil {
		t.Fatalf("SyncSkills() error = %v", err)
	}

	points, ok := pointsPayload["points"].([]any)
	if !ok || len(points) != 1 {
		t.Fatalf("unexpected points payload: %#v", pointsPayload)
	}
	point, ok := points[0].(map[string]any)
	if !ok {
		t.Fatalf("unexpected point shape: %#v", points[0])
	}
	payload, ok := point["payload"].(map[string]any)
	if !ok {
		t.Fatalf("unexpected payload shape: %#v", point["payload"])
	}

	required := []string{"app_id", "repo_id", "environment", "text", "name", "description", "source_system", "source_id", "domain", "ts", "version"}
	for _, key := range required {
		if payload[key] == nil || payload[key] == "" {
			t.Fatalf("missing canonical key %q in %#v", key, payload)
		}
	}
	if payload["section"] == nil || payload["chunk_id"] == nil || payload["checksum"] == nil {
		t.Fatalf("missing chunk metadata in %#v", payload)
	}
	if payload["source_system"] != "skills" || payload["source_id"] != "skill:controle-db" {
		t.Fatalf("unexpected source lineage payload: %#v", payload)
	}
}
