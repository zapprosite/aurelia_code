package tts

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
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

	slog.Debug("tts segmented synthesis starting", slog.Int("chunks", len(chunks)))

	for i, chunk := range chunks {
		select {
		case <-ctx.Done():
			return Audio{}, ctx.Err()
		default:
			audio, err := s.base.Synthesize(ctx, chunk)
			if err != nil {
				return Audio{}, fmt.Errorf("synthesize chunk %d/%d: %w", i+1, len(chunks), err)
			}

			dataToAdd := audio.Data
			// For WAV chunks after the first: strip the header and append only the raw PCM
			// data. The first chunk's header will be corrected at the end to reflect the
			// combined file size — without this, players stop at the first chunk's declared size.
			if i > 0 && isWAV(dataToAdd) {
				dataToAdd = dataToAdd[wavDataOffset(dataToAdd):]
			}

			combinedData = append(combinedData, dataToAdd...)
			lastAudio = audio
		}
	}

	// Update the RIFF/data size fields so the combined file is a valid single WAV.
	if isWAV(combinedData) {
		fixWAVHeader(combinedData)
	}

	slog.Debug("tts segmented synthesis complete", slog.Int("bytes", len(combinedData)))

	return Audio{
		Data:        combinedData,
		ContentType: lastAudio.ContentType,
		Extension:   lastAudio.Extension,
		AsVoiceNote: lastAudio.AsVoiceNote,
	}, nil
}

// WAV RIFF layout constants.
const (
	wavRIFFHeaderSize     = 12 // "RIFF" + 4-byte size + "WAVE"
	wavSubChunkHeaderSize = 8  // 4-byte chunk ID + 4-byte chunk size
	wavFallbackDataOffset = 44 // standard PCM WAV header size
)

// isWAV returns true if data starts with a RIFF/WAVE header.
func isWAV(data []byte) bool {
	return len(data) >= wavRIFFHeaderSize &&
		bytes.Equal(data[0:4], []byte("RIFF")) &&
		bytes.Equal(data[8:12], []byte("WAVE"))
}

// wavDataOffset returns the byte offset where raw PCM data begins (past the "data" sub-chunk header).
// Walks sub-chunks because the fmt chunk size can vary (e.g. when extensible format is used).
func wavDataOffset(data []byte) int {
	offset := wavRIFFHeaderSize
	for offset+wavSubChunkHeaderSize <= len(data) {
		id, size := readRIFFSubChunk(data, offset)
		if id == "data" {
			return offset + wavSubChunkHeaderSize
		}
		offset += wavSubChunkHeaderSize + size
		if size%2 != 0 {
			offset++ // RIFF pads odd-sized chunks to even boundaries
		}
	}
	return wavFallbackDataOffset
}

// fixWAVHeader rewrites the RIFF file-size and "data" chunk-size fields to match
// the actual buffer length after multiple WAV chunks have been concatenated.
func fixWAVHeader(data []byte) {
	if len(data) < wavFallbackDataOffset {
		return
	}
	putU32LE(data[4:], uint32(len(data)-8)) // RIFF size excludes the 8-byte RIFF header itself

	offset := wavRIFFHeaderSize
	for offset+wavSubChunkHeaderSize <= len(data) {
		id, size := readRIFFSubChunk(data, offset)
		if id == "data" {
			putU32LE(data[offset+4:], uint32(len(data)-offset-wavSubChunkHeaderSize))
			return
		}
		offset += wavSubChunkHeaderSize + size
		if size%2 != 0 {
			offset++
		}
	}
}

// readRIFFSubChunk reads a sub-chunk ID and size at the given offset (little-endian).
func readRIFFSubChunk(data []byte, offset int) (id string, size int) {
	id = string(data[offset : offset+4])
	size = int(data[offset+4]) |
		int(data[offset+5])<<8 |
		int(data[offset+6])<<16 |
		int(data[offset+7])<<24
	return
}

// putU32LE writes v as a 4-byte little-endian value into b.
func putU32LE(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
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
	// Prioridade 1: Parágrafos
	if idx := strings.LastIndex(s, "\n\n"); idx != -1 && idx > len(s)/2 {
		return idx + 2
	}
	// Prioridade 2: Quebra de linha simples
	if idx := strings.LastIndex(s, "\n"); idx != -1 && idx > len(s)/2 {
		return idx + 1
	}
	// Prioridade 3: Pontuação forte — busca no último terço para evitar chunks pequenos.
	terminators := []string{". ", "? ", "! ", "; "}
	bestIdx := -1
	for _, term := range terminators {
		if idx := strings.LastIndex(s, term); idx > bestIdx && idx > len(s)/3 {
			bestIdx = idx + len(term)
		}
	}
	if bestIdx != -1 {
		return bestIdx
	}
	// Fallback 1: Vírgula ou dois pontos
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
