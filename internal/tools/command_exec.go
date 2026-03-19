package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

const (
	defaultCommandTimeout = 120 * time.Second
	maxCommandTimeout     = 600 * time.Second
	maxCommandOutputChars = 16000 // Increased from 4000 to 16000 per picoclaw pattern
)

// Compiled deny patterns - more efficient than string matching
var deniedPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\brm\s+-[rf]{1,2}\b`),
	regexp.MustCompile(`\bdel\s+/[fq]\b`),
	regexp.MustCompile(`\brmdir\s+/s\b`),
	regexp.MustCompile(`\b(format|mkfs)\b\s`),
	regexp.MustCompile(`\b(shutdown|poweroff|halt)\b`),
	regexp.MustCompile(`\bdd\s+if=`),
	regexp.MustCompile(`\bgit\s+(reset|clean)\s+.*--(hard|fdx)`),
}

type commandResult struct {
	Command    string `json:"command"`
	Workdir    string `json:"workdir"`
	ExitCode   int    `json:"exit_code"`
	DurationMS int64  `json:"duration_ms"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	TimedOut   bool   `json:"timed_out,omitempty"`
}

func RunCommandHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	command := optionalStringArg(args, "command")
	if command == "" {
		return marshalCommandResult(commandResult{
			ExitCode: -2,
			Stderr:   "blocked: command is required",
		}), nil
	}

	workdir := optionalStringArg(args, "workdir")
	workspaceRoot := optionalStringArg(args, "workspace_root")
	if workdir == "" {
		workdir, _ = agent.WorkdirFromContext(ctx)
	}
	if workdir == "" {
		if _, _, ok := agent.TaskContextFromContext(ctx); ok {
			return marshalCommandResult(commandResult{
				Command:  command,
				ExitCode: -2,
				Stderr:   "blocked: task requires explicit workdir or task default workdir",
			}), nil
		}
		cwd, err := os.Getwd()
		if err == nil {
			workdir = cwd
		} else {
			workdir = "."
		}
	}
	workdir = filepath.Clean(workdir)

	timeout := defaultCommandTimeout
	if rawTimeout, ok := args["timeout_seconds"].(float64); ok && rawTimeout > 0 {
		timeout = time.Duration(rawTimeout * float64(time.Second))
		if timeout > maxCommandTimeout {
			timeout = maxCommandTimeout
		}
	}

	result := commandResult{
		Command: command,
		Workdir: workdir,
	}

	if isBlockedCommand(command) {
		result.ExitCode = -2
		result.Stderr = "blocked: command matches a forbidden pattern"
		return marshalCommandResult(result), nil
	}
	if workspaceRoot != "" && !isWithinWorkspace(workdir, workspaceRoot) {
		result.ExitCode = -2
		result.Stderr = "blocked: workdir is outside workspace root"
		return marshalCommandResult(result), nil
	}

	// Resolve symlinks immediately before execution (TOCTOU protection)
	resolvedDir, err := filepath.EvalSymlinks(workdir)
	if err == nil {
		workdir = resolvedDir
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(runCtx, "sh", "-c", command)
	cmd.Dir = workdir

	// Set process group for proper cleanup
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Panic recovery: prevent tool crashes from killing agent loop
	defer func() {
		if r := recover(); r != nil {
			result.ExitCode = 1
			result.Stderr = "command execution panic: " + toString(r)
		}
	}()

	err = cmd.Run()
	result.DurationMS = time.Since(start).Milliseconds()
	result.Stdout = truncateOutput(stdoutBuf.String())
	result.Stderr = truncateOutput(stderrBuf.String())

	if runCtx.Err() == context.DeadlineExceeded {
		result.ExitCode = -1
		result.TimedOut = true
		// Kill entire process group on timeout
		if cmd.Process != nil && cmd.Process.Pid > 0 {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		if strings.TrimSpace(result.Stderr) == "" {
			result.Stderr = "process killed by timeout"
		}
		return marshalCommandResult(result), nil
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
			if strings.TrimSpace(result.Stderr) == "" {
				result.Stderr = err.Error()
			}
		}
		return marshalCommandResult(result), nil
	}

	result.ExitCode = 0
	return marshalCommandResult(result), nil
}

var dangerousMode = os.Getenv("AURELIA_DANGEROUS_MODE") == "1"

func isBlockedCommand(command string) bool {
	normalized := strings.ToLower(strings.TrimSpace(command))

	// Check deny patterns (compiled regex for efficiency)
	for _, pattern := range deniedPatterns {
		if pattern.MatchString(normalized) {
			return true
		}
	}

	// Allow sudo only in dangerousMode
	if !dangerousMode && strings.Contains(normalized, "sudo") {
		return true
	}

	return false
}

func isWithinWorkspace(workdir, workspaceRoot string) bool {
	workdir = filepath.Clean(workdir)
	workspaceRoot = filepath.Clean(workspaceRoot)

	rel, err := filepath.Rel(workspaceRoot, workdir)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, "..") && rel != "..")
}

func truncateOutput(text string) string {
	text = strings.TrimSpace(text)
	if len(text) <= maxCommandOutputChars {
		return text
	}
	return text[:maxCommandOutputChars] + "\n[truncated]"
}

func marshalCommandResult(result commandResult) string {
	payload, err := json.Marshal(result)
	if err != nil {
		fallback := commandResult{
			Command:  result.Command,
			Workdir:  result.Workdir,
			ExitCode: 1,
			Stderr:   "failed to marshal command result",
		}
		payload, _ = json.Marshal(fallback)
	}
	return string(payload)
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case error:
		return val.Error()
	default:
		return "unknown panic"
	}
}
