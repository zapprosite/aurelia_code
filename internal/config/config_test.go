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
	if cfg.STTProvider != defaultSTTProvider {
		t.Fatalf("STTProvider = %q, want %q", cfg.STTProvider, defaultSTTProvider)
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

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	input := fileConfig{
		LLMProvider:            "kimi",
		STTProvider:            "groq",
		TelegramBotToken:       "telegram-token",
		TelegramAllowedUserIDs: []int64{1, 2, 3},
		KimiAPIKey:             "kimi-key",
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
	if cfg.LLMProvider != input.LLMProvider {
		t.Fatalf("LLMProvider = %q, want %q", cfg.LLMProvider, input.LLMProvider)
	}
	if cfg.STTProvider != input.STTProvider {
		t.Fatalf("STTProvider = %q, want %q", cfg.STTProvider, input.STTProvider)
	}
	if !reflect.DeepEqual(cfg.TelegramAllowedUserIDs, input.TelegramAllowedUserIDs) {
		t.Fatalf("TelegramAllowedUserIDs = %v, want %v", cfg.TelegramAllowedUserIDs, input.TelegramAllowedUserIDs)
	}
	if cfg.KimiAPIKey != input.KimiAPIKey {
		t.Fatalf("KimiAPIKey = %q, want %q", cfg.KimiAPIKey, input.KimiAPIKey)
	}
	if cfg.GroqAPIKey != input.GroqAPIKey {
		t.Fatalf("GroqAPIKey = %q, want %q", cfg.GroqAPIKey, input.GroqAPIKey)
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
	if cfg.STTProvider != defaultSTTProvider {
		t.Fatalf("STTProvider = %q, want %q", cfg.STTProvider, defaultSTTProvider)
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

func TestSaveEditable_PreservesManagedPaths(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	if err := SaveEditable(r, EditableConfig{
		LLMProvider:            "kimi",
		STTProvider:            "groq",
		TelegramBotToken:       "telegram-token",
		TelegramAllowedUserIDs: []int64{7, 8},
		KimiAPIKey:             "kimi-key",
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
	if cfg.LLMProvider != "kimi" || cfg.STTProvider != "groq" {
		t.Fatalf("unexpected providers llm=%q stt=%q", cfg.LLMProvider, cfg.STTProvider)
	}
	if cfg.MCPConfigPath != filepath.Join(tmpDir, "config", "mcp_servers.json") {
		t.Fatalf("MCPConfigPath = %q, want managed default", cfg.MCPConfigPath)
	}
	if !reflect.DeepEqual(cfg.TelegramAllowedUserIDs, []int64{7, 8}) {
		t.Fatalf("TelegramAllowedUserIDs = %v", cfg.TelegramAllowedUserIDs)
	}
}
