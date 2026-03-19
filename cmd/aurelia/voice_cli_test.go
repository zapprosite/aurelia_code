package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
)

func TestRunVoiceEnqueue_QueuesJobInConfiguredSpool(t *testing.T) {
	root := t.TempDir()
	t.Setenv("AURELIA_HOME", root)

	resolver, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() error = %v", err)
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		t.Fatalf("runtime.Bootstrap() error = %v", err)
	}
	if _, err := config.Load(resolver); err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}

	configPath := resolver.AppConfig()
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	payload["voice_spool_path"] = filepath.Join(root, "custom-spool")
	data, err = json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent() error = %v", err)
	}
	if err := os.WriteFile(configPath, append(data, '\n'), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	audioPath := filepath.Join(root, "sample.wav")
	if err := os.WriteFile(audioPath, []byte("fake-audio"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	var out bytes.Buffer
	if err := runVoiceEnqueue([]string{audioPath, "--user-id", "12", "--chat-id", "34", "--requires-audio", "--source", "smoke"}, &out); err != nil {
		t.Fatalf("runVoiceEnqueue() error = %v", err)
	}

	if !strings.Contains(out.String(), "voice job queued: ") {
		t.Fatalf("stdout = %q", out.String())
	}

	entries, err := os.ReadDir(filepath.Join(root, "custom-spool", "inbox"))
	if err != nil {
		t.Fatalf("os.ReadDir() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("queued jobs = %d, want 1", len(entries))
	}
}

func TestRunVoiceCommand_RejectsUnknownSubcommand(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	err := runVoiceCommand([]string{"unknown"}, &out)
	if err == nil {
		t.Fatal("runVoiceCommand() error = nil, want error")
	}
}
