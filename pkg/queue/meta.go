package queue

import (
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/queue/page"
	"github.com/lindb/lindb/pkg/stream"
)

const int64Size = 8

// Meta represents a mmaped bytes for meta data storage.
type Meta interface {
	// ReadInt64 returns int64 starting from offset.
	ReadInt64(offset int) int64
	// WriteInt64 write int64 value to bytes starting from offset.
	WriteInt64(offset int, value int64)
	// Sync syncs mmaped bytes to storage.
	Sync() error
	// Close releases the underlying resources.
	Close() error
}

// meta implements Meta.
type meta struct {
	mappedPage page.MappedPage
}

// NewMeta returns a Meta by mapping file at filePath with size.
func NewMeta(filePath string, size int) (Meta, error) {
	bytes, err := fileutil.RWMap(filePath, size)
	if err != nil {
		return nil, err
	}

	return &meta{
		mappedPage: page.NewMappedPage(filePath, bytes, page.MMapCloseFunc, page.MMapSyncFunc),
	}, nil
}

// ReadInt64 returns int64 starting from offset.
func (m *meta) ReadInt64(offset int) int64 {
	reader := stream.BinaryReader(m.mappedPage.Data(offset, int64Size))
	return reader.ReadInt64()
}

// WriteInt64 write int64 value to bytes starting from offset.
func (m *meta) WriteInt64(offset int, value int64) {
	writer := stream.BinaryBufWriter(m.mappedPage.Data(offset, int64Size))
	writer.PutInt64(value)
}

// Sync syncs mmaped bytes to storage.
func (m *meta) Sync() error {
	return m.mappedPage.Sync()
}

// Close releases the underlying resources.
func (m *meta) Close() error {
	return m.mappedPage.Close()
}
