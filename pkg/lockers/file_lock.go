// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package lockers

import (
	"fmt"
	"os"
	"syscall"

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
	file     *os.File
	logger   *logger.Logger
}

// NewFileLock create new file lock instance
func NewFileLock(fileName string) FileLock {
	return &fileLock{
		fileName: fileName,
		logger:   logger.GetLogger("pkg/lockers", fmt.Sprintf("FileLock(%s)", fileName)),
	}
}

// Lock try locking file, return err if fails.
func (l *fileLock) Lock() error {
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
func (l *fileLock) Unlock() error {
	defer func() {
		if err := os.Remove(l.fileName); nil != err {
			l.logger.Error("remove file lock error", logger.Error(err))
		}
		l.logger.Info("remove file lock successfully")
	}()

	defer func() {
		if err := l.file.Close(); nil != err {
			l.logger.Error("close file lock error", logger.Error(err))
		}
	}()
	return syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
}
