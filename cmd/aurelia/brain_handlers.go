package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type brainPoint struct {
	ID      interface{}            `json:"id"`
	Score   float64                `json:"score,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func qdrantRequest(method, reqURL, apiKey string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("api-key", apiKey)
	}
	return http.DefaultClient.Do(req)
}

// buildBrainSearchHandler returns a handler for /api/brain/search?q=<query>
// S-26: For now, performs a scroll (list) filtered by text if q is provided,
// since embedding is not available without Ollama integration here.
func buildBrainSearchHandler(qdrantURL, collection, apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		q := r.URL.Query().Get("q")

		// Use scroll API to fetch points, then filter by payload text match
		scrollURL := fmt.Sprintf("%s/collections/%s/points/scroll", qdrantURL, url.PathEscape(collection))
		scrollBody := map[string]interface{}{
			"limit":        20,
			"with_payload": true,
			"with_vector":  false,
		}
		bodyBytes, _ := json.Marshal(scrollBody)
		resp, err := qdrantRequest(http.MethodPost, scrollURL, apiKey, bodyBytes)
		if err != nil {
			// Graceful fallback: return empty
			_ = json.NewEncoder(w).Encode([]brainPoint{})
			return
		}
		defer resp.Body.Close()

		var scrollResp struct {
			Result struct {
				Points []struct {
					ID      interface{}            `json:"id"`
					Payload map[string]interface{} `json:"payload"`
				} `json:"points"`
			} `json:"result"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&scrollResp); err != nil {
			_ = json.NewEncoder(w).Encode([]brainPoint{})
			return
		}

		results := make([]brainPoint, 0, len(scrollResp.Result.Points))
		for _, p := range scrollResp.Result.Points {
			if q != "" {
				// Simple text match in payload values
				match := false
				for _, v := range p.Payload {
					if s, ok := v.(string); ok {
						if containsCI(s, q) {
							match = true
							break
						}
					}
				}
				if !match {
					continue
				}
			}
			results = append(results, brainPoint{
				ID:      p.ID,
				Payload: p.Payload,
			})
		}
		_ = json.NewEncoder(w).Encode(results)
	}
}

// buildBrainRecentHandler returns a handler for /api/brain/recent
// Returns the latest 10 points from Qdrant.
func buildBrainRecentHandler(qdrantURL, collection, apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		scrollURL := fmt.Sprintf("%s/collections/%s/points/scroll", qdrantURL, url.PathEscape(collection))
		scrollBody := map[string]interface{}{
			"limit":        10,
			"with_payload": true,
			"with_vector":  false,
		}
		bodyBytes, _ := json.Marshal(scrollBody)
		resp, err := qdrantRequest(http.MethodPost, scrollURL, apiKey, bodyBytes)
		if err != nil {
			_ = json.NewEncoder(w).Encode([]brainPoint{})
			return
		}
		defer resp.Body.Close()

		var scrollResp struct {
			Result struct {
				Points []struct {
					ID      interface{}            `json:"id"`
					Payload map[string]interface{} `json:"payload"`
				} `json:"points"`
			} `json:"result"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&scrollResp); err != nil {
			_ = json.NewEncoder(w).Encode([]brainPoint{})
			return
		}

		results := make([]brainPoint, 0, len(scrollResp.Result.Points))
		for _, p := range scrollResp.Result.Points {
			results = append(results, brainPoint{
				ID:      p.ID,
				Payload: p.Payload,
			})
		}
		_ = json.NewEncoder(w).Encode(results)
	}
}

func containsCI(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	sl := len(s)
	subl := len(substr)
	if subl > sl {
		return false
	}
	for i := 0; i <= sl-subl; i++ {
		match := true
		for j := 0; j < subl; j++ {
			cs := s[i+j]
			cq := substr[j]
			if cs >= 'A' && cs <= 'Z' {
				cs += 32
			}
			if cq >= 'A' && cq <= 'Z' {
				cq += 32
			}
			if cs != cq {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
