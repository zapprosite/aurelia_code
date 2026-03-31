package markdownbrain

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/memory"
)

type SyncStats struct {
	RepoDocs     int
	VaultDocs    int
	SyncedDocs   int
	SyncedChunks int
	RemovedDocs  int
}

func (s SyncStats) Changed() bool {
	return s.SyncedDocs > 0 || s.RemovedDocs > 0
}

type Syncer struct {
	repoRoot     string
	vaultPath    string
	ollamaURL    string
	embedModel   string
	qdrantURL    string
	qdrantAPIKey string
	collection   string
	db           *sql.DB
	httpClient   *http.Client
	logger       *slog.Logger
	syncMu       sync.Mutex
	ensureOnce   sync.Once
	ensureErr    error
}

func NewSyncer(repoRoot, vaultPath, ollamaURL, embedModel, qdrantURL, qdrantAPIKey, collection string, db *sql.DB, logger *slog.Logger) *Syncer {
	if embedModel == "" {
		embedModel = "nomic-embed-text"
	}
	if collection == "" {
		collection = DefaultCollection
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Syncer{
		repoRoot:     strings.TrimSpace(repoRoot),
		vaultPath:    strings.TrimSpace(vaultPath),
		ollamaURL:    strings.TrimRight(strings.TrimSpace(ollamaURL), "/"),
		embedModel:   strings.TrimSpace(embedModel),
		qdrantURL:    strings.TrimRight(strings.TrimSpace(qdrantURL), "/"),
		qdrantAPIKey: strings.TrimSpace(qdrantAPIKey),
		collection:   strings.TrimSpace(collection),
		db:           db,
		httpClient:   memory.NewSemanticHTTPClient(30 * time.Second),
		logger:       logger,
	}
}

func InitSchema(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("markdown brain schema requires a database handle")
	}
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS markdown_brain_sync_state (
			source_system TEXT NOT NULL,
			source_path   TEXT NOT NULL,
			sha256        TEXT NOT NULL,
			last_synced   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (source_system, source_path)
		)
	`)
	return err
}

func (s *Syncer) Sync(ctx context.Context) (SyncStats, error) {
	if s == nil {
		return SyncStats{}, fmt.Errorf("markdown brain syncer is nil")
	}
	s.syncMu.Lock()
	defer s.syncMu.Unlock()
	if s.db == nil {
		return SyncStats{}, fmt.Errorf("markdown brain database is not configured")
	}
	if s.qdrantURL == "" {
		return SyncStats{}, fmt.Errorf("markdown brain qdrant url is not configured")
	}
	if s.ollamaURL == "" {
		return SyncStats{}, fmt.Errorf("markdown brain ollama url is not configured")
	}

	repoDocs, err := ReadRepository(s.repoRoot)
	if err != nil {
		return SyncStats{}, fmt.Errorf("read repository markdown: %w", err)
	}
	vaultDocs, err := ReadObsidianVault(s.vaultPath)
	if err != nil {
		return SyncStats{}, fmt.Errorf("read obsidian vault markdown: %w", err)
	}

	stats := SyncStats{
		RepoDocs:  len(repoDocs),
		VaultDocs: len(vaultDocs),
	}

	allDocs := append(repoDocs, vaultDocs...)
	removed, err := s.purgeMissingDocuments(ctx, allDocs)
	if err != nil {
		return stats, err
	}
	stats.RemovedDocs = removed

	for _, doc := range allDocs {
		hash := documentHash(doc)
		changed, err := s.hasChanged(doc.SourceSystem, doc.RelPath, hash)
		if err != nil {
			s.logger.Warn("markdown brain: sync state check failed", slog.String("source", doc.SourceSystem), slog.String("path", doc.RelPath), slog.Any("err", err))
			continue
		}
		if !changed {
			continue
		}

		chunks, err := s.indexDocument(ctx, doc)
		if err != nil {
			s.logger.Warn("markdown brain: index failed", slog.String("source", doc.SourceSystem), slog.String("path", doc.RelPath), slog.Any("err", err))
			continue
		}
		if err := s.recordSync(doc.SourceSystem, doc.RelPath, hash); err != nil {
			s.logger.Warn("markdown brain: record sync failed", slog.String("source", doc.SourceSystem), slog.String("path", doc.RelPath), slog.Any("err", err))
			continue
		}
		stats.SyncedDocs++
		stats.SyncedChunks += chunks
	}

	s.logger.Info("markdown brain: sync complete",
		slog.Int("repo_docs", stats.RepoDocs),
		slog.Int("vault_docs", stats.VaultDocs),
		slog.Int("synced_docs", stats.SyncedDocs),
		slog.Int("synced_chunks", stats.SyncedChunks),
		slog.Int("removed_docs", stats.RemovedDocs),
	)
	return stats, nil
}

func (s *Syncer) indexDocument(ctx context.Context, doc Document) (int, error) {
	chunks := ChunkDocument(doc)
	if len(chunks) == 0 {
		return 0, nil
	}

	if err := s.deleteSourcePoints(ctx, buildSourceID(doc)); err != nil {
		return 0, fmt.Errorf("delete previous chunks: %w", err)
	}

	points := make([]map[string]any, 0, len(chunks))
	for _, chunk := range chunks {
		vector, err := memory.EmbedText(ctx, s.httpClient, s.ollamaURL, s.embedModel, chunk.Text)
		if err != nil {
			return 0, fmt.Errorf("embed chunk: %w", err)
		}

		s.ensureOnce.Do(func() {
			s.ensureErr = s.ensureCollection(ctx, len(vector))
		})
		if s.ensureErr != nil {
			return 0, fmt.Errorf("ensure collection: %w", s.ensureErr)
		}

		payload := buildPayload(doc, chunk)
		if err := memory.ValidateCanonicalMemoryPayload(payload); err != nil {
			return 0, fmt.Errorf("payload contract: %w", err)
		}

		points = append(points, map[string]any{
			"id":      uuid.NewMD5(uuid.NameSpaceURL, []byte(buildSourceID(doc)+":"+chunk.ID)).String(),
			"vector":  vector,
			"payload": payload,
		})
	}

	body, err := json.Marshal(map[string]any{"points": points})
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.qdrantURL+"/collections/"+s.collection+"/points?wait=true", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.qdrantAPIKey != "" {
		req.Header.Set("api-key", s.qdrantAPIKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return 0, fmt.Errorf("qdrant upsert returned %s", resp.Status)
	}

	return len(points), nil
}

func buildPayload(doc Document, chunk Chunk) map[string]any {
	now := doc.ModifiedAt.UTC()
	payload := map[string]any{
		"app_id":           "aurelia",
		"repo_id":          "github.com/kocar/aurelia",
		"environment":      "local",
		"text":             chunk.Text,
		"canonical_bot_id": "aurelia_code",
		"source_system":    doc.SourceSystem,
		"source_kind":      doc.SourceKind,
		"source_id":        buildSourceID(doc),
		"domain":           "markdown_brain",
		"ts":               now.Unix(),
		"version":          1,
		"title":            doc.Title,
		"section":          chunk.Section,
		"chunk_id":         chunk.ID,
		"chunk_index":      chunk.Index,
		"chunk_count":      chunk.Count,
		"checksum":         chunk.Checksum,
		"synced_at":        time.Now().UTC().Format(time.RFC3339),
	}
	if len(doc.Tags) > 0 {
		payload["tags"] = strings.Join(doc.Tags, ",")
	}
	switch doc.SourceSystem {
	case repoSourceSystem:
		payload["repo_path"] = filepath.ToSlash(doc.RelPath)
	case vaultSourceSystem:
		payload["vault_path"] = filepath.ToSlash(doc.RelPath)
	}
	return payload
}

func buildSourceID(doc Document) string {
	return doc.SourceSystem + ":" + filepath.ToSlash(doc.RelPath)
}

func documentHash(doc Document) string {
	return checksumText(doc.Title + "\n" + strings.Join(doc.Tags, ",") + "\n" + doc.Content)
}

func (s *Syncer) hasChanged(sourceSystem, sourcePath, hash string) (bool, error) {
	var stored string
	err := s.db.QueryRow(
		`SELECT sha256 FROM markdown_brain_sync_state WHERE source_system = ? AND source_path = ?`,
		sourceSystem, sourcePath,
	).Scan(&stored)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return stored != hash, nil
}

func (s *Syncer) recordSync(sourceSystem, sourcePath, hash string) error {
	_, err := s.db.Exec(`
		INSERT INTO markdown_brain_sync_state (source_system, source_path, sha256, last_synced)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(source_system, source_path)
		DO UPDATE SET sha256=excluded.sha256, last_synced=excluded.last_synced
	`, sourceSystem, sourcePath, hash)
	return err
}

func (s *Syncer) purgeMissingDocuments(ctx context.Context, docs []Document) (int, error) {
	current := make(map[string]map[string]struct{})
	for _, doc := range docs {
		paths := current[doc.SourceSystem]
		if paths == nil {
			paths = make(map[string]struct{})
			current[doc.SourceSystem] = paths
		}
		paths[doc.RelPath] = struct{}{}
	}

	sourceSystems := []string{repoSourceSystem, vaultSourceSystem}
	removed := 0
	for _, sourceSystem := range sourceSystems {
		tracked, err := s.listTrackedPaths(sourceSystem)
		if err != nil {
			return removed, err
		}
		currentPaths := current[sourceSystem]
		for sourcePath := range tracked {
			if _, ok := currentPaths[sourcePath]; ok {
				continue
			}
			doc := Document{SourceSystem: sourceSystem, RelPath: sourcePath}
			if err := s.deleteSourcePoints(ctx, buildSourceID(doc)); err != nil {
				return removed, fmt.Errorf("delete stale markdown brain points %s:%s: %w", sourceSystem, sourcePath, err)
			}
			if err := s.deleteTrackedPath(sourceSystem, sourcePath); err != nil {
				return removed, err
			}
			removed++
		}
	}
	return removed, nil
}

func (s *Syncer) listTrackedPaths(sourceSystem string) (map[string]struct{}, error) {
	rows, err := s.db.Query(`SELECT source_path FROM markdown_brain_sync_state WHERE source_system = ?`, sourceSystem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	paths := make(map[string]struct{})
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths[path] = struct{}{}
	}
	return paths, rows.Err()
}

func (s *Syncer) deleteTrackedPath(sourceSystem, sourcePath string) error {
	_, err := s.db.Exec(`DELETE FROM markdown_brain_sync_state WHERE source_system = ? AND source_path = ?`, sourceSystem, sourcePath)
	return err
}

func (s *Syncer) deleteSourcePoints(ctx context.Context, sourceID string) error {
	body, err := json.Marshal(map[string]any{
		"filter": map[string]any{
			"must": []map[string]any{
				{
					"key": "source_id",
					"match": map[string]any{
						"value": sourceID,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.qdrantURL+"/collections/"+s.collection+"/points/delete?wait=true", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.qdrantAPIKey != "" {
		req.Header.Set("api-key", s.qdrantAPIKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("qdrant delete returned %s", resp.Status)
	}
	return nil
}

func (s *Syncer) ensureCollection(ctx context.Context, size int) error {
	body, err := json.Marshal(map[string]any{
		"vectors": map[string]any{
			"size":     size,
			"distance": "Cosine",
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.qdrantURL+"/collections/"+s.collection, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.qdrantAPIKey != "" {
		req.Header.Set("api-key", s.qdrantAPIKey)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusConflict {
		return nil
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil
	}
	return fmt.Errorf("qdrant ensure collection returned %s", resp.Status)
}
