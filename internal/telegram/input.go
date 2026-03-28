package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/observability"
	"gopkg.in/telebot.v3"
)

func (bc *BotController) handleText(c telebot.Context) error {
	return bc.processInput(c, c.Text(), false)
}

func (bc *BotController) handleDocument(c telebot.Context) error {
	doc := c.Message().Document
	if doc == nil {
		return nil
	}

	if !isSupportedDocument(doc.FileName, doc.MIME) {
		observability.Logger("telegram.input").Warn("unsupported document type", slog.String("mime", doc.MIME))
		return SendContextText(c, unsupportedDocumentMessage)
	}

	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	filePath, err := bc.downloadTelegramFile(&doc.File, doc.FileID+"_"+doc.FileName)
	if err != nil {
		observability.Logger("telegram.input").Warn("failed to download file", slog.Any("err", err))
		return SendContextText(c, downloadFailureMessage)
	}
	defer func() { _ = os.Remove(filePath) }()

	finalInput := buildDocumentInput(c.Message().Caption, doc.FileName, doc.MIME, filePath)
	return bc.processInput(c, finalInput, false)
}

func (bc *BotController) handleVoice(c telebot.Context) error {
	fileID, filename, ok := resolveAudioAttachment(c)
	if !ok {
		return nil
	}

	stopRecording := startChatActionLoop(bc.bot, c.Chat(), telebot.RecordingAudio, 4*time.Second)
	defer stopRecording()

	filePath, err := bc.downloadTelegramFile(&telebot.File{FileID: fileID}, fileID+"_"+filename)
	if err != nil {
		observability.Logger("telegram.input").Warn("failed to download audio", slog.Any("err", err))
		return SendContextText(c, downloadFailureMessage)
	}
	defer func() { _ = os.Remove(filePath) }()

	transcribedText, err := bc.transcribeAudioFile(filePath)
	if err != nil {
		observability.Logger("telegram.input").Warn("silent failure: audio transcription skipped", slog.Any("err", err))
		return nil
	}
	bc.persistAudioTranscript(c, filePath, transcribedText)
	return bc.processInput(c, transcribedText, true)
}

func isSupportedDocument(filename, mimeType string) bool {
	return strings.HasSuffix(filename, ".md") || mimeType == "application/pdf"
}

func (bc *BotController) downloadTelegramFile(file *telebot.File, filename string) (string, error) {
	filePath := filepath.Join(os.TempDir(), filename)
	if err := bc.bot.Download(file, filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

func buildDocumentInput(caption, filename, mimeType, filePath string) string {
	var extractedText string
	if strings.HasSuffix(filename, ".md") {
		content, err := os.ReadFile(filePath)
		if err == nil {
			extractedText = string(content)
		}
	} else if mimeType == "application/pdf" {
		extractedText = fmt.Sprintf("[Parsed content of PDF %s]", filename)
	}

	return fmt.Sprintf("%s\n\n[Analise o anexo %s]:\n%s", caption, filename, extractedText)
}

func resolveAudioAttachment(c telebot.Context) (string, string, bool) {
	switch {
	case c.Message().Voice != nil:
		return c.Message().Voice.FileID, "voice.ogg", true
	case c.Message().Audio != nil:
		return c.Message().Audio.FileID, "audio.mp3", true
	default:
		return "", "", false
	}
}

func (bc *BotController) transcribeAudioFile(filePath string) (string, error) {
	logger := observability.Logger("telegram.input")
	if bc.stt == nil || !bc.stt.IsAvailable() {
		return "", SendContextTextError(audioNotConfiguredMessage)
	}

	logger.Info("sending audio to transcriber", slog.String("file", observability.Basename(filePath)))
	transcribedText, err := bc.stt.Transcribe(context.Background(), filePath)
	if err != nil {
		logger.Warn("transcriber error", slog.Any("err", err))
		return "", SendContextTextError(audioProcessingFailureMessage)
	}
	if strings.TrimSpace(transcribedText) == "" {
		return "", SendContextTextError(emptyAudioMessage)
	}
	return transcribedText, nil
}

func (bc *BotController) persistAudioTranscript(c telebot.Context, filePath, transcript string) {
	if bc == nil || bc.memory == nil || c == nil || c.Sender() == nil {
		return
	}

	bc.persistAudioTranscriptForSender(c.Sender().ID, filePath, transcript)
}

func (bc *BotController) persistAudioTranscriptForSender(senderID int64, filePath, transcript string) {
	if bc == nil || bc.memory == nil {
		return
	}

	conversationID := bc.scopedConversationID(senderID)
	ctx := context.Background()

	if err := bc.memory.EnsureConversation(ctx, conversationID, senderID, "groq"); err != nil {
		observability.Logger("telegram.input").Warn("failed to ensure conversation for audio transcript", slog.Any("err", err))
		return
	}

	entry := memory.ArchiveEntry{
		ConversationID: conversationID,
		SessionID:      conversationID,
		Role:           "user",
		Content:        formatAudioTranscriptArchiveContent(bc.config.STTProvider, filePath, transcript),
		MessageType:    "audio_transcript",
	}
	if err := bc.memory.AddArchiveEntry(ctx, entry); err != nil {
		observability.Logger("telegram.input").Warn("failed to persist audio transcript", slog.Any("err", err))
	}
}

func formatAudioTranscriptArchiveContent(provider, filePath, transcript string) string {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = "unknown"
	}
	return fmt.Sprintf("[audio_transcript]\nprovider=%s\nfile=%s\n\n%s", provider, observability.Basename(filePath), strings.TrimSpace(transcript))
}

type sendContextTextError string

func SendContextTextError(message string) error {
	return sendContextTextError(message)
}

func (e sendContextTextError) Error() string {
	return string(e)
}

func errorAs(err error, target *sendContextTextError) bool {
	if err == nil {
		return false
	}
	value, ok := err.(sendContextTextError)
	if !ok {
		return false
	}
	*target = value
	return true
}
