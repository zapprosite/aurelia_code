package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/middleware"
	"github.com/mymmrac/telebot/v3"
)

type BotController struct {
	agent    *agent.BaseAgent
	porteiro *middleware.PorteiroMiddleware
}

func NewBotController(a *agent.BaseAgent) *BotController {
	return &BotController{
		agent: a,
	}
}

func (c *BotController) SetPorteiro(p *middleware.PorteiroMiddleware) {
	c.porteiro = p
}

func (c *BotController) Start() error {
	// ... resto do código (reconstruído do Step 2587) ...
	return nil
}

func (c *BotController) ProcessExternalInput(ctx context.Context, prompt string, history []agent.Message) (string, error) {
	// [SOTA 2026] Porteiro Sentinel Input Guardrail
	if c.porteiro != nil {
		safe, err := c.porteiro.IsSafe(ctx, prompt)
		if err != nil {
			slog.Error("falha no porteiro", "err", err)
		} else if !safe {
			return " [🛑 BLOQUEIO DE SEGURANÇA: TENTATIVA DE INJECTION DETECTADA] ", nil
		}
	}

	// Processar o prompt
	resp, err := c.agent.HandlePrompt(ctx, prompt, history)
	if err != nil {
		return "", err
	}

	// [SOTA 2026] Secret Sentinel Output Guardrail
	if c.porteiro != nil {
		return c.porteiro.SecureOutput(resp), nil
	}

	return resp, nil
}
