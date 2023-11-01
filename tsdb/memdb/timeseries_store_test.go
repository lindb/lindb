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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestTimeSeriesStore_GetOrCreateFStore(t *testing.T) {
	tStore := newTimeSeriesStore()
	f, ok := tStore.GetFStore(10)
	assert.Nil(t, f)
	assert.False(t, ok)
	tStore.InsertFStore(newFieldStore(make([]byte, pageSize), 10))
	// get field store
	f, ok = tStore.GetFStore(10)
	assert.NotNil(t, f)
	assert.True(t, ok)
	// field store not exist
	f, ok = tStore.GetFStore(100)
	assert.Nil(t, f)
	assert.False(t, ok)
	for i := 1; i < 100; i++ {
		tStore.InsertFStore(newFieldStore(make([]byte, pageSize), field.ID(10*i)))
		tStore.InsertFStore(newFieldStore(make([]byte, pageSize), 10))
		f, ok = tStore.GetFStore(10)
		assert.NotNil(t, f)
		assert.True(t, ok)
	}
	assert.True(t, tStore.Capacity() > 0)

	f, ok = tStore.GetFStore(101)
	assert.False(t, ok)
	assert.Nil(t, f)
}

func TestTimeSeriesStore_FlushSeriesTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := metricsdata.NewMockFlusher(ctrl)
	tStore := newTimeSeriesStore()
	// no field store
	err := tStore.FlushFieldsTo(nil, nil)
	assert.NoError(t, err)

	s := tStore.(*timeSeriesStore)
	fStore := NewMockfStoreINTF(ctrl)
	s.InsertFStore(fStore)

	// case : flush data
	gomock.InOrder(
		flusher.EXPECT().GetFieldMetas().Return(field.Metas{{ID: 1}, {ID: 2}, {ID: 3}}),
		fStore.EXPECT().GetFieldID().Return(field.ID(2)),
		flusher.EXPECT().FlushField(nil),
		fStore.EXPECT().GetFieldID().Return(field.ID(2)),
		fStore.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any(), gomock.Any()),
		flusher.EXPECT().FlushField(nil),
	)
	assert.NoError(t, tStore.FlushFieldsTo(flusher, &flushContext{}))

	gomock.InOrder(
		flusher.EXPECT().GetFieldMetas().Return(field.Metas{{ID: 1}, {ID: 2}, {ID: 3}}),
		fStore.EXPECT().GetFieldID().Return(field.ID(2)),
		flusher.EXPECT().FlushField(nil),
		fStore.EXPECT().GetFieldID().Return(field.ID(2)),
		fStore.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err")),
	)
	assert.Error(t, tStore.FlushFieldsTo(flusher, &flushContext{}))
}

func TestTimeSeriesStore_load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore()
	tStore := tStoreInterface.(*timeSeriesStore)
	fs := make(map[field.ID]*MockfStoreINTF)
	for i := 3; i < 10; i++ {
		fStore := NewMockfStoreINTF(ctrl)
		id := field.ID(i)
		fStore.EXPECT().GetFieldID().Return(id).AnyTimes()
		tStore.InsertFStore(fStore)
		fs[id] = fStore
	}

	ctx := &flow.DataLoadContext{}
	// case 2: field id not match
	tStore.load(ctx, 1, field.Metas{{
		ID:   200,
		Type: field.SumField,
	}}, timeutil.SlotRange{})
	// case 3: field id not match
	tStore.load(ctx, 1, field.Metas{{
		ID:   1,
		Type: field.SumField,
	}}, timeutil.SlotRange{})
	// case 4: match one field
	f5 := fs[field.ID(5)]
	f5.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	tStore.load(ctx, 1, field.Metas{{
		ID:   5,
		Type: field.SumField,
	}}, timeutil.SlotRange{})
	// case 4: match two fields
	f5.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	f8 := fs[field.ID(8)]
	f8.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	tStore.load(ctx, 0, field.Metas{{
		ID:   5,
		Type: field.SumField,
	}, {
		ID:   8,
		Type: field.SumField,
	}}, timeutil.SlotRange{})
}
