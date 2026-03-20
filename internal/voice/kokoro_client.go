package voice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type KokoroClient struct {
	BaseURL string // Default: http://localhost:8888
}

type KokoroRequest struct {
	Text    string  `json:"text"`
	VoiceID string  `json:"voice"`
	Speed   float64 `json:"speed"`
}

func NewKokoroClient() *KokoroClient {
	url := os.Getenv("KOKORO_URL")
	if url == "" {
		url = "http://localhost:8888"
	}
	return &KokoroClient{BaseURL: url}
}

func (k *KokoroClient) GenerateSpeech(text string, voiceID string) ([]byte, error) {
	if voiceID == "" {
		voiceID = "pf_dora" // Default profissional feminina pt-br
	}

	url := fmt.Sprintf("%s/v1/audio/speech", k.BaseURL)
	
	reqBody := map[string]interface{}{
		"model": "kokoro",
		"input": text,
		"voice": voiceID,
		"response_format": "mp3",
		"speed": 1.0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to local Kokoro GPU: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kokoro error: status %d, body %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
