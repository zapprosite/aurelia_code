package telegram

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kocar/aurelia/pkg/tts"
	"gopkg.in/telebot.v3"
)

const telegramMessageLimit = 3900

type messageSender interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
}

func SendText(bot *telebot.Bot, chat *telebot.Chat, text string) error {
	return sendTextWithSender(bot, chat, text, telegramMessageLimit)
}

func sendTextWithSender(sender messageSender, chat *telebot.Chat, text string, limit int) error {
	chunks := splitTelegramMarkdown(text, limit)
	for _, chunk := range chunks {
		htmlChunk := MarkdownToHTML(chunk)
		_, err := sender.Send(chat, htmlChunk, &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
		})
		if err == nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		log.Printf("Send chunk with HTML failed (%v). Retrying as plain text...", err)
		_, err = sender.Send(chat, chunk)
		if err != nil {
			if floodErr, ok := err.(*telebot.FloodError); ok {
				log.Printf("Hit rate limit in chunk sending. Retrying in %v...", floodErr.RetryAfter)
				time.Sleep(time.Duration(floodErr.RetryAfter) * time.Second)
				if _, retryErr := sender.Send(chat, chunk); retryErr == nil {
					time.Sleep(200 * time.Millisecond)
					continue
				}
			}
			return err
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func splitTelegramMarkdown(text string, limit int) []string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return []string{""}
	}

	var chunks []string
	remaining := trimmed
	for len([]rune(remaining)) > limit {
		splitAt := bestSplitIndex(remaining, limit)
		chunks = append(chunks, strings.TrimSpace(remaining[:splitAt]))
		remaining = strings.TrimSpace(remaining[splitAt:])
	}
	if remaining != "" {
		chunks = append(chunks, remaining)
	}
	return chunks
}

func bestSplitIndex(text string, limit int) int {
	runes := []rune(text)
	if len(runes) <= limit {
		return len(text)
	}

	candidates := []string{"\n\n", "\n", ". ", " "}
	window := string(runes[:limit])
	for _, candidate := range candidates {
		if idx := strings.LastIndex(window, candidate); idx > 0 {
			return idx
		}
	}
	return len(string(runes[:limit]))
}

func SendDocument(bot *telebot.Bot, chat *telebot.Chat, filename, content string) error {
	tmpDir := os.TempDir()
	path := filepath.Join(tmpDir, filename)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Println("SendDocument tmp write failed, sending as fallback text...")
		return SendText(bot, chat, "Nao consegui gerar arq, segue texto puro:\n\n"+content)
	}
	defer func() { _ = os.Remove(path) }()

	doc := &telebot.Document{
		File:     telebot.FromDisk(path),
		FileName: filename,
		MIME:     "text/markdown",
	}

	_, err = bot.Send(chat, doc)
	return err
}

func SendError(bot *telebot.Bot, chat *telebot.Chat, errMsg string) error {
	return sendErrorWithSender(bot, chat, "Erro", sanitizeUserVisibleErrorMessage(errMsg))
}

func SendAudio(bot *telebot.Bot, chat *telebot.Chat, text string) error {
	return sendAudioWithSender(bot, chat, nil, text)
}

func sendAudioWithSender(sender messageSender, chat *telebot.Chat, synthesizer tts.Synthesizer, text string) error {
	if synthesizer == nil || !synthesizer.IsAvailable() {
		log.Println("TTS synthesizer unavailable. Falling back to text output.")
		return sendTextWithSender(sender, chat, text, telegramMessageLimit)
	}

	speechText := sanitizeTextForSpeech(text)
	if speechText == "" {
		return sendTextWithSender(sender, chat, text, telegramMessageLimit)
	}

	audio, err := synthesizer.Synthesize(context.Background(), speechText)
	if err != nil {
		log.Printf("TTS synthesis failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, text, telegramMessageLimit)
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "aurelia-tts-*"+audio.Extension)
	if err != nil {
		log.Printf("TTS temp file create failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, text, telegramMessageLimit)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.Write(audio.Data); err != nil {
		log.Printf("TTS temp file write failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, text, telegramMessageLimit)
	}
	if err := tmpFile.Close(); err != nil {
		log.Printf("TTS temp file close failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, text, telegramMessageLimit)
	}

	if audio.AsVoiceNote {
		voice := &telebot.Voice{
			File: telebot.FromDisk(tmpPath),
			MIME: "audio/ogg",
		}
		if _, err := sender.Send(chat, voice); err == nil {
			log.Printf("Telegram voice note sent successfully (%d bytes).", len(audio.Data))
			return nil
		} else {
			log.Printf("Telegram voice send failed (%v). Falling back to text output...", err)
			return sendTextWithSender(sender, chat, text, telegramMessageLimit)
		}
	}

	clip := &telebot.Audio{
		File:     telebot.FromDisk(tmpPath),
		MIME:     audio.ContentType,
		FileName: filepath.Base(tmpPath),
		Title:    "Aurelia",
	}
	if _, err := sender.Send(chat, clip); err != nil {
		log.Printf("Telegram audio send failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, text, telegramMessageLimit)
	}
	log.Printf("Telegram audio clip sent successfully (%d bytes).", len(audio.Data))
	return nil
}

var (
	markdownLinkPattern = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	codeFencePattern    = regexp.MustCompile("(?s)```(.*?)```")
	multiSpacePattern   = regexp.MustCompile(`\s+`)
)

func sanitizeTextForSpeech(text string) string {
	sanitized := strings.TrimSpace(text)
	if sanitized == "" {
		return ""
	}
	sanitized = codeFencePattern.ReplaceAllString(sanitized, "$1")
	sanitized = markdownLinkPattern.ReplaceAllString(sanitized, "$1")
	replacer := strings.NewReplacer(
		"`", "",
		"**", "",
		"__", "",
		"*", "",
		"_", "",
		"#", "",
		">", "",
		"|", ", ",
		"\n", ". ",
	)
	sanitized = replacer.Replace(sanitized)
	sanitized = multiSpacePattern.ReplaceAllString(sanitized, " ")
	sanitized = strings.TrimSpace(sanitized)
	runes := []rune(sanitized)
	if len(runes) > 1200 {
		sanitized = strings.TrimSpace(string(runes[:1200])) + "."
	}
	return sanitized
}

func sendErrorWithSender(sender messageSender, chat *telebot.Chat, title, errMsg string) error {
	formatted := ErrorMessage(title, errMsg)
	_, err := sender.Send(chat, formatted, &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	})
	if err == nil {
		return nil
	}

	log.Printf("Send error with HTML failed (%v). Retrying as plain text...", err)
	_, err = sender.Send(chat, title+"\n\n"+errMsg)
	return err
}
