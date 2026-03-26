package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/memory"
)

type brainPoint struct {
	ID      interface{}            `json:"id"`
	Score   float64                `json:"score,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func buildBrainSearchHandler(qdrantURL, collection, apiKey, ollamaURL, embeddingModel string) http.HandlerFunc {
	client := memory.NewSemanticHTTPClient(10 * time.Second)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		query := strings.TrimSpace(r.URL.Query().Get("q"))
		if query == "" {
			points, err := memory.ScrollPoints(r.Context(), client, qdrantURL, collection, apiKey, 20)
			if err != nil {
				writeBrainHeaders(w, "recent", "degraded", err)
				_ = json.NewEncoder(w).Encode([]brainPoint{})
				return
			}
			sortSemanticPoints(points)
			writeBrainHeaders(w, "recent", "ok", nil)
			_ = json.NewEncoder(w).Encode(toBrainPoints(points))
			return
		}

		vector, err := memory.EmbedText(r.Context(), client, ollamaURL, embeddingModel, query)
		if err == nil {
			points, searchErr := memory.SearchSemantic(r.Context(), client, qdrantURL, collection, apiKey, vector, 20)
			if searchErr == nil {
				writeBrainHeaders(w, "semantic", "ok", nil)
				_ = json.NewEncoder(w).Encode(toBrainPoints(points))
				return
			}
			err = searchErr
		}

		points, fallbackErr := lexicalBrainFallback(r.Context(), client, qdrantURL, collection, apiKey, query)
		if fallbackErr != nil {
			writeBrainHeaders(w, "degraded", "degraded", joinErrors(err, fallbackErr))
			_ = json.NewEncoder(w).Encode([]brainPoint{})
			return
		}

		writeBrainHeaders(w, "lexical-fallback", "degraded", err)
		_ = json.NewEncoder(w).Encode(toBrainPoints(points))
	}
}

func buildBrainRecentHandler(qdrantURL, collection, apiKey string) http.HandlerFunc {
	client := memory.NewSemanticHTTPClient(10 * time.Second)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		points, err := memory.ScrollPoints(r.Context(), client, qdrantURL, collection, apiKey, 20)
		if err != nil {
			writeBrainHeaders(w, "recent", "degraded", err)
			_ = json.NewEncoder(w).Encode([]brainPoint{})
			return
		}
		sortSemanticPoints(points)
		writeBrainHeaders(w, "recent", "ok", nil)
		_ = json.NewEncoder(w).Encode(toBrainPoints(points[:min(len(points), 10)]))
	}
}

func lexicalBrainFallback(ctx context.Context, client *http.Client, qdrantURL, collection, apiKey, query string) ([]memory.SemanticPoint, error) {
	points, err := memory.ScrollPoints(ctx, client, qdrantURL, collection, apiKey, 100)
	if err != nil {
		return nil, err
	}
	filtered := memory.FilterPointsLexical(points, query)
	sortSemanticPoints(filtered)
	return filtered, nil
}

func toBrainPoints(points []memory.SemanticPoint) []brainPoint {
	results := make([]brainPoint, 0, len(points))
	for _, point := range points {
		payload := memory.NormalizeSemanticPayload(point.Payload)
		results = append(results, brainPoint{
			ID:      point.ID,
			Score:   point.Score,
			Payload: toStringMap(payload),
		})
	}
	return results
}

func toStringMap(payload map[string]any) map[string]interface{} {
	out := make(map[string]interface{}, len(payload))
	for key, value := range payload {
		out[key] = value
	}
	return out
}

func sortSemanticPoints(points []memory.SemanticPoint) {
	sort.SliceStable(points, func(i, j int) bool {
		left := memory.ExtractPayloadTimestamp(points[i].Payload)
		right := memory.ExtractPayloadTimestamp(points[j].Payload)
		if left.Equal(right) {
			return points[i].Score > points[j].Score
		}
		return left.After(right)
	})
}

func writeBrainHeaders(w http.ResponseWriter, mode, status string, err error) {
	if mode == "" {
		mode = "unknown"
	}
	if status == "" {
		status = "unknown"
	}
	w.Header().Set("X-Aurelia-Brain-Mode", mode)
	w.Header().Set("X-Aurelia-Brain-Status", status)
	if err != nil {
		w.Header().Set("X-Aurelia-Brain-Error", err.Error())
	}
}

func joinErrors(primary, secondary error) error {
	switch {
	case primary == nil:
		return secondary
	case secondary == nil:
		return primary
	default:
		return fmt.Errorf("%v; fallback lexical failed: %w", primary, secondary)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
