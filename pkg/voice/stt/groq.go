package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kocar/aurelia/internal/observability"
)

const (
	defaultGroqSTTBaseURL  = "https://api.groq.com/openai/v1"
	defaultGroqSTTModel    = "whisper-large-v3-turbo"
	defaultGroqSTTLanguage = "pt"
)

// GroqTranscriber implements Transcriber using Groq's high-speed Whisper API.
type GroqTranscriber struct {
	baseURL    string
	model      string
	language   string
	apiKey     string
	httpClient *http.Client
}

// NewGroqTranscriber creates a transcriber pointed at Groq API.
func NewGroqTranscriber(baseURL, model, language, apiKey string) *GroqTranscriber {
	if baseURL == "" {
		baseURL = defaultGroqSTTBaseURL
	}
	if model == "" {
		model = defaultGroqSTTModel
	}
	if language == "" {
		language = defaultGroqSTTLanguage
	}
	return &GroqTranscriber{
		baseURL:  baseURL,
		model:    model,
		language: language,
		apiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Transcribe sends audio to Groq and returns the text.
func (t *GroqTranscriber) Transcribe(ctx context.Context, audioFilePath string) (string, error) {
	logger := observability.Logger("stt.groq")
	logger.Info("starting groq transcription", 
		slog.String("file", observability.Basename(audioFilePath)),
		slog.String("model", t.model))

	audioFile, err := os.Open(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %w", err)
	}
	defer func() { _ = audioFile.Close() }()

	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	part, err := w.CreateFormFile("file", filepath.Base(audioFilePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, audioFile); err != nil {
		return "", fmt.Errorf("failed to copy audio content: %w", err)
	}
	if err := w.WriteField("model", t.model); err != nil {
		return "", fmt.Errorf("failed to write model field: %w", err)
	}
	if err := w.WriteField("language", t.language); err != nil {
		return "", fmt.Errorf("failed to write language field: %w", err)
	}
	if err := w.WriteField("response_format", "json"); err != nil {
		return "", fmt.Errorf("failed to write response_format field: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := t.baseURL + "/audio/transcriptions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if t.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+t.apiKey)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("groq stt request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq stt error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse groq response: %w", err)
	}

	logger.Info("groq transcription completed", slog.String("text_len", fmt.Sprintf("%d chars", len(result.Text))))
	return result.Text, nil
}

// IsAvailable checks if the Groq API is reachable (basic check).
func (t *GroqTranscriber) IsAvailable() bool {
	// For cloud providers, we assume availability if API Key is present,
	// but we could do a ping to the base URL if needed.
	return t.apiKey != ""
}
