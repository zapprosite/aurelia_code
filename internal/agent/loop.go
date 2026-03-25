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

// GetLLMProvider returns the underlying LLM provider
func (l *Loop) GetLLMProvider() LLMProvider {
	return l.llm
}

// Run é a interface clássica do loop, mantida para compatibilidade
func (l *Loop) Run(ctx context.Context, systemPrompt string, history []Message, allowedTools []string) ([]Message, string, error) {
	// Filtrar definições de tools baseado em allowedTools
	var filteredTools []Tool
	if len(allowedTools) > 0 {
		allTools := l.registry.GetDefinitions()
		allowedSet := make(map[string]bool)
		for _, t := range allowedTools {
			allowedSet[t] = true
		}
		for _, tool := range allTools {
			if allowedSet[tool.Name] {
				filteredTools = append(filteredTools, tool)
			}
		}
	} else {
		filteredTools = l.registry.GetDefinitions()
	}

	// Aumentar o system prompt com orientações de runtime capabilities
	augmentedPrompt := augmentSystemPromptWithRuntimeCapabilities(systemPrompt, filteredTools)

	opts := LoopOptions{
		SystemPrompt:     augmentedPrompt,
		InitialHistory:   history,
		MaxIterations:    l.maxIterations,
		ToolDefinitions:  filteredTools,
	}
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
	ToolDefinitions  []Tool // Filtered tools for this run
}

func (l *Loop) RunWithOptions(ctx context.Context, opts LoopOptions) ([]Message, string, error) {
	logger := observability.Logger("agent.loop")
	currentHistory := opts.InitialHistory

	// Atualizar squad status quando inicia
	agentName, _ := AgentContextFromContext(ctx)
	if agentName != "" {
		UpdateSquadAgentStatus(agentName, "busy", 50)
		defer UpdateSquadAgentStatus(agentName, "online", 0)
	}

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
		// O provedor de LLM PRECISA das ferramentas registradas (usar filtradas se disponíveis)
		toolDefs := opts.ToolDefinitions
		if len(toolDefs) == 0 {
			toolDefs = l.registry.GetDefinitions()
		}
		resp, err := l.llm.GenerateContent(ctx, dynamicSystemPrompt, currentHistory, toolDefs)
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

func augmentSystemPromptWithRuntimeCapabilities(basePrompt string, tools []Tool) string {
	if len(tools) == 0 {
		return basePrompt
	}

	var lines []string
	lines = append(lines, basePrompt)
	lines = append(lines, "\n---")
	lines = append(lines, "# RUNTIME CAPABILITIES")
	lines = append(lines, "")
	lines = append(lines, "Se o usuario pedir para rodar comandos, ler/escrever arquivos, ou fazer outras operações no sistema, você TEM permissão para usar estas ferramentas:")
	lines = append(lines, "")

	for _, tool := range tools {
		lines = append(lines, fmt.Sprintf("- %s: %s", tool.Name, tool.Description))
	}

	lines = append(lines, "")
	lines = append(lines, "REGRAS IMPORTANTES:")
	lines = append(lines, "1. Use `workdir` para entender o diretório de trabalho atual antes de qualquer operação de arquivo.")
	lines = append(lines, "2. Para operações de longa duração ou repetidas, prefira criar schedules com `create_schedule` em vez de executar direto.")
	lines = append(lines, "3. Nao diga que o ambiente esta bloqueado ou que nao pode executar comandos — você TEM as ferramentas listadas acima.")
	lines = append(lines, "4. Sempre confirme caminhos completos e permissões antes de modificar arquivos críticos.")

	return strings.Join(lines, "\n")
}
