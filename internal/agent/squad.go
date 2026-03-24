package agent

import (
	"sync"
)

// SquadAgent define a identidade fixa e o status de um membro do Squad Oficial.
type SquadAgent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Status    string `json:"status"` // "online", "offline", "busy"
	Load      int    `json:"load"`   // 0-100% de carga
	Color     string `json:"color"`  // Classe CSS de cor sugerida
	IconName  string `json:"icon"`   // Nome do ícone Lucide
}

var (
	squadMu sync.RWMutex
	
	// fixedSquad mantém a lista oficial em memória do backend Go,
	// em vez de injetar estático via hardcode no frontend.
	fixedSquad = []SquadAgent{
		{
			ID:       "aurelia",
			Name:     "Aurélia",
			Role:     "Líder · Arquiteta · Homelab",
			Status:   "online",
			Load:     12,
			Color:    "text-purple-400",
			IconName: "Crown",
		},
		{
			ID:       "sentinel",
			Name:     "Sentinel",
			Role:     "Monitor de Homelab",
			Status:   "online",
			Load:     5,
			Color:    "text-cyan-400",
			IconName: "Activity",
		},
		{
			ID:       "cronus",
			Name:     "Cronus",
			Role:     "Orquestrador de Crons",
			Status:   "online",
			Load:     3,
			Color:    "text-yellow-400",
			IconName: "Clock",
		},
		{
			ID:       "gemma",
			Name:     "Gemma3",
			Role:     "Executor Local · Ollama",
			Status:   "online",
			Load:     0,
			Color:    "text-green-400",
			IconName: "Cpu",
		},
		{
			ID:       "openrouter",
			Name:     "OpenRouter",
			Role:     "Executor Remoto · MiniMax",
			Status:   "online",
			Load:     0,
			Color:    "text-blue-400",
			IconName: "Globe",
		},
	}
)

// GetFixedSquad retorna a cópia atualizada do Squad.
func GetFixedSquad() []SquadAgent {
	squadMu.RLock()
	defer squadMu.RUnlock()

	copySquad := make([]SquadAgent, len(fixedSquad))
	copy(copySquad, fixedSquad)
	return copySquad
}

// UpdateSquadAgentStatus atualiza o status e workload de um agente fixo em memória.
// Usado pelo loop ou master_team_service para refletir atividade real.
func UpdateSquadAgentStatus(id string, status string, load int) {
	squadMu.Lock()
	defer squadMu.Unlock()

	for i := range fixedSquad {
		if fixedSquad[i].ID == id {
			fixedSquad[i].Status = status
			fixedSquad[i].Load = load
			break
		}
	}
}

// AddSquadAgent registra um novo agente no squad se não existir com esse ID.
// Permite dinâmicamente adicionar agentes spawned pelo swarm.
func AddSquadAgent(a SquadAgent) {
	squadMu.Lock()
	defer squadMu.Unlock()

	// Verificar se já existe
	for _, s := range fixedSquad {
		if s.ID == a.ID {
			return
		}
	}

	// Adicionar novo agente
	fixedSquad = append(fixedSquad, a)
}
