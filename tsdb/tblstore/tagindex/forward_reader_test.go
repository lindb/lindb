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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestForwardReader_GetSeriesIDsForTagKeyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.BitmapUnmarshal = bitmapUnmarshal
		ctrl.Finish()
	}()

	reader := buildForwardReader(ctrl)
	// case 1: read not tagID key
	idSet, err := reader.GetSeriesIDsForTagKeyID(19)
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), idSet)
	// case 2: data is empty
	idSet, err = reader.GetSeriesIDsForTagKeyID(10)
	assert.Error(t, err)
	assert.Nil(t, idSet)
	// case 3: read series ids
	idSet, err = reader.GetSeriesIDsForTagKeyID(20)
	a := roaring.BitmapOf(1, 2, 3, 4, 65535+10, 65535+20, 65535+30, 65535+40)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ToArray(), idSet.ToArray())
	// case 4: unmarshal series ids err
	reader = buildForwardReader(ctrl)
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	idSet, err = reader.GetSeriesIDsForTagKeyID(20)
	assert.Error(t, err)
	assert.Nil(t, idSet)
}

func TestForwardReader_GetGroupingScanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.BitmapUnmarshal = bitmapUnmarshal
		ctrl.Finish()
	}()

	reader := buildForwardReader(ctrl)
	// case 1: read not tagID key
	scanners, err := reader.GetGroupingScanner(19, nil)
	assert.NoError(t, err)
	assert.Empty(t, scanners)
	// case 2: series ids not match
	scanners, err = reader.GetGroupingScanner(20, roaring.BitmapOf(100, 200))
	assert.NoError(t, err)
	assert.Empty(t, scanners)
	// case 3: series ids not match
	scanners, err = reader.GetGroupingScanner(20, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.Len(t, scanners, 1)
	// case 4: unmarshal series ids err
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	scanners, err = reader.GetGroupingScanner(20, roaring.BitmapOf(1, 2, 3))
	assert.Error(t, err)
	assert.Nil(t, scanners)
}

func TestForwardReader_offset_err(t *testing.T) {
	reader, err := NewTagForwardReader([]byte{
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
		4, 4, 4, 4,
		5})
	assert.Error(t, err)
	assert.Nil(t, reader)
}

func TestTagForwardReader_GetGroupingScanner(t *testing.T) {
	allSeriesIDs := roaring.BitmapOf(1, 2, 3, 4, 65535+10, 65535+20, 65535+30, 65535+40)
	block := buildForwardBlock()
	reader, err := NewTagForwardReader(block)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	// case 1: data not found
	seriesIDs, tagValueIDs := reader.GetSeriesAndTagValue(10)
	assert.Nil(t, seriesIDs)
	assert.Nil(t, tagValueIDs)
	// case 2: get container 0 data
	seriesIDs, tagValueIDs = reader.GetSeriesAndTagValue(0)
	assert.EqualValues(t, seriesIDs.ToArray(), allSeriesIDs.GetContainerAtIndex(0).ToArray())
	assert.Equal(t, []uint32{1, 2, 3, 4}, tagValueIDs)
	// case 3: get container 1 data
	seriesIDs, tagValueIDs = reader.GetSeriesAndTagValue(1)
	assert.EqualValues(t, seriesIDs.ToArray(), allSeriesIDs.GetContainerAtIndex(1).ToArray())
	assert.Equal(t, []uint32{10, 20, 30, 40}, tagValueIDs)
}

func TestTagForwardReader_scan(t *testing.T) {
	block := buildForwardBlock()
	reader, _ := NewTagForwardReader(block)
	scanner := newTagForwardScanner(reader)
	var tagValueIDs []uint32
	// case 1: not match
	tagValueIDs = scanner.scan(10, 10, tagValueIDs)
	assert.Len(t, tagValueIDs, 0)
	// case 2: merge data
	scanner = newTagForwardScanner(reader)
	tagValueIDs = scanner.scan(0, 1, tagValueIDs)
	tagValueIDs = scanner.scan(0, 2, tagValueIDs)
	tagValueIDs = scanner.scan(0, 3, tagValueIDs)
	tagValueIDs = scanner.scan(0, 4, tagValueIDs)
	tagValueIDs = scanner.scan(1, 9, tagValueIDs)
	assert.Equal(t, []uint32{1, 2, 3, 4, 10}, tagValueIDs)
	// case 3: scan completed
	tagValueIDs = scanner.scan(3, 9, tagValueIDs)
	assert.Equal(t, []uint32{1, 2, 3, 4, 10}, tagValueIDs)
}

func buildForwardReader(ctrl *gomock.Controller) ForwardReader {
	block := buildForwardBlock()
	// mock readers
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(10)).Return(nil, nil).AnyTimes()
	mockReader.EXPECT().Get(uint32(19)).Return(nil, table.ErrKeyNotExist).AnyTimes()
	mockReader.EXPECT().Get(uint32(20)).Return(block, nil).AnyTimes()
	// build series index inverterReader
	return NewForwardReader([]table.Reader{mockReader})
}

func buildForwardBlock() (block []byte) {
	nopKVFlusher := kv.NewNopFlusher()
	forwardFlusher, _ := NewForwardFlusher(nopKVFlusher)
	forwardFlusher.PrepareTagKey(10)
	_ = forwardFlusher.FlushForwardIndex([]uint32{1, 2, 3, 4})
	_ = forwardFlusher.FlushForwardIndex([]uint32{10, 20, 30, 40})
	_ = forwardFlusher.CommitTagKey(roaring.BitmapOf(1, 2, 3, 4, 65535+10, 65535+20, 65535+30, 65535+40))
	return nopKVFlusher.Bytes()
}
