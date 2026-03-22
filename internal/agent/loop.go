package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"log/slog"

	"github.com/kocar/aurelia/internal/dashboard"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/plan"
)

type Loop struct {
	llm           LLMProvider
	registry      *ToolRegistry
	maxIterations int
}

func NewLoop(llm LLMProvider, registry *ToolRegistry, maxIterations int) *Loop {
	if maxIterations <= 0 {
		maxIterations = 10
	}
	return &Loop{
		llm:           llm,
		registry:      registry,
		maxIterations: maxIterations,
	}
}

func (l *Loop) Registry() *ToolRegistry {
	return l.registry
}

// Run é a interface clássica do loop, mantida para compatibilidade
func (l *Loop) Run(ctx context.Context, systemPrompt string, history []Message, allowedTools []string) ([]Message, string, error) {
	opts := LoopOptions{
		SystemPrompt:  systemPrompt,
		InitialHistory: history,
		MaxIterations: l.maxIterations,
	}
	// TODO: Na versão premium, filtrar registry baseado em allowedTools aqui
	return l.RunWithOptions(ctx, opts)
}

// WithMemoryAssembler wrapper para compatibilidade com app.go
func (l *Loop) WithMemoryAssembler(any) *Loop {
	return l
}

// WithToolCatalog wrapper para compatibilidade com app.go
func (l *Loop) WithToolCatalog(any, int) *Loop {
	return l
}

// WithSemanticRouter wrapper para compatibilidade com app.go
func (l *Loop) WithSemanticRouter(any) *Loop {
	return l
}

type LoopOptions struct {
	Task             string
	SystemPrompt     string
	InitialHistory   []Message
	MaxIterations    int
	InterruptHandler func() bool
}

func (l *Loop) RunWithOptions(ctx context.Context, opts LoopOptions) ([]Message, string, error) {
	logger := observability.Logger("agent.loop")
	currentHistory := opts.InitialHistory

	if opts.MaxIterations <= 0 {
		opts.MaxIterations = l.maxIterations
	}

	for i := 0; i < opts.MaxIterations; i++ {
		// Stop if context cancelled
		select {
		case <-ctx.Done():
			return currentHistory, "", ctx.Err()
		default:
		}

		if opts.InterruptHandler != nil && opts.InterruptHandler() {
			return currentHistory, "Interrompido pelo usuário.", nil
		}

		// Adicionar orientação sobre a fase atual no System Prompt dinâmico
		currentPhase, _ := PrevPhaseFromContext(ctx)
		dynamicSystemPrompt := augmentSystemPromptWithPhase(opts.SystemPrompt, currentPhase)

		// O provedor de LLM PRECISA das ferramentas registradas
		resp, err := l.llm.GenerateContent(ctx, dynamicSystemPrompt, currentHistory, l.registry.GetDefinitions())
		if err != nil {
			return currentHistory, "", fmt.Errorf("generate content error: %w", err)
		}

		// Handle completion
		if len(resp.ToolCalls) == 0 {
			if resp.Content != "" || resp.ReasoningContent != "" {
				currentHistory = append(currentHistory, Message{
					Role:             "assistant",
					Content:          resp.Content,
					ReasoningContent: resp.ReasoningContent,
				})
			}
			return currentHistory, resp.Content, nil
		}

		// AI wants to call tools
		currentHistory = append(currentHistory, Message{
			Role:             "assistant",
			Content:          resp.Content,
			ReasoningContent: resp.ReasoningContent,
			ToolCalls:        resp.ToolCalls,
		})

		for _, call := range resp.ToolCalls {
			logger.Info("executing tool", slog.String("tool_name", call.Name), slog.Any("arg_keys", observability.MapKeys(call.Arguments)))
			
			dashboard.Publish(dashboard.Event{
				Type:      "agent_tool",
				Agent:     "Aurelia",
				Action:    "Executando tool: " + call.Name,
				Payload:   call.Arguments,
				Timestamp: time.Now().Format("15:04:05"),
			})

			resultStr, toolErr := l.registry.Execute(ctx, call.Name, call.Arguments)
			if toolErr != nil {
				errorPayload, _ := json.Marshal(map[string]string{
					"error": toolErr.Error(),
				})
				resultStr = string(errorPayload)
			}

			// Interceptador global do PREV (Hard-Gate)
			if call.Name == "set_phase" {
				if newPhase, ok := call.Arguments["phase"].(string); ok {
					if newPhase == "EXECUTION" {
						if plan.GlobalPlanStore.HasPending() {
							resultStr = "BLOQUEIO DE SEGURANÇA: Existem planos propostos aguardando aprovação humana no Cockpit. Você NÃO pode iniciar a execução sem o OK do usuário."
							dashboard.Publish(dashboard.Event{
								Type: "security_alert",
								Agent: "Aurelia",
								Action: "Tentativa de Execução Bloqueada",
								Payload: "Aguardando aprovação de plano pendente.",
								Timestamp: time.Now().Format("15:04:05"),
							})
							
							currentHistory = append(currentHistory, Message{
								Role:    "tool",
								Content: resultStr,
							})
							continue // Não muda a fase
						}
					}
					ctx = WithPrevPhase(ctx, newPhase)
				}
			}

			const maxToolResultLength = 32768
			if len(resultStr) > maxToolResultLength {
				resultStr = resultStr[:maxToolResultLength] + "\n\n... [TRUNCATED]"
			}

			currentHistory = append(currentHistory, Message{
				Role:       "tool",
				Content:    resultStr,
				ToolCallID: call.ID,
			})

			if call.Name == "handoff_to_agent" {
				return currentHistory, resultStr, nil
			}
		}
	}

	return currentHistory, "Max iterations reached", fmt.Errorf("max iterations reached")
}

func augmentSystemPromptWithPhase(basePrompt, phase string) string {
	var lines []string
	lines = append(lines, basePrompt)
	lines = append(lines, "\n---")
	lines = append(lines, fmt.Sprintf("FASE ATUAL DO WORKFLOW: %s", phase))
	
	switch phase {
	case "PLANNING":
		lines = append(lines, "DIRETRIZ: Você está na fase de desenho técnico. Antes de qualquer mudança, use a tool `propose_plan`.")
		lines = append(lines, "PROXIMO PASSO: Após chamar `propose_plan` e ter o OK, use `set_phase` para EXECUTION.")
	case "EXECUTION":
		lines = append(lines, "DIRETRIZ: Você está executando as mudanças aprovadas. Mantenha o foco no plano.")
		lines = append(lines, "PROXIMO PASSO: Quando terminar, use `set_phase` para VERIFICATION.")
	case "VERIFICATION":
		lines = append(lines, "DIRETRIZ: Valide as mudanças com testes ou leitura de logs. Use `log_verification` para concluir.")
	}
	
	return strings.Join(lines, "\n")
}
