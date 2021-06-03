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

package query

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestStoragePlan_Metric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	metadataDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(uint32(10), nil)
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(field.Meta{
			ID:   10,
			Type: field.SumField,
		}, nil).AnyTimes()

	q, _ := sql.Parse("select f from cpu")
	query := q.(*stmt.Query)
	plan := newStorageExecutePlan("ns", metadata, query)
	err := plan.Plan()
	assert.NoError(t, err)

	metadataDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(uint32(0), constants.ErrNotFound)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Equal(t, constants.ErrNotFound, err)
}

func TestStoragePlan_SelectList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	metadataDB.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(uint32(10), nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
		Return(field.Meta{ID: 10, Type: field.SumField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("a")).
		Return(field.Meta{ID: 11, Type: field.MinField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("b")).
		Return(field.Meta{ID: 12, Type: field.MaxField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("c")).
		Return(field.Meta{ID: 13, Type: field.HistogramField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("e")).
		Return(field.Meta{ID: 14, Type: field.HistogramField}, nil).AnyTimes()

	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("no_f")).
		Return(field.Meta{ID: 99, Type: field.HistogramField}, constants.ErrNotFound).AnyTimes()

	// error
	query := &stmt.Query{MetricName: "cpu"}
	plan := newStorageExecutePlan("ns", metadata, query)
	err := plan.Plan()
	assert.NotNil(t, err)
	q, _ := sql.Parse("select no_f from cpu")
	query = q.(*stmt.Query)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Equal(t, constants.ErrNotFound, err)

	// normal
	q, _ = sql.Parse("select f from cpu")
	query = q.(*stmt.Query)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.NoError(t, err)

	storagePlan := plan.(*storageExecutePlan)
	downSampling := aggregation.NewAggregatorSpec("f", field.SumField)
	downSampling.AddFunctionType(function.Sum)
	assert.Equal(t, downSampling, storagePlan.fields[field.ID(10)].DownSampling)
	assert.Equal(t, field.Metas{{Name: "f", ID: 10, Type: field.SumField}}, storagePlan.getFields())

	q, _ = sql.Parse("select a,b,c as d from cpu")
	query = q.(*stmt.Query)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.NoError(t, err)

	storagePlan = plan.(*storageExecutePlan)
	downSampling1 := aggregation.NewAggregatorSpec("a", field.MinField)
	downSampling1.AddFunctionType(function.Min)
	downSampling2 := aggregation.NewAggregatorSpec("b", field.MaxField)
	downSampling2.AddFunctionType(function.Max)
	downSampling3 := aggregation.NewAggregatorSpec("c", field.HistogramField)
	downSampling3.AddFunctionType(function.Histogram)
	assert.Equal(t, downSampling1, storagePlan.fields[field.ID(11)].DownSampling)
	assert.Equal(t, downSampling2, storagePlan.fields[field.ID(12)].DownSampling)
	assert.Equal(t, downSampling3, storagePlan.fields[field.ID(13)].DownSampling)
	assert.Equal(t,
		field.Metas{
			{Name: "a", ID: 11, Type: field.MinField},
			{Name: "b", ID: 12, Type: field.MaxField},
			{Name: "c", ID: 13, Type: field.HistogramField},
		},
		storagePlan.getFields())

	q, _ = sql.Parse("select min(a),max(sum(c)+avg(c)+e) as d from cpu")
	query = q.(*stmt.Query)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.NoError(t, err)
	storagePlan = plan.(*storageExecutePlan)

	downSampling1 = aggregation.NewAggregatorSpec("a", field.MinField)
	downSampling1.AddFunctionType(function.Min)
	downSampling3 = aggregation.NewAggregatorSpec("c", field.HistogramField)
	downSampling3.AddFunctionType(function.Sum)
	downSampling3.AddFunctionType(function.Avg)
	downSampling4 := aggregation.NewAggregatorSpec("e", field.HistogramField)
	downSampling4.AddFunctionType(function.Histogram)
	assert.Equal(t, downSampling1, storagePlan.fields[field.ID(11)].DownSampling)
	assert.Equal(t, downSampling3, storagePlan.fields[field.ID(13)].DownSampling)
	assert.Equal(t, downSampling4, storagePlan.fields[field.ID(14)].DownSampling)
	assert.Equal(t,
		field.Metas{
			{Name: "a", ID: 11, Type: field.MinField},
			{Name: "c", ID: 13, Type: field.HistogramField},
			{Name: "e", ID: 14, Type: field.HistogramField},
		},
		storagePlan.getFields())
}

func TestStorageExecutePlan_groupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "host").Return(uint32(10), nil),
		metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "path").Return(uint32(11), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
			Return(field.Meta{ID: 12, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("d")).
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
	)

	// normal
	q, _ := sql.Parse("select f,d from disk group by host,path")
	query := q.(*stmt.Query)
	plan := newStorageExecutePlan("ns", metadata, query)
	err := plan.Plan()
	assert.NoError(t, err)

	storagePlan := plan.(*storageExecutePlan)
	aggSpecs := storagePlan.getDownSamplingAggSpecs()
	assert.Equal(t, field.Name("d"), aggSpecs[0].FieldName())
	assert.Equal(t, field.Name("f"), aggSpecs[1].FieldName())

	assert.Equal(t, field.Metas{{Name: "d", ID: 10, Type: field.SumField}, {Name: "f", ID: 12, Type: field.SumField}}, storagePlan.getFields())
	assert.Equal(t, 2, len(storagePlan.groupByTags))
	assert.Equal(t, []tag.Meta{{ID: 10, Key: "host"}, {ID: 11, Key: "path"}}, storagePlan.groupByKeyIDs())

	// get tag key err
	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "host").Return(uint32(0), fmt.Errorf("err")),
	)
	q, _ = sql.Parse("select f from disk group by host,path")
	query = q.(*stmt.Query)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Error(t, err)
}

func TestStorageExecutePlan_empty_select_item(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
	)
	plan := newStorageExecutePlan("ns", metadata, &stmt.Query{MetricName: "disk"})
	err := plan.Plan()
	assert.Equal(t, errEmptySelectList, err)
}

func TestStorageExecutePlan_field_expr_fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
			Return(field.Meta{ID: 10, Type: field.Unknown}, nil),
	)
	q, _ := sql.Parse("select f from disk")
	query := q.(*stmt.Query)
	plan := newStorageExecutePlan("ns", metadata, query)
	err := plan.Plan()
	assert.Error(t, err)

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
	)
	q, _ = sql.Parse("select histogram(f) from disk")
	query = q.(*stmt.Query)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Error(t, err)

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("d")).
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
	)
	q, _ = sql.Parse("select (d+histogram(f)+b) from disk")
	query = q.(*stmt.Query)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Error(t, err)

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("d")).
			Return(field.Meta{ID: 12, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
			Return(field.Meta{ID: 11, Type: field.SumField}, nil),
	)
	q, _ = sql.Parse("select (d+histogram(f)+b),e from disk")
	query = q.(*stmt.Query)
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Error(t, err)
}
