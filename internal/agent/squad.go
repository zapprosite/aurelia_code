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
			Role:     "Arquiteta / Governança",
			Status:   "online",
			Load:     12,
			Color:    "text-purple-400",
			IconName: "Shield",
		},
		{
			ID:       "claude",
			Name:     "Claude 5 Omni",
			Role:     "Implementador Principal",
			Status:   "online",
			Load:     0,
			Color:    "text-orange-400",
			IconName: "Brain",
		},
		{
			ID:       "codex",
			Name:     "Codex",
			Role:     "Executor Rápido",
			Status:   "offline",
			Load:     0,
			Color:    "text-emerald-400",
			IconName: "Terminal",
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
