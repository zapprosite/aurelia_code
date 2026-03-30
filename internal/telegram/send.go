package telegram

import (
	"context"

	"gopkg.in/telebot.v3"
)

type contextSender interface {
	Send(what interface{}, opts ...interface{}) error
	Chat() *telebot.Chat
}

func SendContextText(c contextSender, text string, opts ...interface{}) error {
	// Se o texto for JSON, tentamos um polimento básico antes de enviar.
	// Nota: O polimento ideal ocorre no BotController.
	html := MarkdownToHTML(text)
	sendOpts := append([]interface{}{&telebot.SendOptions{ParseMode: telebot.ModeHTML}}, opts...)
	if err := c.Send(html, sendOpts...); err == nil {
		return nil
	}
	return c.Send(text, opts...)
}

// SendPolishedContextText realiza o polimento via BotController antes do envio.
func (bc *BotController) SendPolishedContextText(c contextSender, text string, opts ...interface{}) error {
	polished := bc.PolishText(context.Background(), text)
	return SendContextText(c, polished, opts...)
}
