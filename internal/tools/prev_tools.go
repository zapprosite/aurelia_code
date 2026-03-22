package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/dashboard"
	"github.com/kocar/aurelia/internal/observability"
)

type SetPhaseTool struct{}

func NewSetPhaseTool() *SetPhaseTool {
	return &SetPhaseTool{}
}

func (t *SetPhaseTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "set_phase",
		Description: "Avanca o loop da Aurelia para a proxima fase da metodologia PREV (PLANNING, REVIEW, EXECUTION, VERIFICATION). Use isso para sinalizar que voce completou a pesquisa (PLANNING) e iniciara a escrita de codigo ou testes pesados.",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"phase": map[string]interface{}{
					"type":        "string",
					"description": "Fase alvo (PLANNING, REVIEW, EXECUTION, VERIFICATION)",
					"enum":        []string{"PLANNING", "REVIEW", "EXECUTION", "VERIFICATION"},
				},
				"reasoning": map[string]interface{}{
					"type":        "string",
					"description": "Explicacao de por que voce esta avancando para esta fase e qual e o plano macro",
				},
			},
			"required": []string{"phase", "reasoning"},
		},
	}
}

func (t *SetPhaseTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	phase, ok := args["phase"].(string)
	if !ok {
		return "", fmt.Errorf("phase is required and must be a string")
	}
	reasoning, _ := args["reasoning"].(string)

	logger := observability.Logger("tools.prev")
	logger.Info("agent requested phase change", "phase", phase, "reasoning", reasoning)

	agentName, _ := agent.AgentContextFromContext(ctx)
	if agentName == "" {
		agentName = "Aurelia"
	}

	payload := map[string]string{
		"phase":     phase,
		"reasoning": reasoning,
	}

	dashboard.Publish(dashboard.Event{
		Type:      "agent_phase",
		Agent:     agentName,
		Action:    "Transitioning to " + phase,
		Payload:   payload,
		Timestamp: time.Now().Format(time.Kitchen),
	})

	return fmt.Sprintf("Phase updated to %s. Reason recorded: %s", phase, reasoning), nil
}
