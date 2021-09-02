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
	"hash/crc32"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
)

const (
	testKVPath = "test_builder"
)

var bitmapMarshal = encoding.BitmapMarshal

func TestFileNumber_Int64(t *testing.T) {
	assert.Equal(t, int64(10), FileNumber(10).Int64())
}

func TestStoreBuilder_magicNumber(t *testing.T) {
	code := []byte("eleme-ci")
	assert.Len(t, code, 8)
	assert.Equal(t, magicNumberOffsetFile, binary.LittleEndian.Uint64(code))
}

func TestStoreBuilder_BuildStore(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	var builder, err = NewStoreBuilder(10, testKVPath+"/000010.sst")
	defer func() {
		_ = os.RemoveAll(testKVPath)
		_ = builder.Close()
	}()

	assert.Nil(t, err)

	err = builder.Add(1, []byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), builder.Count())

	// reject for duplicate key
	err = builder.Add(1, []byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), builder.Count())

	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	assert.Equal(t, uint32(1), builder.MinKey())
	assert.Equal(t, uint32(10), builder.MaxKey())
	assert.Equal(t, FileNumber(10), builder.FileNumber())
	assert.True(t, builder.Size() > 0)
}

func TestStoreBuilder_Build_Err(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	ctrl := gomock.NewController(t)
	defer func() {
		newBufioWriterFunc = bufioutil.NewBufioStreamWriter
		encoding.BitmapMarshal = bitmapMarshal
		_ = os.Remove(testKVPath)
		ctrl.Finish()
	}()
	writer := bufioutil.NewMockBufioWriter(ctrl)
	newBufioWriterFunc = func(fileName string) (bufioutil.BufioWriter, error) {
		return writer, nil
	}
	builder, err := NewStoreBuilder(10, testKVPath+"/000200.sst")
	assert.NoError(t, err)
	writer.EXPECT().Size().Return(int64(10)).AnyTimes()

	// case 1: write value err
	writer.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
	err = builder.Add(10, []byte{1, 2, 3})
	assert.Error(t, err)
	// case 2: close empty keys
	err = builder.Close()
	assert.Equal(t, ErrEmptyKeys, err)
	// case 3: close write offset err
	writer.EXPECT().Write([]byte{1, 2, 3}).Return(10, nil)
	writer.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
	err = builder.Add(10, []byte{1, 2, 3})
	assert.NoError(t, err)
	err = builder.Close()
	assert.Error(t, err)
	// case 4: bitmap marshal err
	encoding.BitmapMarshal = func(bitmap *roaring.Bitmap) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	writer.EXPECT().Write(gomock.Any()).Return(10, nil)
	err = builder.Close()
	assert.Error(t, err)
	// case 5: write keys err
	encoding.BitmapMarshal = bitmapMarshal
	writer.EXPECT().Write(gomock.Any()).Return(10, nil)              // write offset
	writer.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err")) // write keys
	err = builder.Close()
	assert.Error(t, err)
	// case 6: write footer err
	writer.EXPECT().Write(gomock.Any()).Return(10, nil).MaxTimes(2)  // write offset/keys
	writer.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err")) // write footer
	err = builder.Close()
	assert.Error(t, err)
	// case 6: new builder err
	newBufioWriterFunc = func(fileName string) (bufioutil.BufioWriter, error) {
		return nil, fmt.Errorf("err")
	}
	builder, err = NewStoreBuilder(10, testKVPath+"/000200.sst")
	assert.Error(t, err)
	assert.Nil(t, builder)
}

func TestStoreBuilder_Abandon(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	defer func() {
		_ = os.RemoveAll(testKVPath)
	}()
	builder, err := NewStoreBuilder(10, testKVPath+"/000010.sst")
	assert.NoError(t, err)
	_ = builder.Add(1, []byte("test"))
	err = builder.Abandon()
	assert.NoError(t, err)
}

func Test_Builder_Stream_Writer(t *testing.T) {
	var builder, err = NewStoreBuilder(10, filepath.Join(t.TempDir(), "/000010.sst"))
	defer func() {
		_ = builder.Close()
	}()
	err = builder.Add(1, []byte("test"))
	assert.Nil(t, err)
	// bad key
	beforeSize := builder.Size()
	writer := builder.StreamWriter()
	assert.Zero(t, writer.Size())
	writer.Prepare(1)
	_, _ = writer.Write([]byte("aaa"))
	writer.Commit()
	assert.Equal(t, beforeSize, builder.Size())

	// normal increasing key
	writer.Prepare(2)
	beforeBatchSize := writer.Size()
	_, _ = writer.Write([]byte("aaa"))
	writer.Commit()
	// written len
	assert.Equal(t, builder.Size()-beforeSize, int32(3))
	assert.Equal(t, writer.Size()-beforeBatchSize, int32(3))
}

func Test_StreamWriter_CheckSum32(t *testing.T) {
	var builder, _ = NewStoreBuilder(10, filepath.Join(t.TempDir(), "/000011.sst"))
	defer func() {
		_ = builder.Close()
	}()
	writer := builder.StreamWriter()

	writer.Prepare(1)
	assert.Equal(t, uint32(0), writer.CRC32CheckSum())
	_, _ = writer.Write([]byte{1, 2, 3, 4, 5, 6})
	assert.Equal(t, uint32(2180413220), writer.CRC32CheckSum())
	assert.Equal(t, uint32(2180413220), writer.CRC32CheckSum())

	writer.Prepare(2)
	assert.Equal(t, uint32(0), writer.CRC32CheckSum())
	_, _ = writer.Write([]byte{1, 2})
	_, _ = writer.Write([]byte{3, 4})
	_, _ = writer.Write([]byte{5, 6})
	assert.Equal(t, uint32(2180413220), writer.CRC32CheckSum())
	assert.Equal(t, uint32(2180413220), writer.CRC32CheckSum())
}

func Benchmark_CRC32_1MB(b *testing.B) {
	hasher := crc32.New(crc32.IEEETable)
	buf := make([]byte, 1024)

	for i := 0; i < b.N; i++ {
		for y := 0; y < 1024; y++ {
			_, _ = hasher.Write(buf)
		}
	}
}
