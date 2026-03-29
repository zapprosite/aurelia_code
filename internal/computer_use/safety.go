// Package computer_use provides autonomous computer use agent.
package computer_use

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

// SafetyGuard blocks dangerous actions and provides HitL confirmation.
type SafetyGuard struct {
	dangerousPatterns []*regexp.Regexp
	sensitiveKeywords []string
	log              *slog.Logger
	hitl             *HitLGate
}

// NewSafetyGuard creates a new safety guard.
func NewSafetyGuard(hitl *HitLGate) *SafetyGuard {
	return &SafetyGuard{
		dangerousPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\brm\s+-rf\b`),
			regexp.MustCompile(`(?i)\bdd\s+if\b`),
			regexp.MustCompile(`(?i)\bmkfs\b`),
			regexp.MustCompile(`(?i)\bwipefs\b`),
			regexp.MustCompile(`(?i)\bsudo\s+rm\b`),
			regexp.MustCompile(`(?i)\bsudo\s+chmod\s+777\b`),
			regexp.MustCompile(`(?i)\bchmod\s+-R\s+777\b`),
			regexp.MustCompile(`(?i)\bssh\s+.*@`),
			regexp.MustCompile(`(?i)\bcurl\s+.*\|\s*(sh|bash)`),
			regexp.MustCompile(`(?i)\bwget\s+.*\|\s*(sh|bash)`),
		},
		sensitiveKeywords: []string{
			"password", "credit card", "ssn",
			"api key", "secret key", "private key",
		},
		log:  slog.Default(),
		hitl: hitl,
	}
}

// ValidateAction checks if an action is safe to execute.
func (s *SafetyGuard) ValidateAction(action string, params map[string]interface{}) (bool, string) {
	// Check for dangerous patterns in all string params
	for _, v := range params {
		if str, ok := v.(string); ok {
			for _, pattern := range s.dangerousPatterns {
				if pattern.MatchString(str) {
					s.log.Warn("dangerous pattern blocked",
						"pattern", pattern.String(),
						"action", action)
					return false, fmt.Sprintf("blocked dangerous pattern: %s", pattern.String())
				}
			}
		}
	}

	// Check for sensitive keywords
	actionLower := strings.ToLower(action)
	for _, keyword := range s.sensitiveKeywords {
		if strings.Contains(actionLower, keyword) {
			s.log.Warn("sensitive action detected",
				"keyword", keyword,
				"action", action)
			return false, fmt.Sprintf("requires confirmation: sensitive keyword '%s'", keyword)
		}
	}

	return true, ""
}

// HitLGate provides human-in-the-loop confirmation for sensitive actions.
type HitLGate struct {
	enabled   bool
	confirmCh chan Confirmation
	timeout   time.Duration
}

// Confirmation represents a confirmation request.
type Confirmation struct {
	Action  string
	Details string
	Approved bool
}

// NewHitLGate creates a new HitL gate.
func NewHitLGate(enabled bool, timeout time.Duration) *HitLGate {
	return &HitLGate{
		enabled:   enabled,
		confirmCh: make(chan Confirmation, 1),
		timeout:  timeout,
	}
}

// RequestConfirmation waits for user confirmation.
func (h *HitLGate) RequestConfirmation(ctx context.Context, action, details string) error {
	if !h.enabled {
		return nil // Auto-approve if disabled
	}

	select {
	case confirm := <-h.confirmCh:
		if !confirm.Approved {
			return fmt.Errorf("action rejected by user: %s", action)
		}
		return nil

	case <-time.After(h.timeout):
		return fmt.Errorf("hitl timeout: user did not respond in %v", h.timeout)

	case <-ctx.Done():
		return ctx.Err()
	}
}

// Approve approves the pending action.
func (h *HitLGate) Approve() {
	h.confirmCh <- Confirmation{Approved: true}
}

// Reject rejects the pending action.
func (h *HitLGate) Reject() {
	h.confirmCh <- Confirmation{Approved: false}
}

// IsEnabled returns if HitL is enabled.
func (h *HitLGate) IsEnabled() bool {
	return h.enabled
}

// Middleware combines safety and HitL for agent actions.
type Middleware struct {
	guard *SafetyGuard
	hitl   *HitLGate
}

// NewMiddleware creates a new safety middleware.
func NewMiddleware(hitl *HitLGate) *Middleware {
	return &Middleware{
		guard: NewSafetyGuard(hitl),
		hitl:  hitl,
	}
}

// Check validates an action with safety guard.
func (m *Middleware) Check(ctx context.Context, action string, params map[string]interface{}) error {
	// Check safety
	safe, reason := m.guard.ValidateAction(action, params)
	if !safe {
		return fmt.Errorf("safety blocked: %s", reason)
	}

	// Check if sensitive (needs HitL)
	if isSensitive(action, params) {
		if err := m.hitl.RequestConfirmation(ctx, action, formatAction(action, params)); err != nil {
			return fmt.Errorf("hitl rejected: %w", err)
		}
	}

	return nil
}

// isSensitive determines if an action needs HitL confirmation.
func isSensitive(action string, params map[string]interface{}) bool {
	sensitiveActions := map[string]bool{
		"type_text": true,
		"click":     true,
		"extract":   true,
		"navigate":  true,
	}

	if !sensitiveActions[action] {
		return false
	}

	// Check params for sensitive data
	for _, val := range params {
		str := strings.ToLower(fmt.Sprintf("%v", val))
		sensitive := []string{"password", "credit", "ssn", "key", "secret"}
		for _, s := range sensitive {
			if strings.Contains(str, s) {
				return true
			}
		}
	}

	return false
}

// formatAction formats an action for display.
func formatAction(action string, params map[string]interface{}) string {
	return fmt.Sprintf("%s(%v)", action, params)
}
