package voice

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSpoolEnqueueClaimAndComplete(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "voice-spool")
	spool, err := NewSpool(root)
	if err != nil {
		t.Fatalf("NewSpool() error = %v", err)
	}

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("audio"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	job, err := spool.EnqueueAudioFile(Job{Source: "test", UserID: 42, ChatID: 42}, audioPath)
	if err != nil {
		t.Fatalf("EnqueueAudioFile() error = %v", err)
	}
	if job.ID == "" {
		t.Fatal("expected generated job id")
	}

	claimed, err := spool.ClaimOldest(context.Background())
	if err != nil {
		t.Fatalf("ClaimOldest() error = %v", err)
	}
	if claimed == nil {
		t.Fatal("expected claimed job")
	}
	if claimed.Job.ID != job.ID {
		t.Fatalf("claimed job id = %q", claimed.Job.ID)
	}

	if err := spool.Complete(claimed, "ola"); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "done", job.ID, "result.json")); err != nil {
		t.Fatalf("expected result.json in done dir: %v", err)
	}
}

func TestSpoolFailMovesJobToFailed(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "voice-spool")
	spool, err := NewSpool(root)
	if err != nil {
		t.Fatalf("NewSpool() error = %v", err)
	}

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("audio"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	job, err := spool.EnqueueAudioFile(Job{Source: "test"}, audioPath)
	if err != nil {
		t.Fatalf("EnqueueAudioFile() error = %v", err)
	}

	claimed, err := spool.ClaimOldest(context.Background())
	if err != nil {
		t.Fatalf("ClaimOldest() error = %v", err)
	}
	if err := spool.Fail(claimed, os.ErrPermission); err != nil {
		t.Fatalf("Fail() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "failed", job.ID, "result.json")); err != nil {
		t.Fatalf("expected result.json in failed dir: %v", err)
	}
}
