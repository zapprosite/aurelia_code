package runtime

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAcquireInstanceLockWritesMetadata(t *testing.T) {
	root := t.TempDir()
	r := &PathResolver{root: root}

	lock, err := AcquireInstanceLock(r, []string{"aurelia-elite", "--foreground"})
	if err != nil {
		t.Fatalf("AcquireInstanceLock() error = %v", err)
	}
	defer func() {
		if err := lock.Release(); err != nil {
			t.Fatalf("Release() error = %v", err)
		}
	}()

	payload, err := os.ReadFile(r.InstanceLock())
	if err != nil {
		t.Fatalf("ReadFile(lock) error = %v", err)
	}

	var meta InstanceLockMetadata
	if err := json.Unmarshal(payload, &meta); err != nil {
		t.Fatalf("metadata json error = %v", err)
	}
	if meta.PID == 0 {
		t.Fatal("expected lock metadata PID to be set")
	}
	if meta.Root != root {
		t.Fatalf("metadata root = %q, want %q", meta.Root, root)
	}
	if !strings.Contains(meta.Command, "aurelia-elite") {
		t.Fatalf("metadata command = %q", meta.Command)
	}
}

func TestAcquireInstanceLockRejectsSecondHolder(t *testing.T) {
	if os.Getenv("AURELIA_LOCK_HELPER") == "1" {
		r := &PathResolver{root: os.Getenv("AURELIA_LOCK_ROOT")}
		lock, err := AcquireInstanceLock(r, []string{"aurelia-elite", "--helper"})
		if err != nil {
			os.Exit(0)
		}
		_ = lock.Release()
		t.Fatal("helper unexpectedly acquired the lock")
	}

	root := t.TempDir()
	r := &PathResolver{root: root}

	first, err := AcquireInstanceLock(r, []string{"aurelia-elite"})
	if err != nil {
		t.Fatalf("first AcquireInstanceLock() error = %v", err)
	}
	defer func() {
		if err := first.Release(); err != nil {
			t.Fatalf("Release() error = %v", err)
		}
	}()

	cmd := exec.Command(os.Args[0], "-test.run=TestAcquireInstanceLockRejectsSecondHolder")
	cmd.Env = append(os.Environ(),
		"AURELIA_LOCK_HELPER=1",
		"AURELIA_LOCK_ROOT="+root,
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("helper process should exit cleanly after failed lock attempt: %v", err)
	}
}
