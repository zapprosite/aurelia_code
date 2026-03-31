package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/audio"
	"github.com/kocar/aurelia/internal/streaming"
	"github.com/kocar/aurelia/internal/streaming/actors"
	"github.com/kocar/aurelia/pkg/stt"
)

func runLiveCommand(args []string, out io.Writer) error {
	ctx := context.Background()
	log.Printf("[Jarvis] Iniciando Modo Live (Sovereign 2026.1 - SAP Architecture)")

	// 1. Setup Components
	app, err := bootstrapApp(nil)
	if err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}
	defer app.close()

	// 2. Build STT transcriber (local Whisper first, Groq fallback)
	transcriber, err := stt.NewTranscriber(
		app.cfg.STTProvider,
		app.cfg.GroqAPIKey,
		app.cfg.STTBaseURL,
		app.cfg.STTModel,
		app.cfg.STTLanguage,
	)
	if err != nil {
		slog.Warn("STT unavailable, voice input disabled", "err", err)
	}

	// 3. Initialize SAP Actors
	thinker := actors.NewAgentThinker(app.agentLoop)
	weaver := audio.NewSegmentedSynthesizer(app.cfg.TTSBaseURL, app.cfg.TTSVoice, app.cfg.TTSSpeed)
	speaker := audio.NewSimplePlayer()
	vad := actors.NewVADMonitor()

	pipeline := streaming.NewPipeline(thinker, weaver, speaker, vad)

	// Run VAD monitor in background
	go func() {
		if err := vad.Run(ctx); err != nil {
			slog.Error("VAD monitor failed", "err", err)
		}
	}()

	// Run wake word capture loop in background (if configured)
	wakeWordChan := make(chan string, 1)
	captureCmd := os.Getenv("AURELIA_VOICE_CAPTURE_CMD")
	if captureCmd == "" {
		captureCmd = "python3 scripts/voice-capture-openwakeword.py --output-dir /tmp/aurelia-voice"
	}
	go runWakeWordLoop(ctx, captureCmd, wakeWordChan)

	// 4. Input Sources
	keyboardChan := make(chan string)
	go func() {
		for {
			var userInput string
			if _, err := fmt.Scanln(&userInput); err != nil {
				if err == io.EOF {
					return
				}
				continue
			}
			keyboardChan <- userInput
		}
	}()

	log.Printf("[Jarvis] Todos os sistemas prontos. Pressione Ctrl+C para sair.")
	log.Printf("[Jarvis] AGUARDANDO VOZ (Her-Mode ativo)")

	// 5. Main Interaction Loop (Reactive)
	for {
		select {
		case <-ctx.Done():
			return nil

		case input := <-keyboardChan:
			if input == "exit" || input == "sair" {
				return nil
			}
			log.Printf("[Voce (Teclado)] %s", input)
			if err := pipeline.Process(ctx, input); err != nil {
				log.Printf("[Error] Pipeline: %v", err)
			}

		case audioFile := <-wakeWordChan:
			// Wake word detected — transcribe and process
			if transcriber == nil {
				log.Printf("[Jarvis] VOZ DETECTADA mas STT indisponivel")
				continue
			}
			log.Printf("[Jarvis] VOZ DETECTADA — transcrevendo %s", audioFile)
			text, err := transcriber.Transcribe(ctx, audioFile)
			if err != nil {
				log.Printf("[Error] STT: %v", err)
				continue
			}
			text = strings.TrimSpace(text)
			if text == "" {
				log.Printf("[Jarvis] Transcricao vazia, ignorando")
				continue
			}
			log.Printf("[Voce (Voz)] %s", text)
			if err := pipeline.Process(ctx, text); err != nil {
				log.Printf("[Error] Pipeline: %v", err)
			}
			// Clean up temp audio file
			_ = os.Remove(audioFile)

		case <-vad.Trigger:
			// Barge-in signal from socket-based VAD (complementary to wake word)
			log.Printf("[Jarvis] Barge-in / Activation Signal via VAD socket")
		}
	}
}

// runWakeWordLoop continuously runs the wake word detection script.
// When a wake word is detected, the audio file path is sent to the channel.
func runWakeWordLoop(ctx context.Context, captureCmd string, out chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		cmd := exec.CommandContext(ctx, "/bin/sh", "-c", captureCmd)
		output, err := cmd.Output()
		if err != nil {
			// Non-zero exit or context cancelled — retry after brief pause
			if ctx.Err() != nil {
				return
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

		raw := strings.TrimSpace(string(output))
		if raw == "" {
			continue
		}

		var result struct {
			Detected  bool   `json:"detected"`
			AudioFile string `json:"audio_file"`
		}
		if err := json.Unmarshal([]byte(raw), &result); err != nil {
			slog.Warn("failed to parse wake word output", "raw", raw, "err", err)
			continue
		}

		if result.Detected && result.AudioFile != "" {
			select {
			case out <- result.AudioFile:
			default:
				// Channel full, skip this detection
			}
		}
	}
}
