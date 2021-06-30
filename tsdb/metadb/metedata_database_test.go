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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/wal"
)

func TestMetadataDatabase_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		createMetaWAL = wal.NewMetricMetaWAL
		_ = fileutil.RemoveDir(testPath)

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

	// test: create wal err
	mockBackend := NewMockMetadataBackend(ctrl)
	mockBackend.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes()
	createMetadataBackend = func(parent string) (MetadataBackend, error) {
		return mockBackend, nil
	}
	createMetaWAL = func(path string) (wal.MetricMetaWAL, error) {
		return nil, fmt.Errorf("err")
	}
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.Error(t, err)
	assert.Nil(t, db)

	// test: recovery err
	mockWAL := wal.NewMockMetricMetaWAL(ctrl)
	createMetaWAL = func(path string) (wal.MetricMetaWAL, error) {
		return mockWAL, nil
	}
	mockWAL.EXPECT().Recovery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	mockWAL.EXPECT().NeedRecovery().Return(true)
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestMetadataDatabase_SuggestNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	mockBackend.EXPECT().suggestNamespace(gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	values, err := db.SuggestNamespace("ns", 10)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, values)

	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_SuggestMetricName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	mockBackend.EXPECT().suggestMetricName(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	values, err := db.SuggestMetrics("ns", "pp", 10)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, values)

	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GetMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	gomock.InOrder(
		mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(nil, constants.ErrNotFound),
		mockBackend.EXPECT().genMetricID().Return(uint32(1)),
	)
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	metricID, err = db.GetMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	metricID, err = db.GetMetricID("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), metricID)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GetTagKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	meta := NewMockMetricMetadata(ctrl)
	mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(meta, nil)
	meta.EXPECT().getMetricID().Return(uint32(1))
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	// case 1: from memory
	meta.EXPECT().getTagKeyID("tag-key").Return(uint32(10), true)
	tagKeyID, err := db.GetTagKeyID("ns-1", "name1", "tag-key")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), tagKeyID)

	// case 2: from memory not exist
	meta.EXPECT().getTagKeyID("tag-key").Return(uint32(10), false)
	tagKeyID, err = db.GetTagKeyID("ns-1", "name1", "tag-key")
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Equal(t, uint32(0), tagKeyID)

	// case 4: backend, metric not exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(0), constants.ErrNotFound)
	tagKeyID, err = db.GetTagKeyID("ns-1", "name2", "tag-key")
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Equal(t, uint32(0), tagKeyID)

	// case 4: backend exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	mockBackend.EXPECT().getTagKeyID(uint32(10), "tag-key").Return(uint32(20), nil)
	tagKeyID, err = db.GetTagKeyID("ns-1", "name2", "tag-key")
	assert.NoError(t, err)
	assert.Equal(t, uint32(20), tagKeyID)

	// case 5: all tag keys from memory
	meta.EXPECT().getAllTagKeys().Return([]tag.Meta{{ID: 10, Key: "tag-key"}})
	tagKeys, err := db.GetAllTagKeys("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, []tag.Meta{{ID: 10, Key: "tag-key"}}, tagKeys)

	// case 6: backend, metric not exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(0), constants.ErrNotFound)
	tagKeys, err = db.GetAllTagKeys("ns-1", "name2")
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Nil(t, tagKeys)

	// case 7: backend, tag keys exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	mockBackend.EXPECT().getAllTagKeys(uint32(10)).Return([]tag.Meta{{ID: 10, Key: "tag-key"}}, nil)
	tagKeys, err = db.GetAllTagKeys("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, []tag.Meta{{ID: 10, Key: "tag-key"}}, tagKeys)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_SuggestTagKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	mockBackend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(uint32(10), nil).AnyTimes()
	// case 1: suggest tag keys
	mockBackend.EXPECT().getAllTagKeys(gomock.Any()).Return([]tag.Meta{{ID: 10, Key: "tag-key"}}, nil)
	tagKeys, err := db.SuggestTagKeys("ns-1", "name1", "tag", 100)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag-key"}, tagKeys)
	// case 2: get tag values err
	mockBackend.EXPECT().getAllTagKeys(gomock.Any()).Return([]tag.Meta{{ID: 10, Key: "tag-key"}}, fmt.Errorf("err"))
	tagKeys, err = db.SuggestTagKeys("ns-1", "name1", "tag", 100)
	assert.Error(t, err)
	assert.Nil(t, tagKeys)
	// case 2: get tag values limit
	mockBackend.EXPECT().getAllTagKeys(gomock.Any()).Return([]tag.Meta{{ID: 10, Key: "tag-key"}, {ID: 10, Key: "tag-key1"}}, nil)
	tagKeys, err = db.SuggestTagKeys("ns-1", "name1", "tag", 1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag-key"}, tagKeys)
}

func TestMetadataDatabase_GetField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	meta := NewMockMetricMetadata(ctrl)
	mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(meta, nil)
	meta.EXPECT().getMetricID().Return(uint32(1))
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	// case 1: from memory
	meta.EXPECT().getField(field.Name("f1")).Return(field.Meta{ID: 19, Type: field.SumField}, true)
	f, err := db.GetField("ns-1", "name1", "f1")
	assert.NoError(t, err)
	assert.Equal(t, field.Meta{ID: 19, Type: field.SumField}, f)

	// case 2: from memory not exist
	meta.EXPECT().getField(field.Name("f1")).Return(field.Meta{}, false)
	f, err = db.GetField("ns-1", "name1", "f1")
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Equal(t, field.Meta{}, f)

	// case 4: backend, metric not exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(0), constants.ErrNotFound)
	f, err = db.GetField("ns-1", "name2", "f1")
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Equal(t, field.Meta{}, f)

	// case 4: backend exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	mockBackend.EXPECT().getField(uint32(10), field.Name("f1")).Return(field.Meta{ID: 19, Type: field.SumField}, nil)
	f, err = db.GetField("ns-1", "name2", "f1")
	assert.NoError(t, err)
	assert.Equal(t, field.Meta{ID: 19, Type: field.SumField}, f)

	// case 5: all tag keys from memory
	meta.EXPECT().getAllFields().Return([]field.Meta{{ID: 19, Type: field.SumField}})
	fields, err := db.GetAllFields("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, []field.Meta{{ID: 19, Type: field.SumField}}, fields)

	// case 6: backend, metric not exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(0), constants.ErrNotFound)
	fields, err = db.GetAllFields("ns-1", "name2")
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Nil(t, fields)

	// case 7: backend, tag keys exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	mockBackend.EXPECT().getAllFields(uint32(10)).Return([]field.Meta{{ID: 19, Type: field.SumField}}, nil)
	fields, err = db.GetAllFields("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, []field.Meta{{ID: 19, Type: field.SumField}}, fields)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GenMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	gomock.InOrder(
		mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(nil, constants.ErrNotFound),
		mockBackend.EXPECT().genMetricID().Return(uint32(1)),
	)
	// case 1: gen new metric id
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	// cast 2: get metric id from memory
	metricID, err = db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	// case 3: load metric meta err
	mockBackend.EXPECT().loadMetricMetadata("ns-1", "name2").Return(nil, fmt.Errorf("err"))
	metricID, err = db.GenMetricID("ns-1", "name2")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), metricID)

	// case 4: load metric meta ok
	meta := NewMockMetricMetadata(ctrl)
	meta.EXPECT().getMetricID().Return(uint32(100))
	mockBackend.EXPECT().loadMetricMetadata("ns-1", "name2").Return(meta, nil)
	metricID, err = db.GenMetricID("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(100), metricID)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GetMetricID_wal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	metricID, err := db.GenMetricID("ns", "metric")
	assert.Equal(t, uint32(1), metricID)
	assert.NoError(t, err)
	db1 := db.(*metadataDatabase)
	oldWAL := db1.metaWAL
	mockWAL := wal.NewMockMetricMetaWAL(ctrl)
	db1.metaWAL = mockWAL
	mockWAL.EXPECT().AppendMetric(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	metricID, err = db.GenMetricID("ns", "metric2")
	assert.Equal(t, uint32(0), metricID)
	assert.Error(t, err)
	db1.metaWAL = oldWAL
	metricID, err = db.GenMetricID("ns", "metric2")
	assert.Equal(t, uint32(2), metricID)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestMetadataDatabase_GenFieldID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	meta := NewMockMetricMetadata(ctrl)
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	gomock.InOrder(
		mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(meta, nil),
		meta.EXPECT().getMetricID().Return(uint32(100)),
		meta.EXPECT().getField(field.Name("f")).Return(field.Meta{}, false),
		meta.EXPECT().createField(gomock.Any(), gomock.Any()).Return(field.ID(10), nil),
		meta.EXPECT().getMetricID().Return(uint32(1)),
		meta.EXPECT().addField(gomock.Any()),
	)
	// case 1: gen new field id
	_, err = db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	fieldID, err := db.GenFieldID("ns-1", "name1", "f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(10), fieldID)

	// case 2: get field id from memory
	meta.EXPECT().getField(field.Name("f")).Return(field.Meta{ID: 10, Type: field.SumField}, true)
	fieldID, err = db.GenFieldID("ns-1", "name1", "f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(10), fieldID)

	// case 3: get field id from memory, but type not match
	meta.EXPECT().getField(field.Name("f")).Return(field.Meta{ID: 10, Type: field.MinField}, true)
	fieldID, err = db.GenFieldID("ns-1", "name1", "f", field.SumField)
	assert.Equal(t, series.ErrWrongFieldType, err)
	assert.Equal(t, field.ID(0), fieldID)

	// case 4: create fail
	gomock.InOrder(
		meta.EXPECT().getField(field.Name("f")).Return(field.Meta{}, false),
		meta.EXPECT().createField(gomock.Any(), gomock.Any()).Return(field.ID(10), fmt.Errorf("err")),
	)
	fieldID, err = db.GenFieldID("ns-1", "name1", "f", field.SumField)
	assert.Error(t, err)
	assert.Equal(t, field.ID(0), fieldID)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GetField_wal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	_, _ = db.GenMetricID("ns", "metric")
	fieldID, err := db.GenFieldID("ns", "metric", "f", field.SumField)
	assert.Equal(t, field.ID(1), fieldID)
	assert.NoError(t, err)
	db1 := db.(*metadataDatabase)
	oldWAL := db1.metaWAL
	mockWAL := wal.NewMockMetricMetaWAL(ctrl)
	db1.metaWAL = mockWAL
	mockWAL.EXPECT().AppendField(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	fieldID, err = db.GenFieldID("ns", "metric", "f2", field.SumField)
	assert.Equal(t, field.ID(0), fieldID)
	assert.Error(t, err)
	db1.metaWAL = oldWAL
	fieldID, err = db.GenFieldID("ns", "metric", "f2", field.SumField)
	assert.Equal(t, field.ID(2), fieldID)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestMetadataDatabase_GenTagKeyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	meta := NewMockMetricMetadata(ctrl)
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	gomock.InOrder(
		mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(meta, nil),
		meta.EXPECT().getMetricID().Return(uint32(100)),
		meta.EXPECT().getTagKeyID("tag-key").Return(uint32(0), false),
		meta.EXPECT().checkTagKeyCount().Return(nil),
		mockBackend.EXPECT().genTagKeyID().Return(uint32(10)),
		meta.EXPECT().getMetricID().Return(uint32(1)),
		meta.EXPECT().createTagKey("tag-key", uint32(10)),
	)
	// case 1: gen new tag key id
	_, err = db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	tagKeyID, err := db.GenTagKeyID("ns-1", "name1", "tag-key")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), tagKeyID)

	// case 2: get tag key id from memory
	meta.EXPECT().getTagKeyID("tag-key").Return(uint32(10), true)
	tagKeyID, err = db.GenTagKeyID("ns-1", "name1", "tag-key")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), tagKeyID)

	// case 3: too many tags
	gomock.InOrder(
		meta.EXPECT().getTagKeyID("tag-key").Return(uint32(0), false),
		meta.EXPECT().checkTagKeyCount().Return(series.ErrTooManyTagKeys),
	)
	tagKeyID, err = db.GenTagKeyID("ns-1", "name1", "tag-key")
	assert.Equal(t, series.ErrTooManyTagKeys, err)
	assert.Equal(t, uint32(0), tagKeyID)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GenTagKeyID_wal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	_, _ = db.GenMetricID("ns", "metric")
	tagKeyID, err := db.GenTagKeyID("ns", "metric", "tagKey")
	assert.Equal(t, uint32(1), tagKeyID)
	assert.NoError(t, err)
	db1 := db.(*metadataDatabase)
	oldWAL := db1.metaWAL
	mockWAL := wal.NewMockMetricMetaWAL(ctrl)
	db1.metaWAL = mockWAL
	mockWAL.EXPECT().AppendTagKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	tagKeyID, err = db.GenTagKeyID("ns", "metric", "tagKey2")
	assert.Equal(t, uint32(0), tagKeyID)
	assert.Error(t, err)
	db1.metaWAL = oldWAL
	tagKeyID, err = db.GenTagKeyID("ns", "metric", "tagKey2")
	assert.Equal(t, uint32(2), tagKeyID)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestMetadataDatabase_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		createMetaWAL = wal.NewMetricMetaWAL
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	mockWAL := wal.NewMockMetricMetaWAL(ctrl)
	mockWAL.EXPECT().Recovery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	mockWAL.EXPECT().NeedRecovery().Return(false)
	createMetaWAL = func(path string) (wal.MetricMetaWAL, error) {
		return mockWAL, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	mockBackend.EXPECT().Close().Return(fmt.Errorf("err"))
	mockWAL.EXPECT().Close().Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)
}

func TestMetadataDatabase_reopen(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	// mock one metric event
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	metricID, err = db.GenMetricID("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), metricID)
	_, err = db.GenTagKeyID("ns-1", "name1", "tagKey")
	assert.NoError(t, err)
	err = db.Close()
	assert.NoError(t, err)

	// reopen
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	metricID, err = db.GenMetricID("ns-2", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(3), metricID)
	err = db.Close()
	assert.NoError(t, err)
}

func TestMetadataDatabase_Sync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	db := newMockMetadataDatabase(t)
	db1 := db.(*metadataDatabase)
	mockWAL := wal.NewMockMetricMetaWAL(ctrl)
	mockWAL.EXPECT().Sync().Return(fmt.Errorf("err"))
	db1.metaWAL = mockWAL
	err := db.Sync()
	assert.NoError(t, err)
	mockWAL.EXPECT().Close().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_checkSync(t *testing.T) {
	syncInterval = 100
	ctrl := gomock.NewController(t)
	defer func() {
		syncInterval = 2 * timeutil.OneSecond
		_ = fileutil.RemoveDir(testPath)
		createMetaWAL = wal.NewMetricMetaWAL

		ctrl.Finish()
	}()

	var count atomic.Int32
	mockMetaWAL := wal.NewMockMetricMetaWAL(ctrl)
	mockMetaWAL.EXPECT().NeedRecovery().DoAndReturn(func() bool {
		count.Inc()
		return count.Load() != 1
	}).AnyTimes()
	mockMetaWAL.EXPECT().Recovery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	createMetaWAL = func(path string) (wal.MetricMetaWAL, error) {
		return mockMetaWAL, nil
	}

	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	time.Sleep(time.Second)

	mockMetaWAL.EXPECT().Close().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestMetadataDatabase_recovery_metric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()

	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	for i := 0; i < 11000; i++ {
		_, err := db.GenMetricID("ns", fmt.Sprintf("metric-%d", i))
		assert.NoError(t, err)
	}
	err = db.Close()
	assert.NoError(t, err)

	backend := NewMockMetadataBackend(ctrl)
	backend.EXPECT().Close().Return(nil).AnyTimes()
	createMetadataBackend = func(parent string) (MetadataBackend, error) {
		return backend, nil
	}
	backend.EXPECT().saveMetadata(gomock.Any()).Return(fmt.Errorf("err"))
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.Error(t, err)
	assert.Nil(t, db)

	createMetadataBackend = newMetadataBackend
	// recovery success
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	for i := 0; i < 100; i++ {
		_, err := db.GenMetricID("ns", fmt.Sprintf("metric-%d", i))
		assert.NoError(t, err)
	}
	err = db.Close()
	assert.NoError(t, err)

	createMetadataBackend = func(parent string) (MetadataBackend, error) {
		return backend, nil
	}
	backend.EXPECT().saveMetadata(gomock.Any()).Return(fmt.Errorf("err"))
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestMetadataDatabase_recovery_field(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()

	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	for i := 0; i < 9999; i++ {
		_, err := db.GenMetricID("ns", fmt.Sprintf("metric-%d", i))
		assert.NoError(t, err)
	}
	for i := 0; i < 20; i++ {
		_, err := db.GenFieldID("ns", "metric-1", field.Name(fmt.Sprintf("f-%d", i)), field.SumField)
		assert.NoError(t, err)
	}
	err = db.Close()
	assert.NoError(t, err)

	backend := NewMockMetadataBackend(ctrl)
	backend.EXPECT().Close().Return(nil).AnyTimes()
	createMetadataBackend = func(parent string) (MetadataBackend, error) {
		return backend, nil
	}
	backend.EXPECT().saveMetadata(gomock.Any()).Return(fmt.Errorf("err"))
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.Error(t, err)
	assert.Nil(t, db)

	createMetadataBackend = newMetadataBackend
	// recovery success
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Close()
	assert.NoError(t, err)
}

func TestMetadataDatabase_recovery_tagKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackend = newMetadataBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()

	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	for i := 0; i < 9999; i++ {
		_, err := db.GenMetricID("ns", fmt.Sprintf("metric-%d", i))
		assert.NoError(t, err)
	}
	for i := 0; i < 20; i++ {
		_, err := db.GenTagKeyID("ns", "metric-1", fmt.Sprintf("tagKey-%d", i))
		assert.NoError(t, err)
	}
	err = db.Close()
	assert.NoError(t, err)

	backend := NewMockMetadataBackend(ctrl)
	backend.EXPECT().Close().Return(nil).AnyTimes()
	createMetadataBackend = func(parent string) (MetadataBackend, error) {
		return backend, nil
	}
	backend.EXPECT().saveMetadata(gomock.Any()).Return(fmt.Errorf("err"))
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.Error(t, err)
	assert.Nil(t, db)

	createMetadataBackend = newMetadataBackend
	// recovery success
	db, err = NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Close()
	assert.NoError(t, err)
}

func newMockMetadataDatabase(t *testing.T) MetadataDatabase {
	db, err := NewMetadataDatabase(context.TODO(), "test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	return db
}
