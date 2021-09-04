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
)

func TestFlusher_FlushInvertedIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().StreamWriter().Return(mockStreamWriter(ctrl), nil)
	indexFlusher, err := NewInvertedFlusher(mockFlusher)
	assert.Nil(t, err)
	assert.NotNil(t, indexFlusher)
	indexFlusher.PrepareTagKey(3)
	err = indexFlusher.FlushInvertedIndex(1, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	err = indexFlusher.FlushInvertedIndex(2, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	err = indexFlusher.FlushInvertedIndex(3, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	err = indexFlusher.FlushInvertedIndex(5, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	err = indexFlusher.FlushInvertedIndex(6, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)

	assert.NoError(t, indexFlusher.CommitTagKey())

	mockFlusher.EXPECT().Commit().Return(nil)
	err = indexFlusher.Close()
	assert.NoError(t, err)
}

func TestFlusher_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().StreamWriter().Return(nil, io.ErrUnexpectedEOF)

	indexFlusher, err := NewInvertedFlusher(mockFlusher)
	assert.Error(t, err)
	assert.Nil(t, indexFlusher)
}
