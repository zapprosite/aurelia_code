package stt

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CommandTranscriber executes a local shell command that must print the
// transcript to stdout. The input audio path is exposed through the
// AURELIA_AUDIO_FILE environment variable.
type CommandTranscriber struct {
	command string
}

func NewCommandTranscriber(command string) *CommandTranscriber {
	return &CommandTranscriber{command: strings.TrimSpace(command)}
}

func (t *CommandTranscriber) Transcribe(ctx context.Context, audioFilePath string) (string, error) {
	if !t.IsAvailable() {
		return "", fmt.Errorf("command transcriber not configured")
	}

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", t.command)
	cmd.Env = append(os.Environ(), "AURELIA_AUDIO_FILE="+audioFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("fallback command failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	transcript := strings.TrimSpace(string(output))
	if transcript == "" {
		return "", fmt.Errorf("fallback command returned empty transcript")
	}
	return transcript, nil
}

func (t *CommandTranscriber) IsAvailable() bool {
	return t != nil && strings.TrimSpace(t.command) != ""
}
