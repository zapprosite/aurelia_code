package telegram

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kocar/aurelia/internal/media"
	"github.com/kocar/aurelia/internal/observability"
	"gopkg.in/telebot.v3"
)

func (bc *BotController) handleMediaURL(c contextSender, session inputSession) error {
	logger := observability.Logger("telegram.media")
	logger.Info("handling media URL", slog.String("url", session.text))

	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	_ = SendContextText(c, "✨ Identifiquei um link de mídia rica. Iniciando download e transcrição...")

	result, err := bc.mediaProcessor.ProcessURL(session.ctx, session.text)
	if err != nil {
		logger.Error("failed to process media URL", slog.Any("err", err))
		return SendContextText(c, "❌ Falha ao processar a mídia: "+err.Error())
	}

	return bc.deliverMediaResult(c, session, result)
}

func (bc *BotController) handleVideo(c telebot.Context) error {
	video := c.Message().Video
	if video == nil {
		return nil
	}

	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	_ = SendContextText(c, "🎥 Vídeo recebido. Processando áudio e gerando resumo executivo...")
	return bc.processMediaFile(c, &video.File, video.FileID+".mp4")
}

func (bc *BotController) handleTranscreverCommand(c telebot.Context) error {
	args := c.Args()
	if len(args) > 0 {
		// handle as URL
		url := args[0]
		session := newInputSession(c, url)
		return bc.handleMediaURL(c, session)
	}

	// Check if it's a reply to a video/audio
	if c.Message().ReplyTo != nil {
		reply := c.Message().ReplyTo
		if reply.Video != nil {
			// Fake a context for handleVideo
			// We can't easily "fake" telebot.Context, but we can call a common method
			return bc.processMediaFile(c, &reply.Video.File, reply.Video.FileID+".mp4")
		}
		if reply.Voice != nil {
			return bc.processMediaFile(c, &reply.Voice.File, reply.Voice.FileID+".ogg")
		}
		if reply.Audio != nil {
			return bc.processMediaFile(c, &reply.Audio.File, reply.Audio.FileID+".mp3")
		}
	}

	return SendContextText(c, "💡 Use `/transcrever [url]` ou responda a um vídeo/áudio com `/transcrever`.")
}

func (bc *BotController) processMediaFile(c telebot.Context, file *telebot.File, filename string) error {
	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.RecordingAudio, 4*time.Second)
	defer stopTyping()

	_ = SendContextText(c, "🛠️ Processando arquivo solicitado...")

	filePath, err := bc.downloadTelegramFile(file, filename)
	if err != nil {
		return SendContextText(c, downloadFailureMessage)
	}
	defer os.Remove(filePath)

	transcript, err := bc.transcribeAudioFile(filePath)
	if err != nil {
		return SendContextText(c, "❌ Falha na transcrição: "+err.Error())
	}

	session := newInputSession(c, "[Processamento Manual]")
	summary, err := bc.mediaProcessor.Summarize(session.ctx, transcript)
	if err != nil {
		summary = "[Falha ao gerar resumo]"
	}

	result := &media.Result{
		Transcript: transcript,
		Summary:    summary,
		SourceURL:  filename,
	}

	return bc.deliverMediaResult(c, session, result)
}

func (bc *BotController) deliverMediaResult(c contextSender, session inputSession, result *media.Result) error {
	// Persist in memory
	content := fmt.Sprintf("RESUMO DE MÍDIA (%s):\n\n%s\n\nTRANSCRIPÇÃO:\n%s", result.SourceURL, result.Summary, result.Transcript)
	bc.persistAssistantAnswer(session, content)

	// Present to user
	output := fmt.Sprintf("📑 *Resumo Executivo de Mídia*\n\n%s\n\n---\n💡 _Use /transcrever para ver a transcrição completa._", result.Summary)
	return SendContextText(c, output)
}
