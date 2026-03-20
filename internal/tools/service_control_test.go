package tools

import "testing"

func TestBlocksSelfServiceMutation_AllowsStatusAndLogs(t *testing.T) {
	t.Parallel()

	for _, action := range []serviceAction{serviceActionStatus, serviceActionLogs} {
		if err := blocksSelfServiceMutation(action, "aurelia.service"); err != "" {
			t.Fatalf("expected %s on aurelia.service to be allowed, got %q", action, err)
		}
	}
}

func TestBlocksSelfServiceMutation_BlocksStateChanges(t *testing.T) {
	t.Parallel()

	for _, action := range []serviceAction{serviceActionRestart, serviceActionStart, serviceActionStop} {
		if err := blocksSelfServiceMutation(action, "aurelia.service"); err == "" {
			t.Fatalf("expected %s on aurelia.service to be blocked", action)
		}
	}
}
