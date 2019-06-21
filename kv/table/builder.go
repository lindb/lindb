package table

import (
	"encoding/binary"
	"fmt"
	"path/filepath"

	"github.com/eleme/lindb/kv/version"
	"github.com/eleme/lindb/pkg/bufioutil"
	"github.com/eleme/lindb/pkg/encoding"
	"github.com/eleme/lindb/pkg/logger"

	"github.com/RoaringBitmap/roaring"
	"go.uber.org/zap"
)

const (
	// magic-number in the footer of sst file
	magicNumberOffsetFile uint64 = 7308327815838786409
)

// Builder builds sst file
type Builder interface {
	// FileNumber returns file name for store builder
	FileNumber() int64
	// Add puts k/v pair init sst file write buffer
	Add(key uint32, value []byte) error
	// MinKey returns min key in store
	MinKey() uint32
	// MaxKey returns max key in store
	MaxKey() uint32
	// Size returns the length of store file
	Size() int32
	// Count returns the number of k/v pairs contained in the store
	Count() uint64
	// Close closes sst file write buffer
	Close() error
}

// storeBuilder builds store file
type storeBuilder struct {
	fileNumber int64
	fileName   string
	logger     *zap.Logger
	writer     bufioutil.BufioWriter
	offset     *encoding.DeltaBitPackingEncoder

	keys   *roaring.Bitmap
	minKey uint32
	maxKey uint32

	first bool
}

// NewStoreBuilder creates store builder instance for building store file
func NewStoreBuilder(path string, fileNumber int64) (Builder, error) {
	fileName := filepath.Join(path, version.Table(fileNumber))
	keys := roaring.New()
	log := logger.GetLogger()
	writer, err := bufioutil.NewBufioWriter(fileName)
	if err != nil {
		return nil, fmt.Errorf("create file write for store builder error:%s", err)
	}
	return &storeBuilder{
		fileNumber: fileNumber,
		fileName:   fileName,
		keys:       keys,
		logger:     log,
		writer:     writer,
		first:      true,
		offset:     encoding.NewDeltaBitPackingEncoder(),
	}, nil
}

// FileNumber returns file name of store builder.
func (b *storeBuilder) FileNumber() int64 {
	return b.fileNumber
}

// Add adds key/value pair into store file, if write failure return error
func (b *storeBuilder) Add(key uint32, value []byte) error {
	if !b.first && key <= b.maxKey {
		b.logger.Warn("key is smaller then last key ignore current options.",
			zap.String("file", b.fileName), zap.Any("last", b.maxKey), zap.Any("cur", key))
		return nil
	}

	// get write offset
	offset := b.writer.Size()
	if _, err := b.writer.Write(value); err != nil {
		return fmt.Errorf("write data into store file error:%s", err)
	}
	// add offset into offset buffer
	b.offset.Add(int32(offset))
	// add key into index block
	b.keys.Add(key)

	if b.first {
		b.minKey = key
	}

	b.maxKey = key
	b.first = false

	return nil
}

// MinKey returns min key in store
func (b *storeBuilder) MinKey() uint32 {
	return b.minKey
}

// MaxKey returns max key in store
func (b *storeBuilder) MaxKey() uint32 {
	return b.maxKey
}

// Size returns the length of store file
func (b *storeBuilder) Size() int32 {
	return int32(b.writer.Size())
}

// Count returns the number of k/v pairs contained in the store
func (b *storeBuilder) Count() uint64 {
	return b.keys.GetCardinality()
}

// Close writes file footer before closing resources
func (b *storeBuilder) Close() error {
	posOfOffset := b.writer.Size()
	offset, err := b.offset.Bytes()
	if err != nil {
		return fmt.Errorf("marshal store table offsets error:%s", err)
	}
	if _, err = b.writer.Write(offset); err != nil {
		return err
	}

	b.keys.RunOptimize()
	keys, err := b.keys.MarshalBinary()
	if err != nil {
		return err
	}
	posOfKeys := b.writer.Size()
	if _, err = b.writer.Write(keys); err != nil {
		return err
	}

	// for file footer for offsets/keys index, length=1+4+4+8
	var buf [16]byte
	binary.BigEndian.PutUint32(buf[:4], uint32(posOfOffset))
	binary.BigEndian.PutUint32(buf[4:8], uint32(posOfKeys))
	binary.BigEndian.PutUint64(buf[8:], magicNumberOffsetFile)
	if _, err = b.writer.Write(buf[:]); err != nil {
		return err
	}
	return b.writer.Close()
}
