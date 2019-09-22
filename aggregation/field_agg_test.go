package aggregation

//import (
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//
//	"github.com/lindb/lindb/aggregation/function"
//	"github.com/lindb/lindb/aggregation/selector"
//	"github.com/lindb/lindb/pkg/timeutil"
//	"github.com/lindb/lindb/series"
//	"github.com/lindb/lindb/series/field"
//)
//
//func TestFieldAggregator_Aggregate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")
//
//	aggSpec := NewAggregatorSpec("f", field.SumField)
//	aggSpec.AddFunctionType(function.Sum)
//
//	selector1 := selector.NewIndexSlotSelector(0, 60, 1)
//	agg := NewFieldAggregator(baseTime, selector1, true, aggSpec)
//	it := MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
//		5:  5.5,
//		15: 5.6,
//		17: 5.7,
//		16: 5.8,
//		56: 5.9,
//	})
//	agg.Aggregate(it)
//
//	expect := map[int]float64{
//		5: 5.6,
//		6: 5.8,
//		7: 5.7,
//	}
//
//	start, fieldIt := agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	AssertPrimitiveIt(t, fieldIt.Next(), expect)
//	assert.False(t, fieldIt.HasNext())
//
//	it = MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
//		5:  5.5,
//		15: 5.5,
//		19: 5.5,
//		16: 1.1,
//		56: 5.5,
//	})
//	agg.Aggregate(it)
//
//	expect = map[int]float64{
//		5: 11.1,
//		6: 6.9,
//		7: 5.7,
//		9: 5.5,
//	}
//
//	start, fieldIt = agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	AssertPrimitiveIt(t, fieldIt.Next(), expect)
//	assert.False(t, fieldIt.HasNext())
//
//	// not match primitive field
//	agg = NewFieldAggregator(baseTime, selector1, true, aggSpec)
//	it = MockSumFieldIterator(ctrl, uint16(11), map[int]interface{}{})
//	agg.Aggregate(it)
//
//	start, fieldIt = agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	assert.False(t, fieldIt.Next().HasNext())
//	assert.False(t, fieldIt.HasNext())
//
//	// not match query time range case 1
//	agg = NewFieldAggregator(baseTime, selector1, true, aggSpec)
//	it = series.NewMockFieldIterator(ctrl)
//	agg.Aggregate(it)
//	start, fieldIt = agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	assert.False(t, fieldIt.Next().HasNext())
//	assert.False(t, fieldIt.HasNext())
//	// not match query time range case 2
//	agg = NewFieldAggregator(baseTime, selector1, true, aggSpec)
//	it = series.NewMockFieldIterator(ctrl)
//	agg.Aggregate(it)
//	start, fieldIt = agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	assert.False(t, fieldIt.Next().HasNext())
//	assert.False(t, fieldIt.HasNext())
//}
//
//func TestFieldAggregator_Aggregate_2(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")
//
//	aggSpec := NewAggregatorSpec("f", field.SumField)
//	aggSpec.AddFunctionType(function.Sum)
//
//	// query time range 20190729 10:10:00 ~ 20190729 12:10:00
//	selector1 := selector.NewIndexSlotSelector(0, 60, 1)
//	agg := NewFieldAggregator(baseTime, selector1, true, aggSpec)
//	it := MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
//		5:  5.5,
//		15: 5.6,
//		17: 5.7,
//		16: 5.8,
//		46: 5.9,
//	})
//	agg.Aggregate(it)
//
//	expect := map[int]float64{
//		65:  5.5,
//		75:  5.6,
//		77:  5.7,
//		76:  5.8,
//		106: 5.9,
//	}
//
//	start, fieldIt := agg.ResultSet()
//	assert.Equal(t, baseTime, start)
//	assert.True(t, fieldIt.HasNext())
//	AssertPrimitiveIt(t, fieldIt.Next(), expect)
//	assert.False(t, fieldIt.HasNext())
//}
