package tts

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// EdgeTTS implements Synthesizer using Microsoft Edge TTS API
// Free, natural PT-BR voices - no API key needed
type EdgeTTS struct {
	voice   string
	venvBin string
}

// Audio represents synthesized speech
type EdgeAudio struct {
	Data       []byte
	Format     string
	SampleRate int
	Channels   int
}

// NewEdgeSynthesizer creates an Edge TTS synthesizer
// voice: pt-BR-FranciscaNeural, pt-BR-ThalitaMultilingualNeural, etc.
func NewEdgeSynthesizer(venvPath, voice string) *EdgeTTS {
	return &EdgeTTS{
		voice:   voice,
		venvBin: filepath.Join(venvPath, "bin", "python3"),
	}
}

// Synthesize converts text to speech using Edge TTS
func (e *EdgeTTS) Synthesize(ctx context.Context, text string) (Audio, error) {
	// Find the script
	script := "/home/will/aurelia/scripts/edge-tts.py"

	tmpFile, err := os.CreateTemp(os.TempDir(), "aurelia-edge-*.mp3")
	if err != nil {
		return Audio{}, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Call Python script
	cmd := exec.CommandContext(ctx, e.venvBin, script,
		"--text", text,
		"--voice", e.voice,
		"--output", tmpPath,
	)
	cmd.Env = append(os.Environ(), "AURELIA_VENV="+filepath.Dir(filepath.Dir(e.venvBin)))

	if err := cmd.Run(); err != nil {
		return Audio{}, fmt.Errorf("edge-tts failed: %w", err)
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return Audio{}, fmt.Errorf("read output: %w", err)
	}

	return Audio{
		Data:        data,
		ContentType: "audio/mp3",
		AsVoiceNote: true,
		Extension:   ".mp3",
	}, nil
}

// MaxChars returns the maximum characters per request
func (e *EdgeTTS) MaxChars() int {
	return 2000 // Edge TTS limit
}

// IsAvailable checks if Edge TTS is available
func (e *EdgeTTS) IsAvailable() bool {
	cmd := exec.Command(e.venvBin, "-c", "import edge_tts")
	return cmd.Run() == nil
}

// Cleanup removes temporary files
func (e *EdgeTTS) Cleanup() error {
	return nil
}
