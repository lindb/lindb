package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestFieldAggregator_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")

	aggSpec := NewDownSamplingSpec("f", field.SumField)
	aggSpec.AddFunctionType(function.Sum)

	agg := NewFieldAggregator(baseTime, 10*timeutil.OneSecond, 10, 50, 1, aggSpec)
	it := MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.6,
		17: 5.7,
		16: 5.8,
		56: 5.9,
	})
	it.EXPECT().SegmentStartTime().Return(baseTime)
	agg.Aggregate(it)

	expect := map[int]float64{
		5: 5.6,
		6: 5.8,
		7: 5.7,
	}

	fieldIt := agg.Iterator()
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
	it.EXPECT().SegmentStartTime().Return(baseTime)
	agg.Aggregate(it)

	expect = map[int]float64{
		5: 11.1,
		6: 6.9,
		7: 5.7,
		9: 5.5,
	}

	fieldIt = agg.Iterator()
	assert.True(t, fieldIt.HasNext())
	AssertPrimitiveIt(t, fieldIt.Next(), expect)
	assert.False(t, fieldIt.HasNext())

	assert.Equal(t,
		timeutil.TimeRange{Start: baseTime + 10*10*timeutil.OneSecond, End: baseTime + 10*10*4*timeutil.OneSecond},
		agg.TimeRange())

	// not match primitive field
	agg = NewFieldAggregator(baseTime, 10*timeutil.OneSecond, 10, 50, 1, aggSpec)
	it = MockSumFieldIterator(ctrl, uint16(11), map[int]interface{}{})
	it.EXPECT().SegmentStartTime().Return(baseTime + 30*10*timeutil.OneSecond)
	agg.Aggregate(it)

	fieldIt = agg.Iterator()
	assert.True(t, fieldIt.HasNext())
	assert.False(t, fieldIt.Next().HasNext())
	assert.False(t, fieldIt.HasNext())

	// not match query time range case 1
	agg = NewFieldAggregator(baseTime, 10*timeutil.OneSecond, 10, 50, 1, aggSpec)
	it = series.NewMockFieldIterator(ctrl)
	it.EXPECT().SegmentStartTime().Return(baseTime + 10*timeutil.OneSecond)
	agg.Aggregate(it)
	fieldIt = agg.Iterator()
	assert.True(t, fieldIt.HasNext())
	assert.False(t, fieldIt.Next().HasNext())
	assert.False(t, fieldIt.HasNext())
	// not match query time range case 2
	agg = NewFieldAggregator(baseTime, 10*timeutil.OneSecond, 10, 50, 1, aggSpec)
	it = series.NewMockFieldIterator(ctrl)
	it.EXPECT().SegmentStartTime().Return(baseTime - timeutil.OneHour)
	agg.Aggregate(it)
	fieldIt = agg.Iterator()
	assert.True(t, fieldIt.HasNext())
	assert.False(t, fieldIt.Next().HasNext())
	assert.False(t, fieldIt.HasNext())
}

func TestFieldAggregator_Aggregate_2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")

	aggSpec := NewDownSamplingSpec("f", field.SumField)
	aggSpec.AddFunctionType(function.Sum)

	// query time range 20190729 10:10:00 ~ 20190729 12:10:00
	agg := NewFieldAggregator(baseTime, timeutil.OneMinute, 10, 120, 1, aggSpec)
	it := MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.6,
		17: 5.7,
		16: 5.8,
		46: 5.9,
	})
	it.EXPECT().SegmentStartTime().Return(baseTime + timeutil.OneHour)
	agg.Aggregate(it)

	expect := map[int]float64{
		65:  5.5,
		75:  5.6,
		77:  5.7,
		76:  5.8,
		106: 5.9,
	}

	fieldIt := agg.Iterator()
	assert.True(t, fieldIt.HasNext())
	AssertPrimitiveIt(t, fieldIt.Next(), expect)
	assert.False(t, fieldIt.HasNext())
}
