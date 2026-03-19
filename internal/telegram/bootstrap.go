package telegram

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kocar/aurelia/internal/observability"
	"gopkg.in/telebot.v3"
)

func (bc *BotController) setupBootstrapRoutes() {
	bc.bot.Handle("/start", bc.handleStart)
	bc.bot.Handle("\fbtn_coder", bc.handleBootstrapChoice("coder"))
	bc.bot.Handle("\fbtn_assist", bc.handleBootstrapChoice("assist"))
}

func (bc *BotController) handleStart(c telebot.Context) error {
	identityExists := bootstrapIdentityExists(bc.personasDir)

	message, menu := bootstrapStartResponse(identityExists)
	if menu == nil {
		return SendContextText(c, message)
	}
	return SendContextText(c, message, menu)
}

func (bc *BotController) handleBootstrapChoice(choice string) func(telebot.Context) error {
	logger := observability.Logger("telegram.bootstrap")
	return func(c telebot.Context) error {
		_ = bc.bot.Respond(c.Callback(), &telebot.CallbackResponse{})

		preset, err := bootstrapPresetForChoice(choice)
		if err != nil {
			return SendContextText(c, bootstrapFailureMessage)
		}
		if err := writeBootstrapPreset(bc.personasDir, preset); err != nil {
			logger.Error("failed to write bootstrap preset", slog.Any("err", err))
			return SendContextText(c, bootstrapFailureMessage)
		}

		bc.setPendingBootstrap(c.Sender().ID, bootstrapState{Choice: choice})
		if err := bc.seedBootstrapIdentity(c, preset); err != nil {
			logger.Warn("failed to seed bootstrap identity", slog.Any("err", err))
		}
		return SendContextText(c, bootstrapProfileMessage)
	}
}

func bootstrapIdentityExists(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "IDENTITY.md"))
	return err == nil
}
