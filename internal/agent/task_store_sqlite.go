package agent

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *SQLiteTaskStore) initialize() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			key TEXT UNIQUE,
			user_id TEXT,
			lead_agent TEXT,
			status TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS teammates (
			team_id TEXT,
			agent_name TEXT,
			role_description TEXT,
			last_heartbeat DATETIME,
			PRIMARY KEY (team_id, agent_name)
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			team_id TEXT,
			run_id TEXT,
			parent_task_id TEXT,
			title TEXT,
			prompt TEXT,
			working_dir TEXT,
			allowed_tools TEXT,
			assigned_agent TEXT,
			status TEXT,
			result_summary TEXT,
			error_message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			started_at DATETIME,
			finished_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS task_dependencies (
			task_id TEXT,
			depends_on_task_id TEXT,
			PRIMARY KEY (task_id, depends_on_task_id)
		)`,
		`CREATE TABLE IF NOT EXISTS task_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			team_id TEXT,
			task_id TEXT,
			agent_name TEXT,
			event_type TEXT,
			payload TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS mail (
			id TEXT PRIMARY KEY,
			team_id TEXT,
			from_agent TEXT,
			to_agent TEXT,
			task_id TEXT,
			kind TEXT,
			body TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			consumed_at DATETIME
		)`,
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return fmt.Errorf("initialize schema: %w", err)
		}
	}
	return nil
}

// implementation of TeamManager interface

func (s *SQLiteTaskStore) CreateTeam(ctx context.Context, teamKey, userID, leadAgent string) (string, error) {
	id := uuid.New().String()
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO teams (id, key, user_id, lead_agent, status) VALUES (?, ?, ?, ?, ?)`,
		id, teamKey, userID, leadAgent, TeamStatusActive,
	)
	return id, err
}

func (s *SQLiteTaskStore) GetTeamIDByKey(ctx context.Context, teamKey string) (string, error) {
	var id string
	err := s.db.QueryRowContext(ctx, `SELECT id FROM teams WHERE key = ?`, teamKey).Scan(&id)
	return id, err
}

func (s *SQLiteTaskStore) RegisterTeammate(ctx context.Context, teamID, agentName, roleDescription string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO teammates (team_id, agent_name, role_description, last_heartbeat) VALUES (?, ?, ?, ?)
		 ON CONFLICT(team_id, agent_name) DO UPDATE SET role_description = excluded.role_description, last_heartbeat = excluded.last_heartbeat`,
		teamID, agentName, roleDescription, time.Now().UTC(),
	)
	return err
}

func (s *SQLiteTaskStore) CreateTask(ctx context.Context, task TeamTask, dependsOn []string) error {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	toolsStr := strings.Join(task.AllowedTools, ",")

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO tasks (id, team_id, run_id, parent_task_id, title, prompt, working_dir, allowed_tools, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.TeamID, task.RunID, task.ParentTaskID, task.Title, task.Prompt, task.Workdir, toolsStr, TaskPending,
	)
	if err != nil {
		return err
	}

	for _, depID := range dependsOn {
		_, err = tx.ExecContext(ctx, `INSERT INTO task_dependencies (task_id, depends_on_task_id) VALUES (?, ?)`, task.ID, depID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteTaskStore) GetTask(ctx context.Context, teamID, taskID string) (*TeamTask, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, team_id, run_id, parent_task_id, title, prompt, working_dir, allowed_tools, assigned_agent, status, result_summary, error_message, created_at, started_at, finished_at FROM tasks WHERE id = ? AND team_id = ?`, taskID, teamID)
	return scanTask(row)
}

func (s *SQLiteTaskStore) ListTasks(ctx context.Context, teamID string) ([]TeamTask, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, team_id, run_id, parent_task_id, title, prompt, working_dir, allowed_tools, assigned_agent, status, result_summary, error_message, created_at, started_at, finished_at FROM tasks WHERE team_id = ? ORDER BY created_at ASC`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TeamTask
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, nil
}

func (s *SQLiteTaskStore) ClaimNextTask(ctx context.Context, teamID, agentName string) (*TeamTask, error) {
	// S-2026: Simplificado. Pega a primeira task pending que não tenha dependências não resolvidas.
	query := `
		SELECT t.id FROM tasks t
		LEFT JOIN task_dependencies d ON t.id = d.task_id
		WHERE t.team_id = ? AND t.status = ?
		AND (d.depends_on_task_id IS NULL OR EXISTS (SELECT 1 FROM tasks t2 WHERE t2.id = d.depends_on_task_id AND t2.status = ?))
		LIMIT 1
	`
	var taskID string
	err := s.db.QueryRowContext(ctx, query, teamID, TaskPending, TaskCompleted).Scan(&taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	now := time.Now()
	_, err = s.db.ExecContext(ctx, `UPDATE tasks SET status = ?, assigned_agent = ?, started_at = ? WHERE id = ?`, TaskRunning, agentName, now, taskID)
	if err != nil {
		return nil, err
	}

	return s.GetTask(ctx, teamID, taskID)
}

func (s *SQLiteTaskStore) HeartbeatWorker(ctx context.Context, teamID, agentName string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE teammates SET last_heartbeat = ? WHERE team_id = ? AND agent_name = ?`, time.Now().UTC(), teamID, agentName)
	return err
}

func (s *SQLiteTaskStore) CompleteTask(ctx context.Context, teamID, taskID, agentName, result string) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx, `UPDATE tasks SET status = ?, result_summary = ?, finished_at = ? WHERE id = ? AND team_id = ?`, TaskCompleted, result, now, taskID, teamID)
	return err
}

func (s *SQLiteTaskStore) FailTask(ctx context.Context, teamID, taskID, agentName, reason string) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx, `UPDATE tasks SET status = ?, error_message = ?, finished_at = ? WHERE id = ? AND team_id = ?`, TaskFailed, reason, now, taskID, teamID)
	return err
}

func (s *SQLiteTaskStore) GetTeamStatus(ctx context.Context, teamID string) (string, error) {
	var status string
	err := s.db.QueryRowContext(ctx, `SELECT status FROM teams WHERE id = ?`, teamID).Scan(&status)
	return status, err
}

func (s *SQLiteTaskStore) SetTeamStatus(ctx context.Context, teamID, status string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE teams SET status = ? WHERE id = ?`, status, teamID)
	return err
}

func (s *SQLiteTaskStore) CancelActiveTasks(ctx context.Context, teamID, reason string) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx, `UPDATE tasks SET status = ?, error_message = ?, finished_at = ? WHERE team_id = ? AND status IN (?, ?)`, TaskCancelled, reason, now, teamID, TaskPending, TaskRunning)
	return err
}

func (s *SQLiteTaskStore) PostMessage(ctx context.Context, msg MailMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO mail (id, team_id, from_agent, to_agent, task_id, kind, body) VALUES (?, ?, ?, ?, ?, ?, ?)`, msg.ID, msg.TeamID, msg.FromAgent, msg.ToAgent, msg.TaskID, msg.Kind, msg.Body)
	return err
}

func (s *SQLiteTaskStore) PullMessages(ctx context.Context, teamID, agentName string, limit int) ([]MailMessage, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, team_id, from_agent, to_agent, task_id, kind, body, created_at, consumed_at FROM mail WHERE team_id = ? AND to_agent = ? AND consumed_at IS NULL LIMIT ?`, teamID, agentName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MailMessage
	now := time.Now()
	for rows.Next() {
		var m MailMessage
		err := rows.Scan(&m.ID, &m.TeamID, &m.FromAgent, &m.ToAgent, &m.TaskID, &m.Kind, &m.Body, &m.CreatedAt, &m.ConsumedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
		_, _ = s.db.ExecContext(ctx, `UPDATE mail SET consumed_at = ? WHERE id = ?`, now, m.ID)
	}
	return results, nil
}

func (s *SQLiteTaskStore) ListEvents(ctx context.Context, teamID string, limit int) ([]TaskEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, team_id, task_id, agent_name, event_type, payload, created_at FROM task_events WHERE team_id = ? ORDER BY id DESC LIMIT ?`, teamID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TaskEvent
	for rows.Next() {
		var e TaskEvent
		err := rows.Scan(&e.ID, &e.TeamID, &e.TaskID, &e.AgentName, &e.EventType, &e.Payload, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, nil
}

func scanTask(row interface {
	Scan(dest ...any) error
}) (*TeamTask, error) {
	var t TeamTask
	var allowedToolsStr string
	err := row.Scan(
		&t.ID, &t.TeamID, &t.RunID, &t.ParentTaskID, &t.Title, &t.Prompt, &t.Workdir, &allowedToolsStr, &t.AssignedAgent, &t.Status, &t.ResultSummary, &t.ErrorMessage, &t.CreatedAt, &t.StartedAt, &t.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	t.AllowedTools = strings.Split(allowedToolsStr, ",")
	return &t, nil
}

// Métodos transacionais (...Tx) para suporte a operações complexas

func (s *SQLiteTaskStore) insertTaskEventTx(ctx context.Context, tx *sql.Tx, e TaskEvent) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO task_events (team_id, task_id, agent_name, event_type, payload) VALUES (?, ?, ?, ?, ?)`,
		e.TeamID, e.TaskID, e.AgentName, e.EventType, e.Payload,
	)
	return err
}

func (s *SQLiteTaskStore) requeueExpiredRunningTasksTx(ctx context.Context, tx *sql.Tx, teamID string) error {
	// Reenfileira tasks que não tiveram heartbeat recentemente (ex: 5 minutos)
	threshold := time.Now().Add(-5 * time.Minute).UTC()
	_, err := tx.ExecContext(ctx,
		`UPDATE tasks SET status = ?, assigned_agent = NULL, started_at = NULL 
		 WHERE team_id = ? AND status = ? AND started_at < ?`,
		TaskPending, teamID, TaskRunning, threshold,
	)
	return err
}

func (s *SQLiteTaskStore) renewWorkerLeaseTx(ctx context.Context, tx *sql.Tx, teamID, agentName string) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO teammates (team_id, agent_name, last_heartbeat) VALUES (?, ?, ?)
		 ON CONFLICT(team_id, agent_name) DO UPDATE SET last_heartbeat = excluded.last_heartbeat`,
		teamID, agentName, time.Now().UTC(),
	)
	return err
}

func (s *SQLiteTaskStore) validateTaskDependenciesTx(ctx context.Context, tx *sql.Tx, task TeamTask, dependsOn []string) error {
	for _, depID := range dependsOn {
		var exists bool
		err := tx.QueryRowContext(ctx, `SELECT 1 FROM tasks WHERE id = ? AND team_id = ?`, depID, task.TeamID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("dependency %s not found: %w", depID, err)
		}
	}
	return nil
}

func (s *SQLiteTaskStore) evaluateDependencyStateTx(ctx context.Context, tx *sql.Tx, dependsOn []string) (bool, string, error) {
	for _, depID := range dependsOn {
		var status string
		err := tx.QueryRowContext(ctx, `SELECT status FROM tasks WHERE id = ?`, depID).Scan(&status)
		if err != nil {
			return false, "", err
		}
		if status != string(TaskCompleted) {
			return true, fmt.Sprintf("Waiting for dependency %s (status: %s)", depID, status), nil
		}
	}
	return false, "", nil
}

func (s *SQLiteTaskStore) reopenDependentsForRecoveryTx(ctx context.Context, tx *sql.Tx, teamID, parentTaskID, newTaskID string) error {
	// No fluxo de recuperação, se uma task falhou e estamos tentando uma nova, 
	// podemos querer resetar o status de quem dependia da falha original.
	_, err := tx.ExecContext(ctx,
		`UPDATE tasks SET status = ?, error_message = NULL 
		 WHERE team_id = ? AND status = ? AND id IN (
			 SELECT task_id FROM task_dependencies WHERE depends_on_task_id = ?
		 )`,
		TaskPending, teamID, TaskBlocked, parentTaskID,
	)
	return err
}

func (s *SQLiteTaskStore) reopenDependentsTx(ctx context.Context, tx *sql.Tx, taskID string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE tasks SET status = ?, error_message = NULL 
		 WHERE status = ? AND id IN (
			 SELECT task_id FROM task_dependencies WHERE depends_on_task_id = ?
		 )`,
		TaskPending, TaskBlocked, taskID,
	)
	return err
}
