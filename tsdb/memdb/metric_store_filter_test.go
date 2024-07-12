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

package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"
)

func TestFieldEntry(t *testing.T) {
	f := &fieldEntry{}
	f.Reset(nil)

	// write buffer not set
	v, ok := f.GetValue(10)
	assert.Zero(t, v)
	assert.False(t, ok)
	buf := make([]byte, pageSize)
	f.Reset(buf)
	// no data
	v, ok = f.GetValue(10)
	assert.Zero(t, v)
	assert.False(t, ok)

	// no compress store
	b := f.getCompressBuf(0)
	assert.Nil(t, b)

	// no compress data
	f.compressBuf = NewCompressStore()
	b = f.getCompressBuf(0)
	assert.Nil(t, b)
}

func TestMemoryDatabase_filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaDB := NewMockMetadataDatabase(ctrl)
	indexDB := NewMockIndexDatabase(ctrl)
	indexDB.EXPECT().GetMetadataDatabase().Return(metaDB).AnyTimes()
	md := &memoryDatabase{
		indexDB: indexDB,
	}
	t.Run("memory metric not found", func(t *testing.T) {
		metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(nil, false)
		rs, err := md.filter(nil, 100, nil, nil)
		assert.NoError(t, err)
		assert.Nil(t, rs)
	})

	t.Run("field not found", func(t *testing.T) {
		ms := newMetricStore()
		metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(ms, true)
		rs, err := md.filter(&flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{Fields: field.Metas{{Name: "test"}}},
		}, 100, nil, nil)
		assert.Error(t, err)
		assert.Nil(t, rs)
	})
	t.Run("field data not found", func(t *testing.T) {
		ms := newMetricStore()
		_, _ = ms.GenField("test", field.SumField)
		metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(ms, true)
		rs, err := md.filter(&flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{Fields: field.Metas{{Name: "test"}}},
		}, 100, nil, nil)
		assert.Error(t, err)
		assert.Nil(t, rs)
	})
}
