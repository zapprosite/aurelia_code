package telegram

import (
	"gopkg.in/telebot.v3"
)

func (bc *BotController) setupPersonaRoutes() {
	bc.bot.Handle("/junior", bc.handleJunior)
	bc.bot.Handle("/senior", bc.handleSenior)
}

func (bc *BotController) handleJunior(c telebot.Context) error {
	bc.personaID = "junior-developer"
	bc.botName = "Aurélia_Code (Junior) 🐣"
	
	msg := "🐣 **Modo Junior Ativado.**\n\n" +
		"Estou pronto para ajudar, Master! Serei mais didático e pedirei ajuda ao Sênior para coisas complexas. Como posso aprender com você hoje?"
	
	return c.Send(msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
}

func (bc *BotController) handleSenior(c telebot.Context) error {
	bc.personaID = "aurelia-sovereign"
	bc.botName = "Aurélia_Code (Soberana) 💎"
	
	msg := "💎 **Modo Sênior Restaurado.**\n\n" +
		"Soberania reestabelecida. Arquiteta Sênior em comando para operações industriais de alto impacto."
	
	return c.Send(msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
}
