package metrics

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// SystemMetrics holds real-time system and GPU metrics.
type SystemMetrics struct {
	VRAMUsedGiB   *float64 `json:"vram_used_gib"`
	VRAMTotalGiB  *float64 `json:"vram_total_gib"`
	GPUUtilPct    *float64 `json:"gpu_util_percent"`
	CPULoad       float64  `json:"cpu_load"`
	MemUsedGiB    float64  `json:"mem_used_gib"`
	MemTotalGiB   float64  `json:"mem_total_gib"`
}

// Collect gathers current system metrics. GPU fields are nil if nvidia-smi is absent.
func Collect(ctx context.Context) SystemMetrics {
	m := SystemMetrics{}
	m.CPULoad, m.MemUsedGiB, m.MemTotalGiB = collectCPUMem()
	m.VRAMUsedGiB, m.VRAMTotalGiB, m.GPUUtilPct = collectGPU(ctx)
	return m
}

// Handler returns an http.HandlerFunc that serves /api/metrics as JSON.
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		m := Collect(ctx)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_ = json.NewEncoder(w).Encode(m)
	}
}

// collectCPUMem reads /proc/loadavg and /proc/meminfo for CPU load and memory.
func collectCPUMem() (cpuLoad, memUsedGiB, memTotalGiB float64) {
	// CPU load (1-minute average from /proc/loadavg)
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 1 {
			if v, err := strconv.ParseFloat(fields[0], 64); err == nil {
				cpuLoad = v
			}
		}
	}

	// Memory from /proc/meminfo
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return cpuLoad, 0, 0
	}
	defer f.Close()

	var memTotal, memAvail, memFree int64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "MemTotal:"):
			memTotal = parseKiB(line)
		case strings.HasPrefix(line, "MemAvailable:"):
			memAvail = parseKiB(line)
		case strings.HasPrefix(line, "MemFree:"):
			memFree = parseKiB(line)
		}
	}
	if memAvail == 0 {
		memAvail = memFree
	}
	const gib = 1024.0 * 1024.0
	memTotalGiB = float64(memTotal) / gib
	memUsedGiB = float64(memTotal-memAvail) / gib
	return cpuLoad, memUsedGiB, memTotalGiB
}

func parseKiB(line string) int64 {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0
	}
	v, _ := strconv.ParseInt(fields[1], 10, 64)
	return v
}

// collectGPU calls nvidia-smi to get VRAM and utilization. Returns nil pointers if unavailable.
func collectGPU(ctx context.Context) (vramUsed, vramTotal *float64, gpuUtil *float64) {
	cmd := exec.CommandContext(ctx, "nvidia-smi",
		"--query-gpu=memory.used,memory.total,utilization.gpu",
		"--format=csv,noheader,nounits")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, nil
	}

	line := strings.TrimSpace(string(out))
	// nvidia-smi may return multiple lines for multiple GPUs; take first.
	if idx := strings.Index(line, "\n"); idx >= 0 {
		line = line[:idx]
	}
	parts := strings.Split(line, ",")
	if len(parts) < 3 {
		return nil, nil, nil
	}

	usedMiB, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	totalMiB, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	util, err3 := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, nil, nil
	}

	const mibPerGiB = 1024.0
	usedGiB := usedMiB / mibPerGiB
	totalGiB := totalMiB / mibPerGiB
	return &usedGiB, &totalGiB, &util
}
