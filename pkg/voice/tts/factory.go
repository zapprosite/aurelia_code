package tts

import (
	"os"
	"strings"

	"github.com/kocar/aurelia/internal/config"
)

// NewSynthesizer builds a TTS synthesizer based on the provided configuration.
// Supports Edge TTS (free PT-BR voices) and OpenAI-compatible providers (Kokoro).
func NewSynthesizer(baseURL, provider, model, voice, language, format string, speed float64) Synthesizer {
	if provider == "" || provider == "disabled" {
		return nil
	}

	// S-35: Edge TTS is the only TTS (Thalita, PT-BR native, free)
	if provider == "edge" || isEdgeVoice(voice) {
		venvPath := os.Getenv("AURELIA_VENV")
		if venvPath == "" {
			venvPath = "/home/will/aurelia/.venv"
		}
		// Use Thalita as default if voice is empty
		if voice == "" {
			voice = "pt-BR-ThalitaMultilingualNeural"
		}
		base := NewEdgeSynthesizer(venvPath, voice)
		// Edge TTS works better with smaller chunks to avoid timeouts
		return NewSegmentedSynthesizer(base, 1500)
	}

	// OpenAI-compatible providers (Kokoro, etc.)
	base := NewOpenAICompatibleSynthesizer(baseURL, model, voice, language, format, speed)
	return NewSegmentedSynthesizer(base, 3500)
}

// isEdgeVoice checks if the voice is an Edge TTS voice
func isEdgeVoice(voice string) bool {
	edgePrefixes := []string{
		"pt-BR-",
		"pt-PT-",
		"en-",
		"es-",
		"fr-",
		"de-",
		"it-",
		"ja-",
		"ko-",
		"zh-",
	}
	for _, prefix := range edgePrefixes {
		if strings.HasPrefix(voice, prefix) && strings.Contains(voice, "Neural") {
			return true
		}
	}
	return false
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
		"opus",          // Premium always uses high-quality opus
		cfg.TTSSpeed,
	)
}
