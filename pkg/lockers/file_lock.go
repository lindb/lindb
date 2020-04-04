package lockers

import (
	"fmt"
	"os"

	"github.com/gofrs/flock"

	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./file_lock.go -destination=./file_lock_mock.go -package lockers

// FileLock represents file lock
type FileLock interface {
	// Lock try locking file, return err if fails.
	Lock() error
	// Unlock unlock file lock, if fail return err
	Unlock() error
}

// fileLock is file lock
type fileLock struct {
	fileName string
	lock     *flock.Flock
	logger   *logger.Logger
}

// NewFileLock create new file lock instance
func NewFileLock(fileName string) FileLock {
	return &fileLock{
		fileName: fileName,
		lock:     flock.New(fileName),
		logger:   logger.GetLogger("pkg/lockers", fmt.Sprintf("FileLock[%s]", fileName)),
	}
}

// Lock try locking file, return err if fails.
func (l *fileLock) Lock() error {
	if l.lock.Locked() {
		return fmt.Errorf("file: %s is already locked", l.fileName)
	}
	locked, err := l.lock.TryLock()
	if err != nil || !locked {
		return fmt.Errorf("cannot flock file: %s - %s", l.fileName, err)
	}
	return nil
}

// Unlock unlock file lock, if fail return err
func (l *fileLock) Unlock() error {
	defer func() {
		if err := os.Remove(l.fileName); nil != err {
			l.logger.Error("remove file lock error", logger.Error(err))
		}
		l.logger.Info("remove file lock successfully")
	}()

	return l.lock.Unlock()
}
