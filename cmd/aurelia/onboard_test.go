package main

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
)

func TestRunOnboard_SavesInteractiveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	input := strings.Join([]string{
		"kimi-key",
		"groq-key",
		"telegram-token",
		"101,202",
		"700",
		"33",
		"",
	}, "\n")

	var out bytes.Buffer
	if err := runOnboard(strings.NewReader(input), &out); err != nil {
		t.Fatalf("runOnboard() error = %v", err)
	}

	resolver, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() error = %v", err)
	}
	cfg, err := config.Load(resolver)
	if err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}

	if cfg.LLMProvider != "kimi" {
		t.Fatalf("LLMProvider = %q", cfg.LLMProvider)
	}
	if cfg.STTProvider != "groq" {
		t.Fatalf("STTProvider = %q", cfg.STTProvider)
	}
	if cfg.TelegramBotToken != "telegram-token" {
		t.Fatalf("TelegramBotToken = %q", cfg.TelegramBotToken)
	}
	if got, want := cfg.TelegramAllowedUserIDs, []int64{101, 202}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("TelegramAllowedUserIDs = %v", got)
	}
	if cfg.KimiAPIKey != "kimi-key" {
		t.Fatalf("KimiAPIKey = %q", cfg.KimiAPIKey)
	}
	if cfg.GroqAPIKey != "groq-key" {
		t.Fatalf("GroqAPIKey = %q", cfg.GroqAPIKey)
	}
	if cfg.MaxIterations != 700 {
		t.Fatalf("MaxIterations = %d", cfg.MaxIterations)
	}
	if cfg.MemoryWindowSize != 33 {
		t.Fatalf("MemoryWindowSize = %d", cfg.MemoryWindowSize)
	}
	if cfg.DBPath != filepath.Join(tmpDir, "data", "aurelia.db") {
		t.Fatalf("DBPath = %q", cfg.DBPath)
	}
}

func TestRunOnboard_PreservesExistingValuesOnBlankInput(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	resolver, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() error = %v", err)
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		t.Fatalf("runtime.Bootstrap() error = %v", err)
	}
	if err := config.SaveEditable(resolver, config.EditableConfig{
		LLMProvider:            "kimi",
		STTProvider:            "groq",
		TelegramBotToken:       "old-telegram",
		TelegramAllowedUserIDs: []int64{42},
		KimiAPIKey:             "old-kimi",
		GroqAPIKey:             "old-groq",
		MaxIterations:          600,
		MemoryWindowSize:       21,
	}); err != nil {
		t.Fatalf("config.SaveEditable() error = %v", err)
	}

	var out bytes.Buffer
	if err := runOnboard(strings.NewReader("\n\n\n\n\n\n"), &out); err != nil {
		t.Fatalf("runOnboard() error = %v", err)
	}

	cfg, err := config.Load(resolver)
	if err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}
	if cfg.TelegramBotToken != "old-telegram" || cfg.KimiAPIKey != "old-kimi" || cfg.GroqAPIKey != "old-groq" {
		t.Fatalf("expected secrets to be preserved, got %+v", cfg)
	}
	if cfg.LLMProvider != "kimi" || cfg.STTProvider != "groq" {
		t.Fatalf("expected providers to be preserved, got llm=%q stt=%q", cfg.LLMProvider, cfg.STTProvider)
	}
	if len(cfg.TelegramAllowedUserIDs) != 1 || cfg.TelegramAllowedUserIDs[0] != 42 {
		t.Fatalf("expected allowed user IDs to be preserved, got %v", cfg.TelegramAllowedUserIDs)
	}
	if cfg.MaxIterations != 600 || cfg.MemoryWindowSize != 21 {
		t.Fatalf("expected numeric fields to be preserved, got max=%d memory=%d", cfg.MaxIterations, cfg.MemoryWindowSize)
	}
}

func TestParseInt64List_RejectsInvalidInput(t *testing.T) {
	if _, err := parseInt64List("123,abc"); err == nil {
		t.Fatal("expected parseInt64List() to fail on invalid input")
	}
}

func TestRenderOnboardingHeader_IncludesBannerAndColor(t *testing.T) {
	header := renderOnboardingHeader()

	if !strings.Contains(header, "$$$$$$\\") {
		t.Fatal("expected ASCII banner in onboarding header")
	}
	if !strings.Contains(header, colorBlue) || !strings.Contains(header, colorReset) {
		t.Fatal("expected ANSI blue color codes in onboarding header")
	}
	if !strings.Contains(header, "Local onboarding for runtime config") {
		t.Fatal("expected onboarding subtitle in header")
	}
}

func TestOnboardingUI_MenuFlowAndBack(t *testing.T) {
	ui := newOnboardingUI(config.DefaultEditableConfig())

	_, _, err := ui.HandleKey(keyEvent{code: keyEnter})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.step != stepLLMKey {
		t.Fatalf("step = %v, want %v", ui.step, stepLLMKey)
	}
	ui.input = "kimi-key"
	_, _, err = ui.HandleKey(keyEvent{code: keyEnter})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.step != stepSTTProvider {
		t.Fatalf("step = %v, want %v", ui.step, stepSTTProvider)
	}
	_, _, err = ui.HandleKey(keyEvent{code: keyLeft})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.step != stepLLMKey {
		t.Fatalf("step = %v, want %v after back", ui.step, stepLLMKey)
	}
}

func TestReadKey_TreatsQAsInputRune(t *testing.T) {
	ev, err := readKey(bufio.NewReader(strings.NewReader("q")))
	if err != nil {
		t.Fatalf("readKey() error = %v", err)
	}
	if ev.code != keyRune || ev.r != 'q' {
		t.Fatalf("expected q to be treated as input rune, got %+v", ev)
	}
}
