package lockers

import (
	"fmt"

	"github.com/eleme/lindb/pkg/logger"

	"github.com/gofrs/flock"
	"go.uber.org/zap"
)

// FileLock is file lock
type FileLock struct {
	lock *flock.Flock
	log  *zap.Logger
}

// NewFileLock create new file lock instance
func NewFileLock(fileName string) *FileLock {
	return &FileLock{
		lock: flock.New(fileName),
		log:  logger.GetLogger(),
	}
}

// Lock try locking file, return err if fails.
func (l *FileLock) Lock() error {
	if l.lock.Locked() {
		return fmt.Errorf("cannot lock twice, file[%s]", l.lock.Path())
	}
	return l.lock.Lock()
}

// RLock try locking file, return err if fails.
func (l *FileLock) RLock() error {
	if l.lock.RLocked() {
		return fmt.Errorf("cannot rlock twice, file[%s]", l.lock.Path())
	}
	return l.lock.RLock()
}

// IsLocked detects if this file is locked
func (l *FileLock) IsLocked() bool {
	return l.lock.Locked()
}

// IsRLocked detects if this file is read locked
func (l *FileLock) IsRLocked() bool {
	return l.lock.RLocked()
}

// Unlock unlock file lock, if fail return err
func (l *FileLock) Unlock() error {
	if !l.lock.Locked() {
		return fmt.Errorf("no lock to unlock, file[%s]", l.lock.Path())
	}
	return l.lock.Unlock()
}
