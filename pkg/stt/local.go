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
	defaultLocalSTTModel    = "Systran/faster-whisper-large-v3"
	defaultLocalSTTLanguage = "pt"
)

// LocalTranscriber implements Transcriber using a local OpenAI-compatible
// Whisper server (e.g. faster-whisper-server at localhost:8020).
type LocalTranscriber struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewLocalTranscriber creates a transcriber pointed at a local Whisper API.
func NewLocalTranscriber(baseURL, model string) *LocalTranscriber {
	if baseURL == "" {
		baseURL = "http://localhost:8020"
	}
	if model == "" {
		model = defaultLocalSTTModel
	}
	return &LocalTranscriber{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Transcribe sends audio to the local Whisper server and returns the text.
func (t *LocalTranscriber) Transcribe(ctx context.Context, audioFilePath string) (string, error) {
	logger := observability.Logger("stt.local")
	logger.Info("starting local transcription", slog.String("file", observability.Basename(audioFilePath)))

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
	if err := w.WriteField("language", defaultLocalSTTLanguage); err != nil {
		return "", fmt.Errorf("failed to write language field: %w", err)
	}
	if err := w.WriteField("response_format", "json"); err != nil {
		return "", fmt.Errorf("failed to write response_format field: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := t.baseURL + "/v1/audio/transcriptions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("local whisper request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("local whisper error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	logger.Info("local transcription completed", slog.String("text_len", fmt.Sprintf("%d chars", len(result.Text))))
	return result.Text, nil
}

// IsAvailable checks if the local Whisper server is reachable.
func (t *LocalTranscriber) IsAvailable() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(t.baseURL + "/health")
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
