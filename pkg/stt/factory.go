package stt

import "fmt"

// NewTranscriber builds a Transcriber based on the provider name.
// groqAPIKey is used only for the "groq" provider.
// baseURL, model, language accept empty strings to use provider defaults.
func NewTranscriber(provider, groqAPIKey, baseURL, model, language string) (Transcriber, error) {
	switch provider {
	case "local", "faster-whisper":
		return NewLocalTranscriber(baseURL, model, language), nil
	case "", "groq":
		return NewGroqTranscriber(groqAPIKey, baseURL, model, language), nil
	default:
		return nil, fmt.Errorf("unsupported stt provider %q", provider)
	}
}
