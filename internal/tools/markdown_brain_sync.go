package tools

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/markdownbrain"
)

// MarkdownBrainSyncTool triggers synchronization of repository and optional
// vault markdown into the canonical Markdown Brain collection.
type MarkdownBrainSyncTool struct {
	syncer *markdownbrain.Syncer
}

func NewMarkdownBrainSyncTool(
	repoRoot, vaultPath, ollamaURL, embedModel, qdrantURL, qdrantAPIKey, collection string,
	db *sql.DB,
	logger *slog.Logger,
) *MarkdownBrainSyncTool {
	if strings.TrimSpace(repoRoot) == "" && strings.TrimSpace(vaultPath) == "" {
		return nil
	}
	if db == nil || strings.TrimSpace(qdrantURL) == "" || strings.TrimSpace(ollamaURL) == "" {
		return nil
	}
	if err := markdownbrain.InitSchema(db); err != nil {
		slog.Default().Warn("markdown brain: failed to init schema", slog.Any("err", err))
		return nil
	}
	syncer := markdownbrain.NewSyncer(repoRoot, vaultPath, ollamaURL, embedModel, qdrantURL, qdrantAPIKey, collection, db, logger)
	return &MarkdownBrainSyncTool{syncer: syncer}
}

func (t *MarkdownBrainSyncTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "markdown_brain_sync",
		Description: "Sincroniza todo o cerebro Markdown canonico da Aurelia para o Qdrant, incluindo documentos .md do repositorio e notas do vault quando configurado.",
		JSONSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
			"required":   []string{},
		},
	}
}

func (t *MarkdownBrainSyncTool) Sync(ctx context.Context) (markdownbrain.SyncStats, error) {
	if t == nil || t.syncer == nil {
		return markdownbrain.SyncStats{}, fmt.Errorf("markdown brain sync nao configurado")
	}
	return t.syncer.Sync(ctx)
}

func (t *MarkdownBrainSyncTool) Execute(ctx context.Context, _ map[string]any) (string, error) {
	stats, err := t.Sync(ctx)
	if err != nil {
		return "", fmt.Errorf("markdown brain sync: %w", err)
	}
	if !stats.Changed() {
		return fmt.Sprintf(
			"Markdown Brain sincronizado: sem mudancas. Repo=%d documento(s), Vault=%d documento(s).",
			stats.RepoDocs, stats.VaultDocs,
		), nil
	}
	return fmt.Sprintf(
		"Markdown Brain sincronizado: %d documento(s) reindexado(s), %d chunk(s) atualizados e %d documento(s) removido(s). Repo=%d, Vault=%d.",
		stats.SyncedDocs, stats.SyncedChunks, stats.RemovedDocs, stats.RepoDocs, stats.VaultDocs,
	), nil
}
