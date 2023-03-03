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

package metadb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

func TestMetricMetadata_createField(t *testing.T) {
	limits := models.NewDefaultLimits()
	cases := []struct {
		name    string
		prepare func(m MetricMetadata)
		out     struct {
			f   field.Meta
			err error
		}
	}{
		{
			name: "create field",
			out: struct {
				f   field.Meta
				err error
			}{
				f: field.Meta{
					ID:   1,
					Type: field.SumField,
					Name: "test",
				},
				err: nil,
			},
		},
		{

			name: "too many fields",
			prepare: func(m MetricMetadata) {
				m.initialize(nil, limits.MaxFieldsPerMetric, nil)
			},
			out: struct {
				f   field.Meta
				err error
			}{
				f:   field.Meta{},
				err: constants.ErrTooManyFields,
			},
		},
		{

			name: "disable too many fields",
			prepare: func(m MetricMetadata) {
				limits.MaxFieldsPerMetric = 0
				m.initialize(nil, limits.MaxFieldsPerMetric, nil)
			},
			out: struct {
				f   field.Meta
				err error
			}{
				f: field.Meta{
					ID:   1,
					Type: field.SumField,
					Name: "test",
				},
				err: nil,
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := newMetricMetadata(metric.ID(2))
			if tt.prepare != nil {
				tt.prepare(m)
			}
			mid := m.getMetricID()
			assert.Equal(t, metric.ID(2), mid)

			f, err := m.createField("test", field.SumField, limits)
			assert.Equal(t, tt.out.f, f)
			assert.Equal(t, tt.out.err, err)
			if err == nil {
				f1, ok := m.getField("test")
				assert.Equal(t, tt.out.f, f1)
				assert.True(t, ok)
			}
		})
	}
}

func TestMetricMetadata_getField(t *testing.T) {
	m := newMetricMetadata(metric.ID(2))
	assert.Empty(t, m.getAllFields())
	sum := field.Meta{
		ID:   field.ID(1),
		Type: field.SumField,
		Name: "sum",
	}
	m.initialize(field.Metas{
		sum,
		{
			ID:   field.ID(2),
			Type: field.HistogramField,
			Name: "histogram",
		},
	}, 0, nil)
	_, _ = m.createField("max", field.MinField, models.NewDefaultLimits())
	f, ok := m.getField("sum")
	assert.Equal(t, sum, f)
	assert.True(t, ok)
	fs := m.getAllFields()
	assert.Len(t, fs, 3)
	f, ok = m.getField("min")
	assert.Equal(t, field.Meta{}, f)
	assert.False(t, ok)
}

func TestMetricMetadata_createTagKey(t *testing.T) {
	m := newMetricMetadata(metric.ID(2))
	assert.Empty(t, m.getAllTagKeys())
	mid := m.getMetricID()
	assert.Equal(t, metric.ID(2), mid)

	m.createTagKey("key", 1)
	f1, ok := m.getTagKeyID("key")
	assert.Equal(t, tag.KeyID(1), f1)
	assert.True(t, ok)
}

func TestMetricMetadata_getTag(t *testing.T) {
	m := newMetricMetadata(metric.ID(2))
	tag1 := tag.Meta{ID: 2, Key: "key2"}
	m.initialize(nil, 0, tag.Metas{tag1})
	m.createTagKey("key3", 2)
	tags := m.getAllTagKeys()
	assert.Len(t, tags, 2)
	tag2, ok := m.getTagKeyID("key1")
	assert.False(t, ok)
	assert.Equal(t, tag.KeyID(0), tag2)
}

func TestMetricMetadata_checkTagKey(t *testing.T) {
	limits := models.NewDefaultLimits()
	m := newMetricMetadata(metric.ID(2))
	tag1 := tag.Meta{ID: 2, Key: "key2"}
	m.initialize(nil, 0, tag.Metas{tag1})
	assert.NoError(t, m.checkTagKey("", limits))
	limits.MaxTagsPerMetric = 1
	assert.Equal(t, constants.ErrTooManyTagKeys, m.checkTagKey("", limits))
	limits.MaxTagsPerMetric = 0
	assert.NoError(t, m.checkTagKey("", limits))
}
