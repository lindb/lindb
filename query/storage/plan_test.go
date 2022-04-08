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

package storagequery

import (
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestStoragePlan_Metric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	metadataDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil)
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(field.Meta{
			ID:   10,
			Type: field.SumField,
		}, nil).AnyTimes()

	q, _ := sql.Parse("select f from cpu")
	query := q.(*stmt.Query)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	plan := newStorageExecutePlan(ctx)
	err := plan.Plan()
	assert.NoError(t, err)

	metadataDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), constants.ErrNotFound)
	plan = newStorageExecutePlan(ctx)
	err = plan.Plan()
	assert.Equal(t, constants.ErrNotFound, err)
}

func TestStoragePlan_SelectList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	metadataDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(10), nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
		Return(field.Meta{ID: 10, Type: field.SumField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("a")).
		Return(field.Meta{ID: 11, Type: field.MinField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("b")).
		Return(field.Meta{ID: 12, Type: field.MaxField}, nil).AnyTimes()

	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("no_f")).
		Return(field.Meta{ID: 99, Type: field.SumField}, constants.ErrNotFound).AnyTimes()

	// error
	query := &stmt.Query{MetricName: "cpu"}
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	plan := newStorageExecutePlan(ctx)
	err := plan.Plan()
	assert.NotNil(t, err)
	q, _ := sql.Parse("select no_f from cpu")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	plan = newStorageExecutePlan(ctx)
	err = plan.Plan()
	assert.Equal(t, constants.ErrNotFound, err)

	// normal
	q, _ = sql.Parse("select f from cpu")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	storagePlan := newStorageExecutePlan(ctx)
	err = storagePlan.Plan()
	assert.NoError(t, err)

	downSampling := aggregation.NewAggregatorSpec("f", field.SumField)
	downSampling.AddFunctionType(function.Sum)
	assert.Equal(t, downSampling, storagePlan.fields[field.ID(10)].DownSampling)
	assert.Equal(t, field.Metas{{Name: "f", ID: 10, Type: field.SumField}}, ctx.storageExecuteCtx.Fields)

	// function not support
	q, _ = sql.Parse("select stddev(f) from cpu")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	storagePlan = newStorageExecutePlan(ctx)
	err = storagePlan.Plan()
	assert.Error(t, err)

	q, _ = sql.Parse("select a,b as d from cpu")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	storagePlan = newStorageExecutePlan(ctx)
	err = storagePlan.Plan()
	assert.NoError(t, err)

	downSampling1 := aggregation.NewAggregatorSpec("a", field.MinField)
	downSampling1.AddFunctionType(function.Min)
	downSampling2 := aggregation.NewAggregatorSpec("b", field.MaxField)
	downSampling2.AddFunctionType(function.Max)
	assert.Equal(t, downSampling1, storagePlan.fields[field.ID(11)].DownSampling)
	assert.Equal(t, downSampling2, storagePlan.fields[field.ID(12)].DownSampling)
	assert.Equal(t,
		field.Metas{
			{Name: "a", ID: 11, Type: field.MinField},
			{Name: "b", ID: 12, Type: field.MaxField},
		},
		ctx.storageExecuteCtx.Fields)

	q, _ = sql.Parse("select min(a) as d from cpu")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	storagePlan = newStorageExecutePlan(ctx)
	err = storagePlan.Plan()
	assert.NoError(t, err)

	downSampling1 = aggregation.NewAggregatorSpec("a", field.MinField)
	downSampling1.AddFunctionType(function.Min)
	assert.Equal(t, downSampling1, storagePlan.fields[field.ID(11)].DownSampling)
	assert.Equal(t,
		field.Metas{
			{Name: "a", ID: 11, Type: field.MinField},
		},
		ctx.storageExecuteCtx.Fields)
}

func TestStorageExecutePlan_groupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(metric.ID(10), nil),
		metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "host").Return(tag.KeyID(10), nil),
		metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "path").Return(tag.KeyID(11), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
			Return(field.Meta{ID: 12, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("d")).
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
	)

	// normal
	q, _ := sql.Parse("select f,d from disk group by host,path")
	query := q.(*stmt.Query)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	storagePlan := newStorageExecutePlan(ctx)
	err := storagePlan.Plan()
	assert.NoError(t, err)

	aggSpecs := ctx.storageExecuteCtx.DownSamplingSpecs
	assert.Equal(t, field.Name("d"), aggSpecs[0].FieldName())
	assert.Equal(t, field.Name("f"), aggSpecs[1].FieldName())

	assert.Equal(t, field.Metas{
		{Name: "d", ID: 10, Type: field.SumField}, {Name: "f", ID: 12, Type: field.SumField},
	}, ctx.storageExecuteCtx.Fields)
	assert.Equal(t, 2, len(ctx.storageExecuteCtx.GroupByTagKeyIDs))
	assert.Equal(t, 2, len(ctx.storageExecuteCtx.GroupByTags))
	assert.Equal(t, tag.Metas{{ID: 10, Key: "host"}, {ID: 11, Key: "path"}}, ctx.storageExecuteCtx.GroupByTags)

	// get tag key err
	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(metric.ID(10), nil),
		metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "host").Return(tag.KeyID(0), fmt.Errorf("err")),
	)
	q, _ = sql.Parse("select f from disk group by host,path")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	storagePlan = newStorageExecutePlan(ctx)
	err = storagePlan.Plan()
	assert.Error(t, err)
}

func TestStorageExecutePlan_empty_select_item(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(metric.ID(10), nil),
	)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: &stmt.Query{MetricName: "disk"},
		},
	}
	plan := newStorageExecutePlan(ctx)
	err := plan.Plan()
	assert.Equal(t, errEmptySelectList, err)
}

func TestStorageExecutePlan_field_expr_fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(metric.ID(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
			Return(field.Meta{ID: 10, Type: field.Unknown}, nil),
	)
	q, _ := sql.Parse("select f from disk")
	query := q.(*stmt.Query)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	plan := newStorageExecutePlan(ctx)
	err := plan.Plan()
	assert.Error(t, err)

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").
			Return(metric.ID(10), nil).AnyTimes(),
		metadataDB.EXPECT().GetAllHistogramFields(gomock.Any(), gomock.Any()).
			Return(histogramFieldMetas, nil).AnyTimes(),
	)
	// params more than one
	q, _ = sql.Parse("select quantile(0.99,1.0) from disk")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	plan = newStorageExecutePlan(ctx)
	err = plan.Plan()
	assert.Error(t, err)

	// quantile param not float
	q, _ = sql.Parse("select quantile(xxxx) from disk")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	plan = newStorageExecutePlan(ctx)
	err = plan.Plan()
	assert.Error(t, err)

	// quantile value range bad
	q, _ = sql.Parse("select quantile(-0.2) from disk")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	plan = newStorageExecutePlan(ctx)
	err = plan.Plan()
	assert.Error(t, err)
}

var (
	histogramFieldMetas = field.Metas{
		{Name: field.Name(metric.BucketNameOfHistogramExplicitBound(0.1)), ID: 1, Type: field.HistogramField},
		{Name: field.Name(metric.BucketNameOfHistogramExplicitBound(0.2)), ID: 2, Type: field.HistogramField},
		{Name: field.Name(metric.BucketNameOfHistogramExplicitBound(0.4)), ID: 3, Type: field.HistogramField},
		{Name: field.Name(metric.BucketNameOfHistogramExplicitBound(0.8)), ID: 4, Type: field.HistogramField},
		{Name: field.Name(metric.BucketNameOfHistogramExplicitBound(1.0)), ID: 5, Type: field.HistogramField},
		{Name: field.Name(metric.BucketNameOfHistogramExplicitBound(math.MaxFloat64 + 1)), ID: 6, Type: field.HistogramField},
	}
)

func TestStorageExecutePlan_field_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	metadataDB.EXPECT().GetAllHistogramFields(gomock.Any(), gomock.Any()).
		Return(histogramFieldMetas, nil).AnyTimes()

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(metric.ID(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("d")).
			Return(field.Meta{ID: 12, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("b")).
			Return(field.Meta{ID: 11, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("e")).
			Return(field.Meta{ID: 11, Type: field.SumField}, nil),
	)

	q, err := sql.Parse("select (d+quantile(0.1)*10+b),e from disk")
	assert.Nil(t, err)
	query := q.(*stmt.Query)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	plan := newStorageExecutePlan(ctx)
	err = plan.Plan()
	assert.Nil(t, err)
}

func TestStorageExecutePlan_histogramFieldsBad(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	metadataDB.EXPECT().GetAllHistogramFields(gomock.Any(), gomock.Any()).
		Return(nil, io.ErrClosedPipe).AnyTimes()
	metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(metric.ID(10), nil)

	q, _ := sql.Parse("select quantile(0.1) from disk")
	query := q.(*stmt.Query)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	plan := newStorageExecutePlan(ctx)
	err := plan.Plan()
	plan.field(nil, nil)
	assert.Error(t, err)
}
