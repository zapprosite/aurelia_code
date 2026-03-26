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
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/skill"
	"gopkg.in/telebot.v3"
)

const defaultSystemPrompt = "Você é Aurélia, uma assistente virtual de inteligência artificial de alto nível. Sua comunicação deve ser profissional, prestativa, objetiva e estruturada. Utilize formatação Markdown (negritos, listas e tabelas) para garantir clareza máxima. Ao analisar códigos ou realizar tarefas técnicas, priorize a precisão e sempre valide suas ações com 'run_command' antes de confirmar a conclusão."

type inputSession struct {
	senderID  string
	convID    string
	ctx       context.Context
	text      string
	message   agent.Message
	voiceMode bool // true when input came from a voice/audio message
}

func (bc *BotController) processInput(c telebot.Context, text string, requiresAudio bool) error {
	session := bc.newInputSession(c, text)
	return bc.processInputSession(c, session, requiresAudio)
}

func (bc *BotController) processInputSession(c telebot.Context, session inputSession, requiresAudio bool) error {
	logger := observability.Logger("telegram.pipeline")
	session.text = strings.ReplaceAll(session.text, "\x00", "")
	session.voiceMode = requiresAudio

	if state, ok := bc.popPendingBootstrap(c.Sender().ID); ok {
		return bc.completeBootstrapProfile(c, state, session.text)
	}

	if bc.mediaProcessor != nil && bc.mediaProcessor.IsSupportedURL(session.text) {
		return bc.handleMediaURL(c, session)
	}

	if handled, err := bc.handleMemoryCommand(c, session); handled {
		return err
	}

	// Prompt injection guard (gemma3 local pre-filter).
	// Skipped for voice messages: Telegram voice is from an authenticated user,
	// and injection via ASR transcription is unlikely. Skipping saves 1-3s latency.
	if bc.inputGuard != nil && !requiresAudio {
		if blocked, reason := bc.inputGuard.CheckWithUser(session.ctx, c.Sender().ID, bc.allowedUserIDs, session.text); blocked {
			observability.Logger("telegram.pipeline").Warn("input blocked by guard", slog.String("reason", reason))
			_ = SendError(bc.bot, c.Chat(), "Mensagem bloqueada pelo filtro de segurança: "+reason)
			return nil
		}
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
		logger.Error("conversation execution failed",
			slog.Any("err", err),
			slog.String("user", c.Sender().Username),
			slog.String("text", session.text),
			slog.String("error_type", fmt.Sprintf("%T", err)),
		)
		// Se for erro de provider/gateway, tenta usar a resposta mesmo que incompleta
		gatewayFailure := strings.Contains(err.Error(), "provider error") ||
			strings.Contains(err.Error(), "all gateway routes failed") ||
			strings.Contains(err.Error(), "route breaker open") ||
			strings.Contains(err.Error(), "budget exceeded") ||
			strings.Contains(err.Error(), "empty guarded content")
		if gatewayFailure && finalAnswer != "" {
			logger.Info("using partial answer despite gateway failure",
				slog.String("partial_answer_len", fmt.Sprintf("%d", len(finalAnswer))),
			)
			finalAnswer = sanitizeAssistantOutputForUser(finalAnswer)
			bc.persistAssistantAnswer(session, finalAnswer)
			return bc.deliverFinalAnswer(c, finalAnswer, requiresAudio)
		}
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

func (bc *BotController) newInputSession(c telebot.Context, text string) inputSession {
	return newInputSessionWithContext(bc.bindBotContext(context.Background()), c.Sender().ID, text)
}

func (bc *BotController) newInputSessionWithContext(ctx context.Context, senderUserID int64, text string) inputSession {
	return newInputSessionWithContext(bc.bindBotContext(ctx), senderUserID, text)
}

func (bc *BotController) bindBotContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	botID := strings.TrimSpace(bc.botID)
	if botID == "" {
		botID = "aurelia"
	}
	return agent.WithBotContext(ctx, botID)
}

type botChatSender struct {
	bot  *telebot.Bot
	chat *telebot.Chat
}

func (s *botChatSender) Send(what interface{}, opts ...interface{}) error {
	_, err := s.bot.Send(s.chat, what, opts...)
	return err
}

func (s *botChatSender) Chat() *telebot.Chat {
	return s.chat
}

func (bc *BotController) ProcessExternalInput(ctx context.Context, userID, chatID int64, text string, requiresAudio bool) error {
	if bc == nil || bc.bot == nil {
		return fmt.Errorf("telegram bot controller is not available")
	}
	chat := &telebot.Chat{ID: chatID}
	session := bc.newInputSessionWithContext(ctx, userID, text)
	text = strings.ReplaceAll(text, "\x00", "")
	session.text = text
	session.voiceMode = requiresAudio

	if handled, err := bc.handleExternalMemoryCommand(chat, session); handled {
		return err
	}
	if bc.mediaProcessor != nil && bc.mediaProcessor.IsSupportedURL(session.text) {
		sender := &botChatSender{bot: bc.bot, chat: chat}
		return bc.handleMediaURL(sender, session)
	}
	if bc.inputGuard != nil && !requiresAudio {
		if blocked, reason := bc.inputGuard.CheckWithUser(ctx, userID, bc.allowedUserIDs, text); blocked {
			observability.Logger("telegram.pipeline").Warn("external input blocked by guard", slog.String("reason", reason))
			_ = SendError(bc.bot, chat, "Mensagem bloqueada pelo filtro de segurança: "+reason)
			return nil
		}
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
		logger := observability.Logger("telegram.pipeline")
		logger.Error("external conversation execution failed",
			slog.Any("err", err),
			slog.String("user_id", fmt.Sprintf("%d", userID)),
			slog.String("text", session.text),
		)
		// Se for erro de provider/gateway, tenta usar a resposta mesmo que incompleta
		gatewayFailure := strings.Contains(err.Error(), "provider error") ||
			strings.Contains(err.Error(), "all gateway routes failed") ||
			strings.Contains(err.Error(), "route breaker open") ||
			strings.Contains(err.Error(), "budget exceeded") ||
			strings.Contains(err.Error(), "empty guarded content")
		if gatewayFailure && finalAnswer != "" {
			logger.Info("using partial answer despite gateway failure",
				slog.String("partial_answer_len", fmt.Sprintf("%d", len(finalAnswer))),
			)
			finalAnswer = sanitizeAssistantOutputForUser(finalAnswer)
			bc.persistAssistantAnswer(session, finalAnswer)
			return bc.deliverFinalAnswerToChat(chat, finalAnswer, false)
		}
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

// defaultConversationTools is used when no skill or persona specifies a tool list.
// Limits Groq token usage by excluding the 77 MCP schemas sent by default.
// MCP tools (github, playwright, filesystem, etc.) are still available via skills.
var defaultConversationTools = []string{
	"read_file", "write_file", "list_dir", "run_command",
	"web_search", "docker_control", "system_monitor", "service_control",
	"create_schedule", "list_schedules", "cpf_cnpj",
}

// voiceSystemPromptSuffix is appended to the system prompt when the user
// sent a voice message. It instructs the LLM to produce speech-ready text:
// natural sentences, no markdown, numbers written out, conversational tone.
const voiceSystemPromptSuffix = "\n\nATENÇÃO — MODO VOZ: Esta mensagem chegou por áudio e a resposta será sintetizada em voz. Escreva de forma completamente conversacional e natural, como se estivesse falando em voz alta. Regras obrigatórias: (1) Sem markdown — nada de asteriscos, underlines, cerquilhas, colchetes ou backticks. (2) Sem listas numeradas ou com marcadores — use frases conectadas (\"Primeiro... depois... por fim...\"). (3) Sem tabelas. (4) Números por extenso quando for natural (\"dois mil e vinte e seis\", \"trinta por cento\"). (5) Frases curtas e diretas. (6) Se precisar enumerar itens, separe por vírgula ou use \"e\" entre eles."

func (bc *BotController) resolveExecutionPrompt(session inputSession) (string, []string) {
	var prompt string
	var tools []string

	// S-32: Use persona template system prompt when bot has a personaID configured.
	// This ensures ac-vendas, organizadora-obras, agenda-pessoal etc. respond as their
	// specialized personas instead of the default Aurélia prompt.
	if bc.personaID != "" {
		if tmpl := persona.FindTemplate(bc.personaID); tmpl != nil && tmpl.SystemPrompt != "" {
			prompt = tmpl.SystemPrompt
			tools = defaultConversationTools
			if session.voiceMode {
				prompt += voiceSystemPromptSuffix
			}
			return prompt, tools
		}
	}

	if bc.canonical == nil {
		prompt = defaultSystemPrompt
		tools = defaultConversationTools
	} else {
		var err error
		prompt, tools, err = bc.canonical.BuildPromptForQuery(session.ctx, session.senderID, session.convID, session.text)
		if err != nil {
			observability.Logger("telegram.pipeline").Warn("falling back to default prompt", slog.Any("err", err))
			prompt = defaultSystemPrompt
			tools = defaultConversationTools
		}
		// Persona doesn't specify tools → use default core set to avoid 77-tool token bloat.
		if len(tools) == 0 {
			tools = defaultConversationTools
		}
	}

	if session.voiceMode {
		prompt += voiceSystemPromptSuffix
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
	// Send text and synthesize TTS in parallel; audio follows once ready.
	return deliverWithParallelTTS(bc.bot, c.Chat(), bc.tts, finalAnswer)
}

func (bc *BotController) deliverFinalAnswerToChat(chat *telebot.Chat, finalAnswer string, requiresAudio bool) error {
	// Send text and synthesize TTS in parallel; audio follows once ready.
	return deliverWithParallelTTS(bc.bot, chat, bc.tts, finalAnswer)
}
