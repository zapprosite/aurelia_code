package telegram

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/media"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/skill"
	"github.com/kocar/aurelia/pkg/stt"
	"github.com/kocar/aurelia/pkg/tts"
)

// HealthReporter exposes gateway diagnostics as JSON.
type HealthReporter interface {
	GatewayStatusJSON() ([]byte, error)
}

// BotController wires Telegram I/O to the application services.
type BotController struct {
	bot              *telebot.Bot
	botID            string   // S-32: per-bot identity for multi-bot pool
	allowedUserIDs   []int64  // S-32: per-bot override; if nil, uses config.TelegramAllowedUserIDs
	config           *config.AppConfig
	memory           *memory.MemoryManager
	router           *skill.Router
	executor         *skill.Executor
	loader           *skill.Loader
	stt              stt.Transcriber
	tts              tts.Synthesizer
	premiumTTS       tts.Synthesizer
	canonical        *persona.CanonicalIdentityService
	bootstrapMu      sync.Mutex
	pendingBootstrap map[int64]bootstrapState
	albumMu          sync.Mutex
	pendingAlbums    map[string]*pendingAlbum
	mediaMu          sync.Mutex
	recentMedia      map[string]recentMedia
	personasDir      string
	healthReporter   HealthReporter
	inputGuard       *InputGuard
	// S-27: Squad and Cron status reporters for /status command
	squadReporter   SquadStatusReporter
	cronJobReporter CronNextJobReporter
	mediaProcessor  *media.Processor
}

// BotID returns the identifier of this bot instance.
func (bc *BotController) BotID() string {
	return bc.botID
}

// SetHealthReporter wires a gateway health reporter for /status diagnostics.
func (bc *BotController) SetHealthReporter(hr HealthReporter) {
	bc.healthReporter = hr
}

// SetInputGuard wires the prompt injection guard.
func (bc *BotController) SetInputGuard(g *InputGuard) {
	bc.inputGuard = g
}

type pendingAlbum struct {
	ownerMessageID int
	caption        string
	photos         []albumPhoto
}

type albumPhoto struct {
	messageID int
	photo     telebot.Photo
}

type recentMedia struct {
	parts     []agent.ContentPart
	updatedAt time.Time
}

// NewBotController builds the Telegram controller.
func NewBotController(
	cfg *config.AppConfig,
	mem *memory.MemoryManager,
	r *skill.Router,
	e *skill.Executor,
	l *skill.Loader,
	s stt.Transcriber,
	canonical *persona.CanonicalIdentityService,
	personasDir string,
) (*BotController, error) {

	pref := telebot.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	if os.Getenv("RUN_SWARM_E2E") != "" {
		pref.Offline = true
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bc := &BotController{
		bot:              b,
		config:           cfg,
		memory:           mem,
		router:           r,
		executor:         e,
		loader:           l,
		stt:              s,
		tts:              tts.NewDefaultSynthesizer(cfg),
		premiumTTS:       tts.NewPremiumSynthesizer(cfg),
		canonical:        canonical,
		pendingBootstrap: make(map[int64]bootstrapState),
		pendingAlbums:    make(map[string]*pendingAlbum),
		recentMedia:      make(map[string]recentMedia),
		personasDir:      personasDir,
		mediaProcessor:   media.NewProcessor(s, e.GetLoop().GetLLMProvider(), ""),
	}

	bc.setupRoutes()
	return bc, nil
}

// NewBotControllerForBot builds a BotController from a BotConfig entry (S-32 multi-bot).
// The token and allowed users come from botCfg; shared services come from appCfg.
func NewBotControllerForBot(
	appCfg *config.AppConfig,
	botCfg config.BotConfig,
	mem *memory.MemoryManager,
	r *skill.Router,
	e *skill.Executor,
	l *skill.Loader,
	s stt.Transcriber,
	canonical *persona.CanonicalIdentityService,
	personasDir string,
) (*BotController, error) {
	token := botCfg.Token
	if token == "" {
		token = appCfg.TelegramBotToken
	}

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	if os.Getenv("RUN_SWARM_E2E") != "" {
		pref.Offline = true
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot %q: %w", botCfg.ID, err)
	}

	allowedIDs := botCfg.AllowedUserIDs
	if len(allowedIDs) == 0 {
		allowedIDs = appCfg.TelegramAllowedUserIDs
	}

	bc := &BotController{
		bot:              b,
		botID:            botCfg.ID,
		allowedUserIDs:   allowedIDs,
		config:           appCfg,
		memory:           mem,
		router:           r,
		executor:         e,
		loader:           l,
		stt:              s,
		tts:              tts.NewDefaultSynthesizer(appCfg),
		premiumTTS:       tts.NewPremiumSynthesizer(appCfg),
		canonical:        canonical,
		pendingBootstrap: make(map[int64]bootstrapState),
		pendingAlbums:    make(map[string]*pendingAlbum),
		recentMedia:      make(map[string]recentMedia),
		personasDir:      personasDir,
		mediaProcessor:   media.NewProcessor(s, e.GetLoop().GetLLMProvider(), ""),
	}

	bc.setupRoutes()
	return bc, nil
}

// GetBot exposes the underlying Telebot instance.
func (bc *BotController) GetBot() *telebot.Bot {
	return bc.bot
}

// Start begins Telegram polling.
func (bc *BotController) Start() {
	observability.Logger("telegram.bot").Info("starting Aurelia Telegram bot", slog.Int("allowed_users", len(bc.config.TelegramAllowedUserIDs)))
	bc.bot.Start()
}

// Stop ends Telegram polling.
func (bc *BotController) Stop() {
	bc.bot.Stop()
}

func (bc *BotController) isAllowedUser(userID int64) bool {
	if bc == nil || bc.config == nil {
		return false
	}
	// Per-bot override takes precedence over global config.
	list := bc.allowedUserIDs
	if len(list) == 0 {
		list = bc.config.TelegramAllowedUserIDs
	}
	for _, id := range list {
		if id == userID {
			return true
		}
	}
	return false
}

func (bc *BotController) setupRoutes() {
	bc.bot.Use(bc.whitelistMiddleware())

	bc.setupBootstrapRoutes()
	bc.registerContentRoutes()
}
