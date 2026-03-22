package voice

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

type mirrorFunc func(context.Context, TranscriptEvent) error

func (fn mirrorFunc) MirrorTranscript(ctx context.Context, event TranscriptEvent) error {
	return fn(ctx, event)
}
