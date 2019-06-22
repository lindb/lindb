package lockers

import (
	"fmt"
	"os"
	"syscall"

	"github.com/eleme/lindb/pkg/logger"

	"go.uber.org/zap"
)

// FileLock is file lock
type FileLock struct {
	fileName string
	file     *os.File
	logger   *zap.Logger
}

// NewFileLock create new file lock instance
func NewFileLock(fileName string) *FileLock {
	return &FileLock{
		fileName: fileName,
		logger:   logger.GetLogger(),
	}
}

// Lock try locking file, return err if fails.
func (l *FileLock) Lock() error {
	f, err := os.Create(l.fileName)
	if nil != err {
		return fmt.Errorf("cannot create file[%s] for lock err: %s", l.fileName, err)
	}
	l.file = f
	// invoke syscall for file lock
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if nil != err {
		return fmt.Errorf("cannot flock directory %s - %s", l.fileName, err)
	}
	return nil
}

// Unlock unlock file lock, if fail return err
func (l *FileLock) Unlock() error {
	defer func() {
		if err := os.Remove(l.fileName); nil != err {
			l.logger.Error("remove file lock error", zap.String("file", l.fileName), zap.Error(err))
		}
		l.logger.Info("remove file lock successfully", zap.String("file", l.fileName))
	}()

	defer func() {
		if err := l.file.Close(); nil != err {
			l.logger.Error("close file lock error", zap.String("file", l.fileName), zap.Error(err))
		}
	}()
	return syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
}
