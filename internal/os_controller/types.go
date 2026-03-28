package os_controller

import (
	"context"
	"time"
)

// CommandResult representa o resultado da execução de um comando no OS.
type CommandResult struct {
	Stdout   string        `json:"stdout"`
	Stderr   string        `json:"stderr"`
	ExitCode int           `json:"exit_code"`
	Duration time.Duration `json:"duration"`
}

// Executor define a interface para execução de comandos.
type Executor interface {
	Execute(ctx context.Context, command string, args []string) (*CommandResult, error)
}

// OSController orquestra a execução segura de comandos.
type OSController interface {
	RunBash(ctx context.Context, script string) (*CommandResult, error)
	ReadLog(ctx context.Context, path string, lines int) ([]string, error)
}
