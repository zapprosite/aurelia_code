package stt

import "context"

// Transcriber defines the interface for modern, sovereign speech-to-text.
type Transcriber interface {
	// Transcribe converts audio file content into text.
	Transcribe(ctx context.Context, audioFilePath string) (string, error)
	// IsAvailable checks the health of the transcription service.
	IsAvailable() bool
}
