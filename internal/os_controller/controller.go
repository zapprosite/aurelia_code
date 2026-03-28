package os_controller

import (
	"context"
	"fmt"
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
	// Aplicação de patch via sed ou escrita temporária segura
	// Por enquanto, uma implementação via sed simplificada
	// NOTA: Em produção, o ideal é usar atomic writes em Go.
	
	// TODO: Implementar lógica de patch segura
	return nil, fmt.Errorf("apply_patch not yet implemented")
}
