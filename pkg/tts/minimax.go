package tts

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultMiniMaxSpeechEndpoint = "https://api.minimax.io"

type MiniMaxSynthesizer struct {
	baseURL    string
	apiKey     string
	model      string
	voice      string
	format     string
	speed      float64
	httpClient *http.Client
}

func NewMiniMaxSynthesizer(baseURL, apiKey, model, voice, format string, speed float64) *MiniMaxSynthesizer {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = defaultMiniMaxSpeechEndpoint
	}
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "" {
		format = "mp3"
	}
	if speed <= 0 {
		speed = 1.0
	}
	return &MiniMaxSynthesizer{
		baseURL: baseURL,
		apiKey:  strings.TrimSpace(apiKey),
		model:   strings.TrimSpace(model),
		voice:   strings.TrimSpace(voice),
		format:  format,
		speed:   speed,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (s *MiniMaxSynthesizer) IsAvailable() bool {
	return s != nil && s.apiKey != "" && s.model != "" && s.voice != ""
}

func (s *MiniMaxSynthesizer) Synthesize(ctx context.Context, text string) (Audio, error) {
	if !s.IsAvailable() {
		return Audio{}, fmt.Errorf("minimax tts synthesizer is not configured")
	}

	payload := map[string]any{
		"model":          s.model,
		"text":           strings.TrimSpace(text),
		"stream":         false,
		"language_boost": "Portuguese",
		"output_format":  "hex",
		"voice_setting": map[string]any{
			"voice_id": s.voice,
			"speed":    s.speed,
			"vol":      1,
			"pitch":    0,
		},
		"audio_setting": map[string]any{
			"sample_rate": 32000,
			"bitrate":     128000,
			"format":      s.format,
			"channel":     1,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return Audio{}, fmt.Errorf("encode minimax tts payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/v1/t2a_v2", bytes.NewReader(body))
	if err != nil {
		return Audio{}, fmt.Errorf("build minimax tts request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return Audio{}, fmt.Errorf("request minimax tts audio: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Audio{}, fmt.Errorf("read minimax tts response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Audio{}, fmt.Errorf("minimax tts api error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var decoded struct {
		Data struct {
			Audio string `json:"audio"`
		} `json:"data"`
		BaseResp struct {
			StatusCode int    `json:"status_code"`
			StatusMsg  string `json:"status_msg"`
		} `json:"base_resp"`
	}
	if err := json.Unmarshal(respBody, &decoded); err != nil {
		return Audio{}, fmt.Errorf("decode minimax tts response: %w", err)
	}
	if decoded.BaseResp.StatusCode != 0 {
		return Audio{}, fmt.Errorf("minimax tts api error (code %d): %s", decoded.BaseResp.StatusCode, strings.TrimSpace(decoded.BaseResp.StatusMsg))
	}

	audioBytes, err := hex.DecodeString(decoded.Data.Audio)
	if err != nil {
		return Audio{}, fmt.Errorf("decode minimax tts hex audio: %w", err)
	}
	return Audio{
		Data:        audioBytes,
		ContentType: contentTypeForAudioFormat(s.format),
		Extension:   extensionForFormat(s.format),
		AsVoiceNote: strings.EqualFold(s.format, "opus"),
	}, nil
}

func contentTypeForAudioFormat(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "mp3":
		return "audio/mpeg"
	case "wav":
		return "audio/wav"
	case "flac":
		return "audio/flac"
	case "pcm":
		return "audio/L16"
	case "opus":
		return "audio/ogg"
	default:
		return "application/octet-stream"
	}
}
