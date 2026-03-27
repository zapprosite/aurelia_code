package markdownbrain

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSyncer_SyncsRepositoryMarkdownAndPurgesMissingDocuments(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	mustWriteFile(t, filepath.Join(repoRoot, "docs", "brain.md"), "# Brain\n\nInitial content")

	dbPath := filepath.Join(t.TempDir(), "brain.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer db.Close()
	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}

	var deleteCalls int
	var upsertPoints int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/embed":
			_, _ = w.Write([]byte(`{"embeddings":[[0.1,0.2,0.3]]}`))
		case "/collections/aurelia_markdown_brain":
			w.WriteHeader(http.StatusOK)
		case "/collections/aurelia_markdown_brain/points/delete":
			deleteCalls++
			w.WriteHeader(http.StatusOK)
		case "/collections/aurelia_markdown_brain/points":
			var payload struct {
				Points []map[string]any `json:"points"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode upsert payload: %v", err)
			}
			upsertPoints += len(payload.Points)
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	syncer := NewSyncer(repoRoot, "", server.URL, "nomic-embed-text", server.URL, "", DefaultCollection, db, nil)
	stats, err := syncer.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}
	if stats.RepoDocs != 1 || stats.SyncedDocs != 1 || upsertPoints == 0 {
		t.Fatalf("unexpected first sync stats: %#v upsert=%d", stats, upsertPoints)
	}

	if err := os.Remove(filepath.Join(repoRoot, "docs", "brain.md")); err != nil {
		t.Fatalf("remove doc: %v", err)
	}
	stats, err = syncer.Sync(context.Background())
	if err != nil {
		t.Fatalf("second Sync() error = %v", err)
	}
	if stats.RemovedDocs != 1 {
		t.Fatalf("expected one removed doc, got %#v", stats)
	}
	if deleteCalls == 0 {
		t.Fatal("expected qdrant delete to be called")
	}
}
