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
	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
)

func TestForwardMerger_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	nopFlusher1 := kv.NewNopFlusher()
	merge, _ := NewForwardMerger(nopFlusher1)
	merge.Init(nil)
	// case 1: merge data success
	err := merge.Merge(1, mockMergeForwardBlock())
	assert.NoError(t, err)
	reader, err := NewTagForwardReader(nopFlusher1.Bytes())
	assert.NoError(t, err)
	assert.EqualValues(t,
		roaring.BitmapOf(1, 2, 3, 4, 65535+10, 65535+20, 65535+30, 65535+40).ToArray(),
		reader.GetSeriesIDs().ToArray())
	_, tagValueIDs := reader.GetSeriesAndTagValue(0)
	assert.Equal(t, []uint32{1, 2, 3, 4}, tagValueIDs)
	_, tagValueIDs = reader.GetSeriesAndTagValue(1)
	assert.Equal(t, []uint32{10, 20, 30, 40}, tagValueIDs)
	// case 2: new reader err
	nopFlusher2 := kv.NewNopFlusher()
	merge, _ = NewForwardMerger(nopFlusher2)
	err = merge.Merge(1, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Nil(t, nopFlusher2.Bytes())
	// case 3: flush tag key data err
	flusher := NewMockForwardFlusher(ctrl)
	m := merge.(*forwardMerger)
	m.forwardFlusher = flusher
	flusher.EXPECT().FlushForwardIndex(gomock.Any()).AnyTimes()
	flusher.EXPECT().PrepareTagKey(gomock.Any()).AnyTimes()
	flusher.EXPECT().CommitTagKey(gomock.Any()).Return(fmt.Errorf("err"))
	err = merge.Merge(1, mockMergeForwardBlock())
	assert.Error(t, err)
	assert.Nil(t, nopFlusher2.Bytes())
}

func mockMergeForwardBlock() (block [][]byte) {
	nopKVFlusher1 := kv.NewNopFlusher()
	forwardFlusher, _ := NewForwardFlusher(nopKVFlusher1)
	forwardFlusher.PrepareTagKey(10)
	_ = forwardFlusher.FlushForwardIndex([]uint32{1, 3})
	_ = forwardFlusher.FlushForwardIndex([]uint32{10, 20})
	_ = forwardFlusher.CommitTagKey(roaring.BitmapOf(1, 3, 65535+10, 65535+20))
	block = append(block, nopKVFlusher1.Bytes())

	// create new nop flusher, because under nop flusher share buffer
	nopKVFlusher2 := kv.NewNopFlusher()
	forwardFlusher, _ = NewForwardFlusher(nopKVFlusher2)
	forwardFlusher.PrepareTagKey(10)
	_ = forwardFlusher.FlushForwardIndex([]uint32{2, 4})
	_ = forwardFlusher.FlushForwardIndex([]uint32{30, 40})
	_ = forwardFlusher.CommitTagKey(roaring.BitmapOf(2, 4, 65535+30, 65535+40))
	block = append(block, nopKVFlusher2.Bytes())
	return
}
