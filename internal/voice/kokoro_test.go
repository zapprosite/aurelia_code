package voice

import (
	"fmt"
	"testing"
)

func TestKokoroSynthesis(t *testing.T) {
	client := NewKokoroClient()
	text := "Teste de síntese interna Aurélia 2026. Kokoro Local GPU."
	
	fmt.Printf("[TEST] Gerando áudio via GPU para: \"%s\"\n", text)
	audioData, err := client.GenerateSpeech(text, "pf_dora")
	if err != nil {
		t.Fatalf("Falha na síntese Kokoro: %v", err)
	}

	if len(audioData) < 1000 {
		t.Errorf("Áudio gerado muito pequeno: %d bytes", len(audioData))
	}
	
	fmt.Printf("[SUCCESS] Síntese concluída com %d bytes\n", len(audioData))
}
