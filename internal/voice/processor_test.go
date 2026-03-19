package voice

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/health"
	"github.com/kocar/aurelia/pkg/stt"
)

type fakeTranscriber struct {
	available bool
	text      string
	err       error
	calls     int
}

func (f *fakeTranscriber) Transcribe(ctx context.Context, audioFilePath string) (string, error) {
	f.calls++
	return f.text, f.err
}

func (f *fakeTranscriber) IsAvailable() bool {
	return f.available
}

type fakeDispatcher struct {
	userID        int64
	chatID        int64
	text          string
	requiresAudio bool
	calls         int
	err           error
}

func (f *fakeDispatcher) DispatchVoice(ctx context.Context, userID, chatID int64, text string, requiresAudio bool) error {
	f.calls++
	f.userID = userID
	f.chatID = chatID
	f.text = text
	f.requiresAudio = requiresAudio
	return f.err
}

type fakeMirror struct {
	events []TranscriptEvent
}

func (f *fakeMirror) MirrorTranscript(ctx context.Context, event TranscriptEvent) error {
	f.events = append(f.events, event)
	return nil
}

func TestProcessorProcessOnce_DispatchesAcceptedTranscript(t *testing.T) {
	t.Parallel()

	spool, audioPath := newVoiceTestSpoolWithAudio(t)
	_, err := spool.EnqueueAudioFile(Job{Source: "drop", UserID: 7, ChatID: 9}, audioPath)
	if err != nil {
		t.Fatalf("EnqueueAudioFile() error = %v", err)
	}

	primary := &fakeTranscriber{available: true, text: "jarvis verifique o homelab"}
	dispatcher := &fakeDispatcher{}
	mirror := &fakeMirror{}
	processor := NewProcessor(spool, primary, nil, dispatcher, Config{
		HeartbeatPath: filepath.Join(t.TempDir(), "heartbeat.json"),
		Mirror:        mirror,
	})

	if err := processor.ProcessOnce(context.Background()); err != nil {
		t.Fatalf("ProcessOnce() error = %v", err)
	}
	if dispatcher.calls != 1 {
		t.Fatalf("dispatcher calls = %d", dispatcher.calls)
	}
	if dispatcher.text != "verifique o homelab" {
		t.Fatalf("dispatcher text = %q", dispatcher.text)
	}
	if len(mirror.events) != 1 || !mirror.events[0].Accepted {
		t.Fatalf("mirror events = %+v", mirror.events)
	}
}

func TestProcessorProcessOnce_DropsTranscriptWithoutWakePhrase(t *testing.T) {
	t.Parallel()

	spool, audioPath := newVoiceTestSpoolWithAudio(t)
	_, err := spool.EnqueueAudioFile(Job{Source: "drop"}, audioPath)
	if err != nil {
		t.Fatalf("EnqueueAudioFile() error = %v", err)
	}

	primary := &fakeTranscriber{available: true, text: "isso nao deve executar"}
	dispatcher := &fakeDispatcher{}
	processor := NewProcessor(spool, primary, nil, dispatcher, Config{
		HeartbeatPath: filepath.Join(t.TempDir(), "heartbeat.json"),
	})

	if err := processor.ProcessOnce(context.Background()); err != nil {
		t.Fatalf("ProcessOnce() error = %v", err)
	}
	if dispatcher.calls != 0 {
		t.Fatalf("dispatcher calls = %d", dispatcher.calls)
	}
}

func TestProcessorProcessOnce_FallsBackOn429(t *testing.T) {
	t.Parallel()

	spool, audioPath := newVoiceTestSpoolWithAudio(t)
	_, err := spool.EnqueueAudioFile(Job{Source: "drop", UserID: 1, ChatID: 1}, audioPath)
	if err != nil {
		t.Fatalf("EnqueueAudioFile() error = %v", err)
	}

	primary := &fakeTranscriber{available: true, err: errors.New("API error (status 429): too many requests")}
	fallback := &fakeTranscriber{available: true, text: "jarvis usar fallback"}
	dispatcher := &fakeDispatcher{}
	processor := NewProcessor(spool, primary, fallback, dispatcher, Config{
		HeartbeatPath: filepath.Join(t.TempDir(), "heartbeat.json"),
	})

	if err := processor.ProcessOnce(context.Background()); err != nil {
		t.Fatalf("ProcessOnce() error = %v", err)
	}
	if primary.calls != 1 || fallback.calls != 1 {
		t.Fatalf("calls primary=%d fallback=%d", primary.calls, fallback.calls)
	}
	if dispatcher.calls != 1 {
		t.Fatalf("dispatcher calls = %d", dispatcher.calls)
	}
}

func TestProcessorHealthCheck_ReportsStaleHeartbeatAndSoftCap(t *testing.T) {
	t.Parallel()

	processor := NewProcessor(nil, nil, nil, nil, Config{
		HeartbeatPath:      filepath.Join(t.TempDir(), "heartbeat.json"),
		HeartbeatFreshness: 5 * time.Second,
	})
	processor.updateStatus(func(status *HeartbeatStatus) {
		status.LastBeatAt = time.Now().Add(-10 * time.Second)
	})

	check := processor.HealthCheck()
	if check.Status != "error" {
		t.Fatalf("status = %q", check.Status)
	}

	processor.updateStatus(func(status *HeartbeatStatus) {
		status.LastBeatAt = time.Now()
		status.SoftCapReached = true
		status.LastError = ""
	})
	check = processor.HealthCheck()
	if check.Status != "warning" {
		t.Fatalf("status = %q", check.Status)
	}
}

func TestSanitizeTranscript_RedactsSecrets(t *testing.T) {
	t.Parallel()

	input := "minha chave e sk-test123456 e hf_ABCDEF123456 e AIzaSyA35JE5D8p6j_hw"
	got := sanitizeTranscript(input)
	if strings.Contains(got, "sk-test123456") || strings.Contains(got, "hf_ABCDEF123456") || strings.Contains(got, "AIzaSyA35JE5D8p6j_hw") {
		t.Fatalf("sanitizeTranscript() leaked secret: %q", got)
	}
}

func TestStripWakePhrase(t *testing.T) {
	t.Parallel()

	got, ok := stripWakePhrase("Ei Jarvis, abra o painel", "jarvis")
	if !ok || got != "abra o painel" {
		t.Fatalf("stripWakePhrase() = %q %v", got, ok)
	}
}

func newVoiceTestSpoolWithAudio(t *testing.T) (*Spool, string) {
	t.Helper()

	spool, err := NewSpool(filepath.Join(t.TempDir(), "voice-spool"))
	if err != nil {
		t.Fatalf("NewSpool() error = %v", err)
	}
	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("audio"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return spool, audioPath
}

var _ stt.Transcriber = (*fakeTranscriber)(nil)
var _ Dispatcher = (*fakeDispatcher)(nil)
var _ Mirror = (*fakeMirror)(nil)
var _ = health.CheckResult{}
