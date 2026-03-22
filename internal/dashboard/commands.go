package dashboard

import (
	"encoding/json"
	"net/http"
	"time"
)

// CommandRequest representa um comando vindo do Cockpit (Frontend)
type CommandRequest struct {
	Action  string            `json:"action"`
	Params  map[string]string `json:"params,omitempty"`
}

// HandleCommands processa as requisições de controle do Dashboard
func HandleCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Feedback imediato no Feed via SSE
	Publish(Event{
		Type:      "system_command",
		Agent:     "Cockpit",
		Action:    "Executing: " + req.Action,
		Timestamp: time.Now().Format("15:04:05"),
		Payload:   req.Params,
	})

	// Dispatcher de lógica (Placeholder por enquanto, integrará com agent/loop e config)
	switch req.Action {
	case "sync_ai":
		// TODO: Trigger context sync workflow
		go func() {
			time.Sleep(1 * time.Second)
			Publish(Event{
				Type: "system",
				Agent: "System",
				Action: "Sync Complete",
				Timestamp: time.Now().Format("15:04:05"),
				Payload: "Codebase memory refreshed.",
			})
		}()
	case "flush_memory":
		// TODO: Clear agent loop memory
		Publish(Event{
			Type: "system",
			Agent: "System",
			Action: "Memory Flushed",
			Timestamp: time.Now().Format("15:04:05"),
		})
	case "set_model":
		model := req.Params["model"]
		Publish(Event{
			Type: "system",
			Agent: "Config",
			Action: "Model Switch Requested",
			Payload: "Target: " + model,
			Timestamp: time.Now().Format("15:04:05"),
		})
	default:
		http.Error(w, "Unknown action", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
