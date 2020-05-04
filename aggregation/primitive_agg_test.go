package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/selector"
	"github.com/lindb/lindb/series/field"
)

func TestPrimitiveSumFloatAgg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSelector := selector.NewMockSlotSelector(ctrl)
	mockSelector.EXPECT().PointCount().Return(5).AnyTimes()
	mockSelector.EXPECT().Range().Return(10, 15).AnyTimes()

	agg := NewPrimitiveAggregator(1, mockSelector, field.Sum.AggFunc())
	mockSelector.EXPECT().IndexOf(gomock.Any()).Return(1, false)
	agg.Aggregate(1, 10.0)
	mockSelector.EXPECT().IndexOf(gomock.Any()).Return(1, false)
	agg.Aggregate(1, 30.0)
	mockSelector.EXPECT().IndexOf(gomock.Any()).Return(-1, false)
	agg.Aggregate(-1, 30.0)
	mockSelector.EXPECT().IndexOf(gomock.Any()).Return(10, false)
	agg.Aggregate(10, 30.0)
	// completed
	mockSelector.EXPECT().IndexOf(gomock.Any()).Return(10, true)
	agg.Aggregate(10, 30.0)

	expect := map[int]float64{11: 40.0}
	it := agg.Iterator()
	assert.Equal(t, field.PrimitiveID(1), agg.FieldID())
	AssertPrimitiveIt(t, it, expect)

	agg.reset()
	it = agg.Iterator()
	assert.False(t, it.HasNext())
	mockSelector.EXPECT().IndexOf(gomock.Any()).Return(2, false)
	agg.Aggregate(2, 20.0)
	expect = map[int]float64{12: 20.0}
	it = agg.Iterator()
	assert.Equal(t, field.PrimitiveID(1), agg.FieldID())
	AssertPrimitiveIt(t, it, expect)
}
