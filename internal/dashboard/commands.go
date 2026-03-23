package dashboard

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kocar/aurelia/internal/plan"
)

// CommandRequest representa um comando vindo do Cockpit (Frontend)
type CommandRequest struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

// HandleCommands processa as requisições do Dashboard
func HandleCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Action {
	case "set_model":
		model := req.Params["model"]
		Publish(Event{
			Type:      "agent_update",
			Agent:     "System",
			Action:    "Modelo alterado para " + model,
			Payload:   map[string]string{"model": model},
			Timestamp: time.Now().Format("15:04:05"),
		})

	case "approve_plan", "reject_plan":
		planID := req.Params["plan_id"]
		status := "approved"
		if req.Action == "reject_plan" {
			status = "rejected"
		}
		
		// Atualizar o Store global (Acessando pacote neutro para evitar ciclo)
		plan.GlobalPlanStore.UpdateStatus(planID, status)

		// Notificar Cockpit e Agente via SSE
		Publish(Event{
			Type:    "agent_plan",
			Agent:   "User",
			Action:  "Plano " + status,
			Payload: map[string]string{"id": planID, "status": status},
		})

		// Evento de Segurança no log
		Publish(Event{
			Type:    "security_alert",
			Agent:   "Dashboard",
			Action:  "Plan Authority Override",
			Payload: "Plano " + planID + " marcado como " + status,
			Timestamp: time.Now().Format("15:04:05"),
		})

	case "gateway_status":
		// Retornar status atual do gateway (chamada ao endpoint /api/router/status)
		Publish(Event{
			Type:      "gateway_status_request",
			Agent:     "Dashboard",
			Action:    "Consultando status do gateway de roteamento...",
			Timestamp: time.Now().Format("15:04:05"),
		})

	default:
		http.Error(w, "Unknown action", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
