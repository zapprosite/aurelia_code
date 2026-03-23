package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/dashboard"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/skill"
	"gopkg.in/telebot.v3"
)

const defaultSystemPrompt = "Você é Aurélia, uma assistente virtual de inteligência artificial de alto nível. Sua comunicação deve ser profissional, prestativa, objetiva e estruturada. Utilize formatação Markdown (negritos, listas e tabelas) para garantir clareza máxima. Ao analisar códigos ou realizar tarefas técnicas, priorize a precisão e sempre valide suas ações com 'run_command' antes de confirmar a conclusão."

type inputSession struct {
	senderID string
	convID   string
	ctx      context.Context
	text     string
	message  agent.Message
}


func (bc *BotController) processInput(c telebot.Context, text string, requiresAudio bool) error {
	session := newInputSession(c, text)
	return bc.processInputSession(c, session, requiresAudio)
}

func (bc *BotController) processInputSession(c telebot.Context, session inputSession, requiresAudio bool) error {
	logger := observability.Logger("telegram.pipeline")
	session.text = strings.ReplaceAll(session.text, "\x00", "")

	if state, ok := bc.popPendingBootstrap(c.Sender().ID); ok {
		return bc.completeBootstrapProfile(c, state, session.text)
	}

	if handled, err := bc.handleMemoryCommand(c, session); handled {
		return err
	}

	if err := bc.persistIncomingContext(session, c.Sender().ID); err != nil {
		logger.Warn("failed to persist incoming context", slog.Any("err", err))
	}

	dashboard.Publish(dashboard.Event{
		Type:      "user_message",
		Agent:     "User",
		Action:    "Mensagem recebida",
		Payload:   session.text,
		Timestamp: time.Now().Format(time.Kitchen),
	})

	activeSkill, history, systemPrompt, allowedTools, err := bc.prepareExecution(session)
	if err != nil {
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return nil
	}

	finalAnswer, err := bc.executeConversation(c, session, activeSkill, history, systemPrompt, allowedTools)
	if err != nil {
		logger.Error("conversation execution failed", slog.Any("err", err))
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return nil
	}

	finalAnswer = sanitizeAssistantOutputForUser(finalAnswer)
	bc.persistAssistantAnswer(session, finalAnswer)
	return bc.deliverFinalAnswer(c, finalAnswer, requiresAudio)
}

func newInputSession(c telebot.Context, text string) inputSession {
	return newInputSessionWithContext(context.Background(), c.Sender().ID, text)
}

func newInputSessionWithContext(ctx context.Context, senderUserID int64, text string) inputSession {
	senderID := fmt.Sprintf("%d", senderUserID)
	convID := senderID
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = agent.WithRunContext(agent.WithTeamContext(ctx, convID, senderID), uuid.NewString())
	return inputSession{
		senderID: senderID,
		convID:   convID,
		ctx:      ctx,
		text:     text,
		message: agent.Message{
			Role:    "user",
			Content: text,
		},
	}
}


func (bc *BotController) ProcessExternalInput(ctx context.Context, userID, chatID int64, text string, requiresAudio bool) error {
	if bc == nil || bc.bot == nil {
		return fmt.Errorf("telegram bot controller is not available")
	}
	chat := &telebot.Chat{ID: chatID}
	session := newInputSessionWithContext(ctx, userID, text)
	text = strings.ReplaceAll(text, "\x00", "")
	session.text = text

	if handled, err := bc.handleExternalMemoryCommand(chat, session); handled {
		return err
	}
	if err := bc.persistIncomingContext(session, userID); err != nil {
		observability.Logger("telegram.pipeline").Warn("failed to persist external input context", slog.Any("err", err))
	}

	activeSkill, history, systemPrompt, allowedTools, err := bc.prepareExecution(session)
	if err != nil {
		_ = SendError(bc.bot, chat, err.Error())
		return err
	}

	finalAnswer, err := bc.executeExternalConversation(chat, session, activeSkill, history, systemPrompt, allowedTools)
	if err != nil {
		observability.Logger("telegram.pipeline").Error("external conversation execution failed", slog.Any("err", err))
		_ = SendError(bc.bot, chat, err.Error())
		return err
	}
	finalAnswer = sanitizeAssistantOutputForUser(finalAnswer)
	bc.persistAssistantAnswer(session, finalAnswer)
	return bc.deliverFinalAnswerToChat(chat, finalAnswer, requiresAudio)
}

func (bc *BotController) handleMemoryCommand(c telebot.Context, session inputSession) (bool, error) {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(session.text)), "/memory") {
		return false, nil
	}
	if bc.canonical == nil {
		return true, SendContextText(c, "Memoria longa indisponivel neste runtime.")
	}

	reply, err := NewMemoryCommandHandler(bc.canonical).HandleText(context.Background(), session.senderID, session.convID, session.text)
	if err != nil {
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return true, nil
	}
	return true, SendContextText(c, reply)
}

func (bc *BotController) handleExternalMemoryCommand(chat *telebot.Chat, session inputSession) (bool, error) {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(session.text)), "/memory") {
		return false, nil
	}
	if bc.canonical == nil {
		return true, SendText(bc.bot, chat, "Memoria longa indisponivel neste runtime.")
	}

	reply, err := NewMemoryCommandHandler(bc.canonical).HandleText(context.Background(), session.senderID, session.convID, session.text)
	if err != nil {
		_ = SendError(bc.bot, chat, err.Error())
		return true, err
	}
	return true, SendText(bc.bot, chat, reply)
}

func (bc *BotController) persistIncomingContext(session inputSession, senderUserID int64) error {
	if err := bc.memory.EnsureConversation(session.ctx, session.convID, senderUserID, "gemini"); err != nil {
		return err
	}
	if bc.canonical != nil {
		_ = bc.canonical.ApplyFactsAndSync(session.ctx, session.senderID, extractFactsFromConversation(session.text, session.senderID))
	}
	_ = persistConversationNote(session.ctx, bc.memory, session.convID, session.text)
	_ = bc.memory.AddMessage(session.ctx, session.convID, "user", session.text)
	_ = bc.memory.AddArchiveEntry(session.ctx, memory.ArchiveEntry{
		ConversationID: session.convID,
		SessionID:      session.convID,
		Role:           "user",
		Content:        session.text,
		MessageType:    "chat",
	})
	return nil
}

func (bc *BotController) prepareExecution(session inputSession) (*skill.Skill, []agent.Message, string, []string, error) {
	logger := observability.Logger("telegram.pipeline")
	skills, _ := bc.loader.LoadAll()
	targetSkill, err := bc.router.Route(session.ctx, session.text, skills)
	if err != nil {
		logger.Warn("router non-fatal error", slog.Any("err", err))
	}

	history, err := bc.buildAgentHistory(session)
	if err != nil {
		return nil, nil, "", nil, err
	}

	// Se a sessão já possui partes (ex: imagens), usamos a mensagem estruturada
	if len(session.message.Parts) > 0 {
		// Garantimos que o texto principal está sincronizado se houver apenas texto em Parts[0]
		if len(session.message.Parts) == 1 && session.message.Parts[0].Type == agent.ContentPartText {
			session.message.Parts[0].Text = session.text
		}
	} else {
		// Caso contrário, criamos a mensagem simples de texto para o executor
		session.message = agent.Message{Role: "user", Content: session.text}
	}

	activeSkill := resolveActiveSkill(skills, targetSkill)
 
	// Se houver partes multimodais (imagens), removemos a última versão textual
	// salvada no banco e injetamos a versão completa no histórico para o LLM.
	if len(session.message.Parts) > 0 {
		if len(history) > 0 && history[len(history)-1].Role == "user" {
			history = history[:len(history)-1]
		}
		history = append(history, session.message)
	}

	systemPrompt, allowedTools := bc.resolveExecutionPrompt(session)
	return activeSkill, history, systemPrompt, allowedTools, nil
}

func (bc *BotController) buildAgentHistory(session inputSession) ([]agent.Message, error) {
	history, err := bc.memory.GetRecentMessages(session.ctx, session.convID)
	if err != nil {
		return nil, err
	}

	agentHistory := make([]agent.Message, 0, len(history))
	for _, m := range history {
		agentHistory = append(agentHistory, agent.Message{Role: m.Role, Content: m.Content})
	}
	return agentHistory, nil
}

func resolveActiveSkill(skills map[string]skill.Skill, targetSkill string) *skill.Skill {
	if targetSkill == "" {
		return nil
	}
	observability.Logger("telegram.pipeline").Info("router selected skill", slog.String("skill", targetSkill))
	s, ok := skills[targetSkill]
	if !ok {
		return nil
	}
	return &s
}

func (bc *BotController) resolveExecutionPrompt(session inputSession) (string, []string) {
	if bc.canonical == nil {
		return defaultSystemPrompt, nil
	}

	prompt, tools, err := bc.canonical.BuildPromptForQuery(session.ctx, session.senderID, session.convID, session.text)
	if err != nil {
		observability.Logger("telegram.pipeline").Warn("falling back to default prompt", slog.Any("err", err))
		return defaultSystemPrompt, nil
	}
	return prompt, tools
}

func (bc *BotController) executeConversation(c telebot.Context, session inputSession, activeSkill *skill.Skill, history []agent.Message, systemPrompt string, allowedTools []string) (string, error) {
	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	_, finalAnswer, err := bc.executor.Execute(session.ctx, systemPrompt, activeSkill, history, allowedTools)
	return finalAnswer, err
}

func (bc *BotController) executeExternalConversation(chat *telebot.Chat, session inputSession, activeSkill *skill.Skill, history []agent.Message, systemPrompt string, allowedTools []string) (string, error) {
	stopTyping := startChatActionLoop(bc.bot, chat, telebot.Typing, 4*time.Second)
	defer stopTyping()

	_, finalAnswer, err := bc.executor.Execute(session.ctx, systemPrompt, activeSkill, history, allowedTools)
	return finalAnswer, err
}

func (bc *BotController) persistAssistantAnswer(session inputSession, finalAnswer string) {
	_ = bc.memory.AddMessage(session.ctx, session.convID, "assistant", finalAnswer)
	_ = bc.memory.AddArchiveEntry(session.ctx, memory.ArchiveEntry{
		ConversationID: session.convID,
		SessionID:      session.convID,
		Role:           "assistant",
		Content:        finalAnswer,
		MessageType:    "chat",
	})
}

func (bc *BotController) deliverFinalAnswer(c telebot.Context, finalAnswer string, requiresAudio bool) error {
	// 1. Always send the text first
	if err := SendText(bc.bot, c.Chat(), finalAnswer); err != nil {
		return err
	}

	// 2. If TTS is available, send the audio as a follow-up
	// TTS is now Kokoro (pt-br feminine voice)
	if bc.tts != nil && bc.tts.IsAvailable() {
		return sendAudioWithSender(bc.bot, c.Chat(), bc.tts, finalAnswer)
	}
	return nil
}

func (bc *BotController) deliverFinalAnswerToChat(chat *telebot.Chat, finalAnswer string, requiresAudio bool) error {
	// 1. Always send the text first
	if err := SendText(bc.bot, chat, finalAnswer); err != nil {
		return err
	}

	// 2. If TTS is available, send the audio as a follow-up
	// TTS is now Kokoro (pt-br feminine voice)
	if bc.tts != nil && bc.tts.IsAvailable() {
		return sendAudioWithSender(bc.bot, chat, bc.tts, finalAnswer)
	}
	return nil
}
