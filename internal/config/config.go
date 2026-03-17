package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kocar/aurelia/internal/runtime"
)

const (
	defaultMaxIterations    = 500
	defaultMemoryWindowSize = 20
	defaultLLMProvider      = "kimi"
	defaultSTTProvider      = "groq"
)

// AppConfig holds all runtime configuration needed for the application.
type AppConfig struct {
	LLMProvider            string
	STTProvider            string
	TelegramBotToken       string
	TelegramAllowedUserIDs []int64
	KimiAPIKey             string
	GroqAPIKey             string
	MaxIterations          int
	DBPath                 string
	MemoryWindowSize       int
	MCPConfigPath          string
}

type fileConfig struct {
	LLMProvider            string  `json:"llm_provider"`
	STTProvider            string  `json:"stt_provider"`
	TelegramBotToken       string  `json:"telegram_bot_token"`
	TelegramAllowedUserIDs []int64 `json:"telegram_allowed_user_ids"`
	KimiAPIKey             string  `json:"kimi_api_key"`
	GroqAPIKey             string  `json:"groq_api_key"`
	MaxIterations          int     `json:"max_iterations"`
	DBPath                 string  `json:"db_path"`
	MemoryWindowSize       int     `json:"memory_window_size"`
	MCPConfigPath          string  `json:"mcp_servers_config_path"`
}

// EditableConfig represents the user-editable portion of the runtime config.
type EditableConfig struct {
	LLMProvider            string
	STTProvider            string
	TelegramBotToken       string
	TelegramAllowedUserIDs []int64
	KimiAPIKey             string
	GroqAPIKey             string
	MaxIterations          int
	MemoryWindowSize       int
}

// Load reads the instance-local JSON config, creates it with defaults when
// missing, and returns the normalized runtime config.
func Load(r *runtime.PathResolver) (*AppConfig, error) {
	path := r.AppConfig()
	defaults := DefaultFileConfig(r)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := writeConfigFile(path, defaults); err != nil {
				return nil, err
			}
			return toAppConfig(defaults), nil
		}
		return nil, fmt.Errorf("stat app config: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read app config: %w", err)
	}

	cfg := defaults
	if len(data) != 0 {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("decode app config: %w", err)
		}
	}

	normalized := normalizeFileConfig(cfg, r)
	if !sameFileConfig(normalized, cfg) {
		if err := writeConfigFile(path, normalized); err != nil {
			return nil, err
		}
	}

	return toAppConfig(normalized), nil
}

func defaultFileConfig(r *runtime.PathResolver) fileConfig {
	return fileConfig{
		LLMProvider:            defaultLLMProvider,
		STTProvider:            defaultSTTProvider,
		TelegramAllowedUserIDs: []int64{},
		MaxIterations:          defaultMaxIterations,
		DBPath:                 filepath.Join(r.Data(), "aurelia.db"),
		MemoryWindowSize:       defaultMemoryWindowSize,
		MCPConfigPath:          filepath.Join(r.Config(), "mcp_servers.json"),
	}
}

// DefaultEditableConfig returns the default user-editable configuration values.
func DefaultEditableConfig() EditableConfig {
	return EditableConfig{
		LLMProvider:            defaultLLMProvider,
		STTProvider:            defaultSTTProvider,
		TelegramAllowedUserIDs: []int64{},
		MaxIterations:          defaultMaxIterations,
		MemoryWindowSize:       defaultMemoryWindowSize,
	}
}

// DefaultFileConfig returns the full default config including instance paths.
func DefaultFileConfig(r *runtime.PathResolver) fileConfig {
	return defaultFileConfig(r)
}

// LoadEditable returns the editable config subset from the current app config.
func LoadEditable(r *runtime.PathResolver) (*EditableConfig, error) {
	cfg, err := Load(r)
	if err != nil {
		return nil, err
	}
	return &EditableConfig{
		LLMProvider:            cfg.LLMProvider,
		STTProvider:            cfg.STTProvider,
		TelegramBotToken:       cfg.TelegramBotToken,
		TelegramAllowedUserIDs: append([]int64(nil), cfg.TelegramAllowedUserIDs...),
		KimiAPIKey:             cfg.KimiAPIKey,
		GroqAPIKey:             cfg.GroqAPIKey,
		MaxIterations:          cfg.MaxIterations,
		MemoryWindowSize:       cfg.MemoryWindowSize,
	}, nil
}

// SaveEditable updates the user-editable config subset while preserving managed paths.
func SaveEditable(r *runtime.PathResolver, editable EditableConfig) error {
	cfg := normalizeFileConfig(fileConfig{
		LLMProvider:            editable.LLMProvider,
		STTProvider:            editable.STTProvider,
		TelegramBotToken:       editable.TelegramBotToken,
		TelegramAllowedUserIDs: append([]int64(nil), editable.TelegramAllowedUserIDs...),
		KimiAPIKey:             editable.KimiAPIKey,
		GroqAPIKey:             editable.GroqAPIKey,
		MaxIterations:          editable.MaxIterations,
		MemoryWindowSize:       editable.MemoryWindowSize,
	}, r)
	return writeConfigFile(r.AppConfig(), cfg)
}

func normalizeFileConfig(cfg fileConfig, r *runtime.PathResolver) fileConfig {
	defaults := defaultFileConfig(r)
	if cfg.TelegramAllowedUserIDs == nil {
		cfg.TelegramAllowedUserIDs = defaults.TelegramAllowedUserIDs
	}
	if cfg.LLMProvider == "" {
		cfg.LLMProvider = defaults.LLMProvider
	}
	if cfg.STTProvider == "" {
		cfg.STTProvider = defaults.STTProvider
	}
	if cfg.MaxIterations <= 0 {
		cfg.MaxIterations = defaults.MaxIterations
	}
	if cfg.DBPath == "" {
		cfg.DBPath = defaults.DBPath
	}
	if cfg.MemoryWindowSize <= 0 {
		cfg.MemoryWindowSize = defaults.MemoryWindowSize
	}
	if cfg.MCPConfigPath == "" {
		cfg.MCPConfigPath = defaults.MCPConfigPath
	}
	return cfg
}

func writeConfigFile(path string, cfg fileConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create app config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode app config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write app config: %w", err)
	}
	return nil
}

func toAppConfig(cfg fileConfig) *AppConfig {
	return &AppConfig{
		LLMProvider:            cfg.LLMProvider,
		STTProvider:            cfg.STTProvider,
		TelegramBotToken:       cfg.TelegramBotToken,
		TelegramAllowedUserIDs: cfg.TelegramAllowedUserIDs,
		KimiAPIKey:             cfg.KimiAPIKey,
		GroqAPIKey:             cfg.GroqAPIKey,
		MaxIterations:          cfg.MaxIterations,
		DBPath:                 cfg.DBPath,
		MemoryWindowSize:       cfg.MemoryWindowSize,
		MCPConfigPath:          cfg.MCPConfigPath,
	}
}

func sameFileConfig(a, b fileConfig) bool {
	if a.TelegramBotToken != b.TelegramBotToken ||
		a.LLMProvider != b.LLMProvider ||
		a.STTProvider != b.STTProvider ||
		a.KimiAPIKey != b.KimiAPIKey ||
		a.GroqAPIKey != b.GroqAPIKey ||
		a.MaxIterations != b.MaxIterations ||
		a.DBPath != b.DBPath ||
		a.MemoryWindowSize != b.MemoryWindowSize ||
		a.MCPConfigPath != b.MCPConfigPath {
		return false
	}
	if len(a.TelegramAllowedUserIDs) != len(b.TelegramAllowedUserIDs) {
		return false
	}
	for i := range a.TelegramAllowedUserIDs {
		if a.TelegramAllowedUserIDs[i] != b.TelegramAllowedUserIDs[i] {
			return false
		}
	}
	return true
}
