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

package metric

import (
	"bytes"
	"math"
	"sort"
	"testing"

	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/strutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"

	"github.com/stretchr/testify/assert"
)

func makeProtoMetricV1(timestamp int64) *protoMetricsV1.Metric {
	var m protoMetricsV1.Metric
	m.Name = string(strutil.RandStringBytes(10))
	m.Namespace = "default-ns"
	m.Timestamp = timestamp

	var keyValues = tag.KeyValuesFromMap(map[string]string{
		"host": "test",
		"ip":   "1.1.1.1",
		"zone": "sh",
	})
	sort.Sort(keyValues)
	m.Tags = keyValues
	m.TagsHash = tag.XXHashOfKeyValues(keyValues)
	m.SimpleFields = []*protoMetricsV1.SimpleField{
		{
			Name:  "count1",
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
			Value: 100,
		},
		{Name: "gauge",
			Type:  protoMetricsV1.SimpleFieldType_GAUGE,
			Value: 10,
		}, {
			Name:  "min",
			Type:  protoMetricsV1.SimpleFieldType_Min,
			Value: 1,
		},
		{
			Name:  "max",
			Type:  protoMetricsV1.SimpleFieldType_Max,
			Value: 1000,
		}}
	m.CompoundField = &protoMetricsV1.CompoundField{
		Min:            1,
		Max:            1000,
		Count:          10,
		Sum:            10000,
		Values:         []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		ExplicitBounds: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, math.Inf(1)},
	}
	return &m
}

func Test_MarshalProtoMetricsV1List(t *testing.T) {
	var ml protoMetricsV1.MetricList
	for i := 0; i < 10; i++ {
		ml.Metrics = append(ml.Metrics, makeProtoMetricV1(fasttime.UnixMilliseconds()))
	}

	var buf bytes.Buffer
	for i := 0; i < 10; i++ {
		buf.Reset()
		size, err := MarshalProtoMetricsV1ListTo(ml, &buf)
		assert.Equal(t, size, len(buf.Bytes()))
		assert.NoError(t, err)
	}

	var br BatchRows
	br.UnmarshalRows(buf.Bytes())
	assert.Equal(t, 10, br.Len())
	br.UnmarshalRows(buf.Bytes())
	assert.Equal(t, 10, br.Len())

	row := br.Rows()[0]
	itr := row.NewKeyValueIterator()
	assert.True(t, itr.HasNext())
	assert.Equal(t, "host", string(itr.NextKey()))
	assert.Equal(t, "test", string(itr.NextValue()))

	assert.True(t, itr.HasNext())
	assert.Equal(t, "ip", string(itr.NextKey()))
	assert.Equal(t, "1.1.1.1", string(itr.NextValue()))

	assert.True(t, itr.HasNext())
	assert.Equal(t, "zone", string(itr.NextKey()))
	assert.Equal(t, "sh", string(itr.NextValue()))

	assert.Equal(t, tag.XXHashOfKeyValues(ml.Metrics[0].Tags), row.TagsHash())

}
