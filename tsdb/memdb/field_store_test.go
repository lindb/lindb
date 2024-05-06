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
	"go.uber.org/mock/gomock"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

var encodeFunc = encoding.NewTSDEncoder

func TestFieldStore_New(t *testing.T) {
	buf := make([]byte, pageSize)

	store := newFieldStore(buf)
	assert.NotNil(t, store)
	s := store.(*fieldStore)
	assert.Equal(t, uint16(0), s.getStart())
	assert.Equal(t, uint16(15), s.timeWindow())
}

func TestFieldStore_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buf := make([]byte, pageSize)
	store := newFieldStore(buf)
	assert.NotNil(t, store)
	s := store.(*fieldStore)

	capacity := store.Capacity()
	store.Write(field.SumField, 10, 10.1)
	// length not changed
	assert.Zero(t, store.Capacity()-capacity)
	// case 1: get write value
	value, ok := s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 10.1, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// case 2: get not exist value, out of time slot range
	value, ok = s.getCurrentValue(10, 12)
	assert.False(t, ok)
	assert.Equal(t, 0.0, value)
	value, ok = s.getCurrentValue(10, 0)
	assert.False(t, ok)
	assert.Equal(t, 0.0, value)
	// case 3: write exist value, need rollup
	capacity = store.Capacity()
	store.Write(field.SumField, 10, 10.1)
	assert.Zero(t, store.Capacity()-capacity)
	value, ok = s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 20.2, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// case 3: write new value
	capacity = store.Capacity()
	store.Write(field.SumField, 12, 12.1)
	assert.Zero(t, store.Capacity()-capacity)
	value, ok = s.getCurrentValue(10, 12)
	assert.True(t, ok)
	assert.InDelta(t, 12.1, value, 0)
	assert.Equal(t, uint16(2), s.getEnd())
	// case 4: get value in time slot range
	value, ok = s.getCurrentValue(10, 11)
	assert.False(t, ok)
	assert.Equal(t, 0.0, value)
	// case 5: test slot range [10,12]
	thisSlotRange := s.slotRange(s.getStart())
	assert.Equal(t, uint16(10), thisSlotRange.Start)
	assert.Equal(t, uint16(12), thisSlotRange.End)
	// case 6: compact for slot < start time, time range[5,12]
	capacity = store.Capacity()
	store.Write(field.SumField, 5, 5.3)
	assert.True(t, valueSize < store.Capacity()-capacity)
	thisSlotRange = s.slotRange(s.getStart())
	assert.Equal(t, uint16(5), thisSlotRange.Start)
	assert.Equal(t, uint16(12), thisSlotRange.End)
	value, ok = s.getCurrentValue(5, 5)
	assert.True(t, ok)
	assert.InDelta(t, 5.3, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// case 7: write old value
	capacity = store.Capacity()
	store.Write(field.SumField, 10, 10.1)
	assert.Zero(t, store.Capacity()-capacity)
	assert.Equal(t, uint16(5), s.getEnd())
	// case 8: compact for slot > end time, time range[5,12]
	capacity = store.Capacity()
	store.Write(field.SumField, 50, 50.1)
	assert.True(t, valueSize < store.Capacity()-capacity)
	thisSlotRange = s.slotRange(s.getStart())
	assert.Equal(t, uint16(5), thisSlotRange.Start)
	assert.Equal(t, uint16(50), thisSlotRange.End)
	value, ok = s.getCurrentValue(50, 50)
	assert.True(t, ok)
	assert.InDelta(t, 50.1, value, 0.0)
	assert.Equal(t, uint16(0), s.getEnd())
	// case 9: write 10 slot, compact old value
	capacity = store.Capacity()
	store.Write(field.SumField, 10, 10.1)
	assert.True(t, valueSize < store.Capacity()-capacity)
	assert.Equal(t, uint16(0), s.getEnd())
	value, ok = s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 10.1, value, 0)
}

func TestFieldStore_Write2(t *testing.T) {
	buf := make([]byte, pageSize)
	store := newFieldStore(buf)
	s := store.(*fieldStore)
	store.Write(field.SumField, 10, 178)
	capacity := s.Capacity()
	assert.NotZero(t, valueSize+headLen, capacity)
	value, ok := s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 178.0, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// write with old slot
	capacity = s.Capacity()
	store.Write(field.SumField, 10, 178)
	assert.Zero(t, store.Capacity()-capacity)
	value, ok = s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 178.0*2, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
}

func TestFieldStore_Write_Compact_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.TSDEncodeFunc = encodeFunc
		ctrl.Finish()
	}()

	buf := make([]byte, pageSize)
	store := newFieldStore(buf)
	assert.NotNil(t, store)
	s := store.(*fieldStore)

	store.Write(field.SumField, 10, 10.1)
	assert.Zero(t, store.Capacity())
	capacity := store.Capacity()
	store.Write(field.SumField, 100, 100.1)
	assert.Equal(t, 13, store.Capacity()-capacity)
	value, ok := s.getCurrentValue(100, 100)
	assert.True(t, ok)
	assert.InDelta(t, 100.1, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// compress data is nil
	assert.NotNil(t, s.compress)
}

func TestFieldStore_FlushFieldTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.TSDEncodeFunc = encodeFunc
		ctrl.Finish()
	}()
	nopKVFlusher := kv.NewNopFlusher()
	flusher, err := metricsdata.NewFlusher(nopKVFlusher)
	assert.NoError(t, err)
	fields := field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}}
	flusher.PrepareMetric(39, fields)

	slotRange := timeutil.SlotRange{Start: 5, End: 5}
	for idx, f := range fields {
		buf := make([]byte, pageSize)
		store := newFieldStore(buf)
		store.Write(field.SumField, 5, float64(f.ID))
		assert.NoError(t, store.FlushFieldTo(flusher,
			field.Meta{Type: field.SumField},
			&flushContext{SlotRange: slotRange, fieldIdx: idx}))
	}

	err = flusher.FlushSeries(10)
	assert.NoError(t, err)
	assert.NoError(t, flusher.CommitMetric(slotRange))
	data := nopKVFlusher.Bytes()
	r, err := metricsdata.NewReader("1.sst", data)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	seriesIDs := roaring.BitmapOf(10)
	found := 0
	highKeys := seriesIDs.GetHighKeys()
	for idx := range highKeys {
		highKey := highKeys[idx]
		lowSeriesIDs := seriesIDs.GetContainer(highKey)
		ctx := &flow.DataLoadContext{
			SeriesIDHighKey:       highKey,
			LowSeriesIDsContainer: lowSeriesIDs,
			ShardExecuteCtx: &flow.ShardExecuteContext{
				StorageExecuteCtx: &flow.StorageExecuteContext{
					Fields: fields,
					Query:  &stmt.Query{},
				},
			},
			DownSampling: func(slotRange timeutil.SlotRange, _ uint16, fieldIdx int, getter encoding.TSDValueGetter) {
				assert.Equal(t, timeutil.SlotRange{Start: 5, End: 5}, slotRange)
				for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
					value, ok := getter.GetValue(movingSourceSlot)
					if !ok {
						continue
					}
					assert.Equal(t, 5, int(movingSourceSlot))
					assert.Equal(t, value, float64(fields[fieldIdx].ID))
					found++
				}
			},
			Decoder: encoding.GetTSDDecoder(),
		}
		ctx.Grouping()
		loader := r.Load(ctx)
		loader.Load(ctx)
	}
	assert.Equal(t, 2, found)
}

func TestFieldStore_Load(t *testing.T) {
	buf := make([]byte, pageSize)
	f := field.Meta{ID: 1}
	store := newFieldStore(buf)
	store.Write(field.SumField, 5, float64(f.ID))

	ctx := &flow.DataLoadContext{
		DownSampling: func(slotRange timeutil.SlotRange, _ uint16, _ int, getter encoding.TSDValueGetter) {
			for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
				value, ok := getter.GetValue(movingSourceSlot)
				if !ok {
					continue
				}
				assert.Equal(t, 5, int(movingSourceSlot))
				assert.Equal(t, value, 1.0)
			}
		},
		Decoder: encoding.GetTSDDecoder(),
	}
	store.Load(ctx, 0, 0, field.SumField, timeutil.SlotRange{Start: 5, End: 10})
	store1 := store.(*fieldStore)
	store1.compact(field.SumField, 5)
	store.Load(ctx, 0, 0, field.SumField, timeutil.SlotRange{Start: 5, End: 10})
}
