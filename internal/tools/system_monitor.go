package tools

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

type systemMetric string

const (
	metricStats    systemMetric = "stats"
	metricGPU      systemMetric = "gpu"
	metricProcess  systemMetric = "process"
	metricNetwork  systemMetric = "network"
)

type systemResult struct {
	Metric  string `json:"metric"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

func SystemMonitorHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	metric := optionalStringArg(args, "metric")
	if metric == "" {
		metric = "stats"
	}

	timeout := 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result := systemResult{Metric: metric}

	switch systemMetric(metric) {
	case metricStats:
		result = getSystemStats(ctx)
	case metricGPU:
		result = getGPUStatus(ctx)
	case metricProcess:
		result = getProcessList(ctx)
	case metricNetwork:
		result = getNetworkStatus(ctx)
	default:
		result.Error = "unknown metric: " + metric
	}

	payload, _ := json.Marshal(result)
	return string(payload), nil
}

func getSystemStats(ctx context.Context) systemResult {
	result := systemResult{Metric: "stats"}

	// CPU and memory from /proc
	cmd := exec.CommandContext(ctx, "bash", "-c", "echo 'CPU:'; cat /proc/loadavg; echo; echo 'Memory:'; free -h | grep Mem")
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Success = true
	result.Output = strings.TrimSpace(string(output))
	return result
}

func getGPUStatus(ctx context.Context) systemResult {
	result := systemResult{Metric: "gpu"}

	// Try nvidia-smi first
	cmd := exec.CommandContext(ctx, "nvidia-smi", "--query-gpu=index,name,temperature.gpu,utilization.gpu,memory.used,memory.total", "--format=csv,noheader")
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = "nvidia-smi not available or failed: " + err.Error()
		return result
	}

	result.Success = true
	result.Output = strings.TrimSpace(string(output))
	return result
}

func getProcessList(ctx context.Context) systemResult {
	result := systemResult{Metric: "process"}

	cmd := exec.CommandContext(ctx, "bash", "-c", "ps aux --sort=-%cpu,-%mem | head -11")
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Success = true
	result.Output = strings.TrimSpace(string(output))
	return result
}

func getNetworkStatus(ctx context.Context) systemResult {
	result := systemResult{Metric: "network"}

	cmd := exec.CommandContext(ctx, "bash", "-c", "echo 'Interfaces:'; ip addr show; echo; echo 'Connections:'; ss -tlnp 2>/dev/null | head -10")
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Success = true
	result.Output = strings.TrimSpace(string(output))
	return result
}
