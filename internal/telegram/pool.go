package telegram

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/skill"
	"github.com/kocar/aurelia/pkg/stt"
)

// BotPool manages the lifecycle of multiple BotControllers.
type BotPool struct {
	mu          sync.RWMutex
	bots        map[string]*BotController
	configs     map[string]config.BotConfig
	appCfg      *config.AppConfig
	mem         *memory.MemoryManager
	router      *skill.Router
	executor    *skill.Executor
	loader      *skill.Loader
	transcriber stt.Transcriber
	canonical   *persona.CanonicalIdentityService
	personasDir string
}

// NewBotPool creates an empty BotPool with shared dependencies.
func NewBotPool(
	appCfg *config.AppConfig,
	mem *memory.MemoryManager,
	router *skill.Router,
	executor *skill.Executor,
	loader *skill.Loader,
	transcriber stt.Transcriber,
	canonical *persona.CanonicalIdentityService,
	personasDir string,
) *BotPool {
	return &BotPool{
		bots:        make(map[string]*BotController),
		configs:     make(map[string]config.BotConfig),
		appCfg:      appCfg,
		mem:         mem,
		router:      router,
		executor:    executor,
		loader:      loader,
		transcriber: transcriber,
		canonical:   canonical,
		personasDir: personasDir,
	}
}

// Add creates and registers a BotController for the given BotConfig.
// Returns an error if the token is invalid or the ID already exists.
func (p *BotPool) Add(botCfg config.BotConfig) error {
	if !botCfg.Enabled {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.bots[botCfg.ID]; exists {
		return fmt.Errorf("bot %q already registered", botCfg.ID)
	}

	bc, err := NewBotControllerForBot(
		p.appCfg, botCfg,
		p.mem, p.router, p.executor, p.loader,
		p.transcriber, p.canonical, p.personasDir,
	)
	if err != nil {
		return fmt.Errorf("create bot controller for %q: %w", botCfg.ID, err)
	}

	p.bots[botCfg.ID] = bc
	p.configs[botCfg.ID] = botCfg
	return nil
}

// Remove stops and removes a bot by ID.
func (p *BotPool) Remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if bc, ok := p.bots[id]; ok {
		bc.Stop()
		delete(p.bots, id)
		delete(p.configs, id)
	}
}

// Get returns the BotController for the given ID, or nil if not found.
func (p *BotPool) Get(id string) *BotController {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.bots[id]
}

// Primary returns the first enabled bot (canonical primary bot), or nil.
func (p *BotPool) Primary() *BotController {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Prefer "aurelia" as primary
	if bc, ok := p.bots["aurelia"]; ok {
		return bc
	}
	// Fallback: return any bot
	for _, bc := range p.bots {
		return bc
	}
	return nil
}

// All returns a snapshot of all registered BotControllers.
func (p *BotPool) All() map[string]*BotController {
	p.mu.RLock()
	defer p.mu.RUnlock()

	snapshot := make(map[string]*BotController, len(p.bots))
	for id, bc := range p.bots {
		snapshot[id] = bc
	}
	return snapshot
}

// Configs returns a snapshot of all BotConfig entries.
func (p *BotPool) Configs() []config.BotConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]config.BotConfig, 0, len(p.configs))
	for _, cfg := range p.configs {
		out = append(out, cfg)
	}
	return out
}

// StartAll launches each bot's Telegram polling in its own goroutine.
func (p *BotPool) StartAll() {
	logger := observability.Logger("telegram.pool")
	p.mu.RLock()
	defer p.mu.RUnlock()

	for id, bc := range p.bots {
		logger.Info("starting bot", slog.String("bot_id", id))
		go bc.Start()
	}
}

// StopAll stops all bots gracefully.
func (p *BotPool) StopAll() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, bc := range p.bots {
		bc.Stop()
	}
}

// AddController registers a pre-built BotController under the given ID.
// Used for the primary bot fallback when no Bots config is present.
func (p *BotPool) AddController(id string, bc *BotController) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.bots[id]; exists {
		return fmt.Errorf("bot %q already registered", id)
	}
	p.bots[id] = bc
	p.configs[id] = config.BotConfig{ID: id, Name: id, Enabled: true}
	return nil
}

// Size returns the number of registered bots.
func (p *BotPool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.bots)
}
