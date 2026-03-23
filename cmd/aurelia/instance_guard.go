package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// shouldSuppressDuplicateLaunch returns true when a duplicate-launch error
// was caused by an orphan process (PID 1 = systemd adopted it) rather than
// a legitimate service conflict. This prevents noisy exit-code 1 when
// systemd restarts the unit while an old instance is still shutting down.
func shouldSuppressDuplicateLaunch(parentPID int, err error) bool {
	if err == nil {
		return false
	}
	if !strings.Contains(err.Error(), "another Aurelia instance is already running") {
		return false
	}
	// PID 1 means the process was reparented to init/systemd — orphan.
	return parentPID == 1
}

// recordDuplicateLaunch writes a structured log entry so operators can
// audit how often orphan launches happen.
func recordDuplicateLaunch(args []string, err error) {
	home := os.Getenv("AURELIA_HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOME"), ".aurelia")
	}
	logDir := filepath.Join(home, "logs")
	_ = os.MkdirAll(logDir, 0o755)
	logPath := filepath.Join(logDir, "duplicate-launch.log")

	entry := fmt.Sprintf("[%s] args=%v error=%v\n",
		time.Now().Format(time.RFC3339), args, err)
	_ = os.WriteFile(logPath, []byte(entry), 0o644)
}

// exitCodeForBootstrapError decides the process exit code. Orphan
// duplicate launches exit 0 (suppressed); everything else exits 1.
func exitCodeForBootstrapError(logger *slog.Logger, args []string, err error) int {
	if shouldSuppressDuplicateLaunch(os.Getppid(), err) {
		recordDuplicateLaunch(args, err)
		logger.Warn("suppressed orphan duplicate launch", slog.Any("err", err))
		return 0
	}
	logger.Error("bootstrap failed", slog.Any("err", err))
	return 1
}
