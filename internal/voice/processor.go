package voice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/health"
	"github.com/kocar/aurelia/pkg/stt"
)

const voiceCooldownOn429 = 15 * time.Minute

type Dispatcher interface {
	DispatchVoice(ctx context.Context, userID, chatID int64, text string, requiresAudio bool) error
}

type Mirror interface {
	MirrorTranscript(ctx context.Context, event TranscriptEvent) error
}

type TranscriptEvent struct {
	JobID       string
	UserID      int64
	ChatID      int64
	Source      string
	Transcript  string
	Accepted    bool
	RequiresTTS bool
	CreatedAt   time.Time
}

type Config struct {
	PollInterval         time.Duration
	HeartbeatPath        string
	HeartbeatFreshness   time.Duration
	WakePhrase           string
	DefaultUserID        int64
	DefaultChatID        int64
	SoftCapDaily         int
	HardCapDaily         int
	PrimaryLabel         string
	CooldownOn429        time.Duration
	Mirror               Mirror
	MaxHealthyQueueDepth int
}

type Processor struct {
	spool      *Spool
	primary    stt.Transcriber
	fallback   stt.Transcriber
	dispatcher Dispatcher
	cfg        Config
	metrics    *voiceMetrics

	mu          sync.RWMutex
	stopChan    chan struct{}
	status      HeartbeatStatus
	budgetPath  string
	lastStarted time.Time
}

type HeartbeatStatus struct {
	Status          string    `json:"status"`
	LastBeatAt      time.Time `json:"last_beat_at"`
	LastProcessedAt time.Time `json:"last_processed_at,omitempty"`
	LastError       string    `json:"last_error,omitempty"`
	QueueDepth      int       `json:"queue_depth"`
	ProcessedJobs   int       `json:"processed_jobs"`
	FailedJobs      int       `json:"failed_jobs"`
	SoftCapReached  bool      `json:"soft_cap_reached"`
	HardCapReached  bool      `json:"hard_cap_reached"`
}

type budgetState struct {
	Day           string    `json:"day"`
	Requests      int       `json:"requests"`
	CooldownUntil time.Time `json:"cooldown_until,omitempty"`
}

type noopMirror struct{}

func (noopMirror) MirrorTranscript(context.Context, TranscriptEvent) error { return nil }

func NewProcessor(spool *Spool, primary stt.Transcriber, fallback stt.Transcriber, dispatcher Dispatcher, cfg Config) *Processor {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = time.Second
	}
	if cfg.HeartbeatFreshness <= 0 {
		cfg.HeartbeatFreshness = 45 * time.Second
	}
	if cfg.WakePhrase == "" {
		cfg.WakePhrase = "jarvis"
	}
	if cfg.PrimaryLabel == "" {
		cfg.PrimaryLabel = "groq"
	}
	if cfg.CooldownOn429 <= 0 {
		cfg.CooldownOn429 = voiceCooldownOn429
	}
	if cfg.SoftCapDaily <= 0 {
		cfg.SoftCapDaily = 800
	}
	if cfg.HardCapDaily <= 0 {
		cfg.HardCapDaily = 1200
	}
	if cfg.Mirror == nil {
		cfg.Mirror = noopMirror{}
	}
	if cfg.MaxHealthyQueueDepth <= 0 {
		cfg.MaxHealthyQueueDepth = 8
	}
	budgetPath := ""
	if spool != nil {
		budgetPath = filepath.Join(spool.Root(), "budget.json")
	}
	return &Processor{
		spool:      spool,
		primary:    primary,
		fallback:   fallback,
		dispatcher: dispatcher,
		cfg:        cfg,
		metrics:    defaultVoiceMetrics(),
		budgetPath: budgetPath,
		status: HeartbeatStatus{
			Status: "idle",
		},
	}
}

func (p *Processor) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopChan != nil {
		return nil
	}
	p.stopChan = make(chan struct{})
	p.lastStarted = time.Now().UTC()
	go p.runLoop(p.stopChan)
	return nil
}

func (p *Processor) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopChan == nil {
		return
	}
	close(p.stopChan)
	p.stopChan = nil
}

func (p *Processor) IsRunning() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stopChan != nil
}

func (p *Processor) StatusSnapshot() HeartbeatStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.status
}

func (p *Processor) HealthCheck() health.CheckResult {
	status := p.StatusSnapshot()
	if status.LastBeatAt.IsZero() {
		return health.CheckResult{Status: "error", Message: "voice processor has no heartbeat"}
	}
	if time.Since(status.LastBeatAt) > p.cfg.HeartbeatFreshness {
		return health.CheckResult{Status: "error", Message: "voice heartbeat stale"}
	}
	if status.HardCapReached {
		return health.CheckResult{Status: "error", Message: "voice stt hard cap reached"}
	}
	if status.QueueDepth > p.cfg.MaxHealthyQueueDepth {
		return health.CheckResult{Status: "warning", Message: fmt.Sprintf("voice backlog=%d", status.QueueDepth)}
	}
	if status.LastError != "" {
		return health.CheckResult{Status: "warning", Message: status.LastError}
	}
	if status.SoftCapReached {
		return health.CheckResult{Status: "warning", Message: "voice stt soft cap reached"}
	}
	return health.CheckResult{Status: "ok", Message: "voice processor healthy"}
}

func (p *Processor) StatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(p.StatusSnapshot())
	})
}

func (p *Processor) ProcessOnce(ctx context.Context) error {
	if p.spool == nil {
		return fmt.Errorf("voice spool is not configured")
	}

	claimed, err := p.spool.ClaimOldest(ctx)
	if err != nil {
		p.updateStatus(func(status *HeartbeatStatus) {
			status.Status = "error"
			status.LastError = err.Error()
		})
		p.writeHeartbeat()
		return err
	}
	if claimed == nil {
		p.updateStatus(func(status *HeartbeatStatus) {
			status.Status = "idle"
			status.LastError = ""
		})
		p.writeHeartbeat()
		return nil
	}

	transcript, accepted, err := p.handleClaimedJob(ctx, claimed)
	if err != nil {
		_ = p.spool.Fail(claimed, err)
		p.recordJobResult("failed")
		p.updateStatus(func(status *HeartbeatStatus) {
			status.Status = "error"
			status.LastError = err.Error()
			status.FailedJobs++
		})
		p.writeHeartbeat()
		return err
	}

	if err := p.spool.Complete(claimed, transcript); err != nil {
		p.recordJobResult("failed")
		p.updateStatus(func(status *HeartbeatStatus) {
			status.Status = "error"
			status.LastError = err.Error()
			status.FailedJobs++
		})
		p.writeHeartbeat()
		return err
	}

	p.updateStatus(func(status *HeartbeatStatus) {
		status.Status = "idle"
		status.LastProcessedAt = time.Now().UTC()
		status.ProcessedJobs++
		status.LastError = ""
		if !accepted {
			status.Status = "dropped"
		}
	})
	if accepted {
		p.recordJobResult("completed")
	} else {
		p.recordJobResult("dropped")
	}
	p.writeHeartbeat()
	return nil
}

func (p *Processor) runLoop(stopChan chan struct{}) {
	ticker := time.NewTicker(p.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		default:
		}
		_ = p.ProcessOnce(context.Background())
		select {
		case <-stopChan:
			return
		case <-ticker.C:
		}
	}
}

func (p *Processor) handleClaimedJob(ctx context.Context, claimed *ClaimedJob) (string, bool, error) {
	transcript, err := p.transcribeWithBudget(ctx, claimed.AudioDir)
	if err != nil {
		return "", false, err
	}
	sanitized := sanitizeTranscript(transcript)

	userID := claimed.Job.UserID
	chatID := claimed.Job.ChatID
	if userID == 0 {
		userID = p.cfg.DefaultUserID
	}
	if chatID == 0 {
		chatID = p.cfg.DefaultChatID
	}

	command, accepted := stripWakePhrase(sanitized, p.cfg.WakePhrase)
	event := TranscriptEvent{
		JobID:       claimed.Job.ID,
		UserID:      userID,
		ChatID:      chatID,
		Source:      claimed.Job.Source,
		Transcript:  sanitized,
		Accepted:    accepted,
		RequiresTTS: claimed.Job.RequiresAudio,
		CreatedAt:   claimed.Job.CreatedAt,
	}
	if err := p.cfg.Mirror.MirrorTranscript(ctx, event); err != nil {
		p.recordMirrorFailure()
		p.updateStatus(func(status *HeartbeatStatus) {
			status.LastError = "voice mirror failed: " + err.Error()
		})
	}
	if !accepted {
		p.recordDispatch("dropped")
		return sanitized, false, nil
	}
	if p.dispatcher == nil {
		p.recordDispatch("error")
		return sanitized, true, fmt.Errorf("voice dispatcher is not configured")
	}
	if userID == 0 || chatID == 0 {
		p.recordDispatch("error")
		return sanitized, true, fmt.Errorf("voice dispatcher missing default user/chat")
	}
	if err := p.dispatcher.DispatchVoice(ctx, userID, chatID, command, claimed.Job.RequiresAudio); err != nil {
		p.recordDispatch("error")
		return sanitized, true, err
	}
	p.recordDispatch("accepted")
	return sanitized, true, nil
}

func (p *Processor) transcribeWithBudget(ctx context.Context, audioPath string) (string, error) {
	if p.primary == nil && p.fallback == nil {
		return "", fmt.Errorf("no STT transcriber configured")
	}

	state, err := p.loadBudgetState()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	softReached := state.Requests >= p.cfg.SoftCapDaily
	hardReached := state.Requests >= p.cfg.HardCapDaily
	p.updateStatus(func(status *HeartbeatStatus) {
		status.SoftCapReached = softReached
		status.HardCapReached = hardReached
	})

	if hardReached || (!state.CooldownUntil.IsZero() && now.Before(state.CooldownUntil)) {
		if p.fallback != nil && p.fallback.IsAvailable() {
			if hardReached {
				p.recordFallback("hard_cap")
			} else {
				p.recordFallback("cooldown")
			}
			return p.fallback.Transcribe(ctx, audioPath)
		}
		reason := "voice stt hard cap reached"
		if !state.CooldownUntil.IsZero() && now.Before(state.CooldownUntil) {
			reason = "voice stt cooldown active"
		}
		return "", fmt.Errorf("%s", reason)
	}

	if p.primary == nil || !p.primary.IsAvailable() {
		if p.fallback != nil && p.fallback.IsAvailable() {
			p.recordFallback("primary_unavailable")
			return p.fallback.Transcribe(ctx, audioPath)
		}
		return "", fmt.Errorf("primary STT unavailable and fallback not configured")
	}

	state.Requests++
	if err := p.saveBudgetState(state); err != nil {
		return "", err
	}
	transcript, err := p.primary.Transcribe(ctx, audioPath)
	if err == nil {
		return transcript, nil
	}

	if strings.Contains(strings.ToLower(err.Error()), "429") {
		state.CooldownUntil = now.Add(p.cfg.CooldownOn429)
		_ = p.saveBudgetState(state)
		if p.fallback != nil && p.fallback.IsAvailable() {
			p.recordFallback("429")
			return p.fallback.Transcribe(ctx, audioPath)
		}
	}
	return "", err
}

func (p *Processor) writeHeartbeat() {
	if strings.TrimSpace(p.cfg.HeartbeatPath) == "" {
		return
	}
	depth := 0
	if p.spool != nil {
		if queueDepth, err := p.spool.QueueDepth(); err == nil {
			depth = queueDepth
		}
	}

	p.updateStatus(func(status *HeartbeatStatus) {
		status.LastBeatAt = time.Now().UTC()
		status.QueueDepth = depth
	})
	p.recordQueueDepth(depth)
	status := p.StatusSnapshot()

	if err := os.MkdirAll(filepath.Dir(p.cfg.HeartbeatPath), 0o700); err != nil {
		return
	}
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return
	}
	data = append(data, '\n')
	_ = os.WriteFile(p.cfg.HeartbeatPath, data, 0o600)
}

func (p *Processor) loadBudgetState() (budgetState, error) {
	now := time.Now().UTC().Format("2006-01-02")
	if strings.TrimSpace(p.budgetPath) == "" {
		return budgetState{Day: now}, nil
	}
	data, err := os.ReadFile(p.budgetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return budgetState{Day: now}, nil
		}
		return budgetState{}, fmt.Errorf("read voice budget state: %w", err)
	}

	var state budgetState
	if err := json.Unmarshal(data, &state); err != nil {
		return budgetState{}, fmt.Errorf("decode voice budget state: %w", err)
	}
	if state.Day != now {
		state = budgetState{Day: now}
	}
	p.recordBudget(state.Requests, p.cfg.HardCapDaily)
	return state, nil
}

func (p *Processor) saveBudgetState(state budgetState) error {
	if strings.TrimSpace(p.budgetPath) == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(p.budgetPath), 0o700); err != nil {
		return fmt.Errorf("create voice budget dir: %w", err)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode voice budget state: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(p.budgetPath, data, 0o600); err != nil {
		return fmt.Errorf("write voice budget state: %w", err)
	}
	return nil
}

func (p *Processor) updateStatus(fn func(*HeartbeatStatus)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	fn(&p.status)
}

var transcriptSecretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{8,}\b`),
	regexp.MustCompile(`\bgsk_[A-Za-z0-9_-]{8,}\b`),
	regexp.MustCompile(`\bhf_[A-Za-z0-9]{8,}\b`),
	regexp.MustCompile(`\bAIza[0-9A-Za-z\-_]{10,}\b`),
}

func sanitizeTranscript(input string) string {
	sanitized := strings.TrimSpace(strings.ReplaceAll(input, "\x00", ""))
	for _, pattern := range transcriptSecretPatterns {
		sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_SECRET]")
	}
	return sanitized
}

func stripWakePhrase(input, wakePhrase string) (string, bool) {
	text := strings.TrimSpace(input)
	if text == "" {
		return "", false
	}
	wakePhrase = strings.TrimSpace(strings.ToLower(wakePhrase))
	if wakePhrase == "" {
		return text, true
	}

	lower := strings.ToLower(text)
	prefixes := []string{
		wakePhrase,
		"ei " + wakePhrase,
		"hey " + wakePhrase,
		"ok " + wakePhrase,
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			trimmed := strings.TrimSpace(text[len(prefix):])
			trimmed = strings.TrimLeft(trimmed, ",:;- ")
			if trimmed == "" {
				return "", false
			}
			return trimmed, true
		}
	}
	return "", false
}

func ReadHeartbeat(path string) (HeartbeatStatus, error) {
	var status HeartbeatStatus
	data, err := os.ReadFile(path)
	if err != nil {
		return status, err
	}
	if err := json.Unmarshal(data, &status); err != nil {
		return status, err
	}
	return status, nil
}

func formatBudgetMessage(requests, hard int) string {
	return strconv.Itoa(requests) + "/" + strconv.Itoa(hard)
}

func (p *Processor) recordQueueDepth(depth int) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.queueDepth.Set(float64(depth))
}

func (p *Processor) recordJobResult(result string) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.processed.WithLabelValues(result).Inc()
}

func (p *Processor) recordDispatch(result string) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.dispatches.WithLabelValues(result).Inc()
}

func (p *Processor) recordFallback(reason string) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.fallbacks.WithLabelValues(reason).Inc()
}

func (p *Processor) recordMirrorFailure() {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.mirrorFailure.Inc()
}

func (p *Processor) recordBudget(requests, hard int) {
	if p == nil || p.metrics == nil || hard <= 0 {
		return
	}
	p.metrics.budgetUsage.Set(float64(requests) / float64(hard))
}
