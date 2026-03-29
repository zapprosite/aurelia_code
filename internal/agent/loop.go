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
	llm             LLMProvider
	registry        *ToolRegistry
	maxIterations   int
	memoryAssembler MemoryAssembler
	toolCatalog     *ToolCatalog
	toolCatalogTopK int
}

type MemoryAssembler interface {
	AssembleContext(ctx context.Context, query string) string
}

type BotScopedMemoryAssembler interface {
	AssembleContextForBot(ctx context.Context, botID, query string) string
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
	opts := LoopOptions{
		SystemPrompt:    systemPrompt,
		InitialHistory:  history,
		MaxIterations:   l.maxIterations,
		ToolDefinitions: l.resolveToolDefinitions(history, allowedTools),
	}
	return l.RunWithOptions(ctx, opts)
}

func (l *Loop) WithMemoryAssembler(assembler MemoryAssembler) *Loop {
	l.memoryAssembler = assembler
	return l
}

func (l *Loop) WithToolCatalog(catalog *ToolCatalog, topK int) *Loop {
	l.toolCatalog = catalog
	l.toolCatalogTopK = topK
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
	ch, err := l.RunWithOptionsStream(ctx, opts)
	if err != nil {
		return opts.InitialHistory, "", err
	}

	var currentHistory []Message
	var lastContent string
	for resp := range ch {
		if resp.Err != nil {
			return currentHistory, "", resp.Err
		}
		if resp.Done && resp.History != nil {
			currentHistory = resp.History
		}
		if resp.Content != "" {
			lastContent = resp.Content
		}
	}
	return currentHistory, lastContent, nil
}

// RunWithOptionsStream é a nova versão industrial que suporta streaming em tempo real
func (l *Loop) RunWithOptionsStream(ctx context.Context, opts LoopOptions) (<-chan StreamResponse, error) {
	logger := observability.Logger("agent.loop")
	ch := make(chan StreamResponse, 100)

	if opts.MaxIterations <= 0 {
		opts.MaxIterations = l.maxIterations
	}
	if len(opts.ToolDefinitions) == 0 {
		opts.ToolDefinitions = l.resolveToolDefinitions(opts.InitialHistory, nil)
	}
	opts.SystemPrompt = augmentSystemPromptWithRuntimeCapabilities(opts.SystemPrompt, opts.ToolDefinitions)
	opts.SystemPrompt = l.augmentSystemPromptWithMemory(ctx, opts.SystemPrompt, opts.Task, opts.InitialHistory)

	go func() {
		defer close(ch)
		currentHistory := opts.InitialHistory

		for i := 0; i < opts.MaxIterations; i++ {
			select {
			case <-ctx.Done():
				ch <- StreamResponse{Err: ctx.Err()}
				return
			default:
			}

			if opts.InterruptHandler != nil && opts.InterruptHandler() {
				ch <- StreamResponse{Content: "Interrompido pelo usuário.", Done: true, History: currentHistory}
				return
			}

			currentPhase, _ := PrevPhaseFromContext(ctx)
			dynamicSystemPrompt := augmentSystemPromptWithPhase(opts.SystemPrompt, currentPhase)

			stream, err := l.llm.GenerateStream(ctx, dynamicSystemPrompt, currentHistory, opts.ToolDefinitions)
			if err != nil {
				ch <- StreamResponse{Err: fmt.Errorf("generate stream error: %w", err)}
				return
			}

			var fullContent strings.Builder
			var toolCalls []ToolCall

			for resp := range stream {
				if resp.Err != nil {
					ch <- resp
					return
				}
				if resp.Content != "" {
					fullContent.WriteString(resp.Content)
					ch <- resp // Forward token to caller (Jarvis TTS/UI)
				}
				if resp.Done {
					toolCalls = resp.ToolCalls
				}
			}

			// Se não tem tool calls, terminou a jornada
			if len(toolCalls) == 0 {
				currentHistory = append(currentHistory, Message{
					Role:    "assistant",
					Content: fullContent.String(),
				})
				ch <- StreamResponse{Done: true, History: currentHistory}
				return
			}

			// Execução de tools
			currentHistory = append(currentHistory, Message{
				Role:      "assistant",
				Content:   fullContent.String(),
				ToolCalls: toolCalls,
			})

			for _, call := range toolCalls {
				logger.Info("executing tool", slog.String("tool_name", call.Name))
				dashboard.Publish(dashboard.Event{
					Type:      "agent_tool",
					Agent:     "Aurelia",
					Action:    "Executando tool: " + call.Name,
					Payload:   call.Arguments,
					Timestamp: time.Now().Format("15:04:05"),
				})

				resultStr, toolErr := l.registry.Execute(ctx, call.Name, call.Arguments)
				if toolErr != nil {
					errorPayload, _ := json.Marshal(map[string]string{"error": toolErr.Error()})
					resultStr = string(errorPayload)
				}

				// Lógica de fases
				if call.Name == "set_phase" {
					if newPhase, ok := call.Arguments["phase"].(string); ok {
						if newPhase == "EXECUTION" && plan.GlobalPlanStore.HasPending() {
							resultStr = "PLANOS PENDENTES: Aprovação necessária."
						} else {
							ctx = WithPrevPhase(ctx, newPhase)
						}
					}
				}

				currentHistory = append(currentHistory, Message{
					Role:       "tool",
					Content:    resultStr,
					ToolCallID: call.ID,
				})

				if call.Name == "handoff_to_agent" {
					ch <- StreamResponse{Done: true, History: currentHistory}
					return
				}
			}
		}
		ch <- StreamResponse{Err: fmt.Errorf("max iterations reached"), History: currentHistory}
	}()

	return ch, nil
}

func (l *Loop) resolveToolDefinitions(history []Message, allowedTools []string) []Tool {
	if l == nil || l.registry == nil {
		return nil
	}
	if len(allowedTools) > 0 {
		return l.registry.FilterDefinitions(allowedTools)
	}
	if l.toolCatalog != nil {
		if query := latestQueryableMessage(history); query != "" {
			if tools := l.toolCatalog.MatchForTask(query, l.toolCatalogTopK); len(tools) > 0 {
				return tools
			}
		}
	}
	return l.registry.GetDefinitions()
}

func (l *Loop) augmentSystemPromptWithMemory(ctx context.Context, basePrompt, task string, history []Message) string {
	if l == nil || l.memoryAssembler == nil {
		return basePrompt
	}

	query := strings.TrimSpace(task)
	if query == "" {
		query = latestQueryableMessage(history)
	}
	if query == "" {
		return basePrompt
	}

	var memoryContext string
	if scoped, ok := l.memoryAssembler.(BotScopedMemoryAssembler); ok {
		botID, _ := BotContextFromContext(ctx)
		memoryContext = scoped.AssembleContextForBot(ctx, botID, query)
	} else {
		memoryContext = l.memoryAssembler.AssembleContext(ctx, query)
	}
	memoryContext = strings.TrimSpace(memoryContext)
	if memoryContext == "" {
		return basePrompt
	}

	return strings.TrimSpace(basePrompt) + "\n\n---\n# MEMORY CONTEXT\nUse este contexto recuperado apenas quando ele for relevante e consistente com a solicitacao atual.\n\n" + memoryContext
}

func latestQueryableMessage(history []Message) string {
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role != "user" {
			continue
		}
		if text := strings.TrimSpace(history[i].Content); text != "" {
			return text
		}
		for _, part := range history[i].Parts {
			if part.Type != ContentPartText {
				continue
			}
			if text := strings.TrimSpace(part.Text); text != "" {
				return text
			}
		}
	}
	return ""
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
