package stt

import "fmt"

// NewTranscriber builds a Transcriber based on the provider name.
// For SOTA 2026, we prioritize the local "faster-whisper-server".
func NewTranscriber(provider, baseURL, model, language string) (Transcriber, error) {
	local := NewLocalTranscriber(baseURL, model, language)

	switch provider {
	case "local", "faster-whisper", "":
		if local.IsAvailable() {
			return local, nil
		}
		// In a sovereign environment, we only support local.
		return nil, fmt.Errorf("transcription provider 'local' is not available at %s", baseURL)
	default:
		return nil, fmt.Errorf("unsupported or legacy stt provider %q", provider)
	}
}
