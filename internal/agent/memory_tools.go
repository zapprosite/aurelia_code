package agent

import (
	"github.com/kocar/aurelia/internal/memory"
)

// RegisterMemoryTools registra as ferramentas de persistência de conhecimento
func (r *ToolRegistry) RegisterMemoryTools() {
	r.Register(Tool{
		Name:        "persist_knowledge",
		Description: "Grava ou atualiza um Item de Conhecimento (KI) na memória permanente do repositório.",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"slug":    map[string]interface{}{"type": "string", "description": "ID único (ex: 'auth_workflow')"},
				"title":   map[string]interface{}{"type": "string", "description": "Título legível"},
				"summary": map[string]interface{}{"type": "string", "description": "Resumo executivo"},
				"content": map[string]interface{}{"type": "string", "description": "Conteúdo markdown detalhado"},
			},
			"required": []string{"slug", "title", "summary", "content"},
		},
	}, memory.PersistKnowledgeTool)
}
