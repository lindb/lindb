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
	"math"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"
)

type BrokerRowProtoConverter struct {
	flatBuilder *flatbuffers.Builder
	// offsets holding for builder flat buffer
	keys       []flatbuffers.UOffsetT
	values     []flatbuffers.UOffsetT
	kvs        []flatbuffers.UOffsetT
	fieldNames []flatbuffers.UOffsetT
	fields     []flatbuffers.UOffsetT

	// ingestion meta info
	namespace    []byte
	enrichedTags tag.Tags
}

// Reset resets all data-structures
func (rc *BrokerRowProtoConverter) Reset() {
	rc.resetForNextConverter()
	rc.namespace = rc.namespace[:0]
	rc.enrichedTags = rc.enrichedTags[:0]
}

func (rc *BrokerRowProtoConverter) resetForNextConverter() {
	rc.flatBuilder.Reset()
	rc.keys = rc.keys[:0]
	rc.values = rc.values[:0]
	rc.fieldNames = rc.fieldNames[:0]
	rc.kvs = rc.kvs[:0]
	rc.fields = rc.fields[:0]
}

func (rc *BrokerRowProtoConverter) validateMetric(m *protoMetricsV1.Metric) error {
	if m == nil {
		return ErrMetricPBNilMetric
	}
	if len(m.Name) == 0 {
		return ErrMetricPBEmptyMetricName
	}
	m.Name = SanitizeMetricName(m.Name)
	// empty field
	if len(m.SimpleFields) == 0 && m.CompoundField == nil {
		return ErrMetricPBEmptyField
	}
	// re-set timestamp on zero
	if m.Timestamp == 0 {
		m.Timestamp = fasttime.UnixMilliseconds()
	}
	for i := 0; i < len(rc.enrichedTags); i++ {
		// todo: dedup @codingcrush
		m.Tags = append(m.Tags, &protoMetricsV1.KeyValue{
			Key:   string(rc.enrichedTags[i].Key),
			Value: string(rc.enrichedTags[i].Value),
		})
	}
	// replace namespace with enriched
	if len(rc.namespace) > 0 {
		m.Namespace = string(rc.namespace)
	}
	m.Namespace = SanitizeNamespace(m.Namespace)

	// validate empty tags
	if len(m.Tags) > 0 {
		for idx := range m.Tags {
			// nil tag
			if m.Tags[idx] == nil {
				return ErrMetricEmptyTagKeyValue
			}
			// empty key value
			if m.Tags[idx].Key == "" || m.Tags[idx].Value == "" {
				return ErrMetricEmptyTagKeyValue
			}
		}
	}

	// check simple fields
	for idx := range m.SimpleFields {
		// nil value
		if m.SimpleFields[idx] == nil {
			return ErrBadMetricPBFormat
		}
		// field-name empty
		if m.SimpleFields[idx].Name == "" {
			return ErrMetricEmptyFieldName
		}
		// check sanitize
		fieldName := strutil.String2ByteSlice(m.SimpleFields[idx].Name)
		if ShouldSanitizeFieldName(fieldName) {
			m.SimpleFields[idx].Name = string(SanitizeFieldName(fieldName))
		}
		// field type unspecified
		if m.SimpleFields[idx].Type == protoMetricsV1.SimpleFieldType_SIMPLE_UNSPECIFIED {
			return ErrBadMetricPBFormat
		}
		v := m.SimpleFields[idx].Value
		if math.IsNaN(v) {
			return ErrMetricNanField
		}
		if math.IsInf(v, 0) {
			return ErrMetricInfField
		}
	}
	// no more compound field
	if m.CompoundField == nil {
		return nil
	}
	// value length zero or length not match
	if len(m.CompoundField.Values) != len(m.CompoundField.ExplicitBounds) ||
		len(m.CompoundField.Values) <= 2 {
		return ErrBadMetricPBFormat
	}
	// ensure compound field value > 0
	if (m.CompoundField.Max < 0) ||
		m.CompoundField.Min < 0 ||
		m.CompoundField.Sum < 0 ||
		m.CompoundField.Count < 0 {
		return ErrBadMetricPBFormat
	}

	for idx := 0; idx < len(m.CompoundField.Values); idx++ {
		// ensure value > 0
		if m.CompoundField.Values[idx] < 0 || m.CompoundField.ExplicitBounds[idx] < 0 {
			return ErrBadMetricPBFormat
		}
		// ensure explicate bounds increase progressively
		if idx >= 1 && m.CompoundField.ExplicitBounds[idx] < m.CompoundField.ExplicitBounds[idx-1] {
			return ErrBadMetricPBFormat
		}
		// ensure last bound is +Inf
		if idx == len(m.CompoundField.ExplicitBounds)-1 && !math.IsInf(m.CompoundField.ExplicitBounds[idx], 1) {
			return ErrBadMetricPBFormat
		}
	}
	return nil
}

func (rc *BrokerRowProtoConverter) MarshalProtoMetricV1(m *protoMetricsV1.Metric) ([]byte, error) {
	rc.resetForNextConverter()

	if err := rc.validateMetric(m); err != nil {
		return nil, err
	}
	// pre-allocate strings
	for i := 0; i < len(m.Tags); i++ {
		kv := m.Tags[i]
		rc.keys = append(rc.keys, rc.flatBuilder.CreateString(kv.Key))
		rc.values = append(rc.values, rc.flatBuilder.CreateString(kv.Value))
	}

	// building key values vector
	for i := 0; i < len(rc.keys); i++ {
		flatMetricsV1.KeyValueStart(rc.flatBuilder)
		flatMetricsV1.KeyValueAddKey(rc.flatBuilder, rc.keys[i])
		flatMetricsV1.KeyValueAddValue(rc.flatBuilder, rc.values[i])
		rc.kvs = append(rc.kvs, flatMetricsV1.KeyValueEnd(rc.flatBuilder))
	}

	for i := 0; i < len(m.SimpleFields); i++ {
		rc.fieldNames = append(rc.fieldNames, rc.flatBuilder.CreateString(m.SimpleFields[i].Name))
	}

	// building field names
	for i := 0; i < len(m.SimpleFields); i++ {
		sf := m.SimpleFields[i]
		flatMetricsV1.SimpleFieldStart(rc.flatBuilder)
		flatMetricsV1.SimpleFieldAddName(rc.flatBuilder, rc.fieldNames[i])
		switch sf.Type {
		case protoMetricsV1.SimpleFieldType_DELTA_SUM:
			flatMetricsV1.SimpleFieldAddType(rc.flatBuilder, flatMetricsV1.SimpleFieldTypeDeltaSum)
		case protoMetricsV1.SimpleFieldType_GAUGE:
			flatMetricsV1.SimpleFieldAddType(rc.flatBuilder, flatMetricsV1.SimpleFieldTypeGauge)
		case protoMetricsV1.SimpleFieldType_Max:
			flatMetricsV1.SimpleFieldAddType(rc.flatBuilder, flatMetricsV1.SimpleFieldTypeMax)
		case protoMetricsV1.SimpleFieldType_Min:
			flatMetricsV1.SimpleFieldAddType(rc.flatBuilder, flatMetricsV1.SimpleFieldTypeMin)
		}
		flatMetricsV1.SimpleFieldAddValue(rc.flatBuilder, sf.Value)
		rc.fields = append(rc.fields, flatMetricsV1.SimpleFieldEnd(rc.flatBuilder))
	}

	// serialize key values offsets
	flatMetricsV1.MetricStartKeyValuesVector(rc.flatBuilder, len(m.Tags))
	for i := len(rc.kvs) - 1; i >= 0; i-- {
		rc.flatBuilder.PrependUOffsetT(rc.kvs[i])
	}
	kvs := rc.flatBuilder.EndVector(len(rc.kvs))

	// serialize fields
	flatMetricsV1.MetricStartSimpleFieldsVector(rc.flatBuilder, len(rc.fields))
	for i := len(rc.fields) - 1; i >= 0; i-- {
		rc.flatBuilder.PrependUOffsetT(rc.fields[i])
	}
	fields := rc.flatBuilder.EndVector(len(rc.fields))

	var (
		compoundFieldBounds flatbuffers.UOffsetT
		compoundFieldValues flatbuffers.UOffsetT
		compoundField       flatbuffers.UOffsetT
	)

	if m.CompoundField == nil {
		goto Serialize
	}
	// serialize compound fields
	// add compound buckets values
	flatMetricsV1.CompoundFieldStartValuesVector(rc.flatBuilder, len(m.CompoundField.Values))
	for i := len(m.CompoundField.Values) - 1; i >= 0; i-- {
		rc.flatBuilder.PrependFloat64(m.CompoundField.Values[i])
	}
	compoundFieldValues = rc.flatBuilder.EndVector(len(m.CompoundField.Values))
	// add compound buckets explicit bounds
	flatMetricsV1.CompoundFieldStartExplicitBoundsVector(rc.flatBuilder, len(m.CompoundField.ExplicitBounds))
	for i := len(m.CompoundField.ExplicitBounds) - 1; i >= 0; i-- {
		rc.flatBuilder.PrependFloat64(m.CompoundField.ExplicitBounds[i])
	}
	compoundFieldBounds = rc.flatBuilder.EndVector(len(m.CompoundField.ExplicitBounds))

	// add count sum min max
	flatMetricsV1.CompoundFieldStart(rc.flatBuilder)
	flatMetricsV1.CompoundFieldAddCount(rc.flatBuilder, m.CompoundField.Count)
	flatMetricsV1.CompoundFieldAddSum(rc.flatBuilder, m.CompoundField.Sum)
	flatMetricsV1.CompoundFieldAddMin(rc.flatBuilder, m.CompoundField.Min)
	flatMetricsV1.CompoundFieldAddMax(rc.flatBuilder, m.CompoundField.Max)
	flatMetricsV1.CompoundFieldAddValues(rc.flatBuilder, compoundFieldValues)
	flatMetricsV1.CompoundFieldAddExplicitBounds(rc.flatBuilder, compoundFieldBounds)
	compoundField = flatMetricsV1.CompoundFieldEnd(rc.flatBuilder)

Serialize:
	// serialize metric
	metricName := rc.flatBuilder.CreateString(m.Name)
	namespace := rc.flatBuilder.CreateString(m.Namespace)
	flatMetricsV1.MetricStart(rc.flatBuilder)
	flatMetricsV1.MetricAddNamespace(rc.flatBuilder, namespace)
	flatMetricsV1.MetricAddName(rc.flatBuilder, metricName)
	flatMetricsV1.MetricAddTimestamp(rc.flatBuilder, m.Timestamp)
	flatMetricsV1.MetricAddKeyValues(rc.flatBuilder, kvs)
	// sort and computing tags hash
	flatMetricsV1.MetricAddHash(rc.flatBuilder, tag.XXHashOfKeyValues(m.Tags))
	flatMetricsV1.MetricAddSimpleFields(rc.flatBuilder, fields)
	if compoundField != 0 {
		flatMetricsV1.MetricAddCompoundField(rc.flatBuilder, compoundField)
	}

	end := flatMetricsV1.MetricEnd(rc.flatBuilder)
	// size prefix encoding
	rc.flatBuilder.FinishSizePrefixed(end)

	return rc.flatBuilder.FinishedBytes(), nil
}

func (rc *BrokerRowProtoConverter) MarshalProtoMetricV1To(m *protoMetricsV1.Metric, writer io.Writer) (n int, err error) {
	block, err := rc.MarshalProtoMetricV1(m)
	if err != nil {
		return 0, err
	}
	return writer.Write(block)
}

// ConvertTo converts the proto metric into BrokerRow
func (rc *BrokerRowProtoConverter) ConvertTo(m *protoMetricsV1.Metric, row *BrokerRow) error {
	block, err := rc.MarshalProtoMetricV1(m)
	if err != nil {
		return err
	}
	row.FromBlock(block)
	return nil
}

func (rc *BrokerRowProtoConverter) MarshalProtoMetricListV1To(ml protoMetricsV1.MetricList, writer io.Writer) (n int, err error) {
	for _, m := range ml.Metrics {
		size, err := rc.MarshalProtoMetricV1To(m, writer)
		n += size
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

var rowConverterPool sync.Pool

// NewProtoConverter returns a converter for converting proto metrics into flat metric
func NewProtoConverter() *BrokerRowProtoConverter {
	return &BrokerRowProtoConverter{
		flatBuilder: flatbuffers.NewBuilder(1024 + 512),
		keys:        make([]flatbuffers.UOffsetT, 0, 32),
		values:      make([]flatbuffers.UOffsetT, 0, 32),
		fieldNames:  make([]flatbuffers.UOffsetT, 0, 32),
		kvs:         make([]flatbuffers.UOffsetT, 0, 32),
		fields:      make([]flatbuffers.UOffsetT, 0, 32),
	}
}

// NewBrokerRowProtoConverter returns a new converter for converting proto metrics into flat metrics.
// namespace and enrichedTags will also be bound to the metric
func NewBrokerRowProtoConverter(
	namespace []byte,
	enrichedTags tag.Tags,
) (
	cvt *BrokerRowProtoConverter,
	releaseFunc func(cvt *BrokerRowProtoConverter),
) {
	releaseFunc = func(cvt *BrokerRowProtoConverter) { rowConverterPool.Put(cvt) }
	item := rowConverterPool.Get()
	if item == nil {
		cvt = NewProtoConverter()
	} else {
		cvt = item.(*BrokerRowProtoConverter)
	}
	cvt.Reset()
	cvt.namespace = namespace
	cvt.enrichedTags = enrichedTags
	return cvt, releaseFunc
}
