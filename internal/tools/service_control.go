package tools

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

type serviceAction string

const (
	serviceActionStatus  serviceAction = "status"
	serviceActionRestart serviceAction = "restart"
	serviceActionStop    serviceAction = "stop"
	serviceActionStart   serviceAction = "start"
	serviceActionList    serviceAction = "list"
	serviceActionLogs    serviceAction = "logs"
)

type serviceResult struct {
	Action  string `json:"action"`
	Service string `json:"service,omitempty"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

func ServiceControlHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	action := optionalStringArg(args, "action")
	if action == "" {
		action = "list"
	}

	service := optionalStringArg(args, "service")
	timeout := 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result := serviceResult{Action: action, Service: service}

	switch serviceAction(action) {
	case serviceActionStatus:
		if service == "" {
			result.Error = "service name required for status"
		} else {
			cmd := exec.CommandContext(ctx, "systemctl", "status", service)
			output, err := cmd.CombinedOutput()
			if err != nil {
				// status returns non-zero even on success, so we check output
				result.Output = strings.TrimSpace(string(output))
				result.Success = strings.Contains(result.Output, "active")
			} else {
				result.Success = true
				result.Output = strings.TrimSpace(string(output))
			}
		}

	case serviceActionRestart:
		if service == "" {
			result.Error = "service name required for restart"
		} else {
			cmd := exec.CommandContext(ctx, "systemctl", "restart", service)
			_, err := cmd.CombinedOutput()
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Success = true
				result.Output = "Service restarted: " + service
			}
		}

	case serviceActionStop:
		if service == "" {
			result.Error = "service name required for stop"
		} else {
			cmd := exec.CommandContext(ctx, "systemctl", "stop", service)
			_, err := cmd.CombinedOutput()
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Success = true
				result.Output = "Service stopped: " + service
			}
		}

	case serviceActionStart:
		if service == "" {
			result.Error = "service name required for start"
		} else {
			cmd := exec.CommandContext(ctx, "systemctl", "start", service)
			_, err := cmd.CombinedOutput()
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Success = true
				result.Output = "Service started: " + service
			}
		}

	case serviceActionList:
		cmd := exec.CommandContext(ctx, "systemctl", "list-units", "--type=service", "--state=active,failed", "--no-pager")
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
			result.Output = strings.TrimSpace(string(output))
		}

	case serviceActionLogs:
		if service == "" {
			result.Error = "service name required for logs"
		} else {
			cmd := exec.CommandContext(ctx, "journalctl", "-u", service, "-n", "50", "--no-pager")
			output, err := cmd.CombinedOutput()
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Success = true
				result.Output = strings.TrimSpace(string(output))
			}
		}

	default:
		result.Error = "unknown action: " + action
	}

	payload, _ := json.Marshal(result)
	return string(payload), nil
}
