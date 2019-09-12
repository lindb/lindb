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
	}, t1, map[string]*AggregatorSpec{"f1": {
		fieldName: "f1",
		fieldType: field.SumField,
		functions: map[function.FuncType]bool{function.Sum: true},
	}, "f2": {
		fieldName: "f2",
		fieldType: field.SumField,
		functions: map[function.FuncType]bool{function.Sum: true},
	}})

	it.EXPECT().FieldMeta().Return(field.Meta{Name: "f1"})
	it.EXPECT().SegmentStartTime().Return(t1)
	agg.Aggregate(it)
	agg.Aggregate(nil)
	it = MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.5,
		17: 5.5,
		16: 5.5,
		56: 5.5,
	})
	it.EXPECT().FieldMeta().Return(field.Meta{Name: "f2"})
	it.EXPECT().SegmentStartTime().Return(t1)
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

func TestSeriesSegmentAggregator_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.5,
		17: 5.5,
		16: 5.5,
		56: 5.5,
	})

	t1, _ := timeutil.ParseTimestamp("20191201 00:30:00", "20060102 15:04:05")
	t2, _ := timeutil.ParseTimestamp("20191201 04:30:00", "20060102 15:04:05")
	t3, _ := timeutil.ParseTimestamp("20191201 04:00:00", "20060102 15:04:05")
	agg := NewSeriesSegmentAggregator(timeutil.OneMinute, &timeutil.TimeRange{
		Start: t1,
		End:   t2,
	}, map[string]*AggregatorSpec{"f1": {
		fieldName: "f1",
		fieldType: field.SumField,
		functions: map[function.FuncType]bool{function.Sum: true},
	}, "f2": {
		fieldName: "f2",
		fieldType: field.SumField,
		functions: map[function.FuncType]bool{function.Sum: true},
	}})
	it.EXPECT().FieldMeta().Return(field.Meta{Name: "f1"})
	it.EXPECT().SegmentStartTime().Return(t3)
	agg.Aggregate(it)
	it = MockSumFieldIterator(ctrl, uint16(1), map[int]interface{}{
		5:  5.5,
		15: 5.5,
		17: 5.5,
		16: 5.5,
		56: 5.5,
	})
	it.EXPECT().FieldMeta().Return(field.Meta{Name: "f2"})
	it.EXPECT().SegmentStartTime().Return(t3)
	agg.Aggregate(it)

	tags := map[string]string{
		"host": "1.1.1.",
	}
	result := agg.Iterator(tags)
	assert.Equal(t, tags, result.Tags())
	assert.NotNil(t, result)
	f := make(map[string]bool)
	assert.True(t, result.HasNext())
	fIt := result.Next()
	f[fIt.FieldMeta().Name] = true
	expect := map[int]float64{
		215: 5.5,
		225: 5.5,
		227: 5.5,
		226: 5.5,
	}
	assert.True(t, fIt.HasNext())
	AssertPrimitiveIt(t, fIt.Next(), expect)
	assert.True(t, result.HasNext())
	fIt = result.Next()
	f[fIt.FieldMeta().Name] = true
	expect = map[int]float64{
		215: 5.5,
		225: 5.5,
		227: 5.5,
		226: 5.5,
	}
	assert.True(t, fIt.HasNext())
	AssertPrimitiveIt(t, fIt.Next(), expect)
	assert.False(t, result.HasNext())
	assert.Equal(t, 2, len(f))
}
