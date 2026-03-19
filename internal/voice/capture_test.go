package voice

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type fakeCaptureSource struct {
	event *CaptureEvent
	err   error
	calls int
}

func (f *fakeCaptureSource) Capture(ctx context.Context) (*CaptureEvent, error) {
	f.calls++
	return f.event, f.err
}

func (f *fakeCaptureSource) Description() string { return "fake" }

func TestCaptureWorkerCaptureOnce_QueuesDetectedClip(t *testing.T) {
	t.Parallel()

	spool, audioPath := newVoiceTestSpoolWithAudio(t)
	source := &fakeCaptureSource{event: &CaptureEvent{
		Detected:      true,
		AudioFile:     audioPath,
		UserID:        12,
		ChatID:        34,
		RequiresAudio: true,
		Source:        "mic",
	}}
	worker := NewCaptureWorker(spool, source, CaptureConfig{
		HeartbeatPath: filepath.Join(t.TempDir(), "capture-heartbeat.json"),
	})

	if err := worker.CaptureOnce(context.Background()); err != nil {
		t.Fatalf("CaptureOnce() error = %v", err)
	}

	claimed, err := spool.ClaimOldest(context.Background())
	if err != nil {
		t.Fatalf("ClaimOldest() error = %v", err)
	}
	if claimed == nil {
		t.Fatal("expected claimed job")
	}
	if claimed.Job.Source != "mic" {
		t.Fatalf("source = %q", claimed.Job.Source)
	}
	if claimed.Job.UserID != 12 || claimed.Job.ChatID != 34 {
		t.Fatalf("job = %+v", claimed.Job)
	}
	if !claimed.Job.RequiresAudio {
		t.Fatal("expected requires audio flag")
	}
}

func TestCaptureWorkerCaptureOnce_IdleWhenNoDetection(t *testing.T) {
	t.Parallel()

	spool, _ := newVoiceTestSpoolWithAudio(t)
	source := &fakeCaptureSource{}
	worker := NewCaptureWorker(spool, source, CaptureConfig{
		HeartbeatPath: filepath.Join(t.TempDir(), "capture-heartbeat.json"),
	})

	if err := worker.CaptureOnce(context.Background()); err != nil {
		t.Fatalf("CaptureOnce() error = %v", err)
	}
	if source.calls != 1 {
		t.Fatalf("calls = %d", source.calls)
	}
	status := worker.StatusSnapshot()
	if status.Status != "idle" {
		t.Fatalf("status = %q", status.Status)
	}
}

func TestCaptureWorkerHealthCheck_ReportsStaleHeartbeat(t *testing.T) {
	t.Parallel()

	worker := NewCaptureWorker(nil, nil, CaptureConfig{
		HeartbeatPath:      filepath.Join(t.TempDir(), "capture-heartbeat.json"),
		HeartbeatFreshness: 5 * time.Second,
	})
	worker.updateStatus(func(status *CaptureStatus) {
		status.LastBeatAt = time.Now().Add(-10 * time.Second)
	})

	check := worker.HealthCheck()
	if check.Status != "error" {
		t.Fatalf("status = %q", check.Status)
	}
}

func TestCommandCaptureSourceCapture_ParsesJSON(t *testing.T) {
	t.Parallel()

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("audio"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	scriptPath := filepath.Join(t.TempDir(), "capture.sh")
	script := "#!/bin/sh\nprintf '%s' '{\"detected\":true,\"audio_file\":\"" + audioPath + "\",\"source\":\"mic\"}'\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o700); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	source := NewCommandCaptureSource(scriptPath, nil)
	event, err := source.Capture(context.Background())
	if err != nil {
		t.Fatalf("Capture() error = %v", err)
	}
	if event == nil || event.AudioFile != audioPath {
		t.Fatalf("event = %+v", event)
	}
}

func TestCommandCaptureSourceCapture_IgnoresStderrNoiseOnSuccess(t *testing.T) {
	t.Parallel()

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("audio"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	scriptPath := filepath.Join(t.TempDir(), "capture.sh")
	script := "#!/bin/sh\n" +
		"printf '%s\\n' 'runtime warning' >&2\n" +
		"printf '%s' '{\"detected\":true,\"audio_file\":\"" + audioPath + "\",\"source\":\"mic\"}'\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o700); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	source := NewCommandCaptureSource(scriptPath, nil)
	event, err := source.Capture(context.Background())
	if err != nil {
		t.Fatalf("Capture() error = %v", err)
	}
	if event == nil || event.AudioFile != audioPath {
		t.Fatalf("event = %+v", event)
	}
}
