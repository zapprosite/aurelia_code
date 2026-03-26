package voice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/health"
)

type CaptureSource interface {
	Capture(ctx context.Context) (*CaptureEvent, error)
	Description() string
}

type CaptureEvent struct {
	Detected          bool   `json:"detected"`
	AudioFile         string `json:"audio_file"`
	UserID            int64  `json:"user_id,omitempty"`
	ChatID            int64  `json:"chat_id,omitempty"`
	RequiresAudio     bool   `json:"requires_audio,omitempty"`
	Source            string `json:"source,omitempty"`
	DeleteSourceAfter bool   `json:"delete_source_after,omitempty"`
}

type CommandCaptureSource struct {
	command string
	env     []string
}

func NewCommandCaptureSource(command string, env map[string]string) *CommandCaptureSource {
	source := &CommandCaptureSource{command: strings.TrimSpace(command)}
	for key, value := range env {
		if strings.TrimSpace(key) == "" {
			continue
		}
		source.env = append(source.env, key+"="+value)
	}
	return source
}

func (s *CommandCaptureSource) Capture(ctx context.Context) (*CaptureEvent, error) {
	if !s.IsAvailable() {
		return nil, fmt.Errorf("voice capture command not configured")
	}

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", s.command)
	cmd.Env = append(os.Environ(), s.env...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	trimmed := strings.TrimSpace(string(output))
	stderrText := strings.TrimSpace(stderr.String())
	if err != nil {
		detail := stderrText
		if detail == "" {
			detail = trimmed
		}
		return nil, fmt.Errorf("voice capture command failed: %w: %s", err, detail)
	}
	if trimmed == "" {
		return nil, nil
	}

	var event CaptureEvent
	if err := json.Unmarshal([]byte(trimmed), &event); err != nil {
		return nil, fmt.Errorf("voice capture command returned invalid json: %w", err)
	}
	if !event.Detected {
		return nil, nil
	}
	event.AudioFile = strings.TrimSpace(event.AudioFile)
	if event.AudioFile == "" {
		return nil, fmt.Errorf("voice capture command returned detected event without audio_file")
	}
	info, err := os.Stat(event.AudioFile)
	if err != nil {
		return nil, fmt.Errorf("voice capture audio file unavailable: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("voice capture audio file must be a file")
	}
	if strings.TrimSpace(event.Source) == "" {
		event.Source = "capture"
	}
	return &event, nil
}

func (s *CommandCaptureSource) Description() string {
	if s == nil || strings.TrimSpace(s.command) == "" {
		return "capture command not configured"
	}
	return "command"
}

func (s *CommandCaptureSource) IsAvailable() bool {
	return s != nil && strings.TrimSpace(s.command) != ""
}

func MissingCommandPath(command string) string {
	fields := strings.Fields(strings.TrimSpace(command))
	if len(fields) == 0 {
		return ""
	}

	head := fields[0]
	if head == "" {
		return ""
	}
	if !strings.Contains(head, "/") && !strings.HasPrefix(head, ".") {
		return ""
	}
	if _, err := os.Stat(head); err != nil {
		return head
	}
	return ""
}

type CaptureConfig struct {
	PollInterval         time.Duration
	HeartbeatPath        string
	HeartbeatFreshness   time.Duration
	DefaultUserID        int64
	DefaultChatID        int64
	DefaultSource        string
	MaxHealthyQueueDepth int
}

type CaptureStatus struct {
	Status         string    `json:"status"`
	LastBeatAt     time.Time `json:"last_beat_at"`
	LastCapturedAt time.Time `json:"last_captured_at,omitempty"`
	LastError      string    `json:"last_error,omitempty"`
	QueueDepth     int       `json:"queue_depth"`
	CapturedJobs   int       `json:"captured_jobs"`
	RejectedClips  int       `json:"rejected_clips"`
}

type CaptureWorker struct {
	spool   *Spool
	source  CaptureSource
	cfg     CaptureConfig
	metrics *voiceMetrics

	mu       sync.RWMutex
	stopChan chan struct{}
	status   CaptureStatus
}

func NewCaptureWorker(spool *Spool, source CaptureSource, cfg CaptureConfig) *CaptureWorker {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = time.Second
	}
	if cfg.HeartbeatFreshness <= 0 {
		cfg.HeartbeatFreshness = 45 * time.Second
	}
	if strings.TrimSpace(cfg.DefaultSource) == "" {
		cfg.DefaultSource = "capture"
	}
	if cfg.MaxHealthyQueueDepth <= 0 {
		cfg.MaxHealthyQueueDepth = 8
	}
	return &CaptureWorker{
		spool:   spool,
		source:  source,
		cfg:     cfg,
		metrics: defaultVoiceMetrics(),
		status: CaptureStatus{
			Status: "idle",
		},
	}
}

func (w *CaptureWorker) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopChan != nil {
		return nil
	}
	if w.spool == nil {
		return fmt.Errorf("voice capture spool is not configured")
	}
	if w.source == nil {
		return fmt.Errorf("voice capture source is not configured")
	}
	w.stopChan = make(chan struct{})
	go w.runLoop(w.stopChan)
	return nil
}

func (w *CaptureWorker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopChan == nil {
		return
	}
	close(w.stopChan)
	w.stopChan = nil
}

func (w *CaptureWorker) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.stopChan != nil
}

func (w *CaptureWorker) StatusSnapshot() CaptureStatus {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.status
}

func (w *CaptureWorker) HealthCheck() health.CheckResult {
	status := w.StatusSnapshot()
	if status.LastBeatAt.IsZero() {
		return health.CheckResult{Status: "error", Message: "voice capture has no heartbeat"}
	}
	if time.Since(status.LastBeatAt) > w.cfg.HeartbeatFreshness {
		return health.CheckResult{Status: "error", Message: "voice capture heartbeat stale"}
	}
	if status.QueueDepth > w.cfg.MaxHealthyQueueDepth {
		return health.CheckResult{Status: "warning", Message: fmt.Sprintf("voice capture backlog=%d", status.QueueDepth)}
	}
	if status.LastError != "" {
		return health.CheckResult{Status: "warning", Message: status.LastError}
	}
	return health.CheckResult{Status: "ok", Message: "voice capture healthy"}
}

func (w *CaptureWorker) StatusHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.Header().Set("Allow", http.MethodGet)
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(w.StatusSnapshot())
	})
}

func (w *CaptureWorker) CaptureOnce(ctx context.Context) error {
	if w.spool == nil {
		return fmt.Errorf("voice capture spool is not configured")
	}
	if w.source == nil {
		return fmt.Errorf("voice capture source is not configured")
	}

	event, err := w.source.Capture(ctx)
	if err != nil {
		w.recordCapture("error")
		w.updateStatus(func(status *CaptureStatus) {
			status.Status = "error"
			status.LastError = err.Error()
			status.RejectedClips++
		})
		w.writeHeartbeat()
		return err
	}
	if event == nil {
		w.recordCapture("idle")
		w.updateStatus(func(status *CaptureStatus) {
			status.Status = "idle"
			status.LastError = ""
		})
		w.writeHeartbeat()
		return nil
	}

	userID := event.UserID
	chatID := event.ChatID
	if userID == 0 {
		userID = w.cfg.DefaultUserID
	}
	if chatID == 0 {
		chatID = w.cfg.DefaultChatID
	}
	source := strings.TrimSpace(event.Source)
	if source == "" {
		source = w.cfg.DefaultSource
	}

	_, err = w.spool.EnqueueAudioFile(Job{
		Source:        source,
		UserID:        userID,
		ChatID:        chatID,
		RequiresAudio: event.RequiresAudio,
	}, event.AudioFile)
	if err != nil {
		w.recordCapture("error")
		w.updateStatus(func(status *CaptureStatus) {
			status.Status = "error"
			status.LastError = err.Error()
			status.RejectedClips++
		})
		w.writeHeartbeat()
		return err
	}
	if event.DeleteSourceAfter {
		_ = os.Remove(event.AudioFile)
	}

	w.recordCapture("captured")
	w.updateStatus(func(status *CaptureStatus) {
		status.Status = "captured"
		status.LastCapturedAt = time.Now().UTC()
		status.CapturedJobs++
		status.LastError = ""
	})
	w.writeHeartbeat()
	return nil
}

func (w *CaptureWorker) runLoop(stopChan chan struct{}) {
	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		default:
		}
		_ = w.CaptureOnce(context.Background())
		select {
		case <-stopChan:
			return
		case <-ticker.C:
		}
	}
}

func (w *CaptureWorker) writeHeartbeat() {
	if strings.TrimSpace(w.cfg.HeartbeatPath) == "" {
		return
	}
	depth := 0
	if w.spool != nil {
		if queueDepth, err := w.spool.QueueDepth(); err == nil {
			depth = queueDepth
		}
	}

	w.updateStatus(func(status *CaptureStatus) {
		status.LastBeatAt = time.Now().UTC()
		status.QueueDepth = depth
	})
	w.recordQueueDepth(depth)
	status := w.StatusSnapshot()

	if err := os.MkdirAll(filepath.Dir(w.cfg.HeartbeatPath), 0o700); err != nil {
		return
	}
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return
	}
	data = append(data, '\n')
	_ = os.WriteFile(w.cfg.HeartbeatPath, data, 0o600)
}

func (w *CaptureWorker) updateStatus(fn func(*CaptureStatus)) {
	w.mu.Lock()
	defer w.mu.Unlock()
	fn(&w.status)
}

func (w *CaptureWorker) recordCapture(result string) {
	if w == nil || w.metrics == nil {
		return
	}
	w.metrics.capture.WithLabelValues(result).Inc()
}

func (w *CaptureWorker) recordQueueDepth(depth int) {
	if w == nil || w.metrics == nil {
		return
	}
	w.metrics.queueDepth.Set(float64(depth))
}
