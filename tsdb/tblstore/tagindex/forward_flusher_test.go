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
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
)

func mockStreamWriter(ctrl *gomock.Controller) table.StreamWriter {
	sw := table.NewMockStreamWriter(ctrl)
	sw.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()
	sw.EXPECT().Size().Return(uint32(1000)).AnyTimes()
	sw.EXPECT().Prepare(gomock.Any()).AnyTimes()
	sw.EXPECT().CRC32CheckSum().Return(uint32(1)).AnyTimes()
	sw.EXPECT().Commit().Return(nil).AnyTimes()
	return sw
}

func TestForwardFlusher_Flusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	sw := mockStreamWriter(ctrl)
	mockFlusher.EXPECT().StreamWriter().Return(sw, nil)

	indexFlusher, err := NewForwardFlusher(mockFlusher)
	assert.Nil(t, err)
	assert.NotNil(t, indexFlusher)
	assert.Nil(t, indexFlusher.FlushForwardIndex([]uint32{1, 2, 3, 4}))
	assert.Nil(t, indexFlusher.FlushForwardIndex([]uint32{1, 2, 3, 4}))

	indexFlusher.PrepareTagKey(3)
	err = indexFlusher.CommitTagKey(roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	mockFlusher.EXPECT().Commit().Return(nil)
	err = indexFlusher.Close()
	assert.NoError(t, err)
}

func TestForwardFlusher_Flush_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	kvFlusher := kv.NewMockFlusher(ctrl)
	kvFlusher.EXPECT().StreamWriter().Return(nil, io.ErrClosedPipe)

	indexFlusher, err := NewForwardFlusher(kvFlusher)
	assert.Error(t, err)
	assert.Nil(t, indexFlusher)

	mockFlusher := kv.NewMockFlusher(ctrl)
	sw := mockStreamWriter(ctrl)
	mockFlusher.EXPECT().StreamWriter().Return(sw, nil)

	indexFlusher, _ = NewForwardFlusher(mockFlusher)
	indexFlusher.PrepareTagKey(3)
	err = indexFlusher.CommitTagKey(roaring.BitmapOf(1, 2, 3))
	assert.Nil(t, err)
}
