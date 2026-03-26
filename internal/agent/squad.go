package agent

import (
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// SquadAgent define a identidade fixa e o status de um membro do Squad Oficial.
type SquadAgent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Status   string `json:"status"` // "online", "offline", "busy"
	Load     int    `json:"load"`   // 0-100% de carga
	Color    string `json:"color"`  // Classe CSS de cor sugerida
	IconName string `json:"icon"`   // Nome do ícone Lucide
}

var (
	squadMu sync.RWMutex

	// fixedSquad mantém a lista oficial em memória do backend Go,
	// em vez de injetar estático via hardcode no frontend.
	fixedSquad = []SquadAgent{
		{
			ID:       "aurelia",
			Name:     "Aurelia_Code",
			Role:     "Soberana · Comando · Orquestracao",
			Status:   "online",
			Load:     12,
			Color:    "text-purple-400",
			IconName: "Crown",
		},
		{
			ID:       "sentinel",
			Name:     "Sentinel",
			Role:     "Monitor de Homelab",
			Status:   "online",
			Load:     5,
			Color:    "text-cyan-400",
			IconName: "Activity",
		},
		{
			ID:       "cronus",
			Name:     "Cronus",
			Role:     "Orquestrador de Crons",
			Status:   "online",
			Load:     3,
			Color:    "text-yellow-400",
			IconName: "Clock",
		},
		{
			ID:       "gemma",
			Name:     "Gemma3",
			Role:     "Executor Local · Ollama",
			Status:   "online",
			Load:     0,
			Color:    "text-green-400",
			IconName: "Cpu",
		},
		{
			ID:       "openrouter",
			Name:     "OpenRouter",
			Role:     "Executor Remoto · MiniMax",
			Status:   "online",
			Load:     0,
			Color:    "text-blue-400",
			IconName: "Globe",
		},
	}
)

// GetFixedSquad retorna a cópia atualizada do Squad.
func GetFixedSquad() []SquadAgent {
	squadMu.RLock()
	defer squadMu.RUnlock()

	copySquad := make([]SquadAgent, len(fixedSquad))
	copy(copySquad, fixedSquad)
	return copySquad
}

// UpdateSquadAgentStatus atualiza o status e workload de um agente fixo em memória.
// Usado pelo loop ou master_team_service para refletir atividade real.
func UpdateSquadAgentStatus(id string, status string, load int) {
	squadMu.Lock()
	defer squadMu.Unlock()

	for i := range fixedSquad {
		if fixedSquad[i].ID == id {
			fixedSquad[i].Status = status
			fixedSquad[i].Load = load
			break
		}
	}
}

// CronActiveCounter is implemented by cron.Scheduler to expose active job count.
type CronActiveCounter interface {
	ActiveCount() int
}

// StartLiveLoad launches a background goroutine that polls real metrics every 10s
// and updates squad agent statuses accordingly.
func StartLiveLoad(cronScheduler CronActiveCounter, ollamaURL, openrouterKey string) {
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateLiveLoad(cronScheduler, ollamaURL, openrouterKey)
		}
	}()
	// Run once immediately (non-blocking)
	go updateLiveLoad(cronScheduler, ollamaURL, openrouterKey)
}

func updateLiveLoad(cronScheduler CronActiveCounter, ollamaURL, openrouterKey string) {
	// aurelia: goroutine count mapped 0-100 (max ~200 goroutines = 100%)
	numG := runtime.NumGoroutine()
	aureliaLoad := numG * 100 / 200
	if aureliaLoad > 100 {
		aureliaLoad = 100
	}
	UpdateSquadAgentStatus("aurelia", "online", aureliaLoad)

	// cronus: active cron job count (0-100, scaled)
	if cronScheduler != nil {
		count := cronScheduler.ActiveCount()
		cronusLoad := count * 10 // 10 jobs = 100%
		if cronusLoad > 100 {
			cronusLoad = 100
		}
		UpdateSquadAgentStatus("cronus", "online", cronusLoad)
	}

	// gemma: probe Ollama /api/tags
	gemmaStatus, gemmaLoad := probeOllama(ollamaURL)
	UpdateSquadAgentStatus("gemma", gemmaStatus, gemmaLoad)

	// openrouter: probe OpenRouter /api/v1/models
	orStatus, orLoad := probeOpenRouter(openrouterKey)
	UpdateSquadAgentStatus("openrouter", orStatus, orLoad)

	// sentinel: always online, load = goroutine count / 3
	sentinelLoad := numG / 3
	if sentinelLoad > 100 {
		sentinelLoad = 100
	}
	UpdateSquadAgentStatus("sentinel", "online", sentinelLoad)
}

func probeOllama(ollamaURL string) (string, int) {
	client := &http.Client{Timeout: 3 * time.Second}
	start := time.Now()
	resp, err := client.Get(ollamaURL + "/api/tags")
	if err != nil {
		return "offline", 0
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "offline", 0
	}
	latencyMs := time.Since(start).Milliseconds()
	// load: fast response = low load; >2000ms = 100% load
	load := int(latencyMs / 20)
	if load > 100 {
		load = 100
	}
	return "online", load
}

func probeOpenRouter(apiKey string) (string, int) {
	if apiKey == "" {
		return "offline", 0
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		return "offline", 0
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "offline", 0
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "offline", 0
	}
	return "online", 5
}

// BotConfigEntry is a minimal interface to avoid import cycles with config package.
// Callers pass config.BotConfig slices cast through this helper.
type BotConfigEntry struct {
	ID        string
	Name      string
	PersonaID string
	FocusArea string
}

// SyncBotsToSquad adds bots from the multi-bot pool as squad agents if they don't exist.
// Icon and color are derived from PersonaID.
func SyncBotsToSquad(bots []BotConfigEntry) {
	for _, b := range bots {
		icon, color := botPersonaIcon(b.PersonaID)
		role := b.FocusArea
		if role == "" {
			role = "Bot Telegram"
		}
		AddSquadAgent(SquadAgent{
			ID:       b.ID,
			Name:     b.Name,
			Role:     role,
			Status:   "online",
			Load:     0,
			Color:    color,
			IconName: icon,
		})
	}
}

func botPersonaIcon(personaID string) (icon, color string) {
	switch personaID {
	case "aurelia-sovereign":
		return "Crown", "text-purple-400"
	case "aurelia-leader":
		return "Crown", "text-purple-400"
	case "hvac-sales":
		return "Thermometer", "text-blue-400"
	case "project-manager":
		return "ClipboardCheck", "text-yellow-400"
	case "life-organizer":
		return "Calendar", "text-green-400"
	case "secretaria-caixa":
		return "Briefcase", "text-blue-500"
	case "data-governance":
		return "Database", "text-cyan-400"
	case "homelab-ops":
		return "Server", "text-orange-400"
	default:
		return "Bot", "text-white/60"
	}
}

// AddSquadAgent registra um novo agente no squad se não existir com esse ID.
// Permite dinâmicamente adicionar agentes spawned pelo swarm.
func AddSquadAgent(a SquadAgent) {
	squadMu.Lock()
	defer squadMu.Unlock()

	// Verificar se já existe
	for _, s := range fixedSquad {
		if s.ID == a.ID {
			return
		}
	}

	// Adicionar novo agente
	fixedSquad = append(fixedSquad, a)
}
