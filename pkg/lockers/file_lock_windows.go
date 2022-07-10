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

//go:build windows

package lockers

import (
	"golang.org/x/sys/windows"
)

// ref: https://github.com/golang/go/blob/master/src/cmd/go/internal/lockedfile/internal/filelock/filelock_windows.go

type lockType uint32

const (
	writeLock lockType = windows.LOCKFILE_EXCLUSIVE_LOCK
)

const (
	reserved = 0
	allBytes = ^uint32(0)
)

// Lock try locking file, return err if fails.
func (l *fileLock) lock() error {
	ol := new(windows.Overlapped)
	err := windows.LockFileEx(windows.Handle(l.file.Fd()), uint32(writeLock), reserved, allBytes, allBytes, ol)
	if err != nil {
		return err
	}
	return err
}

// Unlock unlock file lock, if fail return err
func (l *fileLock) unlock() error {
	ol := new(windows.Overlapped)
	err := windows.UnlockFileEx(windows.Handle(l.file.Fd()), reserved, allBytes, allBytes, ol)
	if err != nil {
		return err
	}
	return nil
}
