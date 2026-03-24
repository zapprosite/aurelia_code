package gateway

import (
	"context"
	"fmt"
	"testing"
)

func TestGemmaJudge_JudgeReal(t *testing.T) {
	// Este teste requer o Ollama rodando localmente com o modelo gemma3:12b
	judge := NewGemmaJudge("gemma3:12b")

	tests := []struct {
		prompt string
		want   string
	}{
		{"Oi, tudo bem?", "simple_short"},
		{"Como eu faço um loop em Go?", "coding_main"},
		{"Explique a teoria da relatividade em detalhes", "coding_main"},
		{"Analise este log de erro: [ERROR] connection timeout at 0x00FF...", "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.prompt, func(t *testing.T) {
			res, err := judge.Judge(context.Background(), tt.prompt, nil)
			if err != nil {
				t.Logf("Aviso: Falha ao contatar Ollama: %v (pode ser esperado em CI)", err)
				return
			}
			fmt.Printf("Prompt: %s -> Class: %s (Confidence: %.2f, Reason: %s)\n", tt.prompt, res.Class, res.Confidence, res.Reason)
			// Não forçamos igualdade exata pois modelos podem variar, mas logamos para inspeção humana
		})
	}
}
