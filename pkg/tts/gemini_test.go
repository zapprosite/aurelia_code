package tts

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeminiSynthesizer_IsAvailable(t *testing.T) {
	t.Parallel()

	if NewGeminiSynthesizer("", "", "gemini-2.5-flash-preview-tts", "Sulafat").IsAvailable() {
		t.Fatal("expected unavailable gemini synthesizer without api key")
	}
	if !NewGeminiSynthesizer("", "secret", "gemini-2.5-flash-preview-tts", "Sulafat").IsAvailable() {
		t.Fatal("expected configured gemini synthesizer to be available")
	}
}

func TestGeminiSynthesizer_SynthesizeSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/models/gemini-2.5-flash-preview-tts:generateContent" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("x-goog-api-key"); got != "secret" {
			t.Fatalf("unexpected api key header: %q", got)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if body["model"] != "gemini-2.5-flash-preview-tts" {
			t.Fatalf("unexpected model: %v", body["model"])
		}
		gen, _ := body["generationConfig"].(map[string]any)
		voiceCfg := gen["speechConfig"].(map[string]any)["voiceConfig"].(map[string]any)["prebuiltVoiceConfig"].(map[string]any)
		if voiceCfg["voiceName"] != "Sulafat" {
			t.Fatalf("unexpected voice name: %v", voiceCfg["voiceName"])
		}

		pcm := []byte{0x01, 0x02, 0x03, 0x04}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"candidates": []map[string]any{
				{
					"content": map[string]any{
						"parts": []map[string]any{
							{
								"inlineData": map[string]any{
									"data": base64.StdEncoding.EncodeToString(pcm),
								},
							},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	synth := NewGeminiSynthesizer(server.URL, "secret", "gemini-2.5-flash-preview-tts", "Sulafat")
	audio, err := synth.Synthesize(context.Background(), "Olá")
	if err != nil {
		t.Fatalf("Synthesize() error = %v", err)
	}
	if audio.Extension != ".wav" {
		t.Fatalf("extension = %q, want .wav", audio.Extension)
	}
	if audio.ContentType != "audio/wav" {
		t.Fatalf("content type = %q, want audio/wav", audio.ContentType)
	}
	if audio.AsVoiceNote {
		t.Fatal("expected wav output not to be treated as voice note")
	}
	if len(audio.Data) <= 44 {
		t.Fatalf("expected wav payload larger than header, got %d bytes", len(audio.Data))
	}
}

func TestGeminiSynthesizer_SynthesizeAPIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"message": "bad request"},
		})
	}))
	defer server.Close()

	synth := NewGeminiSynthesizer(server.URL, "secret", "gemini-2.5-flash-preview-tts", "Sulafat")
	if _, err := synth.Synthesize(context.Background(), "Olá"); err == nil {
		t.Fatal("expected api error")
	}
}

func TestWrapPCMAsWAV(t *testing.T) {
	t.Parallel()

	got, err := wrapPCMAsWAV([]byte{0x01, 0x02}, 1, 24000, 16)
	if err != nil {
		t.Fatalf("wrapPCMAsWAV() error = %v", err)
	}
	if string(got[:4]) != "RIFF" {
		t.Fatalf("missing RIFF header")
	}
	if string(got[8:12]) != "WAVE" {
		t.Fatalf("missing WAVE header")
	}
}
