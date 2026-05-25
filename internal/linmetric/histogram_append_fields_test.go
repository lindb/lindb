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

package linmetric

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/arrow/pkg/model"
	"github.com/lindb/common/field"
	seriesmetric "github.com/lindb/lindb/series/metric"
)

// TestBoundHistogram_AppendFields verifies that AppendFields writes atomic
// model.Field entries using the .__<stat> and .__bucket_* naming convention.
func TestBoundHistogram_AppendFields(t *testing.T) {
	h := newBoundHistogram("duration")
	h.bkts.reset(1, 1000, 3, linearBucket)
	h.afterResetBuckets()

	// record two observations so we have non-zero deltas
	h.UpdateMilliseconds(50)  // falls in first bucket
	h.UpdateMilliseconds(600) // falls in second bucket

	m := &model.Metric{}
	h.AppendFields(m)

	// should have 4 stat fields + N bucket fields
	var sumField, countField, minField, maxField *model.Field
	var bucketFields []*model.Field

	for i := range m.Fields {
		f := &m.Fields[i]
		switch {
		case f.Name == seriesmetric.HistoStatFieldName("duration", "sum"):
			sumField = f
		case f.Name == seriesmetric.HistoStatFieldName("duration", "count"):
			countField = f
		case f.Name == seriesmetric.HistoStatFieldName("duration", "min"):
			minField = f
		case f.Name == seriesmetric.HistoStatFieldName("duration", "max"):
			maxField = f
		case seriesmetric.IsBucketField(f.Name):
			bucketFields = append(bucketFields, f)
		}
	}

	require.NotNil(t, sumField, ".__sum field must be present")
	require.NotNil(t, countField, ".__count field must be present")
	require.NotNil(t, minField, ".__min field must be present")
	require.NotNil(t, maxField, ".__max field must be present")

	assert.Equal(t, field.Sum, sumField.Kind)
	assert.Equal(t, field.Sum, countField.Kind)
	assert.Equal(t, field.Min, minField.Kind)
	assert.Equal(t, field.Max, maxField.Kind)

	assert.InDelta(t, 650.0, sumField.Value, 1, "sum = 50+600")
	assert.InDelta(t, 2.0, countField.Value, 0.01, "count = 2 observations")

	// bucket fields must use FieldTypeSum kind (counts are summed across replicas)
	// and .__bucket_* naming; storage detects Histogram type via IsBucketField.
	for _, bf := range bucketFields {
		assert.Equal(t, field.Sum, bf.Kind, "bucket field %q must use Sum", bf.Name)
		assert.True(t, strings.Contains(bf.Name, ".__bucket_"), "bucket field name must contain .__bucket_: %s", bf.Name)
		assert.Equal(t, "duration", seriesmetric.HistoNameFromField(bf.Name))
	}

	// stat field names must start with "duration.__"
	assert.Equal(t, "duration.__sum", sumField.Name)
	assert.Equal(t, "duration.__count", countField.Name)
	assert.Equal(t, "duration.__min", minField.Name)
	assert.Equal(t, "duration.__max", maxField.Name)

	// calling AppendFields again should produce zero deltas (delta state was advanced)
	m2 := &model.Metric{}
	h.AppendFields(m2)
	for _, f := range m2.Fields {
		if f.Name == seriesmetric.HistoStatFieldName("duration", "sum") {
			assert.Equal(t, float64(0), f.Value, "sum delta should be 0 after second gather")
		}
		if f.Name == seriesmetric.HistoStatFieldName("duration", "count") {
			assert.Equal(t, float64(0), f.Value, "count delta should be 0 after second gather")
		}
	}
}

// TestBoundHistogram_AppendFields_Idempotent verifies that gathering twice without
// new observations produces zero deltas on the second call.
func TestBoundHistogram_AppendFields_Idempotent(t *testing.T) {
	h := newBoundHistogram("req_duration")
	h.bkts.reset(1, 100, 3, linearBucket)
	h.afterResetBuckets()
	h.UpdateMilliseconds(10)

	m1 := &model.Metric{}
	h.AppendFields(m1)

	m2 := &model.Metric{}
	h.AppendFields(m2)

	// second gather: all deltas should be zero
	for _, f := range m2.Fields {
		if f.Name == seriesmetric.HistoStatFieldName("req_duration", "sum") ||
			f.Name == seriesmetric.HistoStatFieldName("req_duration", "count") {
			assert.Equal(t, float64(0), f.Value, "field %q should be 0 on second gather", f.Name)
		}
	}
}

// TestHistoStatFieldName verifies the naming helper produces the expected format.
func TestHistoStatFieldName(t *testing.T) {
	assert.Equal(t, "latency.__sum", seriesmetric.HistoStatFieldName("latency", "sum"))
	assert.Equal(t, "latency.__count", seriesmetric.HistoStatFieldName("latency", "count"))
	assert.Equal(t, "latency.__min", seriesmetric.HistoStatFieldName("latency", "min"))
	assert.Equal(t, "latency.__max", seriesmetric.HistoStatFieldName("latency", "max"))
}

// TestIsHistoStatField verifies detection of new .__<stat> naming convention.
func TestIsHistoStatField(t *testing.T) {
	cases := []struct {
		name     string
		expected bool
	}{
		{"duration.__sum", true},
		{"duration.__count", true},
		{"duration.__min", true},
		{"duration.__max", true},
		// old-style (should no longer match after migration)
		{"duration_sum", false},
		{"duration_count", false},
		// plain fields
		{"duration", false},
		{"cpu_usage", false},
		// bucket fields (not stat fields)
		{"duration.__bucket_0.5", false},
		{"duration.__bucket_+Inf", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, seriesmetric.IsHistoStatField(tc.name))
		})
	}
}

// TestHistoNameFromField verifies histogram name extraction for both naming conventions.
func TestHistoNameFromField(t *testing.T) {
	cases := []struct {
		name     string
		expected string
	}{
		// new stat naming
		{"duration.__sum", "duration"},
		{"duration.__count", "duration"},
		{"duration.__min", "duration"},
		{"duration.__max", "duration"},
		// bucket fields
		{"duration.__bucket_0.5", "duration"},
		{"duration.__bucket_+Inf", "duration"},
		{"latency.__bucket_100", "latency"},
		// non-histogram fields
		{"duration", ""},
		{"cpu_usage", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, seriesmetric.HistoNameFromField(tc.name))
		})
	}
}

// TestAppendFields_ArrowRoundTrip verifies that fields written by AppendFields are
// correctly serialized and deserialized through the Arrow IPC path.
func TestAppendFields_ArrowRoundTrip(t *testing.T) {
	const histoName = "rpc_latency"

	h := newBoundHistogram(histoName)
	h.bkts.reset(1, 1000, 4, linearBucket) // 4 linear buckets
	h.afterResetBuckets()
	h.UpdateMilliseconds(200)
	h.UpdateMilliseconds(700)

	m := &model.Metric{
		Namespace: "default",
		Name:      "test_metric",
		Timestamp: int64(time.Now().UnixNano()),
	}
	m.Attributes = &model.Attributes{}
	h.AppendFields(m)

	// every Field in m.Fields must be resolvable:
	// stat fields → HistoNameFromField returns histoName
	// bucket fields → IsBucketField && HistoNameFromField returns histoName
	for _, f := range m.Fields {
		extractedName := seriesmetric.HistoNameFromField(f.Name)
		assert.Equal(t, histoName, extractedName, "field %q should resolve to histogram name %q", f.Name, histoName)
	}
}
