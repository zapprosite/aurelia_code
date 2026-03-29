# ADR 20260328: Normalized Coordinates + HitL Safety

## Status
🟢 Proposto (P1 - Crítico)

## Contexto
Implementar sistema de coordenadas normalizadas (0-999) para computer use, evitando problemas de precisão entre diferentes resoluções. Também implementar Human-in-the-Loop (HitL) para ações sensíveis e bloqueio de padrões perigosos.

Baseado em:
- [vnc-use](https://github.com/mayflower/vnc-use) - Normalized coordinates pattern
- [computer-use-mcp](https://github.com/domdomegg/computer-use-mcp) - HitL pattern

## Decisões Arquiteturais

### 1. Coordinate System (0-999 Grid)

```go
// internal/computer_use/coordinates.go
package computeruse

import (
    "math"
    "image"
)

// CoordinateSystem normaliza coordenadas para grid 0-999
// Padrão de vnc-use: escala universal independente de resolução
type CoordinateSystem struct {
    ScreenWidth  int
    ScreenHeight int
}

func NewCoordinateSystem(width, height int) *CoordinateSystem {
    return &CoordinateSystem{
        ScreenWidth:  width,
        ScreenHeight: height,
    }
}

// NormalizeToGrid converte pixels para grid 0-999
// 0 = topo/esquerda, 999 = fundo/direita
func (c *CoordinateSystem) NormalizeToGrid(x, y int) (int, int) {
    normX := int(math.Round(float64(x) / float64(c.ScreenWidth) * 999))
    normY := int(math.Round(float64(y) / float64(c.ScreenHeight) * 999))

    // Clamp para bounds
    normX = clamp(normX, 0, 999)
    normY = clamp(normY, 0, 999)

    return normX, normY
}

// DenormalizeFromGrid converte grid 0-999 para pixels
func (c *CoordinateSystem) DenormalizeFromGrid(normX, normY int) (int, int) {
    x := int(math.Round(float64(normX) / 999 * float64(c.ScreenWidth)))
    y := int(math.Round(float64(normY) / 999 * float64(c.ScreenHeight)))
    return clamp(x, 0, c.ScreenWidth), clamp(y, 0, c.ScreenHeight)
}

// RegionToGrid converte região image.Rectangle para grid
func (c *CoordinateSystem) RegionToGrid(region image.Rectangle) GridRegion {
    x1, y1 := c.NormalizeToGrid(region.Min.X, region.Min.Y)
    x2, y2 := c.NormalizeToGrid(region.Max.X, region.Max.Y)
    return GridRegion{
        X:      x1,
        Y:      y1,
        Width:  x2 - x1,
        Height: y2 - y1,
    }
}

type GridRegion struct {
    X, Y   int
    Width  int
    Height int
}

func clamp(val, min, max int) int {
    if val < min {
        return min
    }
    if val > max {
        return max
    }
    return val
}
```

### 2. HitL Confirmation Gate

```go
// internal/computer_use/hitl.go
package computeruse

import (
    "context"
    "fmt"
)

// HitLGate confirmação humana para ações sensíveis
type HitLGate struct {
    enabled      bool
    confirmCh    chan Confirmation
    timeout      time.Duration
    pendingCount int
}

type Confirmation struct {
    Action   string
    Details  string
    Approved bool
}

func NewHitLGate(enabled bool, timeout time.Duration) *HitLGate {
    return &HitLGate{
        enabled:   enabled,
        confirmCh: make(chan Confirmation, 1),
        timeout:   timeout,
    }
}

// RequestConfirmation bloqueia até usuário confirmar ou timeout
func (h *HitLGate) RequestConfirmation(ctx context.Context, action, details string) error {
    if !h.enabled {
        return nil // HitL desabilitado = auto-aprova
    }

    h.pendingCount++
    defer func() { h.pendingCount-- }()

    select {
    case confirm := <-h.confirmCh:
        if !confirm.Approved {
            return fmt.Errorf("action rejected by user: %s", action)
        }
        return nil

    case <-time.After(h.timeout):
        return fmt.Errorf("hitl timeout: user did not respond in %v", h.timeout)

    case <-ctx.Done():
        return ctx.Err()
    }
}

// Approve/Reject para integração com Telegram/UI
func (h *HitLGate) Approve() {
    h.confirmCh <- Confirmation{Approved: true}
}

func (h *HitLGate) Reject() {
    h.confirmCh <- Confirmation{Approved: false}
}

// IsPending retorna se há confirmação pendente
func (h *HitLGate) IsPending() bool {
    return h.pendingCount > 0
}
```

### 3. Dangerous Pattern Blocking

```go
// internal/computer_use/safety.go
package computeruse

import (
    "regexp"
    "strings"
)

// SafetyGuard bloqueia ações perigosas automaticamente
type SafetyGuard struct {
    dangerousPatterns []*regexp.Regexp
    sensitiveKeywords []string
    log               *slog.Logger
}

func NewSafetyGuard() *SafetyGuard {
    return &SafetyGuard{
        dangerousPatterns: []*regexp.Regexp{
            // Destruição de dados
            regexp.MustCompile(`(?i)\brm\s+-rf\b`),
            regexp.MustCompile(`(?i)\bdd\s+if\b`),
            regexp.MustCompile(`(?i)\bmkfs\b`),
            regexp.MustCompile(`(?i)\bwipefs\b`),

            // Escalação de privilégios
            regexp.MustCompile(`(?i)\bsudo\s+rm\b`),
            regexp.MustCompile(`(?i)\bsudo\s+chmod\s+777\b`),
            regexp.MustCompile(`(?i)\bchmod\s+-R\s+777\b`),

            // Remote execution
            regexp.MustCompile(`(?i)\bssh\s+.*@`),
            regexp.MustCompile(`(?i)\bcurl\s+.*\|\s*(sh|bash|bash\s+-c)`),
            regexp.MustCompile(`(?i)\bwget\s+.*\|\s*(sh|bash|bash\s+-c)`),

            // Credential theft
            regexp.MustCompile(`(?i)(password|passwd|pwd)\s*=\s*['"].*['"]`),
            regexp.MustCompile(`(?i)\.netrc`),
        },
        sensitiveKeywords: []string{
            "credit card",
            "social security",
            "ssn",
            "bank account",
            "api key",
            "secret key",
            "private key",
        },
    }
}

// ValidateAction verifica se ação é segura
func (s *SafetyGuard) ValidateAction(action string, params map[string]any) (bool, string) {
    // Verifica padrões perigosos em todos os params
    for _, val := range params {
        strVal := fmt.Sprintf("%v", val)

        for _, pattern := range s.dangerousPatterns {
            if pattern.MatchString(strVal) {
                s.log.Warn("dangerous pattern blocked",
                    "pattern", pattern.String(),
                    "value", truncate(strVal, 50))
                return false, fmt.Sprintf("blocked: dangerous pattern %s", pattern.String())
            }
        }
    }

    // Verifica palavras sensíveis
    actionLower := strings.ToLower(action)
    for _, keyword := range s.sensitiveKeywords {
        if strings.Contains(actionLower, keyword) {
            s.log.Warn("sensitive action detected",
                "keyword", keyword,
                "action", action)
            return false, fmt.Sprintf("requires confirmation: sensitive keyword '%s'", keyword)
        }
    }

    return true, ""
}

// truncate encurta string para logging
func truncate(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen] + "..."
}
```

### 4. Safety Middleware

```go
// internal/computer_use/middleware.go
package computeruse

import (
    "context"
    "fmt"
)

// SafetyMiddleware integra safety + hitl no agent loop
type SafetyMiddleware struct {
    guard     *SafetyGuard
    hitl      *HitLGate
    next      AgentExecutor
}

type AgentExecutor func(ctx context.Context, action Action) (string, error)

func NewSafetyMiddleware(guard *SafetyGuard, hitl *HitLGate, next AgentExecutor) *SafetyMiddleware {
    return &SafetyMiddleware{
        guard: guard,
        hitl:  hitl,
        next:  next,
    }
}

func (m *SafetyMiddleware) Execute(ctx context.Context, action Action) (string, error) {
    // 1. Valida padrões perigosos
    safe, reason := m.guard.ValidateAction(action.Type, action.Params)
    if !safe {
        return "", fmt.Errorf("safety blocked: %s", reason)
    }

    // 2. Se ação sensíveis, pede confirmação
    if isSensitive(action) {
        if err := m.hitl.RequestConfirmation(ctx, action.Type, formatAction(action)); err != nil {
            return "", fmt.Errorf("hitl rejected: %w", err)
        }
    }

    // 3. Executa ação
    return m.next(ctx, action)
}

func isSensitive(action Action) bool {
    sensitiveTypes := map[string]bool{
        "type_text":    true,
        "click":        true,
        "extract":      true,
        "navigate":     true,
    }

    // Check params for sensitive data
    for _, val := range action.Params {
        str := strings.ToLower(fmt.Sprintf("%v", val))
        sensitive := []string{"password", "credit", "ssn", "key", "secret"}
        for _, s := range sensitive {
            if strings.Contains(str, s) {
                return true
            }
        }
    }

    return sensitiveTypes[action.Type]
}

func formatAction(action Action) string {
    return fmt.Sprintf("%s(%v)", action.Type, action.Params)
}
```

## Consequências

### Positivas
- Normalized coordinates funciona em qualquer resolução
- HitL previne ações destructivas acidentais
- Dangerous pattern blocking é automático
- Logging de segurança para auditoria

### Negativas
- HitL adiciona latência (timeout de ~30s)
- Coordenadas normalizadas perdem precisão
- Algumas ações legítimas podem ser bloqueadas

### Trade-offs
- Grid 0-999 vs pixels: mais portável, menos preciso
- Auto-block vs confirm: mais seguro vs mais lento

## Dependências
- ⚠️ `internal/computer_use/agent.go` - ADR separado
- ⚠️ `internal/computer_use/hitl.go` - neste ADR

## Referências
- [vnc-use - Normalized coordinates](https://github.com/mayflower/vnc-use)
- [computer-use-mcp - HitL pattern](https://github.com/domdomegg/computer-use-mcp)
- [ADR-20260328-mcp-tool-schema-computer-use.md](./20260328-mcp-tool-schema-computer-use.md)
- [ADR-20260328-bua-browser-use-agent-go.md](./20260328-bua-browser-use-agent-go.md)

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
