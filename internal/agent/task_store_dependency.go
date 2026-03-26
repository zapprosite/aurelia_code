package agent

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func (s *SQLiteTaskStore) listDependentsByStatusTx(
	ctx context.Context, tx *sql.Tx, teamID, dependencyTaskID string, statuses ...TaskStatus,
) ([]string, error) {
	if len(statuses) == 0 {
		return nil, nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(statuses)), ",")
	args := make([]any, 0, len(statuses)+2)
	args = append(args, teamID, dependencyTaskID)
	for _, status := range statuses {
		args = append(args, status)
	}
	query := fmt.Sprintf(`
		SELECT DISTINCT t.id FROM tasks t
		JOIN task_dependencies d ON d.task_id = t.id
		WHERE t.team_id = ? AND d.depends_on_task_id = ? AND t.status IN (%s)
	`, placeholders)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query dependents by status: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan dependent task id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *SQLiteTaskStore) updateTaskStatusTx(
	ctx context.Context, tx *sql.Tx, teamID, taskID string, from, to TaskStatus, errorMessage string,
) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE tasks SET status = ?, error_message = ? WHERE id = ? AND team_id = ? AND status = ?`,
		to, errorMessage, taskID, teamID, from,
	)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) getTaskStatusTx(ctx context.Context, tx *sql.Tx, taskID string) (TaskStatus, error) {
	var status TaskStatus
	err := tx.QueryRowContext(ctx, `SELECT status FROM tasks WHERE id = ?`, taskID).Scan(&status)
	if err == sql.ErrNoRows {
		return "", sql.ErrNoRows
	}
	if err != nil {
		return "", fmt.Errorf("read task status: %w", err)
	}
	return status, nil
}

func (s *SQLiteTaskStore) taskExistsInTeamTx(ctx context.Context, tx *sql.Tx, teamID, taskID string) (bool, error) {
	var count int
	err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM tasks WHERE team_id = ? AND id = ?`, teamID, taskID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check task existence: %w", err)
	}
	return count > 0, nil
}

func (s *SQLiteTaskStore) getParentTaskIDTx(ctx context.Context, tx *sql.Tx, taskID string) (string, bool, error) {
	var parentID sql.NullString
	if err := tx.QueryRowContext(ctx, `SELECT parent_task_id FROM tasks WHERE id = ?`, taskID).Scan(&parentID); err != nil {
		return "", false, fmt.Errorf("read parent task id: %w", err)
	}
	if !parentID.Valid || parentID.String == "" {
		return "", false, nil
	}
	return parentID.String, true, nil
}

// ── Graph ─────────────────────────────────────────────────────────────────────

func (s *SQLiteTaskStore) resolveRootTaskIDTx(ctx context.Context, tx *sql.Tx, taskID string) (string, error) {
	currentID := taskID
	seen := map[string]bool{}
	for {
		if seen[currentID] {
			return "", fmt.Errorf("parent task cycle detected at %s", currentID)
		}
		seen[currentID] = true
		var parentID sql.NullString
		err := tx.QueryRowContext(ctx, `SELECT parent_task_id FROM tasks WHERE id = ?`, currentID).Scan(&parentID)
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("parent task %s does not exist", currentID)
		}
		if err != nil {
			return "", fmt.Errorf("resolve root task id: %w", err)
		}
		if !parentID.Valid || strings.TrimSpace(parentID.String) == "" {
			return currentID, nil
		}
		currentID = parentID.String
	}
}

func (s *SQLiteTaskStore) hasDependencyPathTx(ctx context.Context, tx *sql.Tx, fromTaskID, targetTaskID string) (bool, error) {
	if fromTaskID == targetTaskID {
		return true, nil
	}
	queue := []string{fromTaskID}
	visited := map[string]bool{}
	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]
		if visited[currentID] {
			continue
		}
		visited[currentID] = true
		rows, err := tx.QueryContext(ctx, `SELECT depends_on_task_id FROM task_dependencies WHERE task_id = ?`, currentID)
		if err != nil {
			return false, fmt.Errorf("query dependency path: %w", err)
		}
		var nextIDs []string
		for rows.Next() {
			var nextID string
			if err := rows.Scan(&nextID); err != nil {
				_ = rows.Close()
				return false, fmt.Errorf("scan dependency path: %w", err)
			}
			if nextID == targetTaskID {
				_ = rows.Close()
				return true, nil
			}
			nextIDs = append(nextIDs, nextID)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return false, fmt.Errorf("dependency path rows: %w", err)
		}
		_ = rows.Close()
		queue = append(queue, nextIDs...)
	}
	return false, nil
}

func (s *SQLiteTaskStore) listChildTasksTx(ctx context.Context, tx *sql.Tx, parentTaskID string) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id FROM tasks WHERE parent_task_id = ? ORDER BY created_at ASC, id ASC`, parentTaskID)
	if err != nil {
		return nil, fmt.Errorf("query child tasks: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var childIDs []string
	for rows.Next() {
		var childID string
		if err := rows.Scan(&childID); err != nil {
			return nil, fmt.Errorf("scan child task: %w", err)
		}
		childIDs = append(childIDs, childID)
	}
	return childIDs, rows.Err()
}

// ── Resolution ────────────────────────────────────────────────────────────────

func (s *SQLiteTaskStore) resolveDependencyStatusTx(ctx context.Context, tx *sql.Tx, taskID string) (TaskStatus, error) {
	status, err := s.getTaskStatusTx(ctx, tx, taskID)
	if err != nil {
		return "", err
	}
	if status == TaskCompleted {
		return TaskCompleted, nil
	}
	children, err := s.listChildTasksTx(ctx, tx, taskID)
	if err != nil {
		return "", err
	}
	return s.resolveRecoveryStatusTx(ctx, tx, status, children)
}

func (s *SQLiteTaskStore) resolveRecoveryStatusTx(ctx context.Context, tx *sql.Tx, parentStatus TaskStatus, childIDs []string) (TaskStatus, error) {
	hasCompletedRecovery := false
	hasActiveRecovery := false
	for _, childID := range childIDs {
		childStatus, err := s.resolveDependencyStatusTx(ctx, tx, childID)
		if err != nil {
			return "", err
		}
		switch childStatus {
		case TaskCompleted:
			hasCompletedRecovery = true
		case TaskPending, TaskRunning, TaskBlocked:
			hasActiveRecovery = true
		}
	}
	switch {
	case hasCompletedRecovery:
		return TaskCompleted, nil
	case hasActiveRecovery:
		return TaskBlocked, nil
	default:
		return parentStatus, nil
	}
}

func (s *SQLiteTaskStore) areAllDependenciesCompletedTx(ctx context.Context, tx *sql.Tx, taskID string) (bool, error) {
	rows, err := tx.QueryContext(ctx, `SELECT depends_on_task_id FROM task_dependencies WHERE task_id = ?`, taskID)
	if err != nil {
		return false, fmt.Errorf("query dependency completion: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var depID string
		if err := rows.Scan(&depID); err != nil {
			return false, fmt.Errorf("scan dependency completion: %w", err)
		}
		status, err := s.resolveDependencyStatusTx(ctx, tx, depID)
		if err != nil {
			return false, err
		}
		if status != TaskCompleted {
			return false, nil
		}
	}
	return true, rows.Err()
}

// ── Rules / Validation ────────────────────────────────────────────────────────

func (s *SQLiteTaskStore) evaluateDependencyStateTx(ctx context.Context, tx *sql.Tx, dependencyIDs []string) (bool, string, error) {
	for _, depID := range dependencyIDs {
		status, err := s.resolveDependencyStatusTx(ctx, tx, depID)
		if err != nil {
			if err == sql.ErrNoRows {
				return true, fmt.Sprintf("blocked by missing dependency %s", depID), nil
			}
			return false, "", err
		}
		switch status {
		case TaskCompleted:
			continue
		case TaskFailed, TaskCancelled:
			return true, fmt.Sprintf("blocked by failed dependency %s", depID), nil
		default:
			return true, fmt.Sprintf("blocked by dependency %s", depID), nil
		}
	}
	return false, "", nil
}

func (s *SQLiteTaskStore) validateTaskDependenciesTx(ctx context.Context, tx *sql.Tx, task TeamTask, dependencyIDs []string) error {
	protectedIDs := map[string]struct{}{}
	if strings.TrimSpace(task.ID) != "" {
		protectedIDs[task.ID] = struct{}{}
	}
	if task.ParentTaskID != nil && strings.TrimSpace(*task.ParentTaskID) != "" {
		rootID, err := s.resolveRootTaskIDTx(ctx, tx, *task.ParentTaskID)
		if err != nil {
			return err
		}
		protectedIDs[rootID] = struct{}{}
	}
	for _, depID := range dependencyIDs {
		depID = strings.TrimSpace(depID)
		if depID == "" {
			return fmt.Errorf("dependency id cannot be empty")
		}
		if _, ok := protectedIDs[depID]; ok {
			return fmt.Errorf("dependency cycle detected for task %s", task.ID)
		}
		exists, err := s.taskExistsInTeamTx(ctx, tx, task.TeamID, depID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("dependency %s does not exist in team %s", depID, task.TeamID)
		}
		for protectedID := range protectedIDs {
			reaches, err := s.hasDependencyPathTx(ctx, tx, depID, protectedID)
			if err != nil {
				return err
			}
			if reaches {
				return fmt.Errorf("dependency cycle detected between %s and %s", depID, protectedID)
			}
		}
	}
	return nil
}

// ── Cascade ───────────────────────────────────────────────────────────────────

func (s *SQLiteTaskStore) cancelDependentsTx(ctx context.Context, tx *sql.Tx, teamID, taskID, reason string) error {
	dependentIDs, err := s.listDependentsByStatusTx(ctx, tx, teamID, taskID, TaskPending, TaskBlocked)
	if err != nil {
		return err
	}
	for _, dependentID := range dependentIDs {
		if err := s.cancelDependentCascadeTx(ctx, tx, teamID, dependentID, reason); err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLiteTaskStore) cancelDependentCascadeTx(ctx context.Context, tx *sql.Tx, teamID, dependentID, reason string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE tasks SET status = ?, error_message = ?, finished_at = ? WHERE id = ? AND team_id = ? AND status IN (?, ?)`,
		TaskCancelled, reason, time.Now().UTC(), dependentID, teamID, TaskPending, TaskBlocked,
	)
	if err != nil {
		return fmt.Errorf("cancel dependent task: %w", err)
	}
	if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    teamID,
		TaskID:    &dependentID,
		AgentName: MasterAgentName,
		EventType: "task_cancelled",
		Payload:   reason,
	}); err != nil {
		return err
	}
	return s.cancelDependentsTx(ctx, tx, teamID, dependentID, fmt.Sprintf("cancelled because dependency %s was cancelled", dependentID))
}

// ── Flow (unblock) ────────────────────────────────────────────────────────────

func (s *SQLiteTaskStore) unblockDependentsTx(ctx context.Context, tx *sql.Tx, teamID, completedTaskID string) error {
	dependentIDs, err := s.listDependentsByStatusTx(ctx, tx, teamID, completedTaskID, TaskBlocked)
	if err != nil {
		return err
	}
	for _, dependentID := range dependentIDs {
		ready, err := s.areAllDependenciesCompletedTx(ctx, tx, dependentID)
		if err != nil {
			return err
		}
		if !ready {
			continue
		}
		if err := s.updateTaskStatusTx(ctx, tx, teamID, dependentID, TaskBlocked, TaskPending, ""); err != nil {
			return err
		}
		if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
			TeamID:    teamID,
			TaskID:    &dependentID,
			AgentName: MasterAgentName,
			EventType: "task_unblocked",
			Payload:   fmt.Sprintf("dependency %s completed", completedTaskID),
		}); err != nil {
			return err
		}
	}
	return nil
}

// ── Recovery ──────────────────────────────────────────────────────────────────

func (s *SQLiteTaskStore) reopenDependentsForRecoveryTx(ctx context.Context, tx *sql.Tx, teamID, dependencyTaskID, recoveryTaskID string) error {
	parentStatus, err := s.getTaskStatusTx(ctx, tx, dependencyTaskID)
	if err != nil {
		return err
	}
	if parentStatus != TaskFailed && parentStatus != TaskCancelled {
		return nil
	}
	dependentIDs, err := s.listDependentsByStatusTx(ctx, tx, teamID, dependencyTaskID, TaskCancelled)
	if err != nil {
		return err
	}
	for _, dependentID := range dependentIDs {
		reason := fmt.Sprintf("blocked while dependency %s is in recovery via %s", dependencyTaskID, recoveryTaskID)
		_, err := tx.ExecContext(ctx,
			`UPDATE tasks SET status = ?, error_message = ?, finished_at = NULL WHERE id = ? AND team_id = ? AND status = ?`,
			TaskBlocked, reason, dependentID, teamID, TaskCancelled,
		)
		if err != nil {
			return fmt.Errorf("reopen dependent task: %w", err)
		}
		for _, eventType := range []string{"task_reopened", "task_blocked"} {
			if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
				TeamID:    teamID,
				TaskID:    &dependentID,
				AgentName: MasterAgentName,
				EventType: eventType,
				Payload:   reason,
			}); err != nil {
				return err
			}
		}
		if err := s.reopenDependentsForRecoveryTx(ctx, tx, teamID, dependentID, recoveryTaskID); err != nil {
			return err
		}
	}
	return nil
}
