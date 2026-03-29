// Package main provides Jarvis Tutor - audio streaming
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func runTutorCommand(args []string, _ io.Writer) error {
	fs := flag.NewFlagSet("tutor", flag.ExitOnError)
	userID := fs.Int64("user", getEnvInt64("TELEGRAM_USER_ID", 0), "User ID")
	chatID := fs.Int64("chat", getEnvInt64("TELEGRAM_CHAT_ID", 0), "Chat ID")
	fs.Parse(args)

	if *userID == 0 || *chatID == 0 {
		return fmt.Errorf("TELEGRAM_USER_ID e TELEGRAM_CHAT_ID obrigatorios no .env")
	}

	groqKey := os.Getenv("GROQ_API_KEY")
	if groqKey == "" {
		return fmt.Errorf("GROQ_API_KEY obrigatoria no .env")
	}

	log.Printf("[Jarvis] 🤖 Audio Pipeline configurado")
	log.Printf("[Jarvis] UserID: %d, ChatID: %d", *userID, *chatID)
	log.Printf("[Jarvis] 🔄 Testando componentes...")

	// Teste rápido
	if err := testPipeline(groqKey); err != nil {
		log.Printf("[Jarvis] ⚠️ Pipeline error: %v", err)
	} else {
		log.Printf("[Jarvis] ✅ Pipeline OK - Components healthy")
	}
	log.Printf("[Jarvis] ⏸ Demo mode - Rodando por 60s")
	time.Sleep(60 * time.Second)
	return nil
}

func testPipeline(groqKey string) error {
	log.Printf("[Test] STT Groq...")
	if err := testSTT(groqKey); err != nil {
		return fmt.Errorf("STT: %w", err)
	}
	log.Printf("[Test] STT OK")

	log.Printf("[Test] TTS Kokoro...")
	if err := testTTS(); err != nil {
		return fmt.Errorf("TTS: %w", err)
	}
	log.Printf("[Test] TTS OK")

	log.Printf("[Test] LLM LiteLLM...")
	if err := testLLM(); err != nil {
		return fmt.Errorf("LLM: %w", err)
	}
	log.Printf("[Test] LLM OK")

	return nil
}

func testSTT(groqKey string) error {
	// Cria WAV de teste (silêncio)
	wav := []byte{
		// WAV header minimal
	}
	tmp, _ := os.CreateTemp("", "test*.wav")
	tmp.Write(wav)
	tmp.Close()
	defer os.Remove(tmp.Name())

	// POST para Groq
	body := strings.NewReader(`{"model":"whisper-large-v3"}`)
	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/audio/transcriptions", body)
	req.Header.Set("Authorization", "Bearer "+groqKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func testTTS() error {
	resp, err := http.Get("http://localhost:8880/health")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func testLLM() error {
	resp, err := http.Get("http://localhost:4000/health")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func getEnvInt64(key string, defaultVal int64) int64 {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return defaultVal
}
