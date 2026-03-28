package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertAndGetFact(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	err := mm.UpsertFact(ctx, Fact{
		Scope:    "user",
		EntityID: "123",
		Key:      "user.name",
		Value:    "Rafael",
		Source:   "persona",
	})
	require.NoError(t, err)

	fact, ok, err := mm.GetFact(ctx, "user", "123", "user.name")
	require.NoError(t, err)
	require.True(t, ok, "expected fact to exist")
	assert.Equal(t, "Rafael", fact.Value)
}

func TestUpsertFact_UpdatesExistingValue(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	initial := Fact{
		Scope:    "agent",
		EntityID: "default",
		Key:      "agent.role",
		Value:    "Team Lead",
		Source:   "persona",
	}
	updated := Fact{
		Scope:    "agent",
		EntityID: "default",
		Key:      "agent.role",
		Value:    "Chief Architect",
		Source:   "memory",
	}

	require.NoError(t, mm.UpsertFact(ctx, initial))
	require.NoError(t, mm.UpsertFact(ctx, updated))

	fact, ok, err := mm.GetFact(ctx, "agent", "default", "agent.role")
	require.NoError(t, err)
	require.True(t, ok, "expected fact to exist")
	assert.Equal(t, "Chief Architect", fact.Value)
	assert.Equal(t, "memory", fact.Source)
}

func TestGetFact_MissingFact(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	_, ok, err := mm.GetFact(ctx, "user", "999", "user.name")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestAddAndListNotes(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	err := mm.AddNote(ctx, Note{
		ConversationID: "conv-1",
		Topic:          "architecture",
		Summary:        "Decidido manter SQLite com facts e notes.",
		Kind:           "decision",
		Importance:     8,
		Source:         "conversation",
	})
	require.NoError(t, err)

	notes, err := mm.ListRecentNotes(ctx, "conv-1", 5)
	require.NoError(t, err)
	require.Len(t, notes, 1)
	assert.Equal(t, "Decidido manter SQLite com facts e notes.", notes[0].Summary)
}

func TestAddNote_DeduplicatesSameNote(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	note := Note{
		ConversationID: "conv-1",
		Topic:          "architecture",
		Summary:        "Decidido manter SQLite com facts e notes.",
		Kind:           "decision",
		Importance:     8,
		Source:         "conversation",
	}

	require.NoError(t, mm.AddNote(ctx, note))
	require.NoError(t, mm.AddNote(ctx, note))

	notes, err := mm.ListRecentNotes(ctx, "conv-1", 10)
	require.NoError(t, err)
	assert.Len(t, notes, 1, "expected deduplicated notes length 1")
}

func TestAddNote_ConsolidatesByTopicAndKind(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	first := Note{
		ConversationID: "conv-1",
		Topic:          "architecture",
		Summary:        "Decidido manter SQLite.",
		Kind:           "decision",
		Importance:     6,
		Source:         "conversation",
	}
	second := Note{
		ConversationID: "conv-1",
		Topic:          "architecture",
		Summary:        "Vamos usar facts e notes.",
		Kind:           "decision",
		Importance:     8,
		Source:         "conversation",
	}

	require.NoError(t, mm.AddNote(ctx, first))
	require.NoError(t, mm.AddNote(ctx, second))

	notes, err := mm.ListRecentNotes(ctx, "conv-1", 10)
	require.NoError(t, err)
	require.Len(t, notes, 1, "expected one consolidated note")
	assert.Contains(t, notes[0].Summary, "Decidido manter SQLite.")
	assert.Contains(t, notes[0].Summary, "Vamos usar facts e notes.")
	assert.Equal(t, 8, notes[0].Importance)
}

func TestAddAndListArchiveEntries(t *testing.T) {
	t.Parallel()
	mm := setupTestDB(t)
	ctx := context.Background()

	err := mm.AddArchiveEntry(ctx, ArchiveEntry{
		ConversationID: "conv-1",
		SessionID:      "session-a",
		Role:           "user",
		Content:        "Quero manter isso minimalista.",
		MessageType:    "chat",
	})
	require.NoError(t, err)

	err = mm.AddArchiveEntry(ctx, ArchiveEntry{
		ConversationID: "conv-1",
		SessionID:      "session-a",
		Role:           "assistant",
		Content:        "Entendido.",
		MessageType:    "chat",
	})
	require.NoError(t, err)

	entries, err := mm.ListArchiveEntries(ctx, "conv-1", 10)
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "Quero manter isso minimalista.", entries[0].Content)
	assert.Equal(t, "assistant", entries[1].Role)
}
