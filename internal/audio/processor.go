package audio

import (
	"context"
	"io"
	"time"
)

// ProcessEvent representa um segmento de áudio processado (ex: uma frase).
type ProcessEvent struct {
	Data      []byte
	Timestamp time.Time
}

// AudioProcessor analisa o stream de áudio em busca de atividade de voz (VAD).
type AudioProcessor struct {
	input  io.Reader
	output chan ProcessEvent
}

func NewAudioProcessor(input io.Reader) *AudioProcessor {
	return &AudioProcessor{
		input:  input,
		output: make(chan ProcessEvent, 10),
	}
}

// Run inicia o loop de processamento. 
// Em 2026, o VAD deve ser eficiente para não onerar a CPU do Jarvis.
func (p *AudioProcessor) Run(ctx context.Context) error {
	buffer := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := p.input.Read(buffer)
			if err != nil {
				if err == io.EOF {
					time.Sleep(10 * time.Millisecond)
					continue
				}
				return err
			}

			if n > 0 {
				// VAD simplificado: emite evento se houver qualquer dado no buffer circular.
				// A lógica real de silêncio (silence threshold) será refinada na integração com Whisper.
				p.output <- ProcessEvent{
					Data:      buffer[:n],
					Timestamp: time.Now(),
				}
			}
		}
	}
}

func (p *AudioProcessor) Events() <-chan ProcessEvent {
	return p.output
}
