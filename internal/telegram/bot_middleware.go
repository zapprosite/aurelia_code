package telegram

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"gopkg.in/telebot.v3"
)

func (bc *BotController) whitelistMiddleware() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			senderID := c.Sender().ID
			if !bc.isAllowedUser(senderID) {
				slog.Warn("blocked unauthorized user", slog.Int64("user_id", senderID))
				return nil
			}
			return next(c)
		}
	}
}

// safeHandler wraps a telebot handler with panic recovery so a panic never
// leaves the user without a response.
func (bc *BotController) safeHandler(fn func(telebot.Context) error) func(telebot.Context) error {
	return func(c telebot.Context) (retErr error) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered in telegram handler",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)
				retErr = SendError(bc.bot, c.Chat(), fmt.Sprintf("Erro interno inesperado. Tente novamente."))
			}
		}()
		return fn(c)
	}
}

func (bc *BotController) registerContentRoutes() {
	bc.bot.Handle(telebot.OnText, bc.safeHandler(bc.handleText))
	bc.bot.Handle(telebot.OnPhoto, bc.safeHandler(bc.handlePhoto))
	bc.bot.Handle(telebot.OnDocument, bc.safeHandler(bc.handleDocument))
	bc.bot.Handle(telebot.OnVoice, bc.safeHandler(bc.handleVoice))
	bc.bot.Handle(telebot.OnAudio, bc.safeHandler(bc.handleVoice))
	bc.bot.Handle(telebot.OnVideo, bc.safeHandler(bc.handleVideo))
	bc.bot.Handle(telebot.OnAnimation, bc.safeHandler(bc.handleVideo))
	bc.bot.Handle("/transcrever", bc.safeHandler(bc.handleTranscreverCommand))
}


