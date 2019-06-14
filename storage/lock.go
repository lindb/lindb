package storage

import (
	"syscall"
	"os"
	"fmt"
	"go.uber.org/zap"
	"github.com/eleme/lindb/pkg/logger"
)

// File lock
type Lock struct {
	fileName string
	file     *os.File
	logger   *zap.Logger
}

func NewLock(fileName string) *Lock {
	return &Lock{
		fileName: fileName,
		logger:   logger.GetLogger(),
	}
}

// Lock
func (l *Lock) Lock() error {
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

// Unlock
func (l *Lock) Unlock() error {
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
