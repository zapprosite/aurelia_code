package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kocar/aurelia/internal/runtime"
)

const (
	defaultMaxIterations    = 500
	defaultMemoryWindowSize = 20
	defaultLLMProvider      = "openai_compatible"
	defaultLLMModel         = "aurelia-smart"
	defaultSTTProvider      = "groq" // GPU: faster-whisper large-v3 on RTX 4090 (Fallback)
	defaultSTTLanguage      = "pt"
	defaultGroqSTTBaseURL   = "https://api.groq.com/openai/v1"
	defaultGroqSTTModel     = "whisper-large-v3-turbo"
	defaultTTSProvider      = "openai_compatible"
	defaultLocalTTSBaseURL  = "http://127.0.0.1:8012"           // Kokoro/Kodoro TTS GPU
	defaultLocalTTSModel    = "kokoro"                          // Kokoro-82M ONNX
	defaultLocalTTSVoice    = "pt-BR-ThalitaMultilingualNeural" // S-35: ONLY Thalita, PT-BR native
	defaultTTSLanguage      = "pt"
	defaultLocalTTSFormat   = "opus"

	defaultTTSSpeed = 1.0

	defaultHeartbeatEnabled     = true
	defaultHeartbeatIntervalMin = 30
	defaultVoiceEnabled         = false
	defaultVoicePollIntervalMS  = 1000
	defaultVoiceHeartbeatSec    = 45
	defaultVoiceCapturePollMS   = 1000
	defaultVoiceCaptureFreshSec = 45
	defaultVoiceWakePhrase      = "jarvis"
	defaultGroqSoftCapDaily     = 800
	defaultGroqHardCapDaily     = 1200
	defaultQdrantCollection     = "conversation_memory"
	defaultQdrantEmbeddingModel = "nomic-embed-text"
	defaultDashboardPort        = 3334
	defaultHealthPort           = 8484
	defaultRedisURL             = "localhost:6379"
)

// BotConfig holds per-bot Telegram configuration for multi-bot support.
type BotConfig struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Token          string  `json:"token"`
	AllowedUserIDs []int64 `json:"allowed_user_ids"`
	PersonaID      string  `json:"persona_id"`
	FocusArea      string  `json:"focus_area"`
	LLMProvider    string  `json:"llm_provider,omitempty"`
	LLMModel       string  `json:"llm_model,omitempty"`
	Enabled        bool    `json:"enabled"`
}

// AppConfig holds all runtime configuration needed for the application.
type AppConfig struct {
	Features                  Features // FEAT_VOICE, FEAT_DREAM, FEAT_KAIROS, FEAT_COMPUTER_USE
	Bots                      []BotConfig
	LLMProvider               string
	LLMModel                  string
	STTProvider               string
	STTBaseURL                string
	STTModel                  string
	STTLanguage               string
	TTSProvider               string
	TTSBaseURL                string
	TTSModel                  string
	TTSVoice                  string
	TTSLanguage               string
	TTSFormat                 string
	TTSSpeed                  float64
	PremiumTTSProvider        string
	PremiumTTSBaseURL         string
	PremiumTTSModel           string
	PremiumTTSVoice           string
	TelegramBotToken          string
	TelegramNotificationBotID string

	TelegramAllowedUserIDs   []int64
	AnthropicAPIKey          string
	GoogleAPIKey             string
	OpenRouterAPIKey         string
	OpenAIAPIKey             string
	GroqAPIKey               string
	MiniMaxAPIKey            string
	LiteLLMKey               string
	MaxIterations            int
	DBPath                   string
	MemoryWindowSize         int
	MCPConfigPath            string
	HeartbeatEnabled         bool
	HeartbeatIntervalMinutes int
	VoiceEnabled             bool
	VoiceReplyUserID         int64
	VoiceReplyChatID         int64
	VoiceSpoolPath           string
	VoiceDropPath            string
	VoiceHeartbeatPath       string
	VoiceHeartbeatFreshSec   int
	VoicePollIntervalMS      int
	VoiceWakePhrase          string
	VoiceCaptureEnabled      bool
	VoiceCaptureCommand      string
	VoiceCaptureHeartbeat    string
	VoiceCaptureFreshSec     int
	VoiceCapturePollMS       int
	STTFallbackCommand       string
	GroqSoftCapDaily         int
	GroqHardCapDaily         int
	QdrantURL                string
	QdrantAPIKey             string
	QdrantCollection         string
	QdrantEmbeddingModel     string
	OllamaURL                string
	DashboardPort            int
	HealthPort               int
	ObsidianVaultPath        string
	ObsidianSyncEnabled      bool
	AureliaMode              string // "sovereign" or "lite"
	RedisURL                 string
}

type fileConfig struct {
	Bots []BotConfig `json:"bots,omitempty"`

	LLMProvider               string  `json:"llm_provider"`
	LLMModel                  string  `json:"llm_model"`
	STTProvider               string  `json:"stt_provider"`
	STTBaseURL                string  `json:"stt_base_url"`
	STTModel                  string  `json:"stt_model"`
	STTLanguage               string  `json:"stt_language"`
	TTSProvider               string  `json:"tts_provider"`
	TTSBaseURL                string  `json:"tts_base_url"`
	TTSModel                  string  `json:"tts_model"`
	TTSVoice                  string  `json:"tts_voice"`
	TTSLanguage               string  `json:"tts_language"`
	TTSFormat                 string  `json:"tts_format"`
	TTSSpeed                  float64 `json:"tts_speed"`
	PremiumTTSProvider        string  `json:"premium_tts_provider"`
	PremiumTTSBaseURL         string  `json:"premium_tts_base_url"`
	PremiumTTSModel           string  `json:"premium_tts_model"`
	PremiumTTSVoice           string  `json:"premium_tts_voice"`
	TelegramBotToken          string  `json:"telegram_bot_token"`
	TelegramNotificationBotID string  `json:"telegram_notification_bot_id"`

	TelegramAllowedUserIDs   []int64 `json:"telegram_allowed_user_ids"`
	AnthropicAPIKey          string  `json:"anthropic_api_key"`
	GoogleAPIKey             string  `json:"google_api_key"`
	OpenRouterAPIKey         string  `json:"openrouter_api_key"`
	OpenAIAPIKey             string  `json:"openai_api_key"`
	GroqAPIKey               string  `json:"groq_api_key"`
	MiniMaxAPIKey            string  `json:"minimax_api_key"`
	LiteLLMKey               string  `json:"litellm_key"`
	MaxIterations            int     `json:"max_iterations"`
	DBPath                   string  `json:"db_path"`
	MemoryWindowSize         int     `json:"memory_window_size"`
	MCPConfigPath            string  `json:"mcp_servers_config_path"`
	HeartbeatEnabled         bool    `json:"heartbeat_enabled"`
	HeartbeatIntervalMinutes int     `json:"heartbeat_interval_minutes"`
	VoiceEnabled             bool    `json:"voice_enabled"`
	VoiceReplyUserID         int64   `json:"voice_reply_user_id"`
	VoiceReplyChatID         int64   `json:"voice_reply_chat_id"`
	VoiceSpoolPath           string  `json:"voice_spool_path"`
	VoiceDropPath            string  `json:"voice_drop_path"`
	VoiceHeartbeatPath       string  `json:"voice_heartbeat_path"`
	VoiceHeartbeatFreshSec   int     `json:"voice_heartbeat_fresh_seconds"`
	VoicePollIntervalMS      int     `json:"voice_poll_interval_ms"`
	VoiceWakePhrase          string  `json:"voice_wake_phrase"`
	VoiceCaptureEnabled      bool    `json:"voice_capture_enabled"`
	VoiceCaptureCommand      string  `json:"voice_capture_command"`
	VoiceCaptureHeartbeat    string  `json:"voice_capture_heartbeat_path"`
	VoiceCaptureFreshSec     int     `json:"voice_capture_heartbeat_fresh_seconds"`
	VoiceCapturePollMS       int     `json:"voice_capture_poll_interval_ms"`
	STTFallbackCommand       string  `json:"stt_fallback_command"`
	GroqSoftCapDaily         int     `json:"groq_soft_cap_daily"`
	GroqHardCapDaily         int     `json:"groq_hard_cap_daily"`
	QdrantURL                string  `json:"qdrant_url"`
	QdrantAPIKey             string  `json:"qdrant_api_key"`
	QdrantCollection         string  `json:"qdrant_collection"`
	QdrantEmbeddingModel     string  `json:"qdrant_embedding_model"`
	OllamaURL                string  `json:"ollama_url"`
	DashboardPort            int     `json:"dashboard_port"`
	HealthPort               int     `json:"health_port"`
	ObsidianVaultPath        string  `json:"obsidian_vault_path"`
	ObsidianSyncEnabled      bool    `json:"obsidian_sync_enabled"`
	AureliaMode              string  `json:"aurelia_mode"`
	RedisURL                 string  `json:"redis_url"`
}

// EditableConfig represents the user-editable portion of the runtime config.
type EditableConfig struct {
	LLMProvider               string
	LLMModel                  string
	STTProvider               string
	STTLanguage               string
	TTSProvider               string
	TTSBaseURL                string
	TTSModel                  string
	TTSVoice                  string
	TTSLanguage               string
	TTSFormat                 string
	TTSSpeed                  float64
	PremiumTTSProvider        string
	PremiumTTSBaseURL         string
	PremiumTTSModel           string
	PremiumTTSVoice           string
	TelegramBotToken          string
	TelegramNotificationBotID string

	TelegramAllowedUserIDs   []int64
	AnthropicAPIKey          string
	GoogleAPIKey             string
	OpenRouterAPIKey         string
	OpenAIAPIKey             string
	GroqAPIKey               string
	MiniMaxAPIKey            string
	LiteLLMKey               string
	MaxIterations            int
	MemoryWindowSize         int
	HeartbeatEnabled         bool
	HeartbeatIntervalMinutes int
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

	// Apply environment variable overrides for secrets
	applyEnvOverrides(&normalized)

	if !sameFileConfig(normalized, cfg) {
		if err := writeConfigFile(path, normalized); err != nil {
			return nil, err
		}
	}

	return toAppConfig(normalized), nil
}

func applyEnvOverrides(cfg *fileConfig) {
	if env := os.Getenv("DASHBOARD_PORT"); env != "" {
		if v, err := strconv.Atoi(env); err == nil && v > 0 {
			cfg.DashboardPort = v
		}
	}
	if env := os.Getenv("HEALTH_PORT"); env != "" {
		if v, err := strconv.Atoi(env); err == nil && v > 0 {
			cfg.HealthPort = v
		}
	}
	if env := os.Getenv("TELEGRAM_BOT_TOKEN"); env != "" {
		cfg.TelegramBotToken = env
	}
	if env := os.Getenv("ANTHROPIC_API_KEY"); env != "" {
		cfg.AnthropicAPIKey = env
	}
	if env := os.Getenv("GOOGLE_API_KEY"); env != "" {
		cfg.GoogleAPIKey = env
	}
	if env := os.Getenv("OPENAI_API_KEY"); env != "" {
		cfg.OpenAIAPIKey = env
	}
	if env := os.Getenv("GROQ_API_KEY"); env != "" {
		cfg.GroqAPIKey = env
	}
	if env := os.Getenv("OPENROUTER_API_KEY"); env != "" {
		cfg.OpenRouterAPIKey = env
	}
	if env := os.Getenv("QDRANT_API_KEY"); env != "" {
		cfg.QdrantAPIKey = env
	}
	if env := os.Getenv("OLLAMA_URL"); env != "" {
		cfg.OllamaURL = env
	}
	if env := os.Getenv("MINIMAX_API_KEY"); env != "" {
		cfg.MiniMaxAPIKey = env
	}
	if env := os.Getenv("LITELLM_MASTER_KEY"); env != "" {
		cfg.LiteLLMKey = env
	}
	if env := os.Getenv("OBSIDIAN_VAULT_PATH"); env != "" {
		cfg.ObsidianVaultPath = env
	}
	if env := os.Getenv("AURELIA_MODE"); env != "" {
		cfg.AureliaMode = strings.ToLower(env)
	}
	if env := os.Getenv("REDIS_URL"); env != "" {
		cfg.RedisURL = env
	}
	if env := os.Getenv("STT_PROVIDER"); env != "" {
		cfg.STTProvider = env
	}
	if env := os.Getenv("STT_BASE_URL"); env != "" {
		cfg.STTBaseURL = env
	}
	if env := os.Getenv("STT_MODEL"); env != "" {
		cfg.STTModel = env
	}
	if env := os.Getenv("STT_LANGUAGE"); env != "" {
		cfg.STTLanguage = env
	}
	if env := os.Getenv("TTS_PROVIDER"); env != "" {
		cfg.TTSProvider = env
	}
	if env := os.Getenv("TTS_BASE_URL"); env != "" {
		cfg.TTSBaseURL = env
	}
	if env := os.Getenv("TTS_MODEL"); env != "" {
		cfg.TTSModel = env
	}
}

func defaultFileConfig(r *runtime.PathResolver) fileConfig {
	return fileConfig{
		LLMProvider:            defaultLLMProvider,
		LLMModel:               defaultLLMModelForProvider(defaultLLMProvider),
		STTProvider:            defaultSTTProvider,
		STTBaseURL:             defaultGroqSTTBaseURL,
		STTModel:               defaultGroqSTTModel,
		STTLanguage:            defaultSTTLanguage,
		TTSProvider:            defaultTTSProvider,
		TTSBaseURL:             defaultLocalTTSBaseURL, // Kokoro on CPU
		TTSModel:               defaultLocalTTSModel,   // kokoro
		TTSVoice:               defaultLocalTTSVoice,   // pt-br alias
		TTSLanguage:            defaultTTSLanguage,     // pt
		TTSFormat:              defaultLocalTTSFormat,  // opus
		TTSSpeed:               defaultTTSSpeed,
		PremiumTTSProvider:     "disabled", // S-35: Only Edge TTS Thalita
		PremiumTTSBaseURL:      "http://127.0.0.1:8012",
		PremiumTTSModel:        "kokoro",
		PremiumTTSVoice:        "pt-BR-ThalitaMultilingualNeural", // S-35: Force Thalita
		TelegramAllowedUserIDs: []int64{},

		MaxIterations:            defaultMaxIterations,
		DBPath:                   filepath.Join(r.Data(), "aurelia.db"),
		MemoryWindowSize:         defaultMemoryWindowSize,
		MCPConfigPath:            filepath.Join(r.Config(), "mcp_servers.json"),
		HeartbeatEnabled:         defaultHeartbeatEnabled,
		HeartbeatIntervalMinutes: defaultHeartbeatIntervalMin,
		VoiceEnabled:             defaultVoiceEnabled,
		VoiceSpoolPath:           filepath.Join(r.Data(), "voice", "spool"),
		VoiceDropPath:            filepath.Join(r.Data(), "voice", "drop"),
		VoiceHeartbeatPath:       filepath.Join(r.Data(), "voice", "heartbeat.json"),
		VoiceHeartbeatFreshSec:   defaultVoiceHeartbeatSec,
		VoicePollIntervalMS:      defaultVoicePollIntervalMS,
		VoiceWakePhrase:          defaultVoiceWakePhrase,
		VoiceCaptureHeartbeat:    filepath.Join(r.Data(), "voice", "capture-heartbeat.json"),
		VoiceCaptureFreshSec:     defaultVoiceCaptureFreshSec,
		VoiceCapturePollMS:       defaultVoiceCapturePollMS,
		GroqSoftCapDaily:         defaultGroqSoftCapDaily,
		GroqHardCapDaily:         defaultGroqHardCapDaily,
		QdrantCollection:         defaultQdrantCollection,
		QdrantEmbeddingModel:     defaultQdrantEmbeddingModel,
		OllamaURL:                "http://127.0.0.1:4000",
		DashboardPort:            defaultDashboardPort,
		HealthPort:               defaultHealthPort,
		AureliaMode:              "sovereign",
		RedisURL:                 defaultRedisURL,
	}
}

// DefaultEditableConfig returns the default user-editable configuration values.
func DefaultEditableConfig() EditableConfig {
	return EditableConfig{
		LLMProvider:            defaultLLMProvider,
		LLMModel:               defaultLLMModelForProvider(defaultLLMProvider),
		STTProvider:            defaultSTTProvider,
		STTLanguage:            defaultSTTLanguage,
		TTSProvider:            defaultTTSProvider,
		TTSBaseURL:             defaultTTSBaseURLForProvider(defaultTTSProvider),
		TTSModel:               defaultTTSModelForProvider(defaultTTSProvider),
		TTSVoice:               defaultTTSVoiceForProvider(defaultTTSProvider),
		TTSLanguage:            defaultTTSLanguage,
		TTSFormat:              defaultTTSFormatForProvider(defaultTTSProvider),
		TTSSpeed:               defaultTTSSpeed,
		PremiumTTSProvider:     "disabled", // S-35: Only Edge TTS Thalita
		PremiumTTSBaseURL:      "http://127.0.0.1:8012",
		PremiumTTSModel:        "kokoro",
		PremiumTTSVoice:        "pt-BR-ThalitaMultilingualNeural", // S-35: Force Thalita
		TelegramAllowedUserIDs: []int64{},

		MaxIterations:            defaultMaxIterations,
		MemoryWindowSize:         defaultMemoryWindowSize,
		HeartbeatEnabled:         defaultHeartbeatEnabled,
		HeartbeatIntervalMinutes: defaultHeartbeatIntervalMin,
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
		LLMProvider:               cfg.LLMProvider,
		LLMModel:                  cfg.LLMModel,
		STTProvider:               cfg.STTProvider,
		STTLanguage:               cfg.STTLanguage,
		TTSProvider:               cfg.TTSProvider,
		TTSBaseURL:                cfg.TTSBaseURL,
		TTSModel:                  cfg.TTSModel,
		TTSVoice:                  cfg.TTSVoice,
		TTSLanguage:               cfg.TTSLanguage,
		TTSFormat:                 cfg.TTSFormat,
		TTSSpeed:                  cfg.TTSSpeed,
		PremiumTTSProvider:        cfg.PremiumTTSProvider,
		PremiumTTSBaseURL:         cfg.PremiumTTSBaseURL,
		PremiumTTSModel:           cfg.PremiumTTSModel,
		PremiumTTSVoice:           cfg.PremiumTTSVoice,
		TelegramBotToken:          cfg.TelegramBotToken,
		TelegramNotificationBotID: cfg.TelegramNotificationBotID,

		TelegramAllowedUserIDs:   append([]int64(nil), cfg.TelegramAllowedUserIDs...),
		AnthropicAPIKey:          cfg.AnthropicAPIKey,
		GoogleAPIKey:             cfg.GoogleAPIKey,
		OpenRouterAPIKey:         cfg.OpenRouterAPIKey,
		OpenAIAPIKey:             cfg.OpenAIAPIKey,
		GroqAPIKey:               cfg.GroqAPIKey,
		MiniMaxAPIKey:            cfg.MiniMaxAPIKey,
		LiteLLMKey:               cfg.LiteLLMKey,
		MaxIterations:            cfg.MaxIterations,
		MemoryWindowSize:         cfg.MemoryWindowSize,
		HeartbeatEnabled:         cfg.HeartbeatEnabled,
		HeartbeatIntervalMinutes: cfg.HeartbeatIntervalMinutes,
	}, nil
}

// SaveEditable updates the user-editable config subset while preserving managed paths.
func SaveEditable(r *runtime.PathResolver, editable EditableConfig) error {
	cfg := defaultFileConfig(r)
	if data, err := os.ReadFile(r.AppConfig()); err == nil && len(data) != 0 {
		_ = json.Unmarshal(data, &cfg)
	}
	cfg.LLMProvider = editable.LLMProvider
	cfg.LLMModel = editable.LLMModel
	cfg.STTProvider = editable.STTProvider
	cfg.STTLanguage = editable.STTLanguage
	cfg.TTSProvider = editable.TTSProvider
	cfg.TTSBaseURL = editable.TTSBaseURL
	cfg.TTSModel = editable.TTSModel
	cfg.TTSVoice = editable.TTSVoice
	cfg.TTSLanguage = editable.TTSLanguage
	cfg.TTSFormat = editable.TTSFormat
	cfg.TTSSpeed = editable.TTSSpeed
	cfg.PremiumTTSProvider = editable.PremiumTTSProvider
	cfg.PremiumTTSBaseURL = editable.PremiumTTSBaseURL
	cfg.PremiumTTSModel = editable.PremiumTTSModel
	cfg.PremiumTTSVoice = editable.PremiumTTSVoice
	cfg.TelegramBotToken = editable.TelegramBotToken

	cfg.TelegramAllowedUserIDs = append([]int64(nil), editable.TelegramAllowedUserIDs...)
	cfg.AnthropicAPIKey = editable.AnthropicAPIKey
	cfg.GoogleAPIKey = editable.GoogleAPIKey
	cfg.OpenRouterAPIKey = editable.OpenRouterAPIKey
	cfg.OpenAIAPIKey = editable.OpenAIAPIKey
	cfg.GroqAPIKey = editable.GroqAPIKey
	cfg.MiniMaxAPIKey = editable.MiniMaxAPIKey
	cfg.LiteLLMKey = editable.LiteLLMKey
	cfg.MaxIterations = editable.MaxIterations
	cfg.MemoryWindowSize = editable.MemoryWindowSize
	cfg.HeartbeatEnabled = editable.HeartbeatEnabled
	cfg.HeartbeatIntervalMinutes = editable.HeartbeatIntervalMinutes
	cfg = normalizeFileConfig(cfg, r)

	return writeConfigFile(r.AppConfig(), cfg)
}

func normalizeFileConfig(cfg fileConfig, r *runtime.PathResolver) fileConfig {
	defaults := defaultFileConfig(r)
	if cfg.TelegramAllowedUserIDs == nil {
		cfg.TelegramAllowedUserIDs = defaults.TelegramAllowedUserIDs
	}
	if cfg.VoiceReplyUserID == 0 && len(cfg.TelegramAllowedUserIDs) > 0 {
		cfg.VoiceReplyUserID = cfg.TelegramAllowedUserIDs[0]
	}
	if cfg.VoiceReplyChatID == 0 && len(cfg.TelegramAllowedUserIDs) > 0 {
		cfg.VoiceReplyChatID = cfg.TelegramAllowedUserIDs[0]
	}
	if cfg.LLMProvider == "" {
		cfg.LLMProvider = defaults.LLMProvider
	}
	if cfg.LLMModel == "" {
		cfg.LLMModel = defaultLLMModelForProvider(cfg.LLMProvider)
	}
	if cfg.STTProvider == "" {
		cfg.STTProvider = defaults.STTProvider
	}
	if cfg.STTLanguage == "" {
		cfg.STTLanguage = defaults.STTLanguage
	}
	if cfg.TTSProvider == "" {
		cfg.TTSProvider = defaults.TTSProvider
	}
	if cfg.TTSBaseURL == "" || usesLegacyTTSDefaults(cfg.TTSProvider, cfg.TTSBaseURL, cfg.TTSModel, cfg.TTSFormat) {
		cfg.TTSBaseURL = defaultTTSBaseURLForProvider(cfg.TTSProvider)
	}
	if cfg.TTSModel == "" || usesLegacyTTSModel(cfg.TTSProvider, cfg.TTSModel) {
		cfg.TTSModel = defaultTTSModelForProvider(cfg.TTSProvider)
	}
	if cfg.TTSVoice == "" || usesLegacyTTSVoice(cfg.TTSProvider, cfg.TTSVoice) {
		cfg.TTSVoice = defaultTTSVoiceForProvider(cfg.TTSProvider)
	}
	if cfg.TTSFormat == "" || usesLegacyTTSFormat(cfg.TTSProvider, cfg.TTSFormat) {
		cfg.TTSFormat = defaultTTSFormatForProvider(cfg.TTSProvider)
	}
	if cfg.TTSLanguage == "" {
		cfg.TTSLanguage = defaults.TTSLanguage
	}
	if cfg.PremiumTTSProvider == "" {
		cfg.PremiumTTSProvider = "openai_compatible"
	}
	if cfg.PremiumTTSBaseURL == "" {
		cfg.PremiumTTSBaseURL = "http://127.0.0.1:8012"
	}
	if cfg.PremiumTTSModel == "" {
		cfg.PremiumTTSModel = "kokoro"
	}
	if cfg.PremiumTTSVoice == "" {
		cfg.PremiumTTSVoice = "pt-br"
	}
	// Valor oficial: sempre manter PT-BR para a identidade vocal soberana.
	if cfg.TTSLanguage != defaults.TTSLanguage {
		cfg.TTSLanguage = defaults.TTSLanguage
	}
	if cfg.TTSSpeed <= 0 {
		cfg.TTSSpeed = defaults.TTSSpeed
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
	if cfg.HeartbeatIntervalMinutes <= 0 {
		cfg.HeartbeatIntervalMinutes = defaults.HeartbeatIntervalMinutes
	}
	if cfg.MCPConfigPath == "" {
		cfg.MCPConfigPath = defaults.MCPConfigPath
	}
	if cfg.VoiceSpoolPath == "" {
		cfg.VoiceSpoolPath = defaults.VoiceSpoolPath
	}
	if cfg.VoiceDropPath == "" {
		cfg.VoiceDropPath = defaults.VoiceDropPath
	}
	if cfg.VoiceHeartbeatPath == "" {
		cfg.VoiceHeartbeatPath = defaults.VoiceHeartbeatPath
	}
	if cfg.VoiceHeartbeatFreshSec <= 0 {
		cfg.VoiceHeartbeatFreshSec = defaults.VoiceHeartbeatFreshSec
	}
	if cfg.VoicePollIntervalMS <= 0 {
		cfg.VoicePollIntervalMS = defaults.VoicePollIntervalMS
	}
	if cfg.VoiceWakePhrase == "" {
		cfg.VoiceWakePhrase = defaults.VoiceWakePhrase
	}
	if cfg.VoiceCaptureHeartbeat == "" {
		cfg.VoiceCaptureHeartbeat = defaults.VoiceCaptureHeartbeat
	}
	if cfg.VoiceCaptureFreshSec <= 0 {
		cfg.VoiceCaptureFreshSec = defaults.VoiceCaptureFreshSec
	}
	if cfg.VoiceCapturePollMS <= 0 {
		cfg.VoiceCapturePollMS = defaults.VoiceCapturePollMS
	}
	cfg.VoiceCaptureCommand = normalizeVoiceCaptureCommand(cfg.VoiceCaptureCommand)
	if cfg.GroqSoftCapDaily <= 0 {
		cfg.GroqSoftCapDaily = defaults.GroqSoftCapDaily
	}
	if cfg.GroqHardCapDaily <= 0 {
		cfg.GroqHardCapDaily = defaults.GroqHardCapDaily
	}
	if cfg.QdrantCollection == "" {
		cfg.QdrantCollection = defaults.QdrantCollection
	}
	if cfg.QdrantEmbeddingModel == "" {
		cfg.QdrantEmbeddingModel = defaults.QdrantEmbeddingModel
	}
	if cfg.OllamaURL == "" {
		cfg.OllamaURL = defaults.OllamaURL
	}
	if cfg.DashboardPort <= 0 {
		cfg.DashboardPort = defaults.DashboardPort
	}
	if cfg.HealthPort <= 0 {
		cfg.HealthPort = defaults.HealthPort
	}
	if cfg.AureliaMode == "" {
		cfg.AureliaMode = "sovereign"
	}

	// Apply Lite Mode overrides if enabled
	if cfg.AureliaMode == "lite" {
		if cfg.LLMProvider == "ollama" {
			// lite = prefer cloud to save local GPU resources
			if cfg.OpenRouterAPIKey != "" {
				cfg.LLMProvider = "openrouter"
				cfg.LLMModel = "anthropic/claude-3-haiku"
			}
		}
		// In lite mode, disable local voice capture unless explicitly configured
		cfg.VoiceEnabled = cfg.TelegramBotToken != "" && cfg.VoiceEnabled
	}

	return cfg
}

func normalizeVoiceCaptureCommand(command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}

	fields := strings.Fields(command)
	if len(fields) == 0 {
		return command
	}

	legacyPath := fields[0]
	if _, err := os.Stat(legacyPath); err == nil {
		return command
	}
	if !strings.Contains(legacyPath, "voice-capture-openwakeword.sh") {
		return command
	}

	cwd, err := os.Getwd()
	if err != nil || strings.TrimSpace(cwd) == "" {
		return command
	}
	localScript := filepath.Join(cwd, "scripts", "voice-capture-openwakeword.sh")
	if _, err := os.Stat(localScript); err != nil {
		return command
	}

	fields[0] = localScript
	return strings.Join(fields, " ")
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
	heartbeatIntervalMin := cfg.HeartbeatIntervalMinutes
	if heartbeatIntervalMin == 0 {
		heartbeatIntervalMin = defaultHeartbeatIntervalMin
	}
	// Backward compat: if no Bots configured but primary token exists, synthesize aurelia entry.
	bots := cfg.Bots
	if len(bots) == 0 && cfg.TelegramBotToken != "" {
		bots = []BotConfig{{
			ID:             "aurelia",
			Name:           "Aurelia_Code",
			Token:          cfg.TelegramBotToken,
			AllowedUserIDs: append([]int64(nil), cfg.TelegramAllowedUserIDs...),
			PersonaID:      "aurelia-sovereign",
			FocusArea:      "Comando soberano do ecossistema, orquestra time e arbitra prioridades",
			Enabled:        true,
		}}
	}

	return &AppConfig{
		Features:                  LoadFeatures(),
		Bots:                      bots,
		LLMProvider:               cfg.LLMProvider,
		LLMModel:                  cfg.LLMModel,
		STTProvider:               cfg.STTProvider,
		STTBaseURL:                cfg.STTBaseURL,
		STTModel:                  cfg.STTModel,
		STTLanguage:               cfg.STTLanguage,
		TTSProvider:               cfg.TTSProvider,
		TTSBaseURL:                cfg.TTSBaseURL,
		TTSModel:                  cfg.TTSModel,
		TTSVoice:                  cfg.TTSVoice,
		TTSLanguage:               cfg.TTSLanguage,
		TTSFormat:                 cfg.TTSFormat,
		TTSSpeed:                  cfg.TTSSpeed,
		PremiumTTSProvider:        cfg.PremiumTTSProvider,
		PremiumTTSBaseURL:         cfg.PremiumTTSBaseURL,
		PremiumTTSModel:           cfg.PremiumTTSModel,
		PremiumTTSVoice:           cfg.PremiumTTSVoice,
		TelegramBotToken:          cfg.TelegramBotToken,
		TelegramNotificationBotID: cfg.TelegramNotificationBotID,

		TelegramAllowedUserIDs:   cfg.TelegramAllowedUserIDs,
		AnthropicAPIKey:          cfg.AnthropicAPIKey,
		GoogleAPIKey:             cfg.GoogleAPIKey,
		OpenRouterAPIKey:         cfg.OpenRouterAPIKey,
		OpenAIAPIKey:             cfg.OpenAIAPIKey,
		GroqAPIKey:               cfg.GroqAPIKey,
		MiniMaxAPIKey:            cfg.MiniMaxAPIKey,
		LiteLLMKey:               cfg.LiteLLMKey,
		MaxIterations:            cfg.MaxIterations,
		DBPath:                   cfg.DBPath,
		MemoryWindowSize:         cfg.MemoryWindowSize,
		MCPConfigPath:            cfg.MCPConfigPath,
		HeartbeatEnabled:         cfg.HeartbeatEnabled || defaultHeartbeatEnabled,
		HeartbeatIntervalMinutes: heartbeatIntervalMin,
		VoiceEnabled:             cfg.VoiceEnabled,
		VoiceReplyUserID:         cfg.VoiceReplyUserID,
		VoiceReplyChatID:         cfg.VoiceReplyChatID,
		VoiceSpoolPath:           cfg.VoiceSpoolPath,
		VoiceDropPath:            cfg.VoiceDropPath,
		VoiceHeartbeatPath:       cfg.VoiceHeartbeatPath,
		VoiceHeartbeatFreshSec:   cfg.VoiceHeartbeatFreshSec,
		VoicePollIntervalMS:      cfg.VoicePollIntervalMS,
		VoiceWakePhrase:          cfg.VoiceWakePhrase,
		VoiceCaptureEnabled:      cfg.VoiceCaptureEnabled,
		VoiceCaptureCommand:      cfg.VoiceCaptureCommand,
		VoiceCaptureHeartbeat:    cfg.VoiceCaptureHeartbeat,
		VoiceCaptureFreshSec:     cfg.VoiceCaptureFreshSec,
		VoiceCapturePollMS:       cfg.VoiceCapturePollMS,
		STTFallbackCommand:       cfg.STTFallbackCommand,
		GroqSoftCapDaily:         cfg.GroqSoftCapDaily,
		GroqHardCapDaily:         cfg.GroqHardCapDaily,
		QdrantURL:                cfg.QdrantURL,
		QdrantAPIKey:             cfg.QdrantAPIKey,
		QdrantCollection:         cfg.QdrantCollection,
		QdrantEmbeddingModel:     cfg.QdrantEmbeddingModel,
		OllamaURL:                cfg.OllamaURL,
		DashboardPort:            cfg.DashboardPort,
		HealthPort:               cfg.HealthPort,
		ObsidianVaultPath:        cfg.ObsidianVaultPath,
		ObsidianSyncEnabled:      cfg.ObsidianSyncEnabled,
		AureliaMode:              cfg.AureliaMode,
	}
}

func sameFileConfig(a, b fileConfig) bool {
	if a.TelegramBotToken != b.TelegramBotToken ||
		a.TelegramNotificationBotID != b.TelegramNotificationBotID ||
		a.LLMProvider != b.LLMProvider ||
		a.LLMModel != b.LLMModel ||
		a.STTProvider != b.STTProvider ||
		a.STTLanguage != b.STTLanguage ||
		a.TTSProvider != b.TTSProvider ||
		a.TTSBaseURL != b.TTSBaseURL ||
		a.TTSModel != b.TTSModel ||
		a.TTSVoice != b.TTSVoice ||
		a.TTSLanguage != b.TTSLanguage ||
		a.TTSFormat != b.TTSFormat ||
		a.TTSSpeed != b.TTSSpeed ||
		a.PremiumTTSProvider != b.PremiumTTSProvider ||
		a.PremiumTTSBaseURL != b.PremiumTTSBaseURL ||
		a.PremiumTTSModel != b.PremiumTTSModel ||
		a.PremiumTTSVoice != b.PremiumTTSVoice ||

		a.AnthropicAPIKey != b.AnthropicAPIKey ||
		a.GoogleAPIKey != b.GoogleAPIKey ||
		a.OpenRouterAPIKey != b.OpenRouterAPIKey ||
		a.OpenAIAPIKey != b.OpenAIAPIKey ||
		a.GroqAPIKey != b.GroqAPIKey ||
		a.MiniMaxAPIKey != b.MiniMaxAPIKey ||
		a.MaxIterations != b.MaxIterations ||
		a.DBPath != b.DBPath ||
		a.MemoryWindowSize != b.MemoryWindowSize ||
		a.MCPConfigPath != b.MCPConfigPath ||
		a.HeartbeatEnabled != b.HeartbeatEnabled ||
		a.HeartbeatIntervalMinutes != b.HeartbeatIntervalMinutes ||
		a.VoiceEnabled != b.VoiceEnabled ||
		a.VoiceReplyUserID != b.VoiceReplyUserID ||
		a.VoiceReplyChatID != b.VoiceReplyChatID ||
		a.VoiceSpoolPath != b.VoiceSpoolPath ||
		a.VoiceDropPath != b.VoiceDropPath ||
		a.VoiceHeartbeatPath != b.VoiceHeartbeatPath ||
		a.VoiceHeartbeatFreshSec != b.VoiceHeartbeatFreshSec ||
		a.VoicePollIntervalMS != b.VoicePollIntervalMS ||
		a.VoiceWakePhrase != b.VoiceWakePhrase ||
		a.VoiceCaptureEnabled != b.VoiceCaptureEnabled ||
		a.VoiceCaptureCommand != b.VoiceCaptureCommand ||
		a.VoiceCaptureHeartbeat != b.VoiceCaptureHeartbeat ||
		a.VoiceCaptureFreshSec != b.VoiceCaptureFreshSec ||
		a.VoiceCapturePollMS != b.VoiceCapturePollMS ||
		a.STTFallbackCommand != b.STTFallbackCommand ||
		a.GroqSoftCapDaily != b.GroqSoftCapDaily ||
		a.GroqHardCapDaily != b.GroqHardCapDaily ||
		a.QdrantURL != b.QdrantURL ||
		a.QdrantAPIKey != b.QdrantAPIKey ||
		a.QdrantCollection != b.QdrantCollection ||
		a.QdrantEmbeddingModel != b.QdrantEmbeddingModel {
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

func defaultTTSBaseURLForProvider(provider string) string {
	// TTS uses Kokoro on CPU.
	return defaultLocalTTSBaseURL
}

func defaultTTSModelForProvider(provider string) string {
	// Kokoro is the only local TTS runtime.
	return defaultLocalTTSModel
}

func defaultTTSVoiceForProvider(provider string) string {
	// Keep the legacy PT-BR alias at app level; the TTS layer resolves it.
	return defaultLocalTTSVoice
}

func defaultTTSFormatForProvider(provider string) string {
	// Opus format for Kokoro
	return defaultLocalTTSFormat
}

func usesLegacyTTSDefaults(provider, baseURL, model, format string) bool {
	// All TTS providers now use local Kokoro, no legacy cloud providers.
	return false
}

func usesLegacyTTSModel(provider, model string) bool {
	// All TTS providers now use local Kokoro, no legacy cloud models.
	return false
}

func usesLegacyTTSVoice(provider, voice string) bool {
	// All TTS providers now use local Kokoro, no legacy cloud voices.
	return false
}

func usesLegacyTTSFormat(provider, format string) bool {
	// All TTS providers now use local Kokoro with opus format.
	return false
}

func defaultLLMModelForProvider(provider string) string {
	switch provider {
	case "anthropic":
		return "claude-sonnet-4-6"
	case "google":
		return "gemini-2.5-pro"
	case "ollama":
		return "gemma3:27b-it-qat"
	case "openrouter":
		return "openrouter/nvidia/nemotron-3-super-120b-a12b:free"
	case "openai":
		return "gpt-5.4"
	default:
		return defaultLLMModel
	}
}
