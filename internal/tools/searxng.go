package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SearchWebLocalHandler executes a search query against a local SearXNG instance.
func SearchWebLocalHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	query, err := requireStringArg(args, "query")
	if err != nil {
		return "", err
	}

	// We need the config to get the SearXNG URL.
	// Since the tool handler doesn't receive the config directly, we rely on a package-level
	// variable or a closure. In the current wiring.go, we can inject it.
	// For now, we'll try to get it from a global or use the default if not set.
	targetURL := getSearXNGURL()
	if targetURL == "" {
		return "", fmt.Errorf("SEARXNG_URL not configured")
	}

	count := 5
	if c, ok := args["count"].(float64); ok && c > 0 && c <= 10 {
		count = int(c)
	}

	return searchSearXNG(ctx, targetURL, query, count)
}

var globalSearXNGURL string

func SetSearXNGURL(u string) {
	globalSearXNGURL = u
}

func getSearXNGURL() string {
	if globalSearXNGURL != "" {
		return globalSearXNGURL
	}
	return "http://localhost:8888" // Default
}

func searchSearXNG(ctx context.Context, baseURL, query string, count int) (string, error) {
	searchURL, err := url.Parse(baseURL + "/search")
	if err != nil {
		return "", fmt.Errorf("invalid SearXNG URL: %v", err)
	}

	q := searchURL.Query()
	q.Set("q", query)
	q.Set("format", "json")
	q.Set("language", "auto")
	searchURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL.String(), nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("searxng request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("searxng returned HTTP %d", resp.StatusCode)
	}

	var data struct {
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode searxng response: %v", err)
	}

	if len(data.Results) == 0 {
		return "Nenhum resultado encontrado localmente.", nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**Resultados Locais (SearXNG) para:** %s\n\n", query))

	for i, r := range data.Results {
		if i >= count {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. **%s**\n   %s\n", i+1, r.Title, r.URL))
		if r.Content != "" {
			snippet := r.Content
			if len(snippet) > 400 {
				snippet = snippet[:400] + "..."
			}
			sb.WriteString(fmt.Sprintf("   %s\n", snippet))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
