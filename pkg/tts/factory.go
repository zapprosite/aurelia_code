package tts

import (
	"github.com/kocar/aurelia/internal/config"
)

// NewSynthesizer builds a TTS synthesizer based on the provided configuration.
// It prioritizes local OpenAI-compatible providers (like Kokoro).
func NewSynthesizer(baseURL, provider, model, voice, language, format string, speed float64) Synthesizer {
	if provider == "" || provider == "disabled" {
		return nil
	}

	// Currently, all our providers use the OpenAI-compatible API.
	// In the future, we could add native support for other engines here.
	base := NewOpenAICompatibleSynthesizer(baseURL, model, voice, language, format, speed)
	// 3500 chars ≈ 600-700 tokens PT-BR — dentro da janela do Kokoro/Kodoro.
	// Menos chunks = menos stitches de áudio = menos artefatos de corte abrupto.
	return NewSegmentedSynthesizer(base, 3500)
}

// NewDefaultSynthesizer builds the primary TTS engine from the app config.
func NewDefaultSynthesizer(cfg *config.AppConfig) Synthesizer {
	if cfg == nil {
		return nil
	}
	return NewSynthesizer(
		cfg.TTSBaseURL,
		cfg.TTSProvider,
		cfg.TTSModel,
		cfg.TTSVoice,
		cfg.TTSLanguage,
		cfg.TTSFormat,
		cfg.TTSSpeed,
	)
}

// NewPremiumSynthesizer builds the premium TTS engine from the app config.
func NewPremiumSynthesizer(cfg *config.AppConfig) Synthesizer {
	if cfg == nil {
		return nil
	}
	provider := cfg.PremiumTTSProvider
	if provider == "" || provider == "disabled" {
		return nil
	}
	
	return NewSynthesizer(
		cfg.PremiumTTSBaseURL,
		provider,
		cfg.PremiumTTSModel,
		cfg.PremiumTTSVoice,
		cfg.TTSLanguage, // Use same language as default
		"opus",           // Premium always uses high-quality opus
		cfg.TTSSpeed,
	)
}
