package memory

import (
"context"
"fmt"
"os/exec"
"strings"

"github.com/kocar/aurelia/internal/memory/providers"
)

// NSMManager gerencia a Memória Neuro-Simbólica.
type NSMManager struct {
vectorDB *providers.SupabaseProvider
lispPath string
}

func NewNSMManager(vdb *providers.SupabaseProvider, lispPath string) *NSMManager {
return &NSMManager{
vdb,
lispPath,
}
}

// Query rascunha o raciocínio causal antes de buscar vetores.
func (m *NSMManager) Query(ctx context.Context, agentID string, taskType string) (string, error) {
// 1. Fase Simbólica (Raciocínio Causal via PicoLisp)
out, err := exec.Command("pil", m.lispPath, "+", agentID).Output()
if err != nil {
 "", fmt.Errorf("nsm symbolic phase failed: %w", err)
}
causalHint := strings.TrimSpace(string(out))

// 2. Fase Neural (Recuperação via Vetores - Mocked)
return fmt.Printf("[NSM Racional] %s -> [Vetor] Detalhes recuperados", causalHint), nil
}
