package agent

import (
	"context"
	"encoding/json"
	"fmt"
)

// HandoffRequest define os parâmetros para transferir o controle para outro agente
type HandoffRequest struct {
	TargetAgent     string `json:"target_agent"`
	TaskDescription string `json:"task_description"`
	Reason          string `json:"reason"`
}

// HandoffResult é o que a tool retorna para o Loop
type HandoffResult struct {
	Success     bool   `json:"success"`
	NextAgent   string `json:"next_agent"`
	Instruction string `json:"instruction"`
}

// GetHandoffToolDefinition retorna a definição da ferramenta para o LLM
func GetHandoffToolDefinition() Tool {
	return Tool{
		Name:        "handoff_to_agent",
		Description: "Delegue o controle da conversa para outro agente especializado (ex: coder, reviewer, researcher). O histórico atual será preservado.",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"target_agent": map[string]interface{}{
					"type":        "string",
					"description": "Nome do agente destino (ex: 'coder', 'reviewer')",
				},
				"task_description": map[string]interface{}{
					"type":        "string",
					"description": "O que o próximo agente deve fazer especificamente",
				},
				"reason": map[string]interface{}{
					"type":        "string",
					"description": "Por que você está transferindo o controle agora",
				},
			},
			"required": []string{"target_agent", "task_description"},
		},
	}
}

// HandoffHandler retorna a função de execução da ferramenta
func HandoffHandler(service *MasterTeamService) func(context.Context, map[string]interface{}) (string, error) {
	return func(ctx context.Context, args map[string]interface{}) (string, error) {
		var req HandoffRequest
		data, _ := json.Marshal(args)
		if err := json.Unmarshal(data, &req); err != nil {
			return "", err
		}

		err := service.HandleDirectHandoff(ctx, req)
		if err != nil {
			return "", fmt.Errorf("falha no handoff nativo: %w", err)
		}

		res := HandoffResult{
			Success:     true,
			NextAgent:   req.TargetAgent,
			Instruction: "Handoff iniciado com sucesso. O próximo agente assumirá a partir daqui.",
		}
		out, _ := json.Marshal(res)
		return string(out), nil
	}
}
