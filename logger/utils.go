//go:build !linux
// +build !linux

package logger

import "os"

// @author valor.

// logfileDup2 returns a new log file pointer on non-linux platforms.
func logfileDup2(newLogfile *os.File, oldLogfile *os.File) (*os.File, error) {
	if oldLogfile != nil {
		_ = oldLogfile.Close()
	}
	return newLogfile, nil
}

// logfileSync: let the system decide whether to flush to file.
func logfileSync(_ *os.File) error {
	return nil
}
