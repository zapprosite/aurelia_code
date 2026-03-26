package tools

import (
	"context"
	"encoding/json"
	"strings"
	"fmt"

	"github.com/kocar/aurelia/internal/agent"
)

func init() {
	// Registered by RegisterCoreTools — see definitions.go
}

// homelabStatusTool is the agent-callable wrapper around CheckHomelab.
type homelabStatusTool struct {
	ollamaURL string
	qdrantURL string
}

func newHomelabStatusTool(ollamaURL, qdrantURL string) agent.Tool {
	return agent.Tool{
		Name:        "homelab_status",
		Description: "Verifica o estado de saúde do homelab: Ollama, Qdrant, GPU (nvidia-smi), Docker e disco. Retorna status por componente (healthy/degraded/offline).",
		JSONSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
			"required":   []string{},
		},
	}
}

// RegisterHomelabTool registers the homelab_status tool with the given Ollama/Qdrant URLs.
func RegisterHomelabTool(registry *agent.ToolRegistry, ollamaURL, qdrantURL string) {
	t := &homelabStatusTool{ollamaURL: ollamaURL, qdrantURL: qdrantURL}
	registry.Register(newHomelabStatusTool(ollamaURL, qdrantURL), t.execute)
}

func (t *homelabStatusTool) execute(ctx context.Context, _ map[string]any) (string, error) {
	status := CheckHomelab(ctx, t.ollamaURL, t.qdrantURL)

	var sb strings.Builder
	sb.WriteString(status.Summary())
	sb.WriteString("\n\n")
	for _, c := range status.Checks {
		icon := "✓"
		if c.Status == "degraded" {
			icon = "⚠"
		} else if c.Status == "offline" {
			icon = "✗"
		}
		sb.WriteString(icon + " " + c.Component + ": " + c.Summary)
		if c.LatencyMS > 0 {
			sb.WriteString(" (" + fmt.Sprintf("%d", c.LatencyMS) + "ms)")
		}
		sb.WriteString("\n")
	}

	j, _ := json.Marshal(status)
	sb.WriteString("\n```json\n")
	sb.Write(j)
	sb.WriteString("\n```")

	return sb.String(), nil
}

