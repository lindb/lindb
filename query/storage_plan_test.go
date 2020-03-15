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

	query, _ := sql.Parse("select f from cpu")
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
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "f").
		Return(field.Meta{ID: 10, Type: field.SumField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "a").
		Return(field.Meta{ID: 11, Type: field.MinField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "b").
		Return(field.Meta{ID: 12, Type: field.MaxField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "c").
		Return(field.Meta{ID: 13, Type: field.HistogramField}, nil).AnyTimes()
	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "e").
		Return(field.Meta{ID: 14, Type: field.HistogramField}, nil).AnyTimes()

	metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "no_f").
		Return(field.Meta{ID: 99, Type: field.HistogramField}, constants.ErrNotFound).AnyTimes()

	// error
	query := &stmt.Query{MetricName: "cpu"}
	plan := newStorageExecutePlan("ns", metadata, query)
	err := plan.Plan()
	assert.NotNil(t, err)
	query, _ = sql.Parse("select no_f from cpu")
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Equal(t, constants.ErrNotFound, err)

	// normal
	query, _ = sql.Parse("select f from cpu")
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.NoError(t, err)

	storagePlan := plan.(*storageExecutePlan)
	downSampling := aggregation.NewDownSamplingSpec("f", field.SumField)
	downSampling.AddFunctionType(function.Sum)
	assert.Equal(t, map[field.ID]aggregation.AggregatorSpec{field.ID(10): downSampling}, storagePlan.fields)
	assert.Equal(t, []field.ID{10}, storagePlan.getFieldIDs())

	query, _ = sql.Parse("select a,b,c as d from cpu")
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.NoError(t, err)

	storagePlan = plan.(*storageExecutePlan)
	downSampling1 := aggregation.NewDownSamplingSpec("a", field.MinField)
	downSampling1.AddFunctionType(function.Min)
	downSampling2 := aggregation.NewDownSamplingSpec("b", field.MaxField)
	downSampling2.AddFunctionType(function.Max)
	downSampling3 := aggregation.NewDownSamplingSpec("c", field.HistogramField)
	downSampling3.AddFunctionType(function.Histogram)
	expect := map[field.ID]aggregation.AggregatorSpec{
		field.ID(11): downSampling1,
		field.ID(12): downSampling2,
		field.ID(13): downSampling3,
	}
	assert.Equal(t, expect, storagePlan.fields)
	assert.Equal(t, []field.ID{11, 12, 13}, storagePlan.getFieldIDs())

	query, _ = sql.Parse("select min(a),max(sum(c)+avg(c)+e) as d from cpu")
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.NoError(t, err)
	storagePlan = plan.(*storageExecutePlan)

	downSampling1 = aggregation.NewDownSamplingSpec("a", field.MinField)
	downSampling1.AddFunctionType(function.Min)
	downSampling3 = aggregation.NewDownSamplingSpec("c", field.HistogramField)
	downSampling3.AddFunctionType(function.Sum)
	downSampling3.AddFunctionType(function.Avg)
	downSampling4 := aggregation.NewDownSamplingSpec("e", field.HistogramField)
	downSampling4.AddFunctionType(function.Histogram)
	expect = map[field.ID]aggregation.AggregatorSpec{
		field.ID(11): downSampling1,
		field.ID(13): downSampling3,
		field.ID(14): downSampling4,
	}
	assert.Equal(t, expect, storagePlan.fields)
	assert.Equal(t, []field.ID{11, 13, 14}, storagePlan.getFieldIDs())
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
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "f").
			Return(field.Meta{ID: 12, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "d").
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
	)

	// normal
	query, _ := sql.Parse("select f,d from disk group by host,path")
	plan := newStorageExecutePlan("ns", metadata, query)
	err := plan.Plan()
	assert.NoError(t, err)

	storagePlan := plan.(*storageExecutePlan)
	aggSpecs := storagePlan.getDownSamplingAggSpecs()
	assert.Equal(t, "d", aggSpecs[0].FieldName())
	assert.Equal(t, "f", aggSpecs[1].FieldName())

	assert.Equal(t, []field.ID{10, 12}, storagePlan.getFieldIDs())
	assert.Equal(t, 2, len(storagePlan.groupByTags))
	assert.Equal(t, []tag.Meta{{ID: 10, Key: "host"}, {ID: 11, Key: "path"}}, storagePlan.groupByKeyIDs())

	// get tag key err
	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "host").Return(uint32(0), fmt.Errorf("err")),
	)
	query, _ = sql.Parse("select f from disk group by host,path")
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
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "f").
			Return(field.Meta{ID: 10, Type: field.Unknown}, nil),
	)
	query, _ := sql.Parse("select f from disk")
	plan := newStorageExecutePlan("ns", metadata, query)
	err := plan.Plan()
	assert.Error(t, err)

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "f").
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
	)
	query, _ = sql.Parse("select histogram(f) from disk")
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Error(t, err)

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "d").
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "f").
			Return(field.Meta{ID: 10, Type: field.SumField}, nil),
	)
	query, _ = sql.Parse("select (d+histogram(f)+b) from disk")
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Error(t, err)

	gomock.InOrder(
		metadataDB.EXPECT().GetMetricID(gomock.Any(), "disk").Return(uint32(10), nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "d").
			Return(field.Meta{ID: 12, Type: field.SumField}, nil),
		metadataDB.EXPECT().GetField(gomock.Any(), gomock.Any(), "f").
			Return(field.Meta{ID: 11, Type: field.SumField}, nil),
	)
	query, _ = sql.Parse("select (d+histogram(f)+b),e from disk")
	plan = newStorageExecutePlan("ns", metadata, query)
	err = plan.Plan()
	assert.Error(t, err)
}
