package telegram

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kocar/aurelia/pkg/tts"
	"gopkg.in/telebot.v3"
)

type sendCall struct {
	to   telebot.Recipient
	what interface{}
	opts []interface{}
}

type stubSender struct {
	calls        []sendCall
	firstSendErr error
}

func (s *stubSender) Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	s.calls = append(s.calls, sendCall{to: to, what: what, opts: opts})
	if len(s.calls) == 1 && s.firstSendErr != nil {
		return nil, s.firstSendErr
	}
	return &telebot.Message{}, nil
}

type stubSynthesizer struct {
	audio tts.Audio
	err   error
}

func (s stubSynthesizer) Synthesize(context.Context, string) (tts.Audio, error) {
	return s.audio, s.err
}

func (s stubSynthesizer) IsAvailable() bool {
	return s.err == nil || len(s.audio.Data) > 0
}

func TestSendText_SendsTelegramHTML(t *testing.T) {
	sender := &stubSender{}
	chat := &telebot.Chat{ID: 123}

	if err := sendTextWithSender(sender, chat, "## Title\n\n- **item**", 200); err != nil {
		t.Fatalf("sendTextWithSender returned error: %v", err)
	}

	if len(sender.calls) != 1 {
		t.Fatalf("expected 1 send call, got %d", len(sender.calls))
	}

	text, ok := sender.calls[0].what.(string)
	if !ok {
		t.Fatalf("expected sent payload to be string, got %T", sender.calls[0].what)
	}
	if !containsSubstring(text, "<b>Title</b>") {
		t.Fatalf("expected html formatted text, got: %s", text)
	}

	options, ok := sender.calls[0].opts[0].(*telebot.SendOptions)
	if !ok {
		t.Fatalf("expected first option to be *telebot.SendOptions, got %T", sender.calls[0].opts[0])
	}
	if options.ParseMode != telebot.ModeHTML {
		t.Fatalf("expected parse mode %q, got %q", telebot.ModeHTML, options.ParseMode)
	}
}

func TestSendText_FallsBackToPlainTextWhenHTMLSendFails(t *testing.T) {
	sender := &stubSender{firstSendErr: errors.New("bad html")}
	chat := &telebot.Chat{ID: 123}

	if err := sendTextWithSender(sender, chat, "## Title", 200); err != nil {
		t.Fatalf("sendTextWithSender returned error: %v", err)
	}

	if len(sender.calls) != 2 {
		t.Fatalf("expected 2 send calls, got %d", len(sender.calls))
	}

	if _, ok := sender.calls[1].what.(string); !ok {
		t.Fatalf("expected plain text fallback payload, got %T", sender.calls[1].what)
	}
	if len(sender.calls[1].opts) != 0 {
		t.Fatalf("expected fallback send without options, got %d opts", len(sender.calls[1].opts))
	}
}

func TestSplitTelegramMarkdown_PrefersParagraphBoundaries(t *testing.T) {
	text := "primeiro bloco\n\nsegundo bloco muito maior para obrigar split"

	chunks := splitTelegramMarkdown(text, 35)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	if chunks[0] != "primeiro bloco" {
		t.Fatalf("expected first chunk to stop at paragraph boundary, got %q", chunks[0])
	}
	for _, chunk := range chunks {
		if len([]rune(chunk)) > 35 {
			t.Fatalf("chunk exceeded limit: %q", chunk)
		}
	}
}

type mockSender struct {
	onSend func(c *telebot.Chat, message string) error
}

func (m *mockSender) Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	if m.onSend != nil {
		if msg, ok := what.(string); ok {
			return nil, m.onSend(to.(*telebot.Chat), msg)
		}
	}
	return nil, nil
}

func TestSendError_SendsFormattedHTML(t *testing.T) {
	sender := &stubSender{}
	chat := &telebot.Chat{ID: 123}

	if err := sendErrorWithSender(sender, chat, "Erro", sanitizeUserVisibleErrorMessage("provider error: empty guarded content for openrouter:google/gemini-2.5-flash")); err != nil {
		t.Fatalf("sendErrorWithSender returned error: %v", err)
	}

	if len(sender.calls) != 1 {
		t.Fatalf("expected 1 send call, got %d", len(sender.calls))
	}

	text, ok := sender.calls[0].what.(string)
	if !ok {
		t.Fatalf("expected sent payload to be string, got %T", sender.calls[0].what)
	}
	if !containsSubstring(text, "<b>Erro</b>") {
		t.Fatalf("expected bold html title, got: %s", text)
	}
	if !containsSubstring(text, "falha temporaria do runtime") {
		t.Fatalf("expected error body in payload, got: %s", text)
	}
	if containsSubstring(strings.ToLower(text), "gemini-2.5") {
		t.Fatalf("expected provider/model details to be hidden, got: %s", text)
	}

	options, ok := sender.calls[0].opts[0].(*telebot.SendOptions)
	if !ok {
		t.Fatalf("expected first option to be *telebot.SendOptions, got %T", sender.calls[0].opts[0])
	}
	if options.ParseMode != telebot.ModeHTML {
		t.Fatalf("expected parse mode %q, got %q", telebot.ModeHTML, options.ParseMode)
	}
}

func TestSendError_FallsBackToPlainTextWhenHTMLSendFails(t *testing.T) {
	sender := &stubSender{firstSendErr: errors.New("bad html")}
	chat := &telebot.Chat{ID: 123}

	if err := sendErrorWithSender(sender, chat, "Erro", sanitizeUserVisibleErrorMessage("provider error: route breaker open for openrouter:deepseek/deepseek-v3.2")); err != nil {
		t.Fatalf("sendErrorWithSender returned error: %v", err)
	}

	if len(sender.calls) != 2 {
		t.Fatalf("expected 2 send calls, got %d", len(sender.calls))
	}

	payload, ok := sender.calls[1].what.(string)
	if !ok {
		t.Fatalf("expected plain text fallback payload, got %T", sender.calls[1].what)
	}
	if payload != "Erro\n\nNao consegui concluir isso agora por uma falha temporaria do runtime. Tente novamente em alguns segundos." {
		t.Fatalf("unexpected fallback payload: %q", payload)
	}
}

func TestSendAudio_SendsTelegramVoiceWhenTTSSucceeds(t *testing.T) {
	sender := &stubSender{}
	chat := &telebot.Chat{ID: 123}

	err := sendAudioWithSender(sender, chat, stubSynthesizer{
		audio: tts.Audio{
			Data:        []byte("voice"),
			ContentType: "audio/opus",
			Extension:   ".ogg",
			AsVoiceNote: true,
		},
	}, "## Ola")
	if err != nil {
		t.Fatalf("sendAudioWithSender returned error: %v", err)
	}

	if len(sender.calls) != 1 {
		t.Fatalf("expected 1 send call, got %d", len(sender.calls))
	}
	if _, ok := sender.calls[0].what.(*telebot.Voice); !ok {
		t.Fatalf("expected telebot.Voice payload, got %T", sender.calls[0].what)
	}
}

func TestSendAudio_FallsBackToTextWhenTTSFails(t *testing.T) {
	sender := &stubSender{}
	chat := &telebot.Chat{ID: 123}

	err := sendAudioWithSender(sender, chat, stubSynthesizer{err: errors.New("tts down")}, "## Ola")
	if err != nil {
		t.Fatalf("sendAudioWithSender returned error: %v", err)
	}

	if len(sender.calls) != 1 {
		t.Fatalf("expected 1 send call, got %d", len(sender.calls))
	}
	if _, ok := sender.calls[0].what.(string); !ok {
		t.Fatalf("expected text fallback payload, got %T", sender.calls[0].what)
	}
}

func TestSanitizeTextForSpeech_StripsMarkdown(t *testing.T) {
	got := sanitizeTextForSpeech("## Titulo\n\n- **item** com [link](https://example.com)\n\n`codigo`")
	if containsSubstring(got, "#") || containsSubstring(got, "**") || containsSubstring(got, "`") {
		t.Fatalf("unexpected markdown in sanitized text: %q", got)
	}
	if !containsSubstring(got, "Titulo") || !containsSubstring(got, "item com link") {
		t.Fatalf("unexpected sanitized text: %q", got)
	}
}

func TestSanitizeAssistantOutputForUser_BlocksInternalLeak(t *testing.T) {
	got := sanitizeAssistantOutputForUser("provider error: empty guarded content for openrouter:deepseek/deepseek-v3.2")
	if got != genericResponseGuardMessage {
		t.Fatalf("unexpected sanitized output: %q", got)
	}
}

func TestSanitizeAssistantOutputForUser_LeavesNormalAnswer(t *testing.T) {
	got := sanitizeAssistantOutputForUser("Sistema estavel. Docker, GPU e disco estao saudaveis.")
	if got != "Sistema estavel. Docker, GPU e disco estao saudaveis." {
		t.Fatalf("unexpected sanitized output: %q", got)
	}
}
