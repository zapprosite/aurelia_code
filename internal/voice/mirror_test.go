package voice

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestQdrantMirror_MirrorTranscriptEmbedsAndUpserts(t *testing.T) {
	t.Parallel()

	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		switch r.URL.Path {
		case "/embed":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"embeddings":[[0.1,0.2,0.3]]}`))
		case "/collections/test_collection":
			w.WriteHeader(http.StatusOK)
		case "/collections/test_collection/points":
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	mirror := NewQdrantMirror(server.URL, "secret", "test_collection", "bge-m3", server.URL)
	mirror.embedURL = server.URL + "/embed"
	if err := mirror.MirrorTranscript(context.Background(), TranscriptEvent{
		JobID:      "job-1",
		Transcript: "ola mundo",
	}); err != nil {
		t.Fatalf("MirrorTranscript() error = %v", err)
	}
	if len(paths) != 3 {
		t.Fatalf("paths = %v", paths)
	}
}

func TestMultiMirror_StopsOnFirstError(t *testing.T) {
	t.Parallel()

	errMirror := mirrorFunc(func(ctx context.Context, event TranscriptEvent) error {
		return context.DeadlineExceeded
	})
	secondCalled := false
	secondMirror := mirrorFunc(func(ctx context.Context, event TranscriptEvent) error {
		secondCalled = true
		return nil
	})

	err := NewMultiMirror(errMirror, secondMirror).MirrorTranscript(context.Background(), TranscriptEvent{Transcript: "ola"})
	if err == nil {
		t.Fatal("expected MirrorTranscript() to fail")
	}
	if secondCalled {
		t.Fatal("expected second mirror not to be called")
	}
}

func TestQdrantMirror_SkipsEmptyTranscript(t *testing.T) {
	t.Parallel()

	mirror := NewQdrantMirror("http://example.test", "secret", "collection", "bge-m3", "http://example.test")
	if err := mirror.MirrorTranscript(context.Background(), TranscriptEvent{}); err != nil {
		t.Fatalf("MirrorTranscript() error = %v", err)
	}
}

func TestQdrantMirror_UpsertShapeContainsVector(t *testing.T) {
	t.Parallel()

	var pointsBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/embed":
			_, _ = w.Write([]byte(`{"embeddings":[[0.1,0.2]]}`))
		case "/collections/test_collection":
			w.WriteHeader(http.StatusOK)
		case "/collections/test_collection/points":
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			data, _ := json.Marshal(payload)
			pointsBody = string(data)
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	mirror := NewQdrantMirror(server.URL, "secret", "test_collection", "bge-m3", server.URL)
	mirror.embedURL = server.URL + "/embed"
	if err := mirror.MirrorTranscript(context.Background(), TranscriptEvent{JobID: "job-1", Transcript: "ola"}); err != nil {
		t.Fatalf("MirrorTranscript() error = %v", err)
	}
	if !strings.Contains(pointsBody, "\"vector\":[0.1,0.2]") {
		t.Fatalf("points body = %s", pointsBody)
	}
}

func TestBuildVoiceSemanticPayload_AddsCanonicalFieldsAndKeepsLegacyFields(t *testing.T) {
	t.Parallel()

	event := TranscriptEvent{
		JobID:       "job-42",
		UserID:      11,
		ChatID:      22,
		Source:      "groq-whisper",
		Transcript:  "ligar o modo sentinela",
		Accepted:    true,
		RequiresTTS: true,
		CreatedAt:   time.Date(2026, 3, 25, 21, 15, 30, 0, time.UTC),
	}

	payload := buildVoiceSemanticPayload(event)

	if payload["text"] != event.Transcript {
		t.Fatalf("expected text=%q, got %#v", event.Transcript, payload["text"])
	}
	if payload["canonical_bot_id"] != defaultCanonicalBotID {
		t.Fatalf("expected canonical_bot_id=%q, got %#v", defaultCanonicalBotID, payload["canonical_bot_id"])
	}
	if payload["source_system"] != defaultVoiceSourceSystem {
		t.Fatalf("expected source_system=%q, got %#v", defaultVoiceSourceSystem, payload["source_system"])
	}
	if payload["domain"] != defaultVoiceDomain {
		t.Fatalf("expected domain=%q, got %#v", defaultVoiceDomain, payload["domain"])
	}
	if payload["version"] != defaultPayloadVersion {
		t.Fatalf("expected version=%d, got %#v", defaultPayloadVersion, payload["version"])
	}
	if payload["source_id"] != "voice:job-42" {
		t.Fatalf("expected source_id=%q, got %#v", "voice:job-42", payload["source_id"])
	}
	if payload["ts"] != event.CreatedAt.Unix() {
		t.Fatalf("expected ts=%d, got %#v", event.CreatedAt.Unix(), payload["ts"])
	}

	if payload["transcript"] != event.Transcript {
		t.Fatalf("expected transcript legacy field to remain, got %#v", payload["transcript"])
	}
	if payload["source"] != event.Source {
		t.Fatalf("expected source legacy field to remain, got %#v", payload["source"])
	}
	if payload["created_at"] != event.CreatedAt.Format(time.RFC3339Nano) {
		t.Fatalf("expected created_at legacy field to remain, got %#v", payload["created_at"])
	}
}

func TestBuildVoiceSemanticPayload_FallsBackToTimestampSourceID(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 3, 25, 22, 0, 0, 0, time.UTC)
	payload := buildVoiceSemanticPayload(TranscriptEvent{
		Transcript: "modo sentinela",
		CreatedAt:  createdAt,
	})

	if payload["source_id"] != "voice:ts:1774476000" {
		t.Fatalf("unexpected fallback source_id: %#v", payload["source_id"])
	}
}

type mirrorFunc func(context.Context, TranscriptEvent) error

func (fn mirrorFunc) MirrorTranscript(ctx context.Context, event TranscriptEvent) error {
	return fn(ctx, event)
}
