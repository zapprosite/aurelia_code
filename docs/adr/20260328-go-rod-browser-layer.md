# ADR 20260328: go-rod Browser Layer

## Status
🟢 Proposto (P1 - Crítico)

## Contexto
Implementar a camada de browser automation usando **go-rod** (github.com/go-rod/rod), inspirado no BUA. go-rod é um library Go pura que controla Chrome/Chromium via Chrome DevTools Protocol (CDP), sem necessidade de Node.js ou Playwright. É mais leve e idiomático em Go.

Baseado em: [BUA](https://github.com/anxuanzi/bua) que usa go-rod internamente

## Decisões Arquiteturais

### 1. Browser Layer (go-rod)

```go
// internal/browser/rod.go
package browser

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "image"
    "image/png"
    "os"

    "github.com/go-rod/rod"
    "github.com/go-rod/rod/lib/launcher"
    "github.com/go-rod/rod/lib/proto"
)

type RodBrowser struct {
    browser  *rod.Browser
    page     *rod.Page
    stealth  bool
    launched bool
}

type BrowserConfig struct {
    Headless     bool
    StealthMode  bool
    UserDataDir  string
    ProxyURL     string
    WindowSize   image.Rectangle
    UserAgent    string
}

func NewRodBrowser(ctx context.Context, cfg BrowserConfig) (*RodBrowser, error) {
    // Launcher com opções
    l := launcher.New().
        Headless(cfg.Headless).
        Set("disable-blink-features", "AutomationControlled")

    if cfg.UserDataDir != "" {
        l = l.UserDataDir(cfg.UserDataDir)
    }
    if cfg.ProxyURL != "" {
        l = l.Proxy(cfg.ProxyURL)
    }
    if !cfg.WindowSize.Empty() {
        l = l.Window(cfg.WindowSize.Dx(), cfg.WindowSize.Dy())
    }
    if cfg.UserAgent != "" {
        l = l.UserAgent(cfg.UserAgent)
    }

    url, _ := l-launch()
    browser := rod.New().ControlURL(url)

    if err := browser.Connect(); err != nil {
        return nil, fmt.Errorf("connect: %w", err)
    }

    rb := &RodBrowser{
        browser:  browser,
        stealth:  cfg.StealthMode,
        launched: true,
    }

    // Stealth mode: remove traces de automação
    if cfg.StealthMode {
        rb.setupStealth()
    }

    return rb, nil
}
```

### 2. Page Operations

```go
// internal/browser/rod.go (continuação)

// Navigate abre URL
func (b *RodBrowser) Navigate(ctx context.Context, url string) error {
    page, err := b.browser.PageFromExistingBrowser(ctx)
    if err != nil {
        // Nova página se não existir
        page = b.browser.MustPage()
    }
    b.page = page

    if err := page.Navigate(url); err != nil {
        return fmt.Errorf("navigate: %w", err)
    }

    // Espera página carregar
    page.WaitLoad()
    return nil
}

// Click clica em seletor CSS
func (b *RodBrowser) Click(ctx context.Context, selector string) error {
    el, err := b.page.Timeout(10*time.Second).Element(selector)
    if err != nil {
        return fmt.Errorf("element not found: %s: %w", selector, err)
    }

    if err := el.ClickInput(ctx); err != nil {
        return fmt.Errorf("click: %w", err)
    }
    return nil
}

// Type digita texto em elemento
func (b *RodBrowser) Type(ctx context.Context, selector, text string) error {
    el, err := b.page.Timeout(10*time.Second).Element(selector)
    if err != nil {
        return fmt.Errorf("element not found: %s: %w", selector, err)
    }

    // Clear antes de digitar
    el.SelectAllText()
    if err := el.Type(text); err != nil {
        return fmt.Errorf("type: %w", err)
    }
    return nil
}

// Scroll rola a página
func (b *RodBrowser) Scroll(ctx context.Context, direction string, amount float64) error {
    switch direction {
    case "up":
        return b.page.Mouse.Scroll(0, -amount)
    case "down":
        return b.page.Mouse.Scroll(0, amount)
    case "left":
        return b.page.Mouse.Scroll(-amount, 0)
    case "right":
        return b.page.Mouse.Scroll(amount, 0)
    default:
        return fmt.Errorf("invalid direction: %s", direction)
    }
}
```

### 3. Screenshot & State Capture

```go
// internal/browser/rod.go (continuação)

// Screenshot captura screenshot completo
func (b *RodBrowser) Screenshot(ctx context.Context) ([]byte, error) {
    if b.page == nil {
        return nil, fmt.Errorf("no active page")
    }

    img, err := b.page.Screenshot(true, &proto.PageCaptureScreenshot{
        Format:  proto.PageCaptureScreenshotFormatPng,
        Quality: 90,
    })
    if err != nil {
        return nil, fmt.Errorf("screenshot: %w", err)
    }

    return img, nil
}

// ScreenshotElement captura screenshot de elemento específico
func (b *RodBrowser) ScreenshotElement(ctx context.Context, selector string) ([]byte, error) {
    el, err := b.page.Timeout(5*time.Second).Element(selector)
    if err != nil {
        return nil, fmt.Errorf("element: %w", err)
    }

    return el.Screenshot()
}

// GetPageState captura estado completo da página
func (b *RodBrowser) GetPageState(ctx context.Context) (*PageState, error) {
    if b.page == nil {
        return nil, fmt.Errorf("no active page")
    }

    // Captura elementos principais
    elements := []map[string]string{}
    b.page.MustElements("a, button, input, select, textarea").ForEach(func(el *rod.Element) {
        tag := el.MustResourceHTML()
        if len(tag) < 200 {
            elements = append(elements, map[string]string{
                "tag":  tag,
                "text": el.MustText(),
            })
        }
    })

    // DOM simplificado
    dom, _ := b.page.GetMainFrame().Evaluate("() => document.body.innerText")

    // Screenshot
    screenshot, _ := b.Screenshot(ctx)

    // Hash para detecção de mudanças
    hash := sha256.Sum256(screenshot)

    return &PageState{
        URL:       b.page.MustURL(),
        Title:     b.page.MustTitle(),
        DOM:       dom.String(),
        Elements:  elements,
        Screenshot: screenshot,
        Hash:      hex.EncodeToString(hash[:]),
    }, nil
}
```

### 4. Stealth Mode (Anti-Bot)

```go
// internal/browser/stealth.go
package browser

// setupStealth remove traces de automação
func (b *RodBrowser) setupStealth() {
    // Remove navigator.webdriver
    b.page.AddScriptTag(`() => {
        Object.defineProperty(navigator, 'webdriver', { get: () => false });

        // Remove Chrome runtime
        window.chrome = { runtime: {} };

        // Mock plugins
        Object.defineProperty(navigator, 'plugins', {
            get: () => [
                { name: 'Chrome PDF Plugin' },
                { name: 'Chrome PDF Viewer' },
                { name: 'Native Client' }
            ]
        });

        // Mock languages
        Object.defineProperty(navigator, 'languages', {
            get: () => ['pt-BR', 'pt', 'en-US', 'en']
        });

        // Remove automation extensions
        delete window.cdc_adoQpoasnfa76pfcZLmcfl_Array;
        delete window.cdc_adoQpoasnfa76pfcZLmcfl_Promise;
        delete window.cdc_adoQpoasnfa76pfcZLmcfl_Symbol;
    }`)
}
```

### 5. Session Persistence

```go
// internal/browser/session.go
package browser

type Session struct {
    ID          string
    Browser     *RodBrowser
    UserDataDir string
    CreatedAt   time.Time
    LastActive  time.Time
}

// NewSession cria nova sessão com user data dir persistente
func NewSession(ctx context.Context, cfg BrowserConfig) (*Session, error) {
    userDataDir := fmt.Sprintf("/tmp/aurelia-browser-%s", uuid.New().String())

    cfg.UserDataDir = userDataDir
    browser, err := NewRodBrowser(ctx, cfg)
    if err != nil {
        return nil, err
    }

    return &Session{
        ID:          uuid.New().String(),
        Browser:     browser,
        UserDataDir: userDataDir,
        CreatedAt:   time.Now(),
        LastActive:  time.Now(),
    }, nil
}

// Close fecha sessão e limpa user data
func (s *Session) Close() error {
    if s.Browser != nil {
        s.Browser.Close()
    }
    os.RemoveAll(s.UserDataDir)
    return nil
}
```

## Consequências

### Positivas
- Puro Go: sem dependência Node.js
- go-rod é ativamente mantido (2024-2025)
- CDP é bem documentado e estável
- Stealth mode ajuda com anti-bot
- Mais leve que Playwright (~50MB vs ~200MB)

### Negativas
- Apenas Chrome/Chromium (não Firefox/Safari)
- CDP pode ter lag em algumas operações
- Menos abstractions que Playwright

### Trade-offs
- go-rod vs rod: go-rod é o pacote principal
- Stealth vs full automation: stealth pode quebrar algumas sites

## Dependências
- ⚠️ `go-rod/rod` - precisa adicionar ao go.mod
- ⚠️ `go-rod/lib/launcher` - para iniciar browser
- ⚠️ `go-rod/lib/proto` - para tipos CDP
- ❌ `internal/browser/` - NÃO EXISTE

## Referências
- [go-rod - github.com/go-rod/rod](https://github.com/go-rod/rod)
- [BUA - github.com/anxuanzi/bua](https://github.com/anxuanzi/bua)
- [ADR-20260328-bua-browser-use-agent-go.md](./20260328-bua-browser-use-agent-go.md)
- [ADR-20260328-container-steel-browser-isolation.md](./20260328-container-steel-browser-isolation.md)

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
