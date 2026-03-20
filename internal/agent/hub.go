package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Hub struct {
	pressureRouter  *PressureRouter
	contractManager *ContractManager
	mu              sync.RWMutex
	connections     map[string]http.ResponseWriter
}

func NewHub(pr *PressureRouter, cm *ContractManager) *Hub {
	return &Hub{
		connections:     make(map[string]http.ResponseWriter),
		pressureRouter:  pr,
		contractManager: cm,
	}
}

func (h *Hub) Broadcast(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("[HUB-ERROR] Marshal fail: %v\n", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	// In a real production app, we would loop through websocket connections here.
	// For this 2026 Swarm Proof, we output to the system log which is monitored by the Cognitive AR Overlay.
	fmt.Printf("[DASHBOARD-WS-BROADCAST] Kind=%s Payload=%s\n", event.Kind, string(data))
	
	// Bonus: Include Pressure Metrics in periodic Heartbeats
	if event.Kind == EventTaskCompleted {
		h.BroadcastMetrics()
	}
}

func (h *Hub) BroadcastMetrics() {
	metrics := map[string]interface{}{
		"type": "SWARM_METRICS",
		"pressure": h.pressureRouter.agentLoad,
		"reputation": h.contractManager.reputation,
	}
	data, _ := json.Marshal(metrics)
	fmt.Printf("[DASHBOARD-WS-METRICS] %s\n", string(data))
}

func (s *MasterTeamService) StartDashboard(ctx context.Context) {
	h := s.hub()
	s.bus.Subscribe(EventTaskCreated, func(e Event) { h.Broadcast(e) })
	s.bus.Subscribe(EventHelpRequested, func(e Event) { h.Broadcast(e) })
	s.bus.Subscribe(EventHelpOffered, func(e Event) { h.Broadcast(e) })
	s.bus.Subscribe(EventTaskCompleted, func(e Event) { h.Broadcast(e) })
}

func (s *MasterTeamService) hub() *Hub {
	// Simple singleton-like access for the service
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.hubInstance == nil {
		// Mock managers for now if not injected
		pr := NewPressureRouter(nil)
		cm := NewContractManager()
		s.hubInstance = NewHub(pr, cm)
	}
	return s.hubInstance
}
