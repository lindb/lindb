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

import (
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

func TestFlusher_flush_metric(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()
	flusher, err := NewFlusher(nopKVFlusher)
	assert.NoError(t, err)
	flusher.PrepareMetric(39,
		[]field.Meta{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}},
	)
	// no field for series
	assert.NoError(t, flusher.FlushSeries(5))

	assert.NoError(t, flusher.FlushField([]byte{1, 2, 3}))
	assert.NoError(t, flusher.FlushField([]byte{10, 20, 30}))
	assert.NoError(t, flusher.FlushSeries(10))

	// flush has one field
	assert.NoError(t, flusher.FlushField([]byte{10, 20, 30}))
	assert.NoError(t, flusher.FlushField(nil))
	assert.NoError(t, flusher.FlushSeries(100))

	f := flusher.GetFieldMetas()
	assert.Equal(t, field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}}, f)
	assert.NoError(t, flusher.CommitMetric(timeutil.SlotRange{Start: 10, End: 13}))
	assert.NoError(t, err)

	// field not exist, not flush metric
	assert.Empty(t, flusher.GetFieldMetas())
	flusher.PrepareMetric(40, []field.Meta{{ID: 1, Type: field.SumField}})
	assert.NoError(t, flusher.FlushField([]byte{1, 2, 3}))
	assert.NoError(t, flusher.CommitMetric(timeutil.SlotRange{Start: 10, End: 13}))

	// metric hasn't series ids
	flusher.PrepareMetric(50, []field.Meta{{ID: 1, Type: field.SumField}})
	assert.NoError(t, flusher.FlushField(nil))
	assert.NoError(t, flusher.CommitMetric(timeutil.SlotRange{Start: 10, End: 13}))

	// close
	assert.NoError(t, flusher.Close())
}

func TestFlusher_flush_big_series_id(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()
	flusher, _ := NewFlusher(nopKVFlusher)
	flusher.PrepareMetric(39, []field.Meta{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	assert.NoError(t, flusher.FlushField([]byte{1, 2, 3}))
	assert.NoError(t, flusher.FlushSeries(10000))
	assert.NoError(t, flusher.CommitMetric(timeutil.SlotRange{Start: 10, End: 13}))

	assert.Empty(t, flusher.GetFieldMetas())
	assert.NoError(t, flusher.Close())
}

func TestFlusher_TooMany_Data(t *testing.T) {
	flushMoreData(t, field.Metas{{ID: 1, Type: field.SumField}},
		mockSingleField, 1.0, field.Metas{{ID: 1, Type: field.SumField}})
	flushMoreData(t, field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}},
		mockMultiField1, 1.0, field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	flushMoreData(t, field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}},
		mockMultiField2, 2.0, field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	flushMoreData(t, field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}},
		mockMultiField2, 1.0, field.Metas{{ID: 2, Type: field.SumField}})
	flushMoreData(t, field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}},
		mockMultiField2, 1.0, field.Metas{{ID: 1, Type: field.SumField}})
}

func flushMoreData(t *testing.T,
	fields field.Metas,
	mockData func(t *testing.T,
		seriesIDs *roaring.Bitmap, flusher Flusher), assertRatio float64, queryFields field.Metas) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nopKVFlusher := kv.NewNopFlusher()
	flusher, _ := NewFlusher(nopKVFlusher)
	flusher.PrepareMetric(39, fields)

	seriesIDs := roaring.New()
	mockData(t, seriesIDs, flusher)
	assert.NoError(t, flusher.CommitMetric(timeutil.SlotRange{Start: 5, End: 5}))
	data := nopKVFlusher.Bytes()
	r, err := NewReader("1.sst", data)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	found := 0
	highKeys := seriesIDs.GetHighKeys()
	tsdDecoder := encoding.GetTSDDecoder()
	for idx := range highKeys {
		highKey := highKeys[idx]
		lowSeriesIDs := seriesIDs.GetContainer(highKey)
		ctx := &flow.DataLoadContext{
			SeriesIDHighKey:       highKey,
			LowSeriesIDsContainer: lowSeriesIDs,
			ShardExecuteCtx: &flow.ShardExecuteContext{
				StorageExecuteCtx: &flow.StorageExecuteContext{
					Fields: queryFields,
					Query:  &stmt.Query{},
				},
			},
			DownSampling: func(slotRange timeutil.SlotRange, seriesIdx uint16, fieldIdx int, fieldData []byte) {
				assert.Equal(t, timeutil.SlotRange{Start: 5, End: 5}, slotRange)
				tsdDecoder.ResetWithTimeRange(fieldData, slotRange.Start, slotRange.End)
				for movingSourceSlot := tsdDecoder.StartTime(); movingSourceSlot <= tsdDecoder.EndTime(); movingSourceSlot++ {
					if !tsdDecoder.HasValueWithSlot(movingSourceSlot) {
						continue
					}
					value := math.Float64frombits(tsdDecoder.Value())
					assert.Equal(t, 5, int(movingSourceSlot))
					seriesID := float64(int(highKey)*65536 + int(seriesIdx))
					assert.Equal(t, value, seriesID)
					found++
				}
			},
		}
		ctx.Grouping()
		loader := r.Load(ctx)
		loader.Load(ctx)
	}
	assert.Equal(t, int(seriesIDs.GetCardinality())*int(assertRatio), found)
}

func mockSingleField(t *testing.T, seriesIDs *roaring.Bitmap, flusher Flusher) {
	for i := 0; i < 80000; i++ {
		seriesIDs.Add(uint32(i))
		encoder := encoding.NewTSDEncoder(5)
		encoder.AppendTime(bit.One)
		encoder.AppendValue(math.Float64bits(float64(i)))
		data, _ := encoder.BytesWithoutTime()
		assert.NoError(t, flusher.FlushField(data))
		assert.NoError(t, flusher.FlushSeries(uint32(i)))
	}
}

func mockMultiField1(t *testing.T, seriesIDs *roaring.Bitmap, flusher Flusher) {
	for i := 0; i < 80000; i++ {
		seriesIDs.Add(uint32(i))
		encoder := encoding.NewTSDEncoder(5)
		encoder.AppendTime(bit.One)
		encoder.AppendValue(math.Float64bits(float64(i)))
		data, _ := encoder.BytesWithoutTime()
		assert.NoError(t, flusher.FlushField(data))
		assert.NoError(t, flusher.FlushField(nil))
		assert.NoError(t, flusher.FlushSeries(uint32(i)))
	}
}

func mockMultiField2(t *testing.T, seriesIDs *roaring.Bitmap, flusher Flusher) {
	for i := 0; i < 80000; i++ {
		seriesIDs.Add(uint32(i))
		encoder := encoding.NewTSDEncoder(5)
		encoder.AppendTime(bit.One)
		encoder.AppendValue(math.Float64bits(float64(i)))
		data, _ := encoder.BytesWithoutTime()
		assert.NoError(t, flusher.FlushField(data))
		assert.NoError(t, flusher.FlushField(data))
		assert.NoError(t, flusher.FlushSeries(uint32(i)))
	}
}
