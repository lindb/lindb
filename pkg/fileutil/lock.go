package fileutil

import (
	"errors"
	"os"
)

var (
	ErrLocked = errors.New("fileutil: file already locked")
)

type LockedFile struct {
	*os.File
}

func LockFile(path string, flag int, perm os.FileMode) (*LockedFile, error) {
	f, err := lockFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return &LockedFile{f}, nil
}

func TryLockFile(path string, flag int, perm os.FileMode) (*LockedFile, error) {
	f, err := tryLockFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return &LockedFile{f}, nil
}

func (f *LockedFile) Close() error {
	if err := unlockFile(f.File); err != nil {
		return err
	}
	return f.File.Close()
}
