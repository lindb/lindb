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

package metricsdata

// import (
// 	"fmt"
// 	"math"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"
//
// 	"github.com/lindb/roaring"
//
// 	"github.com/lindb/lindb/flow"
// 	"github.com/lindb/lindb/kv"
// 	"github.com/lindb/lindb/pkg/bit"
// 	"github.com/lindb/lindb/pkg/encoding"
// 	"github.com/lindb/lindb/pkg/timeutil"
// 	"github.com/lindb/lindb/series/field"
// 	"github.com/lindb/lindb/sql/stmt"
// )
//
// var bitmapUnmarshal = encoding.BitmapUnmarshal
//
// func TestNewReader(t *testing.T) {
// 	defer func() {
// 		encoding.BitmapUnmarshal = bitmapUnmarshal
// 	}()
// 	// case 1: footer err
// 	r, err := NewReader("1.sst", []byte{1, 2, 3})
// 	assert.Error(t, err)
// 	assert.Nil(t, r)
// 	// case 2: offset err
// 	r, err = NewReader("1.sst", []byte{0, 0, 0, 1, 2, 3, 3, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 1, 2, 3, 4})
// 	assert.Error(t, err)
// 	assert.Nil(t, r)
// 	// case 3: new metricReader success
// 	r, err = NewReader("1.sst", mockMetricBlock())
// 	assert.NoError(t, err)
// 	assert.NotNil(t, r)
// 	assert.Equal(t, "1.sst", r.Path())
// 	timeRange := r.GetTimeRange()
// 	assert.Equal(t, uint16(5), timeRange.Start)
// 	assert.Equal(t, uint16(5), timeRange.End)
// 	assert.Equal(t, field.Metas{
// 		{ID: 2, Type: field.SumField},
// 		{ID: 10, Type: field.MinField},
// 		{ID: 30, Type: field.SumField},
// 		{ID: 100, Type: field.MaxField},
// 	}, r.GetFields())
// 	seriesIDs := roaring.New()
// 	for j := 0; j < 10; j++ {
// 		seriesIDs.Add(uint32(j * 4096))
// 	}
// 	seriesIDs.Add(65536 + 10)
// 	assert.EqualValues(t, seriesIDs.ToArray(), r.GetSeriesIDs().ToArray())
// 	// case 4: unmarshal series ids err
// 	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) (int64, error) {
// 		return 0, fmt.Errorf("err")
// 	}
// 	r, err = NewReader("1.sst", mockMetricBlock())
// 	assert.Error(t, err)
// 	assert.Nil(t, r)
// }
//
// func TestReader_Load(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	r, err := NewReader("1.sst", mockMetricBlock())
// 	assert.NoError(t, err)
// 	assert.NotNil(t, r)
// 	r1 := r.(*metricReader)
// 	assert.Len(t, r1.fieldIndexes(), 4)
// 	// case 1: series high key not found
// 	r.Load(&flow.DataLoadContext{
// 		SeriesIDHighKey: 1000,
// 		ShardExecuteCtx: &flow.ShardExecuteContext{
// 			StorageExecuteCtx: &flow.StorageExecuteContext{
// 				Fields: field.Metas{{ID: 2}, {ID: 30}, {ID: 50}},
// 			},
// 		},
// 	})
// 	// case 3: load data success
// 	r, err = NewReader("1.sst", mockMetricBlock())
// 	assert.NoError(t, err)
// 	scanner := r.Load(&flow.DataLoadContext{
// 		SeriesIDHighKey:       0,
// 		LowSeriesIDsContainer: roaring.BitmapOf(4096, 8192).GetContainer(0),
// 		ShardExecuteCtx: &flow.ShardExecuteContext{
// 			StorageExecuteCtx: &flow.StorageExecuteContext{
// 				Fields: field.Metas{{ID: 2}, {ID: 30}, {ID: 50}},
// 			},
// 		},
// 	})
//
// 	assert.NotNil(t, scanner)
// 	// case 4: series ids not found
// 	r, err = NewReader("1.sst", mockMetricBlock())
// 	assert.NoError(t, err)
// 	scanner = r.Load(&flow.DataLoadContext{
// 		SeriesIDHighKey:       0,
// 		LowSeriesIDsContainer: roaring.BitmapOf(10, 12).GetContainer(0),
// 		ShardExecuteCtx: &flow.ShardExecuteContext{
// 			StorageExecuteCtx: &flow.StorageExecuteContext{
// 				Fields: field.Metas{{ID: 2}, {ID: 30}, {ID: 50}},
// 			},
// 		},
// 	})
// 	assert.Nil(t, scanner)
//
// 	found := 0
// 	ctx := &flow.DataLoadContext{
// 		SeriesIDHighKey:       0,
// 		LowSeriesIDsContainer: roaring.BitmapOf(4096, 8192).GetContainer(0),
// 		ShardExecuteCtx: &flow.ShardExecuteContext{
// 			StorageExecuteCtx: &flow.StorageExecuteContext{
// 				Fields: field.Metas{{ID: 2}, {ID: 30}, {ID: 50}},
// 				Query:  &stmt.Query{},
// 			},
// 		},
// 		DownSampling: func(slotRange timeutil.SlotRange, _ uint16, _ int, getter encoding.TSDValueGetter) {
// 			assert.Equal(t, timeutil.SlotRange{Start: 5, End: 5}, slotRange)
// 			for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
// 				if _, ok := getter.GetValue(movingSourceSlot); !ok {
// 					continue
// 				}
// 				found++
// 			}
// 		},
// 		Decoder: encoding.GetTSDDecoder(),
// 	}
// 	ctx.Grouping()
// 	scanner = r.Load(ctx)
// 	// case 5: load data success, metric has one field
// 	r, err = NewReader("1.sst", mockMetricBlockForOneField())
// 	assert.NotNil(t, r)
// 	assert.NoError(t, err)
// 	ctx.ShardExecuteCtx.StorageExecuteCtx.Fields = field.Metas{{ID: 2}}
// 	ctx.Grouping()
// 	scanner.Load(ctx)
// 	// case 6: high key not exist
// 	r, err = NewReader("1.sst", mockMetricBlockForOneField())
// 	assert.NoError(t, err)
// 	ctx.SeriesIDHighKey = 10
// 	ctx.Grouping()
// 	scanner = r.Load(ctx)
// 	assert.Nil(t, scanner)
// }
//
// func TestReader_scan(t *testing.T) {
// 	r, err := NewReader("1.sst", mockMetricBlock())
// 	assert.NoError(t, err)
// 	assert.NotNil(t, r)
// 	scanner, err := newDataScanner(r)
// 	assert.NoError(t, err)
//
// 	timeRange := scanner.slotRange()
// 	assert.Equal(t, uint16(5), timeRange.Start)
// 	assert.Equal(t, uint16(5), timeRange.End)
//
// 	// case 1: not match
// 	seriesEntry := scanner.scan(10, 10)
// 	assert.Empty(t, seriesEntry)
// 	// case 2: merge data
// 	scanner, _ = newDataScanner(r)
// 	seriesEntry = scanner.scan(0, 0)
// 	assert.True(t, scanner.reader.seriesIDs.Contains(0))
// 	assert.True(t, len(seriesEntry) > 0)
// 	// case 3: scan completed
// 	seriesEntry = scanner.scan(3, 10)
// 	assert.Empty(t, seriesEntry)
// 	// case 4: not match
// 	scanner, _ = newDataScanner(r)
// 	seriesEntry = scanner.scan(0, 10)
// 	assert.Empty(t, seriesEntry)
// }
//
// func mockMetricBlock() []byte {
// 	nopKVFlusher := kv.NewNopFlusher()
// 	flusher, _ := NewFlusher(nopKVFlusher)
// 	flusher.PrepareMetric(10, field.Metas{
// 		{ID: 2, Type: field.SumField},
// 		{ID: 10, Type: field.MinField},
// 		{ID: 30, Type: field.SumField},
// 		{ID: 100, Type: field.MaxField},
// 	})
//
// 	for j := 0; j < 10; j++ {
// 		encoder := encoding.NewTSDEncoder(5)
// 		for i := 0; i < 10; i++ {
// 			encoder.AppendTime(bit.One)
// 			encoder.AppendValue(math.Float64bits(float64(10.0 * i)))
// 		}
// 		data, _ := encoder.BytesWithoutTime()
// 		_ = flusher.FlushField(data)
// 		_ = flusher.FlushField(data)
// 		_ = flusher.FlushField(data)
// 		_ = flusher.FlushField(data)
// 		_ = flusher.FlushSeries(uint32(j * 4096))
// 	}
// 	// mock just has one field
// 	encoder := encoding.NewTSDEncoder(5)
// 	encoder.AppendTime(bit.One)
// 	encoder.AppendValue(math.Float64bits(10.0))
// 	data, _ := encoder.BytesWithoutTime()
// 	_ = flusher.FlushField(data)
// 	_ = flusher.FlushSeries(uint32(65536 + 10))
// 	_ = flusher.CommitMetric(timeutil.SlotRange{Start: 5, End: 5})
//
// 	return nopKVFlusher.Bytes()
// }
//
// func mockMetricBlockForOneField() []byte {
// 	nopKVFlusher := kv.NewNopFlusher()
// 	flusher, _ := NewFlusher(nopKVFlusher)
// 	flusher.PrepareMetric(10, field.Metas{{ID: 2, Type: field.SumField}})
//
// 	for j := 0; j < 10; j++ {
// 		encoder := encoding.NewTSDEncoder(5)
// 		for i := 0; i < 10; i++ {
// 			encoder.AppendTime(bit.One)
// 			encoder.AppendValue(math.Float64bits(float64(10.0 * i)))
// 		}
// 		data, _ := encoder.BytesWithoutTime()
// 		_ = flusher.FlushField(data)
// 		_ = flusher.FlushSeries(uint32(j * 4096))
// 	}
// 	// mock just has one field
// 	encoder := encoding.NewTSDEncoder(5)
// 	encoder.AppendTime(bit.One)
// 	encoder.AppendValue(math.Float64bits(10.0))
// 	data, _ := encoder.BytesWithoutTime()
// 	_ = flusher.FlushField(data)
// 	_ = flusher.FlushSeries(uint32(65536 + 10))
// 	_ = flusher.CommitMetric(timeutil.SlotRange{Start: 5, End: 5})
// 	return nopKVFlusher.Bytes()
// }
//
// func Benchmark_unmarshal_roaring(b *testing.B) {
// 	r := roaring.New()
// 	for i := 0; i < 100000; i += 2 {
// 		r.Add(uint32(i))
// 	}
// 	r.RunOptimize()
// 	data, _ := r.MarshalBinary()
//
// 	r2 := roaring.New()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = r2.UnmarshalBinary(data)
// 	}
// }
//
// func Benchmark_Roaring_FromBuffer(b *testing.B) {
// 	r := roaring.New()
// 	for i := 0; i < 100000; i += 2 {
// 		r.Add(uint32(i))
// 	}
// 	r.RunOptimize()
// 	data, _ := r.MarshalBinary()
//
// 	r2 := roaring.New()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_, _ = r2.FromBuffer(data)
// 	}
// }
//
// func Benchmark_Roaring_FrozenView(b *testing.B) {
// 	r := roaring.New()
// 	for i := 0; i < 100000; i += 2 {
// 		r.Add(uint32(i))
// 	}
// 	r.RunOptimize()
// 	data, _ := r.Freeze()
//
// 	r2 := roaring.New()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = r2.FrozenView(data)
// 	}
// }
