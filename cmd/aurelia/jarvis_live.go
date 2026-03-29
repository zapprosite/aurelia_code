package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

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

	// 4. Main Interaction Loop
	for {
		fmt.Fprintf(out, "\n[Você] (Modo Texto - Jarvis ouvindo...) ")
		
		var userInput string
		if _, err := fmt.Scanln(&userInput); err != nil {
			if err == io.EOF { break }
			continue
		}

		if userInput == "exit" || userInput == "sair" {
			break
		}

		// 5. Execute Pipeline (Think -> Weave -> Speak)
		// O pipeline gerencia o streaming completo de forma reativa.
		err := pipeline.Process(ctx, userInput)
		
		if err != nil {
			log.Printf("[Error] Pipeline: %v", err)
		}

		time.Sleep(200 * time.Millisecond)
	}

	return nil
}
