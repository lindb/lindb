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

package invertedindex

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestForwardFlusher_Flusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewForwardFlusher(mockFlusher)
	assert.NotNil(t, indexFlusher)
	indexFlusher.FlushForwardIndex([]uint32{1, 2, 3, 4})
	indexFlusher.FlushForwardIndex([]uint32{1, 2, 3, 4})
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	err := indexFlusher.FlushTagKeyID(3, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	mockFlusher.EXPECT().Commit().Return(nil)
	err = indexFlusher.Commit()
	assert.NoError(t, err)
}

func TestForwardFlusher_Flush_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.BitmapMarshal = bitMapMarshal
		ctrl.Finish()
	}()

	indexFlusher := NewForwardFlusher(nil)
	assert.NotNil(t, indexFlusher)
	encoding.BitmapMarshal = func(bitmap *roaring.Bitmap) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	err := indexFlusher.FlushTagKeyID(3, roaring.BitmapOf(1, 2, 3))
	assert.Error(t, err)
}
