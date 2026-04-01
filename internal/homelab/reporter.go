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
	hora := r.Timestamp.Format("02/01/2026 às 15:04")

	containers := r.Containers
	if containers == "" {
		containers = "N/A"
	}

	vramPct := calcPercent(r.GPUVRAMUsed, r.GPUVRAMTotal)
	ramPct := calcPercent(r.RAMUsed, r.RAMTotal)

	return fmt.Sprintf(`🎛️ *Aurélia Home Lab — Status* 🖥️

**%s**

---

### 🚀 *Recursos do Sistema*

| Recurso | Uso | Status |
|---------|-----|--------|
| **CPU** | %s | %s |
| **GPU** | %s @ %s | %s |
| **VRAM** | %s / %s (%s%%) | %s |
| **RAM** | %s / %s (%s%%) | %s |

---

### 🐳 *Containers* (%s ativos)

%s

---

### 🧠 *Serviços de IA*

| Serviço | Endpoint | Status |
|---------|----------|--------|
| Ollama | localhost:11434 | %s |
| LiteLLM | localhost:4000 | %s |
| Qdrant | localhost:6333 | %s |
| Redis | localhost:6379 | %s |

---

### 💾 *Armazenamento*

- **ZFS (tank):** %s disponível

---

### 📊 *Links*

[📈 Ver Grafana](%s) • [🖥️ Access CapRover](https://cap.zappro.site)

---

_*Aurélia — Monitoramento Sovereign 2026*_`,
		hora,
		r.CPUTemp, r.statusIconGPU(r.GPUTemp),
		r.GPUUtil, r.GPUPower, r.statusIconGPU(r.GPUTemp),
		r.GPUVRAMUsed, r.GPUVRAMTotal, vramPct, r.statusIconVRAM(vramPct),
		r.RAMUsed, r.RAMTotal, ramPct, r.statusIconRAM(ramPct),
		containers, r.formatContainerList(),
		statusIcon(r.Ollama),
		statusIcon(r.LiteLLM),
		statusIcon(r.Qdrant),
		statusIcon(r.Redis),
		r.ZFS,
		r.GrafanaURL,
	)
}

func (r Report) statusIconGPU(temp string) string {
	if temp == "N/A" {
		return "⚪"
	}
	var t int
	fmt.Sscanf(temp, "%d", &t)
	switch {
	case t < 50:
		return "🟢"
	case t < 70:
		return "🟡"
	case t < 80:
		return "🟠"
	default:
		return "🔴"
	}
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
