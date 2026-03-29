# ADR 20260328: MCP Tool Schema para Computer Use

## Status
🟢 Proposto (P1 - Crítico)

## Contexto
Definir o schema de tools MCP compatível com o padrão computer-use. Baseado no projeto **computer-use-mcp** (domdomegg/computer-use-mcp, 171 ⭐) que é a reimplementação open source mais próxima do Anthropic computer use tool. Este ADR garante interoperabilidade com Claude Desktop e outros clientes MCP.

Baseado em: [computer-use-mcp](https://github.com/domdomegg/computer-use-mcp) (171 ⭐, v1.7.1, Mar 2026)

## Decisões Arquiteturais

### 1. MCP Tool Schema (computer-use-mcp Pattern)

```go
// internal/computer_use/mcp_tools.go
package computeruse

import "github.com/nickysemenza/gomcpm/mcp/tool"

// ComputerUseTools define as tools MCP compatíveis com computer use
var ComputerUseTools = []tool.Definition{
    // Navegação
    {
        Name:        "navigate",
        Description: "Navega o browser para URL específica",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "url": map[string]any{
                    "type":        "string",
                    "description": "URL completa para navegar",
                },
            },
            "required": []string{"url"},
        },
    },

    // Ações de mouse
    {
        Name:        "mouse_move",
        Description: "Move o cursor para coordenadas normalizadas (0-999)",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "x": map[string]any{
                    "type":        "number",
                    "description": "Coordenada X normalizada (0-999)",
                    "minimum":     0,
                    "maximum":     999,
                },
                "y": map[string]any{
                    "type":        "number",
                    "description": "Coordenada Y normalizada (0-999)",
                    "minimum":     0,
                    "maximum":     999,
                },
            },
            "required": []string{"x", "y"},
        },
    },
    {
        Name:        "mouse_click",
        Description: "Clique do mouse em coordenadas normalizadas",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "x":      map[string]any{"type": "number"},
                "y":      map[string]any{"type": "number"},
                "button": map[string]any{
                    "type":        "string",
                    "enum":        []string{"left", "right", "middle"},
                    "description": "Botão do mouse",
                },
            },
            "required": []string{"x", "y"},
        },
    },

    // Teclado
    {
        Name:        "key_type",
        Description: "Digita texto",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "text": map[string]any{
                    "type":        "string",
                    "description": "Texto para digitar",
                },
            },
            "required": []string{"text"},
        },
    },
    {
        Name:        "key_shortcut",
        Description: "Atalho de teclado",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "keys": map[string]any{
                    "type":        "string",
                    "description": "Atalho (ex: Cmd+C, Ctrl+V)",
                },
            },
            "required": []string{"keys"},
        },
    },

    // Screenshot
    {
        Name:        "screenshot",
        Description: "Captura screenshot da tela atual",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "region": map[string]any{
                    "type": "object",
                    "properties": map[string]any{
                        "x":      map[string]any{"type": "number"},
                        "y":      map[string]any{"type": "number"},
                        "width":  map[string]any{"type": "number"},
                        "height": map[string]any{"type": "number"},
                    },
                    "description": "Região opcional para capturar",
                },
            },
        },
    },

    // Extração
    {
        Name:        "extract",
        Description: "Extrai dados estruturados da página atual",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "query": map[string]any{
                    "type":        "string",
                    "description": "Query em linguagem natural para extração",
                },
            },
            "required": []string{"query"},
        },
    },
}
```

### 2. Normalized Coordinates System (de vnc-use)

```go
// internal/computer_use/coordinates.go
package computeruse

import (
    "image"
    "math"
)

// CoordinateSystem normaliza coordenadas para 0-999
// Evita problemas de precisão entre diferentes resoluções
type CoordinateSystem struct {
    ScreenWidth  int
    ScreenHeight int
}

// NormalizeToGrid converte coordenadas reais para grid 0-999
func (c *CoordinateSystem) NormalizeToGrid(x, y int) (int, int) {
    normX := int(math.Round(float64(x) / float64(c.ScreenWidth) * 999))
    normY := int(math.Round(float64(y) / float64(c.ScreenHeight) * 999))

    // Clamp para 0-999
    normX = max(0, min(999, normX))
    normY = max(0, min(999, normY))

    return normX, normY
}

// DenormalizeFromGrid converte grid 0-999 de volta para pixels
func (c *CoordinateSystem) DenormalizeFromGrid(normX, normY int) (int, int) {
    x := int(math.Round(float64(normX) / 999 * float64(c.ScreenWidth)))
    y := int(math.Round(float64(normY) / 999 * float64(c.ScreenHeight)))
    return x, y
}

// RegionToNormalized converte região para grid
func (c *CoordinateSystem) RegionToNormalized(region image.Rectangle) map[string]int {
    nx1, ny1 := c.NormalizeToGrid(region.Min.X, region.Min.Y)
    nx2, ny2 := c.NormalizeToGrid(region.Max.X, region.Max.Y)
    return map[string]int{
        "x":      nx1,
        "y":      ny1,
        "width":  nx2 - nx1,
        "height": ny2 - ny1,
    }
}
```

### 3. MCP Client Integration

```go
// internal/mcp/computer_use.go
package mcp

import (
    "context"
    "encoding/json"
)

// ComputerUseMCPClient wrapper para tools de computer use
type ComputerUseMCPClient struct {
    client  *Client
    coords  *computeruse.CoordinateSystem
}

func (c *ComputerUseMCPClient) Navigate(ctx context.Context, url string) error {
    return c.client.CallTool(ctx, "navigate", map[string]any{"url": url})
}

func (c *ComputerUseMCPClient) MouseMove(ctx context.Context, x, y int) error {
    // Normaliza antes de enviar
    nx, ny := c.coords.NormalizeToGrid(x, y)
    return c.client.CallTool(ctx, "mouse_move", map[string]any{"x": nx, "y": ny})
}

func (c *ComputerUseMCPClient) Screenshot(ctx context.Context) ([]byte, error) {
    result, err := c.client.CallTool(ctx, "screenshot", nil)
    if err != nil {
        return nil, err
    }
    return c.decodeResult(result)
}

func (c *ComputerUseMCPClient) Extract(ctx context.Context, query string) (string, error) {
    result, err := c.client.CallTool(ctx, "extract", map[string]any{"query": query})
    if err != nil {
        return "", err
    }
    var data struct {
        Content []json.RawMessage `json:"content"`
    }
    json.Unmarshal(result, &data)
    if len(data.Content) > 0 {
        return string(data.Content[0]), nil
    }
    return "", nil
}
```

## Consequências

### Positivas
- Compatível com Claude Desktop MCP
- Schema testado em produção (171 ⭐)
- Normalized coordinates evita bugs de resolução
- Tool definitions auto-documentadas

### Negativas
- Adiciona camada de abstração
- Conversão de coordenadas adiciona latência
- Nem todos os clientes MCP suportam todas as features

### Trade-offs
- Grid 0-999 vs pixels reais: grid é mais portável mas menos preciso
- Tool por ação vs macro actions: mais granular mas mais LLM calls

## Dependências
- ⚠️ `internal/mcp/client.go` - existente
- ⚠️ `github.com/nickysemenza/gomcpm` - MCP SDK Go
- ⚠️ `internal/computer_use/tools.go` - ADR separado

## Referências
- [computer-use-mcp - github.com/domdomegg/computer-use-mcp](https://github.com/domdomegg/computer-use-mcp)
- [vnc-use - Normalized coordinates](https://github.com/mayflower/vnc-use)
- [ADR-20260328-bua-browser-use-agent-go.md](./20260328-bua-browser-use-agent-go.md)
- [ADR-20260328-go-rod-browser-layer.md](./20260328-go-rod-browser-layer.md)

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
