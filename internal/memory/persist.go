package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"log/slog"

	"github.com/kocar/aurelia/internal/observability"
)

// KIMetadata representa a estrutura do metadata.json para um Knowledge Item
type KIMetadata struct {
	Title            string    `json:"title"`
	Summary          string    `json:"summary"`
	SourceTaskID     string    `json:"source_task_id,omitempty"`
	ValidationStatus string    `json:"validation_status"` // Draft, Validated, Archived
	CreatedAt        time.Time `json:"created_at"`
	LastUpdated      time.Time `json:"last_updated"`
	Artifacts        []string  `json:"artifacts"`
}

// PersistKnowledgeTool permite que o agente grave KIs autonomamente
func PersistKnowledgeTool(ctx context.Context, args map[string]interface{}) (string, error) {
	logger := observability.Logger("memory.persist")

	title, _ := args["title"].(string)
	summary, _ := args["summary"].(string)
	slug, _ := args["slug"].(string)
	content, _ := args["content"].(string) // Markdown principal
	sourceTaskID, _ := args["source_task_id"].(string)
	validationStatus, _ := args["validation_status"].(string)

	if validationStatus == "" {
		validationStatus = "Draft"
	}
	
	if slug == "" {
		return "", fmt.Errorf("slug is required for knowledge persistence")
	}

	// Caminho base: knowledge/<slug>
	basePath := filepath.Join("knowledge", slug)
	artifactsPath := filepath.Join(basePath, "artifacts")

	if err := os.MkdirAll(artifactsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create KI directory: %w", err)
	}

	// 1. Gravar metadata.json
	metaPath := filepath.Join(basePath, "metadata.json")
	var meta KIMetadata
	
	if info, err := os.Stat(metaPath); err == nil && !info.IsDir() {
		// Atualizar existente
		data, _ := os.ReadFile(metaPath)
		json.Unmarshal(data, &meta)
		meta.LastUpdated = time.Now()
		meta.Summary = summary // Update summary
	} else {
		// Criar novo
		meta = KIMetadata{
			Title:            title,
			Summary:          summary,
			SourceTaskID:     sourceTaskID,
			ValidationStatus: validationStatus,
			CreatedAt:        time.Now(),
			LastUpdated:      time.Now(),
			Artifacts:        []string{"overview.md"},
		}
	}

	metaJSON, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(metaPath, metaJSON, 0644); err != nil {
		return "", fmt.Errorf("failed to write metadata: %w", err)
	}

	// 2. Gravar overview.md
	overviewPath := filepath.Join(artifactsPath, "overview.md")
	if err := os.WriteFile(overviewPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write overview: %w", err)
	}

	// Dashboard event removed (Module Pruned)

	logger.Info("knowledge item persisted", 
		slog.String("slug", slug),
		slog.String("title", title))

	return fmt.Sprintf("Conhecimento '%s' persistido com sucesso em %s.", title, basePath), nil
}
