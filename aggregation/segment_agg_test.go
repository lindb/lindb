package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func TestSegmentAggregator_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.5,
		17: 5.5,
		16: 5.5,
		56: 5.5,
	})
	t1, _ := timeutil.ParseTimestamp("20191201 00:00:00", "20060102 15:04:05")
	t2, _ := timeutil.ParseTimestamp("20191201 01:00:00", "20060102 15:04:05")
	agg := NewSegmentAggregator(20000, 10000, &timeutil.TimeRange{
		Start: t1,
		End:   t2,
	}, t1, map[uint16]*AggregatorSpec{uint16(1): {
		fieldID:   uint16(1),
		fieldName: "f1",
		fieldType: field.SumField,
		functions: map[function.FuncType]bool{function.Sum: true},
	}, uint16(2): {
		fieldID:   uint16(1),
		fieldName: "f2",
		fieldType: field.SumField,
		functions: map[function.FuncType]bool{function.Sum: true},
	}})

	it.EXPECT().FieldMeta().Return(field.Meta{ID: 1})
	agg.Aggregate(it)
	agg.Aggregate(nil)
	it = MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.5,
		17: 5.5,
		16: 5.5,
		56: 5.5,
	})
	it.EXPECT().FieldMeta().Return(field.Meta{ID: 2})
	agg.Aggregate(it)

	result := agg.Iterator(nil)
	assert.NotNil(t, result)
	f := make(map[string]bool)
	assert.True(t, result.HasNext())
	fIt := result.Next()
	f[fIt.FieldMeta().Name] = true
	expect := map[int]float64{
		2:  5.5,
		7:  5.5,
		8:  11.0,
		28: 5.5,
	}
	assert.True(t, fIt.HasNext())
	AssertPrimitiveIt(t, fIt.Next(), expect)
	assert.True(t, result.HasNext())
	fIt = result.Next()
	f[fIt.FieldMeta().Name] = true
	expect = map[int]float64{
		2:  5.5,
		7:  5.5,
		8:  11.0,
		28: 5.5,
	}
	assert.True(t, fIt.HasNext())
	AssertPrimitiveIt(t, fIt.Next(), expect)
	assert.False(t, result.HasNext())
	assert.Equal(t, 2, len(f))
}
