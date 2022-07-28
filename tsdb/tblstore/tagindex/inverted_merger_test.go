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
	"github.com/lindb/lindb/pkg/encoding"
)

func TestInvertedMerger_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nopFlusher := kv.NewNopFlusher()
	merge, _ := NewInvertedMerger(nopFlusher)
	merge.Init(nil)
	// case 1: merge data success
	err := merge.Merge(1, mockInvertedMergeData())
	assert.NoError(t, err)
	reader, err := newTagInvertedReader(nopFlusher.Bytes())
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8000000, 9000000).ToArray(), reader.keys.ToArray())
	seriesIDs, _ := reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(1))
	assert.EqualValues(t, roaring.BitmapOf(1, 10).ToArray(), seriesIDs.ToArray())
	seriesIDs, _ = reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(2))
	assert.EqualValues(t, roaring.BitmapOf(2).ToArray(), seriesIDs.ToArray())
	seriesIDs, _ = reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(8000000))
	assert.EqualValues(t, roaring.BitmapOf(8000000).ToArray(), seriesIDs.ToArray())
	// case 2: new reader err
	_ = nopFlusher.Commit()
	err = merge.Merge(2, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Len(t, nopFlusher.Bytes(), 0)
	// case 3: flush tag value data err
	flusher := NewMockInvertedFlusher(ctrl)
	m := merge.(*invertedMerger)
	m.invertedFlusher = flusher
	flusher.EXPECT().PrepareTagKey(gomock.Any()).Return()
	flusher.EXPECT().FlushInvertedIndex(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = merge.Merge(3, mockInvertedMergeData())
	assert.Error(t, err)
	assert.Len(t, nopFlusher.Bytes(), 0)
	// case 4: flush tag data err
	flusher.EXPECT().PrepareTagKey(gomock.Any()).Return().AnyTimes()
	flusher.EXPECT().FlushInvertedIndex(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	flusher.EXPECT().CommitTagKey().Return(fmt.Errorf("err"))
	err = merge.Merge(4, mockInvertedMergeData())
	assert.Error(t, err)
	assert.Len(t, nopFlusher.Bytes(), 0)
	// case 5: scan data err
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		d, _ := roaring.BitmapOf(1).ToBytes()
		if reflect.DeepEqual(d, data[:len(d)]) {
			// mock get series ids data err
			return fmt.Errorf("err")
		}
		// for other unmarshal
		return bitmap.UnmarshalBinary(data)
	}
	err = merge.Merge(5, mockInvertedMergeData())
	assert.Error(t, err)
	assert.Len(t, nopFlusher.Bytes(), 0)
}

func mockInvertedMergeData() (data [][]byte) {
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher, _ := NewInvertedFlusher(nopKVFlusher)
	zoneMapping := map[uint32]*roaring.Bitmap{
		1: roaring.BitmapOf(1),
		2: roaring.BitmapOf(2),
		3: roaring.BitmapOf(3),
	}
	hostMapping := map[uint32]*roaring.Bitmap{
		1:       roaring.BitmapOf(10),
		3:       roaring.BitmapOf(30),
		4:       roaring.BitmapOf(4),
		5:       roaring.BitmapOf(5),
		6:       roaring.BitmapOf(6),
		7:       roaring.BitmapOf(7),
		8000000: roaring.BitmapOf(8000000),
		9000000: roaring.BitmapOf(9000000),
	}
	flush := func(tagValueIDs []uint32, mapping map[uint32]*roaring.Bitmap) {
		for _, tagValueID := range tagValueIDs {
			_ = seriesFlusher.FlushInvertedIndex(tagValueID, mapping[tagValueID])
		}
	}
	seriesFlusher.PrepareTagKey(20)
	flush([]uint32{1, 2, 3}, zoneMapping)
	_ = seriesFlusher.CommitTagKey()
	data = append(data, append([]byte{}, nopKVFlusher.Bytes()...))
	seriesFlusher.PrepareTagKey(22)
	flush([]uint32{1, 3, 4, 5, 6, 7, 8000000, 9000000}, hostMapping)
	// pick the hostBlock buffer
	_ = seriesFlusher.CommitTagKey()
	data = append(data, append([]byte{}, nopKVFlusher.Bytes()...))
	return data
}

func TestInvertedMerger_Merge_same_tagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encoding.BitmapUnmarshal = bitmapUnmarshal
	nopFlusher := kv.NewNopFlusher()
	merge, _ := NewInvertedMerger(nopFlusher)
	merge.Init(nil)
	// case 1: merge data success
	err := merge.Merge(1, [][]byte{
		mockInvertedData(1, []uint32{1, 2, 3}, map[uint32]*roaring.Bitmap{
			1: roaring.BitmapOf(1),
			2: roaring.BitmapOf(2),
			3: roaring.BitmapOf(3),
		}),
		mockInvertedData(1, []uint32{4}, map[uint32]*roaring.Bitmap{
			4: roaring.BitmapOf(4),
		}),
	})
	assert.NoError(t, err)
	reader, err := newTagInvertedReader(append([]byte{}, nopFlusher.Bytes()...))
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1, 2, 3, 4).ToArray(), reader.keys.ToArray())
	seriesIDs, _ := reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(1))
	assert.EqualValues(t, roaring.BitmapOf(1).ToArray(), seriesIDs.ToArray())
	seriesIDs, _ = reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(2))
	assert.EqualValues(t, roaring.BitmapOf(2).ToArray(), seriesIDs.ToArray())
	seriesIDs, _ = reader.getSeriesIDsByTagValueIDs(roaring.BitmapOf(4))
	assert.EqualValues(t, roaring.BitmapOf(4).ToArray(), seriesIDs.ToArray())
	err = merge.Merge(1, [][]byte{
		nopFlusher.Bytes(),
		mockInvertedData(1, []uint32{40}, map[uint32]*roaring.Bitmap{
			40: roaring.BitmapOf(40),
		}),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, nopFlusher.Bytes())
}

func mockInvertedData(tagKeyID uint32, tagValueIDs []uint32, tagValues map[uint32]*roaring.Bitmap) (data []byte) {
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher, _ := NewInvertedFlusher(nopKVFlusher)
	flush := func(tagValueIDs []uint32, mapping map[uint32]*roaring.Bitmap) {
		for _, tagValueID := range tagValueIDs {
			_ = seriesFlusher.FlushInvertedIndex(tagValueID, mapping[tagValueID])
		}
	}
	seriesFlusher.PrepareTagKey(tagKeyID)
	flush(tagValueIDs, tagValues)
	_ = seriesFlusher.CommitTagKey()
	return nopKVFlusher.Bytes()
}
