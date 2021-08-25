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
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func TestSeriesMerger_compact_merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := NewMockFlusher(ctrl)
	merger := newSeriesMerger(flusher)
	decodeStreams := make([]*encoding.TSDDecoder, 3)
	reader1 := NewMockFieldReader(ctrl)
	reader2 := NewMockFieldReader(ctrl)
	reader1.EXPECT().close().AnyTimes()
	reader2.EXPECT().close().AnyTimes()
	readers := []FieldReader{reader1, nil, reader2}

	encodeStream := encoding.NewTSDEncoder(5)
	// case 1: merge success and rollup
	reader1.EXPECT().getFieldData(gomock.Any()).Return(mockField(10))
	reader1.EXPECT().slotRange().Return(uint16(10), uint16(10))
	reader2.EXPECT().getFieldData(gomock.Any()).Return(mockField(10))
	reader2.EXPECT().slotRange().Return(uint16(10), uint16(10))
	var result []byte
	flusher.EXPECT().FlushField(gomock.Any()).DoAndReturn(func(data []byte) {
		result = data
	})
	err := merger.merge(
		&mergerContext{
			targetFields: field.Metas{{ID: 1, Type: field.SumField}},
			sourceRange:  timeutil.SlotRange{Start: 5, End: 15},
			targetRange:  timeutil.SlotRange{Start: 5, End: 15},
			ratio:        1,
		}, decodeStreams, encodeStream, readers)
	assert.NoError(t, err)
	tsd := encoding.GetTSDDecoder()
	tsd.ResetWithTimeRange(result, 5, 15)
	slot := uint16(0)
	for i := uint16(5); i <= 15; i++ {
		if tsd.HasValueWithSlot(i) {
			slot = i
			assert.Equal(t, 20.0, math.Float64frombits(tsd.Value()))
		}
	}
	assert.Equal(t, uint16(10), slot)
	// case 2: merge success with diff slot range
	reader1.EXPECT().getFieldData(gomock.Any()).Return(mockField(10))
	reader1.EXPECT().slotRange().Return(uint16(10), uint16(10))
	reader2.EXPECT().getFieldData(gomock.Any()).Return(mockField(12))
	reader2.EXPECT().slotRange().Return(uint16(12), uint16(12))
	flusher.EXPECT().FlushField(gomock.Any()).DoAndReturn(func(data []byte) {
		result = data
	})
	err = merger.merge(
		&mergerContext{
			targetFields: field.Metas{{ID: 1, Type: field.SumField}},
			sourceRange:  timeutil.SlotRange{Start: 5, End: 15},
			targetRange:  timeutil.SlotRange{Start: 5, End: 15},
			ratio:        1,
		}, decodeStreams, encodeStream, readers)
	assert.NoError(t, err)
	tsd.ResetWithTimeRange(result, 5, 15)
	c := 0
	for i := uint16(5); i <= 15; i++ {
		if tsd.HasValueWithSlot(i) && (i == 10 || i == 12) {
			c++
			assert.Equal(t, 10.0, math.Float64frombits(tsd.Value()))
		}
	}
	assert.Equal(t, 2, c)
	// case 3: encode stream err
	encodeStream2 := encoding.NewMockTSDEncoder(ctrl)
	reader1.EXPECT().getFieldData(gomock.Any()).Return(mockField(10))
	reader1.EXPECT().slotRange().Return(uint16(10), uint16(10))
	reader2.EXPECT().getFieldData(gomock.Any()).Return(mockField(12))
	reader2.EXPECT().slotRange().Return(uint16(12), uint16(12))
	encodeStream2.EXPECT().EmitDownSamplingValue(gomock.Any(), gomock.Any()).AnyTimes()
	encodeStream2.EXPECT().BytesWithoutTime().Return(nil, fmt.Errorf("err"))
	err = merger.merge(
		&mergerContext{
			targetFields: field.Metas{{ID: 1, Type: field.SumField}},
			sourceRange:  timeutil.SlotRange{Start: 5, End: 15},
			targetRange:  timeutil.SlotRange{Start: 5, End: 15},
			ratio:        1,
		}, decodeStreams, encodeStream2, readers)
	assert.Error(t, err)
}

func TestSeriesMerger_rollup_merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := NewMockFlusher(ctrl)
	merger := newSeriesMerger(flusher)
	decodeStreams := make([]*encoding.TSDDecoder, 3)
	reader1 := NewMockFieldReader(ctrl)
	reader2 := NewMockFieldReader(ctrl)
	reader1.EXPECT().close().AnyTimes()
	reader2.EXPECT().close().AnyTimes()
	readers := []FieldReader{reader1, reader2, nil}

	encodeStream := encoding.NewTSDEncoder(5)
	// case 1: merge success and rollup
	reader1.EXPECT().getFieldData(gomock.Any()).Return(mockField(10))
	reader1.EXPECT().slotRange().Return(uint16(10), uint16(10))
	reader2.EXPECT().getFieldData(gomock.Any()).Return(mockField(10))
	reader2.EXPECT().slotRange().Return(uint16(12), uint16(12))
	var result []byte
	flusher.EXPECT().FlushField(gomock.Any()).DoAndReturn(func(data []byte) {
		result = data
	})
	// source:[5,15] target:[0,0], interval: 10s => 5min
	err := merger.merge(
		&mergerContext{
			targetFields: field.Metas{{ID: 1, Type: field.SumField}},
			sourceRange:  timeutil.SlotRange{Start: 5, End: 15},
			targetRange:  timeutil.SlotRange{Start: 0, End: 0},
			ratio:        30,
		}, decodeStreams, encodeStream, readers)
	assert.NoError(t, err)
	tsd := encoding.GetTSDDecoder()
	tsd.ResetWithTimeRange(result, 0, 0)
	slot := uint16(0)
	for i := uint16(0); i <= 0; i++ {
		if tsd.HasValueWithSlot(i) {
			slot = i
			assert.Equal(t, 20.0, math.Float64frombits(tsd.Value()))
		}
	}
	assert.Equal(t, uint16(0), slot)
	// case 2: merge success and rollup
	reader1.EXPECT().getFieldData(gomock.Any()).Return(mockField(10))
	reader1.EXPECT().slotRange().Return(uint16(10), uint16(10))
	reader2.EXPECT().getFieldData(gomock.Any()).Return(mockField(10))
	reader2.EXPECT().slotRange().Return(uint16(182), uint16(182))
	flusher.EXPECT().FlushField(gomock.Any()).DoAndReturn(func(data []byte) {
		result = data
	})
	// source:[5,182] target:[0,6], interval: 10s => 5min
	err = merger.merge(
		&mergerContext{
			targetFields: field.Metas{{ID: 1, Type: field.SumField}},
			sourceRange:  timeutil.SlotRange{Start: 5, End: 182},
			targetRange:  timeutil.SlotRange{Start: 0, End: 6},
			ratio:        30,
		}, decodeStreams, encodeStream, readers)
	assert.NoError(t, err)
	tsd = encoding.GetTSDDecoder()
	tsd.ResetWithTimeRange(result, 0, 6)
	c := 0
	for i := uint16(0); i <= 6; i++ {
		if tsd.HasValueWithSlot(i) && (i == 0 || i == 6) {
			assert.Equal(t, 10.0, math.Float64frombits(tsd.Value()))
			c++
		}
	}
	assert.Equal(t, 2, c)
}

func mockField(start uint16) []byte {
	encoder := encoding.NewTSDEncoder(start)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	data, _ := encoder.BytesWithoutTime()
	return data
}
