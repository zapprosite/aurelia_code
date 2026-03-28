package memory

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/kocar/aurelia/internal/purity/alog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

type mockSummarizer struct{}

func (m *mockSummarizer) Summarize(ctx context.Context, history string) (string, error) {
	return "summary", nil
}

func setupTestDB(t *testing.T) *MemoryManager {
	alog.Configure(alog.Options{Format: "text"})
	
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)

	// Inicializar tabelas conforme SOTA 2026.1
	err = initializeDB(db)
	require.NoError(t, err, "falha ao inicializar esquema de tabelas")

	mm := NewMemoryManager(db, &mockSummarizer{})
	mm.memoryWindowSize = 5 // Override para teste
	
	t.Cleanup(func() {
		_ = mm.Close()
		_ = os.RemoveAll(tempDir)
	})
	return mm
}

func TestEnsureConversation(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	err := mm.EnsureConversation(ctx, "conv-1", 123, "provider-a")
	assert.NoError(t, err)

	// Upsert with new provider should work
	err = mm.EnsureConversation(ctx, "conv-1", 123, "provider-b")
	assert.NoError(t, err)
}

func TestAddAndGetMessages(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	_ = mm.EnsureConversation(ctx, "conv-1", 123, "provider-a")

	// Add 6 messages
	for i := 1; i <= 6; i++ {
		err := mm.AddMessage(ctx, "conv-1", "user", "msg "+string(rune('0'+i)))
		require.NoError(t, err)
	}

	msgs, err := mm.GetRecentMessages(ctx, "conv-1")
	require.NoError(t, err)

	// We set window size to 5, so we should get the last 5 messages
	assert.Len(t, msgs, 5)

	// Should be in chronological order
	assert.Equal(t, "msg 2", msgs[0].Content)
	assert.Equal(t, "msg 6", msgs[4].Content)
}

func TestNullBytesStripping(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	_ = mm.EnsureConversation(ctx, "conv-1", 123, "provider-a")

	dirtyMsg := "hello" + string([]byte{0}) + "world"
	err := mm.AddMessage(ctx, "conv-1", "user", dirtyMsg)
	require.NoError(t, err)

	msgs, err := mm.GetRecentMessages(ctx, "conv-1")
	require.NoError(t, err)
	require.Len(t, msgs, 1)

	assert.NotContains(t, msgs[0].Content, "\x00", "null byte was not stripped")
	assert.Equal(t, "helloworld", msgs[0].Content)
}

func TestGetRecentMessages_IsolatesConversations(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	_ = mm.EnsureConversation(ctx, "conv-a", 123, "provider-a")
	_ = mm.EnsureConversation(ctx, "conv-b", 456, "provider-a")

	require.NoError(t, mm.AddMessage(ctx, "conv-a", "user", "from-a"))
	require.NoError(t, mm.AddMessage(ctx, "conv-b", "user", "from-b"))

	msgsA, err := mm.GetRecentMessages(ctx, "conv-a")
	require.NoError(t, err)
	assert.Len(t, msgsA, 1)
	assert.Equal(t, "from-a", msgsA[0].Content)

	msgsB, err := mm.GetRecentMessages(ctx, "conv-b")
	require.NoError(t, err)
	assert.Len(t, msgsB, 1)
	assert.Equal(t, "from-b", msgsB[0].Content)
}
