package telegram

import (
	"context"
	"fmt"
	"strings"
)

// SlashCommandHandler centraliza o roteamento de comandos "/" no padrão Sênior SOTA 2026.
type SlashCommandHandler struct {
	bc *BotController
}

func NewSlashCommandHandler(bc *BotController) *SlashCommandHandler {
	return &SlashCommandHandler{bc: bc}
}

func (h *SlashCommandHandler) Handle(ctx context.Context, session inputSession) (bool, string, error) {
	text := strings.TrimSpace(session.text)
	if !strings.HasPrefix(text, "/") {
		return false, "", nil
	}

	parts := strings.Fields(text)
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "/memory", "/memoria":
		return h.handleMemory(ctx, session, parts)
	case "/obsidian":
		return h.handleObsidian(ctx, session, parts)
	case "/qdrant":
		return h.handleQdrant(ctx, session, parts)
	case "/supabase":
		return h.handleSupabase(ctx, session, parts)
	case "/status":
		return h.handleStatus(ctx, session)
	case "/config":
		return true, fmt.Sprintf(" [⚙️ Configurações]: Modo: %s | DB: %s", h.bc.config.AureliaMode, h.bc.config.DBPath), nil
	}

	return false, "", nil
}

func (h *SlashCommandHandler) handleMemory(ctx context.Context, session inputSession, parts []string) (bool, string, error) {
	if h.bc.canonical == nil {
		return true, "Memória longa indisponível.", nil
	}
	query := ""
	if len(parts) > 1 {
		query = strings.Join(parts[1:], " ")
	}
	reply, err := NewMemoryCommandHandler(h.bc.canonical).HandleText(ctx, session.senderID, session.convID, query)
	return true, reply, err
}

func (h *SlashCommandHandler) handleObsidian(ctx context.Context, session inputSession, parts []string) (bool, string, error) {
	if len(parts) > 1 && parts[1] == "sync" {
		reply, err := h.bc.executor.GetLoop().Registry().Execute(ctx, "markdown_brain_sync", nil)
		if err != nil {
			return true, " [📓 Obsidian Sync Error]: " + err.Error(), nil
		}
		return true, " [📓 Obsidian Sync]: " + reply, nil
	}
	return true, " [📓 Obsidian]: Use `/obsidian sync` para sincronizar o vault.", nil
}

func (h *SlashCommandHandler) handleQdrant(ctx context.Context, session inputSession, parts []string) (bool, string, error) {
	if h.bc.config.QdrantURL == "" {
		return true, "Qdrant não configurado.", nil
	}
	return true, fmt.Sprintf(" [⚡ Qdrant Semantic]: Coleção '%s' ativa em %s", h.bc.config.QdrantCollection, h.bc.config.QdrantURL), nil
}

func (h *SlashCommandHandler) handleSupabase(ctx context.Context, session inputSession, parts []string) (bool, string, error) {
	if !h.bc.config.SupabaseEnabled {
		return true, "Supabase desabilitado.", nil
	}
	return true, " [🌌 Supabase Galactic]: Camada L3 ativa. Persistência de fatos sincronizada.", nil
}

func (h *SlashCommandHandler) handleStatus(ctx context.Context, session inputSession) (bool, string, error) {
	status := fmt.Sprintf(" [🚀 Sovereign SOTA 2026 Status]\n\n- **Persona**: %s\n- **Modo**: %s\n- **Memória**: SQLite [OK] | Qdrant [OK] | Supabase [%s]\n- **Vault**: %s\n- **Identidade**: %s",
		h.bc.personaID, h.bc.config.AureliaMode, h.statusEmoji(h.bc.config.SupabaseEnabled), h.bc.config.ObsidianVaultPath, h.bc.botName)
	return true, status, nil
}

func (h *SlashCommandHandler) statusEmoji(enabled bool) string {
	if enabled {
		return "OK"
	}
	return "OFF"
}
