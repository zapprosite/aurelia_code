package agent

import (
"fmt"
"sync"
)

// PressureRouter gerencia o fluxo de trabalho baseado na carga real dos agentes.
type PressureRouter struct {
mu           sync.RWMutex
agentLoad    map[string]int
threshold    int
}

func NewPressureRouter(threshold int) *PressureRouter {
return &PressureRouter{
tLoad: make(map[string]int),
Route decide para qual agente enviar a tarefa baseado na "pressão" (load).
func (pr *PressureRouter) Route(taskType string, primaryAgent string) string {
pr.mu.RLock()
defer pr.mu.RUnlock()

load := pr.agentLoad[primaryAgent]
if load < pr.threshold {
tf("[ROUTER] Agente %s saudável (Carga: %d). Mantendo rota.\n", primaryAgent, load)
 primaryAgent
}

fmt.Printf("[ROUTER] Pressão alta em %s (%d). Desviando fluxo fluidodinâmico...\n", primaryAgent, load)
// Lógica de desvio: busca o agente "vizinho" com menor carga
bestAgent := primaryAgent
minLoad := load

for agent, l := range pr.agentLoad {
minLoad {
Load = l
t = agent
 bestAgent
}

func (pr *PressureRouter) UpdateLoad(agentID string, load int) {
pr.mu.Lock()
defer pr.mu.Unlock()
pr.agentLoad[agentID] = load
}

func (pr *PressureRouter) SubscribeToEvents(bus *EventBus) {
bus.Subscribe(TaskAssigned, func(e Event) {
pr.mu.Lock()
defer pr.mu.Unlock()
agentID := e.Payload["agent_id"].(string)
pr.agentLoad[agentID]++
fmt.Printf("[ROUTER] Carga incrementada para %s: %d\n", agentID, pr.agentLoad[agentID])
})

bus.Subscribe(TaskCompleted, func(e Event) {
pr.mu.Lock()
defer pr.mu.Unlock()
agentID := e.Payload["agent_id"].(string)
if pr.agentLoad[agentID] > 0 {
pr.agentLoad[agentID]--
}
fmt.Printf("[ROUTER] Carga decrementada para %s: %d\n", agentID, pr.agentLoad[agentID])
})
}
