package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kocar/aurelia/internal/runtime"
)

func TestLoad_CreatesDefaultAppConfigWhenMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)
	t.Setenv("TELEGRAM_BOT_TOKEN", "") // isolate from real env

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.MaxIterations != defaultMaxIterations {
		t.Fatalf("MaxIterations = %d, want %d", cfg.MaxIterations, defaultMaxIterations)
	}
	if cfg.LLMProvider != defaultLLMProvider {
		t.Fatalf("LLMProvider = %q, want %q", cfg.LLMProvider, defaultLLMProvider)
	}
	if cfg.LLMModel != defaultLLMModelForProvider(defaultLLMProvider) {
		t.Fatalf("LLMModel = %q, want %q", cfg.LLMModel, defaultLLMModelForProvider(defaultLLMProvider))
	}
	if cfg.STTProvider != defaultSTTProvider {
		t.Fatalf("STTProvider = %q, want %q", cfg.STTProvider, defaultSTTProvider)
	}
	if cfg.TTSProvider != defaultTTSProvider {
		t.Fatalf("TTSProvider = %q, want %q", cfg.TTSProvider, defaultTTSProvider)
	}
	if cfg.TTSBaseURL != defaultLocalTTSBaseURL {
		t.Fatalf("TTSBaseURL = %q, want %q", cfg.TTSBaseURL, defaultLocalTTSBaseURL)
	}
	if cfg.TTSVoice != defaultLocalTTSVoice {
		t.Fatalf("TTSVoice = %q, want %q", cfg.TTSVoice, defaultLocalTTSVoice)
	}

	if cfg.MemoryWindowSize != defaultMemoryWindowSize {
		t.Fatalf("MemoryWindowSize = %d, want %d", cfg.MemoryWindowSize, defaultMemoryWindowSize)
	}

	if _, err := os.Stat(r.AppConfig()); err != nil {
		t.Fatalf("expected app config file to be created: %v", err)
	}
}

func TestLoad_UsesJSONConfigValues(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)
	// isolate from real env so test uses only the JSON file values
	for _, k := range []string{
		"TELEGRAM_BOT_TOKEN", "ANTHROPIC_API_KEY", "GOOGLE_API_KEY",
		"OPENAI_API_KEY", "GROQ_API_KEY", "OPENROUTER_API_KEY",
	} {
		t.Setenv(k, "")
	}

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	input := fileConfig{
		Bots: []BotConfig{{
			ID:             "controle-db",
			Name:           "CONTROLE DB",
			Token:          "telegram-token-2",
			AllowedUserIDs: []int64{7},
			PersonaID:      "data-governance",
			FocusArea:      "Governanca de dados",
			LLMProvider:    "openrouter",
			LLMModel:       "minimax/minimax-m2.7",
			Enabled:        true,
		}},
		LLMProvider:            "ollama",
		LLMModel:               "qwen3.5",
		STTProvider:            "groq",
		TTSProvider:            "openai_compatible",
		TTSBaseURL:             "http://127.0.0.1:8011",
		TTSModel:               "tts-1-hd",
		TTSVoice:               "aurelia",
		TTSFormat:              "opus",
		TTSSpeed:               1.1,
		TelegramBotToken:       "telegram-token",
		TelegramAllowedUserIDs: []int64{1, 2, 3},
		AnthropicAPIKey:        "anthropic-key",
		GoogleAPIKey:           "google-key",
		OpenRouterAPIKey:       "openrouter-key",
		OpenAIAPIKey:           "openai-key",
		GroqAPIKey:             "groq-key",
		MaxIterations:          321,
		DBPath:                 filepath.Join(tmpDir, "data", "custom.db"),
		MemoryWindowSize:       42,
		MCPConfigPath:          filepath.Join(tmpDir, "config", "custom-mcp.json"),
	}
	if err := os.MkdirAll(filepath.Dir(r.AppConfig()), 0o700); err != nil {
		t.Fatalf("MkdirAll() unexpected error: %v", err)
	}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal() unexpected error: %v", err)
	}
	if err := os.WriteFile(r.AppConfig(), data, 0o600); err != nil {
		t.Fatalf("WriteFile() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.TelegramBotToken != input.TelegramBotToken {
		t.Fatalf("TelegramBotToken = %q, want %q", cfg.TelegramBotToken, input.TelegramBotToken)
	}
	if len(cfg.Bots) != 1 {
		t.Fatalf("Bots len = %d, want 1", len(cfg.Bots))
	}
	if cfg.Bots[0].ID != "controle-db" || cfg.Bots[0].LLMProvider != "openrouter" || cfg.Bots[0].LLMModel != "minimax/minimax-m2.7" {
		t.Fatalf("unexpected bot config loaded: %+v", cfg.Bots[0])
	}
	if cfg.LLMProvider != input.LLMProvider {
		t.Fatalf("LLMProvider = %q, want %q", cfg.LLMProvider, input.LLMProvider)
	}
	if cfg.LLMModel != input.LLMModel {
		t.Fatalf("LLMModel = %q, want %q", cfg.LLMModel, input.LLMModel)
	}
	if cfg.STTProvider != input.STTProvider {
		t.Fatalf("STTProvider = %q, want %q", cfg.STTProvider, input.STTProvider)
	}
	if cfg.TTSProvider != input.TTSProvider {
		t.Fatalf("TTSProvider = %q, want %q", cfg.TTSProvider, input.TTSProvider)
	}
	if cfg.TTSBaseURL != input.TTSBaseURL {
		t.Fatalf("TTSBaseURL = %q, want %q", cfg.TTSBaseURL, input.TTSBaseURL)
	}
	if cfg.TTSVoice != input.TTSVoice {
		t.Fatalf("TTSVoice = %q, want %q", cfg.TTSVoice, input.TTSVoice)
	}
	if !reflect.DeepEqual(cfg.TelegramAllowedUserIDs, input.TelegramAllowedUserIDs) {
		t.Fatalf("TelegramAllowedUserIDs = %v, want %v", cfg.TelegramAllowedUserIDs, input.TelegramAllowedUserIDs)
	}
	if cfg.GroqAPIKey != input.GroqAPIKey {
		t.Fatalf("GroqAPIKey = %q, want %q", cfg.GroqAPIKey, input.GroqAPIKey)
	}
	if cfg.OpenRouterAPIKey != input.OpenRouterAPIKey {
		t.Fatalf("OpenRouterAPIKey = %q, want %q", cfg.OpenRouterAPIKey, input.OpenRouterAPIKey)
	}
	if cfg.OpenAIAPIKey != input.OpenAIAPIKey {
		t.Fatalf("OpenAIAPIKey = %q, want %q", cfg.OpenAIAPIKey, input.OpenAIAPIKey)
	}
	if cfg.AnthropicAPIKey != input.AnthropicAPIKey {
		t.Fatalf("AnthropicAPIKey = %q, want %q", cfg.AnthropicAPIKey, input.AnthropicAPIKey)
	}
	if cfg.GoogleAPIKey != input.GoogleAPIKey {
		t.Fatalf("GoogleAPIKey = %q, want %q", cfg.GoogleAPIKey, input.GoogleAPIKey)
	}
	if cfg.MaxIterations != input.MaxIterations {
		t.Fatalf("MaxIterations = %d, want %d", cfg.MaxIterations, input.MaxIterations)
	}
	if cfg.DBPath != input.DBPath {
		t.Fatalf("DBPath = %q, want %q", cfg.DBPath, input.DBPath)
	}
	if cfg.MemoryWindowSize != input.MemoryWindowSize {
		t.Fatalf("MemoryWindowSize = %d, want %d", cfg.MemoryWindowSize, input.MemoryWindowSize)
	}
	if cfg.MCPConfigPath != input.MCPConfigPath {
		t.Fatalf("MCPConfigPath = %q, want %q", cfg.MCPConfigPath, input.MCPConfigPath)
	}
}

func TestLoad_NormalizesMissingFieldsWithDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)
	t.Setenv("TELEGRAM_BOT_TOKEN", "") // isolate from real env

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(r.AppConfig()), 0o700); err != nil {
		t.Fatalf("MkdirAll() unexpected error: %v", err)
	}
	if err := os.WriteFile(r.AppConfig(), []byte(`{"telegram_bot_token":"abc"}`), 0o600); err != nil {
		t.Fatalf("WriteFile() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.TelegramBotToken != "abc" {
		t.Fatalf("TelegramBotToken = %q, want %q", cfg.TelegramBotToken, "abc")
	}
	if cfg.LLMProvider != defaultLLMProvider {
		t.Fatalf("LLMProvider = %q, want %q", cfg.LLMProvider, defaultLLMProvider)
	}
	if cfg.LLMModel != defaultLLMModelForProvider(defaultLLMProvider) {
		t.Fatalf("LLMModel = %q, want %q", cfg.LLMModel, defaultLLMModelForProvider(defaultLLMProvider))
	}
	if cfg.STTProvider != defaultSTTProvider {
		t.Fatalf("STTProvider = %q, want %q", cfg.STTProvider, defaultSTTProvider)
	}
	if cfg.TTSProvider != defaultTTSProvider {
		t.Fatalf("TTSProvider = %q, want %q", cfg.TTSProvider, defaultTTSProvider)
	}
	if cfg.MaxIterations != defaultMaxIterations {
		t.Fatalf("MaxIterations = %d, want %d", cfg.MaxIterations, defaultMaxIterations)
	}
	if cfg.DBPath != filepath.Join(tmpDir, "data", "aurelia.db") {
		t.Fatalf("DBPath = %q, want instance default", cfg.DBPath)
	}
	if cfg.MemoryWindowSize != defaultMemoryWindowSize {
		t.Fatalf("MemoryWindowSize = %d, want %d", cfg.MemoryWindowSize, defaultMemoryWindowSize)
	}
	if cfg.MCPConfigPath != filepath.Join(tmpDir, "config", "mcp_servers.json") {
		t.Fatalf("MCPConfigPath = %q, want instance default", cfg.MCPConfigPath)
	}
}

func TestLoad_NormalizesLegacyVoiceCaptureCommandToRepoScript(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	repoRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repoRoot, "scripts"), 0o755); err != nil {
		t.Fatalf("MkdirAll() unexpected error: %v", err)
	}
	localScript := filepath.Join(repoRoot, "scripts", "voice-capture-openwakeword.sh")
	if err := os.WriteFile(localScript, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("WriteFile() unexpected error: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() unexpected error: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() unexpected error: %v", err)
	}

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	input := fileConfig{
		VoiceCaptureEnabled: true,
		VoiceCaptureCommand: "/home/will/aurelia-24x7/scripts/voice-capture-openwakeword.sh --output-dir /tmp/out",
	}
	if err := os.MkdirAll(filepath.Dir(r.AppConfig()), 0o700); err != nil {
		t.Fatalf("MkdirAll() unexpected error: %v", err)
	}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal() unexpected error: %v", err)
	}
	if err := os.WriteFile(r.AppConfig(), data, 0o600); err != nil {
		t.Fatalf("WriteFile() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if got := cfg.VoiceCaptureCommand; got != localScript+" --output-dir /tmp/out" {
		t.Fatalf("VoiceCaptureCommand = %q", got)
	}
}

func TestSaveEditable_PreservesManagedPaths(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	if err := SaveEditable(r, EditableConfig{
		LLMProvider:            "ollama",
		LLMModel:               "qwen3.5",
		STTProvider:            "groq",
		TelegramBotToken:       "telegram-token",
		TelegramAllowedUserIDs: []int64{7, 8},
		AnthropicAPIKey:        "anthropic-key",
		GoogleAPIKey:           "google-key",
		OpenRouterAPIKey:       "openrouter-key",
		OpenAIAPIKey:           "openai-key",
		GroqAPIKey:             "groq-key",
		MaxIterations:          900,
		MemoryWindowSize:       25,
	}); err != nil {
		t.Fatalf("SaveEditable() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.DBPath != filepath.Join(tmpDir, "data", "aurelia.db") {
		t.Fatalf("DBPath = %q, want managed default", cfg.DBPath)
	}
	if cfg.LLMProvider != "ollama" || cfg.LLMModel != "gemma3:27b-it-qat" || cfg.STTProvider != "groq" {
		t.Fatalf("unexpected providers llm=%q model=%q stt=%q", cfg.LLMProvider, cfg.LLMModel, cfg.STTProvider)
	}
	if cfg.MCPConfigPath != filepath.Join(tmpDir, "config", "mcp_servers.json") {
		t.Fatalf("MCPConfigPath = %q, want managed default", cfg.MCPConfigPath)
	}
	if !reflect.DeepEqual(cfg.TelegramAllowedUserIDs, []int64{7, 8}) {
		t.Fatalf("TelegramAllowedUserIDs = %v", cfg.TelegramAllowedUserIDs)
	}
}
