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

package table

import (
	"encoding/binary"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"sync"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./builder.go -destination=./builder_mock.go -package table

// for testing
var (
	newBufioWriterFunc         = bufioutil.NewBufioStreamWriter
	_once4Builder              sync.Once
	_instanceBuilderStatistics *builderStatistics
)

func getBuilderStatistics() *builderStatistics {
	_once4Builder.Do(func() {
		tableBuilderScope := linmetric.NewScope("lindb.kv.table.builder")
		_instanceBuilderStatistics = &builderStatistics{
			AddBadKeys: tableBuilderScope.NewDeltaCounter("bad_keys"),
			AddKeys:    tableBuilderScope.NewDeltaCounter("add_keys"),
			AddBytes:   tableBuilderScope.NewDeltaCounter("add_bytes"),
		}
	})
	return _instanceBuilderStatistics
}

type builderStatistics struct {
	AddBadKeys *linmetric.BoundDeltaCounter
	AddKeys    *linmetric.BoundDeltaCounter
	AddBytes   *linmetric.BoundDeltaCounter
}

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
	// StreamWriter returns a writer for streaming writing data
	StreamWriter() StreamWriter
	// MinKey returns min key in store
	MinKey() uint32
	// MaxKey returns max key in store
	MaxKey() uint32
	// Size returns the length of store file
	Size() uint32
	// Count returns the number of k/v pairs contained in the store
	Count() uint64
	// Abandon abandons current store build for some reason
	Abandon() error
	// Close closes sst file write buffer
	Close() error
}

// StreamWriter writes multi buffer into the builder continuously
// Call Prepare, Write, Commit in order.
// sw.Prepare(1)
// sw.Write(...)
// sw.Write(...)
// sw.Commit()
type StreamWriter interface {
	// Prepare the writer with specified key
	// Resets the size and checksum
	Prepare(key uint32)
	// Writer writes buffer into the underlying file
	io.Writer
	// Size returns total written size of Write
	// Prepare will resets it to zero.
	Size() uint32
	// CRC32CheckSum returns a IEEE checksum of written bytes
	CRC32CheckSum() uint32
	// Commit marks the key/value pair has been written
	Commit() error
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
		offset:     encoding.NewFixedOffsetEncoder(true),
	}, nil
}

// FileNumber returns file name of store builder.
func (b *storeBuilder) FileNumber() FileNumber {
	return b.fileNumber
}

func (b *storeBuilder) ensureIncreasingKey(key uint32) bool {
	if b.first {
		return true
	}
	if key <= b.maxKey {
		getBuilderStatistics().AddBadKeys.Incr()
		tableLogger.Warn("key is smaller then last key ignore current options.",
			logger.String("file", b.fileName),
			logger.Uint32("last", b.maxKey),
			logger.Uint32("cur", key))
		return false
	}
	return true
}

func (b *storeBuilder) afterWrite(key uint32, offset int) {
	// add offset into offset buffer
	b.offset.Add(offset)
	// add key into index block
	b.keys.Add(key)
	if b.first {
		b.minKey = key
	}
	b.maxKey = key
	b.first = false
}

// Add adds key/value pair into store file, if write failure return error
func (b *storeBuilder) Add(key uint32, value []byte) error {
	if !b.ensureIncreasingKey(key) {
		return nil
	}

	// get write offset
	offset := b.writer.Size()
	if _, err := b.writer.Write(value); err != nil {
		return fmt.Errorf("write data into store file error:%s", err)
	}
	getBuilderStatistics().AddKeys.Incr()
	getBuilderStatistics().AddBytes.Add(float64(len(value)))
	b.afterWrite(key, int(offset))
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
func (b *storeBuilder) Size() uint32 {
	return uint32(b.writer.Size())
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

func (b *storeBuilder) StreamWriter() StreamWriter {
	return newStreamWriter(b)
}

func newStreamWriter(builder *storeBuilder) *streamWriter {
	return &streamWriter{
		builder: builder,
		badKey:  true,
		crc32:   crc32.New(crc32.IEEETable),
	}
}

type streamWriter struct {
	builder *storeBuilder
	size    uint32
	key     uint32
	offset  int64
	badKey  bool
	crc32   hash.Hash32
}

func (sw *streamWriter) Prepare(key uint32) {
	sw.badKey = !sw.builder.ensureIncreasingKey(key)
	sw.offset = sw.builder.writer.Size()
	sw.key = key
	sw.size = 0
	sw.crc32.Reset()
}

func (sw *streamWriter) Write(data []byte) (int, error) {
	if sw.badKey {
		return 0, nil
	}
	n, err := sw.builder.writer.Write(data)
	_, _ = sw.crc32.Write(data)
	if err == nil {
		sw.size += uint32(n)
	}
	getBuilderStatistics().AddBytes.Add(float64(len(data)))
	return n, err
}

func (sw *streamWriter) Size() uint32 {
	return sw.size
}

func (sw *streamWriter) CRC32CheckSum() uint32 {
	return sw.crc32.Sum32()
}

func (sw *streamWriter) Commit() error {
	if sw.badKey {
		return nil
	}
	sw.builder.afterWrite(sw.key, int(sw.offset))
	// preventing committing twice
	sw.badKey = true
	return nil
}
