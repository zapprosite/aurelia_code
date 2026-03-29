// Package computer_use provides autonomous computer use agent.
package computer_use

import (
	"encoding/json"
)

// MCPTool represents a tool in MCP format.
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
}

// MCPTools returns all browser tools in MCP format.
func MCPTools() []MCPTool {
	return []MCPTool{
		// Navegação
		{
			Name:        "navigate",
			Description: "Navega o browser para URL específica",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "URL completa para navegar",
					},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "go_back",
			Description: "Volta no histórico do browser",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "go_forward",
			Description: "Avança no histórico do browser",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "reload",
			Description: "Recarrega a página atual",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},

		// Ações de mouse
		{
			Name:        "mouse_move",
			Description: "Move o cursor para coordenadas normalizadas (0-999)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"x": map[string]interface{}{
						"type":        "number",
						"description": "Coordenada X normalizada (0-999)",
						"minimum":     0,
						"maximum":     999,
					},
					"y": map[string]interface{}{
						"type":        "number",
						"description": "Coordenada Y normalizada (0-999)",
						"minimum":     0,
						"maximum":     999,
					},
				},
				"required": []string{"x", "y"},
			},
		},
		{
			Name:        "mouse_click",
			Description: "Clique do mouse em coordenadas normalizadas",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"x": map[string]interface{}{
						"type": "number",
					},
					"y": map[string]interface{}{
						"type": "number",
					},
					"button": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"left", "right", "middle"},
						"description": "Botão do mouse",
					},
				},
				"required": []string{"x", "y"},
			},
		},

		// Teclado
		{
			Name:        "key_type",
			Description: "Digita texto em elemento focado",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "Texto para digitar",
					},
				},
				"required": []string{"text"},
			},
		},
		{
			Name:        "key_press",
			Description: "Pressiona tecla especial",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key": map[string]interface{}{
						"type":        "string",
						"description": "Nome da tecla (enter, escape, tab, etc)",
					},
				},
				"required": []string{"key"},
			},
		},

		// Scroll
		{
			Name:        "scroll",
			Description: "Rola a página",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"direction": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"up", "down", "left", "right"},
						"description": "Direção do scroll",
					},
					"amount": map[string]interface{}{
						"type":        "number",
						"description": "Quantidade em pixels",
					},
				},
				"required": []string{"direction"},
			},
		},

		// Screenshot
		{
			Name:        "screenshot",
			Description: "Captura screenshot da tela atual",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"region": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"x":      map[string]interface{}{"type": "number"},
							"y":      map[string]interface{}{"type": "number"},
							"width":  map[string]interface{}{"type": "number"},
							"height": map[string]interface{}{"type": "number"},
						},
					},
				},
			},
		},

		// Extração
		{
			Name:        "extract",
			Description: "Extrai dados da página atual",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Query em linguagem natural",
					},
				},
				"required": []string{"query"},
			},
		},

		// Done
		{
			Name:        "done",
			Description: "Finaliza tarefa e retorna resultado",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"summary": map[string]interface{}{
						"type":        "string",
						"description": "Resumo do que foi feito",
					},
				},
			},
		},
	}
}

// ToJSON returns tools as JSON for MCP.
func (t MCPTool) ToJSON() ([]byte, error) {
	return json.MarshalIndent(t, "", "  ")
}

// AllToolsJSON returns all tools as formatted JSON.
func AllToolsJSON() string {
	tools := MCPTools()
	b, _ := json.MarshalIndent(tools, "", "  ")
	return string(b)
}

// ToolNames returns all tool names.
func ToolNames() []string {
	tools := MCPTools()
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.Name
	}
	return names
}
