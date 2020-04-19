package page

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./mpage.go -destination ./mpage_mock.go -package page

// for testing
var (
	mapFileFunc = fileutil.RWMap
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
	size        int
	// false -> opened, true -> closed
	closed atomic.Bool
}

// NewMappedPage returns a new MappedPage wrapping the give bytes.
func NewMappedPage(fileName string, size int) (MappedPage, error) {
	bytes, err := mapFileFunc(fileName, size)
	if err != nil {
		return nil, err
	}

	return &mappedPage{
		fileName:    fileName,
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

// PutUint64 puts uint64 into buffer
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

// Sync syncs page to persist storage.
func (mp *mappedPage) Sync() error {
	return MMapSyncFunc(mp.mappedBytes)
}

// Close releases underlying bytes.
func (mp *mappedPage) Close() error {
	if mp.closed.CAS(false, true) {
		return MMapCloseFunc(mp.mappedBytes)
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
