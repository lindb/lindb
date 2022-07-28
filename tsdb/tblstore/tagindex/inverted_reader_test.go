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

package tagindex

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
)

var bitmapUnmarshal = encoding.BitmapUnmarshal

func buildInvertedIndexBlock() (zoneBlock, ipBlock, hostBlock []byte) {
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher, _ := NewInvertedFlusher(nopKVFlusher)
	zoneMapping := map[uint32]*roaring.Bitmap{
		1: roaring.BitmapOf(1),
		2: roaring.BitmapOf(2),
		3: roaring.BitmapOf(3),
	}
	hostMapping := map[uint32]*roaring.Bitmap{
		1: roaring.BitmapOf(1),
		2: roaring.BitmapOf(2),
		3: roaring.BitmapOf(3),
		4: roaring.BitmapOf(4),
		5: roaring.BitmapOf(5),
		6: roaring.BitmapOf(6),
		7: roaring.BitmapOf(7),
		8: roaring.BitmapOf(8),
		9: roaring.BitmapOf(9),
	}
	flush := func(tagValueIDs []uint32, mapping map[uint32]*roaring.Bitmap) {
		for _, tagValueID := range tagValueIDs {
			_ = seriesFlusher.FlushInvertedIndex(tagValueID, mapping[tagValueID])
		}
	}
	/////////////////////////
	// flush zone tag, tagID: 20
	/////////////////////////
	seriesFlusher.PrepareTagKey(20)
	_ = seriesFlusher.FlushInvertedIndex(0, roaring.BitmapOf(1, 2, 3))
	flush([]uint32{1, 2, 3}, zoneMapping)
	// pick the zoneBlock buffer
	_ = seriesFlusher.CommitTagKey()
	zoneBlock = append(zoneBlock, nopKVFlusher.Bytes()...)

	/////////////////////////
	// flush ip tag, tagID: 21
	/////////////////////////
	// flush(ipMapping)
	seriesFlusher.PrepareTagKey(21)
	_ = seriesFlusher.FlushInvertedIndex(0, roaring.BitmapOf(1, 2, 3, 4000000, 5000000, 6000000, 7000000, 8000000, 9000000))
	_ = seriesFlusher.FlushInvertedIndex(1, roaring.BitmapOf(1))
	_ = seriesFlusher.FlushInvertedIndex(2, roaring.BitmapOf(2))
	_ = seriesFlusher.FlushInvertedIndex(3, roaring.BitmapOf(3))
	_ = seriesFlusher.FlushInvertedIndex(4000000, roaring.BitmapOf(4000000))
	_ = seriesFlusher.FlushInvertedIndex(5000000, roaring.BitmapOf(5000000))
	_ = seriesFlusher.FlushInvertedIndex(6000000, roaring.BitmapOf(6000000))
	_ = seriesFlusher.FlushInvertedIndex(7000000, roaring.BitmapOf(7000000))
	_ = seriesFlusher.FlushInvertedIndex(8000000, roaring.BitmapOf(8000000))
	_ = seriesFlusher.FlushInvertedIndex(9000000, roaring.BitmapOf(9000000))

	// pick the ipBlock buffer
	_ = seriesFlusher.CommitTagKey()
	ipBlock = append(ipBlock, nopKVFlusher.Bytes()...)

	/////////////////////////
	// flush host tag, tagID: 22
	/////////////////////////
	seriesFlusher.PrepareTagKey(22)
	_ = seriesFlusher.FlushInvertedIndex(0, roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9))
	flush([]uint32{1, 2, 3, 4, 5, 6, 7, 8, 9}, hostMapping)
	// pick the hostBlock buffer
	_ = seriesFlusher.CommitTagKey()
	hostBlock = append(hostBlock, nopKVFlusher.Bytes()...)
	return zoneBlock, ipBlock, hostBlock
}

func buildInvertedIndexReader(ctrl *gomock.Controller) InvertedReader {
	zoneBlock, ipBlock, hostBlock := buildInvertedIndexBlock()
	// mock readers
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(10)).Return(nil, nil).AnyTimes()
	mockReader.EXPECT().Get(uint32(19)).Return(nil, table.ErrKeyNotExist).AnyTimes()
	mockReader.EXPECT().Get(uint32(20)).Return(zoneBlock, nil).AnyTimes()
	mockReader.EXPECT().Get(uint32(21)).Return(ipBlock, nil).AnyTimes()
	mockReader.EXPECT().Get(uint32(22)).Return(hostBlock, nil).AnyTimes()
	// build series index inverterReader
	return NewInvertedReader([]table.Reader{mockReader})
}

func TestReader_FindSeriesIDsByTagValueIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := buildInvertedIndexReader(ctrl)
	// read not tag key id
	idSet, err := reader.GetSeriesIDsByTagValueIDs(19, roaring.BitmapOf(1))
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), idSet)
	// tag value ids is empty
	idSet, err = reader.GetSeriesIDsByTagValueIDs(19, nil)
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), idSet)
	// not found
	idSet, err = reader.GetSeriesIDsByTagValueIDs(10, roaring.BitmapOf(1))
	assert.Error(t, err)
	assert.Nil(t, idSet)

	// read zone block
	idSet, err = reader.GetSeriesIDsByTagValueIDs(21, roaring.BitmapOf(2, 49, 6000000, 6000033, 7000000))
	a := roaring.BitmapOf(2, 6000000, 7000000)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ToArray(), idSet.ToArray())

	idSet, err = reader.GetSeriesIDsByTagValueIDs(20, roaring.BitmapOf(1, 2))
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1, 2).ToArray(), idSet.ToArray())

	// unmarshal series ids err
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		d, _ := roaring.BitmapOf(1).ToBytes()
		if reflect.DeepEqual(d, data[:len(d)]) {
			// mock scan data err
			return fmt.Errorf("err")
		}
		// for other unmarshal
		return bitmap.UnmarshalBinary(data)
	}
	idSet, err = reader.GetSeriesIDsByTagValueIDs(20, roaring.BitmapOf(1, 2))
	assert.Error(t, err)
	assert.Nil(t, idSet)
}

func TestReader_InvertedIndex_reader_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.BitmapUnmarshal = bitmapUnmarshal
		ctrl.Finish()
	}()

	zoneBlock, _, _ := buildInvertedIndexBlock()
	reader, err := newTagInvertedReader(zoneBlock)
	assert.NoError(t, err)
	assert.NotNil(t, reader)

	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	// case 1: unmarshal series id err
	idSet, err := reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(1, 2))
	assert.Error(t, err)
	assert.Nil(t, idSet)
	// case 2: init inverted inverterReader err
	reader, err = newTagInvertedReader(zoneBlock)
	assert.Error(t, err)
	assert.Nil(t, reader)
	// case 3: validation offset err
	reader, err = newTagInvertedReader([]byte{
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
		4, 4, 4, 4,
		5})
	assert.Error(t, err)
	assert.Nil(t, reader)
}

func TestTagInvertedReader_scan(t *testing.T) {
	zoneBlock, ipBlock, _ := buildInvertedIndexBlock()
	reader, _ := newTagInvertedReader(ipBlock)
	scanner, err := newTagInvertedScanner(reader)
	assert.Nil(t, err)
	seriesIDs := roaring.New()
	// case 1: not match
	err = scanner.scan(10, 10, seriesIDs)
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), seriesIDs)
	// case 2: merge data
	scanner, _ = newTagInvertedScanner(reader)
	err = scanner.scan(0, 1, seriesIDs)
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1).ToArray(), seriesIDs.ToArray())
	// case 3: unmarshal series data err
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	scanner, _ = newTagInvertedScanner(reader)
	err = scanner.scan(0, 1, seriesIDs)
	assert.Error(t, err)
	// case 4: scanner is completed
	encoding.BitmapUnmarshal = bitmapUnmarshal
	reader, _ = newTagInvertedReader(zoneBlock)
	scanner, _ = newTagInvertedScanner(reader)
	seriesIDs.Clear()
	err = scanner.scan(10, 1, seriesIDs)
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.New().ToArray(), seriesIDs.ToArray())
}
