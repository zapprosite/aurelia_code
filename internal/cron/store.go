package cron

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// ── Store ─────────────────────────────────────────────────────────────────────

type SQLiteCronStore struct {
	db *sql.DB
}

func NewSQLiteCronStore(dbPath string) (*SQLiteCronStore, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open cron sqlite store: %w", err)
	}
	store := &SQLiteCronStore{db: db}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SQLiteCronStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

// ── Schema ────────────────────────────────────────────────────────────────────

func (s *SQLiteCronStore) initialize() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS cron_jobs (
		id TEXT PRIMARY KEY,
		owner_user_id TEXT NOT NULL,
		target_chat_id INTEGER NOT NULL,
		schedule_type TEXT NOT NULL,
		cron_expr TEXT NOT NULL DEFAULT '',
		run_at DATETIME,
		prompt TEXT NOT NULL,
		active INTEGER NOT NULL DEFAULT 1,
		last_run_at DATETIME,
		next_run_at DATETIME,
		last_status TEXT NOT NULL DEFAULT 'idle',
		last_error TEXT NOT NULL DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS cron_executions (
		id TEXT PRIMARY KEY,
		job_id TEXT NOT NULL,
		started_at DATETIME NOT NULL,
		finished_at DATETIME,
		status TEXT NOT NULL,
		output_summary TEXT NOT NULL DEFAULT '',
		error_message TEXT NOT NULL DEFAULT ''
	);`)
	if err != nil {
		return fmt.Errorf("initialize cron schema: %w", err)
	}
	return nil
}

// ── Scan helpers ──────────────────────────────────────────────────────────────

func scanCronJob(scanner interface{ Scan(dest ...any) error }) (*CronJob, error) {
	var job CronJob
	var runAt, lastRunAt, nextRunAt sql.NullTime
	var active int
	err := scanner.Scan(
		&job.ID, &job.OwnerUserID, &job.TargetChatID, &job.ScheduleType, &job.CronExpr,
		&runAt, &job.Prompt, &active, &lastRunAt, &nextRunAt,
		&job.LastStatus, &job.LastError, &job.CreatedAt, &job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	job.Active = active == 1
	if runAt.Valid { ts := runAt.Time; job.RunAt = &ts }
	if lastRunAt.Valid { ts := lastRunAt.Time; job.LastRunAt = &ts }
	if nextRunAt.Valid { ts := nextRunAt.Time; job.NextRunAt = &ts }
	return &job, nil
}

func scanCronJobs(rows *sql.Rows) ([]CronJob, error) {
	var jobs []CronJob
	for rows.Next() {
		job, err := scanCronJob(rows)
		if err != nil {
			return nil, fmt.Errorf("scan cron job row: %w", err)
		}
		jobs = append(jobs, *job)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cron job rows: %w", err)
	}
	return jobs, nil
}

func boolToInt(v bool) int {
	if v { return 1 }
	return 0
}

const jobSelectCols = `id, owner_user_id, target_chat_id, schedule_type, cron_expr, run_at, prompt, active,
	last_run_at, next_run_at, last_status, last_error, created_at, updated_at`

// ── Jobs CRUD ─────────────────────────────────────────────────────────────────

func (s *SQLiteCronStore) CreateJob(ctx context.Context, job CronJob) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO cron_jobs (id, owner_user_id, target_chat_id, schedule_type, cron_expr, run_at, prompt, active,
			last_run_at, next_run_at, last_status, last_error)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		job.ID, job.OwnerUserID, job.TargetChatID, job.ScheduleType, job.CronExpr, job.RunAt, job.Prompt,
		boolToInt(job.Active), job.LastRunAt, job.NextRunAt, job.LastStatus, job.LastError,
	)
	if err != nil {
		return fmt.Errorf("insert cron job: %w", err)
	}
	return nil
}

func (s *SQLiteCronStore) UpdateJob(ctx context.Context, job CronJob) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE cron_jobs
		SET owner_user_id=?, target_chat_id=?, schedule_type=?, cron_expr=?, run_at=?, prompt=?, active=?,
			last_run_at=?, next_run_at=?, last_status=?, last_error=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?`,
		job.OwnerUserID, job.TargetChatID, job.ScheduleType, job.CronExpr, job.RunAt, job.Prompt,
		boolToInt(job.Active), job.LastRunAt, job.NextRunAt, job.LastStatus, job.LastError, job.ID,
	)
	if err != nil {
		return fmt.Errorf("update cron job: %w", err)
	}
	return nil
}

func (s *SQLiteCronStore) DeleteJob(ctx context.Context, jobID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM cron_jobs WHERE id = ?`, jobID)
	if err != nil {
		return fmt.Errorf("delete cron job: %w", err)
	}
	return nil
}

func (s *SQLiteCronStore) GetJob(ctx context.Context, jobID string) (*CronJob, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+jobSelectCols+` FROM cron_jobs WHERE id = ?`, jobID)
	job, err := scanCronJob(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get cron job: %w", err)
	}
	return job, nil
}

func (s *SQLiteCronStore) ListJobsByChat(ctx context.Context, chatID int64) ([]CronJob, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT `+jobSelectCols+` FROM cron_jobs WHERE target_chat_id = ? ORDER BY created_at ASC, id ASC`, chatID)
	if err != nil {
		return nil, fmt.Errorf("list cron jobs by chat: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return scanCronJobs(rows)
}

func (s *SQLiteCronStore) ListDueJobs(ctx context.Context, now time.Time, limit int) ([]CronJob, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT `+jobSelectCols+` FROM cron_jobs WHERE active=1 AND next_run_at IS NOT NULL AND next_run_at<=?
		 ORDER BY next_run_at ASC, id ASC LIMIT ?`, now.UTC(), limit)
	if err != nil {
		return nil, fmt.Errorf("list due cron jobs: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return scanCronJobs(rows)
}

// ── Executions ────────────────────────────────────────────────────────────────

func (s *SQLiteCronStore) RecordExecution(ctx context.Context, exec CronExecution) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO cron_executions (id, job_id, started_at, finished_at, status, output_summary, error_message)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		exec.ID, exec.JobID, exec.StartedAt, exec.FinishedAt, exec.Status, exec.OutputSummary, exec.ErrorMessage,
	)
	if err != nil {
		return fmt.Errorf("insert cron execution: %w", err)
	}
	return nil
}

func (s *SQLiteCronStore) ListExecutionsByJob(ctx context.Context, jobID string) ([]CronExecution, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, job_id, started_at, finished_at, status, output_summary, error_message
		FROM cron_executions WHERE job_id = ? ORDER BY started_at ASC, id ASC`, jobID)
	if err != nil {
		return nil, fmt.Errorf("list cron executions: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var executions []CronExecution
	for rows.Next() {
		var exec CronExecution
		var finishedAt sql.NullTime
		if err := rows.Scan(&exec.ID, &exec.JobID, &exec.StartedAt, &finishedAt, &exec.Status, &exec.OutputSummary, &exec.ErrorMessage); err != nil {
			return nil, fmt.Errorf("scan cron execution row: %w", err)
		}
		if finishedAt.Valid { ts := finishedAt.Time; exec.FinishedAt = &ts }
		executions = append(executions, exec)
	}
	return executions, rows.Err()
}

// ── Scan (due jobs helper used by scheduler) ──────────────────────────────────

func (s *SQLiteCronStore) scanDueJobs(now time.Time, limit int) ([]CronJob, error) {
	return s.ListDueJobs(context.Background(), now, limit)
}
