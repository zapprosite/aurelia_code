package agent

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

func (s *MasterTeamService) BuildStatusSnapshot(ctx context.Context, teamKey string) (TeamStatusSnapshot, error) {
	return s.buildSnapshot(ctx, teamKey, "")
}

func (s *MasterTeamService) BuildExecutionStatusSnapshot(ctx context.Context, teamKey, runID string) (TeamStatusSnapshot, error) {
	return s.buildSnapshot(ctx, teamKey, runID)
}

func (s *MasterTeamService) buildSnapshot(ctx context.Context, teamKey, runID string) (TeamStatusSnapshot, error) {
	s.mu.Lock()
	teamID := s.teamByKey[teamKey]
	s.mu.Unlock()
	if teamID == "" {
		return TeamStatusSnapshot{}, fmt.Errorf("team not found for key %s", teamKey)
	}

	tasks, err := s.manager.ListTasks(ctx, teamID)
	if err != nil {
		return TeamStatusSnapshot{}, err
	}

	filtered := filterTasksByRunID(tasks, runID)
	teamStatus, err := s.manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		return TeamStatusSnapshot{}, err
	}
	snapshot := TeamStatusSnapshot{TeamKey: teamKey, TeamID: teamID, TeamStatus: teamStatus}
	if len(filtered) == 0 {
		return snapshot, nil
	}

	for _, statuses := range groupTaskStatusesByRoot(filtered) {
		snapshot.TotalTasks++
		switch resolveLogicalStatus(statuses) {
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

func filterTasksByRunID(tasks []TeamTask, runID string) []TeamTask {
	if runID == "" {
		return tasks
	}
	filtered := make([]TeamTask, 0, len(tasks))
	for _, task := range tasks {
		if task.RunID == runID {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func groupTaskStatusesByRoot(tasks []TeamTask) map[string][]TaskStatus {
	taskByID := make(map[string]TeamTask, len(tasks))
	for _, task := range tasks {
		taskByID[task.ID] = task
	}

	groupStatuses := make(map[string][]TaskStatus)
	for _, task := range tasks {
		rootID := task.ID
		current := task
		seen := map[string]bool{current.ID: true}
		for current.ParentTaskID != nil && *current.ParentTaskID != "" {
			parent, ok := taskByID[*current.ParentTaskID]
			if !ok || seen[parent.ID] {
				break
			}
			rootID = parent.ID
			current = parent
			seen[current.ID] = true
		}
		groupStatuses[rootID] = append(groupStatuses[rootID], task.Status)
	}
	return groupStatuses
}

func resolveLogicalStatus(statuses []TaskStatus) TaskStatus {
	if len(statuses) == 0 {
		return TaskPending
	}

	switch {
	case slices.Contains(statuses, TaskRunning):
		return TaskRunning
	case slices.Contains(statuses, TaskBlocked):
		return TaskBlocked
	case slices.Contains(statuses, TaskPending):
		return TaskPending
	case slices.Contains(statuses, TaskCompleted):
		return TaskCompleted
	case slices.Contains(statuses, TaskFailed):
		return TaskFailed
	case slices.Contains(statuses, TaskCancelled):
		return TaskCancelled
	default:
		return statuses[0]
	}
}

func (s *MasterTeamService) formatMasterNotification(snapshot TeamStatusSnapshot, processedCount int, lines []string) string {
	// Cabeçalho Premium
	header := "━━━━━━ Aurelia Sovereign 2026 ━━━━━━"
	if snapshot.Pending == 0 && snapshot.Running == 0 && snapshot.Blocked == 0 && snapshot.TotalTasks > 0 {
		header = "━━━━━━ Missão Concluída ━━━━━━"
	}

	statusIcons := formatStatusIcons(snapshot)
	humanStatus := s.formatHumanStatus(snapshot)
	body := formatBodyLines(lines)

	if snapshot.Pending == 0 && snapshot.Running == 0 && snapshot.Blocked == 0 && snapshot.TotalTasks > 0 {
		return fmt.Sprintf(
			"%s\n\n%s %s\n\n%s\n\n🎯 **Resumo da Operação**\n%s\n\n*Recursos transientes liberados.*\n\n📊 [Acesse o Dashboard Completo](https://aurelia.zappro.site/)",
			header,
			statusIcons,
			classifyFinalSnapshot(snapshot),
			humanStatus,
			body,
		)
	}

	return fmt.Sprintf(
		"%s\n\n%s **Acompanhamento de Time**\n\n%s\n\n⚙️ **Ações Recentes (%d)**\n%s\n\n🔗 [Live Script Dashboard](https://aurelia.zappro.site/)",
		header,
		statusIcons,
		humanStatus,
		processedCount,
		body,
	)
}

func formatStatusIcons(snapshot TeamStatusSnapshot) string {
	if snapshot.Failed > 0 {
		return "⚠️"
	}
	if snapshot.Blocked > 0 {
		return "🛑"
	}
	if snapshot.Running > 5 {
		return "🔥" // Stress mode
	}
	if snapshot.Running > 0 {
		return "⚡"
	}
	return "🛰️"
}

func (s *MasterTeamService) formatHumanStatus(snapshot TeamStatusSnapshot) string {
	statusStr := fallbackTeamStatus(snapshot.TeamStatus)
	var parts []string
	if snapshot.Pending > 0 {
		parts = append(parts, fmt.Sprintf("📂 %d pendentes", snapshot.Pending))
	}
	if snapshot.Running > 0 {
		parts = append(parts, fmt.Sprintf("🔄 %d em execução", snapshot.Running))
	}
	if snapshot.Completed > 0 {
		parts = append(parts, fmt.Sprintf("✅ %d concluídas", snapshot.Completed))
	}
	if snapshot.Failed > 0 {
		parts = append(parts, fmt.Sprintf("❌ %d falhas", snapshot.Failed))
	}
	if snapshot.Blocked > 0 {
		parts = append(parts, fmt.Sprintf("🧱 %d bloqueios", snapshot.Blocked))
	}

	return fmt.Sprintf("Estado: **%s**\n%s", strings.ToUpper(statusStr), strings.Join(parts, "  ·  "))
}

func formatBodyLines(lines []string) string {
	if len(lines) == 0 {
		return "_Nenhuma mudança estrutural detectada._"
	}
	var res []string
	for _, line := range lines {
		// Limpeza de logs técnicos crus que podem vir no body
		clean := strings.TrimPrefix(line, "- ")
		if strings.Contains(clean, "{") && strings.Contains(clean, "}") {
			continue // Pula linhas com JSON bruto
		}
		res = append(res, "• "+clean)
	}
	if len(res) == 0 {
		return "_Metadados de execução refinados (silenciados no Telegram)._"
	}
	return strings.Join(res, "\n")
}

func (s *MasterTeamService) finalizeTeamRunIfIdle(ctx context.Context, teamKey, runID string, snapshot *TeamStatusSnapshot) {
	if snapshot == nil || snapshot.TeamID == "" {
		return
	}
	if snapshot.Pending != 0 || snapshot.Running != 0 || snapshot.Blocked != 0 || snapshot.TotalTasks == 0 {
		return
	}

	finalStatus := TeamStatusCompleted
	if snapshot.Failed > 0 || snapshot.Cancelled > 0 {
		finalStatus = TeamStatusAttentionRequired
	}
	if err := s.manager.SetTeamStatus(ctx, snapshot.TeamID, finalStatus); err == nil {
		snapshot.TeamStatus = finalStatus
	}

	s.cancelWorkerLoops(snapshot.TeamID)
	s.clearRunRuntimeState(ctx, snapshot.TeamID, runID)
}

func (s *MasterTeamService) clearRunRuntimeState(ctx context.Context, teamID, runID string) {
	tasks, err := s.manager.ListTasks(ctx, teamID)
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, task := range tasks {
		if runID != "" && task.RunID != runID {
			continue
		}
		delete(s.recoveryCount, task.ID)
	}
}

func classifyFinalSnapshot(snapshot TeamStatusSnapshot) string {
	switch {
	case snapshot.Failed == 0 && snapshot.Cancelled == 0:
		return "sucesso total"
	case snapshot.Completed > 0:
		return "conclusao parcial"
	default:
		return "bloqueio terminal"
	}
}

func fallbackTeamStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return "active"
	}
	return status
}

func classicStatusLine(snapshot TeamStatusSnapshot) string {
	return fmt.Sprintf(
		"status: pending=%d running=%d blocked=%d completed=%d failed=%d cancelled=%d total=%d",
		snapshot.Pending,
		snapshot.Running,
		snapshot.Blocked,
		snapshot.Completed,
		snapshot.Failed,
		snapshot.Cancelled,
		snapshot.TotalTasks,
	)
}
