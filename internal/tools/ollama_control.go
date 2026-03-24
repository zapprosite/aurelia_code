package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type ollamaAction string

const (
	ollamaActionStatus ollamaAction = "status"
	ollamaActionList   ollamaAction = "list"
	ollamaActionPull   ollamaAction = "pull"
	ollamaActionRun    ollamaAction = "run"
)

type ollamaResult struct {
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

const defaultOllamaAPI = "http://localhost:11434"

// NewOllamaControlHandler returns an OllamaControlHandler bound to the given baseURL.
// If baseURL is empty, it falls back to the default localhost address.
func NewOllamaControlHandler(baseURL string) func(ctx context.Context, args map[string]interface{}) (string, error) {
	if baseURL == "" {
		baseURL = defaultOllamaAPI
	}
	return func(ctx context.Context, args map[string]interface{}) (string, error) {
		return ollamaControlHandle(ctx, args, baseURL)
	}
}

// OllamaControlHandler is the legacy handler that uses the default Ollama URL.
func OllamaControlHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	return ollamaControlHandle(ctx, args, defaultOllamaAPI)
}

func ollamaControlHandle(ctx context.Context, args map[string]interface{}, baseURL string) (string, error) {
	action := optionalStringArg(args, "action")
	if action == "" {
		action = "status"
	}

	result := ollamaResult{Action: action}

	switch ollamaAction(action) {
	case ollamaActionStatus:
		result = ollamaStatus(ctx, baseURL)
	case ollamaActionList:
		result = ollamaListModels(ctx, baseURL)
	case ollamaActionPull:
		model := optionalStringArg(args, "model")
		if model == "" {
			result.Error = "model name required for pull"
		} else {
			result = ollamaPull(ctx, baseURL, model)
		}
	case ollamaActionRun:
		model := optionalStringArg(args, "model")
		prompt := optionalStringArg(args, "prompt")
		if model == "" || prompt == "" {
			result.Error = "model and prompt required for run"
		} else {
			result = ollamaRun(ctx, baseURL, model, prompt)
		}
	default:
		result.Error = "unknown action: " + action
	}

	payload, _ := json.Marshal(result)
	return string(payload), nil
}

func ollamaStatus(ctx context.Context, baseURL string) ollamaResult {
	result := ollamaResult{Action: "status"}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := http.Get(baseURL + "/api/tags")
	if err != nil {
		result.Error = "Ollama not available: " + err.Error()
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result.Error = "Ollama returned status: " + resp.Status
		return result
	}

	result.Success = true
	result.Output = "Ollama is running"
	return result
}

func ollamaListModels(ctx context.Context, baseURL string) ollamaResult {
	result := ollamaResult{Action: "list"}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := http.Get(baseURL + "/api/tags")
	if err != nil {
		result.Error = "Failed to list models: " + err.Error()
		return result
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result.Success = true
	result.Output = string(body)
	return result
}

func ollamaPull(ctx context.Context, baseURL, model string) ollamaResult {
	result := ollamaResult{Action: "pull"}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/pull", strings.NewReader(`{"name":"`+model+`"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result.Success = resp.StatusCode == 200
	result.Output = strings.TrimSpace(string(body))
	if !result.Success {
		result.Error = "Pull failed with status: " + resp.Status
	}
	return result
}

func ollamaRun(ctx context.Context, baseURL, model, prompt string) ollamaResult {
	result := ollamaResult{Action: "run"}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	payloadStr := `{"model":"` + model + `","prompt":"` + strings.TrimSpace(prompt) + `","stream":false}`
	req, _ := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/generate", strings.NewReader(payloadStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result.Success = resp.StatusCode == 200
	result.Output = strings.TrimSpace(string(body))
	if !result.Success {
		result.Error = "Generate failed with status: " + resp.Status
	}
	return result
}
