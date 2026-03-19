package main

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExitCodeForBootstrapErrorSuppressesExternalDuplicateLaunch(t *testing.T) {
	t.Setenv("AURELIA_HOME", t.TempDir())

	if !shouldSuppressDuplicateLaunch(1, errors.New("acquire instance lock: another Aurelia instance is already running")) {
		t.Fatal("expected orphan duplicate launch to be suppressed")
	}

	recordDuplicateLaunch([]string{"aurelia-elite"}, errors.New("acquire instance lock: another Aurelia instance is already running"))
	logPath := filepath.Join(os.Getenv("AURELIA_HOME"), "logs", "duplicate-launch.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("expected duplicate launch log to be written: %v", err)
	}
	if !strings.Contains(string(content), "another Aurelia instance is already running") {
		t.Fatalf("expected duplicate launch log to contain lock error, got %q", string(content))
	}
}

func TestShouldSuppressDuplicateLaunchDoesNotHideServiceStartConflict(t *testing.T) {
	if shouldSuppressDuplicateLaunch(2117, errors.New("acquire instance lock: another Aurelia instance is already running")) {
		t.Fatal("expected non-orphan duplicate launch to remain fatal")
	}
}

func TestExitCodeForBootstrapErrorKeepsOtherBootstrapErrorsFatal(t *testing.T) {
	t.Setenv("AURELIA_HOME", t.TempDir())

	code := exitCodeForBootstrapError(testLogger(), []string{"aurelia-elite"}, errors.New("load config: boom"))
	if code != 1 {
		t.Fatalf("expected generic bootstrap error to exit 1, got %d", code)
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
