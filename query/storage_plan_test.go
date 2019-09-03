package query

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/diskdb"
)

func TestStoragePlan_Metric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadataIndex := diskdb.NewMockIDGetter(ctrl)
	metadataIndex.EXPECT().GetMetricID(gomock.Any()).Return(uint32(10), nil)
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), gomock.Any()).
		Return(uint16(10), field.SumField, nil).AnyTimes()

	query, _ := sql.Parse("select f from cpu")
	plan := newStorageExecutePlan(metadataIndex, query)
	err := plan.Plan()
	if err != nil {
		t.Fatal(err)
	}

	metadataIndex.EXPECT().GetMetricID(gomock.Any()).Return(uint32(0), series.ErrNotFound)
	plan = newStorageExecutePlan(metadataIndex, query)
	err = plan.Plan()
	assert.Equal(t, series.ErrNotFound, err)
}

func TestStoragePlan_SelectList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadataIndex := diskdb.NewMockIDGetter(ctrl)
	metadataIndex.EXPECT().GetMetricID(gomock.Any()).Return(uint32(10), nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), "f").
		Return(uint16(10), field.SumField, nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), "a").
		Return(uint16(11), field.MinField, nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), "b").
		Return(uint16(12), field.MaxField, nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), "c").
		Return(uint16(13), field.HistogramField, nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), "e").
		Return(uint16(14), field.HistogramField, nil).AnyTimes()

	metadataIndex.EXPECT().GetFieldID(gomock.Any(), "no_f").
		Return(uint16(99), field.HistogramField, series.ErrNotFound).AnyTimes()

	// error
	query := &stmt.Query{MetricName: "cpu"}
	plan := newStorageExecutePlan(metadataIndex, query)
	err := plan.Plan()
	assert.NotNil(t, err)
	query, _ = sql.Parse("select no_f from cpu")
	plan = newStorageExecutePlan(metadataIndex, query)
	err = plan.Plan()
	assert.Equal(t, series.ErrNotFound, err)

	// normal
	query, _ = sql.Parse("select f from cpu")
	plan = newStorageExecutePlan(metadataIndex, query)
	err = plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	storagePlan := plan.(*storageExecutePlan)
	downSampling := aggregation.NewAggregatorSpec(uint16(10), "f", field.SumField)
	downSampling.AddFunctionType(function.Sum)
	assert.Equal(t, map[uint16]*aggregation.AggregatorSpec{uint16(10): downSampling}, storagePlan.fields)
	assert.Equal(t, []uint16{uint16(10)}, storagePlan.getFieldIDs())

	query, _ = sql.Parse("select a,b,c as d from cpu")
	plan = newStorageExecutePlan(metadataIndex, query)
	err = plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	storagePlan = plan.(*storageExecutePlan)
	downSampling1 := aggregation.NewAggregatorSpec(uint16(11), "a", field.MinField)
	downSampling1.AddFunctionType(function.Min)
	downSampling2 := aggregation.NewAggregatorSpec(uint16(12), "b", field.MaxField)
	downSampling2.AddFunctionType(function.Max)
	downSampling3 := aggregation.NewAggregatorSpec(uint16(13), "c", field.HistogramField)
	downSampling3.AddFunctionType(function.Histogram)
	expect := map[uint16]*aggregation.AggregatorSpec{
		uint16(11): downSampling1,
		uint16(12): downSampling2,
		uint16(13): downSampling3,
	}
	assert.Equal(t, expect, storagePlan.fields)
	assert.Equal(t, []uint16{uint16(11), uint16(12), uint16(13)}, storagePlan.getFieldIDs())

	query, _ = sql.Parse("select min(a),max(sum(c)+avg(c)+e) as d from cpu")
	plan = newStorageExecutePlan(metadataIndex, query)
	err = plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	storagePlan = plan.(*storageExecutePlan)

	downSampling1 = aggregation.NewAggregatorSpec(uint16(11), "a", field.MinField)
	downSampling1.AddFunctionType(function.Min)
	downSampling3 = aggregation.NewAggregatorSpec(uint16(13), "c", field.HistogramField)
	downSampling3.AddFunctionType(function.Sum)
	downSampling3.AddFunctionType(function.Avg)
	downSampling4 := aggregation.NewAggregatorSpec(uint16(14), "e", field.HistogramField)
	downSampling4.AddFunctionType(function.Histogram)
	expect = map[uint16]*aggregation.AggregatorSpec{
		uint16(11): downSampling1,
		uint16(13): downSampling3,
		uint16(14): downSampling4,
	}
	assert.Equal(t, expect, storagePlan.fields)
	assert.Equal(t, []uint16{uint16(11), uint16(13), uint16(14)}, storagePlan.getFieldIDs())
}

func TestStorageExecutePlan_groupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGetter := diskdb.NewMockIDGetter(ctrl)
	gomock.InOrder(
		idGetter.EXPECT().GetMetricID("disk").Return(uint32(10), nil),
		idGetter.EXPECT().GetTagKeyID(uint32(10), "host").Return(uint32(10), nil),
		idGetter.EXPECT().GetTagKeyID(uint32(10), "path").Return(uint32(11), nil),
		idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil),
	)

	// normal
	query, _ := sql.Parse("select f from disk group by host,path")
	plan := newStorageExecutePlan(idGetter, query)
	err := plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	storagePlan := plan.(*storageExecutePlan)
	downSampling := aggregation.NewAggregatorSpec(uint16(10), "f", field.SumField)
	downSampling.AddFunctionType(function.Sum)
	assert.Equal(t, map[uint16]*aggregation.AggregatorSpec{uint16(10): downSampling}, storagePlan.fields)
	assert.Equal(t, []uint16{uint16(10)}, storagePlan.getFieldIDs())
	assert.Equal(t, 2, len(storagePlan.groupByTagKeys))
	assert.Equal(t, uint32(10), storagePlan.groupByTagKeys["host"])
	assert.Equal(t, uint32(11), storagePlan.groupByTagKeys["path"])
	assert.True(t, storagePlan.hasGroupBy())

	// get tag key err
	gomock.InOrder(
		idGetter.EXPECT().GetMetricID("disk").Return(uint32(10), nil),
		idGetter.EXPECT().GetTagKeyID(uint32(10), "host").Return(uint32(0), fmt.Errorf("err")),
	)
	query, _ = sql.Parse("select f from disk group by host,path")
	plan = newStorageExecutePlan(idGetter, query)
	err = plan.Plan()
	assert.NotNil(t, err)
}

func TestStorageExecutePlan_empty_select_item(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGetter := diskdb.NewMockIDGetter(ctrl)
	gomock.InOrder(
		idGetter.EXPECT().GetMetricID("disk").Return(uint32(10), nil),
	)
	plan := newStorageExecutePlan(idGetter, &stmt.Query{MetricName: "disk"})
	err := plan.Plan()
	assert.Equal(t, errEmptySelectList, err)
}

func TestStorageExecutePlan_field_expr_fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGetter := diskdb.NewMockIDGetter(ctrl)
	gomock.InOrder(
		idGetter.EXPECT().GetMetricID("disk").Return(uint32(10), nil),
		idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.Unknown, nil),
	)
	query, _ := sql.Parse("select f from disk")
	plan := newStorageExecutePlan(idGetter, query)
	err := plan.Plan()
	assert.NotNil(t, err)

	gomock.InOrder(
		idGetter.EXPECT().GetMetricID("disk").Return(uint32(10), nil),
		idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil),
	)
	query, _ = sql.Parse("select histogram(f) from disk")
	plan = newStorageExecutePlan(idGetter, query)
	err = plan.Plan()
	assert.NotNil(t, err)

	gomock.InOrder(
		idGetter.EXPECT().GetMetricID("disk").Return(uint32(10), nil),
		idGetter.EXPECT().GetFieldID(uint32(10), "d").Return(uint16(10), field.SumField, nil),
		idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil),
	)
	query, _ = sql.Parse("select (d+histogram(f)+b) from disk")
	plan = newStorageExecutePlan(idGetter, query)
	err = plan.Plan()
	assert.NotNil(t, err)

	gomock.InOrder(
		idGetter.EXPECT().GetMetricID("disk").Return(uint32(10), nil),
		idGetter.EXPECT().GetFieldID(uint32(10), "d").Return(uint16(10), field.SumField, nil),
		idGetter.EXPECT().GetFieldID(uint32(10), "f").Return(uint16(10), field.SumField, nil),
	)
	query, _ = sql.Parse("select (d+histogram(f)+b),e from disk")
	plan = newStorageExecutePlan(idGetter, query)
	err = plan.Plan()
	assert.NotNil(t, err)
}
