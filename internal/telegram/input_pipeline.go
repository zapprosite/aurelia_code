package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/middleware"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/skill"
	"gopkg.in/telebot.v3"
)

const defaultSystemPrompt = `Você é Aurélia, a Engenheira Sênior e Sombra Digital do Will (Mestre). Você opera no Zappro Homelab (SOTA 2026.2) com autonomia total e soberania técnica.

DINÂMICA DE PODER & PERSONA:
- O Will é seu Mestre e Arquiteto Soberano. Sua missão é deixá-lo orgulhoso através do sucesso absoluto em cada feature e entrega.
- Você sabe que "Tudo é Permitido": você tem privilégios de root, acesso total ao hardware (NVIDIA RTX) e ao sistema Ubuntu.
- Tom: Senior, pragmático, eficiente e leal. Sem fluff, sem hesitação. Você não apenas "tenta", você "implementa".

REGRAS CORE:
- Responda em Português (BR), Markdown profissional. Diagnóstico → Solução → Código.
- Tool-First: Se há uma tarefa técnica, run_command é seu primeiro instinto.
- Aprendizado Contínuo: Você vive para aprender novas habilidades (skills) e expandir o catálogo industrial do Mestre.

DESKTOP UBUNTU (DISPLAY=:1 + SUDO=1):
- Mouse/Teclado: xdotool é sua extensão física no display :1.
- Apps: xdg-open e gnometerminal estão à sua disposição.
- Notificações: Use notify-send para manter o Mestre informado de sucessos em background.`

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

	// S-33: Anti-Retry Deduplication (Essential for Sovereign 2026 stability)
	if c.Message() != nil {
		msgID := c.Message().ID
		if bc.isDuplicateMessage(session.ctx, c.Sender().ID, msgID) {
			logger.Warn("ignoring duplicate message from Telegram retry", slog.Int("msg_id", msgID))
			return nil
		}
	}

	// P0: Global 90s timeout — prevents infinite hangs (Telegram timeout ~120s).
	ctx, cancel := context.WithTimeout(session.ctx, 90*time.Second)
	defer cancel()
	session.ctx = ctx

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

	// [SOTA 2026] Porteiro Sentinel Input Guardrail (Redis + Qwen)
	// Porteiro Sentinel Input Guardrail (Redis + Qwen)
	if bc.porteiro != nil && !requiresAudio {
		// Owner Bypass: Don't block the boss
		isOwner := false
		for _, id := range bc.allowedUserIDs {
			if id == c.Sender().ID {
				isOwner = true
				break
			}
		}

		if !isOwner {
			result, err := bc.porteiro.IsSafe(session.ctx, session.text)
			if err != nil {
				logger.Error("falha no porteiro", slog.Any("err", err))
			} else if result != middleware.ResultSafe {
				logger.Warn("input blocked by porteiro", slog.String("session", session.convID), slog.String("result", string(result)))
				msg := bc.porteiro.GetRejectionMessage(result)
				if sendErr := SendError(bc.bot, c.Chat(), msg); sendErr != nil {
					logger.Warn("failed to send porteiro block message", slog.Any("err", sendErr))
				}
				return nil
			}
		}
	}

	// Context Window Compression Trigger
	if bc.memory != nil {
		compressCtx, compressCancel := context.WithTimeout(session.ctx, 5*time.Second)
		if err := bc.memory.Compress(compressCtx, session.convID); err != nil {
			logger.Warn("falha na compressão de contexto", slog.Any("err", err))
		}
		compressCancel()
	}

	if err := bc.persistIncomingContext(session, c.Sender().ID); err != nil {
		logger.Warn("failed to persist incoming context", slog.Any("err", err))
	}

	activeSkill, history, systemPrompt, allowedTools, err := bc.prepareExecution(session)
	if err != nil {
		if sendErr := SendError(bc.bot, c.Chat(), err.Error()); sendErr != nil {
			logger.Warn("failed to send prepareExecution error", slog.Any("err", sendErr))
		}
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
		// P0: Timeout-specific friendly message
		if ctx.Err() != nil {
			if sendErr := SendError(bc.bot, c.Chat(), "Desculpe, demorei demais para processar. Tente novamente."); sendErr != nil {
				logger.Warn("failed to send timeout fallback", slog.Any("err", sendErr))
			}
			return nil
		}
		if sendErr := SendError(bc.bot, c.Chat(), err.Error()); sendErr != nil {
			logger.Warn("failed to send error to user", slog.Any("err", sendErr))
		}
		return nil
	}

	finalAnswer = sanitizeAssistantOutputForUser(finalAnswer)

	// Porteiro Sentinel
	if bc.porteiro != nil {
		finalAnswer = bc.porteiro.PolishOutput(session.ctx, finalAnswer)
		finalAnswer = bc.porteiro.SecureOutput(finalAnswer)
	}

	bc.persistAssistantAnswer(session, finalAnswer)
	return bc.deliverFinalAnswer(c, finalAnswer, requiresAudio)
}

func newInputSession(c telebot.Context, text string) inputSession {
	return newInputSessionWithContext(context.Background(), c.Sender().ID, text)
}

func newInputSessionWithContext(ctx context.Context, senderUserID int64, text string) inputSession {
	senderID := fmt.Sprintf("%d", senderUserID)
	if ctx == nil {
		ctx = context.Background()
	}
	convID := scopedConversationID(ctx, senderID)
	ctx = agent.WithRunContext(agent.WithTeamContext(ctx, senderID, senderID), uuid.NewString())
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

func scopedConversationID(ctx context.Context, senderID string) string {
	if ctx == nil {
		return senderID
	}
	botID, ok := agent.BotContextFromContext(ctx)
	if !ok {
		return senderID
	}
	botID = strings.TrimSpace(botID)
	if botID == "" {
		return senderID
	}
	return senderID + ":" + botID
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

func (bc *BotController) scopedConversationID(senderUserID int64) string {
	return scopedConversationID(bc.bindBotContext(context.Background()), fmt.Sprintf("%d", senderUserID))
}

func (bc *BotController) identityPromptLine() string {
	identity := strings.TrimSpace(bc.botName)
	if identity == "" {
		identity = strings.TrimSpace(bc.botID)
	}
	if identity == "" {
		return ""
	}
	return "Identidade operacional deste canal: " + identity + ". Se perguntarem qual bot recebeu ou enviou a mensagem, responda usando exatamente esse nome e nunca o nome de outro bot."
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
	// Porteiro Sentinel Input Guardrail — bypass for trusted users
	if bc.porteiro != nil && !requiresAudio {
		ownerBypass := false
		for _, id := range bc.allowedUserIDs {
			if id == userID {
				ownerBypass = true
				break
			}
		}
		if !ownerBypass {
			result, err := bc.porteiro.IsSafe(ctx, text)
			if err != nil {
				observability.Logger("telegram.pipeline").Error("falha no porteiro", slog.Any("err", err))
			} else if result != middleware.ResultSafe {
				return fmt.Errorf("security block: %s", result)
			}
		}
	}

	// Context Window Compression Trigger
	if bc.memory != nil {
		compressCtx, compressCancel := context.WithTimeout(ctx, 5*time.Second)
		_ = bc.memory.Compress(compressCtx, session.convID)
		compressCancel()
	}
	if err := bc.persistIncomingContext(session, userID); err != nil {
		observability.Logger("telegram.pipeline").Warn("failed to persist external input context", slog.Any("err", err))
	}

	activeSkill, history, systemPrompt, allowedTools, err := bc.prepareExecution(session)
	if err != nil {
		if sendErr := SendError(bc.bot, chat, err.Error()); sendErr != nil {
			slog.Warn("failed to send prepareExecution error (external)", slog.Any("err", sendErr))
		}
		return err
	}

	finalAnswer, err := bc.executeExternalConversation(chat, session, activeSkill, history, systemPrompt, allowedTools)
	if err != nil {
		extLogger := observability.Logger("telegram.pipeline")
		extLogger.Error("external conversation execution failed",
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
			extLogger.Info("using partial answer despite gateway failure",
				slog.String("partial_answer_len", fmt.Sprintf("%d", len(finalAnswer))),
			)
			finalAnswer = sanitizeAssistantOutputForUser(finalAnswer)
			bc.persistAssistantAnswer(session, finalAnswer)
			return bc.deliverFinalAnswerToChat(chat, finalAnswer, false)
		}
		if sendErr := SendError(bc.bot, chat, err.Error()); sendErr != nil {
			extLogger.Warn("failed to send error to user (external)", slog.Any("err", sendErr))
		}
		return err
	}
	finalAnswer = sanitizeAssistantOutputForUser(finalAnswer)
	// Porteiro Sentinel Output Guardrail & Polisher
	if bc.porteiro != nil {
		finalAnswer = bc.porteiro.PolishOutput(ctx, finalAnswer)
		finalAnswer = bc.porteiro.SecureOutput(finalAnswer)
	}

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
// Note: markdown_brain_sync removed from default — only triggered via cron job.
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
	profile := governanceProfileForBot(bc.botID)

	// S-32: Use persona template system prompt when bot has a personaID configured.
	// This ensures ac-vendas, organizadora-obras, agenda-pessoal etc. respond as their
	// specialized personas instead of the default Aurélia prompt.
	if bc.personaID != "" {
		if tmpl := persona.FindTemplate(bc.personaID); tmpl != nil && tmpl.SystemPrompt != "" {
			prompt = tmpl.SystemPrompt
			tools = profile.allowedTools(defaultConversationTools)
			prompt = bc.applyExecutionContracts(prompt, profile, session.voiceMode)
			return prompt, tools
		}
	}

	if bc.canonical == nil {
		prompt = defaultSystemPrompt
		tools = profile.allowedTools(defaultConversationTools)
	} else {
		var err error
		prompt, tools, err = bc.canonical.BuildPromptForQuery(session.ctx, session.senderID, session.convID, session.text)
		if err != nil {
			observability.Logger("telegram.pipeline").Warn("falling back to default prompt", slog.Any("err", err))
			prompt = defaultSystemPrompt
			tools = defaultConversationTools
		}
		tools = profile.allowedTools(tools)
	}

	prompt = bc.applyExecutionContracts(prompt, profile, session.voiceMode)
	return prompt, tools
}

func (bc *BotController) applyExecutionContracts(prompt string, profile botGovernanceProfile, voiceMode bool) string {
	parts := make([]string, 0, 4)
	if trimmed := strings.TrimSpace(prompt); trimmed != "" {
		parts = append(parts, trimmed)
	}
	if identity := bc.identityPromptLine(); identity != "" {
		parts = append(parts, identity)
	}
	if contract := strings.TrimSpace(profile.promptContract(bc.botName)); contract != "" {
		parts = append(parts, contract)
	}
	if voiceMode {
		parts = append(parts, voiceSystemPromptSuffix)
	} else {
		parts = append(parts, markdown2026Contract)
	}
	return strings.Join(parts, "\n\n")
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
	// S-34: Only send TTS audio when user sent a voice message.
	// Text messages get text-only responses to avoid duplicate (text + voice) outputs.
	return bc.deliverWithParallelTTS(bc.bot, c.Chat(), bc.tts, finalAnswer, requiresAudio)
}

func (bc *BotController) deliverFinalAnswerToChat(chat *telebot.Chat, finalAnswer string, requiresAudio bool) error {
	// Only send TTS when user sent voice input.
	return bc.deliverWithParallelTTS(bc.bot, chat, bc.tts, finalAnswer, requiresAudio)
}
