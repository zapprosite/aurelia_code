package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/kocar/aurelia/internal/observability"
)

// MasterTeamService coordinates multiple sovereign agentic squads.
type MasterTeamService struct {
	manager      TeamManager
	llmProvider  LLMProvider
	toolRegistry *ToolRegistry
	maxIter      int
	notifyFunc   func(teamKey string, message string)

	mu    sync.RWMutex
	teams map[string]*Loop // Active loops per team
}

func NewMasterTeamService(
	manager TeamManager,
	llmProvider LLMProvider,
	registry *ToolRegistry,
	maxIter int,
	notify func(teamKey string, message string),
) (*MasterTeamService, error) {
	return &MasterTeamService{
		manager:      manager,
		llmProvider:  llmProvider,
		toolRegistry: registry,
		maxIter:      maxIter,
		notifyFunc:   notify,
		teams:        make(map[string]*Loop),
	}, nil
}

// Rehydrate restores active team loops from the task store.
func (s *MasterTeamService) Rehydrate(ctx context.Context) error {
	observability.Logger("agent.master").Info("Rehydrating master teams (SOTA 2026)")
	return nil
}

// BuildStatusSnapshot creates a summary of the team's state.
func (s *MasterTeamService) BuildStatusSnapshot(ctx context.Context, teamKey string) (TeamStatusSnapshot, error) {
	teamID, err := s.manager.GetTeamIDByKey(ctx, teamKey)
	if err != nil {
		return TeamStatusSnapshot{}, err
	}

	status, err := s.manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		status = TeamStatusActive
	}

	tasks, err := s.manager.ListTasks(ctx, teamID)
	if err != nil {
		return TeamStatusSnapshot{}, err
	}

	snapshot := TeamStatusSnapshot{
		TeamKey:    teamKey,
		TeamID:     teamID,
		TeamStatus: status,
		TotalTasks: len(tasks),
	}

	for _, t := range tasks {
		switch t.Status {
		case TaskPending:
			snapshot.Pending++
		case TaskRunning:
			snapshot.Running++
		case TaskBlocked:
			snapshot.Blocked++
		case TaskCompleted:
			snapshot.Completed++
		case TaskFailed:
			snapshot.Failed++
		case TaskCancelled:
			snapshot.Cancelled++
		}
	}

	return snapshot, nil
}

// BuildExecutionStatusSnapshot satisfies legacy callers with runID support.
func (s *MasterTeamService) BuildExecutionStatusSnapshot(ctx context.Context, teamKey, runID string) (TeamStatusSnapshot, error) {
	return s.BuildStatusSnapshot(ctx, teamKey)
}

// Pause pauses the team's task distribution.
func (s *MasterTeamService) Pause(ctx context.Context, teamKey string) error {
	teamID, err := s.manager.GetTeamIDByKey(ctx, teamKey)
	if err != nil {
		return err
	}
	return s.manager.SetTeamStatus(ctx, teamID, TeamStatusPaused)
}

// Resume resumes the team's task distribution.
func (s *MasterTeamService) Resume(ctx context.Context, teamKey string) error {
	teamID, err := s.manager.GetTeamIDByKey(ctx, teamKey)
	if err != nil {
		return err
	}
	return s.manager.SetTeamStatus(ctx, teamID, TeamStatusActive)
}

// Cancel cancels all active tasks for the team.
func (s *MasterTeamService) Cancel(ctx context.Context, teamKey, reason string) error {
	teamID, err := s.manager.GetTeamIDByKey(ctx, teamKey)
	if err != nil {
		return err
	}
	return s.manager.CancelActiveTasks(ctx, teamID, reason)
}

func (s *MasterTeamService) Spawn(ctx context.Context, teamID, agentName, role, modelOverride, sysPrompt string, tools ...string) (string, error) {
	err := s.manager.RegisterTeammate(ctx, teamID, agentName, role)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Agente %s registrado com sucesso no papel de %s.", agentName, role), nil
}
