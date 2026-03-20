package tts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAICompatibleSynthesizer_IsAvailable(t *testing.T) {
	t.Parallel()

	if NewOpenAICompatibleSynthesizer("", "chatterbox", "Aurelia.wav", "", "opus", 1).IsAvailable() {
		t.Fatal("expected unavailable synthesizer without base url")
	}
	if !NewOpenAICompatibleSynthesizer("http://127.0.0.1:8011", "chatterbox", "Aurelia.wav", "", "opus", 1).IsAvailable() {
		t.Fatal("expected configured synthesizer to be available")
	}
}

func TestOpenAICompatibleSynthesizer_SynthesizeSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/audio/speech" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if body["model"] != "chatterbox" {
			t.Fatalf("unexpected model: %v", body["model"])
		}
		if body["voice"] != "Aurelia.wav" {
			t.Fatalf("unexpected voice: %v", body["voice"])
		}
		w.Header().Set("Content-Type", "audio/opus")
		_, _ = w.Write([]byte("opus-bytes"))
	}))
	defer server.Close()

	synth := NewOpenAICompatibleSynthesizer(server.URL, "chatterbox", "Aurelia.wav", "", "opus", 1)
	audio, err := synth.Synthesize(context.Background(), "Ola")
	if err != nil {
		t.Fatalf("Synthesize() error = %v", err)
	}
	if string(audio.Data) != "opus-bytes" {
		t.Fatalf("unexpected audio bytes: %q", string(audio.Data))
	}
	if !audio.AsVoiceNote {
		t.Fatal("expected opus response to be marked as voice note")
	}
	if audio.Extension != ".ogg" {
		t.Fatalf("extension = %q, want .ogg", audio.Extension)
	}
}

func TestOpenAICompatibleSynthesizer_SynthesizeAPIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "broken", http.StatusBadGateway)
	}))
	defer server.Close()

	synth := NewOpenAICompatibleSynthesizer(server.URL, "chatterbox", "Aurelia.wav", "", "opus", 1)
	if _, err := synth.Synthesize(context.Background(), "Ola"); err == nil {
		t.Fatal("expected api error")
	}
}
