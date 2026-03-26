package voice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/memory"
)

const (
	defaultQdrantMirrorCollection = "conversation_memory"
	defaultQdrantMirrorEmbedding  = "nomic-embed-text"
	defaultCanonicalBotID         = "aurelia_code"
	defaultVoiceDomain            = "system"
	defaultVoiceSourceSystem      = "voice"
	defaultPayloadVersion         = 1
)

type MultiMirror struct {
	mirrors []Mirror
}

func NewMultiMirror(mirrors ...Mirror) Mirror {
	filtered := make([]Mirror, 0, len(mirrors))
	for _, mirror := range mirrors {
		if mirror != nil {
			filtered = append(filtered, mirror)
		}
	}
	if len(filtered) == 0 {
		return noopMirror{}
	}
	return &MultiMirror{mirrors: filtered}
}

func (m *MultiMirror) MirrorTranscript(ctx context.Context, event TranscriptEvent) error {
	for _, mirror := range m.mirrors {
		if err := mirror.MirrorTranscript(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

type QdrantMirror struct {
	baseURL        string
	apiKey         string
	collection     string
	embeddingModel string
	embedURL       string
	client         *http.Client
	ensureOnce     sync.Once
	ensureErr      error
}

func NewQdrantMirror(baseURL, apiKey, collection, embeddingModel, ollamaURL string) *QdrantMirror {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	collection = strings.TrimSpace(collection)
	if collection == "" {
		collection = defaultQdrantMirrorCollection
	}
	embeddingModel = strings.TrimSpace(embeddingModel)
	if embeddingModel == "" {
		embeddingModel = defaultQdrantMirrorEmbedding
	}
	return &QdrantMirror{
		baseURL:        baseURL,
		apiKey:         strings.TrimSpace(apiKey),
		collection:     collection,
		embeddingModel: embeddingModel,
		embedURL:       strings.TrimRight(strings.TrimSpace(ollamaURL), "/") + "/api/embed",
		client:         &http.Client{Timeout: 15 * time.Second},
	}
}

func (m *QdrantMirror) MirrorTranscript(ctx context.Context, event TranscriptEvent) error {
	if m == nil || m.baseURL == "" || m.apiKey == "" || strings.TrimSpace(event.Transcript) == "" {
		return nil
	}

	vector, err := m.embed(ctx, event.Transcript)
	if err != nil {
		return err
	}
	m.ensureOnce.Do(func() {
		m.ensureErr = m.ensureCollection(ctx, len(vector))
	})
	if m.ensureErr != nil {
		return m.ensureErr
	}

	payload := buildVoiceSemanticPayload(event)
	if err := memory.ValidateCanonicalMemoryPayload(payload); err != nil {
		return fmt.Errorf("voice semantic payload rejected: %w", err)
	}
	pointID := strings.TrimSpace(event.JobID)
	if pointID == "" {
		pointID = buildVoiceSourceID(event, payload["ts"].(int64))
	}

	body, err := json.Marshal(map[string]any{
		"points": []map[string]any{{
			"id":      pointID,
			"vector":  vector,
			"payload": payload,
		}},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, m.baseURL+"/collections/"+m.collection+"/points?wait=true", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", m.apiKey)

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("qdrant mirror returned %s", resp.Status)
	}
	return nil
}

func buildVoiceSemanticPayload(event TranscriptEvent) map[string]any {
	createdAt := event.CreatedAt.UTC()
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	ts := createdAt.Unix()

	payload := map[string]any{
		// Canonical payload fields for Qdrant/schema alignment.
		"app_id":           "aurelia",
		"repo_id":          "aurelia",
		"environment":      "local",
		"text":             event.Transcript,
		"canonical_bot_id": defaultCanonicalBotID,
		"source_system":    defaultVoiceSourceSystem,
		"source_id":        buildVoiceSourceID(event, ts),
		"domain":           defaultVoiceDomain,
		"ts":               ts,
		"version":          defaultPayloadVersion,

		// Legacy/backward-compatible fields.
		"user_id":      event.UserID,
		"chat_id":      event.ChatID,
		"source":       event.Source,
		"transcript":   event.Transcript,
		"accepted":     event.Accepted,
		"requires_tts": event.RequiresTTS,
		"created_at":   createdAt.Format(time.RFC3339Nano),
	}
	return payload
}

func buildVoiceSourceID(event TranscriptEvent, ts int64) string {
	jobID := strings.TrimSpace(event.JobID)
	if jobID != "" {
		return "voice:" + jobID
	}
	return "voice:ts:" + strconv.FormatInt(ts, 10)
}

func (m *QdrantMirror) embed(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(map[string]any{
		"model": m.embeddingModel,
		"input": text,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.embedURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ollama embed returned %s", resp.Status)
	}

	var payload struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload.Embeddings) == 0 || len(payload.Embeddings[0]) == 0 {
		return nil, fmt.Errorf("ollama embed returned no vectors")
	}
	return payload.Embeddings[0], nil
}

func (m *QdrantMirror) ensureCollection(ctx context.Context, size int) error {
	body, err := json.Marshal(map[string]any{
		"vectors": map[string]any{
			"size":     size,
			"distance": "Cosine",
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, m.baseURL+"/collections/"+m.collection, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", m.apiKey)

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusConflict {
		return fmt.Errorf("qdrant ensure collection returned %s", resp.Status)
	}
	return nil
}
