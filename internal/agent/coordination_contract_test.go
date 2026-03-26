package agent

import "testing"

func TestCoordinationLabel_DefaultModes(t *testing.T) {
	t.Parallel()

	modes := DefaultCoordinationModes()
	if len(modes) != 3 {
		t.Fatalf("DefaultCoordinationModes() len = %d", len(modes))
	}
	if got := CoordinationLabel(modes); got != "delegation + handoff + assist" {
		t.Fatalf("CoordinationLabel() = %q", got)
	}
}
