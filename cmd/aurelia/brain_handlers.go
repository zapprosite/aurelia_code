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
	return buildBrainSearchHandlerForCollections(qdrantURL, []string{collection}, apiKey, ollamaURL, embeddingModel)
}

func buildBrainSearchHandlerForCollections(qdrantURL string, collections []string, apiKey, ollamaURL, embeddingModel string) http.HandlerFunc {
	client := memory.NewSemanticHTTPClient(10 * time.Second)
	collections = normalizeBrainCollections(collections)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		query := strings.TrimSpace(r.URL.Query().Get("q"))
		if query == "" {
			points, err := scrollBrainCollections(r.Context(), client, qdrantURL, collections, apiKey, 20)
			if len(points) == 0 && err != nil {
				writeBrainHeaders(w, "recent", "degraded", err)
				_ = json.NewEncoder(w).Encode([]brainPoint{})
				return
			}
			sortSemanticPoints(points)
			status := "ok"
			if err != nil {
				status = "degraded"
			}
			writeBrainHeaders(w, "recent", status, err)
			_ = json.NewEncoder(w).Encode(toBrainPoints(points[:min(len(points), 20)]))
			return
		}

		vector, err := memory.EmbedText(r.Context(), client, ollamaURL, embeddingModel, query)
		if err == nil {
			points, searchErr := searchBrainCollections(r.Context(), client, qdrantURL, collections, apiKey, vector, 20)
			if len(points) > 0 {
				sortSemanticHits(points)
				status := "ok"
				if searchErr != nil {
					status = "degraded"
				}
				writeBrainHeaders(w, "semantic", status, searchErr)
				_ = json.NewEncoder(w).Encode(toBrainPoints(points[:min(len(points), 20)]))
				return
			}
			if searchErr != nil {
				err = searchErr
			}
		}

		points, fallbackErr := lexicalBrainFallback(r.Context(), client, qdrantURL, collections, apiKey, query)
		if len(points) == 0 && fallbackErr != nil {
			writeBrainHeaders(w, "degraded", "degraded", joinErrors(err, fallbackErr))
			_ = json.NewEncoder(w).Encode([]brainPoint{})
			return
		}

		status := "degraded"
		if err == nil && fallbackErr == nil {
			status = "ok"
		}
		writeBrainHeaders(w, "lexical-fallback", status, joinErrors(err, fallbackErr))
		_ = json.NewEncoder(w).Encode(toBrainPoints(points[:min(len(points), 20)]))
	}
}

func buildBrainRecentHandler(qdrantURL, collection, apiKey string) http.HandlerFunc {
	return buildBrainRecentHandlerForCollections(qdrantURL, []string{collection}, apiKey)
}

func buildBrainRecentHandlerForCollections(qdrantURL string, collections []string, apiKey string) http.HandlerFunc {
	client := memory.NewSemanticHTTPClient(10 * time.Second)
	collections = normalizeBrainCollections(collections)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		points, err := scrollBrainCollections(r.Context(), client, qdrantURL, collections, apiKey, 20)
		if len(points) == 0 && err != nil {
			writeBrainHeaders(w, "recent", "degraded", err)
			_ = json.NewEncoder(w).Encode([]brainPoint{})
			return
		}
		sortSemanticPoints(points)
		status := "ok"
		if err != nil {
			status = "degraded"
		}
		writeBrainHeaders(w, "recent", status, err)
		_ = json.NewEncoder(w).Encode(toBrainPoints(points[:min(len(points), 10)]))
	}
}

func lexicalBrainFallback(ctx context.Context, client *http.Client, qdrantURL string, collections []string, apiKey, query string) ([]memory.SemanticPoint, error) {
	points, err := scrollBrainCollections(ctx, client, qdrantURL, collections, apiKey, 100)
	if len(points) == 0 && err != nil {
		return nil, err
	}
	filtered := memory.FilterPointsLexical(points, query)
	sortSemanticPoints(filtered)
	return filtered, err
}

func searchBrainCollections(ctx context.Context, client *http.Client, qdrantURL string, collections []string, apiKey string, vector []float32, limit int) ([]memory.SemanticPoint, error) {
	var out []memory.SemanticPoint
	var combinedErr error
	for _, collection := range collections {
		points, err := memory.SearchSemantic(ctx, client, qdrantURL, collection, apiKey, vector, limit)
		if err != nil {
			combinedErr = joinErrors(combinedErr, fmt.Errorf("%s: %w", collection, err))
			continue
		}
		out = append(out, points...)
	}
	if len(out) == 0 {
		return nil, combinedErr
	}
	return out, combinedErr
}

func scrollBrainCollections(ctx context.Context, client *http.Client, qdrantURL string, collections []string, apiKey string, limit int) ([]memory.SemanticPoint, error) {
	var out []memory.SemanticPoint
	var combinedErr error
	for _, collection := range collections {
		points, err := memory.ScrollPoints(ctx, client, qdrantURL, collection, apiKey, limit)
		if err != nil {
			combinedErr = joinErrors(combinedErr, fmt.Errorf("%s: %w", collection, err))
			continue
		}
		out = append(out, points...)
	}
	if len(out) == 0 {
		return nil, combinedErr
	}
	return out, combinedErr
}

func normalizeBrainCollections(collections []string) []string {
	out := make([]string, 0, len(collections))
	seen := make(map[string]struct{}, len(collections))
	for _, collection := range collections {
		collection = strings.TrimSpace(collection)
		if collection == "" {
			continue
		}
		if _, ok := seen[collection]; ok {
			continue
		}
		seen[collection] = struct{}{}
		out = append(out, collection)
	}
	return out
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

func sortSemanticHits(points []memory.SemanticPoint) {
	sort.SliceStable(points, func(i, j int) bool {
		if points[i].Score == points[j].Score {
			left := memory.ExtractPayloadTimestamp(points[i].Payload)
			right := memory.ExtractPayloadTimestamp(points[j].Payload)
			return left.After(right)
		}
		return points[i].Score > points[j].Score
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
		return fmt.Errorf("%v; additional error: %w", primary, secondary)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
