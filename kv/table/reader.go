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
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package table

// for testing
var (
	ErrKeyNotExist            = errors.New("key not exist in kv table")
	mapFunc                   = fileutil.Map
	unmapFunc                 = fileutil.Unmap
	uint64Func                = binary.LittleEndian.Uint64
	_once4ReaderStatistics    sync.Once
	_instanceReaderStatistics *readerStatistics
)

func getReaderStatistics() *readerStatistics {
	_once4ReaderStatistics.Do(func() {
		tableReaderScope := linmetric.NewScope("lindb.kv.table.reader")
		_instanceReaderStatistics = &readerStatistics{
			getErrors:    tableReaderScope.NewCounter("get_errors"),
			getCounts:    tableReaderScope.NewCounter("get_counts"),
			getBytes:     tableReaderScope.NewCounter("get_bytes"),
			mmapCounts:   tableReaderScope.NewCounter("mmap_counts"),
			mmapErros:    tableReaderScope.NewCounter("mmap_errors"),
			unmmapCounts: tableReaderScope.NewCounter("unmmap_counts"),
			unmmapErrors: tableReaderScope.NewCounter("unmmap_errors"),
		}
	})
	return _instanceReaderStatistics
}

type readerStatistics struct {
	getErrors    *linmetric.BoundDeltaCounter
	getCounts    *linmetric.BoundDeltaCounter
	getBytes     *linmetric.BoundDeltaCounter
	mmapCounts   *linmetric.BoundDeltaCounter
	mmapErros    *linmetric.BoundDeltaCounter
	unmmapCounts *linmetric.BoundDeltaCounter
	unmmapErrors *linmetric.BoundDeltaCounter
}

// Reader represents reader which reads k/v pair from store file
type Reader interface {
	// Path returns the file path
	Path() string
	// Get returns value for giving key,
	// if key not exist, return nil, ErrKeyNotExist
	Get(key uint32) ([]byte, error)
	// Iterator iterates over a store's key/value pairs in key order.
	Iterator() Iterator
	// Close closes reader, release related resources
	Close() error
}

// storeMMapReader represents mmap store file reader
type storeMMapReader struct {
	path         string                       // path of sst-file
	fullBlock    []byte                       // mmaped file content
	entriesBlock []byte                       // mmaped file content without footer
	keys         *roaring.Bitmap              // bitmap of keys
	offsets      *encoding.FixedOffsetDecoder // offset of values
}

// newMMapStoreReader creates mmap store file reader
func newMMapStoreReader(path string) (r Reader, err error) {
	data, err := mapFunc(path)
	defer func() {
		if err != nil && len(data) > 0 {
			// if init err and map data exist, need unmap it
			if e := unmapFunc(data); e != nil {
				getReaderStatistics().unmmapErrors.Incr()
				tableLogger.Warn("unmap error when new store reader fail",
					logger.String("path", path), logger.Error(err))
			} else {
				getReaderStatistics().unmmapCounts.Incr()
			}
		}
	}()
	if err != nil {
		getReaderStatistics().mmapErros.Incr()
		return
	}
	getReaderStatistics().mmapCounts.Incr()

	if len(data) < sstFileFooterSize {
		err = fmt.Errorf("length of sstfile:%s length is too short", path)
		return
	}
	reader := &storeMMapReader{
		path:      path,
		fullBlock: data,
		keys:      roaring.New(),
	}

	if err := reader.initialize(); err != nil {
		return nil, err
	}

	return reader, nil
}

// initialize initializes store reader, reads index block(keys,offset etc.), then caches it
func (r *storeMMapReader) initialize() error {
	// decode footer
	footerStart := len(r.fullBlock) - sstFileFooterSize
	// validate magic-number
	if uint64Func(r.fullBlock[footerStart+magicNumberAtFooter:]) != magicNumberOffsetFile {
		return fmt.Errorf("verify magic-number of sstfile:%s failure", r.path)
	}
	posOfOffset := int(binary.LittleEndian.Uint32(r.fullBlock[footerStart : footerStart+4]))
	posOfKeys := int(binary.LittleEndian.Uint32(r.fullBlock[footerStart+4 : footerStart+8]))
	if !sort.IntsAreSorted([]int{
		0, posOfOffset, posOfKeys, footerStart}) {
		return fmt.Errorf("bad footer data, posOfOffsets: %d posOfKeys: %d,"+
			" footerStart: %d", posOfOffset, posOfKeys, footerStart)
	}
	// decode offsets
	offsetsBlock := r.fullBlock[posOfOffset:posOfKeys]
	r.offsets = encoding.NewFixedOffsetDecoder()
	if _, err := r.offsets.Unmarshal(offsetsBlock); err != nil {
		return fmt.Errorf("unmarshal fixed-offsets decoder with error: %s", err)
	}
	// decode keys
	if err := encoding.BitmapUnmarshal(r.keys, r.fullBlock[posOfKeys:]); err != nil {
		return fmt.Errorf("unmarshal keys data from file[%s] error:%s", r.path, err)
	}
	// validate keys and offsets
	if r.offsets.Size() != int(r.keys.GetCardinality()) {
		return fmt.Errorf("num. of keys != num. of offsets in file[%s]", r.path)
	}
	// read entries block
	r.entriesBlock = r.fullBlock[:posOfOffset]
	return nil
}

// Path returns the file path
func (r *storeMMapReader) Path() string {
	return r.path
}

// Get return value for key, if not exist return nil,false
func (r *storeMMapReader) Get(key uint32) ([]byte, error) {
	if !r.keys.Contains(key) {
		return nil, ErrKeyNotExist
	}
	// bitmap data's index from 1, so idx= get index - 1
	idx := r.keys.Rank(key)
	return r.getBlock(int(idx) - 1)
}

func (r *storeMMapReader) getBlock(idx int) ([]byte, error) {
	block, err := r.offsets.GetBlock(idx, r.entriesBlock)
	if err == nil {
		getReaderStatistics().getCounts.Incr()
		getReaderStatistics().getBytes.Add(float64(len(block)))
	} else {
		getReaderStatistics().getErrors.Incr()
	}
	return block, err
}

// Iterator iterates over a store's key/value pairs in key order.
func (r *storeMMapReader) Iterator() Iterator {
	return newMMapIterator(r)
}

// Close store reader, release resource
func (r *storeMMapReader) Close() error {
	r.entriesBlock = nil
	err := fileutil.Unmap(r.fullBlock)
	if err == nil {
		getReaderStatistics().unmmapCounts.Incr()
	} else {
		getReaderStatistics().unmmapErrors.Incr()
	}
	return err
}

// storeMMapIterator iterates k/v pair using mmap store reader
type storeMMapIterator struct {
	reader *storeMMapReader
	keyIt  roaring.IntIterable

	idx int
}

// newMMapIterator creates store iterator using mmap store reader
func newMMapIterator(reader *storeMMapReader) Iterator {
	return &storeMMapIterator{
		reader: reader,
		keyIt:  reader.keys.Iterator(),
	}
}

// HasNext returns if the iteration has more element.
// It returns false if the iterator is exhausted.
func (it *storeMMapIterator) HasNext() bool {
	return it.keyIt.HasNext()
}

// Key returns the key of the current key/value pair
func (it *storeMMapIterator) Key() uint32 {
	key := it.keyIt.Next()
	return key
}

// Value returns the value of the current key/value pair
func (it *storeMMapIterator) Value() []byte {
	block, _ := it.reader.getBlock(it.idx)
	it.idx++
	return block
}
