// +build !windows

package lockers

import (
	"os"
	"syscall"
)

func lockFile(path string, flag int) (*os.File, error) {
	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return nil, err
	}
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		return nil, err
	}
	return f, err
}

func unlockFile(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
