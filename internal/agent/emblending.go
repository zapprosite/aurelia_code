package agent

import (
"fmt"
"os/exec"
)

type EmblendingEngine struct{}

func NewEmblendingEngine() *EmblendingEngine {
return &EmblendingEngine{}
}

func (e *EmblendingEngine) Fuse(intent string, facts map[string]string) string {
fmt.Printf("[EMBLENDING] Calculando pesos via PicoLisp para intenção: %s\n", intent)

// Chamada ao PicoLisp para lógica simbólica de pesagem
cmd := exec.Command("pil", "internal/agent/emblending_logic.l", "-main", intent)
out, _ := cmd.CombinedOutput()
fmt.Printf("[EMBLENDING] Pesos calculados: %s\n", string(out))

// Verificação de Contradição
if facts["Vision"] != facts["MCP"] && facts["MCP"] != "" {
fmt.Printf("[EMBLENDING-ALERT] CONTRADIÇÃO DETECTADA!\n")
return "ERROR_CONTRADICTION"
}

return "SUCCESS_FUSED_CONTEXT"
}
