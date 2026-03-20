package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"github.com/kocar/aurelia/internal/voice"
)

func main() {
	fmt.Println("=== SMOKE TEST: KOKORO (KODORO) LOCAL GPU TTS PT-BR ===")

	client := voice.NewKokoroClient()
	text := "Olá! Sou a Aurélia, sua assistente onisciente. O motor local de voz PT-BR está operacional via GPU."
	voiceID := "pf_dora"

	fmt.Printf("[SIM] Gerando áudio via GPU para: \"%s\"\n", text)
	audioData, err := client.GenerateSpeech(text, voiceID)
	if err != nil {
		fmt.Printf("[ERROR] Falha ao gerar áudio local: %v\n", err)
		fmt.Println("Certifique-se de que o container 'kokoro-fastapi' está rodando em http://localhost:8888")
		os.Exit(1)
	}

	outputFile := "sample_kokoro_br.mp3"
	err = ioutil.WriteFile(outputFile, audioData, 0644)
	if err != nil {
		fmt.Printf("[ERROR] Falha ao salvar arquivo: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[SUCCESS] Áudio local gerado: %s (%d bytes)\n", outputFile, len(audioData))
}
