package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildBrainSearchHandlerUsesSemanticSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/embed":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"embeddings": [][]float32{{0.1, 0.2, 0.3}},
			})
		case "/collections/conversation_memory/points/search":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": []map[string]any{{
					"id":    "point-1",
					"score": 0.93,
					"payload": map[string]any{
						"text":       "Contrato HVAC aprovado",
						"bot_id":     "ac_vendas",
						"domain":     "business",
						"created_at": "2026-03-25T20:19:00Z",
					},
				}},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/brain/search?q=contrato", nil)
	rec := httptest.NewRecorder()

	buildBrainSearchHandler(server.URL, "conversation_memory", "", server.URL, "bge-m3").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("X-Aurelia-Brain-Mode"); got != "semantic" {
		t.Fatalf("expected semantic mode, got %q", got)
	}
	if got := rec.Header().Get("X-Aurelia-Brain-Status"); got != "ok" {
		t.Fatalf("expected ok status, got %q", got)
	}

	var points []brainPoint
	if err := json.NewDecoder(rec.Body).Decode(&points); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(points) != 1 {
		t.Fatalf("expected 1 point, got %d", len(points))
	}
	if text := points[0].Payload["text"]; text != "Contrato HVAC aprovado" {
		t.Fatalf("expected normalized text, got %#v", text)
	}
	if botID := points[0].Payload["canonical_bot_id"]; botID != "ac_vendas" {
		t.Fatalf("expected canonical bot id from payload, got %#v", botID)
	}
}

func TestBuildBrainSearchHandlerFallsBackLexicallyWithExplicitDegradation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/embed":
			http.Error(w, "embed offline", http.StatusBadGateway)
		case "/collections/conversation_memory/points/scroll":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"points": []map[string]any{
						{
							"id": "legacy-1",
							"payload": map[string]any{
								"transcript": "Caixa com DDA pendente",
								"created_at": "2026-03-25T20:19:00Z",
							},
						},
						{
							"id": "legacy-2",
							"payload": map[string]any{
								"text": "Obras sem relação",
							},
						},
					},
				},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/brain/search?q=caixa", nil)
	rec := httptest.NewRecorder()

	buildBrainSearchHandler(server.URL, "conversation_memory", "", server.URL, "bge-m3").ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Aurelia-Brain-Mode"); got != "lexical-fallback" {
		t.Fatalf("expected lexical-fallback mode, got %q", got)
	}
	if got := rec.Header().Get("X-Aurelia-Brain-Status"); got != "degraded" {
		t.Fatalf("expected degraded status, got %q", got)
	}
	if got := rec.Header().Get("X-Aurelia-Brain-Error"); !strings.Contains(got, "ollama embed returned") {
		t.Fatalf("expected explicit degradation header, got %q", got)
	}

	var points []brainPoint
	if err := json.NewDecoder(rec.Body).Decode(&points); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(points) != 1 {
		t.Fatalf("expected 1 fallback point, got %d", len(points))
	}
	if text := points[0].Payload["text"]; text != "Caixa com DDA pendente" {
		t.Fatalf("expected transcript to be normalized into text, got %#v", text)
	}
}

func TestBuildBrainRecentHandlerSortsByTimestampDescending(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/collections/conversation_memory/points/scroll" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"points": []map[string]any{
					{"id": "older", "payload": map[string]any{"text": "antigo", "created_at": "2026-03-25T19:00:00Z"}},
					{"id": "newer", "payload": map[string]any{"text": "novo", "created_at": "2026-03-25T20:19:00Z"}},
				},
			},
		})
	}))
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/brain/recent", nil)
	rec := httptest.NewRecorder()

	buildBrainRecentHandler(server.URL, "conversation_memory", "").ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Aurelia-Brain-Mode"); got != "recent" {
		t.Fatalf("expected recent mode, got %q", got)
	}
	if got := rec.Header().Get("X-Aurelia-Brain-Status"); got != "ok" {
		t.Fatalf("expected ok status, got %q", got)
	}

	var points []brainPoint
	if err := json.NewDecoder(rec.Body).Decode(&points); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}
	if points[0].ID != "newer" {
		t.Fatalf("expected newest point first, got %#v", points[0].ID)
	}
}
