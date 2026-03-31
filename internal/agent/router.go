package agent

import (
	"context"
	"fmt"
	"log/slog"
)

// TieredRouter implementa LLMProvider com lógica de Fallback (SOTA 2026).
type TieredRouter struct {
	tiers []LLMProvider
}

// NewTieredRouter cria um roteador com a ordem de preferência fornecida.
func NewTieredRouter(providers ...LLMProvider) *TieredRouter {
	return &TieredRouter{
		tiers: providers,
	}
}

// GenerateContent tenta cada tier em ordem até obter uma resposta válida.
func (r *TieredRouter) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	var lastErr error
	for i, provider := range r.tiers {
		slog.Debug("Router: Tentando execução", "tier", i+1)
		resp, err := provider.GenerateContent(ctx, systemPrompt, history, tools)
		if err == nil {
			if resp.Metadata == nil {
				resp.Metadata = make(map[string]string)
			}
			resp.Metadata["aurelia-tier"] = fmt.Sprintf("%d", i+1)
			return resp, nil
		}
		slog.Warn("Router: Falha no tier", "tier", i+1, "error", err)
		lastErr = err
	}
	return nil, fmt.Errorf("todos os tiers de LLM falharam: %w", lastErr)
}

// GenerateStream redireciona para o primeiro tier disponível que suporte streaming.
// Nota: Em cenários de erro durante o stream, o fallback é mais complexo e 
// geralmente requer reinicialização da resposta pelo agente superior.
func (r *TieredRouter) GenerateStream(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (<-chan StreamResponse, error) {
	// Por simplicidade industrial, tentamos o Tier 1 primeiro.
	// Se falhar antes de iniciar o stream, tentamos o próximo.
	for i, provider := range r.tiers {
		ch, err := provider.GenerateStream(ctx, systemPrompt, history, tools)
		if err == nil {
			return ch, nil
		}
		slog.Warn("Router: Falha ao iniciar stream no tier", "tier", i+1, "error", err)
	}
	return nil, fmt.Errorf("falha ao iniciar streaming em todos os provedores")
}
