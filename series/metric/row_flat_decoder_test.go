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
	"sync"
	"testing"

	commontimeutil "github.com/lindb/common/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/tag"
)

func Test_NewBrokerRowFlatDecoder(t *testing.T) {
	converter1 := NewProtoConverter(models.NewDefaultLimits())
	converter2 := NewProtoConverter(models.NewDefaultLimits())
	var buf bytes.Buffer

	now := commontimeutil.Now()
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

	data2, err := converter2.MarshalProtoMetricV1(&protoMetricsV1.Metric{Name: "test",
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

	decoder, releaseFunc := NewBrokerRowFlatDecoder(nil, nil, nil, models.NewDefaultLimits())
	assert.False(t, decoder.HasNext())
	releaseFunc(decoder)
	assert.Zero(t, decoder.ReadLen())

	reader := bytes.NewReader(buf.Bytes())
	decoder, releaseFunc = NewBrokerRowFlatDecoder(
		reader,
		[]byte("lindb-ns"),
		tag.Tags{
			tag.NewTag([]byte("a"), []byte("b")),
		}, models.NewDefaultLimits())
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

func Test_NewBrokerRowFlatDecoder_pool(t *testing.T) {
	defer func() {
		brokerRowFlatDecoderPool = sync.Pool{}
	}()
	brokerRowFlatDecoderPool = sync.Pool{
		New: func() any {
			return nil
		},
	}
	decoder, releaseFunc := NewBrokerRowFlatDecoder(nil, nil, nil, models.NewDefaultLimits())
	assert.False(t, decoder.HasNext())
	releaseFunc(decoder)
	assert.Zero(t, decoder.ReadLen())

	brokerRowFlatDecoderPool = sync.Pool{
		New: func() any {
			return &BrokerRowFlatDecoder{}
		},
	}
	decoder, releaseFunc = NewBrokerRowFlatDecoder(nil, nil, nil, models.NewDefaultLimits())
	assert.False(t, decoder.HasNext())
	releaseFunc(decoder)
	assert.Zero(t, decoder.ReadLen())
}

func Test_BrokerRowFlatDecoder_Decode_Fail(t *testing.T) {
	decoder, _ := NewBrokerRowFlatDecoder(nil, nil,
		tag.Tags{
			tag.NewTag([]byte("a"), []byte("b")),
		}, models.NewDefaultLimits())
	var row BrokerRow
	assert.Error(t, decoder.DecodeTo(&row))

	cases := []struct {
		name    string
		prepare func(limits *models.Limits)
		wantErr bool
		err     error
	}{
		{
			name: "too many tags",
			prepare: func(limits *models.Limits) {
				limits.MaxTagsPerMetric = 1
			},
			wantErr: true,
			err:     constants.ErrTooManyTagKeys,
		},
		{
			name: "disable too many tags",
			prepare: func(limits *models.Limits) {
				limits.MaxTagsPerMetric = 0
			},
		},
		{
			name: "too many fields",
			prepare: func(limits *models.Limits) {
				limits.MaxFieldsPerMetric = 1
			},
			wantErr: true,
			err:     constants.ErrTooManyFields,
		},
		{
			name: "disable too many fields",
			prepare: func(limits *models.Limits) {
				limits.MaxFieldsPerMetric = 0
			},
		},
		{
			name: "tag name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxTagNameLength = 1
			},
			wantErr: true,
			err:     constants.ErrTagKeyTooLong,
		},
		{
			name: "disable tag name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxTagNameLength = 0
			},
		},
		{
			name: "tag value too long",
			prepare: func(limits *models.Limits) {
				limits.MaxTagValueLength = 1
			},
			wantErr: true,
			err:     constants.ErrTagValueTooLong,
		},
		{
			name: "disable tag value too long",
			prepare: func(limits *models.Limits) {
				limits.MaxTagValueLength = 0
			},
		},
		{
			name: "field name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxFieldNameLength = 1
			},
			wantErr: true,
			err:     constants.ErrFieldNameTooLong,
		},
		{
			name: "disable field name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxFieldNameLength = 0
			},
		},
		{
			name: "metric name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxMetricNameLength = 1
			},
			wantErr: true,
			err:     constants.ErrMetricNameTooLong,
		},
		{
			name: "disable metric name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxMetricNameLength = 0
			},
		},
		{
			name: "namespace too long",
			prepare: func(limits *models.Limits) {
				limits.MaxNamespaceLength = 1
			},
			wantErr: true,
			err:     constants.ErrNamespaceTooLong,
		},
		{
			name: "disable namespace too long",
			prepare: func(limits *models.Limits) {
				limits.MaxNamespaceLength = 0
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			limits := models.NewDefaultLimits()
			tt.prepare(limits)
			decoder = mockDecoder(limits)
			assert.True(t, decoder.HasNext())
			err := decoder.DecodeTo(&row)
			if tt.wantErr && err != nil {
				assert.Equal(t, tt.err, err)
			}
		})
	}
}

func mockDecoder(limits *models.Limits) *BrokerRowFlatDecoder {
	converter2 := NewProtoConverter(models.NewDefaultLimits())
	data2, err := converter2.MarshalProtoMetricV1(&protoMetricsV1.Metric{
		Name:      "test",
		Namespace: "ns",
		Tags: []*protoMetricsV1.KeyValue{
			{Key: "key", Value: "value"},
		},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "F1", Type: protoMetricsV1.SimpleFieldType_Min, Value: 1},
			{Name: "F2", Type: protoMetricsV1.SimpleFieldType_Min, Value: 1},
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
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	_, _ = buf.Write(data2)

	reader := bytes.NewReader(buf.Bytes())
	decoder, _ := NewBrokerRowFlatDecoder(
		reader,
		nil,
		tag.Tags{
			tag.NewTag([]byte("a"), []byte("b")),
		}, limits)
	return decoder
}
