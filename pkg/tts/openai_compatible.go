package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Synthesizer interface {
	Synthesize(ctx context.Context, text string) (Audio, error)
	IsAvailable() bool
	MaxChars() int
}

type Audio struct {
	Data        []byte
	ContentType string
	Extension   string
	AsVoiceNote bool
}

type OpenAICompatibleSynthesizer struct {
	baseURL    string
	model      string
	voice      string
	language   string
	format     string
	speed      float64
	httpClient *http.Client
}

func NewOpenAICompatibleSynthesizer(baseURL, model, voice, language, format string, speed float64) *OpenAICompatibleSynthesizer {
	if strings.TrimSpace(format) == "" {
		format = "opus"
	}

	if speed <= 0 {
		speed = 1.0
	}
	return &OpenAICompatibleSynthesizer{
		baseURL:  strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		model:    strings.TrimSpace(model),
		voice:    normalizeVoiceAlias(model, voice),
		language: strings.TrimSpace(language),
		format:   strings.TrimSpace(format),
		speed:    speed,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func normalizeVoiceAlias(model, voice string) string {
	normalizedVoice := strings.TrimSpace(voice)
	if normalizedVoice == "" {
		return ""
	}

	switch strings.ToLower(normalizedVoice) {
	case "pt-br", "pt_br", "ptbr":
		// Preserve the historic app-level alias while translating it to a real
		// Kokoro feminine PT voice that the current image exposes.
		switch strings.ToLower(strings.TrimSpace(model)) {
		case "kokoro", "tts-1", "tts-1-hd":
			return "pf_dora"
		}
	}

	return normalizedVoice
}

func (s *OpenAICompatibleSynthesizer) IsAvailable() bool {
	return s != nil && s.baseURL != "" && s.model != "" && s.voice != ""
}

func (s *OpenAICompatibleSynthesizer) MaxChars() int {
	if s == nil {
		return 3000
	}
	// Kokoro (GPU) handles long text with internal chunking.
	if strings.Contains(strings.ToLower(s.model), "kokoro") {
		return 50000
	}
	// Default safety limit for external APIs (OpenAI standard)
	return 3000
}

func (s *OpenAICompatibleSynthesizer) Synthesize(ctx context.Context, text string) (Audio, error) {
	if !s.IsAvailable() {
		return Audio{}, fmt.Errorf("tts synthesizer is not configured")
	}
	payload := map[string]any{
		"model":           s.model,
		"input":           strings.TrimSpace(text),
		"voice":           s.voice,
		"language":        s.language,
		"response_format": s.format,
		"speed":           s.speed,
	}
	if s.language == "" {
		delete(payload, "language")
	}
	body, err := json.Marshal(payload)

	if err != nil {
		return Audio{}, fmt.Errorf("encode tts payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/v1/audio/speech", bytes.NewReader(body))
	if err != nil {
		return Audio{}, fmt.Errorf("build tts request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return Audio{}, fmt.Errorf("request tts audio: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Audio{}, fmt.Errorf("read tts response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Audio{}, fmt.Errorf("tts api error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	audio := Audio{
		Data:        respBody,
		ContentType: resp.Header.Get("Content-Type"),
		Extension:   extensionForFormat(s.format),
		AsVoiceNote: strings.EqualFold(s.format, "opus"),
	}
	if audio.Extension == "" {
		audio.Extension = extensionFromContentType(audio.ContentType)
	}
	if audio.Extension == "" {
		audio.Extension = ".bin"
	}
	return audio, nil
}

func extensionForFormat(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "opus":
		return ".ogg"
	case "mp3":
		return ".mp3"
	case "wav":
		return ".wav"
	default:
		return ""
	}
}

func extensionFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "audio/opus"), strings.Contains(contentType, "audio/ogg"):
		return ".ogg"
	case strings.Contains(contentType, "audio/mpeg"), strings.Contains(contentType, "audio/mp3"):
		return ".mp3"
	case strings.Contains(contentType, "audio/wav"), strings.Contains(contentType, "audio/x-wav"):
		return ".wav"
	default:
		return ""
	}
}
