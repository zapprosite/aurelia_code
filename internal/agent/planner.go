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

// ProposePlanTool permite que o agente apresente um plano de ação para aprovação
func ProposePlanTool(ctx context.Context, args map[string]interface{}) (string, error) {
	logger := observability.Logger("agent.planner")

	// Extrair argumentos de forma segura
	var p plan.ActionPlan
	p.Title, _ = args["title"].(string)
	p.Description, _ = args["description"].(string)
	p.RiskLevel, _ = args["risk_level"].(string)
	p.EstimatedTime, _ = args["estimated_time"].(string)
	p.BackoutPlan, _ = args["backout_plan"].(string)
	
	if stepsRaw, ok := args["steps"].([]interface{}); ok {
		for _, s := range stepsRaw {
			if str, ok := s.(string); ok {
				p.Steps = append(p.Steps, str)
			}
		}
	}

	p.ID = fmt.Sprintf("plan_%d", time.Now().Unix())
	p.CreatedAt = time.Now()
	p.Status = "proposed"

	// Registrar no Store global para bloqueio (Hard-Gate)
	plan.GlobalPlanStore.Add(&p)

	// Publicar para o Dashboard via SSE
	dashboard.Publish(dashboard.Event{
		Type:      "agent_plan",
		Agent:     "Aurelia",
		Action:    "Plano Proposto: " + p.Title,
		Payload:   p,
		Timestamp: time.Now().Format("15:04:05"),
	})

	logger.Info("action plan proposed", 
		slog.String("id", p.ID),
		slog.String("risk", p.RiskLevel))

	return fmt.Sprintf("Plano '%s' (ID: %s) foi enviado para o Cockpit. AGUARDE a aprovação humana antes de proceder para EXECUTION.", p.Title, p.ID), nil
}

// RegisterPlannerTools registra as ferramentas de planejamento
func (r *ToolRegistry) RegisterPlannerTools() {
	r.Register(Tool{
		Name:        "propose_plan",
		Description: "Propõe um plano de ação detalhado para revisão humana no Cockpit. Use antes de tarefas complexas ou de alto risco.",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title":          map[string]interface{}{"type": "string", "description": "Resumo curto da missão"},
				"description":    map[string]interface{}{"type": "string", "description": "Explicação do que será feito"},
				"risk_level":     map[string]interface{}{"type": "string", "enum": []string{"Low", "Medium", "High", "Critical"}},
				"steps":          map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
				"estimated_time": map[string]interface{}{"type": "string", "description": "Ex: 15 min"},
				"backout_plan":   map[string]interface{}{"type": "string", "description": "Como reverter em caso de falha"},
			},
			"required": []string{"title", "description", "risk_level", "steps", "backout_plan"},
		},
	}, ProposePlanTool)
}
