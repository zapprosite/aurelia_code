package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kocar/aurelia/internal/runtime"
)

const (
	defaultMaxIterations        = 500
	defaultMemoryWindowSize     = 20
	defaultLLMProvider          = "kimi"
	defaultLLMModel             = "kimi-k2-thinking"
	defaultSTTProvider          = "groq"
	defaultTTSProvider          = "openai_compatible"
	defaultLocalTTSBaseURL = "http://127.0.0.1:8011"
	defaultLocalTTSModel   = "chatterbox"
	defaultLocalTTSVoice   = "Aurelia.wav" // PT-BR voice via gTTS (sweet, educated)
	defaultLocalTTSFormat  = "opus"
	defaultTTSSpeed        = 1.0
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
	defaultSupabaseEventsTable  = "aurelia_voice_events"
	defaultQdrantCollection     = "conversation_memory"
	defaultQdrantEmbeddingModel = "bge-m3"
)

// AppConfig holds all runtime configuration needed for the application.
type AppConfig struct {
	LLMProvider              string
	LLMModel                 string
	STTProvider              string
	TTSProvider              string
	TTSBaseURL               string
	TTSModel                 string
	TTSVoice                 string
	TTSFormat                string
	TTSSpeed                 float64
	TelegramBotToken         string
	TelegramAllowedUserIDs   []int64
	AnthropicAPIKey          string
	GoogleAPIKey             string
	KiloAPIKey               string
	KimiAPIKey               string
	OpenRouterAPIKey         string
	ZAIAPIKey                string
	AlibabaAPIKey            string
	OpenAIAPIKey             string
	OpenAIAuthMode           string
	GroqAPIKey               string
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
	SupabaseURL              string
	SupabaseServiceRoleKey   string
	SupabaseEventsTable      string
	QdrantURL                string
	QdrantAPIKey             string
	QdrantCollection         string
	QdrantEmbeddingModel     string
}

type fileConfig struct {
	LLMProvider              string  `json:"llm_provider"`
	LLMModel                 string  `json:"llm_model"`
	STTProvider              string  `json:"stt_provider"`
	TTSProvider              string  `json:"tts_provider"`
	TTSBaseURL               string  `json:"tts_base_url"`
	TTSModel                 string  `json:"tts_model"`
	TTSVoice                 string  `json:"tts_voice"`
	TTSFormat                string  `json:"tts_format"`
	TTSSpeed                 float64 `json:"tts_speed"`
	TelegramBotToken         string  `json:"telegram_bot_token"`
	TelegramAllowedUserIDs   []int64 `json:"telegram_allowed_user_ids"`
	AnthropicAPIKey          string  `json:"anthropic_api_key"`
	GoogleAPIKey             string  `json:"google_api_key"`
	KiloAPIKey               string  `json:"kilo_api_key"`
	KimiAPIKey               string  `json:"kimi_api_key"`
	OpenRouterAPIKey         string  `json:"openrouter_api_key"`
	ZAIAPIKey                string  `json:"zai_api_key"`
	AlibabaAPIKey            string  `json:"alibaba_api_key"`
	OpenAIAPIKey             string  `json:"openai_api_key"`
	OpenAIAuthMode           string  `json:"openai_auth_mode"`
	GroqAPIKey               string  `json:"groq_api_key"`
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
	SupabaseURL              string  `json:"supabase_url"`
	SupabaseServiceRoleKey   string  `json:"supabase_service_role_key"`
	SupabaseEventsTable      string  `json:"supabase_events_table"`
	QdrantURL                string  `json:"qdrant_url"`
	QdrantAPIKey             string  `json:"qdrant_api_key"`
	QdrantCollection         string  `json:"qdrant_collection"`
	QdrantEmbeddingModel     string  `json:"qdrant_embedding_model"`
}

// EditableConfig represents the user-editable portion of the runtime config.
type EditableConfig struct {
	LLMProvider              string
	LLMModel                 string
	STTProvider              string
	TTSProvider              string
	TTSBaseURL               string
	TTSModel                 string
	TTSVoice                 string
	TTSFormat                string
	TTSSpeed                 float64
	TelegramBotToken         string
	TelegramAllowedUserIDs   []int64
	AnthropicAPIKey          string
	GoogleAPIKey             string
	KiloAPIKey               string
	KimiAPIKey               string
	OpenRouterAPIKey         string
	ZAIAPIKey                string
	AlibabaAPIKey            string
	OpenAIAPIKey             string
	OpenAIAuthMode           string
	GroqAPIKey               string
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
	if !sameFileConfig(normalized, cfg) {
		if err := writeConfigFile(path, normalized); err != nil {
			return nil, err
		}
	}

	return toAppConfig(normalized), nil
}

func defaultFileConfig(r *runtime.PathResolver) fileConfig {
	return fileConfig{
		LLMProvider:              defaultLLMProvider,
		LLMModel:                 defaultLLMModelForProvider(defaultLLMProvider),
		OpenAIAuthMode:           "api_key",
		STTProvider:              defaultSTTProvider,
		TTSProvider:              defaultTTSProvider,
		TTSBaseURL:               defaultTTSBaseURLForProvider(defaultTTSProvider),
		TTSModel:                 defaultTTSModelForProvider(defaultTTSProvider),
		TTSVoice:                 defaultTTSVoiceForProvider(defaultTTSProvider),
		TTSFormat:                defaultTTSFormatForProvider(defaultTTSProvider),
		TTSSpeed:                 defaultTTSSpeed,
		TelegramAllowedUserIDs:   []int64{},
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
		SupabaseEventsTable:      defaultSupabaseEventsTable,
		QdrantCollection:         defaultQdrantCollection,
		QdrantEmbeddingModel:     defaultQdrantEmbeddingModel,
	}
}

// DefaultEditableConfig returns the default user-editable configuration values.
func DefaultEditableConfig() EditableConfig {
	return EditableConfig{
		LLMProvider:              defaultLLMProvider,
		LLMModel:                 defaultLLMModelForProvider(defaultLLMProvider),
		OpenAIAuthMode:           "api_key",
		STTProvider:              defaultSTTProvider,
		TTSProvider:              defaultTTSProvider,
		TTSBaseURL:               defaultTTSBaseURLForProvider(defaultTTSProvider),
		TTSModel:                 defaultTTSModelForProvider(defaultTTSProvider),
		TTSVoice:                 defaultTTSVoiceForProvider(defaultTTSProvider),
		TTSFormat:                defaultTTSFormatForProvider(defaultTTSProvider),
		TTSSpeed:                 defaultTTSSpeed,
		TelegramAllowedUserIDs:   []int64{},
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
		LLMProvider:              cfg.LLMProvider,
		LLMModel:                 cfg.LLMModel,
		STTProvider:              cfg.STTProvider,
		TTSProvider:              cfg.TTSProvider,
		TTSBaseURL:               cfg.TTSBaseURL,
		TTSModel:                 cfg.TTSModel,
		TTSVoice:                 cfg.TTSVoice,
		TTSFormat:                cfg.TTSFormat,
		TTSSpeed:                 cfg.TTSSpeed,
		TelegramBotToken:         cfg.TelegramBotToken,
		TelegramAllowedUserIDs:   append([]int64(nil), cfg.TelegramAllowedUserIDs...),
		AnthropicAPIKey:          cfg.AnthropicAPIKey,
		GoogleAPIKey:             cfg.GoogleAPIKey,
		KiloAPIKey:               cfg.KiloAPIKey,
		KimiAPIKey:               cfg.KimiAPIKey,
		OpenRouterAPIKey:         cfg.OpenRouterAPIKey,
		ZAIAPIKey:                cfg.ZAIAPIKey,
		AlibabaAPIKey:            cfg.AlibabaAPIKey,
		OpenAIAPIKey:             cfg.OpenAIAPIKey,
		OpenAIAuthMode:           cfg.OpenAIAuthMode,
		GroqAPIKey:               cfg.GroqAPIKey,
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
	cfg.TTSProvider = editable.TTSProvider
	cfg.TTSBaseURL = editable.TTSBaseURL
	cfg.TTSModel = editable.TTSModel
	cfg.TTSVoice = editable.TTSVoice
	cfg.TTSFormat = editable.TTSFormat
	cfg.TTSSpeed = editable.TTSSpeed
	cfg.TelegramBotToken = editable.TelegramBotToken
	cfg.TelegramAllowedUserIDs = append([]int64(nil), editable.TelegramAllowedUserIDs...)
	cfg.AnthropicAPIKey = editable.AnthropicAPIKey
	cfg.GoogleAPIKey = editable.GoogleAPIKey
	cfg.KiloAPIKey = editable.KiloAPIKey
	cfg.KimiAPIKey = editable.KimiAPIKey
	cfg.OpenRouterAPIKey = editable.OpenRouterAPIKey
	cfg.ZAIAPIKey = editable.ZAIAPIKey
	cfg.AlibabaAPIKey = editable.AlibabaAPIKey
	cfg.OpenAIAPIKey = editable.OpenAIAPIKey
	cfg.OpenAIAuthMode = editable.OpenAIAuthMode
	cfg.GroqAPIKey = editable.GroqAPIKey
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
	if cfg.OpenAIAuthMode == "" {
		cfg.OpenAIAuthMode = defaults.OpenAIAuthMode
	}
	if cfg.STTProvider == "" {
		cfg.STTProvider = defaults.STTProvider
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
	if cfg.GroqSoftCapDaily <= 0 {
		cfg.GroqSoftCapDaily = defaults.GroqSoftCapDaily
	}
	if cfg.GroqHardCapDaily <= 0 {
		cfg.GroqHardCapDaily = defaults.GroqHardCapDaily
	}
	if cfg.SupabaseEventsTable == "" {
		cfg.SupabaseEventsTable = defaults.SupabaseEventsTable
	}
	if cfg.QdrantCollection == "" {
		cfg.QdrantCollection = defaults.QdrantCollection
	}
	if cfg.QdrantEmbeddingModel == "" {
		cfg.QdrantEmbeddingModel = defaults.QdrantEmbeddingModel
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
	heartbeatIntervalMin := cfg.HeartbeatIntervalMinutes
	if heartbeatIntervalMin == 0 {
		heartbeatIntervalMin = defaultHeartbeatIntervalMin
	}
	return &AppConfig{
		LLMProvider:              cfg.LLMProvider,
		LLMModel:                 cfg.LLMModel,
		STTProvider:              cfg.STTProvider,
		TTSProvider:              cfg.TTSProvider,
		TTSBaseURL:               cfg.TTSBaseURL,
		TTSModel:                 cfg.TTSModel,
		TTSVoice:                 cfg.TTSVoice,
		TTSFormat:                cfg.TTSFormat,
		TTSSpeed:                 cfg.TTSSpeed,
		TelegramBotToken:         cfg.TelegramBotToken,
		TelegramAllowedUserIDs:   cfg.TelegramAllowedUserIDs,
		AnthropicAPIKey:          cfg.AnthropicAPIKey,
		GoogleAPIKey:             cfg.GoogleAPIKey,
		KiloAPIKey:               cfg.KiloAPIKey,
		KimiAPIKey:               cfg.KimiAPIKey,
		OpenRouterAPIKey:         cfg.OpenRouterAPIKey,
		ZAIAPIKey:                cfg.ZAIAPIKey,
		AlibabaAPIKey:            cfg.AlibabaAPIKey,
		OpenAIAPIKey:             cfg.OpenAIAPIKey,
		OpenAIAuthMode:           cfg.OpenAIAuthMode,
		GroqAPIKey:               cfg.GroqAPIKey,
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
		SupabaseURL:              cfg.SupabaseURL,
		SupabaseServiceRoleKey:   cfg.SupabaseServiceRoleKey,
		SupabaseEventsTable:      cfg.SupabaseEventsTable,
		QdrantURL:                cfg.QdrantURL,
		QdrantAPIKey:             cfg.QdrantAPIKey,
		QdrantCollection:         cfg.QdrantCollection,
		QdrantEmbeddingModel:     cfg.QdrantEmbeddingModel,
	}
}

func sameFileConfig(a, b fileConfig) bool {
	if a.TelegramBotToken != b.TelegramBotToken ||
		a.LLMProvider != b.LLMProvider ||
		a.LLMModel != b.LLMModel ||
		a.STTProvider != b.STTProvider ||
		a.TTSProvider != b.TTSProvider ||
		a.TTSBaseURL != b.TTSBaseURL ||
		a.TTSModel != b.TTSModel ||
		a.TTSVoice != b.TTSVoice ||
		a.TTSFormat != b.TTSFormat ||
		a.TTSSpeed != b.TTSSpeed ||
		a.AnthropicAPIKey != b.AnthropicAPIKey ||
		a.GoogleAPIKey != b.GoogleAPIKey ||
		a.KiloAPIKey != b.KiloAPIKey ||
		a.KimiAPIKey != b.KimiAPIKey ||
		a.OpenRouterAPIKey != b.OpenRouterAPIKey ||
		a.ZAIAPIKey != b.ZAIAPIKey ||
		a.AlibabaAPIKey != b.AlibabaAPIKey ||
		a.OpenAIAPIKey != b.OpenAIAPIKey ||
		a.OpenAIAuthMode != b.OpenAIAuthMode ||
		a.GroqAPIKey != b.GroqAPIKey ||
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
		a.SupabaseURL != b.SupabaseURL ||
		a.SupabaseServiceRoleKey != b.SupabaseServiceRoleKey ||
		a.SupabaseEventsTable != b.SupabaseEventsTable ||
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
	// All TTS providers use local voice-proxy service
	return defaultLocalTTSBaseURL
}

func defaultTTSModelForProvider(provider string) string {
	// All TTS providers use local voice-proxy service
	return defaultLocalTTSModel
}

func defaultTTSVoiceForProvider(provider string) string {
	// All TTS providers use local voice-proxy with Olivia.wav
	return defaultLocalTTSVoice
}

func defaultTTSFormatForProvider(provider string) string {
	// All TTS providers use local voice-proxy with opus format
	return defaultLocalTTSFormat
}

func usesLegacyTTSDefaults(provider, baseURL, model, format string) bool {
	// All TTS providers now use local voice-proxy, no legacy cloud providers
	return false
}

func usesLegacyTTSModel(provider, model string) bool {
	// All TTS providers now use local voice-proxy, no legacy cloud models
	return false
}

func usesLegacyTTSVoice(provider, voice string) bool {
	// All TTS providers now use local voice-proxy with Olivia.wav, no legacy cloud voices
	return false
}

func usesLegacyTTSFormat(provider, format string) bool {
	// All TTS providers now use local voice-proxy with opus format, no legacy cloud formats
	return false
}

func defaultLLMModelForProvider(provider string) string {
	switch provider {
	case "anthropic":
		return "claude-sonnet-4-6"
	case "google":
		return "gemini-2.5-pro"
	case "kilo":
		return "gpt-5.4"
	case "ollama":
		return "qwen3.5:9b"
	case "openrouter":
		return "openrouter/auto"
	case "zai":
		return "glm-5"
	case "alibaba":
		return "qwen3-coder-plus"
	case "openai":
		return "gpt-5.4"
	default:
		return defaultLLMModel
	}
}
