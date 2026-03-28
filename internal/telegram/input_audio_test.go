package telegram

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/memory"
	_ "modernc.org/sqlite"
)

func TestFormatAudioTranscriptArchiveContent(t *testing.T) {
	got := formatAudioTranscriptArchiveContent("groq", "/tmp/sample.ogg", "  ola mundo  ")

	if !strings.Contains(got, "provider=groq") {
		t.Fatalf("expected provider marker, got %q", got)
	}
	if !strings.Contains(got, "file=sample.ogg") {
		t.Fatalf("expected basename only, got %q", got)
	}
	if !strings.HasSuffix(got, "ola mundo") {
		t.Fatalf("expected trimmed transcript suffix, got %q", got)
	}
}

func TestPersistAudioTranscript_AddsArchiveEntry(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	mem := memory.NewMemoryManager(db, nil)
	t.Cleanup(func() {
		_ = mem.Close()
	})

	bc := &BotController{
		config: &config.AppConfig{STTProvider: "groq"},
		memory: mem,
	}

	bc.persistAudioTranscriptForSender(42, "/tmp/voice.ogg", "ola teste")

	entries, err := mem.ListArchiveEntries(context.Background(), "42:aurelia", 10)
	if err != nil {
		t.Fatalf("ListArchiveEntries() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 archive entry, got %d", len(entries))
	}
	if entries[0].MessageType != "audio_transcript" {
		t.Fatalf("unexpected message type %q", entries[0].MessageType)
	}
	if !strings.Contains(entries[0].Content, "provider=groq") {
		t.Fatalf("expected provider in content, got %q", entries[0].Content)
	}
	if !strings.Contains(entries[0].Content, "ola teste") {
		t.Fatalf("expected transcript text, got %q", entries[0].Content)
	}
}
