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

package cumulativecache

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/fasttime"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/memdb"
)

func Test_Cache_panics(t *testing.T) {
	assert.Panics(t, func() {
		_ = NewCache(20, time.Second, time.Second, linmetric.NewScope("11"))
	})
}

func newTestPoint() *memdb.MetricPoint {
	return &memdb.MetricPoint{
		MetricID:  1,
		SeriesID:  2,
		SlotIndex: 1,
		FieldIDs:  []field.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Proto: &protoMetricsV1.Metric{
			Namespace: "default-ns",
			Name:      "test-metric",
			Timestamp: fasttime.UnixMilliseconds(),
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "0", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
				{Name: "1", Type: protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM, Value: 1},
				{Name: "2", Type: protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM, Value: 1},
				{Name: "3", Type: protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM, Value: 1},
				{Name: "4", Type: protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM, Value: 1},
			},
			CompoundField: &protoMetricsV1.CompoundField{
				Type:           protoMetricsV1.CompoundFieldType_CUMULATIVE_HISTOGRAM,
				Sum:            1,
				Count:          1,
				ExplicitBounds: []float64{1, 2, 3, 4, math.Inf(1) + 1},
				Values:         []float64{1, 1, 1, 1, 1},
			},
		},
	}
}

func Test_CumulativePointToDelta(t *testing.T) {
	cache := NewCache(32, time.Second, time.Second, linmetric.NewScope("11"))
	mp := newTestPoint()
	assert.Equal(t, cache.Capacity(), 0)
	assert.False(t, cache.CumulativePointToDelta(mp))
	assert.Equal(t, cache.Capacity(), 1)
	assert.True(t, cache.CumulativePointToDelta(mp))
	assert.Equal(t, cache.Capacity(), 1)

	time.Sleep(time.Millisecond * 2000)
	assert.Equal(t, 0, cache.Capacity())
	cache.Close()
}

func Test_CumulativePointToDelta_SimpleFields(t *testing.T) {
	mp := newTestPoint()
	mp.Proto.CompoundField = nil
	assert.Equal(t, countCumulativeFields(mp), 4)
	cache := NewCache(32, time.Second, time.Second, linmetric.NewScope("11"))
	assert.Equal(t, cache.Capacity(), 0)
	assert.False(t, cache.CumulativePointToDelta(mp))
	assert.Equal(t, cache.Capacity(), 1)
	assert.True(t, cache.CumulativePointToDelta(mp))
}

func Test_CumulativePointToDelta_SimpleFields2(t *testing.T) {
	mp := newTestPoint()
	mp.Proto.CompoundField.Type = protoMetricsV1.CompoundFieldType_DELTA_HISTOGRAM
	assert.Equal(t, countCumulativeFields(mp), 4)
	cache := NewCache(32, time.Second, time.Second, linmetric.NewScope("12"))
	assert.Equal(t, cache.Capacity(), 0)
	assert.False(t, cache.CumulativePointToDelta(mp))
	assert.Equal(t, cache.Capacity(), 1)
	assert.True(t, cache.CumulativePointToDelta(mp))
}

func Test_CumulativePointToDelta_CompoundFields(t *testing.T) {
	mp := newTestPoint()
	mp.Proto.SimpleFields = nil

	assert.Equal(t, countCumulativeFields(mp), 7)
	cache := NewCache(32, time.Second, time.Second, linmetric.NewScope("13"))
	assert.Equal(t, 0, cache.Capacity())
	assert.False(t, cache.CumulativePointToDelta(mp))
	assert.Equal(t, 1, cache.Capacity())
	assert.True(t, cache.CumulativePointToDelta(mp))
}

func Test_decodeCumulativeFieldsInto(t *testing.T) {
	mp := newTestPoint()

	assert.False(t, decodeCumulativeFieldsInto(mp, nil))
	cache := NewCache(32, time.Second, time.Second, linmetric.NewScope("14"))
	assert.False(t, cache.CumulativePointToDelta(mp))
	// fieldid not match
	mp2 := newTestPoint()
	mp2.FieldIDs[1] = field.ID(222)
	assert.False(t, cache.CumulativePointToDelta(mp2))
	mp3 := newTestPoint()
	mp3.Proto.SimpleFields[1].Value = 0
	assert.False(t, cache.CumulativePointToDelta(mp3))

	// sum
	mp4 := newTestPoint()
	mp4.Proto.CompoundField.Sum = 0
	assert.False(t, cache.CumulativePointToDelta(mp4))

	// count
	mp5 := newTestPoint()
	mp5.Proto.CompoundField.Count = 0
	assert.False(t, cache.CumulativePointToDelta(mp5))

	// value
	mp6 := newTestPoint()
	mp6.Proto.CompoundField.Values[0] = 0
	assert.False(t, cache.CumulativePointToDelta(mp6))
}
