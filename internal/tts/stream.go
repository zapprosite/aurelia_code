// Package tts provides streaming TTS capabilities for the Jarvis voice loop
// ADR: 20260328-e2e-jarvis-loop-wake-tts

package tts

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Stream provides streaming TTS synthesis
// ADR: 20260328-e2e-jarvis-loop-wake-tts
type Stream struct {
	baseURL   string
	model     string
	voice     string
	language  string
	client    *http.Client
	chunkSize int
}

// NewStream creates a new TTS streaming instance
func NewStream(baseURL, model, voice, language string) *Stream {
	return &Stream{
		baseURL:   baseURL,
		model:     model,
		voice:     voice,
		language:  language,
		client:    &http.Client{Timeout: 30 * time.Second},
		chunkSize: 1200, // Characters per TTS chunk
	}
}

// StreamChunk represents a chunk of synthesized audio
type StreamChunk struct {
	Audio []byte
	Text  string
	Error error
}

// Stream synthesizes text and streams audio chunks
// ADR: 20260328-e2e-jarvis-loop-wake-tts
func (s *Stream) Stream(ctx context.Context, text string) <-chan StreamChunk {
	chunks := make(chan StreamChunk, 10)

	go func() {
		defer close(chunks)

		// Segment text into chunks
		segmented := segmentText(text, s.chunkSize)

		for _, segment := range segmented {
			select {
			case <-ctx.Done():
				chunks <- StreamChunk{Error: ctx.Err()}
				return
			default:
			}

			// Synthesize this chunk
			audio, err := s.synthesize(ctx, segment)
			if err != nil {
				chunks <- StreamChunk{Text: segment, Error: err}
				continue
			}

			chunks <- StreamChunk{
				Audio: audio,
				Text:  segment,
			}
		}
	}()

	return chunks
}

// synthesize sends a single text segment to TTS
func (s *Stream) synthesize(ctx context.Context, text string) ([]byte, error) {
	// Build request
	reqBody := map[string]string{
		"model": s.model,
		"input": text,
		"voice": s.voice,
	}

	if s.language != "" {
		reqBody["language"] = s.language
	}

	// Simplified - in production would use proper HTTP client
	resp, err := s.client.Post(s.baseURL+"/v1/audio/speech", "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("tts request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tts returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// segmentText splits text into chunks that fit within TTS limits
func segmentText(text string, maxChars int) []string {
	if len(text) <= maxChars {
		return []string{text}
	}

	var chunks []string
	sentences := strings.Split(text, ".")

	var current strings.Builder
	for _, sentence := range sentences {
		if current.Len()+len(sentence)+1 > maxChars {
			if current.Len() > 0 {
				chunks = append(chunks, strings.TrimSpace(current.String()))
				current.Reset()
			}
		}
		if current.Len() > 0 {
			current.WriteString(". ")
		}
		current.WriteString(sentence)
	}

	if current.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(current.String()))
	}

	return chunks
}

// SynthesizeURL returns a URL to the synthesized audio (non-streaming)
func (s *Stream) SynthesizeURL(ctx context.Context, text string) (string, error) {
	// For non-streaming, just return empty for now
	// In production, would upload to storage and return URL
	return "", nil
}
