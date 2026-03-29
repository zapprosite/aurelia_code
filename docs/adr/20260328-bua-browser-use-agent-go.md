# ADR 20260328: BUA-Style Browser Use Agent (Pure Go)

## Status
🟢 Proposto (P1 - Crítico)

## Contexto
Implementar computer use agent puro em Go inspirado no **BUA** (github.com/anxuanzi/bua). BUA é o projeto open source Go mais maduro para browser automation agentic, com loop observe→act completo e 20+ ferramentas nativas. Este ADR define a arquitetura primary agent que substitui a dependência do Stagehand MCP para tasks principais.

Baseado em: [BUA - AI-powered browser automation for Go](https://github.com/anxuanzi/bua) (86 ⭐, Jan 2026)

## Decisões Arquiteturais

### 1. Agent Loop (BUA Pattern)

```go
// internal/computer_use/agent.go
package computeruse

import (
    "context"
    "time"
)

type Agent struct {
    browser   *browser.Browser
    llm       *llm.Provider
    maxSteps  int
    timeout   time.Duration
    preset    Preset
}

type AgentState struct {
    Intent      string
    Steps      []ActionStep
    Screenshot string
    DOMState   string
    Done       bool
    Error      error
}

type ActionStep struct {
    Action      string
    Params      map[string]any
    Result      string
    Screenshot  string
    Timestamp   time.Time
}

const agentLoopTemplate = `
Você é um agente de computer use. Analise o estado atual e decida a próxima ação.

Histórico:
{{range .Steps}}
- {{.Action}}: {{.Params}} → {{.Result}}
{{end}}

Estado atual:
- Screenshot: [imagem]
- DOM: {{.DOMState}}

Ações disponíveis:
- navigate(url): Ir para URL
- click(selector): Clicar elemento
- type(selector, text): Digitar texto
- scroll(direction, amount): Rolar página
- extract(query): Extrair dados
- done(summary): Finalizar

Responda em JSON com action e reasoning.
`

func (a *Agent) Run(ctx context.Context, intent string) (*AgentState, error) {
    state := &AgentState{Intent: intent}

    for step := 0; step < a.maxSteps; step++ {
        // 1. OBSERVE - Captura estado
        screenshot, dom, err := a.browser.CaptureState(ctx)
        if err != nil {
            state.Error = fmt.Errorf("capture: %w", err)
            break
        }

        // 2. REASON - LLM decide ação
        action, err := a.llm.DecideAction(ctx, state.Intent, screenshot, dom, state.Steps)
        if err != nil {
            state.Error = fmt.Errorf("llm decision: %w", err)
            break
        }

        // 3. ACT - Executa ação
        if action.Type == "done" {
            state.Done = true
            break
        }

        result, err := a.browser.ExecuteAction(ctx, action)
        state.Steps = append(state.Steps, ActionStep{
            Action:     action.Type,
            Params:     action.Params,
            Result:     result,
            Screenshot: screenshot,
            Timestamp:  time.Now(),
        })
    }

    return state, nil
}
```

### 2. Browser Tools (20+ como BUA)

```go
// internal/computer_use/tools.go
package computeruse

var BrowserTools = []ToolDef{
    // Navegação
    {Name: "navigate", Description: "Navega para URL",
     Params: z.Object({"url": z.String().url()})},
    {Name: "go_back", Description: "Volta no histórico"},
    {Name: "go_forward", Description: "Avança no histórico"},
    {Name: "reload", Description: "Recarrega página"},

    // Interação
    {Name: "click", Description: "Clica no elemento",
     Params: z.Object({"selector": z.String()})},
    {Name: "double_click", Description: "Duplo clique",
     Params: z.Object({"selector": z.String()})},
    {Name: "hover", Description: "Mouse sobre elemento",
     Params: z.Object({"selector": z.String()})},
    {Name: "type_text", Description: "Digita texto",
     Params: z.Object({"selector": z.String(), "text": z.String()})},
    {Name: "press_key", Description: "Pressiona tecla",
     Params: z.Object({"key": z.String()})},
    {Name: "scroll", Description: "Rola página",
     Params: z.Object({"direction": z.String(), "amount": z.Number()})},

    // Extração
    {Name: "get_page_state", Description: "Captura estado completo"},
    {Name: "extract_content", Description: "Extrai conteúdo",
     Params: z.Object({"selector": z.String()})},
    {Name: "extract_text", Description: "Extrai texto",
     Params: z.Object({"query": z.String()})},
    {Name: "evaluate_js", Description: "Executa JS",
     Params: z.Object({"code": z.String()})},

    // Abas
    {Name: "new_tab", Description: "Nova aba",
     Params: z.Object({"url": z.String().url()})},
    {Name: "switch_tab", Description: "Muda aba",
     Params: z.Object({"index": z.Number()})},
    {Name: "close_tab", Description: "Fecha aba"},

    // Screenshots
    {Name: "screenshot", Description: "Captura screenshot"},
    {Name: "screenshot_element", Description: "Screenshot de elemento",
     Params: z.Object({"selector": z.String()})},

    // Done
    {Name: "done", Description: "Finaliza tarefa",
     Params: z.Object({"summary": z.String()})},
}

// Total: 20 ferramentas como BUA
```

### 3. Preset System (BUA Pattern)

```go
// internal/computer_use/preset.go
package computeruse

type Preset string

const (
    PresetSpeed   Preset = "speed"   // Fast, less thorough
    PresetCost    Preset = "cost"    // Conservative, cheaper
    PresetQuality Preset = "quality"  // Thorough, slower
)

type PresetConfig struct {
    MaxSteps         int
    ScreenshotBudget int
    StealthMode      bool
    HitLEnabled      bool
    Timeout          time.Duration
}

var Presets = map[Preset]PresetConfig{
    PresetSpeed: {
        MaxSteps:         10,
        ScreenshotBudget: 20,
        StealthMode:      false,
        HitLEnabled:      false,
        Timeout:          2 * time.Minute,
    },
    PresetCost: {
        MaxSteps:         15,
        ScreenshotBudget: 10,
        StealthMode:      true,
        HitLEnabled:      true,
        Timeout:          5 * time.Minute,
    },
    PresetQuality: {
        MaxSteps:         30,
        ScreenshotBudget: 50,
        StealthMode:      true,
        HitLEnabled:      false,
        Timeout:          10 * time.Minute,
    },
}
```

### 4. Screen State Tracking

```go
// internal/computer_use/state.go
package computeruse

type ScreenState struct {
    URL        string
    Title      string
    DOM        string
    Screenshot string
    Hash       string    // SHA256 para detectar mudanças
    Cursor     CursorPos
    Focus      string
    ScrollY    int
    Timestamp  time.Time
}

type ScreenHistory struct {
    States     []ScreenState
    maxHistory int
}

func (h *ScreenHistory) Add(state ScreenState) {
    h.States = append(h.States, state)
    if len(h.States) > h.maxHistory {
        h.States = h.States[len(h.States)-h.maxHistory:]
    }
}

func (h *ScreenHistory) HasChanged(newHash string) bool {
    if len(h.States) == 0 {
        return true
    }
    return h.States[len(h.States)-1].Hash != newHash
}
```

## Consequências

### Positivas
- Puro Go: sem dependência Node.js/Playwright
- 20+ ferramentas nativas como BUA
- Loop observe→act testado e funcional
- Presets permitem tuning fine-grained
- go-rod é mais leve que Playwright

### Negativas
- gemma3:27b não suporta visão nativa (precisa VL fallback)
- Implementação do zero leva tempo
- Sem comunidade extensa como Playwright

### Trade-offs
- go-rod vs Playwright: go-rod mais leve, Playwright mais maduro
- Agent loop interno vs MCP: comunicação direta é mais rápida

## Dependências
- ⚠️ `internal/browser/rod.go` - ADR separado
- ❌ `go-rod` - precisa adicionar ao go.mod
- ⚠️ `internal/llm/` - existente via gateway

## Referências
- [BUA - github.com/anxuanzi/bua](https://github.com/anxuanzi/bua)
- [ADR-20260328-mcp-tool-schema-computer-use.md](./20260328-mcp-tool-schema-computer-use.md)
- [ADR-20260328-go-rod-browser-layer.md](./20260328-go-rod-browser-layer.md)
- [ADR-20260328-computer-use-dependency-map.md](./20260328-computer-use-dependency-map.md)

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
