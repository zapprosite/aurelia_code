package voice

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestSQLiteMirror_MirrorTranscript(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "voice.db")
	mirror := NewSQLiteMirror(dbPath)
	if mirror == nil {
		t.Fatal("expected sqlite mirror")
	}
	defer func() { _ = mirror.Close() }()

	event := TranscriptEvent{
		JobID:       "job-1",
		UserID:      7,
		ChatID:      9,
		Source:      "mic",
		Transcript:  "jarvis revisar os logs",
		Accepted:    true,
		RequiresTTS: false,
		CreatedAt:   time.Now().UTC(),
	}
	if err := mirror.MirrorTranscript(context.Background(), event); err != nil {
		t.Fatalf("MirrorTranscript() error = %v", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer func() { _ = db.Close() }()

	var transcript string
	var accepted int
	if err := db.QueryRow(`SELECT transcript, accepted FROM voice_events WHERE job_id = ?`, event.JobID).Scan(&transcript, &accepted); err != nil {
		t.Fatalf("QueryRow().Scan() error = %v", err)
	}
	if transcript != event.Transcript || accepted != 1 {
		t.Fatalf("row = transcript=%q accepted=%d", transcript, accepted)
	}
}
