package page

import (
	"sync/atomic"

	"github.com/eleme/lindb/pkg/fileutil"
)

// MappedPage represents a holder for mmap bytes.
// One MappedPage corresponds to a mapped file.
type MappedPage interface {
	// FilePath returns mapped filePath.
	FilePath() string
	// Buffer returns remaining mmap bytes starting from position.
	Buffer(position int) []byte
	// Data returns mmap bytes starting from position with length bytes.
	Data(position, length int) []byte
	// Sync syncs page to persist storage.
	Sync() error
	// Close releases underlying bytes.
	Close() error
	// Closed returns if the mmap bytes is closed.
	Closed() bool
	// Size returns the size of underlying bytes.
	Size() int
}

// CloseFunc defines the way to release underlying bytes.
// This makes it easy to use normal non-mmapped bytes for test.
type CloseFunc func(mappedBytes []byte) error

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
	mappedBytes []byte
	// 0 -> opened, 1 -> closed
	closed    int32
	closeFunc CloseFunc
	syncFunc  SyncFunc
}

// NewMappedPage returns a new MappedPage wrapping the give bytes.
func NewMappedPage(fileName string, bytes []byte, closeFunc CloseFunc, syncFunc SyncFunc) MappedPage {
	return &mappedPage{
		fileName:    fileName,
		mappedBytes: bytes,
		closed:      0,
		closeFunc:   closeFunc,
		syncFunc:    syncFunc,
	}
}

// FilePath returns mapped filePath.
func (mp *mappedPage) FilePath() string {
	return mp.fileName
}

// Buffer returns remaining mmap bytes starting from position.
func (mp *mappedPage) Buffer(position int) []byte {
	return mp.mappedBytes[position:]
}

// Data returns length mmap bytes starting from position.
func (mp *mappedPage) Data(position, length int) []byte {
	return mp.mappedBytes[position : position+length]
}

// Sync syncs page to persist storage.
func (mp *mappedPage) Sync() error {
	return mp.syncFunc(mp.mappedBytes)
}

// Close releases underlying bytes.
func (mp *mappedPage) Close() error {
	if atomic.CompareAndSwapInt32(&mp.closed, 0, 1) {
		return mp.closeFunc(mp.mappedBytes)
	}
	return nil
}

// Closed returns if the mmap bytes is closed.
func (mp *mappedPage) Closed() bool {
	return atomic.LoadInt32(&mp.closed) == 1
}

// Size returns the size of underlying bytes.
func (mp *mappedPage) Size() int {
	return cap(mp.mappedBytes)
}
