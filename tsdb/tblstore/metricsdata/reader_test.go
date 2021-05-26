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
	"fmt"
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

var bitmapUnmarshal = encoding.BitmapUnmarshal

func TestNewReader(t *testing.T) {
	defer func() {
		encoding.BitmapUnmarshal = bitmapUnmarshal
	}()
	// case 1: footer err
	r, err := NewReader("1.sst", []byte{1, 2, 3})
	assert.Error(t, err)
	assert.Nil(t, r)
	// case 2: offset err
	r, err = NewReader("1.sst", []byte{0, 0, 0, 1, 2, 3, 3, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 1, 2, 3, 4})
	assert.Error(t, err)
	assert.Nil(t, r)
	// case 3: new metricReader success
	r, err = NewReader("1.sst", mockMetricBlock())
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "1.sst", r.Path())
	timeRange := r.GetTimeRange()
	assert.Equal(t, uint16(5), timeRange.Start)
	assert.Equal(t, uint16(5), timeRange.End)
	assert.Equal(t, field.Metas{
		{ID: 2, Type: field.SumField},
		{ID: 10, Type: field.MinField},
		{ID: 30, Type: field.SumField},
		{ID: 100, Type: field.MaxField},
	}, r.GetFields())
	seriesIDs := roaring.New()
	for j := 0; j < 10; j++ {
		seriesIDs.Add(uint32(j * 4096))
	}
	seriesIDs.Add(65536 + 10)
	assert.EqualValues(t, seriesIDs.ToArray(), r.GetSeriesIDs().ToArray())
	// case 4: unmarshal series ids err
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	r, err = NewReader("1.sst", mockMetricBlock())
	assert.Error(t, err)
	assert.Nil(t, r)
}

func TestReader_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r, err := NewReader("1.sst", mockMetricBlock())
	assert.NoError(t, err)
	assert.NotNil(t, r)
	// case 1: series high key not found
	r.Load(1000, nil, field.Metas{{ID: 2}, {ID: 30}, {ID: 50}})
	// case 3: load data success
	r, err = NewReader("1.sst", mockMetricBlock())
	assert.NoError(t, err)
	scanner := r.Load(0, roaring.BitmapOf(4096, 8192).GetContainer(0), field.Metas{{ID: 2}, {ID: 30}, {ID: 50}})
	assert.NotNil(t, scanner)
	// case 4: series ids not found
	r, err = NewReader("1.sst", mockMetricBlock())
	assert.NoError(t, err)
	scanner = r.Load(0,
		roaring.BitmapOf(10, 12).GetContainer(0),
		field.Metas{{ID: 2}, {ID: 30}, {ID: 50}})
	assert.Nil(t, scanner)

	// case 5: load data success, but time slot not in query range
	r, err = NewReader("1.sst", mockMetricBlock())
	assert.NoError(t, err)
	scanner = r.Load(0,
		roaring.BitmapOf(4096, 8192).GetContainer(0),
		field.Metas{{ID: 2}, {ID: 30}, {ID: 50}})
	scanner.Load(4096)
	scanner.Load(8192)

	// case 6: load data success, metric has one field
	r, err = NewReader("1.sst", mockMetricBlockForOneField())
	assert.NoError(t, err)
	scanner = r.Load(0, roaring.BitmapOf(4096, 8192).GetContainer(0), field.Metas{{ID: 2}})
	scanner.Load(4096)
	scanner.Load(8192)
	// case 7: high key not exist
	r, err = NewReader("1.sst", mockMetricBlockForOneField())
	assert.NoError(t, err)
	scanner = r.Load(10, roaring.BitmapOf(4096, 8192).GetContainer(0), field.Metas{{ID: 2}})
	assert.Nil(t, scanner)
}

func TestReader_scan(t *testing.T) {
	defer func() {
		getOffsetFunc = getOffset
	}()

	r, err := NewReader("1.sst", mockMetricBlock())
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner := newDataScanner(r)
	timeRange := scanner.slotRange()
	assert.Equal(t, uint16(5), timeRange.Start)
	assert.Equal(t, uint16(5), timeRange.End)
	// case 1: not match
	seriesPos := scanner.scan(10, 10)
	assert.True(t, seriesPos < 0)
	// case 2: merge data
	scanner = newDataScanner(r)
	seriesPos = scanner.scan(0, 0)
	assert.True(t, seriesPos >= 0)
	seriesPos = scanner.scan(1, 10)
	assert.True(t, seriesPos >= 0)
	// case 3: scan completed
	seriesPos = scanner.scan(3, 10)
	assert.True(t, seriesPos < 0)
	// case 4: not match
	scanner = newDataScanner(r)
	seriesPos = scanner.scan(0, 10)
	assert.True(t, seriesPos < 0)
	// case 6: get wrong offset
	scanner = newDataScanner(r)
	getOffsetFunc = func(seriesOffsets *encoding.FixedOffsetDecoder, idx int) (int, bool) {
		return 0, false
	}
	seriesPos = scanner.scan(0, 0)
	assert.True(t, seriesPos < 0)
	fields := scanner.fieldIndexes()
	assert.Len(t, fields, 4)
}

func mockMetricBlock() []byte {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas(field.Metas{
		{ID: 2, Type: field.SumField},
		{ID: 10, Type: field.MinField},
		{ID: 30, Type: field.SumField},
		{ID: 100, Type: field.MaxField},
	})
	for j := 0; j < 10; j++ {
		encoder := encoding.NewTSDEncoder(5)
		for i := 0; i < 10; i++ {
			encoder.AppendTime(bit.One)
			encoder.AppendValue(math.Float64bits(float64(10.0 * i)))
		}
		data, _ := encoder.BytesWithoutTime()
		flusher.FlushField(data)
		flusher.FlushField(data)
		flusher.FlushField(data)
		flusher.FlushField(data)
		flusher.FlushSeries(uint32(j * 4096))
	}
	// mock just has one field
	encoder := encoding.NewTSDEncoder(5)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	data, _ := encoder.BytesWithoutTime()
	flusher.FlushField(data)
	flusher.FlushSeries(uint32(65536 + 10))
	_ = flusher.FlushMetric(uint32(10), 5, 5)

	return nopKVFlusher.Bytes()
}

func mockMetricBlockForOneField() []byte {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas(field.Metas{
		{ID: 2, Type: field.SumField},
	})
	for j := 0; j < 10; j++ {
		encoder := encoding.NewTSDEncoder(5)
		for i := 0; i < 10; i++ {
			encoder.AppendTime(bit.One)
			encoder.AppendValue(math.Float64bits(float64(10.0 * i)))
		}
		data, _ := encoder.BytesWithoutTime()
		flusher.FlushField(data)
		flusher.FlushSeries(uint32(j * 4096))
	}
	// mock just has one field
	encoder := encoding.NewTSDEncoder(5)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	data, _ := encoder.BytesWithoutTime()
	flusher.FlushField(data)
	flusher.FlushSeries(uint32(65536 + 10))
	_ = flusher.FlushMetric(uint32(10), 5, 5)
	return nopKVFlusher.Bytes()
}
