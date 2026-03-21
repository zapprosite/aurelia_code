package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/dashboard"
	"github.com/kocar/aurelia/internal/observability"
)

// Loop executes the ReAct logic
type Loop struct {
	llm           LLMProvider
	registry      *ToolRegistry
	maxIterations int
	// catalog permite filtrar tools por relevância léxica ao prompt.
	// Se nil, todas as tools são enviadas à LLM (comportamento legado).
	catalog *ToolCatalog
	// catalogTopK é o número máximo de tools retornadas pelo catalog filter.
	catalogTopK int
}

// NewLoop constructs an agent loop. Pass -1 for unlimited iterations.
func NewLoop(llm LLMProvider, registry *ToolRegistry, maxIterations int) *Loop {
	if maxIterations == 0 {
		maxIterations = 5 // Fallback to PRD standard
	}
	return &Loop{
		llm:           llm,
		registry:      registry,
		maxIterations: maxIterations,
	}
}

// WithToolCatalog ativa o filtro inteligente de tools por relevância léxica.
// k define quantas tools são enviadas no máximo para a LLM por chamada.
func (l *Loop) WithToolCatalog(catalog *ToolCatalog, k int) *Loop {
	l.catalog = catalog
	l.catalogTopK = k
	return l
}

// Run executes the agent resolving loop on a given state of messages
func (l *Loop) Run(ctx context.Context, systemPrompt string, history []Message, allowedTools []string) ([]Message, string, error) {
	if _, ok := RunContextFromContext(ctx); !ok {
		ctx = WithRunContext(ctx, uuid.NewString())
	}
	logger := observability.Logger("agent.loop")
	if runID, ok := RunContextFromContext(ctx); ok {
		logger = logger.With(slog.String("run_id", runID))
	}

	currentHistory := make([]Message, len(history))
	copy(currentHistory, history)

	tools := l.registry.FilterDefinitions(allowedTools)
	// Sub-1: ToolCatalog — filtra tools por relevância léxica ao prompt da mensagem mais recente.
	// Quando o catalog está ativo, reduz drasticamente o número de tools no contexto da LLM.
	if l.catalog != nil && l.catalogTopK > 0 {
		var queryHint string
		for i := len(currentHistory) - 1; i >= 0; i-- {
			if currentHistory[i].Role == "user" {
				queryHint = currentHistory[i].Content
				break
			}
		}
		if queryHint != "" {
			filtered := l.catalog.MatchForTask(queryHint, l.catalogTopK)
			if len(filtered) > 0 {
				tools = filtered
			}
		}
	}
	// HARD OVERRIDE: Forçar identidade Linux/Ubuntu no topo de cada prompt
	systemPrompt = "### ENVIRONMENT CONTEXT\n- OS: Linux (Ubuntu 24.04 LTS)\n- SHELL: Bash (/bin/bash)\n- ARCH: amd64\n- RESTRICTION: NUNCA use PowerShell ou comandos Windows. Use apenas Bash nativo.\n\n" + systemPrompt

	systemPrompt = augmentSystemPromptWithToolGuidance(systemPrompt, tools)
	systemPrompt = augmentSystemPromptWithRuntimeCapabilities(systemPrompt, tools)

	// DIANOGSTIC LOG: Check what tools are actually being sent to the LLM
	var toolNames []string
	for _, t := range tools {
		toolNames = append(toolNames, t.Name)
	}
	logger.Debug("tools passed to LLM provider", slog.Any("tool_names", toolNames))

	iterations := 0

	for l.maxIterations < 0 || iterations < l.maxIterations {
		iterations++
		if l.maxIterations < 0 {
			logger.Debug("agent loop iteration", slog.Int("iteration", iterations), slog.String("max_iterations", "unbounded"))
		} else {
			logger.Debug("agent loop iteration", slog.Int("iteration", iterations), slog.Int("max_iterations", l.maxIterations))
		}

		// Check for context cancellation from the Time Tracker before hitting the LLM
		if ctx.Err() != nil {
			return currentHistory, "", fmt.Errorf("context cancelled by timer: %w", ctx.Err())
		}

		resp, err := l.llm.GenerateContent(ctx, systemPrompt, currentHistory, tools)
		if err != nil {
			return currentHistory, "", fmt.Errorf("provider error: %w", err)
		}

		// Publishes thought to dashboard
		dashboard.Publish(dashboard.Event{
			Type:      "agent_thought",
			Agent:     "Aurelia",
			Action:    "Pensando...",
			Payload:   resp.ReasoningContent,
			Timestamp: time.Now().Format(time.Kitchen),
		})

		// AI provided a final answer without tools
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

		// AI wants to call tools. Append its "Thought" / Request
		currentHistory = append(currentHistory, Message{
			Role:             "assistant",
			Content:          resp.Content, // Potentially empty or containing 'thought'
			ReasoningContent: resp.ReasoningContent,
			ToolCalls:        resp.ToolCalls, // Very important for provider API consistency
		})

		// Execute tools sequentially per PRD spec "- NG-02: as tool calls serão tratadas resolutivamente em cascata iterativa/síncrona na mesma Goroutine"
		for _, call := range resp.ToolCalls {
			logger.Info("executing tool", slog.String("tool_name", call.Name), slog.Any("arg_keys", observability.MapKeys(call.Arguments)))
			
			dashboard.Publish(dashboard.Event{
				Type:      "agent_tool",
				Agent:     "Aurelia",
				Action:    "Executando tool: " + call.Name,
				Payload:   call.Arguments,
				Timestamp: time.Now().Format(time.Kitchen),
			})

			resultStr, toolErr := l.registry.Execute(ctx, call.Name, call.Arguments)

			// If tool fails, return the error as text to the LLM
			if toolErr != nil {
				errorPayload, _ := json.Marshal(map[string]string{
					"error": toolErr.Error(),
				})
				resultStr = string(errorPayload)
			}

			// Append tool observation
			// SAFETY: Truncate tool output if it exceeds context limits (approx 8k-10k tokens)
			const maxToolResultLength = 32768
			if len(resultStr) > maxToolResultLength {
				logger.Warn("truncating tool output", slog.String("tool_name", call.Name), slog.Int("original_chars", len(resultStr)), slog.Int("max_chars", maxToolResultLength))
				resultStr = resultStr[:maxToolResultLength] + "\n\n... [TRUNCATED: O resultado desta ferramenta excedeu o limite de segurança do contexto. Se precisar de mais detalhes, leia partes específicas do arquivo ou diretório.]"
			}

			currentHistory = append(currentHistory, Message{
				Role:       "tool",
				Content:    resultStr,
				ToolCallID: call.ID,
			})

			// Se a ferramenta executada foi o Handoff, interrompemos o loop imediatamente
			if call.Name == "handoff_to_agent" {
				logger.Info("handoff detected, exiting loop", slog.String("target", call.Name))
				return currentHistory, resultStr, nil
			}
		}
	}

	return currentHistory, "Desculpe, desisti ou deu timeout no processamento pois falhei nas chamadas em MAX iteracoes.", fmt.Errorf("max iterations reached")
}

func augmentSystemPromptWithToolGuidance(systemPrompt string, tools []Tool) string {
	toolNames := make(map[string]bool, len(tools))
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	var sections []string

	if toolNames["run_command"] {
		sections = append(sections, "Se o usuario pedir para rodar, testar, iniciar, buildar, validar, verificar healthcheck ou inspecionar um projeto local, voce deve tentar usar `run_command` antes de responder com passos manuais. So ofereca execucao manual se `run_command` falhar, for bloqueado ou nao existir.")
		sections = append(sections, "ESTE AMBIENTE É LINUX (Ubuntu 24.04). NUNCA use PowerShell. Use apenas Bash/Sh para `run_command` e ferramentas de sistema.")
		sections = append(sections, "Nao diga que o ambiente esta bloqueado, que nao consegue executar processos ou que a execucao deve ser manual sem antes receber esse resultado explicitamente de uma tool. Se `run_command` nao retornou bloqueio ou erro, continue usando ferramentas.")
		sections = append(sections, "Se a tarefa exigir varias etapas locais, execute em sequencia: por exemplo subir o servico com `run_command`, depois testar endpoint com outro `run_command`, depois sintetizar o resultado observado.")
	}

	hasFilesystem := toolNames["read_file"] || toolNames["write_file"] || toolNames["list_dir"]
	if hasFilesystem {
		sections = append(sections, "As tools `read_file`, `write_file` e `list_dir` aceitam `workdir`. Sempre que estiver trabalhando em outro projeto ou pasta fora da raiz atual, informe `workdir` e use caminhos relativos a esse diretorio.")
	}

	if toolNames["run_command"] && hasFilesystem {
		sections = append(sections, "Se voce descobrir um diretorio de projeto via `run_command`, reutilize o mesmo `workdir` nas tools de filesystem para nao ler ou escrever no repositorio errado.")
	}
	if toolNames["spawn_agent"] {
		sections = append(sections, "Ao delegar trabalho com `spawn_agent` para outro projeto, passe o `workdir` canonico do projeto alvo. Nunca deixe subagente assumir por padrao a pasta do Aurelia como diretorio de trabalho.")
		sections = append(sections, "Se o usuario quiser interromper, pausar, retomar ou inspecionar a operacao do time, prefira usar `cancel_team`, `pause_team`, `resume_team` e `team_status` em vez de responder apenas em texto.")
	}

	if toolNames["create_schedule"] {
		sections = append(sections, "Se o usuario pedir lembretes, rotinas, tarefas recorrentes, avisos futuros ou qualquer acao para acontecer depois, voce deve considerar usar `create_schedule` em vez de apenas responder com texto.")
	}
	if toolNames["list_schedules"] || toolNames["pause_schedule"] || toolNames["resume_schedule"] || toolNames["delete_schedule"] {
		sections = append(sections, "Se o usuario perguntar quais agendamentos existem ou pedir para pausar, retomar ou remover uma rotina, use `list_schedules`, `pause_schedule`, `resume_schedule` e `delete_schedule` conforme a intencao.")
	}
	if toolNames["create_schedule"] || toolNames["list_schedules"] {
		sections = append(sections, "Nao exija comandos como `/cron`. A interface correta e linguagem natural; as tools de scheduling existem para voce transformar a intencao do usuario em operacoes reais.")
	}

	if len(sections) == 0 {
		return systemPrompt
	}

	return strings.TrimSpace(systemPrompt) + "\n\n# TOOL USAGE GUIDE\n" + strings.Join(sections, "\n")
}

func augmentSystemPromptWithRuntimeCapabilities(systemPrompt string, tools []Tool) string {
	var lines []string
	lines = append(lines, "# RUNTIME CAPABILITIES")
	if len(tools) == 0 {
		lines = append(lines, "Nenhuma tool esta disponivel neste runtime para esta execucao.")
		lines = append(lines, "Se houver duvida sobre capacidades, considere esta secao como fonte canonica em vez de assumir ferramentas inexistentes.")
	} else {
		lines = append(lines, "Tools disponiveis nesta execucao:")
		names := make([]string, 0, len(tools))
		for _, tool := range tools {
			names = append(names, tool.Name)
		}
		sort.Strings(names)
		for _, name := range names {
			lines = append(lines, "- "+name)
		}
		lines = append(lines, "Considere esta lista como a fonte canonica das capacidades reais deste runtime.")
	}

	base := strings.TrimSpace(systemPrompt)
	if base == "" {
		return strings.Join(lines, "\n")
	}
	return base + "\n\n" + strings.Join(lines, "\n")
}
