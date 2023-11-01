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

package tagkeymeta

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
)

func TestFlusher_NewError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVFlusher := kv.NewMockFlusher(ctrl)
	mockKVFlusher.EXPECT().StreamWriter().Return(nil, io.ErrClosedPipe)
	flusher, err := NewFlusher(mockKVFlusher)
	assert.NotNil(t, err)
	assert.Nil(t, flusher)
}

func TestFlusher_Commit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVFlusher := kv.NewMockFlusher(ctrl)
	sw := table.NewMockStreamWriter(ctrl)
	mockKVFlusher.EXPECT().StreamWriter().Return(sw, nil).AnyTimes()

	flusher, err := NewFlusher(mockKVFlusher)
	assert.Nil(t, err)
	assert.NotNil(t, flusher)

	// mock Close error
	mockKVFlusher.EXPECT().Commit().Return(fmt.Errorf("commit error"))
	assert.NotNil(t, flusher.Close())

	// mock commit ok
	err = flusher.FlushTagKeyID(333, 100)
	assert.Nil(t, err)
}

func TestFlushTagKeyID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVFlusher := kv.NewMockFlusher(ctrl)
	sw := table.NewMockStreamWriter(ctrl)
	sw.EXPECT().Prepare(gomock.Any()).AnyTimes()
	sw.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()
	sw.EXPECT().Size().Return(uint32(10000)).AnyTimes()
	sw.EXPECT().CRC32CheckSum().Return(uint32(10000)).AnyTimes()
	sw.EXPECT().Commit().Return(nil)
	mockKVFlusher.EXPECT().StreamWriter().Return(sw, nil).AnyTimes()

	flusher, _ := NewFlusher(mockKVFlusher)

	// flush tagValue1
	flusher.EnsureSize((1 << 8) * (1 << 8))
	for x := 1; x < 1<<8; x++ {
		for y := 1; y < 1<<8; y++ {
			flusher.FlushTagValue([]byte(fmt.Sprintf("192.168.%d.%d", x, y)), uint32(x*y))
		}
	}
	// flush tagKeyID
	assert.Nil(t, flusher.FlushTagKeyID(1, 10))
}
