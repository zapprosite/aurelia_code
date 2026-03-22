package agent

import (
	"context"
	"fmt"
	"time"
	"log/slog"

	"github.com/kocar/aurelia/internal/dashboard"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/plan"
)

// LogVerificationTool permite ao agente registrar o sucesso de uma tarefa
func LogVerificationTool(ctx context.Context, args map[string]interface{}) (string, error) {
	logger := observability.Logger("agent.verifier")

	planID, _ := args["plan_id"].(string)
	evidence, _ := args["evidence"].(string)
	success, _ := args["success"].(bool)

	// Atualizar status no Store
	status := "completed"
	if !success {
		status = "failed"
	}
	plan.GlobalPlanStore.UpdateStatus(planID, status)

	// Publicar evento de verificação no Dashboard
	dashboard.Publish(dashboard.Event{
		Type:      "agent_verification",
		Agent:     "Aurelia",
		Action:    "Verificação Final: " + status,
		Payload:   map[string]interface{}{
			"plan_id":  planID,
			"evidence": evidence,
			"success":  success,
		},
		Timestamp: time.Now().Format("15:04:05"),
	})

	logger.Info("verification logged", 
		slog.String("plan_id", planID),
		slog.Bool("success", success))

	return fmt.Sprintf("Verificação registrada para o plano %s. Status final: %s.", planID, status), nil
}

// RegisterVerifierTools registra as ferramentas de verificação
func (r *ToolRegistry) RegisterVerifierTools() {
	r.Register(Tool{
		Name:        "log_verification",
		Description: "Registra os resultados da fase de verificação após a execução de um plano. OBRIGATÓRIO para fechar o ciclo de trabalho.",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"plan_id":  map[string]interface{}{"type": "string", "description": "ID do plano que foi executado"},
				"evidence": map[string]interface{}{"type": "string", "description": "Descrição do que foi testado/validado"},
				"success":  map[string]interface{}{"type": "boolean", "description": "Indica se o objetivo foi atingido"},
			},
			"required": []string{"plan_id", "evidence", "success"},
		},
	}, LogVerificationTool)
}
