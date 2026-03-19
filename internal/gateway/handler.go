package gateway

import (
	"encoding/json"
	"net/http"
)

func NewDryRunHandler(planner *Planner) http.Handler {
	if planner == nil {
		planner = NewPlanner()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req DryRunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		if req.Task == "" && req.TaskClass == "" {
			http.Error(w, "task or task_class is required", http.StatusBadRequest)
			return
		}

		decision := planner.Decide(req)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(decision)
	})
}
