//go:build linux && !arm64
// +build linux,!arm64

package logger

import (
	"os"
	"syscall"
)

// @author valor.

// logfileDup2 returns stderr on linux platforms.
// And let logs output to stderr, let stderr redirect to log file.
func logfileDup2(newLogfile *os.File, _ *os.File) (*os.File, error) {
	err := syscall.Dup2(int(newLogfile.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return nil, err
	}
	_ = newLogfile.Close()
	return os.Stderr, nil
}

// logfileSync: let the system decide whether to flush to file.
func logfileSync(_ *os.File) error {
	return nil
}