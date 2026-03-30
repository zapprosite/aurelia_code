package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// TavilySearchURL is the Tavily search API endpoint.
const tavilySearchURL = "https://api.tavily.com/search"

// WebSearchHandler executes a Tavily search (if TAVILY_API_KEY is set) or falls
// back to DuckDuckGo HTML scraping.
func WebSearchHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	query, err := requireStringArg(args, "query")
	if err != nil {
		return "", err
	}

	count := 5
	if c, ok := args["count"].(float64); ok && c > 0 && c <= 10 {
		count = int(c)
	}

	if key := os.Getenv("TAVILY_API_KEY"); key != "" {
		return tavilySearch(ctx, key, query, count)
	}
	return duckDuckGoSearch(ctx, query, count)
}

// tavilySearch calls the Tavily Search API and returns structured results.
func tavilySearch(ctx context.Context, apiKey, query string, count int) (string, error) {
	payload := map[string]interface{}{
		"api_key":              apiKey,
		"query":                query,
		"search_depth":         "basic",
		"max_results":          count,
		"include_answer":       true,
		"include_raw_content":  false,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", tavilySearchURL, bytes.NewReader(body))
	if err != nil {
		return duckDuckGoSearch(ctx, query, count) // fallback
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return duckDuckGoSearch(ctx, query, count) // fallback
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		Answer  string `json:"answer"`
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
			Score   float64 `json:"score"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return duckDuckGoSearch(ctx, query, count)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**Resultados Tavily para:** %s\n\n", query))

	if result.Answer != "" {
		sb.WriteString(fmt.Sprintf("**Resposta direta:** %s\n\n", result.Answer))
	}

	for i, r := range result.Results {
		if i >= count {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. **%s**\n   %s\n", i+1, r.Title, r.URL))
		if r.Content != "" {
			snippet := r.Content
			if len(snippet) > 300 {
				snippet = snippet[:300] + "..."
			}
			sb.WriteString(fmt.Sprintf("   %s\n", snippet))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// duckDuckGoSearch is the legacy fallback — HTML scraping, no API key needed.
var duckDuckGoBaseURL = "https://html.duckduckgo.com/html/?q=%s"

// DuckDuckGoBaseURL is exported so tests can override the endpoint.
var DuckDuckGoBaseURL = duckDuckGoBaseURL

func duckDuckGoSearch(ctx context.Context, query string, count int) (string, error) {
	searchURL := fmt.Sprintf(DuckDuckGoBaseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return fmt.Sprintf("request creation error: %v", err), nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("request failed: %v", err), nil
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("HTTP Error %d", resp.StatusCode), nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("failed reading response: %v", err), nil
	}

	htmlStr := string(bodyBytes)

	reLink := regexp.MustCompile(`<a[^>]*class="[^"]*result__a[^"]*"[^>]*href="([^"]+)"[^>]*>([\s\S]*?)</a>`)
	matches := reLink.FindAllStringSubmatch(htmlStr, count+5)
	if len(matches) == 0 {
		return "No results found.", nil
	}

	reSnippet := regexp.MustCompile(`<a class="result__snippet[^"]*".*?>([\s\S]*?)</a>`)
	snippetMatches := reSnippet.FindAllStringSubmatch(htmlStr, count+5)

	maxItems := count
	if len(matches) < count {
		maxItems = len(matches)
	}

	var results []string
	results = append(results, fmt.Sprintf("Results for: %s", query))

	for i := 0; i < maxItems; i++ {
		urlStr := matches[i][1]
		title := stripTags(matches[i][2])
		results = append(results, fmt.Sprintf("%d. %s\n   %s", i+1, title, decodeDuckDuckGoURL(urlStr)))
		if i < len(snippetMatches) {
			if snippet := stripTags(snippetMatches[i][1]); snippet != "" {
				results = append(results, fmt.Sprintf("   %s", snippet))
			}
		}
	}

	return strings.Join(results, "\n"), nil
}

func stripTags(content string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return strings.TrimSpace(re.ReplaceAllString(content, ""))
}

func decodeDuckDuckGoURL(raw string) string {
	if strings.Contains(raw, "uddg=") {
		if decoded, err := url.QueryUnescape(raw); err == nil {
			if idx := strings.Index(decoded, "uddg="); idx != -1 {
				return decoded[idx+5:]
			}
		}
	}
	return raw
}
