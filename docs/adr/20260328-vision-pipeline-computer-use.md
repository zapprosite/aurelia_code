# ADR 20260328: Vision Pipeline para Computer Use (Screenshot → LLM)

## Status
🟡 Proposto (P2)

## Contexto
Para que o gemma3 27b (ou Tier 1) possa fazer computer use, ele precisa ver o estado da tela. O pipeline de visão é: screenshot → base64 encoding → tool call → LLM processing → decision.

## Decisões Arquiteturais

### 1. Screenshot Capture

```go
// internal/vision/screenshot.go
package vision

import (
    "context"
    "encoding/base64"
    "os/exec"
)

type ScreenshotCapture struct {
    stagehandClient *mcp.StagehandClient
    format         string  // "png" or "jpeg"
    quality        int     // 1-100 for jpeg
}

func (s *ScreenshotCapture) Capture(ctx context.Context) (string, error) {
    // Via MCP tool ou via Playwright API direto
    screenshot, err := s.stagehandClient.Screenshot(ctx)
    if err != nil {
        return "", err
    }
    // Encode para base64 para enviar ao LLM
    return base64.StdEncoding.EncodeToString(screenshot), nil
}

func (s *ScreenshotCapture) CaptureRegion(ctx context.Context, x, y, w, h int) (string, error) {
    // Screenshot parcial para economia de tokens
    screenshot, err := s.stagehandClient.ScreenshotRegion(ctx, x, y, w, h)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(screenshot), nil
}
```

### 2. Vision Tool no Gateway

```go
// internal/gateway/vision.go
package gateway

type VisionTool struct {
    capture *vision.ScreenshotCapture
    encoder *llm.VisionEncoder
}

func (v *VisionTool) Name() string { return "computer_vision" }

func (v *VisionTool) Handle(ctx context.Context, params map[string]any) (string, error) {
    // Captura screenshot
    img, err := v.capture.Capture(ctx)
    if err != nil {
        return "", fmt.Errorf("screenshot capture: %w", err)
    }

    // Gera descrição textual via visão LLM
    description, err := v.encoder.Describe(ctx, img)
    if err != nil {
        return "", fmt.Errorf("vision encoding: %w", err)
    }

    return description, nil
}

// Tool definition para o LLM
var ComputerVisionToolDef = ToolDef{
    Name:        "computer_vision",
    Description: "Captura screenshot da tela atual e descreve o estado visual",
    Params: z.Object{
        "region": z.Optional(z.Object{
            "x": z.Number(),
            "y": z.Number(),
            "width": z.Number(),
            "height": z.Number(),
        }).Describe("Região opcional da tela para capturar"),
    },
}
```

### 3. Ollama Vision Support

O gemma3 27b via Ollama **não suporta** imagens diretamente. Opções:

**Opção A**: Usar modelo Ollama com suporte a visão:
```bash
ollama pull llava  # Modelo de visão leve
```

**Opção B**: Usar Tier 1 (OpenRouter) com visão:
```yaml
#configs/litellm/config.yaml
- model_name: "qwen-vl-32b"
  litellm_params:
    model: "openrouter/qwen/qwen-2.5-vl-32b-instruct"
```

**Opção C (Recomendada)**: LLaVA local como processador de visão:
```go
// Detectar automaticamente se Ollama tem llava
func DetectVisionModel(ctx context.Context) string {
    models, _ := listOllamaModels(ctx)
    for _, m := range models {
        if strings.HasPrefix(m.ID, "llava") {
            return m.ID
        }
    }
    return ""  // Fallback para OpenRouter
}
```

### 4. Screen State Tracking

```go
// internal/vision/state.go
type ScreenState struct {
    History    []ScreenSnapshot
    Cursor     Point
    Focus      string  // Elemento em foco
    ScrollY    int
    URL        string
}

type ScreenSnapshot struct {
    Timestamp  time.Time
    Base64    string
    Hash       string  // SHA256 para detectar mudanças
}

func (s *ScreenState) HasChanged(newHash string) bool {
    if len(s.History) == 0 {
        return true
    }
    return s.History[len(s.History)-1].Hash != newHash
}

func (s *ScreenState) Prune(maxHistory int) {
    if len(s.History) > maxHistory {
        s.History = s.History[len(s.History)-maxHistory:]
    }
}
```

### 5. Vision Budget

Para evitar custos excessivos com screenshots:

```go
const (
    MaxScreenshotsPerSession = 50
    ScreenshotCooldown       = 2 * time.Second
    MaxImageTokens          = 512  // Limite para gemma3 27b
)

func (v *VisionTool) EnforceBudget() error {
    // Tracking de screenshots por sessão
}
```

## Consequências

### Positivas
- Agent pode "ver" o estado da UI antes de agir
- Screenshots permite debugging visual
- Compressão otimizada reduz custo de tokens

### Negativas
- Screenshots são caros em tokens (~512-1024 por imagem)
- gemma3 27b não suporta visão nativamente
- Latência: captura + encode + processamento

### Trade-offs
- LLaVA local vs OpenRouter VL: LLaVA é sovereign mas menor
- Full screenshot vs region: Region economiza tokens mas perde contexto

## Dependências
- ⚠️ `internal/mcp/client.go` (precisa existir primeiro - ADR separado)
- ❌ `internal/vision/screenshot.go` (NÃO EXISTE)
- ❌ `internal/vision/state.go` (NÃO EXISTE)
- ❌ Ollama com llava ou OpenRouter API key com visão

## Referências
- [ADR-20260328-mcp-go-client-stagehand-computer-use.md](./20260328-mcp-go-client-stagehand-computer-use.md)
- [ADR-20260328-smart-router-litellm-gemma3-27b.md](./20260328-smart-router-litellm-gemma3-27b.md)
- [pkg/llm/vision.go](../../pkg/llm/vision.go)
- [LLaVA Ollama](https://ollama.ai/library/llava)

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
