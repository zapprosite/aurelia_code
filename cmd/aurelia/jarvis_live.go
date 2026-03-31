package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"

	"github.com/kocar/aurelia/internal/audio"
	"github.com/kocar/aurelia/internal/streaming"
	"github.com/kocar/aurelia/internal/streaming/actors"
)

func runLiveCommand(args []string, out io.Writer) error {
	ctx := context.Background()
	log.Printf("[Jarvis] 🚀 Iniciando Modo Live (Sovereign 2026.1 - SAP Architecture)")

	// 1. Setup Components
	app, err := bootstrapApp(nil)
	if err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}
	defer app.close()

	// 2. Initialize SAP Actors
	thinker := actors.NewAgentThinker(app.agentLoop)
	weaver := audio.NewSegmentedSynthesizer(app.cfg.TTSBaseURL, "pt-br_isabela", 1.0)
	speaker := audio.NewSimplePlayer()
	vad := actors.NewVADMonitor()

	pipeline := streaming.NewPipeline(thinker, weaver, speaker, vad)

	// Rodar o VAD monitor em background
	go func() {
		if err := vad.Run(ctx); err != nil {
			slog.Error("VAD monitor failed", "err", err)
		}
	}()

	log.Printf("[Jarvis] ✅ Todos os sistemas prontos. Pressione Ctrl+C para sair.")

	// 4. Input Sources
	keyboardChan := make(chan string)
	go func() {
		for {
			var userInput string
			if _, err := fmt.Scanln(&userInput); err != nil {
				if err == io.EOF { return }
				continue
			}
			keyboardChan <- userInput
		}
	}()

	log.Printf("[Jarvis] ✅ Todos os sistemas prontos. Pressione Ctrl+C para sair.")
	log.Printf("[Jarvis] 🎙️ AGUARDANDO VOZ (Her-Mode ativo no background...)")

	// 4. Main Interaction Loop (Reactive)
	for {
		select {
		case <-ctx.Done():
			return nil
		
		case input := <-keyboardChan:
			if input == "exit" || input == "sair" {
				return nil
			}
			log.Printf("[Você (Teclado)] %s", input)
			if err := pipeline.Process(ctx, input); err != nil {
				log.Printf("[Error] Pipeline: %v", err)
			}

		case <-vad.Trigger:
			// No Her-Mode real, aqui dispararíamos o STT.
			// Por enquanto, apenas logamos a ativação reativa.
			log.Printf("[Jarvis] 🔔 VOZ DETECTADA (Barge-in / Activation Signal)")
			// Futuro: STT -> task -> pipeline.Process(ctx, task)
		}
	}

	return nil
}
