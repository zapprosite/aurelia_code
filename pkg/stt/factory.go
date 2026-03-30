package stt

import "fmt"

// NewTranscriber builds a Transcriber based on the provider name.
// groqAPIKey is used only for the "groq" provider.
// baseURL, model, language accept empty strings to use provider defaults.
//
// Priority: local (faster-whisper-server on :8020) → Groq cloud fallback.
func NewTranscriber(provider, groqAPIKey, baseURL, model, language string) (Transcriber, error) {
	local := NewLocalTranscriber(baseURL, model, language)
	groq := NewGroqTranscriber(groqAPIKey, baseURL, model, language)

	switch provider {
	case "local", "faster-whisper":
		// LOCAL FIRST: faster-whisper-server → Groq fallback
		if local.IsAvailable() {
			return local, nil
		}
		if groq.IsAvailable() {
			return groq, nil
		}
		return nil, fmt.Errorf("no stt provider available: faster-whisper-server unreachable and groq api key not set")
	case "", "groq":
		// GROQ CLOUD FIRST: Groq → faster-whisper fallback
		if groq.IsAvailable() {
			return groq, nil
		}
		if local.IsAvailable() {
			return local, nil
		}
		return nil, fmt.Errorf("no stt provider available: groq api key not set and faster-whisper-server unreachable")
	default:
		return nil, fmt.Errorf("unsupported stt provider %q", provider)
	}
}
