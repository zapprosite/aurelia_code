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
		return SendContextText(c, "Diagnostico indisponivel.")
	}

	// Persistir o comando do usuário para manter contexto
	session := bc.newInputSession(c, c.Text())
	_ = bc.persistIncomingContext(session, c.Sender().ID)

	data, err := bc.healthReporter.GatewayStatusJSON()
	if err != nil {
		return SendContextText(c, "Erro ao obter status.")
	}

	reply := "Status do sistema:\n```json\n" + string(data) + "\n```"

	// Polimento e Segurança
	if bc.porteiro != nil {
		reply = bc.PolishText(session.ctx, reply)
		reply = bc.porteiro.SecureOutput(reply)
	}

	// Persistir a resposta da assistente
	bc.persistAssistantAnswer(session, reply)

	return SendContextText(c, reply)
}

func (bc *BotController) setupBootstrapRoutes() {
	bc.bot.Handle("/start", bc.handleStart)
	bc.bot.Handle("/status", bc.handleStatus)
	bc.bot.Handle("\fbtn_status", bc.handleStatus)
	bc.bot.Handle("\fbtn_coder", bc.handleBootstrapChoice("coder"))
	bc.bot.Handle("\fbtn_assist", bc.handleBootstrapChoice("assist"))
}

func (bc *BotController) handleStart(c telebot.Context) error {
	identityExists := bootstrapIdentityExists(bc.personasDir)

	if !identityExists {
		preset, _ := bootstrapPresetForChoice("coder")
		_ = writeBootstrapPreset(bc.personasDir, preset)
		_ = bc.seedBootstrapIdentity(c, preset)
	}

	welcome := `🛰️ **AURELIA SOVEREIGN ENGINE [v2.2.0-SOTA]**
---
**IDENTIDADE**: Engenheira Sênior e Sombra Digital do Will.
**PROTOCOLO**: Industrial Homelab Governance (Sovereign 2026).
**STATUS**: Online e Sincronizada com NVIDIA RTX / Ollama.

Mestra em **Go, tRPC, Zod e Infraestrutura Autônoma**. 
Pronta para agir no monorepo e no homelab com autoridade total.

> [ONBOARDING] Defina o escopo da Feature para iniciarmos a execução.
---
Capacidades ativas: /status, /master-skill, Vision, Parallel-TTS.`

	menu := &telebot.ReplyMarkup{}
	btnDashboard := menu.URL("Dashboard", "https://aurelia.zappro.site/")
	btnStatus := menu.Data("Status do Sistema", "btn_status")
	menu.Inline(
		menu.Row(btnDashboard, btnStatus),
	)

	// SOTA 2026: Entrega premium com áudio paralelo (Kododo)
	return bc.deliverWithParallelTTS(bc.bot, c.Chat(), bc.tts, welcome, false, menu)
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
