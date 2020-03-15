package table

import (
	"encoding/binary"
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./builder.go -destination=./builder_mock.go -package table

// for testing
var (
	newBufioWriterFunc = bufioutil.NewBufioWriter
)

// FileNumber represents sst file number
type FileNumber int64

// Int64 returns the int64 value of file number
func (i FileNumber) Int64() int64 {
	return int64(i)
}

// Builder represents sst file builder
type Builder interface {
	// FileNumber returns file name for store builder
	FileNumber() FileNumber
	// Add puts k/v pair init sst file write buffer
	// NOTICE: key must key in sort by desc
	Add(key uint32, value []byte) error
	// MinKey returns min key in store
	MinKey() uint32
	// MaxKey returns max key in store
	MaxKey() uint32
	// Size returns the length of store file
	Size() int32
	// Count returns the number of k/v pairs contained in the store
	Count() uint64
	// Abandon abandons current store build for some reason
	Abandon() error
	// Close closes sst file write buffer
	Close() error
}

// storeBuilder builds store file
type storeBuilder struct {
	fileNumber FileNumber
	fileName   string
	writer     bufioutil.BufioWriter
	offset     *encoding.FixedOffsetEncoder

	// see paper of roaring bitmap: https://arxiv.org/pdf/1603.06549.pdf
	keys   *roaring.Bitmap
	minKey uint32
	maxKey uint32

	first bool
}

// NewStoreBuilder creates store builder instance for building store file
func NewStoreBuilder(fileNumber FileNumber, fileName string) (Builder, error) {
	writer, err := newBufioWriterFunc(fileName)
	if err != nil {
		return nil, fmt.Errorf("create file write for store builder error:%s", err)
	}
	return &storeBuilder{
		fileNumber: fileNumber,
		fileName:   fileName,
		keys:       roaring.New(),
		writer:     writer,
		first:      true,
		offset:     encoding.NewFixedOffsetEncoder(),
	}, nil
}

// FileNumber returns file name of store builder.
func (b *storeBuilder) FileNumber() FileNumber {
	return b.fileNumber
}

// Add adds key/value pair into store file, if write failure return error
func (b *storeBuilder) Add(key uint32, value []byte) error {
	if !b.first && key <= b.maxKey {
		tableLogger.Warn("key is smaller then last key ignore current options.",
			logger.String("file", b.fileName),
			logger.Uint32("last", b.maxKey), logger.Uint32("cur", key))
		return nil
	}

	// get write offset
	offset := b.writer.Size()
	if _, err := b.writer.Write(value); err != nil {
		return fmt.Errorf("write data into store file error:%s", err)
	}
	// add offset into offset buffer
	b.offset.Add(uint32(offset))
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

// Abandon abandons current store build for some reason, for example compaction job fail or memory store dump error
func (b *storeBuilder) Abandon() error {
	return b.writer.Close()
}

// Close writes file footer before closing resources
func (b *storeBuilder) Close() error {
	if b.keys.IsEmpty() {
		return ErrEmptyKeys
	}
	posOfOffset := b.writer.Size()
	offset := b.offset.MarshalBinary()
	if _, err := b.writer.Write(offset); err != nil {
		return err
	}

	b.keys.RunOptimize()
	keys, err := encoding.BitmapMarshal(b.keys)
	if err != nil {
		return err
	}
	posOfKeys := b.writer.Size()
	if _, err = b.writer.Write(keys); err != nil {
		return err
	}

	// for file footer for offsets/keys index, length=1+4+4+8
	var buf [17]byte
	binary.LittleEndian.PutUint32(buf[:4], uint32(posOfOffset))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(posOfKeys))
	buf[8] = version0
	binary.LittleEndian.PutUint64(buf[9:], magicNumberOffsetFile)
	if _, err = b.writer.Write(buf[:]); err != nil {
		return err
	}
	return b.writer.Close()
}
