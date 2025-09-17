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

// import (
// 	"fmt"
// 	"testing"
//
// 	"github.com/lindb/roaring"
// 	"github.com/stretchr/testify/assert"
//
// 	"github.com/lindb/lindb/flow"
// 	"github.com/lindb/lindb/kv"
// 	"github.com/lindb/lindb/pkg/bit"
// 	"github.com/lindb/lindb/pkg/encoding"
// 	"github.com/lindb/lindb/pkg/timeutil"
// 	"github.com/lindb/lindb/series/field"
// 	"github.com/lindb/lindb/sql/stmt"
// 	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
// )
//
// var tsdFlushFn = encoding.FlushFunc
//
// func TestFieldWriter_write(t *testing.T) {
// 	buf := make([]byte, pageSize)
// 	md := &memoryDatabase{}
//
// 	// case 1: get write value
// 	write(md, buf, 1, 0, field.SumField, 10, 10.1)
// 	value, ok := getCurrentValue(buf, 10, 10)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 10.1, value, 0)
// 	assert.Equal(t, uint16(0), getEnd(buf))
// 	// case 2: get not exist value, out of time slot range
// 	value, ok = getCurrentValue(buf, 10, 12)
// 	assert.False(t, ok)
// 	assert.Equal(t, 0.0, value)
// 	value, ok = getCurrentValue(buf, 10, 0)
// 	assert.False(t, ok)
// 	assert.Equal(t, 0.0, value)
// 	thisSlotRange := slotRange(getStart(buf), buf, md.getFieldCompressBuffer(1, 0))
// 	assert.Equal(t, uint16(10), thisSlotRange.Start)
// 	assert.Equal(t, uint16(10), thisSlotRange.End)
// 	// case 3: write exist value, need rollup
// 	write(md, buf, 1, 0, field.SumField, 10, 10.1)
// 	value, ok = getCurrentValue(buf, 10, 10)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 20.2, value, 0)
// 	assert.Equal(t, uint16(0), getEnd(buf))
// 	// case 3: write new value
// 	write(md, buf, 1, 0, field.SumField, 12, 12.1)
// 	value, ok = getCurrentValue(buf, 10, 12)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 12.1, value, 0)
// 	assert.Equal(t, uint16(2), getEnd(buf))
// 	// case 4: get value in time slot range
// 	value, ok = getCurrentValue(buf, 10, 11)
// 	assert.False(t, ok)
// 	assert.Equal(t, 0.0, value)
// 	// case 5: test slot range [10,12]
// 	thisSlotRange = slotRange(getStart(buf), buf, md.getFieldCompressBuffer(1, 0))
// 	assert.Equal(t, uint16(10), thisSlotRange.Start)
// 	assert.Equal(t, uint16(12), thisSlotRange.End)
// 	// case 6: compact for slot < start time, time range[5,12]
// 	write(md, buf, 1, 0, field.SumField, 5, 5.3)
// 	thisSlotRange = slotRange(getStart(buf), buf, md.getFieldCompressBuffer(1, 0))
// 	assert.Equal(t, uint16(5), thisSlotRange.Start)
// 	assert.Equal(t, uint16(12), thisSlotRange.End)
// 	value, ok = getCurrentValue(buf, 5, 5)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 5.3, value, 0)
// 	assert.Equal(t, uint16(0), getEnd(buf))
// 	// case 7: write old value
// 	write(md, buf, 1, 0, field.SumField, 10, 10.1)
// 	assert.Equal(t, uint16(5), getEnd(buf))
// 	// case 8: compact for slot > end time, time range[5,12]
// 	write(md, buf, 1, 0, field.SumField, 50, 50.1)
// 	thisSlotRange = slotRange(getStart(buf), buf, md.getFieldCompressBuffer(1, 0))
// 	assert.Nil(t, md.getFieldCompressBuffer(10, 0))
// 	assert.Equal(t, uint16(5), thisSlotRange.Start)
// 	assert.Equal(t, uint16(50), thisSlotRange.End)
// 	value, ok = getCurrentValue(buf, 50, 50)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 50.1, value, 0.0)
// 	assert.Equal(t, uint16(0), getEnd(buf))
// 	// case 9: write 10 slot, compact old value
// 	write(md, buf, 1, 0, field.SumField, 10, 10.1)
// 	assert.Equal(t, uint16(0), getEnd(buf))
// 	value, ok = getCurrentValue(buf, 10, 10)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 10.1, value, 0)
// }
//
// func TestFieldWRiter_write2(t *testing.T) {
// 	md := &memoryDatabase{}
// 	buf := make([]byte, pageSize)
// 	write(md, buf, 1, 0, field.SumField, 10, 178)
// 	value, ok := getCurrentValue(buf, 10, 10)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 178.0, value, 0)
// 	assert.Equal(t, uint16(0), getEnd(buf))
// 	// write with old slot
// 	write(md, buf, 1, 0, field.SumField, 10, 178)
// 	value, ok = getCurrentValue(buf, 10, 10)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 178.0*2, value, 0)
// 	assert.Equal(t, uint16(0), getEnd(buf))
// }
//
// func TestFieldWriter_write_compact_err(t *testing.T) {
// 	defer func() {
// 		encoding.FlushFunc = tsdFlushFn
// 	}()
//
// 	encoding.FlushFunc = func(writer *bit.Writer) error {
// 		return fmt.Errorf("err")
// 	}
// 	buf := make([]byte, pageSize)
// 	md := &memoryDatabase{}
// 	write(md, buf, 1, 0, field.SumField, 10, 10.1)
// 	write(md, buf, 1, 0, field.SumField, 100, 100.1)
// 	value, ok := getCurrentValue(buf, 100, 100)
// 	assert.True(t, ok)
// 	assert.InDelta(t, 100.1, value, 0)
// 	assert.Equal(t, uint16(0), getEnd(buf))
// 	// compress data is nil
// 	assert.Nil(t, md.getFieldCompressBuffer(1, 0))
// }
//
// func TestFieldWriter_flush(t *testing.T) {
// 	nopKVFlusher := kv.NewNopFlusher()
// 	flusher, err := metricsdata.NewFlusher(nopKVFlusher)
// 	assert.NoError(t, err)
// 	fields := field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}}
// 	flusher.PrepareMetric(39, fields)
//
// 	md := &memoryDatabase{}
// 	slotRange := timeutil.SlotRange{Start: 5, End: 100}
// 	for idx, f := range fields {
// 		f.Index = uint8(idx)
// 		buf := make([]byte, pageSize)
// 		write(md, buf, 1, f.Index, field.SumField, 5, float64(f.ID))
// 		write(md, buf, 1, f.Index, field.SumField, 100, 100.1)
// 		assert.NoError(t, flushFieldTo(md, 1, buf, slotRange, flusher, idx,
// 			field.Meta{Type: field.SumField, Index: f.Index}))
// 	}
//
// 	err = flusher.FlushSeries(10)
// 	assert.NoError(t, err)
// 	assert.NoError(t, flusher.CommitMetric(slotRange))
// 	data := nopKVFlusher.Bytes()
// 	r, err := metricsdata.NewReader("1.sst", data)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, r)
//
// 	seriesIDs := roaring.BitmapOf(10)
// 	found := 0
// 	highKeys := seriesIDs.GetHighKeys()
// 	for idx := range highKeys {
// 		highKey := highKeys[idx]
// 		lowSeriesIDs := seriesIDs.GetContainer(highKey)
// 		ctx := &flow.DataLoadContext{
// 			SeriesIDHighKey:       highKey,
// 			LowSeriesIDsContainer: lowSeriesIDs,
// 			ShardExecuteCtx: &flow.ShardExecuteContext{
// 				StorageExecuteCtx: &flow.StorageExecuteContext{
// 					Fields: fields,
// 					Query:  &stmt.Query{},
// 				},
// 			},
// 			DownSampling: func(slotRange timeutil.SlotRange, _ uint16, fieldIdx int, getter encoding.TSDValueGetter) {
// 				assert.Equal(t, timeutil.SlotRange{Start: 5, End: 100}, slotRange)
// 				for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
// 					value, ok := getter.GetValue(movingSourceSlot)
// 					if !ok {
// 						continue
// 					}
// 					if movingSourceSlot != 100 {
// 						assert.Equal(t, 5, int(movingSourceSlot))
// 						assert.Equal(t, value, float64(fields[fieldIdx].ID))
// 					} else {
// 						assert.Equal(t, 100, int(movingSourceSlot))
// 						assert.InDelta(t, 100.1, value, 0)
// 					}
// 					found++
// 				}
// 			},
// 			Decoder: encoding.GetTSDDecoder(),
// 		}
// 		ctx.Grouping()
// 		loader := r.Load(ctx)
// 		loader.Load(ctx)
// 	}
// 	assert.Equal(t, 4, found)
// }
//
// func TestFieldWriter_flush_error(t *testing.T) {
// 	defer func() {
// 		encoding.FlushFunc = tsdFlushFn
// 	}()
//
// 	encoding.FlushFunc = func(writer *bit.Writer) error {
// 		return fmt.Errorf("err")
// 	}
// 	nopKVFlusher := kv.NewNopFlusher()
// 	flusher, err := metricsdata.NewFlusher(nopKVFlusher)
// 	assert.NoError(t, err)
// 	fields := field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}}
// 	flusher.PrepareMetric(39, fields)
// 	buf := make([]byte, pageSize)
// 	md := &memoryDatabase{}
// 	slotRange := timeutil.SlotRange{Start: 5, End: 100}
// 	write(md, buf, 1, 0, field.SumField, 10, 10.1)
// 	write(md, buf, 1, 0, field.SumField, 100, 100.1)
// 	assert.NoError(t, flushFieldTo(md, 1, buf, slotRange, flusher, 0,
// 		field.Meta{Type: field.SumField, Index: 0}))
// }
