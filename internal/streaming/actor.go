package streaming

import (
	"context"
	"log/slog"

	"github.com/kocar/aurelia/internal/agent"
)

// Actor é a base para todos os componentes do pipeline SAP.
type Actor interface {
	Run(ctx context.Context) error
	Name() string
}

// Thinker é o ator que gera o fluxo de pensamento (LLM).
type Thinker interface {
	Think(ctx context.Context, task string, history []agent.Message) (<-chan agent.StreamResponse, error)
}

// Weaver é o ator que transforma pensamento em áudio (TTS).
type Weaver interface {
	Weave(ctx context.Context, tokenStream <-chan agent.StreamResponse) (<-chan []byte, <-chan error)
}

// Speaker é o ator que dá voz física ao sistema (Player).
type Speaker interface {
	Speak(ctx context.Context, audioStream <-chan []byte) error
}

// Interrupter é o contrato para fatias de Barge-in (SOTA 2026.2).
type Interrupter interface {
	Actor
	OnInterrupt(callback func())
}

// BaseActor fornece utilitários comuns para implementações de atores.
type BaseActor struct {
	actorName string
	Logger    *slog.Logger
}

func NewBaseActor(name string) BaseActor {
	return BaseActor{
		actorName: name,
		Logger:    ActorLogger(name),
	}
}

func (b *BaseActor) Name() string {
	return b.actorName
}

// ActorLogger fornece telemetria estruturada padrão SOTA 2026.
func ActorLogger(actor string) *slog.Logger {
	return slog.Default().With("actor", actor, "component", "sap")
}
