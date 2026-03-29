// Package browser provides browser automation using go-rod (CDP).
package browser

import (
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// Browser provides browser automation via Chrome DevTools Protocol.
type Browser struct {
	browser  *rod.Browser
	page     *rod.Page
	stealth  bool
	launched bool
	UserAgent string
}

// Config holds browser configuration.
type Config struct {
	Headless    bool
	StealthMode bool
	UserDataDir string
	ProxyURL    string
	WindowSize  image.Rectangle
	UserAgent   string
	Timeout     time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Headless:    true,
		StealthMode: true,
		WindowSize:  image.Rect(0, 0, 1920, 1080),
		Timeout:     30 * time.Second,
	}
}

// New creates a new browser instance.
func New(cfg Config) (*Browser, error) {
	// Create launcher
	l := launcher.New().
		Headless(cfg.Headless)

	// Set window size
	if !cfg.WindowSize.Empty() {
		l = l.Set("window-size",
			fmt.Sprintf("%d,%d", cfg.WindowSize.Dx(), cfg.WindowSize.Dy()))
	}

	// Set user agent if provided
	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	}
	l = l.Set("user-agent", userAgent)

	// Launch browser
	url, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("launch browser: %w", err)
	}

	// Create browser instance
	b := rod.New().ControlURL(url)

	if err := b.Connect(); err != nil {
		return nil, fmt.Errorf("connect to browser: %w", err)
	}

	browser := &Browser{
		browser:   b,
		stealth:   cfg.StealthMode,
		launched:  true,
		UserAgent: userAgent,
	}

	return browser, nil
}

// Page returns or creates a page.
func (b *Browser) Page() (*rod.Page, error) {
	if b.page != nil {
		return b.page, nil
	}

	page := b.browser.MustPage()
	b.page = page

	// Apply stealth mode
	if b.stealth {
		b.setupStealth(page)
	}

	return page, nil
}

// setupStealth applies stealth mode to hide automation.
func (b *Browser) setupStealth(page *rod.Page) {
	page.AddScriptTag("stealth.js", `
		Object.defineProperty(navigator, 'webdriver', { get: () => false });
		window.chrome = { runtime: {} };
	`)
}

// Navigate opens a URL.
func (b *Browser) Navigate(url string) error {
	page, err := b.Page()
	if err != nil {
		return err
	}

	if err := page.Navigate(url); err != nil {
		return fmt.Errorf("navigate to %s: %w", url, err)
	}

	page.WaitLoad()
	return nil
}

// Click clicks on an element by CSS selector.
func (b *Browser) Click(selector string) error {
	page, err := b.Page()
	if err != nil {
		return err
	}

	el, err := page.Timeout(10 * time.Second).Element(selector)
	if err != nil {
		return fmt.Errorf("element not found: %s: %w", selector, err)
	}

	el.MustClick()
	return nil
}

// ClickAt clicks at specific coordinates.
func (b *Browser) ClickAt(x, y float64) error {
	page, err := b.Page()
	if err != nil {
		return err
	}

	page.Mouse.MustMoveTo(x, y)
	page.Mouse.MustClick(proto.InputMouseButtonLeft)
	return nil
}

// Type types text into an element.
func (b *Browser) Type(selector, text string) error {
	page, err := b.Page()
	if err != nil {
		return err
	}

	el, err := page.Timeout(10 * time.Second).Element(selector)
	if err != nil {
		return fmt.Errorf("element not found: %s: %w", selector, err)
	}

	el.MustFocus()
	el.MustInput(text)
	return nil
}

// TypeKeys types keyboard keys.
func (b *Browser) TypeKeys(keys ...input.Key) error {
	page, err := b.Page()
	if err != nil {
		return err
	}
	page.Keyboard.MustType(keys...)
	return nil
}

// Scroll scrolls the page.
func (b *Browser) Scroll(direction string, amount float64) error {
	page, err := b.Page()
	if err != nil {
		return err
	}

	switch strings.ToLower(direction) {
	case "up":
		page.Mouse.MustScroll(0, -amount)
	case "down":
		page.Mouse.MustScroll(0, amount)
	case "left":
		page.Mouse.MustScroll(-amount, 0)
	case "right":
		page.Mouse.MustScroll(amount, 0)
	default:
		return fmt.Errorf("invalid direction: %s (use up/down/left/right)", direction)
	}
	return nil
}

// Hover moves mouse over an element.
func (b *Browser) Hover(selector string) error {
	page, err := b.Page()
	if err != nil {
		return err
	}

	el, err := page.Timeout(5 * time.Second).Element(selector)
	if err != nil {
		return fmt.Errorf("element not found: %s: %w", selector, err)
	}

	el.MustHover()
	return nil
}

// Screenshot captures a screenshot of the current page.
func (b *Browser) Screenshot() ([]byte, error) {
	page, err := b.Page()
	if err != nil {
		return nil, err
	}

	return page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
	})
}

// ScreenshotElement captures a screenshot of a specific element.
func (b *Browser) ScreenshotElement(selector string) ([]byte, error) {
	page, err := b.Page()
	if err != nil {
		return nil, err
	}

	el, err := page.Timeout(5 * time.Second).Element(selector)
	if err != nil {
		return nil, fmt.Errorf("element not found: %s: %w", selector, err)
	}

	return el.Screenshot(proto.PageCaptureScreenshotFormatPng, 90)
}

// PageState holds the current state of the page.
type PageState struct {
	URL   string
	Title string
	DOM   string
}

// GetPageState captures the full state of the page.
func (b *Browser) GetPageState() (*PageState, error) {
	page, err := b.Page()
	if err != nil {
		return nil, err
	}

	info := page.MustInfo()
	return &PageState{
		URL:   info.URL,
		Title: info.Title,
	}, nil
}

// GetURL returns the current page URL.
func (b *Browser) GetURL() (string, error) {
	page, err := b.Page()
	if err != nil {
		return "", err
	}
	return page.MustInfo().URL, nil
}

// GetTitle returns the current page title.
func (b *Browser) GetTitle() (string, error) {
	page, err := b.Page()
	if err != nil {
		return "", err
	}
	return page.MustInfo().Title, nil
}

// PressKey presses a keyboard key.
func (b *Browser) PressKey(key string) error {
	page, err := b.Page()
	if err != nil {
		return err
	}

	keyMap := map[string]input.Key{
		"enter":     input.Enter,
		"escape":    input.Escape,
		"tab":       input.Tab,
		"backspace": input.Backspace,
		"delete":    input.Delete,
		"ctrl":      input.ControlLeft,
		"alt":       input.AltLeft,
		"shift":     input.ShiftLeft,
		"up":        input.ArrowUp,
		"down":      input.ArrowDown,
		"left":      input.ArrowLeft,
		"right":     input.ArrowRight,
	}

	if k, ok := keyMap[strings.ToLower(key)]; ok {
		page.Keyboard.MustType(k)
	}
	return nil
}

// WaitFor waits for an element to appear.
func (b *Browser) WaitFor(selector string, timeout time.Duration) error {
	page, err := b.Page()
	if err != nil {
		return err
	}
	page.Timeout(timeout).MustElement(selector)
	return nil
}

// Close closes the browser.
func (b *Browser) Close() error {
	if b.page != nil {
		b.page.Close()
	}
	if b.browser != nil {
		b.browser.Close()
	}
	return nil
}

// Health checks if the browser is still running.
func (b *Browser) Health() error {
	if b.browser == nil {
		return errors.New("browser not initialized")
	}
	return nil
}
