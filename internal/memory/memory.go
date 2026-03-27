package memory

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Message represents a single chat message
type Message struct {
	ID             int64
	ConversationID string
	Role           string
	Content        string
	CreatedAt      time.Time
}

// Fact represents a durable canonical fact stored between sessions.
type Fact struct {
	ID        int64
	Scope     string
	EntityID  string
	Key       string
	Value     string
	Source    string
	UpdatedAt time.Time
}

// Note represents a compact long-term memory note.
type Note struct {
	ID             int64
	ConversationID string
	Topic          string
	Summary        string
	Kind           string
	Importance     int
	Source         string
	CreatedAt      time.Time
}

// ArchiveEntry stores raw long-term conversation history outside the working window.
type ArchiveEntry struct {
	ID             int64
	ConversationID string
	SessionID      string
	Role           string
	Content        string
	MessageType    string
	CreatedAt      time.Time
}

// Summarizer define a interface mínima para compressão de contexto sem importar o pacote agent.
type Summarizer interface {
	Summarize(ctx context.Context, history string) (string, error)
}

// MemoryManager handles interactions with the SQLite database
type MemoryManager struct {
	db               *sql.DB
	memoryWindowSize int
	llm              Summarizer
}

// NewMemoryManager creates a new MemoryManager instance
func NewMemoryManager(db *sql.DB, llm Summarizer) *MemoryManager {
	return &MemoryManager{
		db:               db,
		memoryWindowSize: 200,
		llm:              llm,
	}
}

// Compress condensa o histórico de mensagens para economizar tokens.
func (m *MemoryManager) Compress(ctx context.Context, sessionID string) error {
	if m.llm == nil {
		return nil
	}

	messages, err := m.GetRecentMessages(ctx, sessionID)
	if err != nil {
		return err
	}

	if len(messages) < 15 {
		return nil
	}

	slog.Info("iniciando compressão de contexto SOTA 2026", "session", sessionID)

	toSummarize := messages[:10]
	var historyBuilder strings.Builder
	for _, msg := range toSummarize {
		historyBuilder.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	summary, err := m.llm.Summarize(ctx, historyBuilder.String())
	if err != nil {
		return fmt.Errorf("falha ao sumarizar para compressão: %w", err)
	}

	summaryMsg := fmt.Sprintf("[📦 RESUMO DE CONTEXTO ANTERIOR]: %s", summary)
	return m.AddMessage(ctx, sessionID, "system", summaryMsg)
}

func (m *MemoryManager) Close() error { return m.db.Close() }
func (m *MemoryManager) DB() *sql.DB  { return m.db }
