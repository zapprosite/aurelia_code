package tts

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultGeminiTTSBaseURL = "https://generativelanguage.googleapis.com"

type GeminiSynthesizer struct {
	baseURL    string
	apiKey     string
	model      string
	voice      string
	httpClient *http.Client
}

func NewGeminiSynthesizer(baseURL, apiKey, model, voice string) *GeminiSynthesizer {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = defaultGeminiTTSBaseURL
	}
	return &GeminiSynthesizer{
		baseURL: baseURL,
		apiKey:  strings.TrimSpace(apiKey),
		model:   strings.TrimSpace(model),
		voice:   strings.TrimSpace(voice),
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (s *GeminiSynthesizer) IsAvailable() bool {
	return s != nil && s.apiKey != "" && s.model != "" && s.voice != ""
}

func (s *GeminiSynthesizer) Synthesize(ctx context.Context, text string) (Audio, error) {
	if !s.IsAvailable() {
		return Audio{}, fmt.Errorf("gemini tts synthesizer is not configured")
	}

	payload := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{
						"text": geminiSpeechPrompt(text),
					},
				},
			},
		},
		"generationConfig": map[string]any{
			"responseModalities": []string{"AUDIO"},
			"speechConfig": map[string]any{
				"voiceConfig": map[string]any{
					"prebuiltVoiceConfig": map[string]any{
						"voiceName": s.voice,
					},
				},
			},
		},
		"model": s.model,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Audio{}, fmt.Errorf("encode gemini tts payload: %w", err)
	}

	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent", s.baseURL, s.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Audio{}, fmt.Errorf("build gemini tts request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return Audio{}, fmt.Errorf("request gemini tts audio: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var decoded struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					InlineData struct {
						Data string `json:"data"`
					} `json:"inlineData"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return Audio{}, fmt.Errorf("decode gemini tts response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if decoded.Error != nil && strings.TrimSpace(decoded.Error.Message) != "" {
			return Audio{}, fmt.Errorf("gemini tts api error (status %d): %s", resp.StatusCode, strings.TrimSpace(decoded.Error.Message))
		}
		return Audio{}, fmt.Errorf("gemini tts api error (status %d)", resp.StatusCode)
	}
	if len(decoded.Candidates) == 0 || len(decoded.Candidates[0].Content.Parts) == 0 {
		return Audio{}, fmt.Errorf("gemini tts response missing audio candidate")
	}
	pcmBase64 := decoded.Candidates[0].Content.Parts[0].InlineData.Data
	pcmBytes, err := base64.StdEncoding.DecodeString(pcmBase64)
	if err != nil {
		return Audio{}, fmt.Errorf("decode gemini tts pcm: %w", err)
	}
	wavBytes, err := wrapPCMAsWAV(pcmBytes, 1, 24000, 16)
	if err != nil {
		return Audio{}, fmt.Errorf("wrap gemini tts wav: %w", err)
	}
	return Audio{
		Data:        wavBytes,
		ContentType: "audio/wav",
		Extension:   ".wav",
		AsVoiceNote: false,
	}, nil
}

func geminiSpeechPrompt(text string) string {
	spoken := strings.TrimSpace(text)
	if spoken == "" {
		return ""
	}
	return "Fale em português do Brasil com voz feminina, tom doce, calmo e acolhedor, dicção clara e elegante, sem gírias, sem regionalismos informais e sem portunhol. Ritmo pausado, equilibrado e profissional. Texto: " + spoken
}

func wrapPCMAsWAV(pcm []byte, channels, sampleRate, bitsPerSample int) ([]byte, error) {
	if channels <= 0 || sampleRate <= 0 || bitsPerSample <= 0 {
		return nil, fmt.Errorf("invalid wav parameters")
	}
	byteRate := sampleRate * channels * bitsPerSample / 8
	blockAlign := channels * bitsPerSample / 8
	dataLen := uint32(len(pcm))

	buf := bytes.NewBuffer(make([]byte, 0, 44+len(pcm)))
	writeString := func(value string) error {
		_, err := buf.WriteString(value)
		return err
	}
	writeLE := func(value any) error {
		return binary.Write(buf, binary.LittleEndian, value)
	}

	if err := writeString("RIFF"); err != nil {
		return nil, err
	}
	if err := writeLE(uint32(36) + dataLen); err != nil {
		return nil, err
	}
	if err := writeString("WAVE"); err != nil {
		return nil, err
	}
	if err := writeString("fmt "); err != nil {
		return nil, err
	}
	if err := writeLE(uint32(16)); err != nil {
		return nil, err
	}
	if err := writeLE(uint16(1)); err != nil {
		return nil, err
	}
	if err := writeLE(uint16(channels)); err != nil {
		return nil, err
	}
	if err := writeLE(uint32(sampleRate)); err != nil {
		return nil, err
	}
	if err := writeLE(uint32(byteRate)); err != nil {
		return nil, err
	}
	if err := writeLE(uint16(blockAlign)); err != nil {
		return nil, err
	}
	if err := writeLE(uint16(bitsPerSample)); err != nil {
		return nil, err
	}
	if err := writeString("data"); err != nil {
		return nil, err
	}
	if err := writeLE(dataLen); err != nil {
		return nil, err
	}
	if _, err := buf.Write(pcm); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
