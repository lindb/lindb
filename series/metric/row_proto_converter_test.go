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
	"strconv"
	"sync"
	"testing"

	"github.com/lindb/common/pkg/fasttime"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/tag"
)

func makeProtoMetricV1(timestamp int64) *protoMetricsV1.Metric {
	var m protoMetricsV1.Metric
	m.Name = string(strutil.RandStringBytes(10))
	m.Namespace = "default-ns"
	m.Timestamp = timestamp

	var keyValues = tag.KeyValuesFromMap(map[string]string{
		"host": strconv.FormatInt(timestamp, 10),
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
		{
			Name:  "last",
			Type:  protoMetricsV1.SimpleFieldType_LAST,
			Value: 10,
		},
		{
			Name:  "first",
			Type:  protoMetricsV1.SimpleFieldType_FIRST,
			Value: 10,
		},
		{
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

	converter := NewProtoConverter(models.NewDefaultLimits())

	var buf bytes.Buffer
	for i := 0; i < 10; i++ {
		buf.Reset()
		size, err := converter.MarshalProtoMetricListV1To(ml, &buf)
		assert.Equal(t, size, len(buf.Bytes()))
		assert.NoError(t, err)
	}

	var br StorageBatchRows
	br.UnmarshalRows(buf.Bytes())
	assert.Equal(t, 10, br.Len())
	br.UnmarshalRows(buf.Bytes())
	assert.Equal(t, 10, br.Len())

	row := br.Rows()[0]

	itr := row.NewKeyValueIterator()
	assert.True(t, itr.HasNext())
	assert.Equal(t, "host", string(itr.NextKey()))
	assert.NotEmpty(t, "test", string(itr.NextValue()))

	assert.True(t, itr.HasNext())
	assert.Equal(t, "ip", string(itr.NextKey()))
	assert.Equal(t, "1.1.1.1", string(itr.NextValue()))

	assert.True(t, itr.HasNext())
	assert.Equal(t, "zone", string(itr.NextKey()))
	assert.Equal(t, "sh", string(itr.NextValue()))

	assert.Equal(t, tag.XXHashOfKeyValues(ml.Metrics[0].Tags), row.TagsHash())
}

func Test_BrokerRowProtoConverter_ValidateMetric(t *testing.T) {
	converter, releaseFunc := NewBrokerRowProtoConverter(
		[]byte("lindb-ns"), tag.Tags{
			tag.NewTag([]byte("a"), []byte("b")),
		}, models.NewDefaultLimits())
	defer releaseFunc(converter)

	// nil pb
	assert.Error(t, converter.validateMetric(nil))
	// empty name
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{}))
	// empty field
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
	}))
	// nil tag
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{
				Name:  "f1",
				Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				Value: 1,
			}},
		Tags: []*protoMetricsV1.KeyValue{nil, nil},
	}))
	// empty tag
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{
				Name:  "f1",
				Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				Value: 1,
			}},
		Tags: []*protoMetricsV1.KeyValue{{Key: "", Value: ""}},
	}))

	// nil simple field
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name:         "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{nil, nil},
	}))
	// empty simple field name
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name:         "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{{Name: ""}},
	}))
	// unspecified simple field type
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name:         "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{{Name: "f1"}},
	}))
	// NaN simple field value
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{
				Name:  "__bucket_1",
				Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				Value: math.NaN(),
			}},
	}))
	// Inf simple field value
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{
				Name:  "f1",
				Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				Value: math.Inf(1),
			}},
	}))
	// ok with none compound field
	assert.NoError(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{
				Name:  "__bucket_1",
				Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				Value: 1,
			}},
	}))

	// compound field values not match explicit-bounds
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		CompoundField: &protoMetricsV1.CompoundField{
			ExplicitBounds: []float64{1, 2, 3, 4, 5},
			Values:         []float64{1, 2, 3, 4},
		},
	}))
	// invalid mmsc
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		CompoundField: &protoMetricsV1.CompoundField{
			ExplicitBounds: []float64{1, 2, 3, 4, 5},
			Values:         []float64{1, 2, 3, 4, 5},
			Count:          -1,
		},
	}))
	// negative value
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		CompoundField: &protoMetricsV1.CompoundField{
			ExplicitBounds: []float64{-1, 2, 3, 4, 5},
			Values:         []float64{1, 2, 3, 4, 5},
		},
	}))
	// not increasing
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		CompoundField: &protoMetricsV1.CompoundField{
			ExplicitBounds: []float64{1, 2, 1, 4, math.Inf(1)},
			Values:         []float64{1, 2, 3, 4, 5},
		},
	}))
	// last explicit bound not inf
	assert.Error(t, converter.validateMetric(&protoMetricsV1.Metric{
		Name: "test-metric",
		CompoundField: &protoMetricsV1.CompoundField{
			ExplicitBounds: []float64{1, 2, 3, 4, 5},
			Values:         []float64{1, 2, 3, 4, 5},
		},
	}))
}

func Test_BrokerRowProtoConverter_MarshalProtoMetricV1(t *testing.T) {
	converter, releaseFunc := NewBrokerRowProtoConverter(
		[]byte("lindb-ns"), tag.Tags{
			tag.NewTag([]byte("a"), []byte("b")),
		}, models.NewDefaultLimits())
	defer releaseFunc(converter)

	data, err := converter.MarshalProtoMetricV1(nil)
	assert.Error(t, err)
	assert.Len(t, data, 0)

	var ml = protoMetricsV1.MetricList{
		Metrics: []*protoMetricsV1.Metric{{Name: ""}},
	}
	var buf bytes.Buffer
	var row BrokerRow
	_, err = converter.MarshalProtoMetricListV1To(ml, &buf)
	assert.Error(t, err)
	assert.Error(t, converter.ConvertTo(ml.Metrics[0], &row))

	// marshal ok
	m := &protoMetricsV1.Metric{
		Name: "test-metric",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{
				Name:  "__bucket_1",
				Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				Value: 1,
			}},
	}
	data, err = converter.MarshalProtoMetricV1(m)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	assert.NoError(t, converter.ConvertTo(m, &row))

	ml = protoMetricsV1.MetricList{
		Metrics: []*protoMetricsV1.Metric{m},
	}
	_, err = converter.MarshalProtoMetricListV1To(ml, &buf)
	assert.NoError(t, err)
}

func Test_BrokerRowProtoConverter_deDupTags(t *testing.T) {
	converter, releaseFunc := NewBrokerRowProtoConverter(
		nil, nil, models.NewDefaultLimits())
	defer releaseFunc(converter)

	m := &protoMetricsV1.Metric{
		Tags: []*protoMetricsV1.KeyValue{
			{Key: "a", Value: "1"},
			{Key: "b", Value: "2"},
			{Key: "a", Value: "3"},
			{Key: "b", Value: "4"},
			{Key: "c", Value: "5"},
		},
	}
	converter.deDupTags(m)
	assert.EqualValues(t, []*protoMetricsV1.KeyValue{
		{Key: "a", Value: "3"},
		{Key: "b", Value: "4"},
		{Key: "c", Value: "5"},
	}, m.Tags)

	m2 := &protoMetricsV1.Metric{
		Tags: []*protoMetricsV1.KeyValue{
			{Key: "a", Value: "1"},
		},
	}
	assert.EqualValues(t, []*protoMetricsV1.KeyValue{
		{Key: "a", Value: "1"},
	}, m2.Tags)
}

func TestNewProtoCoverter(t *testing.T) {
	defer func() {
		rowConverterPool = sync.Pool{}
	}()
	rowConverterPool = sync.Pool{
		New: func() any {
			return nil
		},
	}
	converter, releaseFunc := NewBrokerRowProtoConverter(
		nil, nil, models.NewDefaultLimits())
	releaseFunc(converter)
	rowConverterPool = sync.Pool{
		New: func() any {
			return NewProtoConverter(models.NewDefaultLimits())
		},
	}
	converter, releaseFunc = NewBrokerRowProtoConverter(
		nil, nil, models.NewDefaultLimits())
	releaseFunc(converter)
}

func TestProtoCoverter_Limits(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(limits *models.Limits)
		wantErr bool
		err     error
	}{
		{
			name: "metric name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxMetricNameLength = 1
			},
			wantErr: true,
			err:     constants.ErrMetricNameTooLong,
		},
		{
			name: "too many tags",
			prepare: func(limits *models.Limits) {
				limits.MaxTagsPerMetric = 1
			},
			wantErr: true,
			err:     constants.ErrTooManyTagKeys,
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
			name: "field name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxFieldNameLength = 1
			},
			wantErr: true,
			err:     constants.ErrFieldNameTooLong,
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
			name: "tag value too long",
			prepare: func(limits *models.Limits) {
				limits.MaxTagValueLength = 1
			},
			wantErr: true,
			err:     constants.ErrTagValueTooLong,
		},
		{
			name: "disable metric name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxMetricNameLength = 0
			},
		},
		{
			name: "disable too many tags",
			prepare: func(limits *models.Limits) {
				limits.MaxTagsPerMetric = 0
			},
		},
		{
			name: "disable too many fields",
			prepare: func(limits *models.Limits) {
				limits.MaxFieldsPerMetric = 0
			},
		},
		{
			name: "disable field name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxFieldNameLength = 0
			},
		},
		{
			name: "disable tag name too long",
			prepare: func(limits *models.Limits) {
				limits.MaxTagNameLength = 0
			},
		},
		{
			name: "disable tag value too long",
			prepare: func(limits *models.Limits) {
				limits.MaxTagValueLength = 0
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			limits := models.NewDefaultLimits()
			// marshal ok
			m := &protoMetricsV1.Metric{
				Name: "test-metric",
				Tags: []*protoMetricsV1.KeyValue{
					{Key: "host", Value: "host_name"},
					{Key: "host1", Value: "host_name"},
				},
				SimpleFields: []*protoMetricsV1.SimpleField{
					{
						Name:  "__bucket_1",
						Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
						Value: 1,
					},
					{
						Name:  "__bucket_2",
						Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
						Value: 1,
					},
				},
			}
			tt.prepare(limits)
			converter, releaseFunc := NewBrokerRowProtoConverter(
				nil, nil, limits)
			defer releaseFunc(converter)
			if tt.wantErr {
				assert.Equal(t, tt.err, converter.validateMetric(m))
			} else {
				assert.NoError(t, converter.validateMetric(m))
			}
		})
	}
}
