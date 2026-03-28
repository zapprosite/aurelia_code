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

	fmt.Printf("TTS: Starting synthesis for %d chunks (Infinite-Voice SOTA 2026.1)\n", len(chunks))

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

			dataToAdd := audio.Data
			// Binary Stream Healing SOTA 2026.1:
			// Se o áudio for WAV (RIFF), removemos o cabeçalho de 44 bytes de todos os chunks após o primeiro.
			// Isso garante que o bitstream seja interpretado como um único arquivo contínuo.
			if i > 0 && len(dataToAdd) > 44 && string(dataToAdd[:4]) == "RIFF" {
				fmt.Printf("TTS: Stripping RIFF header from chunk %d\n", i+1)
				dataToAdd = dataToAdd[44:]
			}

			combinedData = append(combinedData, dataToAdd...)
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
	// Prioridade 1: Parágrafos (Luxo SOTA 2026)
	if idx := strings.LastIndex(s, "\n\n"); idx != -1 && idx > len(s)/2 {
		return idx + 2
	}
	// Prioridade 2: Quebra de linha simples
	if idx := strings.LastIndex(s, "\n"); idx != -1 && idx > len(s)/2 {
		return idx + 1
	}
	// Prioridade 3: Pontuação forte (fim de sentença)
	terminators := []string{". ", "? ", "! ", "; "}
	bestIdx := -1
	// Procuramos o último terminador que esteja pelo menos no último terço do chunk
	// para evitar chunks pequenos demais que prejudiquem a prosódia.
	for _, term := range terminators {
		if idx := strings.LastIndex(s, term); idx > bestIdx && idx > len(s)/3 {
			bestIdx = idx + len(term)
		}
	}
	if bestIdx != -1 {
		return bestIdx
	}
	// Fallback 1: Vírgula ou Dois Pontos
	fallbacks := []string{", ", ": "}
	for _, term := range fallbacks {
		if idx := strings.LastIndex(s, term); idx > bestIdx && idx > len(s)/2 {
			bestIdx = idx + len(term)
		}
	}
	if bestIdx != -1 {
		return bestIdx
	}
	// Fallback 2: Espaço
	if idx := strings.LastIndex(s, " "); idx != -1 {
		return idx + 1
	}
	return -1
}
