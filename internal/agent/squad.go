package agent

// SquadMember represents a fixed member of the sovereign agentic squad.
type SquadMember struct {
	Name     string
	IconName string
	Role     string
	Status   string // "online", "busy", "offline"
	Load     int    // percentage 0-100
}

// GetFixedSquad returns the default squad configuration for the Aurelia ecosystem.
// This is used for status reporting and initial coordination.
func GetFixedSquad() []SquadMember {
	return []SquadMember{
		{
			Name:     "Aurelia",
			IconName: "crown",
			Role:     "Lead Coordinator",
			Status:   "online",
			Load:     0,
		},
		{
			Name:     "Sentinel",
			IconName: "shield",
			Role:     "Security & Guardrails",
			Status:   "online",
			Load:     0,
		},
		{
			Name:     "Researcher",
			IconName: "magnifying-glass",
			Role:     "Deep Knowledge Discovery",
			Status:   "online",
			Load:     0,
		},
		{
			Name:     "Coder",
			IconName: "code",
			Role:     "System Implementation",
			Status:   "online",
			Load:     0,
		},
	}
}

// UpdateSquadAgentStatus updates the status and load of a squad member in memory.
// In SOTA 2026, this satisfies health check reporting.
func UpdateSquadAgentStatus(agentName, status string, load int) error {
	// For now, this is a no-op that just satisfies the build.
	// Future: update a global shared state for the dashboard/telemetry.
	return nil
}
