package main

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/internal/voice"
	"github.com/kocar/aurelia/pkg/stt"
)

func runVoiceCommand(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("voice command is required")
	}

	switch args[0] {
	case "enqueue":
		return runVoiceEnqueue(args[1:], stdout)
	case "capture-once":
		return runVoiceCaptureOnce(args[1:], stdout)
	case "process-once":
		return runVoiceProcessOnce(args[1:], stdout)
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

func runVoiceCaptureOnce(args []string, stdout io.Writer) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: aurelia voice capture-once")
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
	if !cfg.VoiceEnabled {
		return fmt.Errorf("voice is not enabled")
	}
	if !cfg.VoiceCaptureEnabled {
		return fmt.Errorf("voice capture is not enabled")
	}
	if cfg.VoiceCaptureCommand == "" {
		return fmt.Errorf("voice_capture_command is not configured")
	}
	spool, err := voice.NewSpool(cfg.VoiceSpoolPath)
	if err != nil {
		return err
	}
	worker := voice.NewCaptureWorker(
		spool,
		voice.NewCommandCaptureSource(cfg.VoiceCaptureCommand, map[string]string{
			"AURELIA_VOICE_WAKE_PHRASE": cfg.VoiceWakePhrase,
			"AURELIA_VOICE_DROP_PATH":   cfg.VoiceDropPath,
			"AURELIA_VOICE_USER_ID":     strconv.FormatInt(cfg.VoiceReplyUserID, 10),
			"AURELIA_VOICE_CHAT_ID":     strconv.FormatInt(cfg.VoiceReplyChatID, 10),
		}),
		voice.CaptureConfig{
			PollInterval:       time.Duration(cfg.VoiceCapturePollMS) * time.Millisecond,
			HeartbeatPath:      cfg.VoiceCaptureHeartbeat,
			HeartbeatFreshness: time.Duration(cfg.VoiceCaptureFreshSec) * time.Second,
			DefaultUserID:      cfg.VoiceReplyUserID,
			DefaultChatID:      cfg.VoiceReplyChatID,
			DefaultSource:      "capture",
		},
	)
	if err := worker.CaptureOnce(context.Background()); err != nil {
		return err
	}
	_, err = fmt.Fprintln(stdout, "voice capture tick completed")
	return err
}

func runVoiceProcessOnce(args []string, stdout io.Writer) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: aurelia voice process-once")
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
	if !cfg.VoiceEnabled {
		return fmt.Errorf("voice is not enabled")
	}
	spool, err := voice.NewSpool(cfg.VoiceSpoolPath)
	if err != nil {
		return err
	}
	transcriber, err := stt.NewTranscriber(cfg.STTProvider, cfg.GroqAPIKey, cfg.STTBaseURL, cfg.STTModel, cfg.STTLanguage)
	if err != nil {
		return err
	}
	var fallback stt.Transcriber
	if cfg.STTFallbackCommand != "" {
		fallback = stt.NewCommandTranscriber(cfg.STTFallbackCommand)
	}
	processor := voice.NewProcessor(spool, transcriber, fallback, nil, voice.Config{
		PollInterval:       time.Duration(cfg.VoicePollIntervalMS) * time.Millisecond,
		HeartbeatPath:      cfg.VoiceHeartbeatPath,
		HeartbeatFreshness: time.Duration(cfg.VoiceHeartbeatFreshSec) * time.Second,
		WakePhrase:         cfg.VoiceWakePhrase,
		DefaultUserID:      cfg.VoiceReplyUserID,
		DefaultChatID:      cfg.VoiceReplyChatID,
		SoftCapDaily:       cfg.GroqSoftCapDaily,
		HardCapDaily:       cfg.GroqHardCapDaily,
		PrimaryLabel:       cfg.STTProvider,
		Mirror:             voice.NewSQLiteMirror(cfg.DBPath),
	})
	if err := processor.ProcessOnce(context.Background()); err != nil {
		return err
	}
	_, err = fmt.Fprintln(stdout, "voice processor tick completed")
	return err
}
