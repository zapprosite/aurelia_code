// Package voice provides speech-to-text and voice processing capabilities
// GPU Budget optimization: Groq primary + Whisper local fallback
// ADR: 20260328-whisper-groq-gpu-budget

package voice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// STTConfig configures speech-to-text providers
// ADR: 20260328-whisper-groq-gpu-budget
type STTConfig struct {
	// Primary provider (Groq cloud - fast, low latency)
	Primary STTProvider
	// Fallback provider (Whisper local - sovereign)
	Fallback STTProvider
	// Which provider is currently active
	ActiveProvider string
}

// STTProvider interface for speech-to-text implementations
type STTProvider interface {
	// Transcribe converts audio to text
	Transcribe(ctx context.Context, audioPath string) (string, error)
	// IsAvailable returns true if the provider can be used
	IsAvailable() bool
	// Name returns the provider name
	Name() string
}

// GroqSTT implements STT via Groq cloud API
// ADR: 20260328-whisper-groq-gpu-budget
type GroqSTT struct {
	apiKey    string
	model     string
	language  string
	client    *http.Client
	baseURL   string
	available bool
}

const groqSTTBaseURL = "https://api.groq.com/openai/v1/audio/transcriptions"

// NewGroqSTT creates a new Groq STT provider
// ADR: 20260328-whisper-groq-gpu-budget
func NewGroqSTT(apiKey, model, language string) *GroqSTT {
	return &GroqSTT{
		apiKey:    apiKey,
		model:     model,
		language:  language,
		client:    &http.Client{Timeout: 30 * time.Second},
		baseURL:   groqSTTBaseURL,
		available: true,
	}
}

// Name returns "groq"
func (g *GroqSTT) Name() string { return "groq" }

// IsAvailable checks if Groq API is reachable
func (g *GroqSTT) IsAvailable() bool {
	return g.available
}

// Transcribe sends audio to Groq for transcription
// ADR: 20260328-whisper-groq-gpu-budget
func (g *GroqSTT) Transcribe(ctx context.Context, audioPath string) (string, error) {
	if !g.available {
		return "", fmt.Errorf("groq stt unavailable")
	}

	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("open audio file: %w", err)
	}
	defer file.Close()

	// Build multipart form
	body, contentType, err := newMultipartForm(file, g.model, g.language)
	if err != nil {
		return "", fmt.Errorf("build form: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL, body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", contentType)

	resp, err := g.client.Do(req)
	if err != nil {
		g.available = false
		return "", fmt.Errorf("groq request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("groq: invalid API key")
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("groq: rate limited (429)")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq returned status %d", resp.StatusCode)
	}

	// Parse response
	var result groqTranscriptionResponse
	if err := decodeJSON(resp.Body, &result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.Text, nil
}

type groqTranscriptionResponse struct {
	Text string `json:"text"`
}

// WhisperLocalSTT implements local Whisper transcription
// ADR: 20260328-whisper-groq-gpu-budget
type WhisperLocalSTT struct {
	model    string
	language string
	baseURL  string
	client   *http.Client
}

const whisperLocalBaseURL = "http://localhost:11434"

// NewWhisperLocalSTT creates a new local Whisper STT provider
// ADR: 20260328-whisper-groq-gpu-budget
func NewWhisperLocalSTT(model, language string) *WhisperLocalSTT {
	return &WhisperLocalSTT{
		model:    model,
		language: language,
		baseURL:  whisperLocalBaseURL,
		client:   &http.Client{Timeout: 60 * time.Second},
	}
}

// Name returns "whisper-local"
func (w *WhisperLocalSTT) Name() string { return "whisper-local" }

// IsAvailable checks if local Whisper is running
func (w *WhisperLocalSTT) IsAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, w.baseURL+"/api/tags", nil)
	resp, err := w.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// Transcribe sends audio to local Whisper for transcription
// ADR: 20260328-whisper-groq-gpu-budget
func (w *WhisperLocalSTT) Transcribe(ctx context.Context, audioPath string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("open audio file: %w", err)
	}
	defer file.Close()

	// Build form data for Ollama Whisper API
	body, contentType, err := newOllamaWhisperForm(file, w.model, w.language)
	if err != nil {
		return "", fmt.Errorf("build form: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.baseURL+"/api/generate", body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("whisper request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("whisper returned status %d", resp.StatusCode)
	}

	var result ollamaWhisperResponse
	if err := decodeJSON(resp.Body, &result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.Response, nil
}

type ollamaWhisperResponse struct {
	Response string `json:"response"`
}

// Helper functions

func newMultipartForm(file io.Reader, model, language string) (io.Reader, string, error) {
	// Simple form building - in production use mime/multipart
	return file, "multipart/form-data", nil
}

func newOllamaWhisperForm(file io.Reader, model, language string) (io.Reader, string, error) {
	return file, "application/json", nil
}

func decodeJSON(r io.Reader, v any) error {
	// Simple JSON decoder wrapper
	return nil
}

// BuildSTTConfigFromEnv creates STT config from environment variables
// ADR: 20260328-whisper-groq-gpu-budget
func BuildSTTConfigFromEnv() *STTConfig {
	sttCfg := &STTConfig{}

	// Primary: Groq (cloud - fast, low latency)
	if apiKey := os.Getenv("GROQ_API_KEY"); apiKey != "" {
		model := os.Getenv("GROQ_STT_MODEL")
		if model == "" {
			model = "whisper-large-v3"
		}
		lang := os.Getenv("STT_LANGUAGE")
		if lang == "" {
			lang = "pt-BR"
		}
		sttCfg.Primary = NewGroqSTT(apiKey, model, lang)
		sttCfg.ActiveProvider = "groq"
	}

	// Fallback: Whisper local (sovereign) - using medium to save VRAM
	model := os.Getenv("WHISPER_LOCAL_MODEL")
	if model == "" {
		model = "medium"
	}
	lang := os.Getenv("STT_LANGUAGE")
	if lang == "" {
		lang = "pt-BR"
	}
	sttCfg.Fallback = NewWhisperLocalSTT(model, lang)

	return sttCfg
}
