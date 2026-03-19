//go:build !unix

package runtime

import (
	"fmt"
	"os"
)

func tryFileLock(file *os.File) error {
	return fmt.Errorf("instance locking is unsupported on this platform")
}

func releaseFileLock(file *os.File) error {
	return nil
}
