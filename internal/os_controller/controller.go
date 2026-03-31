package os_controller

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

type Controller struct {
	executor *BashExecutor
	guard    *ExecutionGuard
}

func NewController(timeout time.Duration, unsafeAuto bool) *Controller {
	return &Controller{
		executor: NewBashExecutor(timeout),
		guard:    NewExecutionGuard(unsafeAuto),
	}
}

func (c *Controller) RunBash(ctx context.Context, script string) (*CommandResult, error) {
	// 1. Validar via Execution Guard
	if err := c.guard.Validate(script); err != nil {
		return nil, err
	}

	// 2. Executar via BashExecutor
	return c.executor.Execute(ctx, script)
}

func (c *Controller) ReadLog(ctx context.Context, path string, lines int) ([]string, error) {
	if lines <= 0 {
		lines = 50
	}
	
	// Prevenir leitura de arquivos fora de contexto ou via caminhos perigosos
	if err := c.guard.Validate("read_log " + path); err != nil {
		return nil, err
	}

	script := fmt.Sprintf("tail -n %d %s", lines, path)
	result, err := c.executor.Execute(ctx, script)
	if err != nil {
		return nil, err
	}

	if result.ExitCode != 0 {
		return nil, fmt.Errorf("read log error (exit %d): %s", result.ExitCode, result.Stderr)
	}

	return strings.Split(strings.TrimSpace(result.Stdout), "\n"), nil
}

func (c *Controller) ApplyPatch(ctx context.Context, path, oldContent, newContent string) (*CommandResult, error) {
	if err := c.guard.Validate("apply_patch " + path); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file for patch: %w", err)
	}

	original := string(data)
	if !strings.Contains(original, oldContent) {
		return nil, fmt.Errorf("old content not found in %s", path)
	}

	patched := strings.Replace(original, oldContent, newContent, 1)

	tmpPath := path + ".aurelia-patch.tmp"
	if err := os.WriteFile(tmpPath, []byte(patched), 0644); err != nil {
		return nil, fmt.Errorf("write temp patch file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("atomic rename patch: %w", err)
	}

	return &CommandResult{
		Stdout:   fmt.Sprintf("patched %s: replaced %d bytes", path, len(oldContent)),
		ExitCode: 0,
	}, nil
}
