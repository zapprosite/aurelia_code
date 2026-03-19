package voice

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const defaultSQLiteMirrorTable = "voice_events"

type SQLiteMirror struct {
	db    *sql.DB
	table string
}

func NewSQLiteMirror(dbPath string) *SQLiteMirror {
	dbPath = strings.TrimSpace(dbPath)
	if dbPath == "" {
		return nil
	}
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil
	}
	mirror := &SQLiteMirror{db: db, table: defaultSQLiteMirrorTable}
	if err := mirror.init(); err != nil {
		_ = db.Close()
		return nil
	}
	return mirror
}

func (m *SQLiteMirror) Close() error {
	if m == nil || m.db == nil {
		return nil
	}
	return m.db.Close()
}

func (m *SQLiteMirror) init() error {
	if m == nil || m.db == nil {
		return fmt.Errorf("sqlite mirror is not configured")
	}
	_, err := m.db.Exec(`
	CREATE TABLE IF NOT EXISTS voice_events (
		job_id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL DEFAULT 0,
		chat_id INTEGER NOT NULL DEFAULT 0,
		source TEXT NOT NULL DEFAULT '',
		transcript TEXT NOT NULL DEFAULT '',
		accepted INTEGER NOT NULL DEFAULT 0,
		requires_tts INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		mirrored_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_voice_events_created_at ON voice_events(created_at);
	CREATE INDEX IF NOT EXISTS idx_voice_events_user_chat ON voice_events(user_id, chat_id);
	`)
	if err != nil {
		return fmt.Errorf("init sqlite voice mirror: %w", err)
	}
	return nil
}

func (m *SQLiteMirror) MirrorTranscript(ctx context.Context, event TranscriptEvent) error {
	if m == nil || m.db == nil {
		return nil
	}
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO voice_events (
			job_id, user_id, chat_id, source, transcript, accepted, requires_tts, created_at, mirrored_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(job_id) DO UPDATE SET
			user_id = excluded.user_id,
			chat_id = excluded.chat_id,
			source = excluded.source,
			transcript = excluded.transcript,
			accepted = excluded.accepted,
			requires_tts = excluded.requires_tts,
			created_at = excluded.created_at,
			mirrored_at = excluded.mirrored_at
	`,
		event.JobID,
		event.UserID,
		event.ChatID,
		event.Source,
		event.Transcript,
		boolToInt(event.Accepted),
		boolToInt(event.RequiresTTS),
		event.CreatedAt.UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("sqlite voice mirror: %w", err)
	}
	return nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
