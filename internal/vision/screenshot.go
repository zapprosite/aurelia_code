package vision

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
)

// ScreenshotCapture handles screen capture via browser.
type ScreenshotCapture struct {
	browser *rod.Browser
}

// NewScreenshotCapture creates a new capturer.
func NewScreenshotCapture(b *rod.Browser) *ScreenshotCapture {
	return &ScreenshotCapture{browser: b}
}

// CapturePage takes a screenshot of a specific page.
func (s *ScreenshotCapture) CapturePage(page *rod.Page, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output dir: %w", err)
	}

	filename := fmt.Sprintf("screenshot_%d.png", time.Now().Unix())
	path := filepath.Join(outputDir, filename)

	img, err := page.Screenshot(true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to take screenshot: %w", err)
	}

	if err := os.WriteFile(path, img, 0644); err != nil {
		return "", fmt.Errorf("failed to write screenshot file: %w", err)
	}

	return path, nil
}

// CaptureVisibleArea captures the current viewport.
func (s *ScreenshotCapture) CaptureVisibleArea(page *rod.Page) ([]byte, error) {
	return page.Screenshot(false, nil)
}
