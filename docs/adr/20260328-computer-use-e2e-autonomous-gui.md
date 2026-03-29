# ADR 20260328: Computer Use E2E (Autonomous GUI Navigation)

## Status
🟢 Proposto (P2)

## Contexto
O Computer Use completo é: **Intent → Vision (screenshot) → LLM reasoning → Action (act/navigate) → Observation (screenshot) → Repeat**. Este ADR documenta o fluxo E2E e os padrões de segurança.

## Decisões Arquiteturais

### 1. Agent Loop para Computer Use

```go
// internal/computer_use/agent.go
type ComputerUseAgent struct {
    llm       *gateway.Provider
    vision    *vision.ScreenshotCapture
    stagehand *mcp.StagehandClient
    maxSteps  int
}

type AgentState struct {
    Intent     string
    Steps     []ActionStep
    Screenshot string
    Done      bool
}

type ActionStep struct {
    Action   string  // "navigate" | "act" | "extract"
    Params   map[string]any
    Result   string
    Screenshot string
}

func (a *ComputerUseAgent) Run(ctx context.Context, intent string) (*AgentState, error) {
    state := &AgentState{Intent: intent}

    for step := 0; step < a.maxSteps; step++ {
        // 1. Observar (screenshot)
        screenshot, err := a.vision.Capture(ctx)
        if err != nil {
            return nil, fmt.Errorf("vision capture: %w", err)
        }

        // 2. Raciocinar (LLM decide próxima ação)
        action, err := a.llm.DecideAction(ctx, state.Intent, screenshot, state.Steps)
        if err != nil {
            return nil, fmt.Errorf("llm decision: %w", err)
        }

        if action.Type == "done" {
            state.Done = true
            break
        }

        // 3. Agir (executa via Stagehand)
        result, err := a.executeAction(ctx, action)
        if err != nil {
            logger.Warn("action failed", "action", action.Type, "error", err)
            // Continua para permitir retry
        }

        state.Steps = append(state.Steps, ActionStep{
            Action:     action.Type,
            Params:     action.Params,
            Result:     result,
            Screenshot: screenshot,
        })
        state.Screenshot = screenshot
    }

    return state, nil
}
```

### 2. LLM Prompt para Decision

```go
// internal/computer_use/prompt.go
const computerUsePromptTemplate = `
Você é um agente de computer use. Analise o screenshot atual e decida a próxima ação.

Histórico de ações:
{{range .Steps}}
- {{.Action}}: {{.Params}} → {{.Result}}
{{end}}

Screenshot atual: [imagem em base64]

Decida a próxima ação:
- navigate(url): Ir para URL específica
- act(instruction): Executar ação (clique, digitação, etc)
- extract(instruction): Extrair informação da página
- done(summary): Finalizar e retornar resultado

Responda em JSON:
{
    "action": "act",
    "params": {"instruction": "clique no botão de login"},
    "reasoning": "O usuário quer fazer login..."
}
`
```

### 3. Safety Guardrails

```go
// internal/computer_use/guardrails.go
var dangerousPatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)(rm\s+-rf|dd\s+|mkfs|wipefs)`),
    regexp.MustCompile(`(?i)(sudo\s+rm|chmod\s+777|ssh\s+.*@)`),
    regexp.MustCompile(`(?i)(curl\s+.*\|\s*sh|wget\s+.*\|\s*bash)`),
}

var sensitiveActions = []string{
    "login", "password", "credit card", "ssn",
}

func (a *ComputerUseAgent) validateAction(action Action) error {
    // Check for dangerous commands
    for _, pattern := range dangerousPatterns {
        if pattern.MatchString(action.Params["instruction"]) {
            return fmt.Errorf("blocked dangerous action: %s", action.Type)
        }
    }

    // Always confirm sensitive actions
    for _, sensitive := range sensitiveActions {
        if strings.Contains(strings.ToLower(action.Params["instruction"]), sensitive) {
            logger.Warn("sensitive action detected", "action", action)
            // Envia confirmação para o usuário
        }
    }

    return nil
}
```

### 4. Computer Use Tools para o LLM

```go
// internal/computer_use/tools.go
var ComputerUseTools = []ToolDef{
    {
        Name:        "computer_navigate",
        Description: "Navega para uma URL no browser",
        Params: z.Object{
            "url": z.String().url().Describe("URL completa"),
        },
    },
    {
        Name:        "computer_act",
        Description: "Executa ação no browser (clique, digite, role)",
        Params: z.Object{
            "instruction": z.String().Describe("Instrução em linguagem natural"),
        },
    },
    {
        Name:        "computer_extract",
        Description: "Extrai dados estruturados da página",
        Params: z.Object{
            "query": z.String().Describe("O que extrair"),
        },
    },
    {
        Name:        "computer_screenshot",
        Description: "Captura screenshot atual da tela",
        Params: z.Object{},
    },
}
```

### 5. Telegram Integration

```go
// internal/computer_use/telegram.go
func (b *TelegramBot) handleComputerUseCommand(ctx context.Context, msg *tgbot.Message) {
    // Exemplo: "/navega para github.com"
    intent := strings.TrimPrefix(msg.Text, "/navega ")

    agent := computeruse.NewComputerUseAgent(b.llm, b.vision, b.stagehand)

    // Envia status inicial
    b.sendMessage(ctx, msg.Chat.ID, "🔍 Navegando...")

    state, err := agent.Run(ctx, intent)

    if state.Done {
        b.sendMessage(ctx, msg.Chat.ID, fmt.Sprintf("✅ %s", state.Steps[len(state.Steps)-1].Result))
    } else {
        b.sendScreenshot(ctx, msg.Chat.ID, state.Screenshot)
        b.sendMessage(ctx, msg.Chat.ID, fmt.Sprintf("⚠️ Não consegui completar. Resultado parcial:\n%s", formatSteps(state.Steps)))
    }
}
```

## Consequências

### Positivas
- Automação de tarefas web complexas
- Agent pode navegar, extrair dados, preencher formulários
- Integração com Telegram para comandos naturais

### Negativas
- Custo de tokens pode ser alto (screenshots + LLM calls)
- Browser automation é frágil (mudanças de UI quebram)
- Tempo de execução pode ser longo

### Trade-offs
- Autonomous vs Guided: Autonomous é mais rápido mas arriscado
- Screenshot per step vs Conditional: Sempre screenshot é mais seguro mas mais caro

## Dependências
- ⚠️ `internal/mcp/client.go` (ADR separado)
- ⚠️ `internal/vision/screenshot.go` (ADR separado)
- ⚠️ `internal/gateway/provider.go` (precisa expor DecideAction)
- ❌ Ollama com llava ou OpenRouter VL para visão

## Referências
- [ADR-20260328-mcp-go-client-stagehand-computer-use.md](./20260328-mcp-go-client-stagehand-computer-use.md)
- [ADR-20260328-vision-pipeline-computer-use.md](./20260328-vision-pipeline-computer-use.md)
- [ADR-20260328-computer-use-e2e-autonomous-gui.md](./20260328-computer-use-e2e-autonomous-gui.md)
- [ADR-20260328-implementacao-jarvis-voice-e-computer-use.md](./20260328-implementacao-jarvis-voice-e-computer-use.md)

## Links Obrigatórios
- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/neon-sentinel
**Progress**: 0%
