package voice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const defaultOllamaEmbedURL = "http://127.0.0.1:11434/api/embed"
const (
	defaultSupabaseMirrorTable    = "aurelia_voice_events"
	defaultQdrantMirrorCollection = "conversation_memory"
	defaultQdrantMirrorEmbedding  = "bge-m3"
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

type SupabaseMirror struct {
	baseURL string
	apiKey  string
	table   string
	client  *http.Client
}

func NewSupabaseMirror(baseURL, apiKey, table string) *SupabaseMirror {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	table = strings.TrimSpace(table)
	if table == "" {
		table = defaultSupabaseMirrorTable
	}
	return &SupabaseMirror{
		baseURL: baseURL,
		apiKey:  strings.TrimSpace(apiKey),
		table:   table,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *SupabaseMirror) MirrorTranscript(ctx context.Context, event TranscriptEvent) error {
	if m == nil || m.baseURL == "" || m.apiKey == "" {
		return nil
	}

	body, err := json.Marshal(map[string]any{
		"job_id":       event.JobID,
		"user_id":      event.UserID,
		"chat_id":      event.ChatID,
		"source":       event.Source,
		"transcript":   event.Transcript,
		"accepted":     event.Accepted,
		"requires_tts": event.RequiresTTS,
		"processed_at": time.Now().UTC().Format(time.RFC3339Nano),
		"conversation": fmt.Sprintf("%d", event.UserID),
		"message_type": "audio_transcript",
		"created_at":   event.CreatedAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL+"/rest/v1/"+m.table, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", m.apiKey)
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Prefer", "return=minimal")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("supabase mirror returned %s", resp.Status)
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

func NewQdrantMirror(baseURL, apiKey, collection, embeddingModel string) *QdrantMirror {
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
		embedURL:       defaultOllamaEmbedURL,
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

	body, err := json.Marshal(map[string]any{
		"points": []map[string]any{{
			"id":     event.JobID,
			"vector": vector,
			"payload": map[string]any{
				"user_id":      event.UserID,
				"chat_id":      event.ChatID,
				"source":       event.Source,
				"transcript":   event.Transcript,
				"accepted":     event.Accepted,
				"requires_tts": event.RequiresTTS,
				"created_at":   event.CreatedAt.Format(time.RFC3339Nano),
			},
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
