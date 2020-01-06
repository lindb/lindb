package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func TestNewFieldAggregates(t *testing.T) {
	agg := NewFieldAggregates(
		timeutil.Interval(timeutil.OneSecond),
		1,
		timeutil.TimeRange{
			Start: 10,
			End:   20,
		}, true,
		AggregatorSpecs{
			NewDownSamplingSpec("b", field.SumField),
			NewDownSamplingSpec("a", field.SumField),
		})
	assert.Equal(t, "b", agg[0].FieldName())
	assert.Equal(t, "a", agg[1].FieldName())
	assert.Equal(t, field.SumField, agg[0].GetFieldType())
	assert.Equal(t, field.SumField, agg[1].GetFieldType())

	agg = NewFieldAggregates(
		timeutil.Interval(timeutil.OneSecond),
		1,
		timeutil.TimeRange{
			Start: 10,
			End:   20,
		}, false,
		AggregatorSpecs{
			NewDownSamplingSpec("b", field.SumField),
			NewDownSamplingSpec("a", field.SumField),
		})
	assert.Equal(t, "a", agg[0].FieldName())
	assert.Equal(t, "b", agg[1].FieldName())

	it := agg.ResultSet("")
	assert.True(t, it.HasNext())
	sIt := it.Next()
	assert.Equal(t, "a", sIt.FieldName())
	assert.Equal(t, field.SumField, sIt.FieldType())
	assert.True(t, it.HasNext())
	sIt = it.Next()
	assert.Equal(t, "b", sIt.FieldName())
	assert.Equal(t, field.SumField, sIt.FieldType())
	assert.False(t, it.HasNext())

	agg.Reset()
}

func TestNewSeriesAggregator(t *testing.T) {
	now, _ := timeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
	familyTime, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
	agg := NewSeriesAggregator(
		timeutil.Interval(timeutil.OneSecond),
		1,
		timeutil.TimeRange{
			Start: now,
			End:   now + 3*timeutil.OneHour,
		}, true,
		NewDownSamplingSpec("b", field.SumField),
	)

	fAgg, ok := agg.GetAggregator(familyTime)
	assert.True(t, ok)
	assert.NotNil(t, fAgg)

	fAgg, ok = agg.GetAggregator(familyTime - timeutil.OneHour)
	assert.False(t, ok)
	assert.Nil(t, fAgg)
	fAgg, ok = agg.GetAggregator(familyTime + 3*timeutil.OneHour)
	assert.True(t, ok)
	assert.NotNil(t, fAgg)
	fAgg, ok = agg.GetAggregator(familyTime + 4*timeutil.OneHour)
	assert.False(t, ok)
	assert.Nil(t, fAgg)

	rs := agg.ResultSet()
	assert.Equal(t, "b", rs.FieldName())
	assert.True(t, rs.HasNext())
	startTime, fIt := rs.Next()
	assert.Equal(t, familyTime, startTime)
	assert.NotNil(t, fIt)
	assert.True(t, rs.HasNext())
	startTime, fIt = rs.Next()
	assert.Equal(t, familyTime+3*timeutil.OneHour, startTime)
	assert.NotNil(t, fIt)
	assert.False(t, rs.HasNext())
	rs = agg.ResultSet()
	d, err := rs.MarshalBinary()
	assert.NoError(t, err)
	assert.True(t, len(d) > 0)

	agg.Reset()

	agg = NewSeriesAggregator(
		timeutil.Interval(timeutil.OneSecond),
		1,
		timeutil.TimeRange{
			Start: now,
			End:   now - 3*timeutil.OneHour,
		}, true,
		NewDownSamplingSpec("b", field.SumField),
	)
	fAgg, ok = agg.GetAggregator(familyTime)
	assert.False(t, ok)
	assert.Nil(t, fAgg)

	rs = agg.ResultSet()
	assert.Nil(t, rs)
}
