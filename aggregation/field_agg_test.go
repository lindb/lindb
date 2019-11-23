package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/aggregation/selector"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestFieldAggregator_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")

	aggSpec := NewAggregatorSpec("f", field.SumField)
	aggSpec.AddFunctionType(function.Sum)

	selector1 := selector.NewIndexSlotSelector(15, 55, 1)
	agg := NewFieldAggregator(baseTime, selector1)
	it := MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.6,
		17: 5.7,
		16: 5.8,
		56: 5.9,
	})
	agg.Aggregate(it)
	assert.Equal(t, 1, len(agg.GetAllAggregators()))

	expect := map[int]float64{
		15: 5.6,
		16: 5.8,
		17: 5.7,
	}

	start, fieldIt := agg.ResultSet()
	assert.Equal(t, baseTime, start)
	assert.True(t, fieldIt.HasNext())
	AssertPrimitiveIt(t, fieldIt.Next(), expect)
	assert.False(t, fieldIt.HasNext())

	it = MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.5,
		19: 5.5,
		16: 1.1,
		56: 5.5,
	})
	agg.Aggregate(it)

	expect = map[int]float64{
		15: 11.1,
		16: 6.9,
		17: 5.7,
		19: 5.5,
	}

	start, fieldIt = agg.ResultSet()
	assert.Equal(t, baseTime, start)
	assert.True(t, fieldIt.HasNext())
	AssertPrimitiveIt(t, fieldIt.Next(), expect)
	assert.False(t, fieldIt.HasNext())

	// not match query time range case 1
	agg.reset()
	it = MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		4:  1.1,
		56: 5.5,
	})
	agg.Aggregate(it)
	start, fieldIt = agg.ResultSet()
	assert.Equal(t, baseTime, start)
	assert.True(t, fieldIt.HasNext())
	assert.False(t, fieldIt.Next().HasNext())
	assert.False(t, fieldIt.HasNext())
}

func TestDownSamplingFieldAggregator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")

	aggSpec := NewAggregatorSpec("f", field.SummaryField)
	aggSpec.AddFunctionType(function.Sum)
	aggSpec.AddFunctionType(function.Max)
	aggSpec.AddFunctionType(function.Avg)

	// query time range 20190729 10:10:00 ~ 20190729 12:10:00
	selector1 := selector.NewIndexSlotSelector(0, 60, 1)
	agg := NewDownSamplingFieldAggregator(baseTime, selector1, aggSpec)
	assert.Equal(t, 3, len(agg.GetAllAggregators()))
	agg.Aggregate(nil)
	it := series.NewMockFieldIterator(ctrl)
	agg.Aggregate(it)
	startTime, rs := agg.ResultSet()
	assert.NotNil(t, rs)
	assert.Equal(t, baseTime, startTime)
	agg.reset()
}
