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
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
)

var bitmapUnmarshal = encoding.BitmapUnmarshal

func TestReader_Fail(t *testing.T) {
	defer func() {
		mapFunc = fileutil.Map
		unmapFunc = fileutil.Unmap
		openFileFn = os.Open
		_ = os.RemoveAll(testKVPath)
	}()
	// open file fail
	openFileFn = func(name string) (*os.File, error) {
		return nil, fmt.Errorf("err")
	}
	reader, err := newMMapStoreReader(testKVPath, "000010.sst")
	assert.Error(t, err)
	assert.Nil(t, reader)
	// case 1: map err
	mapFunc = func(path *os.File) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	reader, err = newMMapStoreReader(testKVPath, "000010.sst")
	assert.Error(t, err)
	assert.Nil(t, reader)
	// case 2: footer length err
	mapFunc = func(path *os.File) (bytes []byte, err error) {
		return []byte{1, 2, 3}, nil
	}
	unmapFunc = func(_ *os.File, _ []byte) error {
		return fmt.Errorf("err")
	}
	reader, err = newMMapStoreReader(testKVPath, "000010.sst")
	assert.Error(t, err)
	assert.Nil(t, reader)
	// case 3: init err
	mapFunc = func(path *os.File) (bytes []byte, err error) {
		return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5}, nil
	}
	reader, err = newMMapStoreReader(testKVPath, "000010.sst")
	assert.Error(t, err)
	assert.Nil(t, reader)
}

func TestStoreMMapReader_readBytes_Err(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	defer func() {
		uint64Func = binary.LittleEndian.Uint64
		intsAreSortedFunc = sort.IntsAreSorted
		encoding.BitmapUnmarshal = bitmapUnmarshal
		unmarshalFixedOffsetFunc = unmarshalFixedOffset
		assert.NoError(t, os.RemoveAll(testKVPath))
	}()
	builder, err := NewStoreBuilder(10, filepath.Join(testKVPath, "000010.sst"))
	assert.NoError(t, err)

	_ = builder.Add(1, []byte("test"))
	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	err = builder.Close()
	assert.Nil(t, err)

	// case1, ok
	r, err := newMMapStoreReader(filepath.Join(testKVPath, "000010.sst"), "000010.sst")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "000010.sst", r.FileName())

	block, err := r.(*storeMMapReader).getBlock(0)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(block))

	block, err = r.(*storeMMapReader).getBlock(1)
	assert.NoError(t, err)
	assert.Equal(t, "test10", string(block))

	block, err = r.(*storeMMapReader).getBlock(2)
	assert.Error(t, err)
	assert.Len(t, block, 0)
	assert.NoError(t, r.Close())
	// case 2: offset not in sort
	uint64Func = binary.LittleEndian.Uint64
	intsAreSortedFunc = func(x []int) bool {
		return false
	}
	r, err = newMMapStoreReader(filepath.Join(testKVPath, "000010.sst"), "000010.sst")
	assert.Error(t, err)
	assert.Nil(t, r)
	intsAreSortedFunc = sort.IntsAreSorted

	// case 3: unmarshal fixed offset failure
	unmarshalFixedOffsetFunc = func(decoder *encoding.FixedOffsetDecoder, data []byte) error {
		return fmt.Errorf("err")
	}
	r, err = newMMapStoreReader(testKVPath, "000010.sst")
	assert.Error(t, err)
	assert.Nil(t, r)

	unmarshalFixedOffsetFunc = unmarshalFixedOffset

	// case 4: unmarshal keys err
	uint64Func = binary.LittleEndian.Uint64
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	r, err = newMMapStoreReader(filepath.Join(testKVPath, "000010.sst"), "000010.sst")
	assert.Error(t, err)
	assert.Nil(t, r)

	// case 5: offset's size != key's size
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		bitmap.AddRange(1, 1000)
		return nil
	}
	r, err = newMMapStoreReader(filepath.Join(testKVPath, "000010.sst"), "000010.sst")
	assert.Error(t, err)
	assert.Nil(t, r)
}

func TestReader(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	defer func() {
		_ = os.RemoveAll(testKVPath)
	}()

	builder, err := NewStoreBuilder(10, filepath.Join(testKVPath, "000010.sst"))
	assert.NoError(t, err)

	_ = builder.Add(1, []byte("test"))
	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	err = builder.Close()
	assert.Nil(t, err)

	cache := NewCache(testKVPath, time.Hour)

	reader, err := cache.GetReader("", "000010.sst")
	assert.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()
	assert.Equal(t, filepath.Join(testKVPath, "000010.sst"), reader.Path())

	// get from store cache
	reader, err = cache.GetReader("", "000010.sst")
	assert.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()
	value, err := reader.Get(100)
	assert.Error(t, err)
	assert.Nil(t, value)

	value, _ = reader.Get(1)
	assert.Equal(t, []byte("test"), value)
	value, _ = reader.Get(10)
	assert.Equal(t, []byte("test10"), value)
	cache.Evict("000100.sst")
	_ = reader.Close()
	cache.Evict("000010.sst")
	_ = cache.Close()
}

func TestStoreIterator(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	defer func() {
		_ = os.RemoveAll(testKVPath)
	}()
	builder, err := NewStoreBuilder(10, filepath.Join(testKVPath, "000010.sst"))
	assert.NoError(t, err)

	_ = builder.Add(1, []byte("test"))
	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	err = builder.Close()
	assert.Nil(t, err)

	cache := NewCache(testKVPath, time.Hour)
	reader, err := cache.GetReader("", "000010.sst")
	assert.NoError(t, err)

	defer func() {
		_ = reader.Close()
	}()
	it := reader.Iterator()
	assert.True(t, it.HasNext())
	assert.Equal(t, uint32(1), it.Key())
	assert.Equal(t, []byte("test"), it.Value())

	assert.True(t, it.HasNext())
	assert.Equal(t, uint32(10), it.Key())
	assert.Equal(t, []byte("test10"), it.Value())

	assert.False(t, it.HasNext())
}
