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

package indexdb

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/unique"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestNewIndexDatabase(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadata := metadb.NewMockMetadata(ctrl)
	mockMetadata.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, mockMetadata, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	// can't new duplicate
	db2, err := NewIndexDatabase(context.TODO(), testPath, nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db2)

	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_SuggestTagValues(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaDB := metadb.NewMockMetadata(ctrl)
	metaDB.EXPECT().DatabaseName().Return("test").AnyTimes()
	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metaDB.EXPECT().TagMetadata().Return(tagMeta)
	db, err := NewIndexDatabase(context.TODO(), testPath, metaDB, nil, nil)
	assert.NoError(t, err)
	tagMeta.EXPECT().SuggestTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a", "b"})
	tagValues := db.SuggestTagValues(10, "test", 100)
	assert.Equal(t, []string{"a", "b"}, tagValues)

	err = db.Close()
	assert.NoError(t, err)
}

func mockTagKeyValueIterator(kvs map[string]string) *metric.KeyValueIterator {
	var ml protoMetricsV1.MetricList
	var m = protoMetricsV1.Metric{
		Namespace: "ns",
		Name:      "name",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_Min, Value: 1},
		},
	}
	for k, v := range kvs {
		m.Tags = append(m.Tags, &protoMetricsV1.KeyValue{Key: k, Value: v})
	}

	ml.Metrics = append(ml.Metrics, &m)
	var buf bytes.Buffer
	converter := metric.NewProtoConverter()
	_, _ = converter.MarshalProtoMetricListV1To(ml, &buf)
	var br metric.StorageBatchRows
	br.UnmarshalRows(buf.Bytes())
	return br.Rows()[0].NewKeyValueIterator()
}

func TestIndexDatabase_BuildInvertIndex(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db1 := db.(*indexDatabase)
	index := NewMockInvertedIndex(ctrl)
	db1.index = index
	index.EXPECT().buildInvertIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	db.BuildInvertIndex("ns", "cpu", mockTagKeyValueIterator(map[string]string{"ip": "1.1.1.1"}), 10)

	index.EXPECT().Flush().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_GetOrCreateSeriesID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sequence := unique.NewMockSequence(ctrl)
	backend := NewMockIDMappingBackend(ctrl)
	mapping := NewMockMetricIDMapping(ctrl)
	db := &indexDatabase{
		backend: backend,
		metricID2Mapping: map[metric.ID]MetricIDMapping{
			2: mapping,
		},
	}

	cases := []struct {
		name     string
		metricID metric.ID
		tagsHash uint64
		prepare  func()
		out      struct {
			seriesID uint32
			isCreate bool
			err      error
		}
	}{
		{
			name:     "get series from cache",
			metricID: 2,
			tagsHash: 3,
			prepare: func() {
				mapping.EXPECT().GetSeriesID(gomock.Any()).Return(uint32(3), true)
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: uint32(3),
				isCreate: false,
				err:      nil,
			},
		},
		{
			name:     "get series id failure from backend",
			metricID: 2,
			tagsHash: 30,
			prepare: func() {
				mapping.EXPECT().GetSeriesID(gomock.Any()).Return(series.EmptySeriesID, false)
				backend.EXPECT().getSeriesID(gomock.Any(), gomock.Any()).Return(series.EmptySeriesID, fmt.Errorf("err"))
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: series.EmptySeriesID,
				isCreate: false,
				err:      fmt.Errorf("err"),
			},
		},
		{
			name:     "get series id successfully from backend",
			metricID: 2,
			tagsHash: 30,
			prepare: func() {
				mapping.EXPECT().GetSeriesID(gomock.Any()).Return(series.EmptySeriesID, false)
				mapping.EXPECT().AddSeriesID(gomock.Any(), gomock.Any())
				backend.EXPECT().getSeriesID(gomock.Any(), gomock.Any()).Return(uint32(30), nil)
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: uint32(30),
				isCreate: false,
				err:      nil,
			},
		},
		{
			name:     "save series id failure from backend",
			metricID: 2,
			tagsHash: 33,
			prepare: func() {
				mapping.EXPECT().GetSeriesID(gomock.Any()).Return(series.EmptySeriesID, false)
				mapping.EXPECT().GenSeriesID(gomock.Any()).Return(uint32(33))
				mapping.EXPECT().SeriesSequence().Return(sequence)
				sequence.EXPECT().HasNext().Return(true)
				backend.EXPECT().getSeriesID(gomock.Any(), gomock.Any()).Return(series.EmptySeriesID, constants.ErrNotFound)
				backend.EXPECT().genSeriesID(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: series.EmptySeriesID,
				isCreate: false,
				err:      fmt.Errorf("err"),
			},
		},
		{
			name:     "batch sequence failure from backend",
			metricID: 2,
			tagsHash: 3399,
			prepare: func() {
				mapping.EXPECT().GetSeriesID(gomock.Any()).Return(series.EmptySeriesID, false)
				backend.EXPECT().getSeriesID(metric.ID(2), uint64(3399)).Return(series.EmptySeriesID, constants.ErrNotFound)
				mapping.EXPECT().SeriesSequence().Return(sequence)
				sequence.EXPECT().HasNext().Return(false)
				sequence.EXPECT().Current().Return(uint32(20))
				backend.EXPECT().saveSeriesSequence(metric.ID(2), 20+config.GlobalStorageConfig().TSDB.SeriesSequenceCache).Return(fmt.Errorf("err"))
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: series.EmptySeriesID,
				isCreate: false,
				err:      fmt.Errorf("err"),
			},
		},
		{
			name:     "batch sequence successfully from backend, but store series id failure",
			metricID: 2,
			tagsHash: 339999,
			prepare: func() {
				mapping.EXPECT().GetSeriesID(gomock.Any()).Return(series.EmptySeriesID, false)
				backend.EXPECT().getSeriesID(metric.ID(2), uint64(339999)).Return(series.EmptySeriesID, constants.ErrNotFound)
				mapping.EXPECT().SeriesSequence().Return(sequence)
				sequence.EXPECT().HasNext().Return(false)
				sequence.EXPECT().Current().Return(uint32(20))
				backend.EXPECT().saveSeriesSequence(metric.ID(2), 20+config.GlobalStorageConfig().TSDB.SeriesSequenceCache).Return(nil)
				sequence.EXPECT().Limit(20 + config.GlobalStorageConfig().TSDB.SeriesSequenceCache)
				mapping.EXPECT().GenSeriesID(gomock.Any())
				backend.EXPECT().genSeriesID(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: series.EmptySeriesID,
				isCreate: false,
				err:      fmt.Errorf("err"),
			},
		},
		{

			name:     "save series id successfully from backend",
			metricID: 2,
			tagsHash: 333,
			prepare: func() {
				mapping.EXPECT().GetSeriesID(gomock.Any()).Return(series.EmptySeriesID, false)
				mapping.EXPECT().GenSeriesID(gomock.Any()).Return(uint32(333))
				mapping.EXPECT().SeriesSequence().Return(sequence)
				sequence.EXPECT().HasNext().Return(true)
				backend.EXPECT().getSeriesID(gomock.Any(), gomock.Any()).Return(series.EmptySeriesID, constants.ErrNotFound)
				backend.EXPECT().genSeriesID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: uint32(333),
				isCreate: true,
				err:      nil,
			},
		},
		{
			name:     "load mapping failure",
			metricID: 3,
			prepare: func() {
				backend.EXPECT().loadMetricIDMapping(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: series.EmptySeriesID,
				isCreate: false,
				err:      fmt.Errorf("err"),
			},
		},
		{
			name:     "load mapping successfully",
			metricID: 3,
			prepare: func() {
				mapping1 := NewMockMetricIDMapping(ctrl)
				backend.EXPECT().loadMetricIDMapping(gomock.Any()).Return(mapping1, nil)
				mapping1.EXPECT().AddSeriesID(gomock.Any(), gomock.Any())
				backend.EXPECT().getSeriesID(gomock.Any(), gomock.Any()).Return(uint32(30), nil)
			},
			out: struct {
				seriesID uint32
				isCreate bool
				err      error
			}{
				seriesID: uint32(30),
				isCreate: false,
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

			seriesID, isCreate, err := db.GetOrCreateSeriesID(tt.metricID, tt.tagsHash)
			assert.Equal(t, tt.out.seriesID, seriesID)
			assert.Equal(t, tt.out.isCreate, isCreate)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestIndexDatabase_GetGroupingContext(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	index := NewMockInvertedIndex(ctrl)
	db1 := db.(*indexDatabase)
	db1.index = index
	index.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, nil)
	ctx, err := db.GetGroupingContext([]uint32{1, 2}, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.Nil(t, ctx)

	index.EXPECT().Flush().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_GetSeriesIDs(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	index := NewMockInvertedIndex(ctrl)
	metaDB := metadb.NewMockMetadataDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	meta.EXPECT().MetadataDatabase().Return(metaDB).AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	db2 := db.(*indexDatabase)
	db2.index = index
	db2.metadata = meta
	assert.NoError(t, err)

	// case 1: get series ids by tag key
	index.EXPECT().GetSeriesIDsForTag(uint32(1)).Return(roaring.BitmapOf(1, 2), nil)
	seriesIDs, err := db.GetSeriesIDsForTag(1)
	assert.NoError(t, err)
	assert.NotNil(t, seriesIDs)
	// case 2: get series ids by tag value ids
	index.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), roaring.BitmapOf(1, 2, 3)).Return(roaring.BitmapOf(1, 2), nil)
	seriesIDs, err = db.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.NotNil(t, seriesIDs)
	// case 3: get tags err
	metaDB.EXPECT().GetAllTagKeys(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err = db.GetSeriesIDsForMetric("ns", "name")
	assert.Equal(t, fmt.Errorf("err"), err)
	assert.Nil(t, seriesIDs)
	// case 4: get empty tags
	metaDB.EXPECT().GetAllTagKeys(gomock.Any(), gomock.Any()).Return(nil, nil)
	seriesIDs, err = db.GetSeriesIDsForMetric("ns", "name")
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(0), seriesIDs)
	// case 5: get series ids for metric
	metaDB.EXPECT().GetAllTagKeys(gomock.Any(), gomock.Any()).Return([]tag.Meta{{ID: 1}}, nil)
	index.EXPECT().GetSeriesIDsForTags([]uint32{1}).Return(roaring.BitmapOf(1, 2, 3), nil)
	seriesIDs, err = db.GetSeriesIDsForMetric("ns", "name")
	assert.NoError(t, err)
	assert.NotNil(t, seriesIDs)

	index.EXPECT().Flush().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_Close(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer func() {
		createBackendFn = newIDMappingBackend
		ctrl.Finish()
	}()

	backend := NewMockIDMappingBackend(ctrl)
	createBackendFn = func(parent string) (IDMappingBackend, error) {
		return backend, nil
	}
	backend.EXPECT().sync().Return(nil)

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)

	assert.NoError(t, err)
	backend.EXPECT().Close().Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)

	backend.EXPECT().sync().Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)
}

func TestIndexDatabase_Flush(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer func() {

		createBackendFn = newIDMappingBackend
		ctrl.Finish()
	}()
	backend := NewMockIDMappingBackend(ctrl)
	createBackendFn = func(parent string) (IDMappingBackend, error) {
		return backend, nil
	}

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	backend.EXPECT().sync().Return(nil)
	assert.NoError(t, db.Flush())

	backend.EXPECT().sync().Return(fmt.Errorf("err"))
	assert.Error(t, db.Flush())
}
