package audio

import (
	"context"
	"log/slog"
	"os/exec"
)

// SimplePlayer utiliza o mpv para reproduzir chunks de áudio em tempo real via stdin.
type SimplePlayer struct {
	command string
	logger  *slog.Logger
}

func (p *SimplePlayer) Name() string {
	return "speaker"
}

func NewSimplePlayer() *SimplePlayer {
	return &SimplePlayer{
		command: "mpv",
		logger:  slog.Default().With("actor", "speaker", "component", "sap"),
	}
}

// Speak consome chunks de áudio e os envia para o player nativo.
func (p *SimplePlayer) Speak(ctx context.Context, chunks <-chan []byte) error {
	p.logger.Info("Starting audio playback")
	// --no-terminal: Silencia o mpv
	// --vo=null: Desativa saída de vídeo
	// --cache=no: Minimiza latência
	// -: Lê do stdin
	cmd := exec.CommandContext(ctx, p.command, "--no-terminal", "--vo=null", "--cache=no", "-")
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		p.logger.Error("Failed to create stdin pipe", "err", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		p.logger.Error("Failed to start player command", "err", err)
		return err
	}

	// Goroutine para alimentar o stdin sem bloquear o loop principal
	go func() {
		defer stdin.Close()
		for {
			select {
			case <-ctx.Done():
				p.logger.Warn("Playback context cancelled, stopping input")
				return
			case chunk, ok := <-chunks:
				if !ok {
					return
				}
				if _, err := stdin.Write(chunk); err != nil {
					p.logger.Warn("Write to player stdin failed", "err", err)
					return
				}
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		p.logger.Error("Player command exited with error", "err", err)
	} else {
		p.logger.Info("Audio playback completed")
	}
	return err
}
