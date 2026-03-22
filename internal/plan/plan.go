package plan

import (
	"sync"
	"time"
)

// ActionPlan representa uma estratégia estruturada de execução
type ActionPlan struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	RiskLevel     string    `json:"risk_level"` // Low, Medium, High, Critical
	Steps         []string  `json:"steps"`
	EstimatedTime string    `json:"estimated_time"`
	BackoutPlan   string    `json:"backout_plan"`
	Status        string    `json:"status"` // proposed, approved, rejected, completed, failed
	CreatedAt     time.Time `json:"created_at"`
}

// PlanStore gerencia o estado global de aprovações de planos
type PlanStore struct {
	mu    sync.RWMutex
	plans map[string]*ActionPlan
}

func NewPlanStore() *PlanStore {
	return &PlanStore{
		plans: make(map[string]*ActionPlan),
	}
}

var GlobalPlanStore = NewPlanStore()

func (s *PlanStore) Add(p *ActionPlan) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.plans[p.ID] = p
}

func (s *PlanStore) UpdateStatus(id string, status string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := s.plans[id]; ok {
		p.Status = status
		return true
	}
	return false
}

func (s *PlanStore) Get(id string) *ActionPlan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.plans[id]
}

func (s *PlanStore) HasPending() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.plans {
		if p.Status == "proposed" {
			return true
		}
	}
	return false
}
