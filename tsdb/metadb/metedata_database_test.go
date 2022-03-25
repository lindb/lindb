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
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

func TestMetadataDatabase_New(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend

		ctrl.Finish()
	}()

	// test: new success
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// test: can't re-open
	db1, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.Error(t, err)
	assert.Nil(t, db1)

	// close db
	err = db.Close()
	assert.NoError(t, err)

	// test: create backend err
	createMetadataBackendFn = func(parent string) (MetadataBackend, error) {
		return nil, fmt.Errorf("err")
	}
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestMetadataDatabase_SuggestNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockMetadataBackend(ctrl)
	db := &metadataDatabase{
		backend: mockBackend,
	}
	mockBackend.EXPECT().suggestNamespace(gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	values, err := db.SuggestNamespace("ns", 10)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, values)
}

func TestMetadataDatabase_SuggestMetricName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockMetadataBackend(ctrl)
	db := &metadataDatabase{
		backend: mockBackend,
	}
	mockBackend.EXPECT().suggestMetricName(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	values, err := db.SuggestMetrics("ns", "pp", 10)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, values)
}

func TestMetadataDatabase_GetMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend
		ctrl.Finish()
	}()
	backend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (MetadataBackend, error) {
		return backend, nil
	}

	db := newMockMetadataDatabase(t, t.TempDir())
	db2 := db.(*metadataDatabase)
	db2.rwMux.Lock()
	db2.metrics[metric.JoinNamespaceMetric("ns-1", "name2")] = newMetricMetadata(metric.ID(2))
	db2.rwMux.Unlock()

	backend.EXPECT().getMetricID("ns-1", "name1").Return(metric.EmptyMetricID, fmt.Errorf("err"))
	metricID, err := db.GetMetricID("ns-1", "name1")
	assert.Error(t, err)
	assert.Equal(t, metric.EmptyMetricID, metricID)

	metricID, err = db.GetMetricID("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, metric.ID(2), metricID)
}

func TestMetadataDatabase_GetAllTagKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend
		ctrl.Finish()
	}()
	backend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (MetadataBackend, error) {
		return backend, nil
	}
	tags := tag.Metas{{ID: 1, Key: "key1"}, {ID: 2, Key: "key2"}}
	db := newMockMetadataDatabase(t, t.TempDir())

	cases := []struct {
		name       string
		metricName string
		prepare    func()
		out        struct {
			tags tag.Metas
			err  error
		}
	}{
		{
			name:       "get tag keys from memory cache",
			metricName: "name2",
			prepare: func() {
				db2 := db.(*metadataDatabase)
				db2.rwMux.Lock()
				metricMeta := newMetricMetadata(metric.ID(2))
				metricMeta.initialize(nil, 0, tags)
				db2.metrics[metric.JoinNamespaceMetric("ns-1", "name2")] = metricMeta
				db2.rwMux.Unlock()
			},
			out: struct {
				tags tag.Metas
				err  error
			}{tags: tags, err: nil},
		},
		{
			name:       "get metric id failure",
			metricName: "metric-name",
			prepare: func() {
				backend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(metric.EmptyMetricID, fmt.Errorf("err"))
			},
			out: struct {
				tags tag.Metas
				err  error
			}{tags: nil, err: fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, "metric-name")},
		},
		{
			name:       "get tag keys from backend storage",
			metricName: "metric-name",
			prepare: func() {
				id := metric.ID(3)
				backend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(id, nil)
				backend.EXPECT().getAllTagKeys(id).Return(tags, nil)
			},
			out: struct {
				tags tag.Metas
				err  error
			}{tags: tags, err: nil},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			tags, err := db.GetAllTagKeys("ns-1", tt.metricName)

			assert.Equal(t, tt.out.tags, tags)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataDatabase_GetTagKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	tags := tag.Metas{{ID: 1, Key: "key1"}, {ID: 2, Key: "key2"}}
	db := newMockMetadataDatabase(t, t.TempDir())
	db2 := db.(*metadataDatabase)
	db2.rwMux.Lock()
	metricMeta := newMetricMetadata(metric.ID(2))
	metricMeta.initialize(nil, 0, tags)
	db2.metrics[metric.JoinNamespaceMetric("ns-1", "name2")] = metricMeta
	db2.rwMux.Unlock()
	cases := []struct {
		name       string
		key        string
		metricName string
		prepare    func()
		out        struct {
			tagKeyID tag.KeyID
			err      error
		}
	}{
		{
			name:       "get all tag keys failure",
			metricName: "name",
			prepare: func() {
				mockBackend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(metric.EmptyMetricID, fmt.Errorf("err"))
			},
			out: struct {
				tagKeyID tag.KeyID
				err      error
			}{
				tagKeyID: tag.EmptyTagKeyID,
				err:      fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, "name"),
			},
		},
		{
			name:       "tag key not found",
			metricName: "name2",
			key:        "key3",
			out: struct {
				tagKeyID tag.KeyID
				err      error
			}{
				tagKeyID: tag.EmptyTagKeyID,
				err:      constants.ErrTagKeyIDNotFound,
			},
		},
		{
			name:       "get tag key successfully",
			metricName: "name2",
			key:        "key2",
			out: struct {
				tagKeyID tag.KeyID
				err      error
			}{
				tagKeyID: tag.KeyID(2),
				err:      nil,
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			tagKeyID, err := db.GetTagKeyID("ns-1", tt.metricName, tt.key)

			assert.Equal(t, tt.out.tagKeyID, tagKeyID)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataDatabase_GetAllFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db := newMockMetadataDatabase(t, t.TempDir())
	fields := field.Metas{
		{ID: 1, Type: field.SumField, Name: "sum"},
		{ID: 2, Type: field.HistogramField, Name: "histogram"},
	}
	cases := []struct {
		name       string
		metricName string
		prepare    func()
		out        struct {
			fields field.Metas
			err    error
		}
	}{
		{
			name:       "get fields from memory cache",
			metricName: "cache",
			prepare: func() {
				db2 := db.(*metadataDatabase)
				db2.rwMux.Lock()
				metricMeta := newMetricMetadata(metric.ID(2))
				metricMeta.initialize(fields, 0, nil)
				db2.metrics[metric.JoinNamespaceMetric("ns-1", "cache")] = metricMeta
				db2.rwMux.Unlock()
			},
			out: struct {
				fields field.Metas
				err    error
			}{fields: fields, err: nil},
		},
		{
			name:       "get metric id failure",
			metricName: "metric-name",
			prepare: func() {
				mockBackend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(metric.EmptyMetricID, fmt.Errorf("err"))
			},
			out: struct {
				fields field.Metas
				err    error
			}{fields: nil, err: fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, "metric-name")},
		},
		{
			name:       "get fields from backend storage",
			metricName: "metric-name",
			prepare: func() {
				id := metric.ID(3)
				mockBackend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(id, nil)
				mockBackend.EXPECT().getAllFields(id).Return(fields, field.ID(3), nil)
			},
			out: struct {
				fields field.Metas
				err    error
			}{fields: fields, err: nil},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			fields, err := db.GetAllFields("ns-1", tt.metricName)
			assert.Equal(t, tt.out.fields, fields)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataDatabase_GetField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db := newMockMetadataDatabase(t, t.TempDir())
	fields := field.Metas{
		{ID: 1, Type: field.SumField, Name: "sum"},
		{ID: 2, Type: field.HistogramField, Name: "histogram"},
	}
	db2 := db.(*metadataDatabase)
	db2.rwMux.Lock()
	metricMeta := newMetricMetadata(metric.ID(2))
	metricMeta.initialize(fields, 0, nil)
	db2.metrics[metric.JoinNamespaceMetric("ns-1", "cache")] = metricMeta
	db2.rwMux.Unlock()

	cases := []struct {
		name       string
		metricName string
		fieldName  field.Name
		prepare    func()
		out        struct {
			f   field.Meta
			err error
		}
	}{
		{
			name:       "get metric id failure",
			metricName: "metric-name",
			prepare: func() {
				mockBackend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(metric.EmptyMetricID, fmt.Errorf("err"))
			},
			out: struct {
				f   field.Meta
				err error
			}{f: field.Meta{}, err: fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, "metric-name")},
		},
		{
			name:       "field not found",
			metricName: "cache",
			fieldName:  field.Name("max"),
			out: struct {
				f   field.Meta
				err error
			}{f: field.Meta{}, err: constants.ErrNotFound},
		},
		{
			name:       "get field successfully",
			metricName: "cache",
			fieldName:  field.Name("sum"),
			out: struct {
				f   field.Meta
				err error
			}{f: field.Meta{ID: 1, Type: field.SumField, Name: "sum"}, err: nil},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			f, err := db.GetField("ns-1", tt.metricName, tt.fieldName)
			assert.Equal(t, tt.out.f, f)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataDatabase_GetAllHistogramFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db := newMockMetadataDatabase(t, t.TempDir())
	fields := field.Metas{
		{ID: 1, Type: field.SumField, Name: "sum"},
		{ID: 2, Type: field.HistogramField, Name: "histogram"},
	}
	db2 := db.(*metadataDatabase)
	db2.rwMux.Lock()
	metricMeta := newMetricMetadata(metric.ID(2))
	metricMeta.initialize(fields, 0, nil)
	db2.metrics[metric.JoinNamespaceMetric("ns-1", "cache")] = metricMeta
	db2.rwMux.Unlock()

	cases := []struct {
		name       string
		metricName string
		fieldName  field.Name
		prepare    func()
		out        struct {
			f   field.Metas
			err error
		}
	}{
		{
			name:       "get metric id failure",
			metricName: "metric-name",
			prepare: func() {
				mockBackend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(metric.EmptyMetricID, fmt.Errorf("err"))
			},
			out: struct {
				f   field.Metas
				err error
			}{f: nil, err: fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, "metric-name")},
		},
		{
			name:       "get histogram field successfully",
			metricName: "cache",
			fieldName:  field.Name("sum"),
			out: struct {
				f   field.Metas
				err error
			}{f: field.Metas{{ID: 2, Type: field.HistogramField, Name: "histogram"}}, err: nil},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			f, err := db.GetAllHistogramFields("ns-1", tt.metricName)
			assert.Equal(t, tt.out.f, f)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataDatabase_GenMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db := newMockMetadataDatabase(t, t.TempDir())
	cases := []struct {
		name       string
		metricName string
		prepare    func()
		out        struct {
			id  metric.ID
			err error
		}
	}{
		{
			name:       "get metric id from memory cache",
			metricName: "cache",
			prepare: func() {
				db2 := db.(*metadataDatabase)
				db2.rwMux.Lock()
				metricMeta := newMetricMetadata(metric.ID(2))
				db2.metrics[metric.JoinNamespaceMetric("ns-1", "cache")] = metricMeta
				db2.rwMux.Unlock()
			},
			out: struct {
				id  metric.ID
				err error
			}{id: metric.ID(2), err: nil},
		},
		{
			name:       "get metric metadata failure",
			metricName: "name",
			prepare: func() {
				mockBackend.EXPECT().getOrCreateMetricMetadata("ns-1", "name").Return(nil, fmt.Errorf("err"))
			},
			out: struct {
				id  metric.ID
				err error
			}{id: metric.EmptyMetricID, err: fmt.Errorf("err")},
		},
		{
			name:       "get metric metadata from backend storage",
			metricName: "name",
			prepare: func() {
				metadata := NewMockMetricMetadata(ctrl)
				metadata.EXPECT().getMetricID().Return(metric.ID(3))
				mockBackend.EXPECT().getOrCreateMetricMetadata("ns-1", "name").Return(metadata, nil)
			},
			out: struct {
				id  metric.ID
				err error
			}{id: metric.ID(3), err: nil},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			id, err := db.GenMetricID("ns-1", tt.metricName)
			assert.Equal(t, tt.out.id, id)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataDatabase_GenFieldID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend
		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	meta := NewMockMetricMetadata(ctrl)
	db := newMockMetadataDatabase(t, t.TempDir())
	db2 := db.(*metadataDatabase)
	db2.rwMux.Lock()
	db2.metrics[metric.JoinNamespaceMetric("ns-1", "cache")] = meta
	db2.rwMux.Unlock()
	cases := []struct {
		name       string
		metricName string
		f          field.Meta
		prepare    func()
		out        struct {
			id  field.ID
			err error
		}
	}{
		{
			name: "field type unspecified",
			f:    field.Meta{},
			out: struct {
				id  field.ID
				err error
			}{id: field.EmptyFieldID, err: series.ErrFieldTypeUnspecified},
		},
		{
			name:       "wrong field type",
			metricName: "cache",
			f:          field.Meta{Name: "sum", Type: field.MaxField},
			prepare: func() {
				meta.EXPECT().getField(field.Name("sum")).Return(field.Meta{Type: field.SumField}, true)
			},
			out: struct {
				id  field.ID
				err error
			}{id: field.EmptyFieldID, err: series.ErrWrongFieldType},
		},
		{
			name:       "get field from memory cache",
			metricName: "cache",
			f:          field.Meta{Name: "sum", Type: field.SumField},
			prepare: func() {
				meta.EXPECT().getField(field.Name("sum")).Return(field.Meta{Type: field.SumField, ID: 3}, true)
			},
			out: struct {
				id  field.ID
				err error
			}{id: field.ID(3), err: nil},
		},
		{
			name:       "create field failure",
			metricName: "cache",
			f:          field.Meta{Name: "sum", Type: field.SumField},
			prepare: func() {
				meta.EXPECT().getField(field.Name("sum")).Return(field.Meta{}, false)
				meta.EXPECT().createField(gomock.Any(), gomock.Any()).Return(field.Meta{}, fmt.Errorf("err"))
			},
			out: struct {
				id  field.ID
				err error
			}{id: field.EmptyFieldID, err: fmt.Errorf("err")},
		},
		{
			name:       "save field into backend storage failure",
			metricName: "cache",
			f:          field.Meta{Name: "sum", Type: field.SumField},
			prepare: func() {
				meta.EXPECT().getField(field.Name("sum")).Return(field.Meta{}, false)
				meta.EXPECT().createField(gomock.Any(), gomock.Any()).Return(field.Meta{}, nil)
				meta.EXPECT().getMetricID().Return(metric.ID(3))
				mockBackend.EXPECT().saveField(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			out: struct {
				id  field.ID
				err error
			}{id: field.EmptyFieldID, err: fmt.Errorf("err")},
		},
		{
			name:       "save field into backend storage successfully",
			metricName: "cache",
			f:          field.Meta{Name: "sum", Type: field.SumField},
			prepare: func() {
				meta.EXPECT().getField(field.Name("sum")).Return(field.Meta{}, false)
				meta.EXPECT().createField(gomock.Any(), gomock.Any()).Return(field.Meta{ID: 3}, nil)
				meta.EXPECT().getMetricID().Return(metric.ID(3))
				mockBackend.EXPECT().saveField(gomock.Any(), gomock.Any()).Return(nil)
			},
			out: struct {
				id  field.ID
				err error
			}{id: field.ID(3), err: nil},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			id, err := db.GenFieldID("ns-1", tt.metricName, tt.f.Name, tt.f.Type)
			assert.Equal(t, tt.out.id, id)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataDatabase_GenTagKeyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend
		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	meta := NewMockMetricMetadata(ctrl)
	db := newMockMetadataDatabase(t, t.TempDir())
	db2 := db.(*metadataDatabase)
	db2.rwMux.Lock()
	db2.metrics[metric.JoinNamespaceMetric("ns-1", "cache")] = meta
	db2.rwMux.Unlock()
	cases := []struct {
		name       string
		metricName string
		prepare    func()
		out        struct {
			id  tag.KeyID
			err error
		}
	}{
		{
			name:       "get tag memory cache",
			metricName: "cache",
			prepare: func() {
				meta.EXPECT().getTagKeyID(gomock.Any()).Return(tag.KeyID(3), true)
			},
			out: struct {
				id  tag.KeyID
				err error
			}{id: tag.KeyID(3), err: nil},
		},
		{
			name:       "create tag failure",
			metricName: "cache",
			prepare: func() {
				meta.EXPECT().getTagKeyID(gomock.Any()).Return(tag.EmptyTagKeyID, false)
				meta.EXPECT().checkTagKey(gomock.Any()).Return(fmt.Errorf("err"))
			},
			out: struct {
				id  tag.KeyID
				err error
			}{id: tag.EmptyTagKeyID, err: fmt.Errorf("err")},
		},
		{
			name:       "save tag into backend storage failure",
			metricName: "cache",
			prepare: func() {
				meta.EXPECT().getTagKeyID(gomock.Any()).Return(tag.EmptyTagKeyID, false)
				meta.EXPECT().checkTagKey(gomock.Any()).Return(nil)
				meta.EXPECT().getMetricID().Return(metric.ID(3))
				mockBackend.EXPECT().saveTagKey(gomock.Any(), gomock.Any()).Return(tag.EmptyTagKeyID, fmt.Errorf("err"))
			},
			out: struct {
				id  tag.KeyID
				err error
			}{id: tag.EmptyTagKeyID, err: fmt.Errorf("err")},
		},
		{
			name:       "save tag into backend storage successfully",
			metricName: "cache",
			prepare: func() {
				meta.EXPECT().getTagKeyID(gomock.Any()).Return(tag.EmptyTagKeyID, false)
				meta.EXPECT().checkTagKey(gomock.Any()).Return(nil)
				meta.EXPECT().getMetricID().Return(metric.ID(3))
				mockBackend.EXPECT().saveTagKey(gomock.Any(), gomock.Any()).Return(tag.KeyID(3), nil)
				meta.EXPECT().createTagKey(gomock.Any(), gomock.Any())
			},
			out: struct {
				id  tag.KeyID
				err error
			}{id: tag.KeyID(3), err: nil},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			id, err := db.GenTagKeyID("ns-1", tt.metricName, "key")
			assert.Equal(t, tt.out.id, id)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataDatabase_Close_Sync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackendFn = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	mockBackend.EXPECT().Close().Return(nil)
	mockBackend.EXPECT().sync().Return(nil)
	db := newMockMetadataDatabase(t, t.TempDir())
	assert.NoError(t, db.Close())
	assert.NoError(t, db.Sync())
}

func newMockMetadataDatabase(t *testing.T, dir string) MetadataDatabase {
	db, err := NewMetadataDatabase(context.TODO(), "test", dir)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	return db
}
