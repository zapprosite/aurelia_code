package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// InstanceLockMetadata captures the process currently holding the runtime lock.
type InstanceLockMetadata struct {
	PID       int       `json:"pid"`
	Command   string    `json:"command"`
	Root      string    `json:"root"`
	StartedAt time.Time `json:"started_at"`
}

// InstanceLock keeps the lock file descriptor alive for the process lifetime.
type InstanceLock struct {
	file     *os.File
	path     string
	metadata InstanceLockMetadata
}

// AcquireInstanceLock enforces a single active Aurelia runtime per instance root.
func AcquireInstanceLock(r *PathResolver, argv []string) (*InstanceLock, error) {
	lockPath := r.InstanceLock()
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("runtime: open instance lock: %w", err)
	}

	if err := tryFileLock(file); err != nil {
		meta := readInstanceLockMetadata(file)
		_ = file.Close()
		if meta.PID != 0 {
			return nil, fmt.Errorf(
				"another Aurelia instance is already running (pid=%d, command=%q, started_at=%s, lock=%s)",
				meta.PID,
				meta.Command,
				meta.StartedAt.UTC().Format(time.RFC3339),
				lockPath,
			)
		}
		return nil, fmt.Errorf("another Aurelia instance is already running (lock=%s)", lockPath)
	}

	lock := &InstanceLock{
		file: file,
		path: lockPath,
		metadata: InstanceLockMetadata{
			PID:       os.Getpid(),
			Command:   strings.TrimSpace(strings.Join(argv, " ")),
			Root:      r.Root(),
			StartedAt: time.Now().UTC(),
		},
	}
	if err := lock.writeMetadata(); err != nil {
		_ = releaseFileLock(file)
		_ = file.Close()
		return nil, err
	}

	return lock, nil
}

// Path returns the lockfile path for diagnostics.
func (l *InstanceLock) Path() string {
	if l == nil {
		return ""
	}
	return l.path
}

// Metadata returns the current lock metadata.
func (l *InstanceLock) Metadata() InstanceLockMetadata {
	if l == nil {
		return InstanceLockMetadata{}
	}
	return l.metadata
}

// Release unlocks and closes the lock file.
func (l *InstanceLock) Release() error {
	if l == nil || l.file == nil {
		return nil
	}

	unlockErr := releaseFileLock(l.file)
	closeErr := l.file.Close()
	l.file = nil

	if unlockErr != nil {
		return fmt.Errorf("runtime: release instance lock: %w", unlockErr)
	}
	if closeErr != nil {
		return fmt.Errorf("runtime: close instance lock: %w", closeErr)
	}
	return nil
}

func (l *InstanceLock) writeMetadata() error {
	if l == nil || l.file == nil {
		return fmt.Errorf("runtime: instance lock is not initialized")
	}

	if err := l.file.Truncate(0); err != nil {
		return fmt.Errorf("runtime: truncate instance lock: %w", err)
	}
	if _, err := l.file.Seek(0, 0); err != nil {
		return fmt.Errorf("runtime: seek instance lock: %w", err)
	}

	payload, err := json.MarshalIndent(l.metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("runtime: marshal instance lock metadata: %w", err)
	}
	payload = append(payload, '\n')
	if _, err := l.file.Write(payload); err != nil {
		return fmt.Errorf("runtime: write instance lock metadata: %w", err)
	}
	if err := l.file.Sync(); err != nil {
		return fmt.Errorf("runtime: sync instance lock metadata: %w", err)
	}
	return nil
}

func readInstanceLockMetadata(file *os.File) InstanceLockMetadata {
	if file == nil {
		return InstanceLockMetadata{}
	}
	if _, err := file.Seek(0, 0); err != nil {
		return InstanceLockMetadata{}
	}
	payload, err := os.ReadFile(file.Name())
	if err != nil {
		return InstanceLockMetadata{}
	}
	var meta InstanceLockMetadata
	if err := json.Unmarshal(payload, &meta); err != nil {
		return InstanceLockMetadata{}
	}
	return meta
}
