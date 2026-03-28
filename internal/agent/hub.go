package agent

import (
	"github.com/kocar/aurelia/internal/purity/alog"
	"os"
	"path/filepath"
)

// Skill representa uma capacidade do agente (inspirado no ClawHub)
type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

// Hub gerencia o ecossistema de skills soberanas
type Hub struct {
	rootPath string
}

// NewHub inicia o hub apontando para o diretório de skills
func NewHub(root string) *Hub {
	return &Hub{rootPath: root}
}

// ListSkills escaneia o diretório .agent/skills (Sovereign Discovery)
func (h *Hub) ListSkills() ([]Skill, error) {
	var skills []Skill
	
	entries, err := os.ReadDir(h.rootPath)
	if err != nil {
		alog.Error("failed to read skills directory", alog.With("path", h.rootPath), alog.With("err", err))
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			skills = append(skills, Skill{
				Name: entry.Name(),
				Path: filepath.Join(h.rootPath, entry.Name()),
			})
		}
	}

	alog.Info("skills discovered", alog.With("count", len(skills)))
	return skills, nil
}

// SearchSemantic busca skills via Qdrant (The Dream SOTA 2026.1)
func (h *Hub) SearchSemantic(query string) ([]Skill, error) {
	alog.Info("searching skills semantically", alog.With("query", query))
	// No futuro: Integrar com cliente Qdrant em Go.
	// Por enquanto, delegar para a ferramenta do agente.
	return nil, nil
}
