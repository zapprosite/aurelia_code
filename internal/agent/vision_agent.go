package agent

import (
"fmt"
)

// VisionAgent processa solicitações de contexto visual (Grounding).
type VisionAgent struct {
ModelPrompt string
}

func NewVisionAgent() *VisionAgent {
return &VisionAgent{
ModelPrompt: "Analise esta imagem. Identifique elementos de UI, estados e dados funcionais. Retorne JSON estruturado.",
}
}

// InspectRegion captura e analisa uma área da tela.
func (v *VisionAgent) InspectRegion(x, y, w, h int) string {
fmt.Printf("[VISION] Capturando região (%d,%d) Tamanho: %dx%d...\n", x, y, w, h)
// Simulação de output JSON do modelo VL
return `{"elements": [{"type": "button", "label": "Pagar Fatura", "coord": [450, 200], "status": "enabled"}]}`
}

func (v *VisionAgent) GroundingFacts() string {
return "[VISION-FACTS] Fatura visualizada na tela: Valor R$ 1.500,00, Vencimento 20/03/2026."
}
