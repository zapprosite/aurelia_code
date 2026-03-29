// Package voice provides the Jarvis Tutor mode - escuta tudo, 24/7
// ADR: 20260328-jarvis-tutor-24-7

package voice

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// TutorProcessor processa áudio SEM wake word - escuta tudo
// ADR: 20260328-jarvis-tutor-24-7
type TutorProcessor struct {
	*Processor
	tutorMode bool
	systemPrompt string
}

// Tutor system prompt - paciente, explicativo
const TutorSystemPrompt = `Você é o Jarvis, um tutor inteligente e prestativo.

Características:
- Explica conceitos de forma clara e simples
- Pergunta para confirmar entendimento
- Dá exemplos práticos do dia a dia
- Corrige erros com gentileza, sem julgamento
- Usa analogias quando necessário
- Responde em português brasileiro
- Mantém o contexto da conversa

Se não souber algo, seja honesto e diga que vai pesquisar.
`

// NewTutorProcessor cria um processador em modo tutor
// Escuta tudo, sem wake word, 24/7
// ADR: 20260328-jarvis-tutor-24-7
func NewTutorProcessor(spool *Spool, primary, fallback STTProvider, dispatcher Dispatcher, cfg Config) *TutorProcessor {
	p := NewProcessor(spool, primary, fallback, dispatcher, cfg)

	return &TutorProcessor{
		Processor:  p,
		tutorMode:  true,
		systemPrompt: TutorSystemPrompt,
	}
}

// processTutorAudio processa áudio em modo tutor (sem wake word)
// ADR: 20260328-jarvis-tutor-24-7
func (tp *TutorProcessor) processTutorAudio(ctx context.Context, audioPath string, userID, chatID int64) (string, error) {
	// Transcreve direto, sem wake word
	transcript, err := tp.transcribeAudio(ctx, audioPath)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	// Limpa texto
	transcript = strings.TrimSpace(transcript)
	if transcript == "" {
		return "", nil
	}

	// Dispara resposta sempre (escuta tudo)
	if tp.dispatcher != nil {
		_ = tp.dispatcher.DispatchVoice(ctx, userID, chatID, transcript, true)
	}

	return transcript, nil
}

// transcribeAudio transcreve áudio usando primary ou fallback
func (tp *TutorProcessor) transcribeAudio(ctx context.Context, audioPath string) (string, error) {
	// Tenta primary (Groq)
	if tp.primary != nil && tp.primary.IsAvailable() {
		result, err := tp.primary.Transcribe(ctx, audioPath)
		if err == nil {
			return result, nil
		}
	}

	// Fallback para local (Whisper)
	if tp.fallback != nil && tp.fallback.IsAvailable() {
		return tp.fallback.Transcribe(ctx, audioPath)
	}

	return "", fmt.Errorf("no STT available")
}

// TutorConfig configuration for tutor mode
type TutorConfig struct {
	// AlwaysOn ativa escuta contínua sem wake word
	AlwaysOn bool
	// ResponseDelay tempo entre receber áudio e responder
	ResponseDelay time.Duration
	// MaxSilenceSeconds máximo de silêncio antes de timeout
	MaxSilenceSeconds int
}

// DefaultTutorConfig configuração padrão para BK 24/7
var DefaultTutorConfig = TutorConfig{
	AlwaysOn:        true,
	ResponseDelay:    500 * time.Millisecond,
	MaxSilenceSeconds: 30,
}
