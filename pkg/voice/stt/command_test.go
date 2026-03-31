package stt

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandTranscriber_IsAvailable(t *testing.T) {
	t.Parallel()

	if NewCommandTranscriber("").IsAvailable() {
		t.Fatal("expected empty command to be unavailable")
	}
	if !NewCommandTranscriber("printf ola").IsAvailable() {
		t.Fatal("expected configured command to be available")
	}
}

func TestCommandTranscriber_UsesAudioEnvAndStdout(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	audioPath := filepath.Join(dir, "sample.wav")
	if err := os.WriteFile(audioPath, []byte("fake"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	transcriber := NewCommandTranscriber(`basename "$AURELIA_AUDIO_FILE"`)
	got, err := transcriber.Transcribe(context.Background(), audioPath)
	if err != nil {
		t.Fatalf("Transcribe() error = %v", err)
	}
	if got != "sample.wav" {
		t.Fatalf("transcript = %q", got)
	}
}

func TestCommandTranscriber_ReturnsCommandFailure(t *testing.T) {
	t.Parallel()

	transcriber := NewCommandTranscriber(`echo boom >&2; exit 4`)
	_, err := transcriber.Transcribe(context.Background(), "/tmp/audio.wav")
	if err == nil {
		t.Fatal("expected Transcribe() to fail")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error = %v", err)
	}
}
