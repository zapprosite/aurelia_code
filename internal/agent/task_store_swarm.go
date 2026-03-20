package agent

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SwarmChannel struct {
	ID        string
	TeamID    string
	Name      string
	Kind      string
	CreatedAt time.Time
}

type SwarmThread struct {
	ID        string
	TeamID    string
	ChannelID string
	Title     string
	Status    string
	CreatedAt time.Time
}

type SwarmThreadMessage struct {
	ID         string
	TeamID     string
	ThreadID   string
	ChannelID  string
	Sender     string
	Kind       string
	Body       string
	CreatedAt  time.Time
	ConsumedAt *time.Time
}

type AssistanceTask struct {
	ID          string
	TeamID      string
	ThreadID    *string
	OwnerAgent  string
	HelperAgent *string
	Title       string
	Body        string
	Status      string
	LeaseUntil  *time.Time
	CreatedAt   time.Time
}

func (s *SQLiteTaskStore) createSwarmChannel(ctx context.Context, channel SwarmChannel) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO swarm_channels (id, team_id, name, kind)
		VALUES (?, ?, ?, ?)
	`, channel.ID, channel.TeamID, channel.Name, channel.Kind)
	if err != nil {
		return fmt.Errorf("insert swarm channel: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) createSwarmThread(ctx context.Context, thread SwarmThread) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO swarm_threads (id, team_id, channel_id, title, status)
		VALUES (?, ?, ?, ?, ?)
	`, thread.ID, thread.TeamID, thread.ChannelID, thread.Title, thread.Status)
	if err != nil {
		return fmt.Errorf("insert swarm thread: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) postSwarmThreadMessage(ctx context.Context, msg SwarmThreadMessage) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO swarm_thread_messages (id, team_id, thread_id, channel_id, sender_agent, kind, body, consumed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NULL)
	`, msg.ID, msg.TeamID, msg.ThreadID, msg.ChannelID, msg.Sender, msg.Kind, msg.Body)
	if err != nil {
		return fmt.Errorf("insert swarm thread message: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) enqueueAssistanceTask(ctx context.Context, task AssistanceTask) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO assistance_tasks (id, team_id, thread_id, owner_agent, helper_agent, title, body, status, lease_until)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL)
	`, task.ID, task.TeamID, task.ThreadID, task.OwnerAgent, task.HelperAgent, task.Title, task.Body, "pending")
	if err != nil {
		return fmt.Errorf("insert assistance task: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) claimAssistanceTask(ctx context.Context, teamID, helperAgent string) (*AssistanceTask, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin claim assistance tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now().UTC()
	rows, err := tx.QueryContext(ctx, `
		SELECT id, team_id, thread_id, owner_agent, helper_agent, title, body, status, lease_until, created_at
		FROM assistance_tasks
		WHERE team_id = ?
		  AND status = 'pending'
		ORDER BY created_at ASC, id ASC
		LIMIT 1
	`, teamID)
	if err != nil {
		return nil, fmt.Errorf("query assistance task: %w", err)
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit empty assistance claim tx: %w", err)
		}
		return nil, nil
	}

	task, err := scanAssistanceTask(rows)
	if err != nil {
		return nil, err
	}
	leaseUntil := now.Add(s.leaseDuration)
	res, err := tx.ExecContext(ctx, `
		UPDATE assistance_tasks
		SET helper_agent = ?, status = 'claimed', lease_until = ?
		WHERE id = ? AND status = 'pending'
	`, helperAgent, leaseUntil, task.ID)
	if err != nil {
		return nil, fmt.Errorf("claim assistance task: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("claim assistance task rows affected: %w", err)
	}
	if rowsAffected != 1 {
		return nil, nil
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit assistance claim tx: %w", err)
	}
	task.HelperAgent = &helperAgent
	task.Status = "claimed"
	task.LeaseUntil = &leaseUntil
	return task, nil
}

func (s *SQLiteTaskStore) listSwarmChannels(ctx context.Context, teamID string) ([]SwarmChannel, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, team_id, name, kind, created_at
		FROM swarm_channels
		WHERE team_id = ?
		ORDER BY created_at ASC, id ASC
	`, teamID)
	if err != nil {
		return nil, fmt.Errorf("list swarm channels: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var channels []SwarmChannel
	for rows.Next() {
		var item SwarmChannel
		if err := rows.Scan(&item.ID, &item.TeamID, &item.Name, &item.Kind, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan swarm channel: %w", err)
		}
		channels = append(channels, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("swarm channel rows: %w", err)
	}
	return channels, nil
}

func scanAssistanceTask(scanner interface{ Scan(dest ...any) error }) (*AssistanceTask, error) {
	var task AssistanceTask
	var threadID, helperAgent sql.NullString
	var leaseUntil sql.NullTime
	if err := scanner.Scan(
		&task.ID,
		&task.TeamID,
		&threadID,
		&task.OwnerAgent,
		&helperAgent,
		&task.Title,
		&task.Body,
		&task.Status,
		&leaseUntil,
		&task.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan assistance task: %w", err)
	}
	if threadID.Valid {
		task.ThreadID = &threadID.String
	}
	if helperAgent.Valid {
		task.HelperAgent = &helperAgent.String
	}
	if leaseUntil.Valid {
		ts := leaseUntil.Time
		task.LeaseUntil = &ts
	}
	return &task, nil
}
