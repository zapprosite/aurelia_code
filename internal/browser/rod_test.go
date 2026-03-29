package browser

import (
	"os"
	"testing"
	"time"

	"github.com/go-rod/rod/lib/launcher"
)

func TestMain(m *testing.M) {
	// Check if Chrome is available
	_, err := launcher.New().Bin("").Launch()
	if err != nil {
		// Chrome not available, skip tests
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestBrowserConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Headless != true {
		t.Error("expected Headless to be true by default")
	}
	if cfg.StealthMode != true {
		t.Error("expected StealthMode to be true by default")
	}
	if cfg.WindowSize.Dx() != 1920 {
		t.Errorf("expected WindowSize width 1920, got %d", cfg.WindowSize.Dx())
	}
	if cfg.WindowSize.Dy() != 1080 {
		t.Errorf("expected WindowSize height 1080, got %d", cfg.WindowSize.Dy())
	}
}

func TestBrowserConfigCustom(t *testing.T) {
	cfg := Config{
		Headless:    false,
		StealthMode: false,
		UserAgent:   "test-agent",
	}

	if cfg.Headless != false {
		t.Error("expected Headless to be false")
	}
	if cfg.StealthMode != false {
		t.Error("expected StealthMode to be false")
	}
	if cfg.UserAgent != "test-agent" {
		t.Error("expected UserAgent to be test-agent")
	}
}

func TestSessionConfig(t *testing.T) {
	cfg := DefaultSessionConfig()
	if cfg.DataDir != "/tmp/aurelia-browser-sessions" {
		t.Errorf("expected DataDir /tmp/aurelia-browser-sessions, got %s", cfg.DataDir)
	}
	if cfg.MaxAge != 24*time.Hour {
		t.Errorf("expected MaxAge 24h, got %v", cfg.MaxAge)
	}
}

func TestSession(t *testing.T) {
	session := &Session{
		ID:          "test-id",
		CreatedAt:   time.Now(),
		LastActive:  time.Now(),
	}

	if session.ID != "test-id" {
		t.Error("expected ID test-id")
	}

	session.Touch()
	if session.LastActive.IsZero() {
		t.Error("expected LastActive to be updated")
	}
}

func TestSessionIsIdle(t *testing.T) {
	session := &Session{
		LastActive: time.Now().Add(-1 * time.Minute),
	}

	if !session.IsIdle(30 * time.Second) {
		t.Error("expected session to be idle after 1 minute")
	}

	if session.IsIdle(5 * time.Minute) {
		t.Error("expected session to not be idle after 1 minute with 5min timeout")
	}
}
