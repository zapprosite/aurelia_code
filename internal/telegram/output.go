package telegram

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/pkg/tts"
	"gopkg.in/telebot.v3"
)

func SendMediaParts(bot *telebot.Bot, chat *telebot.Chat, parts []agent.ContentPart) error {
	if len(parts) == 0 {
		return nil
	}

	var textParts []string
	var imageParts []agent.ContentPart

	for _, p := range parts {
		if p.Type == agent.ContentPartText && p.Text != "" {
			textParts = append(textParts, p.Text)
		} else if p.Type == agent.ContentPartImage {
			imageParts = append(imageParts, p)
		}
	}

	fullText := strings.Join(textParts, "\n\n")

	if len(imageParts) == 0 {
		if fullText != "" {
			return SendText(bot, chat, fullText)
		}
		return nil
	}

	// Handle first image with the full text as caption if it fits
	firstImage := imageParts[0]
	remainingImages := imageParts[1:]

	tmpFile, err := os.CreateTemp(os.TempDir(), "aurelia-media-*.jpg")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(firstImage.Data); err != nil {
		_ = tmpFile.Close()
		return err
	}
	_ = tmpFile.Close()

	if len(fullText) > 1024 {
		// Telegram caption limit is 1024. If it's longer, send text separately.
		if _, err := bot.Send(chat, &telebot.Photo{File: telebot.FromDisk(tmpPath)}); err != nil {
			return err
		}
		if err := SendText(bot, chat, fullText); err != nil {
			return err
		}
	} else {
		photo := &telebot.Photo{
			File:    telebot.FromDisk(tmpPath),
			Caption: MarkdownToHTML(fullText),
		}
		if _, err := bot.Send(chat, photo, &telebot.SendOptions{ParseMode: telebot.ModeHTML}); err != nil {
			// Fallback to plain text caption if HTML fails
			photo.Caption = fullText
			if _, err := bot.Send(chat, photo); err != nil {
				return err
			}
		}
	}

	// Send remaining images
	for _, img := range remainingImages {
		_ = func() error {
			t, e := os.CreateTemp(os.TempDir(), "aurelia-media-*.jpg")
			if e != nil {
				return e
			}
			p := t.Name()
			defer os.Remove(p)
			if _, e := t.Write(img.Data); e != nil {
				_ = t.Close()
				return e
			}
			_ = t.Close()
			_, e = bot.Send(chat, &telebot.Photo{File: telebot.FromDisk(p)})
			return e
		}()
	}

	return nil
}

type ttsAsyncResult struct {
	audio tts.Audio
	err   error
}

const telegramMessageLimit = 3900

type messageSender interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
}

func SendText(bot *telebot.Bot, chat *telebot.Chat, text string) error {
	if os.Getenv("RUN_SWARM_E2E") != "" {
		log.Printf("\n[LLM SWARM RESPONSE E2E]:\n%s\n\n", text)
		return nil
	}
	// S-30: Polish Premium — sanitização final antes do envio
	text = sanitizeAssistantOutputForUser(text)
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
		return nil // No audio to send, text already sent
	}

	limit := synthesizer.MaxChars()
	speechText := sanitizeTextForSpeech(text, limit)
	if speechText == "" {
		return nil
	}

	ttsCtx, ttsCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ttsCancel()
	audio, err := synthesizer.Synthesize(ttsCtx, speechText)
	if err != nil {
		log.Printf("TTS synthesis failed (%v).", err)
		return nil // Don't fail the whole response if only audio fails
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

// deliverWithParallelTTS sends text to the user while concurrently synthesizing
// TTS audio, then sends the audio once both are ready. This hides Kokoro latency
// (~0.5-1.5s) behind the Telegram SendText round-trip for lower perceived latency.
func deliverWithParallelTTS(sender messageSender, chat *telebot.Chat, synthesizer tts.Synthesizer, text string) error {
	var ttsCh chan ttsAsyncResult
	if synthesizer != nil && synthesizer.IsAvailable() {
		limit := synthesizer.MaxChars()
		if speechText := sanitizeTextForSpeech(text, limit); speechText != "" {
			ttsCh = make(chan ttsAsyncResult, 1)
			go func() {
				ttsCtx, ttsCancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer ttsCancel()
				audio, err := synthesizer.Synthesize(ttsCtx, speechText)
				ttsCh <- ttsAsyncResult{audio, err}
			}()
		}
	}

	// Send text while TTS is generating in background.
	if err := sendTextWithSender(sender, chat, text, telegramMessageLimit); err != nil {
		return err
	}

	if ttsCh == nil {
		return nil
	}

	r := <-ttsCh
	if r.err != nil {
		log.Printf("TTS synthesis failed (%v).", r.err)
		return nil
	}
	return sendAudioBytes(sender, chat, r.audio, text)
}

// sendAudioBytes sends a pre-synthesized tts.Audio as a Telegram voice note or audio clip.
func sendAudioBytes(sender messageSender, chat *telebot.Chat, audio tts.Audio, fallbackText string) error {
	tmpFile, err := os.CreateTemp(os.TempDir(), "aurelia-tts-*"+audio.Extension)
	if err != nil {
		log.Printf("TTS temp file create failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, fallbackText, telegramMessageLimit)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.Write(audio.Data); err != nil {
		log.Printf("TTS temp file write failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, fallbackText, telegramMessageLimit)
	}
	if err := tmpFile.Close(); err != nil {
		log.Printf("TTS temp file close failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, fallbackText, telegramMessageLimit)
	}

	if audio.AsVoiceNote {
		voice := &telebot.Voice{
			File: telebot.FromDisk(tmpPath),
			MIME: "audio/ogg",
		}
		if _, err := sender.Send(chat, voice); err == nil {
			log.Printf("Telegram voice note sent successfully (%d bytes).", len(audio.Data))
			return nil
		}
		log.Printf("Telegram voice send failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, fallbackText, telegramMessageLimit)
	}

	clip := &telebot.Audio{
		File:     telebot.FromDisk(tmpPath),
		MIME:     audio.ContentType,
		FileName: filepath.Base(tmpPath),
		Title:    "Aurelia",
	}
	if _, err := sender.Send(chat, clip); err != nil {
		log.Printf("Telegram audio clip send failed (%v). Falling back to text output...", err)
		return sendTextWithSender(sender, chat, fallbackText, telegramMessageLimit)
	}
	log.Printf("Telegram audio clip sent successfully (%d bytes).", len(audio.Data))
	return nil
}

var (
	markdownLinkPattern       = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	codeFencePattern          = regexp.MustCompile("(?s)```[a-z]*\n?(.*?)```")
	inlineCodePattern         = regexp.MustCompile("`([^`]+)`")
	multiSpacePattern         = regexp.MustCompile(`\s+`)

	// speechSymbolReplacer converts symbols/punctuation to spoken equivalents.
	speechSymbolReplacer = strings.NewReplacer(
		"R$", "reais ",
		"%", " por cento",
		"&", " e ",
		"/", " ou ",
		"->", " para ",
		"=>", " resulta em ",
		"!=", " diferente de ",
		"==", " igual a ",
		">=", " maior ou igual a ",
		"<=", " menor ou igual a ",
		"...", ". ",
		"\n\n", ". ",
		"\n", ". ",
		"|", ", ",
	)

	// markdownDecorationReplacer strips visual markdown that has no spoken value.
	markdownDecorationReplacer = strings.NewReplacer(
		"**", "",
		"__", "",
		"***", "",
		"*", "",
		"_", "",
		"##", "",
		"#", "",
		">", "",
		"`", "",
	)
)

func sanitizeTextForSpeech(text string, limit int) string {
	sanitized := strings.TrimSpace(text)
	if sanitized == "" {
		return ""
	}

	// Remove code blocks entirely — reading code aloud is useless.
	sanitized = codeFencePattern.ReplaceAllString(sanitized, "")

	// Keep link label, discard URL.
	sanitized = markdownLinkPattern.ReplaceAllString(sanitized, "$1")

	// Inline code: strip backticks but keep content.
	sanitized = inlineCodePattern.ReplaceAllString(sanitized, "$1")

	// Expand common symbols to spoken equivalents.
	sanitized = speechSymbolReplacer.Replace(sanitized)

	// Strip remaining markdown decoration.
	sanitized = markdownDecorationReplacer.Replace(sanitized)

	// Collapse whitespace and fix double-period artifacts.
	sanitized = multiSpacePattern.ReplaceAllString(sanitized, " ")
	sanitized = strings.ReplaceAll(sanitized, ".. ", ". ")
	sanitized = strings.ReplaceAll(sanitized, "..", ".")
	sanitized = strings.TrimSpace(sanitized)

	// Truncate at requested limit for TTS synthesis.
	// Cut at the last sentence boundary to avoid stopping mid-phrase.
	runes := []rune(sanitized)
	if len(runes) > limit {
		window := string(runes[:limit])
		// Optimization: only search for sentence ends in the last 1000 characters 
		// of the window to avoid discarding huge chunks of text.
		searchStart := 0
		if limit > 1000 {
			searchStart = limit - 1000
		}
		
		found := false
		for _, sep := range []string{". ", "! ", "? ", ".\n", "!\n", "?\n"} {
			if idx := strings.LastIndex(window, sep); idx >= searchStart {
				sanitized = strings.TrimSpace(window[:idx+1])
				found = true
				break
			}
		}
		if !found {
			// fallback: hard cut at limit
			sanitized = strings.TrimSpace(string(runes[:limit])) + "."
		}
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
