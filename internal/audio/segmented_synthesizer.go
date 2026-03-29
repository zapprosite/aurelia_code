package audio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"log/slog"
	"strings"
	"unicode"

	"github.com/kocar/aurelia/internal/agent"
)

// SegmentedSynthesizer fraciona o stream de texto em segmentos para o Kokoro TTS.
type SegmentedSynthesizer struct {
	kokoroURL string
	voice     string
	speed     float64
	logger    *slog.Logger
}

func NewSegmentedSynthesizer(url, voice string, speed float64) *SegmentedSynthesizer {
	if url == "" {
		url = "http://localhost:8012"
	}
	if voice == "" {
		voice = "pt-br_isabela"
	}
	if speed == 0 {
		speed = 1.0
	}
	return &SegmentedSynthesizer{
		kokoroURL: url,
		voice:     voice,
		speed:     speed,
		logger:    slog.Default().With("actor", "weaver", "component", "sap"),
	}
}

func (s *SegmentedSynthesizer) Name() string {
	return "weaver"
}

// Weave consome o stream de tokens e emite chunks de áudio (opus/mp3).
func (s *SegmentedSynthesizer) Weave(ctx context.Context, tokenStream <-chan agent.StreamResponse) (<-chan []byte, <-chan error) {
	audioCh := make(chan []byte, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(audioCh)
		defer close(errCh)

		s.logger.Info("Starting audio weaving process")
		var currentSegment strings.Builder
		for resp := range tokenStream {
			if resp.Err != nil {
				s.logger.Error("Token stream error", "err", resp.Err)
				errCh <- resp.Err
				return
			}
			if resp.Content == "" {
				continue
			}

			currentSegment.WriteString(resp.Content)
			text := currentSegment.String()

			// Verifica se o segmento terminou de forma natural (pontuação)
			if isEndOfSegment(text) {
				s.logger.Debug("Synthesizing segment", "length", len(text))
				audio, err := s.synthesize(ctx, text)
				if err != nil {
					s.logger.Error("Synthesize failure", "err", err)
					errCh <- fmt.Errorf("synthesize error: %w", err)
				} else if len(audio) > 0 {
					audioCh <- audio
				}
				currentSegment.Reset()
			}
		}

		// Processa o resto se sobrar algo
		if currentSegment.Len() > 0 {
			audio, err := s.synthesize(ctx, currentSegment.String())
			if err != nil {
				errCh <- err
			} else if len(audio) > 0 {
				audioCh <- audio
			}
		}
		s.logger.Info("Audio weaving completed")
	}()

	return audioCh, errCh
}

func isEndOfSegment(text string) bool {
	trimmed := strings.TrimSpace(text)
	if len(trimmed) < 15 { // Evita sintetizar palavras soltas muito curtas
		return false
	}

	lastChar := rune(text[len(text)-1])
	// Se terminou com espaço, verifica a pontuação anterior
	if unicode.IsSpace(lastChar) {
		last := rune(trimmed[len(trimmed)-1])
		return isPunctuation(last) || last == '\n'
	}

	// Se for uma sentença longa (+ de 100 chars), força o split no próximo espaço
	if len(text) > 100 && unicode.IsSpace(lastChar) {
		return true
	}

	return lastChar == '\n'
}

func isPunctuation(r rune) bool {
	return r == '.' || r == '?' || r == '!' || r == ':' || r == ';'
}

// analyzeEmotion detecta o tom emocional do texto para ajustar o Kokoro.
func analyzeEmotion(text string) (speed float64, voice string) {
	lower := strings.ToLower(text)
	
	// Default SOTA 2026: Isabela (Serene)
	speed = 1.0
	voice = "pt-br_isabela"

	// Emoções e Ajustes (SOTA 2026.2)
	emotions := map[string]struct {
		speed float64
		voice string
	}{
		"perigo": {1.2, "pt-br_isabela"},
		"alerta": {1.15, "pt-br_isabela"},
		"urgência": {1.2, "pt-br_isabela"},
		"sucesso": {1.05, "pt-br_isabela"},
		"concluído": {1.05, "pt-br_isabela"},
		"desculpe": {0.9, "pt-br_isabela"},
		"lamento": {0.9, "pt-br_isabela"},
		"calma": {0.85, "pt-br_isabela"},
	}

	for key, val := range emotions {
		if strings.Contains(lower, key) {
			return val.speed, val.voice
		}
	}

	return speed, voice
}

func (s *SegmentedSynthesizer) synthesize(ctx context.Context, text string) ([]byte, error) {
	// Protocolo de Perfeição Portuguesa (Aurelia Serene Audio)
	// 1. Tail Padding para decaimento natural
	paddedText := strings.TrimSpace(text) + " . . . . . "

	// 2. Análise de Prosódia (Slice 3: Serene Prosody)
	speed, voice := analyzeEmotion(text)

	payload := map[string]interface{}{
		"model":           "kokoro",
		"input":           paddedText,
		"voice":           voice,
		"response_format": "opus",
		"speed":           speed,
		"lang_code":       "pt-br", // Ativa espeak-ng nativo para PT-BR
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.kokoroURL+"/v1/audio/speech", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kokoro error %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
