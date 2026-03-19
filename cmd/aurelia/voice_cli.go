package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/internal/voice"
)

func runVoiceCommand(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("voice command is required")
	}

	switch args[0] {
	case "enqueue":
		return runVoiceEnqueue(args[1:], stdout)
	default:
		return fmt.Errorf("unknown voice command %q", args[0])
	}
}

func runVoiceEnqueue(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: aurelia voice enqueue <audio-file> [--user-id N] [--chat-id N] [--requires-audio] [--source name]")
	}

	audioPath := args[0]
	var (
		userID        int64
		chatID        int64
		requiresAudio bool
		source        = "cli"
	)
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--user-id":
			i++
			if i >= len(args) {
				return fmt.Errorf("--user-id requires a value")
			}
			value, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid --user-id: %w", err)
			}
			userID = value
		case "--chat-id":
			i++
			if i >= len(args) {
				return fmt.Errorf("--chat-id requires a value")
			}
			value, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid --chat-id: %w", err)
			}
			chatID = value
		case "--source":
			i++
			if i >= len(args) {
				return fmt.Errorf("--source requires a value")
			}
			source = args[i]
		case "--requires-audio":
			requiresAudio = true
		default:
			return fmt.Errorf("unknown flag %q", args[i])
		}
	}

	resolver, err := runtime.New()
	if err != nil {
		return err
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		return err
	}
	cfg, err := config.Load(resolver)
	if err != nil {
		return err
	}
	spool, err := voice.NewSpool(cfg.VoiceSpoolPath)
	if err != nil {
		return err
	}
	if userID == 0 {
		userID = cfg.VoiceReplyUserID
	}
	if chatID == 0 {
		chatID = cfg.VoiceReplyChatID
	}
	job, err := spool.EnqueueAudioFile(voice.Job{
		Source:        source,
		UserID:        userID,
		ChatID:        chatID,
		RequiresAudio: requiresAudio,
	}, audioPath)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(stdout, "voice job queued: %s\n", job.ID)
	return err
}
