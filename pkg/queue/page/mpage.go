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

package page

import (
	"os"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./mpage.go -destination ./mpage_mock.go -package page

// for testing
var (
	mapFileFunc  = fileutil.RWMap
	openFileFunc = os.OpenFile
)

// MappedPage represents a holder for mmap bytes,
// one MappedPage corresponds to a mapped file.
type MappedPage interface {
	// FilePath returns mapped filePath
	FilePath() string
	// WriteBytes writes bytes data into buffer
	WriteBytes(data []byte, offset int)
	// ReadBytes reads bytes data from buffer
	ReadBytes(offset, length int) []byte
	// PutUint64 puts uint64 into buffer
	PutUint64(value uint64, offset int)
	// ReadUint64 reads uint64 from buffer
	ReadUint64(offset int) uint64
	// PutUint32 puts uint32 into buffer
	PutUint32(value uint32, offset int)
	// ReadUint32 reads uint32 from buffer
	ReadUint32(offset int) uint32
	// PutUint8 puts uint8 into buffer
	PutUint8(value uint8, offset int)
	// ReadUint8 reads uint8 from buffer
	ReadUint8(offset int) uint8
	// Sync syncs page to persist storage
	Sync() error
	// Close releases underlying bytes
	Close() error
	// Closed returns if the mmap bytes is closed
	Closed() bool
	// Size returns the size of underlying bytes
	Size() int
}

// CloseFunc defines the way to release underlying bytes.
// This makes it easy to use normal non-mmapped bytes for test.
type CloseFunc func(f *os.File, mappedBytes []byte) error

// MMapCloseFunc defines the way to release mmaped bytes.
var MMapCloseFunc CloseFunc = fileutil.Unmap

// SyncFunc defines the way to flush bytes to storage.
// This makes it easy to use normal non-mmapped bytes for test.
type SyncFunc func(mappedBytes []byte) error

// MMapSyncFunc defines the way to flush mmaped bytes.
var MMapSyncFunc SyncFunc = fileutil.Sync

// mappedPage implements MappedPage
type mappedPage struct {
	fileName    string
	f           *os.File
	mappedBytes []byte
	size        int
	// false -> opened, true -> closed
	closed atomic.Bool
}

// NewMappedPage returns a new MappedPage wrapping the give bytes.
func NewMappedPage(fileName string, size int) (MappedPage, error) {
	f, err := openFileFunc(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	bytes, err := mapFileFunc(f, size)
	if err != nil {
		// need close file, if map file failure
		_ = f.Close()
		return nil, err
	}

	return &mappedPage{
		fileName:    fileName,
		f:           f,
		mappedBytes: bytes,
		size:        size,
	}, nil
}

// FilePath returns mapped filePath.
func (mp *mappedPage) FilePath() string {
	return mp.fileName
}

// WriteBytes writes bytes data into buffer
func (mp *mappedPage) WriteBytes(data []byte, offset int) {
	copy(mp.mappedBytes[offset:], data)
}

// ReadBytes reads bytes data from buffer
func (mp *mappedPage) ReadBytes(offset, length int) []byte {
	return mp.mappedBytes[offset : offset+length]
}

// PutUint64 puts uint64 into buffer
func (mp *mappedPage) PutUint64(value uint64, offset int) {
	stream.PutUint64(mp.mappedBytes, offset, value)
}

// ReadUint64 puts uint64 into buffer
func (mp *mappedPage) ReadUint64(offset int) uint64 {
	return stream.ReadUint64(mp.mappedBytes, offset)
}

// PutUint32 puts uint32 into buffer
func (mp *mappedPage) PutUint32(value uint32, offset int) {
	stream.PutUint32(mp.mappedBytes, offset, value)
}

// ReadUint32 reads uint32 from buffer
func (mp *mappedPage) ReadUint32(offset int) uint32 {
	return stream.ReadUint32(mp.mappedBytes, offset)
}

// PutUint8 puts uint8 into buffer
func (mp *mappedPage) PutUint8(value uint8, offset int) {
	mp.mappedBytes[offset] = value
}

// ReadUint8 reads uint8 from buffer
func (mp *mappedPage) ReadUint8(offset int) uint8 {
	return mp.mappedBytes[offset]
}

// Sync syncs page to persist storage.
func (mp *mappedPage) Sync() error {
	return MMapSyncFunc(mp.mappedBytes)
}

// Close releases underlying bytes.
func (mp *mappedPage) Close() error {
	if mp.closed.CompareAndSwap(false, true) {
		// close file after unmap file.
		defer func() {
			_ = mp.f.Close()
		}()
		return MMapCloseFunc(mp.f, mp.mappedBytes)
	}

	return nil
}

// Closed returns if the mmap bytes is closed.
func (mp *mappedPage) Closed() bool {
	return mp.closed.Load()
}

// Size returns the size of underlying bytes.
func (mp *mappedPage) Size() int {
	return mp.size
}
