package streaming

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// Pipeline coordena o fluxo de streaming reativo.
type Pipeline struct {
	Thinker Thinker
	Weaver  Weaver
	Speaker Speaker
	VAD     Interrupter

	logger *slog.Logger
	cancel context.CancelFunc
	mu     sync.Mutex
}

func NewPipeline(t Thinker, w Weaver, s Speaker, v Interrupter) *Pipeline {
	p := &Pipeline{
		Thinker: t,
		Weaver:  w,
		Speaker: s,
		VAD:     v,
		logger:  slog.Default().With("component", "pipeline"),
	}

	// Configurar callback de interrupção
	if p.VAD != nil {
		p.VAD.OnInterrupt(func() {
			p.Interrupt()
		})
	}

	return p
}

// Interrupt cancela a rodada atual imediatamente.
func (p *Pipeline) Interrupt() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cancel != nil {
		p.logger.Warn("Pipeline interrupted by VAD")
		p.cancel()
		p.cancel = nil
	}
}

// Process inicia um novo ciclo de pensamento e voz.
func (p *Pipeline) Process(ctx context.Context, input string) error {
	p.mu.Lock()
	// Cancelar qualquer processo anterior se ainda estiver rodando.
	if p.cancel != nil {
		p.cancel()
	}
	
	// Criar novo contexto cancelável para esta rodada.
	runCtx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.mu.Unlock()

	p.logger.Info("Starting streaming pipeline processing", "input", input)

	// 1. Think (LLM)
	tokenStream, err := p.Thinker.Think(runCtx, input, nil)
	if err != nil {
		return fmt.Errorf("thinker failure: %w", err)
	}

	// 2. Weave (TTS)
	audioStream, errStream := p.Weaver.Weave(runCtx, tokenStream)

	// 3. Speak (Player)
	// O pipeline de áudio roda em paralelo com a geração de tokens.
	err = p.Speaker.Speak(runCtx, audioStream)
	
	// Verifica erros assíncronos do Weaver
	select {
	case err := <-errStream:
		if err != nil {
			p.logger.Error("Weaver error", "err", err)
			return err
		}
	default:
	}

	p.logger.Info("Interaction loop completed successfully")
	return err
}
