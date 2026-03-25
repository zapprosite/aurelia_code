package media

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/pkg/stt"
)

var (
	youtubeRegex = regexp.MustCompile(`(?i)https?://(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/|youtube\.com/shorts/)([a-zA-Z0-9_-]{11})`)
)

type Processor struct {
	stt      stt.Transcriber
	llm      agent.LLMProvider
	tempDir  string
}

type Result struct {
	Transcript string
	Summary    string
	SourceURL  string
	Duration   time.Duration
}

func NewProcessor(stt stt.Transcriber, llm agent.LLMProvider, tempDir string) *Processor {
	if tempDir == "" {
		tempDir = os.TempDir()
	}
	return &Processor{
		stt:     stt,
		llm:     llm,
		tempDir: tempDir,
	}
}

func (p *Processor) IsSupportedURL(url string) bool {
	return youtubeRegex.MatchString(url)
}

func (p *Processor) ProcessURL(ctx context.Context, url string) (*Result, error) {
	logger := observability.Logger("media.processor")
	logger.Info("processing media URL", slog.String("url", url))

	// Check if yt-dlp is available for URLs
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return nil, fmt.Errorf("YouTube download requires 'yt-dlp' to be installed on the host system")
	}

	// 1. Download audio via yt-dlp
	audioPath, err := p.downloadAudio(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("download audio: %w", err)
	}
	defer os.Remove(audioPath)

	// 2. Transcribe
	logger.Info("transcribing media", slog.String("file", filepath.Base(audioPath)))
	transcript, err := p.stt.Transcribe(ctx, audioPath)
	if err != nil {
		return nil, fmt.Errorf("transcribe: %w", err)
	}

	// 3. Summarize
	logger.Info("summarizing media transcript")
	summary, err := p.Summarize(ctx, transcript)
	if err != nil {
		logger.Warn("summarization failed, returning only transcript", slog.Any("err", err))
		summary = "[O resumo falhou, mas a transcrição está disponível abaixo]"
	}

	return &Result{
		Transcript: transcript,
		Summary:    summary,
		SourceURL:  url,
	}, nil
}

func (p *Processor) downloadAudio(ctx context.Context, url string) (string, error) {
	outputPattern := filepath.Join(p.tempDir, fmt.Sprintf("aurelia-media-%d.%%(ext)s", time.Now().UnixNano()))
	
	// yt-dlp -f bestaudio --extract-audio --audio-format wav --audio-quality 0 -o <path> <url>
	// We use wav to be safe for STT providers
	cmd := exec.CommandContext(ctx, "yt-dlp", 
		"-f", "bestaudio", 
		"--extract-audio", 
		"--audio-format", "wav", 
		"-o", outputPattern, 
		"--no-playlist",
		url,
	)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Find the actual file (yt-dlp adds extension)
	files, _ := filepath.Glob(filepath.Join(p.tempDir, "aurelia-media-*.wav"))
	if len(files) == 0 {
		return "", fmt.Errorf("yt-dlp output file not found")
	}
	
	// Return the most recent one matching our prefix
	return files[len(files)-1], nil
}

func (p *Processor) Summarize(ctx context.Context, transcript string) (string, error) {
	if p.llm == nil {
		return "", fmt.Errorf("no LLM provider configured for summarization")
	}

	systemPrompt := `Você é a Aurélia, uma assistente industrial soberana e técnica.
Sua tarefa é resumir a transcrição de uma mídia rica (vídeo ou áudio) fornecida pelo usuário.
O resumo deve ser executivo, técnico e direto ao ponto, destacando os insights principais.
Use o formato de lista e mantenha um tom formal e profissional.
NÃO use gírias ou introduções desnecessárias. Vá direto ao conteúdo.
USE APENAS MARKDOWN LIMPO. NÃO use blocos de código JSON.`

	prompt := fmt.Sprintf("Por favor, resuma a seguinte transcrição de mídia:\n\n%s", transcript)
	
	// Forçar soberania semântica (Local Only) para resumos de mídia
	ctx = agent.WithRunOptions(ctx, agent.RunOptions{LocalOnly: true})

	resp, err := p.llm.GenerateContent(ctx, systemPrompt, []agent.Message{
		{Role: "user", Content: prompt},
	}, nil)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}
