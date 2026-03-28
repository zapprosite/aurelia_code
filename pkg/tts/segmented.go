package tts

import (
	"context"
	"fmt"
	"strings"
)

// SegmentedSynthesizer is a decorator that splits long text into smaller chunks
// before synthesis to prevent memory exhaustion or timeouts in the TTS engine.
type SegmentedSynthesizer struct {
	base     Synthesizer
	maxChars int
}

// NewSegmentedSynthesizer wraps a synthesizer with segmentation logic.
func NewSegmentedSynthesizer(base Synthesizer, maxChars int) *SegmentedSynthesizer {
	if maxChars <= 0 {
		maxChars = 2500 // Optimized default for Kodoro/Kokoro SOTA 2026
	}
	return &SegmentedSynthesizer{
		base:     base,
		maxChars: maxChars,
	}
}

func (s *SegmentedSynthesizer) IsAvailable() bool {
	return s.base != nil && s.base.IsAvailable()
}

func (s *SegmentedSynthesizer) MaxChars() int {
	return 1000000 
}

func (s *SegmentedSynthesizer) Synthesize(ctx context.Context, text string) (Audio, error) {
	if !s.IsAvailable() {
		return Audio{}, fmt.Errorf("base synthesizer is not available")
	}

	chunks := splitText(text, s.maxChars)
	if len(chunks) == 0 {
		return Audio{}, nil
	}

	if len(chunks) == 1 {
		return s.base.Synthesize(ctx, chunks[0])
	}

	var combinedData []byte
	var lastAudio Audio

	fmt.Printf("TTS: Starting synthesis for %d chunks (Infinite-Voice SOTA 2026)\n", len(chunks))

	for i, chunk := range chunks {
		select {
		case <-ctx.Done():
			return Audio{}, ctx.Err()
		default:
			fmt.Printf("TTS: Processing chunk %d/%d (%d chars)...\n", i+1, len(chunks), len(chunk))
			audio, err := s.base.Synthesize(ctx, chunk)
			if err != nil {
				return Audio{}, fmt.Errorf("synthesize chunk %d/%d: %w", i+1, len(chunks), err)
			}
			combinedData = append(combinedData, audio.Data...)
			lastAudio = audio
		}
	}

	fmt.Printf("TTS: Synthesis complete. Total audio size: %d bytes\n", len(combinedData))

	return Audio{
		Data:        combinedData,
		ContentType: lastAudio.ContentType,
		Extension:   lastAudio.Extension,
		AsVoiceNote: lastAudio.AsVoiceNote,
	}, nil
}

// splitText divides text into chunks of at most maxChars, trying to respect sentences.
func splitText(text string, maxChars int) []string {
	if len(text) <= maxChars {
		if t := strings.TrimSpace(text); t == "" {
			return nil
		}
		return []string{text}
	}

	var chunks []string
	remaining := text

	for len(remaining) > 0 {
		if len(remaining) <= maxChars {
			chunks = append(chunks, remaining)
			break
		}

		limit := maxChars
		if limit > len(remaining) {
			limit = len(remaining)
		}
		
		splitIdx := findSplitPoint(remaining[:limit])
		if splitIdx <= 0 {
			splitIdx = limit
		}

		chunks = append(chunks, remaining[:splitIdx])
		remaining = strings.TrimLeft(remaining[splitIdx:], " \t\n\r")
		
		if len(remaining) > 0 && strings.TrimSpace(remaining) == "" {
			chunks[len(chunks)-1] += remaining
			break
		}
	}

	return chunks
}

func findSplitPoint(s string) int {
	if idx := strings.LastIndex(s, "\n\n"); idx != -1 {
		return idx + 2
	}
	if idx := strings.LastIndex(s, "\n"); idx != -1 {
		return idx + 1
	}
	terminators := []string{". ", "? ", "! ", "; "}
	bestIdx := -1
	for _, term := range terminators {
		if idx := strings.LastIndex(s, term); idx > bestIdx {
			bestIdx = idx + len(term)
		}
	}
	if bestIdx != -1 {
		return bestIdx
	}
	if idx := strings.LastIndex(s, " "); idx != -1 {
		return idx + 1
	}
	return -1
}
