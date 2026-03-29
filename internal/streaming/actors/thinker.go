package actors

import (
	"context"
	"log/slog"

	"github.com/kocar/aurelia/internal/agent"
)

// AgentThinker adapta o Loop do agente para o padrão de atores SAP.
type AgentThinker struct {
	loop   *agent.Loop
	logger *slog.Logger
}

func NewAgentThinker(l *agent.Loop) *AgentThinker {
	return &AgentThinker{
		loop:   l,
		logger: slog.Default().With("actor", "thinker", "component", "sap"),
	}
}

func (at *AgentThinker) Name() string {
	return "thinker"
}

func (at *AgentThinker) Think(ctx context.Context, task string, history []agent.Message) (<-chan agent.StreamResponse, error) {
	at.logger.Info("Thinking started", "task", task)
	
	opts := agent.LoopOptions{
		SystemPrompt:   "Você é o Jarvis, o assistente oficial de Will. Seja ultra-eficiente, sutil e fiel aos Protocolos Soberanos.",
		Task:           task,
		InitialHistory: history,
	}

	return at.loop.RunWithOptionsStream(ctx, opts)
}
