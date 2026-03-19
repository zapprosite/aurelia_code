//go:build unix

package runtime

import (
	"os"
	"syscall"
)

func tryFileLock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

func releaseFileLock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}

// isProcessAlive checks if a process with the given PID still exists.
// It returns true if the process exists, false if it doesn't or the check fails.
func isProcessAlive(pid int) bool {
	// Send signal 0 to check process existence without killing it.
	// syscall.Kill returns nil if process exists, ESRCH if it doesn't.
	err := syscall.Kill(pid, 0)
	return err == nil
}
