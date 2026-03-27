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

	if NewOpenAICompatibleSynthesizer("", "tts-1-hd", "aurelia", "", "opus", 1).IsAvailable() {
		t.Fatal("expected unavailable synthesizer without base url")
	}
	if !NewOpenAICompatibleSynthesizer("http://127.0.0.1:8011", "tts-1-hd", "aurelia", "", "opus", 1).IsAvailable() {
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
		if body["model"] != "tts-1-hd" {
			t.Fatalf("unexpected model: %v", body["model"])
		}
		if body["voice"] != "aurelia" {
			t.Fatalf("unexpected voice: %v", body["voice"])
		}
		w.Header().Set("Content-Type", "audio/opus")
		_, _ = w.Write([]byte("opus-bytes"))
	}))
	defer server.Close()

	synth := NewOpenAICompatibleSynthesizer(server.URL, "tts-1-hd", "aurelia", "", "opus", 1)
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

	synth := NewOpenAICompatibleSynthesizer(server.URL, "tts-1-hd", "aurelia", "", "opus", 1)
	if _, err := synth.Synthesize(context.Background(), "Ola"); err == nil {
		t.Fatal("expected api error")
	}
}

func TestOpenAICompatibleSynthesizer_NormalizesLegacyPTBRAlias(t *testing.T) {
	t.Parallel()

	synth := NewOpenAICompatibleSynthesizer("http://127.0.0.1:8012", "kokoro", "pt-br", "pt", "opus", 1)
	if synth.voice != "pf_dora" {
		t.Fatalf("voice = %q, want pf_dora", synth.voice)
	}
}
func TestOpenAICompatibleSynthesizer_MaxChars(t *testing.T) {
	t.Parallel()

	tests := []struct {
		model    string
		expected int
	}{
		{"kokoro", 50000},
		{"kodoro", 50000},
		{"KOKORO-v1", 50000},
		{"tts-1", 3000},
		{"", 3000},
	}

	for _, tc := range tests {
		s := NewOpenAICompatibleSynthesizer("http://local", tc.model, "pf_dora", "pt", "opus", 1.0)
		if got := s.MaxChars(); got != tc.expected {
			t.Errorf("MaxChars(%q) = %d, want %d", tc.model, got, tc.expected)
		}
	}
}
