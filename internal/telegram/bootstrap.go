package telegram

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kocar/aurelia/internal/observability"
	"gopkg.in/telebot.v3"
)

func (bc *BotController) handleStatus(c telebot.Context) error {
	// S-27: If squad reporter is available, show rich status
	if bc.squadReporter != nil {
		return bc.handleSquadStatus(c)
	}

	// Fallback to gateway status
	if bc.healthReporter == nil {
		return SendContextText(c, "Diagnostico indisponivel neste runtime.")
	}

	// Persistir o comando do usuário para manter contexto (ex: "Resolva !")
	session := bc.newInputSession(c, c.Text())
	_ = bc.persistIncomingContext(session, c.Sender().ID)

	data, err := bc.healthReporter.GatewayStatusJSON()
	if err != nil {
		return SendContextText(c, "Erro ao obter status do gateway.")
	}

	reply := "Gateway Status:\n```json\n" + string(data) + "\n```"

	// Persistir a resposta da assistente
	bc.persistAssistantAnswer(session, reply)

	return SendContextText(c, reply)
}

func (bc *BotController) setupBootstrapRoutes() {
	bc.bot.Handle("/start", bc.handleStart)
	bc.bot.Handle("/status", bc.handleStatus)
	bc.bot.Handle("\fbtn_coder", bc.handleBootstrapChoice("coder"))
	bc.bot.Handle("\fbtn_assist", bc.handleBootstrapChoice("assist"))
}

func (bc *BotController) handleStart(c telebot.Context) error {
	identityExists := bootstrapIdentityExists(bc.personasDir)

	if !identityExists {
		// Inicializa silenciosamente com o preset de coder (padrão sênior)
		preset, _ := bootstrapPresetForChoice("coder")
		_ = writeBootstrapPreset(bc.personasDir, preset)
		_ = bc.seedBootstrapIdentity(c, preset)
	}

	welcome := "━━━━━━ Aurelia Sovereign 2026 ━━━━━━\n\n" +
		"🛰️ **Soberania Ativa.**\n" +
		"Bem-vindo ao cockpit de comando, Master. O sistema está operando em regime industrial com **Gemma 3** e **RTX 4090**.\n\n" +
		"Toda a verbosidade técnica foi movida para o dashboard para manter este canal limpo e estratégico."

	menu := &telebot.ReplyMarkup{}
	btnDashboard := menu.URL("🛰️ Dashboard Operacional", "https://aurelia.zappro.site/")
	btnStatus := menu.Data("📊 Status do Sistema", "btn_status")
	menu.Inline(
		menu.Row(btnDashboard),
		menu.Row(btnStatus),
	)

	return c.Send(welcome, menu, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
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
