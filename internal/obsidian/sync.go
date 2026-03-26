package obsidian

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Syncer indexes Obsidian vault notes into Qdrant via Ollama embeddings.
// sync state is tracked in SQLite to avoid re-embedding unchanged files.
type Syncer struct {
	vaultPath    string
	ollamaURL    string
	embedModel   string
	qdrantURL    string
	qdrantAPIKey string
	collection   string
	db           *sql.DB
	httpClient   *http.Client
	logger       *slog.Logger
}

// NewSyncer creates an Obsidian syncer. db must have obsidian_sync_state table (see InitSchema).
func NewSyncer(vaultPath, ollamaURL, embedModel, qdrantURL, qdrantAPIKey, collection string, db *sql.DB, logger *slog.Logger) *Syncer {
	if embedModel == "" {
		embedModel = "bge-m3"
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Syncer{
		vaultPath:    vaultPath,
		ollamaURL:    strings.TrimRight(ollamaURL, "/"),
		embedModel:   embedModel,
		qdrantURL:    strings.TrimRight(qdrantURL, "/"),
		qdrantAPIKey: qdrantAPIKey,
		collection:   collection,
		db:           db,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		logger:       logger,
	}
}

// InitSchema creates the obsidian_sync_state table if it does not exist.
func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS obsidian_sync_state (
			path        TEXT PRIMARY KEY,
			sha256      TEXT NOT NULL,
			last_synced DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// Sync reads the vault and indexes new or changed notes into Qdrant.
// Returns the number of notes upserted.
func (s *Syncer) Sync(ctx context.Context) (int, error) {
	notes, err := ReadVault(s.vaultPath)
	if err != nil {
		return 0, fmt.Errorf("obsidian: read vault: %w", err)
	}

	upserted := 0
	for _, note := range notes {
		hash := contentHash(note.Content)
		changed, err := s.hasChanged(note.RelPath, hash)
		if err != nil {
			s.logger.Warn("obsidian: sync state check failed", slog.String("path", note.RelPath), slog.Any("err", err))
			continue
		}
		if !changed {
			continue
		}

		if err := s.indexNote(ctx, note, hash); err != nil {
			s.logger.Warn("obsidian: index failed", slog.String("path", note.RelPath), slog.Any("err", err))
			continue
		}
		upserted++
	}

	if upserted > 0 {
		s.logger.Info("obsidian: sync complete", slog.Int("upserted", upserted), slog.Int("total", len(notes)))
	}
	return upserted, nil
}

func (s *Syncer) indexNote(ctx context.Context, note VaultNote, hash string) error {
	text := buildIndexText(note)
	if strings.TrimSpace(text) == "" {
		return nil
	}

	vec, err := s.embed(ctx, text)
	if err != nil {
		return fmt.Errorf("embed: %w", err)
	}

	payload := map[string]any{
		"app_id":           "aurelia",
		"repo_id":          "github.com/kocar/aurelia",
		"environment":      "production",
		"text":             text,
		"canonical_bot_id": "aurelia_code",
		"source_system":    "obsidian",
		"source_id":        "obsidian:" + note.RelPath,
		"domain":           "knowledge",
		"ts":               note.ModifiedAt.Unix(),
		"version":          1,
		"title":            note.Title,
		"tags":             strings.Join(note.Tags, ","),
		"vault_path":       note.RelPath,
	}

	if err := s.upsertPoint(ctx, note.RelPath, vec, payload); err != nil {
		return fmt.Errorf("qdrant upsert: %w", err)
	}

	return s.recordSync(note.RelPath, hash)
}

func (s *Syncer) embed(ctx context.Context, text string) ([]float32, error) {
	body, _ := json.Marshal(map[string]any{"model": s.embedModel, "input": text})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.ollamaURL+"/api/embed", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ollama HTTP %s", resp.Status)
	}

	var payload struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload.Embeddings) == 0 {
		return nil, fmt.Errorf("ollama returned empty embeddings")
	}
	return payload.Embeddings[0], nil
}

func (s *Syncer) upsertPoint(ctx context.Context, relPath string, vec []float32, payload map[string]any) error {
	body, err := json.Marshal(map[string]any{
		"points": []map[string]any{
			{
				"id":      deterministicID(relPath),
				"vector":  vec,
				"payload": payload,
			},
		},
	})
	if err != nil {
		return err
	}

	url := s.qdrantURL + "/collections/" + s.collection + "/points?wait=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.qdrantAPIKey != "" {
		req.Header.Set("api-key", s.qdrantAPIKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qdrant: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("qdrant HTTP %s", resp.Status)
	}
	return nil
}

func (s *Syncer) hasChanged(relPath, hash string) (bool, error) {
	var stored string
	err := s.db.QueryRow(
		`SELECT sha256 FROM obsidian_sync_state WHERE path = ?`, relPath,
	).Scan(&stored)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return stored != hash, nil
}

func (s *Syncer) recordSync(relPath, hash string) error {
	_, err := s.db.Exec(`
		INSERT INTO obsidian_sync_state (path, sha256, last_synced)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(path) DO UPDATE SET sha256=excluded.sha256, last_synced=excluded.last_synced
	`, relPath, hash)
	return err
}

func buildIndexText(note VaultNote) string {
	parts := []string{}
	if note.Title != "" {
		parts = append(parts, note.Title)
	}
	if len(note.Tags) > 0 {
		parts = append(parts, strings.Join(note.Tags, " "))
	}
	if note.Content != "" {
		// Limit to first 2000 chars to keep embedding cost reasonable
		content := note.Content
		if len(content) > 2000 {
			content = content[:2000]
		}
		parts = append(parts, content)
	}
	return strings.Join(parts, "\n")
}

func contentHash(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])
}

// deterministicID converts a relative vault path to a stable uint64 for Qdrant point IDs.
// Uses the first 8 bytes of sha256 as a uint64.
func deterministicID(relPath string) uint64 {
	h := sha256.Sum256([]byte(relPath))
	var id uint64
	for i := 0; i < 8; i++ {
		id = (id << 8) | uint64(h[i])
	}
	return id
}
