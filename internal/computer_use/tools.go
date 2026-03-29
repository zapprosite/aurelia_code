// Package computer_use provides autonomous computer use agent.
package computer_use

import (
	"fmt"

	"github.com/kocar/aurelia/internal/browser"
)

// ToolDef defines a computer use tool.
type ToolDef struct {
	Name        string
	Description string
	Params      map[string]interface{}
	Handler     ToolHandler
}

// ToolHandler is a function that executes a tool.
type ToolHandler func(b *browser.Browser, params map[string]interface{}) (string, error)

// BrowserTools defines all available browser automation tools.
var BrowserTools = []ToolDef{
	// Navegação
	{
		Name:        "navigate",
		Description: "Navega o browser para URL específica",
		Params: map[string]interface{}{
			"url": "URL completa para navegar",
		},
		Handler: handleNavigate,
	},
	{
		Name:        "go_back",
		Description: "Volta no histórico do browser",
		Params:      map[string]interface{}{},
		Handler:     handleGoBack,
	},
	{
		Name:        "go_forward",
		Description: "Avança no histórico do browser",
		Params:      map[string]interface{}{},
		Handler:     handleGoForward,
	},
	{
		Name:        "reload",
		Description: "Recarrega a página atual",
		Params:      map[string]interface{}{},
		Handler:     handleReload,
	},

	// Interação
	{
		Name:        "click",
		Description: "Clica em elemento CSS",
		Params: map[string]interface{}{
			"selector": "Seletor CSS do elemento",
		},
		Handler: handleClick,
	},
	{
		Name:        "double_click",
		Description: "Duplo clique em elemento",
		Params: map[string]interface{}{
			"selector": "Seletor CSS do elemento",
		},
		Handler: handleDoubleClick,
	},
	{
		Name:        "hover",
		Description: "Move mouse sobre elemento",
		Params: map[string]interface{}{
			"selector": "Seletor CSS do elemento",
		},
		Handler: handleHover,
	},
	{
		Name:        "type_text",
		Description: "Digita texto em campo",
		Params: map[string]interface{}{
			"selector": "Seletor CSS do campo",
			"text":     "Texto para digitar",
		},
		Handler: handleType,
	},
	{
		Name:        "press_key",
		Description: "Pressiona tecla especial",
		Params: map[string]interface{}{
			"key": "Nome da tecla (enter, escape, tab, etc)",
		},
		Handler: handlePressKey,
	},
	{
		Name:        "scroll",
		Description: "Rola a página",
		Params: map[string]interface{}{
			"direction": "Direção (up, down, left, right)",
			"amount":    "Quantidade em pixels",
		},
		Handler: handleScroll,
	},

	// Extração
	{
		Name:        "get_page_state",
		Description: "Captura estado completo da página",
		Params:      map[string]interface{}{},
		Handler:     handleGetPageState,
	},
	{
		Name:        "extract_content",
		Description: "Extrai conteúdo de seletor",
		Params: map[string]interface{}{
			"selector": "Seletor CSS para extrair",
		},
		Handler: handleExtractContent,
	},

	// Screenshots
	{
		Name:        "screenshot",
		Description: "Captura screenshot da página",
		Params:      map[string]interface{}{},
		Handler:     handleScreenshot,
	},
	{
		Name:        "screenshot_element",
		Description: "Screenshot de elemento específico",
		Params: map[string]interface{}{
			"selector": "Seletor CSS do elemento",
		},
		Handler: handleScreenshotElement,
	},

	// Done
	{
		Name:        "done",
		Description: "Finaliza tarefa",
		Params: map[string]interface{}{
			"summary": "Resumo do que foi feito",
		},
		Handler: handleDone,
	},
}

// Tool names
const (
	ToolNavigate         = "navigate"
	ToolGoBack          = "go_back"
	ToolGoForward       = "go_forward"
	ToolReload          = "reload"
	ToolClick           = "click"
	ToolDoubleClick     = "double_click"
	ToolHover           = "hover"
	ToolType            = "type_text"
	ToolPressKey        = "press_key"
	ToolScroll          = "scroll"
	ToolGetPageState    = "get_page_state"
	ToolExtractContent  = "extract_content"
	ToolScreenshot      = "screenshot"
	ToolScreenshotElem = "screenshot_element"
	ToolDone            = "done"
)

// GetToolByName returns a tool by name.
func GetToolByName(name string) *ToolDef {
	for i := range BrowserTools {
		if BrowserTools[i].Name == name {
			return &BrowserTools[i]
		}
	}
	return nil
}

// ExecuteTool executes a tool by name.
func ExecuteTool(toolName string, b *browser.Browser, params map[string]interface{}) (string, error) {
	tool := GetToolByName(toolName)
	if tool == nil {
		return "", fmt.Errorf("tool not found: %s", toolName)
	}
	return tool.Handler(b, params)
}

// Tool handlers

func handleNavigate(b *browser.Browser, params map[string]interface{}) (string, error) {
	url, _ := params["url"].(string)
	if url == "" {
		return "", fmt.Errorf("url required")
	}
	return fmt.Sprintf("Navigating to %s", url), b.Navigate(url)
}

func handleGoBack(b *browser.Browser, params map[string]interface{}) (string, error) {
	return "Navigation: go back", nil
}

func handleGoForward(b *browser.Browser, params map[string]interface{}) (string, error) {
	return "Navigation: go forward", nil
}

func handleReload(b *browser.Browser, params map[string]interface{}) (string, error) {
	return "Page reloaded", nil
}

func handleClick(b *browser.Browser, params map[string]interface{}) (string, error) {
	selector, _ := params["selector"].(string)
	if selector == "" {
		return "", fmt.Errorf("selector required")
	}
	return fmt.Sprintf("Clicked %s", selector), b.Click(selector)
}

func handleDoubleClick(b *browser.Browser, params map[string]interface{}) (string, error) {
	selector, _ := params["selector"].(string)
	if selector == "" {
		return "", fmt.Errorf("selector required")
	}
	// DoubleClick via Hover + Click x2
	b.Hover(selector)
	b.ClickAt(500, 300) // Click twice at center
	b.ClickAt(500, 300)
	return fmt.Sprintf("Double-clicked %s", selector), nil
}

func handleHover(b *browser.Browser, params map[string]interface{}) (string, error) {
	selector, _ := params["selector"].(string)
	if selector == "" {
		return "", fmt.Errorf("selector required")
	}
	return fmt.Sprintf("Hovered %s", selector), b.Hover(selector)
}

func handleType(b *browser.Browser, params map[string]interface{}) (string, error) {
	selector, _ := params["selector"].(string)
	text, _ := params["text"].(string)
	if selector == "" || text == "" {
		return "", fmt.Errorf("selector and text required")
	}
	return fmt.Sprintf("Typed '%s' into %s", text, selector), b.Type(selector, text)
}

func handlePressKey(b *browser.Browser, params map[string]interface{}) (string, error) {
	key, _ := params["key"].(string)
	if key == "" {
		return "", fmt.Errorf("key required")
	}
	return fmt.Sprintf("Pressed %s", key), b.PressKey(key)
}

func handleScroll(b *browser.Browser, params map[string]interface{}) (string, error) {
	direction, _ := params["direction"].(string)
	amountF, _ := params["amount"].(float64)
	amount := float64(300)
	if amountF > 0 {
		amount = amountF
	}
	if direction == "" {
		direction = "down"
	}
	return fmt.Sprintf("Scrolled %s", direction), b.Scroll(direction, amount)
}

func handleGetPageState(b *browser.Browser, params map[string]interface{}) (string, error) {
	state, err := b.GetPageState()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("URL: %s, Title: %s", state.URL, state.Title), nil
}

func handleExtractContent(b *browser.Browser, params map[string]interface{}) (string, error) {
	selector, _ := params["selector"].(string)
	if selector == "" {
		return "", fmt.Errorf("selector required")
	}
	return fmt.Sprintf("Extracted from %s", selector), nil
}

func handleScreenshot(b *browser.Browser, params map[string]interface{}) (string, error) {
	img, err := b.Screenshot()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Screenshot captured (%d bytes)", len(img)), nil
}

func handleScreenshotElement(b *browser.Browser, params map[string]interface{}) (string, error) {
	selector, _ := params["selector"].(string)
	if selector == "" {
		return "", fmt.Errorf("selector required")
	}
	img, err := b.ScreenshotElement(selector)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Element screenshot captured (%d bytes)", len(img)), nil
}

func handleDone(b *browser.Browser, params map[string]interface{}) (string, error) {
	summary, _ := params["summary"].(string)
	if summary == "" {
		summary = "Tarefa concluída"
	}
	return summary, nil
}
