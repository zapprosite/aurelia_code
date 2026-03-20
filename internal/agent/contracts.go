package agent

import (
	"fmt"
	"sync"
)

type ContractManager struct {
	mu         sync.RWMutex
	reputation map[string]float64 // AgentID -> Score
	credits    map[string]int     // AgentID -> Balance
}

func NewContractManager() *ContractManager {
	return &ContractManager{
		reputation: make(map[string]float64),
		credits:    make(map[string]int),
	}
}

func (cm *ContractManager) RecordSuccess(agentID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.reputation[agentID] += 0.1
	if cm.reputation[agentID] > 1.0 {
		cm.reputation[agentID] = 1.0
	}
	fmt.Printf("[CONTRACT] Reputação de %s subiu para %.2f\n", agentID, cm.reputation[agentID])
}

func (cm *ContractManager) RecordFailure(agentID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.reputation[agentID] -= 0.2
	if cm.reputation[agentID] < 0 {
		cm.reputation[agentID] = 0
	}
	fmt.Printf("[CONTRACT] Reputação de %s caiu para %.2f\n", agentID, cm.reputation[agentID])
}

func (cm *ContractManager) ProcessBid(from, to string, cost int) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.credits[from] >= cost {
		cm.credits[from] -= cost
		cm.credits[to] += cost
		return true
	}
	return false
}
