package tts

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiniMaxSynthesizer_IsAvailable(t *testing.T) {
	t.Parallel()

	if NewMiniMaxSynthesizer("", "", "speech-2.8-hd", "aurelia-voice", "mp3", 1).IsAvailable() {
		t.Fatal("expected unavailable minimax synthesizer without api key")
	}
	if !NewMiniMaxSynthesizer("", "secret", "speech-2.8-hd", "aurelia-voice", "mp3", 1).IsAvailable() {
		t.Fatal("expected configured minimax synthesizer to be available")
	}
}

func TestMiniMaxSynthesizer_SynthesizeSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/t2a_v2" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if body["model"] != "speech-2.8-hd" {
			t.Fatalf("unexpected model: %v", body["model"])
		}
		if body["language_boost"] != "Portuguese" {
			t.Fatalf("unexpected language_boost: %v", body["language_boost"])
		}
		voice, _ := body["voice_setting"].(map[string]any)
		if voice["voice_id"] != "aurelia-ptbr-formal-doce-v1" {
			t.Fatalf("unexpected voice_id: %v", voice["voice_id"])
		}
		audioBytes := []byte("mp3-bytes")
		resp := map[string]any{
			"data": map[string]any{
				"audio": hex.EncodeToString(audioBytes),
			},
			"base_resp": map[string]any{
				"status_code": 0,
				"status_msg":  "success",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	synth := NewMiniMaxSynthesizer(server.URL, "secret", "speech-2.8-hd", "aurelia-ptbr-formal-doce-v1", "mp3", 1)
	audio, err := synth.Synthesize(context.Background(), "Olá, eu sou a Aurélia.")
	if err != nil {
		t.Fatalf("Synthesize() error = %v", err)
	}
	if string(audio.Data) != "mp3-bytes" {
		t.Fatalf("unexpected audio bytes: %q", string(audio.Data))
	}
	if audio.AsVoiceNote {
		t.Fatal("expected mp3 response not to be treated as voice note")
	}
	if audio.Extension != ".mp3" {
		t.Fatalf("extension = %q, want .mp3", audio.Extension)
	}
}

func TestMiniMaxSynthesizer_SynthesizeAPIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "broken", http.StatusBadGateway)
	}))
	defer server.Close()

	synth := NewMiniMaxSynthesizer(server.URL, "secret", "speech-2.8-hd", "aurelia-voice", "mp3", 1)
	if _, err := synth.Synthesize(context.Background(), "Olá"); err == nil {
		t.Fatal("expected api error")
	}
}

func TestMiniMaxSynthesizer_SynthesizeBaseRespError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"audio": ""},
			"base_resp": map[string]any{
				"status_code": 1004,
				"status_msg":  "voice not found",
			},
		})
	}))
	defer server.Close()

	synth := NewMiniMaxSynthesizer(server.URL, "secret", "speech-2.8-hd", "missing-voice", "mp3", 1)
	if _, err := synth.Synthesize(context.Background(), "Olá"); err == nil {
		t.Fatal("expected base_resp error")
	}
}
