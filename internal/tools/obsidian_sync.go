package tools

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/obsidian"
)

// ObsidianSyncTool is an agent-callable tool that triggers an Obsidian vault sync to Qdrant.
type ObsidianSyncTool struct {
	syncer *obsidian.Syncer
}

// NewObsidianSyncTool creates the tool. Returns nil if ObsidianSyncEnabled is false or vaultPath is empty.
func NewObsidianSyncTool(
	vaultPath, ollamaURL, embedModel, qdrantURL, qdrantAPIKey, collection string,
	db *sql.DB,
	logger *slog.Logger,
) *ObsidianSyncTool {
	if vaultPath == "" {
		return nil
	}
	if err := obsidian.InitSchema(db); err != nil {
		slog.Default().Warn("obsidian: failed to init schema", slog.Any("err", err))
		return nil
	}
	syncer := obsidian.NewSyncer(vaultPath, ollamaURL, embedModel, qdrantURL, qdrantAPIKey, collection, db, logger)
	return &ObsidianSyncTool{syncer: syncer}
}

func (t *ObsidianSyncTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "obsidian_sync",
		Description: "Sincroniza o vault do Obsidian com o índice vetorial do Qdrant. Lê arquivos .md, gera embeddings via Ollama e indexa apenas notas novas ou modificadas.",
		JSONSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
			"required":   []string{},
		},
	}
}

func (t *ObsidianSyncTool) Execute(ctx context.Context, _ map[string]any) (string, error) {
	if t.syncer == nil {
		return "", fmt.Errorf("obsidian sync não configurado")
	}
	n, err := t.syncer.Sync(ctx)
	if err != nil {
		return "", fmt.Errorf("obsidian sync: %w", err)
	}
	if n == 0 {
		return "Obsidian vault sincronizado: nenhuma nota nova ou modificada.", nil
	}
	return fmt.Sprintf("Obsidian vault sincronizado: %d nota(s) indexada(s) no Qdrant.", n), nil
}
