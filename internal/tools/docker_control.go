package tools

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

type dockerAction string

const (
	dockerActionPS       dockerAction = "ps"
	dockerActionRestart  dockerAction = "restart"
	dockerActionLogs     dockerAction = "logs"
	dockerActionStats    dockerAction = "stats"
	dockerActionComposeUp dockerAction = "compose_up"
)

type dockerResult struct {
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

func DockerControlHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	action := optionalStringArg(args, "action")
	if action == "" {
		return marshalResult(dockerResult{
			Success: false,
			Error:   "action is required (ps, restart, logs, stats, compose_up)",
		}), nil
	}

	container := optionalStringArg(args, "container")
	workdir := optionalStringArg(args, "workdir")

	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result := dockerResult{Action: action}

	switch dockerAction(action) {
	case dockerActionPS:
		cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "--format", "table {{.ID}}\t{{.Image}}\t{{.Status}}\t{{.Names}}")
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
			result.Output = strings.TrimSpace(string(output))
		}

	case dockerActionRestart:
		if container == "" {
			result.Error = "container name/ID required for restart"
		} else {
			cmd := exec.CommandContext(ctx, "docker", "restart", container)
			_, err := cmd.CombinedOutput()
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Success = true
				result.Output = "Container restarted: " + container
			}
		}

	case dockerActionLogs:
		if container == "" {
			result.Error = "container name/ID required for logs"
		} else {
			cmd := exec.CommandContext(ctx, "docker", "logs", "--tail", "50", container)
			output, err := cmd.CombinedOutput()
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Success = true
				result.Output = strings.TrimSpace(string(output))
			}
		}

	case dockerActionStats:
		cmd := exec.CommandContext(ctx, "docker", "stats", "--no-stream", "--format", "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}")
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
			result.Output = strings.TrimSpace(string(output))
		}

	case dockerActionComposeUp:
		if workdir == "" {
			workdir = "."
		}
		cmd := exec.CommandContext(ctx, "docker-compose", "up", "-d")
		cmd.Dir = workdir
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
			result.Output = strings.TrimSpace(string(output))
		}

	default:
		result.Error = "unknown action: " + action
	}

	payload, _ := json.Marshal(result)
	return string(payload), nil
}

func marshalResult(v interface{}) string {
	payload, _ := json.Marshal(v)
	return string(payload)
}
