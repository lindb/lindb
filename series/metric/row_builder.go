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
	"io"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
)

var rowBuilderPool sync.Pool

type rowBuilder struct {
	flatBuilder *flatbuffers.Builder
	// offsets holding for builder flat buffer
	keys       []flatbuffers.UOffsetT
	values     []flatbuffers.UOffsetT
	kvs        []flatbuffers.UOffsetT
	fieldNames []flatbuffers.UOffsetT
	fields     []flatbuffers.UOffsetT
}

func (rb *rowBuilder) Reset() {
	rb.flatBuilder.Reset()
	rb.keys = rb.keys[:0]
	rb.values = rb.values[:0]
	rb.fieldNames = rb.fieldNames[:0]
	rb.kvs = rb.kvs[:0]
	rb.fields = rb.fields[:0]
}

func (rb *rowBuilder) MarshalProtoMetricV1(m *protoMetricsV1.Metric) {
	rb.Reset()

	// pre-allocate strings
	for i := 0; i < len(m.Tags); i++ {
		kv := m.Tags[i]
		rb.keys = append(rb.keys, rb.flatBuilder.CreateString(kv.Key))
		rb.values = append(rb.values, rb.flatBuilder.CreateString(kv.Value))
	}
	for i := 0; i < len(m.SimpleFields); i++ {
		rb.fieldNames = append(rb.fieldNames, rb.flatBuilder.CreateString(m.SimpleFields[i].Name))
	}
	// building key values vector
	for i := 0; i < len(rb.keys); i++ {
		flatMetricsV1.KeyValueStart(rb.flatBuilder)
		flatMetricsV1.KeyValueAddKey(rb.flatBuilder, rb.keys[i])
		flatMetricsV1.KeyValueAddValue(rb.flatBuilder, rb.values[i])
		rb.kvs = append(rb.kvs, flatMetricsV1.KeyValueEnd(rb.flatBuilder))
	}

	// building field names
	for i := 0; i < len(m.SimpleFields); i++ {
		sf := m.SimpleFields[i]
		flatMetricsV1.SimpleFieldStart(rb.flatBuilder)
		flatMetricsV1.SimpleFieldAddName(rb.flatBuilder, rb.fieldNames[i])
		switch sf.Type {
		case protoMetricsV1.SimpleFieldType_DELTA_SUM:
			flatMetricsV1.SimpleFieldAddType(rb.flatBuilder, flatMetricsV1.SimpleFieldTypeDeltaSum)
		case protoMetricsV1.SimpleFieldType_GAUGE:
			flatMetricsV1.SimpleFieldAddType(rb.flatBuilder, flatMetricsV1.SimpleFieldTypeGauge)
		case protoMetricsV1.SimpleFieldType_Max:
			flatMetricsV1.SimpleFieldAddType(rb.flatBuilder, flatMetricsV1.SimpleFieldTypeMax)
		case protoMetricsV1.SimpleFieldType_Min:
			flatMetricsV1.SimpleFieldAddType(rb.flatBuilder, flatMetricsV1.SimpleFieldTypeMin)
		default:
			flatMetricsV1.SimpleFieldAddType(rb.flatBuilder, flatMetricsV1.SimpleFieldTypeUnSpecified)
		}
		flatMetricsV1.SimpleFieldAddValue(rb.flatBuilder, sf.Value)
		rb.fields = append(rb.fields, flatMetricsV1.SimpleFieldEnd(rb.flatBuilder))
	}

	// serialize key values offsets
	flatMetricsV1.MetricStartKeyValuesVector(rb.flatBuilder, len(m.Tags))
	for i := len(rb.kvs) - 1; i >= 0; i-- {
		rb.flatBuilder.PrependUOffsetT(rb.kvs[i])
	}
	kvs := rb.flatBuilder.EndVector(len(rb.kvs))

	// serialize fields
	flatMetricsV1.MetricStartSimpleFieldsVector(rb.flatBuilder, len(rb.fields))
	for i := len(rb.fields) - 1; i >= 0; i-- {
		rb.flatBuilder.PrependUOffsetT(rb.fields[i])
	}
	fields := rb.flatBuilder.EndVector(len(rb.fields))

	var (
		compoundFieldBounds flatbuffers.UOffsetT
		compoundFieldValues flatbuffers.UOffsetT
		compoundField       flatbuffers.UOffsetT
	)

	if m.CompoundField == nil {
		goto Serialize
	}
	// serialize compound fields
	// add compound buckets explicit bounds
	flatMetricsV1.CompoundFieldStartValuesVector(rb.flatBuilder, len(m.CompoundField.Values))
	for i := len(m.CompoundField.Values) - 1; i >= 0; i-- {
		rb.flatBuilder.PrependFloat64(m.CompoundField.Values[i])
	}
	compoundFieldValues = rb.flatBuilder.EndVector(len(m.CompoundField.Values))
	// add compound buckets values
	flatMetricsV1.CompoundFieldStartExplicitBoundsVector(rb.flatBuilder, len(m.CompoundField.ExplicitBounds))
	for i := len(m.CompoundField.ExplicitBounds) - 1; i >= 0; i-- {
		rb.flatBuilder.PrependFloat64(m.CompoundField.ExplicitBounds[i])
	}
	compoundFieldBounds = rb.flatBuilder.EndVector(len(m.CompoundField.ExplicitBounds))

	// add count sum min max
	flatMetricsV1.CompoundFieldStart(rb.flatBuilder)
	flatMetricsV1.CompoundFieldAddCount(rb.flatBuilder, m.CompoundField.Count)
	flatMetricsV1.CompoundFieldAddSum(rb.flatBuilder, m.CompoundField.Sum)
	flatMetricsV1.CompoundFieldAddMin(rb.flatBuilder, m.CompoundField.Min)
	flatMetricsV1.CompoundFieldAddMax(rb.flatBuilder, m.CompoundField.Max)
	flatMetricsV1.CompoundFieldAddValues(rb.flatBuilder, compoundFieldValues)
	flatMetricsV1.CompoundFieldAddExplicitBounds(rb.flatBuilder, compoundFieldBounds)
	compoundField = flatMetricsV1.CompoundFieldEnd(rb.flatBuilder)

Serialize:
	// serialize metric
	metricName := rb.flatBuilder.CreateString(m.Name)
	namespace := rb.flatBuilder.CreateString(m.Namespace)
	flatMetricsV1.MetricStart(rb.flatBuilder)
	flatMetricsV1.MetricAddNamespace(rb.flatBuilder, namespace)
	flatMetricsV1.MetricAddName(rb.flatBuilder, metricName)
	flatMetricsV1.MetricAddTimestamp(rb.flatBuilder, m.Timestamp)
	flatMetricsV1.MetricAddKeyValues(rb.flatBuilder, kvs)
	flatMetricsV1.MetricAddHash(rb.flatBuilder, m.TagsHash)
	flatMetricsV1.MetricAddSimpleFields(rb.flatBuilder, fields)
	if compoundField != 0 {
		flatMetricsV1.MetricAddCompoundField(rb.flatBuilder, compoundField)
	}

	end := flatMetricsV1.MetricEnd(rb.flatBuilder)
	// size prefix encoding
	rb.flatBuilder.FinishSizePrefixed(end)
}

func newRowBuilder() *rowBuilder {
	return &rowBuilder{
		flatBuilder: flatbuffers.NewBuilder(1024 + 512),
		keys:        make([]flatbuffers.UOffsetT, 0, 32),
		values:      make([]flatbuffers.UOffsetT, 0, 32),
		fieldNames:  make([]flatbuffers.UOffsetT, 0, 32),
		kvs:         make([]flatbuffers.UOffsetT, 0, 32),
		fields:      make([]flatbuffers.UOffsetT, 0, 32),
	}
}

func getRowBuilder() (rb *rowBuilder, releaseFunc func(builder *rowBuilder)) {
	putBack := func(builder *rowBuilder) { rowBuilderPool.Put(builder) }
	item := rowBuilderPool.Get()
	if item != nil {
		builder := item.(*rowBuilder)
		builder.Reset()
		return builder, putBack
	}
	return newRowBuilder(), putBack
}

// MarshalProtoMetricsV1ListTo is a temporary marshal method for marshaling proto metrics with flat format
// todo: @codingcrush, use flat buffer on broker side
func MarshalProtoMetricsV1ListTo(pm protoMetricsV1.MetricList, writer io.Writer) (n int, err error) {
	builder, releaseFunc := getRowBuilder()
	defer releaseFunc(builder)

	for _, m := range pm.Metrics {
		builder.MarshalProtoMetricV1(m)
		size, err := writer.Write(builder.flatBuilder.FinishedBytes())
		n += size
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
