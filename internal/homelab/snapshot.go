package homelab

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	changeLogPath       = "/srv/ops/ai-governance/logs/CHANGE_LOG.txt"
	defaultDiskPath     = "/srv"
	defaultZpoolName    = "tank"
	commandTimeout      = 5 * time.Second
	probeTimeout        = 2 * time.Second
	maxSnapshots        = 15
	maxChangeLogLines   = 20
	maxZpoolStatusLines = 10

	// service_S_ubabase removido (migrado para SQLite)
	homelabServiceQdrant   = "qdrant"
	homelabServiceCaprover = "caprover"

	homelabStatusHealthy  = "healthy"
	homelabStatusDegraded = "degraded"
	homelabStatusOffline  = "offline"
)

type Snapshot struct {
	Containers  []ContainerStatus `json:"containers"`
	Services    []HomelabService  `json:"services"`
	Summary     HomelabSummary    `json:"summary"`
	Counts      HomelabCounts     `json:"counts"`
	Health      []HealthEndpoint  `json:"health"`
	Snapshots   []ZFSSnapshot     `json:"snapshots"`
	Changelog   string            `json:"changelog"`
	AgentState  AgentState        `json:"agent_state"`
	ZpoolStatus string            `json:"zpool_status"`
	DiskUsage   DiskUsage         `json:"disk_usage"`
	Timestamp   time.Time         `json:"timestamp"`
}

type ContainerStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Ports  string `json:"ports"`
	Up     bool   `json:"up"`
}

type HomelabService struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Label      string            `json:"label"`
	Status     string            `json:"status"`
	Summary    string            `json:"summary"`
	Detail     string            `json:"detail"`
	Counts     ServiceCounts     `json:"counts"`
	Containers []ContainerStatus `json:"containers"`
}

type HomelabSummary struct {
	Status   string `json:"status"`
	Healthy  int    `json:"healthy"`
	Degraded int    `json:"degraded"`
	Offline  int    `json:"offline"`
}

type ServiceCounts struct {
	Total      int `json:"total"`
	Up         int `json:"up"`
	Restarting int `json:"restarting"`
	Exited     int `json:"exited"`
	Dead       int `json:"dead"`
	Other      int `json:"other"`
}

type HomelabCounts struct {
	Services             int `json:"services"`
	Containers           int `json:"containers"`
	UpContainers         int `json:"up_containers"`
	RestartingContainers int `json:"restarting_containers"`
	ExitedContainers     int `json:"exited_containers"`
	DeadContainers       int `json:"dead_containers"`
	OtherContainers      int `json:"other_containers"`
}

type HealthEndpoint struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	OK   bool   `json:"ok"`
	Code int    `json:"code"`
}

type ZFSSnapshot struct {
	Name     string `json:"name"`
	Creation string `json:"creation"`
}

type AgentState struct {
	LastAction string `json:"lastAction"`
	Timestamp  string `json:"timestamp"`
	Status     string `json:"status"`
	Agent      string `json:"agent"`
	Model      string `json:"model"`
}

type DiskUsage struct {
	Raw         string `json:"raw"`
	Filesystem  string `json:"filesystem"`
	Size        string `json:"size"`
	Used        string `json:"used"`
	Available   string `json:"available"`
	UsedPercent int    `json:"used_percent"`
	Mount       string `json:"mount"`
}

type commandRunner func(ctx context.Context, name string, args ...string) (string, error)

type Collector struct {
	run             commandRunner
	readFile        func(string) ([]byte, error)
	now             func() time.Time
	client          *http.Client
	healthTargets   []HealthEndpoint
	agentStatePaths []string
	changeLogPath   string
	diskPath        string
	zpoolName       string
}

func NewCollector() *Collector {
	homeDir, _ := os.UserHomeDir()
	agentStatePaths := []string{
		"/srv/ops/ai-governance/state/agent_state.json",
	}
	if homeDir != "" {
		agentStatePaths = append([]string{
			filepath.Join(homeDir, "vrv-dashboard", "state", "agent_state.json"),
		}, agentStatePaths...)
	}

	return &Collector{
		run:      runCommand,
		readFile: os.ReadFile,
		now:      time.Now,
		client:   &http.Client{Timeout: probeTimeout},
		healthTargets: []HealthEndpoint{
			{Name: "Qdrant", URL: "http://localhost:6333/healthz"},
			{Name: "n8n", URL: "http://localhost:5678/healthz"},
		},
		agentStatePaths: agentStatePaths,
		changeLogPath:   changeLogPath,
		diskPath:        defaultDiskPath,
		zpoolName:       defaultZpoolName,
	}
}

func (c *Collector) Collect(ctx context.Context) Snapshot {
	containers := c.collectContainers(ctx)
	services, summary, counts := buildHomelabServices(containers)
	return Snapshot{
		Containers:  containers,
		Services:    services,
		Summary:     summary,
		Counts:      counts,
		Health:      c.collectHealth(ctx),
		Snapshots:   c.collectSnapshots(ctx),
		Changelog:   c.collectChangelog(),
		AgentState:  c.collectAgentState(),
		ZpoolStatus: c.collectZpoolStatus(ctx),
		DiskUsage:   c.collectDiskUsage(ctx),
		Timestamp:   c.now().UTC(),
	}
}

func (c *Collector) collectContainers(ctx context.Context) []ContainerStatus {
	raw, err := c.run(ctx, "docker", "ps", "-a", "--format", "{{.Names}}\t{{.Status}}\t{{.Ports}}")
	if err != nil {
		return nil
	}
	return parseContainerStatuses(raw)
}

func buildHomelabServices(containers []ContainerStatus) ([]HomelabService, HomelabSummary, HomelabCounts) {
	specs := []struct {
		name  string
		label string
		match func(ContainerStatus) bool
	}{
		// Antigo_DB removido (migrado para SQLite)
		{
			name:  homelabServiceQdrant,
			label: "Qdrant",
			match: func(container ContainerStatus) bool {
				return container.Name == "qdrant"
			},
		},
		{
			name:  homelabServiceCaprover,
			label: "CapRover",
			match: func(container ContainerStatus) bool {
				return strings.HasPrefix(container.Name, "captain-")
			},
		},
	}

	services := make([]HomelabService, 0, len(specs))
	serviceSummary := HomelabSummary{}
	containerCounts := HomelabCounts{
		Services: len(specs),
	}

	for _, container := range containers {
		containerCounts.Containers++
		switch classifyContainerState(container.Status) {
		case homelabStatusHealthy:
			containerCounts.UpContainers++
		case "restarting":
			containerCounts.RestartingContainers++
		case "exited":
			containerCounts.ExitedContainers++
		case "dead":
			containerCounts.DeadContainers++
		default:
			containerCounts.OtherContainers++
		}
	}

	for _, spec := range specs {
		group := make([]ContainerStatus, 0)
		for _, container := range containers {
			if spec.match(container) {
				group = append(group, container)
			}
		}

		serviceCounts := ServiceCounts{
			Total: len(group),
		}
		healthy := 0
		for _, container := range group {
			switch classifyContainerState(container.Status) {
			case homelabStatusHealthy:
				healthy++
				serviceCounts.Up++
			case "restarting":
				serviceCounts.Restarting++
			case "exited":
				serviceCounts.Exited++
			case "dead":
				serviceCounts.Dead++
			default:
				serviceCounts.Other++
			}
		}

		status := homelabStatusOffline
		switch {
		case len(group) == 0 || healthy == 0:
			status = homelabStatusOffline
		case healthy == len(group):
			status = homelabStatusHealthy
		default:
			status = homelabStatusDegraded
		}
		if status == homelabStatusHealthy {
			serviceSummary.Healthy++
		} else if status == homelabStatusDegraded {
			serviceSummary.Degraded++
		} else {
			serviceSummary.Offline++
		}

		services = append(services, HomelabService{
			ID:         spec.name,
			Name:       spec.name,
			Label:      spec.label,
			Status:     status,
			Summary:    serviceSummaryText(spec.label, status),
			Detail:     serviceDetailText(serviceCounts),
			Counts:     serviceCounts,
			Containers: group,
		})
	}

	serviceSummary.Status = deriveHomelabStatus(serviceSummary.Healthy, serviceSummary.Degraded, serviceSummary.Offline)
	return services, serviceSummary, containerCounts
}

func serviceSummaryText(label string, status string) string {
	switch status {
	case homelabStatusHealthy:
		return label + " operacional."
	case homelabStatusDegraded:
		return label + " parcialmente operacional; requer atenção."
	default:
		return label + " offline ou sem containers saudáveis."
	}
}

func serviceDetailText(counts ServiceCounts) string {
	parts := make([]string, 0, 4)
	if counts.Total > 0 {
		parts = append(parts, strconv.Itoa(counts.Up)+"/"+strconv.Itoa(counts.Total)+" up")
	}
	if counts.Restarting > 0 {
		parts = append(parts, strconv.Itoa(counts.Restarting)+" restarting")
	}
	if counts.Exited > 0 {
		parts = append(parts, strconv.Itoa(counts.Exited)+" exited")
	}
	if counts.Dead > 0 {
		parts = append(parts, strconv.Itoa(counts.Dead)+" dead")
	}
	if counts.Other > 0 {
		parts = append(parts, strconv.Itoa(counts.Other)+" other")
	}
	if len(parts) == 0 {
		return "nenhum container encontrado"
	}
	return strings.Join(parts, " · ")
}

func deriveHomelabStatus(healthy, degraded, offline int) string {
	switch {
	case healthy > 0 && degraded == 0 && offline == 0:
		return homelabStatusHealthy
	case healthy == 0 && degraded == 0 && offline > 0:
		return homelabStatusOffline
	default:
		return homelabStatusDegraded
	}
}

func classifyContainerState(status string) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	switch {
	case strings.HasPrefix(normalized, "up"):
		return homelabStatusHealthy
	case strings.HasPrefix(normalized, "restarting"):
		return "restarting"
	case strings.HasPrefix(normalized, "exited"):
		return "exited"
	case strings.HasPrefix(normalized, "dead"):
		return "dead"
	default:
		return "other"
	}
}

func (c *Collector) collectHealth(ctx context.Context) []HealthEndpoint {
	items := make([]HealthEndpoint, 0, len(c.healthTargets))
	for _, target := range c.healthTargets {
		endpoint := target
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.URL, nil)
		if err != nil {
			items = append(items, endpoint)
			continue
		}
		resp, err := c.client.Do(req)
		if err != nil {
			items = append(items, endpoint)
			continue
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		endpoint.Code = resp.StatusCode
		endpoint.OK = resp.StatusCode >= 200 && resp.StatusCode < 300
		items = append(items, endpoint)
	}
	return items
}

func (c *Collector) collectSnapshots(ctx context.Context) []ZFSSnapshot {
	raw, err := c.run(ctx, "zfs", "list", "-H", "-t", "snapshot", "-o", "name,creation", "-s", "creation")
	if err != nil {
		return nil
	}
	return parseZFSSnapshots(raw)
}

func (c *Collector) collectChangelog() string {
	body, err := c.readFile(c.changeLogPath)
	if err != nil {
		return "Nenhum changelog encontrado"
	}
	lines := strings.Split(strings.ReplaceAll(string(body), "\r\n", "\n"), "\n")
	return tailLines(lines, maxChangeLogLines)
}

func (c *Collector) collectAgentState() AgentState {
	for _, path := range c.agentStatePaths {
		body, err := c.readFile(path)
		if err != nil {
			continue
		}
		var state AgentState
		if err := json.Unmarshal(body, &state); err == nil {
			return state
		}
	}
	return AgentState{
		LastAction: "unknown",
		Timestamp:  "-",
		Status:     "unknown",
		Agent:      "-",
		Model:      "-",
	}
}

func (c *Collector) collectZpoolStatus(ctx context.Context) string {
	raw, err := c.run(ctx, "zpool", "status", c.zpoolName)
	if err != nil {
		return "ZFS pool não disponível"
	}
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	return tailLines(lines[:min(len(lines), maxZpoolStatusLines)], maxZpoolStatusLines)
}

func (c *Collector) collectDiskUsage(ctx context.Context) DiskUsage {
	raw, err := c.run(ctx, "df", "-h", "--output=source,size,used,avail,pcent,target", c.diskPath)
	if err != nil {
		return DiskUsage{Raw: "Não disponível"}
	}
	return parseDiskUsage(raw)
}

func parseContainerStatuses(raw string) []ContainerStatus {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	items := make([]ContainerStatus, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 3)
		name := strings.TrimSpace(parts[0])
		status := ""
		ports := "-"
		if len(parts) > 1 {
			status = strings.TrimSpace(parts[1])
		}
		if len(parts) > 2 && strings.TrimSpace(parts[2]) != "" {
			ports = strings.TrimSpace(parts[2])
		}
		items = append(items, ContainerStatus{
			Name:   name,
			Status: status,
			Ports:  ports,
			Up:     strings.HasPrefix(strings.ToLower(status), "up"),
		})
	}
	return items
}

func parseZFSSnapshots(raw string) []ZFSSnapshot {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		filtered = append(filtered, line)
	}
	if len(filtered) > maxSnapshots {
		filtered = filtered[len(filtered)-maxSnapshots:]
	}
	items := make([]ZFSSnapshot, 0, len(filtered))
	for _, line := range filtered {
		parts := strings.SplitN(line, "\t", 2)
		name := strings.TrimSpace(parts[0])
		creation := ""
		if len(parts) > 1 {
			creation = strings.TrimSpace(parts[1])
		}
		items = append(items, ZFSSnapshot{Name: name, Creation: creation})
	}
	return items
}

func parseDiskUsage(raw string) DiskUsage {
	lines := strings.Split(strings.TrimSpace(strings.ReplaceAll(raw, "\r\n", "\n")), "\n")
	if len(lines) < 2 {
		return DiskUsage{Raw: strings.TrimSpace(raw)}
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 6 {
		return DiskUsage{Raw: strings.TrimSpace(lines[1])}
	}
	usedPercent, _ := strconv.Atoi(strings.TrimSuffix(fields[4], "%"))
	return DiskUsage{
		Raw:         strings.TrimSpace(lines[1]),
		Filesystem:  fields[0],
		Size:        fields[1],
		Used:        fields[2],
		Available:   fields[3],
		UsedPercent: usedPercent,
		Mount:       fields[5],
	}
}

func tailLines(lines []string, count int) string {
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		filtered = append(filtered, strings.TrimRight(line, " \t"))
	}
	if len(filtered) == 0 {
		return ""
	}
	if len(filtered) > count {
		filtered = filtered[len(filtered)-count:]
	}
	return strings.Join(filtered, "\n")
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	runCtx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	out, err := exec.CommandContext(runCtx, name, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
