package agent

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func newTestSQLiteTaskStore(t *testing.T) *SQLiteTaskStore {
	t.Helper()
	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "swarm.db"))
	if err != nil {
		t.Fatalf("NewSQLiteTaskStore() error = %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})
	return store
}

func TestSQLiteTaskStoreSwarmChannelLifecycle(t *testing.T) {
	ctx := context.Background()
	store := newTestSQLiteTaskStore(t)

	teamID := uuid.NewString()
	if err := store.createTeam(ctx, teamID, "team-key", "user-1", MasterAgentName); err != nil {
		t.Fatalf("createTeam() error = %v", err)
	}

	channel := SwarmChannel{
		ID:     uuid.NewString(),
		TeamID: teamID,
		Name:   "ops-live",
		Kind:   "ops",
	}
	if err := store.createSwarmChannel(ctx, channel); err != nil {
		t.Fatalf("createSwarmChannel() error = %v", err)
	}

	channels, err := store.listSwarmChannels(ctx, teamID)
	if err != nil {
		t.Fatalf("listSwarmChannels() error = %v", err)
	}
	if len(channels) != 1 || channels[0].Name != "ops-live" {
		t.Fatalf("unexpected channels %#v", channels)
	}
}

func TestSQLiteTaskStoreClaimAssistanceTask(t *testing.T) {
	ctx := context.Background()
	store := newTestSQLiteTaskStore(t)

	teamID := uuid.NewString()
	if err := store.createTeam(ctx, teamID, "team-key-2", "user-2", MasterAgentName); err != nil {
		t.Fatalf("createTeam() error = %v", err)
	}

	task := AssistanceTask{
		ID:         uuid.NewString(),
		TeamID:     teamID,
		OwnerAgent: "aurelia-chief",
		Title:      "help summarize incident",
		Body:       "Need a short summary for ops-live",
	}
	if err := store.enqueueAssistanceTask(ctx, task); err != nil {
		t.Fatalf("enqueueAssistanceTask() error = %v", err)
	}

	claimed, err := store.claimAssistanceTask(ctx, teamID, "librarian")
	if err != nil {
		t.Fatalf("claimAssistanceTask() error = %v", err)
	}
	if claimed == nil {
		t.Fatal("expected claimed assistance task")
	}
	if claimed.HelperAgent == nil || *claimed.HelperAgent != "librarian" {
		t.Fatalf("unexpected helper agent %#v", claimed.HelperAgent)
	}
	if claimed.Status != "claimed" {
		t.Fatalf("unexpected status %q", claimed.Status)
	}
}
