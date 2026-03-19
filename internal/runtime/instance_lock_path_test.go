package runtime

import (
	"path/filepath"
	"testing"
)

func TestPathResolver_InstanceLock(t *testing.T) {
	base := t.TempDir()
	r := &PathResolver{root: base}

	if got := r.InstanceLock(); got != filepath.Join(base, "instance.lock") {
		t.Fatalf("InstanceLock() = %q", got)
	}
}
