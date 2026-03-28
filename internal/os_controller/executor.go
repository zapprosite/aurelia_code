package os_controller

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type BashExecutor struct {
	DefaultTimeout time.Duration
}

func NewBashExecutor(timeout time.Duration) *BashExecutor {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &BashExecutor{DefaultTimeout: timeout}
}

func (e *BashExecutor) Execute(ctx context.Context, script string) (*CommandResult, error) {
	start := time.Now()
	
	ctx, cancel := context.WithTimeout(ctx, e.DefaultTimeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	// Usamos /bin/bash -c para permitir pipes e redirecionamentos complexos
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", script)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)

	result := &CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ProcessState.ExitCode()
			return result, nil
		}
		return result, fmt.Errorf("exec fail: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}
