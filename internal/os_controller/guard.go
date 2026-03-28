package os_controller

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Padrões destrutivos que exigem intervenção ou bloqueio
	destructivePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\brm\s+-[rf]+\s+/`),
		regexp.MustCompile(`(?i)\bmkfs\b`),
		regexp.MustCompile(`(?i)\bdd\s+if=`),
		regexp.MustCompile(`(?i)>[ ]*/dev/sd[a-z]`),
		regexp.MustCompile(`(?i)\biptables\s+-F\b`),
		regexp.MustCompile(`(?i)\bshutdown\b`),
		regexp.MustCompile(`(?i)\breboot\b`),
	}
)

type ExecutionGuard struct {
	UnsafeAuto bool
}

func NewExecutionGuard(unsafeAuto bool) *ExecutionGuard {
	return &ExecutionGuard{UnsafeAuto: unsafeAuto}
}

// Validate verifica se o comando é seguro.
func (g *ExecutionGuard) Validate(script string) error {
	for _, pattern := range destructivePatterns {
		if pattern.MatchString(script) {
			if !g.UnsafeAuto {
				return fmt.Errorf("BLOQUEIO DE SEGURANÇA: Comando detectado como destrutivo pelo Execution Guard: %s", script)
			}
			// Se estiver em modo unsafe, podemos logar um aviso mas permitir
		}
	}

	// Bloqueio extra para tentativas de escapar do path raiz se necessário
	if strings.Contains(script, "/etc/shadow") || strings.Contains(script, "/etc/sudoers") {
		if !g.UnsafeAuto {
			return fmt.Errorf("BLOQUEIO DE SEGURANÇA: Acesso a arquivos de sistema sensíveis negado")
		}
	}

	return nil
}
