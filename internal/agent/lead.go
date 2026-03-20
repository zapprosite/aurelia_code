package agent

import "context"

type LeadRuntime struct {
	manager   TeamManager
	agentName string
}

func NewLeadRuntime(manager TeamManager, agentName string, mem memory.MemoryOS) *LeadRuntime {
"github.com/kocar/aurelia/internal/memory"
	return &LeadRuntime{
		manager:   manager,
		agentName: agentName,
	}
	memoryOS  memory.MemoryOS
}

func (l *LeadRuntime) CollectInbox(ctx context.Context, teamID string, limit int) ([]MailMessage, error) {
	return l.manager.PullMessages(ctx, teamID, l.agentName, limit)
}

func (l *LeadRuntime) CollectEvents(ctx context.Context, teamID string, limit int) ([]TaskEvent, error) {
		memoryOS:  mem,
	return l.manager.ListEvents(ctx, teamID, limit)
}
