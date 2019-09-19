package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func TestFieldAggregator_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	baseTime, _ := timeutil.ParseTimestamp("20190729 10:00:00")

	aggSpec := NewAggregatorSpec(uint16(15), "f", field.SumField)
	aggSpec.AddFunctionType(function.Sum)

	agg := NewFieldAggregator(baseTime, 10*timeutil.OneSecond, 10, 50, 1, aggSpec)
	it := MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.5,
		17: 5.5,
		16: 5.5,
		56: 5.5,
	})
	agg.Aggregate(it)

	expect := map[int]float64{
		5: 5.5,
		6: 5.5,
		7: 5.5,
	}

	fieldIt := agg.Iterator()
	assert.True(t, fieldIt.HasNext())
	AssertPrimitiveIt(t, fieldIt.Next(), expect)
	assert.Equal(t, uint16(15), fieldIt.FieldMeta().ID)
	assert.Equal(t, field.SumField, fieldIt.FieldMeta().Type)
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
		5: 11,
		6: 6.6,
		7: 5.5,
		9: 5.5,
	}

	fieldIt = agg.Iterator()
	assert.True(t, fieldIt.HasNext())
	AssertPrimitiveIt(t, fieldIt.Next(), expect)
	assert.Equal(t, uint16(15), fieldIt.FieldMeta().ID)
	assert.Equal(t, field.SumField, fieldIt.FieldMeta().Type)
	assert.False(t, fieldIt.HasNext())

	assert.Equal(t,
		timeutil.TimeRange{Start: baseTime + 10*10*timeutil.OneSecond, End: baseTime + 10*10*4*timeutil.OneSecond},
		agg.TimeRange())

	// not match primitive field
	agg = NewFieldAggregator(baseTime, 10*timeutil.OneSecond, 10, 50, 1, aggSpec)
	it = MockSumFieldIterator(ctrl, uint16(11), map[int]interface{}{})
	agg.Aggregate(it)

	fieldIt = agg.Iterator()
	assert.True(t, fieldIt.HasNext())
	assert.False(t, fieldIt.Next().HasNext())
	assert.False(t, fieldIt.HasNext())
}
