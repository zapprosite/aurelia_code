package agent

import (
	"fmt"
	"net/url"
	"strings"
)

type BrowserLoginStage string

const (
	BrowserLoginStageStart    BrowserLoginStage = "start"
	BrowserLoginStageUsername BrowserLoginStage = "username"
	BrowserLoginStagePassword BrowserLoginStage = "password"
	BrowserLoginStageTwoFA    BrowserLoginStage = "two_factor"
	BrowserLoginStageSuccess  BrowserLoginStage = "success"
	BrowserLoginStageAbort    BrowserLoginStage = "abort"
)

type BrowserLoginPolicy struct {
	AllowedHosts []string
	MaxSteps     int
}

type BrowserLoginRequest struct {
	TargetURL string
	Stage     BrowserLoginStage
	StepCount int
}

func DefaultBrowserLoginPolicy() BrowserLoginPolicy {
	return BrowserLoginPolicy{
		AllowedHosts: nil,
		MaxSteps:     8,
	}
}

func (p BrowserLoginPolicy) Validate(req BrowserLoginRequest) error {
	if strings.TrimSpace(req.TargetURL) == "" {
		return fmt.Errorf("target URL is required")
	}
	host, err := normalizeLoginHost(req.TargetURL)
	if err != nil {
		return err
	}
	if len(p.AllowedHosts) > 0 && !hostAllowed(host, p.AllowedHosts) {
		return fmt.Errorf("host %q is not allowed for guided login", host)
	}
	if req.StepCount < 0 {
		return fmt.Errorf("step count cannot be negative")
	}
	maxSteps := p.MaxSteps
	if maxSteps <= 0 {
		maxSteps = DefaultBrowserLoginPolicy().MaxSteps
	}
	if req.StepCount > maxSteps {
		return fmt.Errorf("step budget exceeded: %d > %d", req.StepCount, maxSteps)
	}
	switch req.Stage {
	case BrowserLoginStageStart,
		BrowserLoginStageUsername,
		BrowserLoginStagePassword,
		BrowserLoginStageTwoFA,
		BrowserLoginStageSuccess,
		BrowserLoginStageAbort:
		return nil
	default:
		return fmt.Errorf("unsupported login stage %q", req.Stage)
	}
}

func (p BrowserLoginPolicy) NeedsHumanGate(req BrowserLoginRequest) bool {
	return req.Stage == BrowserLoginStagePassword || req.Stage == BrowserLoginStageTwoFA
}

func normalizeLoginHost(rawURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", fmt.Errorf("parse target URL: %w", err)
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "" {
		return "", fmt.Errorf("target URL must include a valid host")
	}
	return host, nil
}

func hostAllowed(host string, allowedHosts []string) bool {
	for _, item := range allowedHosts {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if normalized == "" {
			continue
		}
		if host == normalized || strings.HasSuffix(host, "."+normalized) {
			return true
		}
	}
	return false
}
