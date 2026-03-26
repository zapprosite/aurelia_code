package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// HomelabCheck holds the result of a single homelab health probe.
type HomelabCheck struct {
	Component string `json:"component"`
	Status    string `json:"status"` // healthy | degraded | offline
	Summary   string `json:"summary"`
	LatencyMS int64  `json:"latency_ms,omitempty"`
}

// HomelabStatus is the aggregate result of all homelab checks.
type HomelabStatus struct {
	Checks    []HomelabCheck `json:"checks"`
	CheckedAt time.Time      `json:"checked_at"`
}

// CheckHomelab runs all homelab checks and returns the aggregate status.
// It never panics — failed probes are reported as "degraded" or "offline".
func CheckHomelab(ctx context.Context, ollamaURL, qdrantURL string) HomelabStatus {
	if ollamaURL == "" {
		ollamaURL = "http://127.0.0.1:11434"
	}
	if qdrantURL == "" {
		qdrantURL = "http://127.0.0.1:6333"
	}

	status := HomelabStatus{CheckedAt: time.Now().UTC()}
	status.Checks = []HomelabCheck{
		checkOllama(ctx, ollamaURL),
		checkQdrant(ctx, qdrantURL),
		checkGPU(ctx),
		checkDocker(ctx),
		checkDisk(ctx),
	}
	return status
}

// Summary returns a human-readable one-line summary.
func (s HomelabStatus) Summary() string {
	healthy, degraded, offline := 0, 0, 0
	for _, c := range s.Checks {
		switch c.Status {
		case "healthy":
			healthy++
		case "degraded":
			degraded++
		default:
			offline++
		}
	}
	if offline > 0 || degraded > 0 {
		return fmt.Sprintf("homelab: %d healthy / %d degraded / %d offline", healthy, degraded, offline)
	}
	return fmt.Sprintf("homelab: todos %d componentes saudáveis", healthy)
}

// checkOllama pings Ollama /api/tags
func checkOllama(ctx context.Context, ollamaURL string) HomelabCheck {
	start := time.Now()
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		strings.TrimRight(ollamaURL, "/")+"/api/tags", nil)
	if err != nil {
		return HomelabCheck{Component: "ollama", Status: "offline", Summary: err.Error()}
	}
	resp, err := client.Do(req)
	if err != nil {
		return HomelabCheck{Component: "ollama", Status: "offline", Summary: "unreachable: " + err.Error()}
	}
	defer resp.Body.Close()
	latency := time.Since(start).Milliseconds()
	if resp.StatusCode != http.StatusOK {
		return HomelabCheck{Component: "ollama", Status: "degraded",
			Summary: fmt.Sprintf("HTTP %d", resp.StatusCode), LatencyMS: latency}
	}
	return HomelabCheck{Component: "ollama", Status: "healthy",
		Summary: "reachable", LatencyMS: latency}
}

// checkQdrant pings Qdrant /readyz
func checkQdrant(ctx context.Context, qdrantURL string) HomelabCheck {
	start := time.Now()
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		strings.TrimRight(qdrantURL, "/")+"/readyz", nil)
	if err != nil {
		return HomelabCheck{Component: "qdrant", Status: "offline", Summary: err.Error()}
	}
	resp, err := client.Do(req)
	if err != nil {
		return HomelabCheck{Component: "qdrant", Status: "offline", Summary: "unreachable: " + err.Error()}
	}
	defer resp.Body.Close()
	latency := time.Since(start).Milliseconds()
	if resp.StatusCode != http.StatusOK {
		return HomelabCheck{Component: "qdrant", Status: "degraded",
			Summary: fmt.Sprintf("HTTP %d", resp.StatusCode), LatencyMS: latency}
	}
	return HomelabCheck{Component: "qdrant", Status: "healthy",
		Summary: "reachable", LatencyMS: latency}
}

// checkGPU runs nvidia-smi and extracts temperature + memory usage.
func checkGPU(ctx context.Context) HomelabCheck {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(tCtx,
		"nvidia-smi",
		"--query-gpu=temperature.gpu,memory.used,memory.total",
		"--format=csv,noheader,nounits",
	).Output()
	if err != nil {
		return HomelabCheck{Component: "gpu", Status: "offline", Summary: "nvidia-smi unavailable: " + err.Error()}
	}

	line := strings.TrimSpace(string(out))
	parts := strings.Split(line, ",")
	if len(parts) < 3 {
		return HomelabCheck{Component: "gpu", Status: "degraded", Summary: "unexpected nvidia-smi output: " + line}
	}

	tempStr := strings.TrimSpace(parts[0])
	usedStr := strings.TrimSpace(parts[1])
	totalStr := strings.TrimSpace(parts[2])

	temp, _ := strconv.Atoi(tempStr)
	status := "healthy"
	if temp >= 85 {
		status = "degraded"
	}

	return HomelabCheck{
		Component: "gpu",
		Status:    status,
		Summary:   fmt.Sprintf("temp=%s°C vram=%s/%s MiB", tempStr, usedStr, totalStr),
	}
}

// checkDocker lists containers and reports any that are unhealthy or stopped unexpectedly.
func checkDocker(ctx context.Context) HomelabCheck {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(tCtx,
		"docker", "ps", "--format", `{"name":"{{.Names}}","status":"{{.Status}}"}`).Output()
	if err != nil {
		return HomelabCheck{Component: "docker", Status: "offline", Summary: "docker unavailable: " + err.Error()}
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	total, unhealthy := 0, 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var container struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal([]byte(line), &container); err != nil {
			continue
		}
		total++
		if strings.Contains(strings.ToLower(container.Status), "unhealthy") {
			unhealthy++
		}
	}

	if unhealthy > 0 {
		return HomelabCheck{Component: "docker", Status: "degraded",
			Summary: fmt.Sprintf("%d/%d containers unhealthy", unhealthy, total)}
	}
	return HomelabCheck{Component: "docker", Status: "healthy",
		Summary: fmt.Sprintf("%d containers running", total)}
}

// checkDisk checks available space on the primary data mount.
func checkDisk(ctx context.Context) HomelabCheck {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(tCtx, "df", "-h", "--output=avail,pcent", "/").Output()
	if err != nil {
		return HomelabCheck{Component: "disk", Status: "offline", Summary: "df unavailable: " + err.Error()}
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return HomelabCheck{Component: "disk", Status: "degraded", Summary: "unexpected df output"}
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return HomelabCheck{Component: "disk", Status: "degraded", Summary: "cannot parse df output"}
	}

	avail := fields[0]
	pctStr := strings.TrimSuffix(fields[1], "%")
	pct, _ := strconv.Atoi(pctStr)

	status := "healthy"
	if pct >= 90 {
		status = "degraded"
	}

	return HomelabCheck{Component: "disk", Status: status,
		Summary: fmt.Sprintf("avail=%s used=%s%%", avail, pctStr)}
}
