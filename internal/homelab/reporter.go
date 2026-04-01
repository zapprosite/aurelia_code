package homelab

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Report struct {
	Timestamp   time.Time
	CPUTemp     string
	GPUTemp     string
	GPUVRAMUsed string
	RAMUsed     string
	Containers  string
	Ollama      string
	LiteLLM     string
	Qdrant      string
	ZFS         string
	GrafanaURL  string
}

func Collect(ctx context.Context) Report {
	grafanaURL := os.Getenv("GRAFANA_URL")
	if grafanaURL == "" {
		grafanaURL = "https://monitor.zappro.site"
	}

	return Report{
		Timestamp:   time.Now(),
		CPUTemp:     runCmd("sensors | grep 'Tctl' | awk '{print $2}'"),
		GPUTemp:     runCmd("nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader,nounits | awk '{print $1\"°C\"}'"),
		GPUVRAMUsed: runCmd("nvidia-smi --query-gpu=memory.used,memory.total --format=csv,noheader | tr ',' '/' | sed 's/ //g'"),
		RAMUsed:     runCmd("free -h | awk '/^Mem/{print $3\"/\"$2}'"),
		Containers:  runCmd("docker ps --format '{{.Names}}' | wc -l | tr -d ' '"),
		Ollama:      checkHTTP(ctx, "http://localhost:11434/api/tags"),
		LiteLLM:     checkHTTP(ctx, "http://localhost:4000/health"),
		Qdrant:      checkHTTP(ctx, "http://localhost:6333/health"),
		ZFS:         runCmd("zfs list -H -o name,avail 2>/dev/null | grep tank | awk '{print $2}' | head -1"),
		GrafanaURL:  grafanaURL,
	}
}

func (r Report) Format() string {
	hora := r.Timestamp.Format("02/01 15:04")

	containers := r.Containers
	if containers == "" {
		containers = "N/A"
	}

	return fmt.Sprintf(
		"🏠 *Home Lab* — %s\n\n"+
			"🔥 CPU: `%s` | GPU: `%s`\n"+
			"💾 VRAM: `%s` | RAM: `%s`\n"+
			"🐳 Containers: `%s` ativos\n\n"+
			"🟢 Ollama:   %s\n"+
			"🟢 LiteLLM:  %s\n"+
			"🟢 Qdrant:   %s\n"+
			"💽 ZFS livre: `%s`\n\n"+
			"📊 [Ver Grafana](%s)",
		hora,
		r.CPUTemp, r.GPUTemp,
		r.GPUVRAMUsed, r.RAMUsed,
		containers,
		statusIcon(r.Ollama),
		statusIcon(r.LiteLLM),
		statusIcon(r.Qdrant),
		r.ZFS,
		r.GrafanaURL,
	)
}

func runCmd(cmdStr string) string {
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(out))
}

func checkHTTP(ctx context.Context, url string) string {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "curl", "-sf", "-o", "/dev/null", "-w", "%{http_code}", url)
	out, err := cmd.Output()
	if err != nil {
		return "error"
	}
	code := strings.TrimSpace(string(out))
	if code == "200" || code == "401" {
		return "ok"
	}
	return "error:" + code
}

func statusIcon(status string) string {
	if status == "ok" {
		return "✅ OK"
	}
	return "⚠️ " + status
}
