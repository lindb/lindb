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
	"testing"

	"github.com/stretchr/testify/assert"

	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/tag"
)

func Test_NewBrokerRowFlatDecoder(t *testing.T) {
	converter1 := NewProtoConverter()
	converter2 := NewProtoConverter()
	var buf bytes.Buffer

	now := timeutil.Now()
	data1, err := converter1.MarshalProtoMetricV1(&protoMetricsV1.Metric{
		Name:      "test",
		Timestamp: now,
		Tags: []*protoMetricsV1.KeyValue{
			{Key: "key", Value: "value"},
		},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "F1", Type: protoMetricsV1.SimpleFieldType_Min, Value: 1},
		},
	})
	assert.NoError(t, err)
	_, _ = buf.Write(data1)

	data2, err := converter2.MarshalProtoMetricV1(&protoMetricsV1.Metric{
		Name:      "test",
		Namespace: "ns",
		Timestamp: now,
		Tags: []*protoMetricsV1.KeyValue{
			{Key: "key", Value: "value"},
		},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "F1", Type: protoMetricsV1.SimpleFieldType_Min, Value: 1},
		},
		CompoundField: &protoMetricsV1.CompoundField{
			Min:            1,
			Max:            1,
			Sum:            1,
			Count:          1,
			ExplicitBounds: []float64{1, 2, 3, 4, 5, math.Inf(1)},
			Values:         []float64{0, 0, 0, 0, 0, 0},
		},
	})
	assert.NoError(t, err)
	_, _ = buf.Write(data2)

	decoder, releaseFunc := NewBrokerRowFlatDecoder(nil, nil, nil)
	assert.False(t, decoder.HasNext())
	releaseFunc(decoder)
	assert.Zero(t, decoder.ReadLen())

	reader := bytes.NewReader(buf.Bytes())
	decoder, releaseFunc = NewBrokerRowFlatDecoder(
		reader,
		[]byte("lindb-ns"),
		tag.Tags{
			tag.NewTag([]byte("a"), []byte("b")),
		})
	defer releaseFunc(decoder)

	var row BrokerRow
	assert.True(t, decoder.HasNext())
	assert.NoError(t, decoder.DecodeTo(&row))

	assert.True(t, decoder.HasNext())
	assert.NoError(t, decoder.DecodeTo(&row))

	assert.False(t, decoder.HasNext())
	assert.Error(t, decoder.DecodeTo(&row))
	assert.Equal(t, len(buf.Bytes()), decoder.ReadLen())
	metric := row.Metric()
	assert.Equal(t, now, metric.Timestamp())
}
