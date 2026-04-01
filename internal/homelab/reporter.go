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
	Timestamp    time.Time
	CPUTemp      string
	GPUTemp      string
	GPUUtil      string
	GPUPower     string
	GPUVRAMUsed  string
	GPUVRAMTotal string
	RAMUsed      string
	RAMTotal     string
	Containers   string
	Ollama       string
	LiteLLM      string
	Qdrant       string
	Redis        string
	ZFS          string
	GrafanaURL   string
}

func Collect(ctx context.Context) Report {
	grafanaURL := os.Getenv("GRAFANA_URL")
	if grafanaURL == "" {
		grafanaURL = "https://monitor.zappro.site"
	}

	vram := runCmd("nvidia-smi --query-gpu=memory.used,memory.total --format=csv,noheader | tr ',' ' / ' | sed 's/ //g'")
	ram := runCmd("free -h | awk '/^Mem/{print $3\" / \"$2}'")

	return Report{
		Timestamp:    time.Now(),
		CPUTemp:      runCmd("sensors | grep 'Tctl' | awk '{print $2}'"),
		GPUTemp:      runCmd("nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader,nounits | awk '{print $1\"°C\"}'"),
		GPUUtil:      runCmd("nvidia-smi --query-gpu=utilization.gpu --format=csv,noheader,nounits | awk '{print $1\"%\"}'"),
		GPUPower:     runCmd("nvidia-smi --query-gpu=power.draw --format=csv,noheader,nounits | awk '{print $1\"W\"}'"),
		GPUVRAMUsed:  strings.Split(vram, " / ")[0],
		GPUVRAMTotal: strings.Split(vram, " / ")[1],
		RAMUsed:      strings.Split(ram, " / ")[0],
		RAMTotal:     strings.Split(ram, " / ")[1],
		Containers:   runCmd("docker ps --format '{{.Names}}' | wc -l | tr -d ' '"),
		Ollama:       checkHTTP(ctx, "http://localhost:11434/api/tags"),
		LiteLLM:      checkHTTP(ctx, "http://localhost:4000/health"),
		Qdrant:       checkHTTP(ctx, "http://localhost:6333/health"),
		Redis:        checkHTTP(ctx, "http://localhost:6379/health"),
		ZFS:          runCmd("zfs list -H -o name,avail 2>/dev/null | grep tank | awk '{print $2}' | head -1"),
		GrafanaURL:   grafanaURL,
	}
}

func (r Report) Format() string {
	grafana := r.GrafanaURL
	if grafana == "" {
		grafana = "https://monitor.zappro.site"
	}

	return fmt.Sprintf(`*🏠 Home Lab — %s*

*GPU*
├ Temp: %s | Util: %s | Power: %s
└ VRAM: %s / %s

*Sistema*
├ RAM: %s / %s
└ Containers: %s

*Serviços*
├ Ollama: %s
├ LiteLLM: %s
├ Qdrant: %s
└ Redis: %s

%s

[📊 Grafana](%s)`,
		r.Timestamp.Format("02/01 15:04"),
		r.GPUTemp, r.GPUUtil, r.GPUPower,
		r.GPUVRAMUsed, r.GPUVRAMTotal,
		r.RAMUsed, r.RAMTotal,
		r.Containers,
		statusEmoji(r.Ollama),
		statusEmoji(r.LiteLLM),
		statusEmoji(r.Qdrant),
		statusEmoji(r.Redis),
		zfsLine(r.ZFS),
		grafana,
	)
}

func statusEmoji(s string) string {
	if s == "ok" || s == "healthy" {
		return "✅"
	}
	return "🔴 " + s
}

func zfsLine(zfs string) string {
	if zfs == "" {
		return ""
	}
	return fmt.Sprintf("*ZFS* livre: %s", zfs)
}

func (r Report) statusIconVRAM(pct string) string {
	var p int
	fmt.Sscanf(pct, "%d", &p)
	switch {
	case p < 60:
		return "🟢"
	case p < 80:
		return "🟡"
	default:
		return "🔴"
	}
}

func (r Report) statusIconRAM(pct string) string {
	var p int
	fmt.Sscanf(pct, "%d", &p)
	switch {
	case p < 70:
		return "🟢"
	case p < 85:
		return "🟡"
	default:
		return "🔴"
	}
}

func (r Report) formatContainerList() string {
	out, err := exec.Command("bash", "-c", "docker ps --format 'table {{.Names}}\t{{.Status}}' | tail -n +2").Output()
	if err != nil {
		return "N/A"
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 8 {
		lines = lines[:8]
		lines = append(lines, "  ... e mais")
	}
	var result []string
	for _, line := range lines {
		if line != "" {
			result = append(result, "  "+strings.ReplaceAll(line, "\t", " • "))
		}
	}
	return strings.Join(result, "\n")
}

func calcPercent(used, total string) string {
	var u, t int
	fmt.Sscanf(used, "%d", &u)
	fmt.Sscanf(total, "%d", &t)
	if t == 0 {
		return "N/A"
	}
	pct := (u * 100) / t
	return fmt.Sprintf("%d", pct)
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
		return "✅ online"
	}
	return "⚠️ " + status
}
