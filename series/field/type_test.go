package field

import (
	"testing"

	"github.com/lindb/lindb/aggregation/function"

	"github.com/stretchr/testify/assert"
)

func TestDownSamplingFunc(t *testing.T) {
	assert.Equal(t, function.Sum, SumField.DownSamplingFunc())
	assert.Equal(t, function.Min, MinField.DownSamplingFunc())
	assert.Equal(t, function.Max, MaxField.DownSamplingFunc())
	assert.Equal(t, function.Histogram, HistogramField.DownSamplingFunc())
	assert.Equal(t, function.Unknown, Unknown.DownSamplingFunc())
}

func TestType_String(t *testing.T) {
	assert.Equal(t, "sum", SumField.String())
	assert.Equal(t, "max", MaxField.String())
	assert.Equal(t, "min", MinField.String())
	assert.Equal(t, "summary", SummaryField.String())
	assert.Equal(t, "histogram", HistogramField.String())
	assert.Equal(t, "unknown", Unknown.String())
}

func Test_GetPrimitiveFields(t *testing.T) {
	assert.NotNil(t, SumField.GetPrimitiveFields(function.Sum))
	assert.NotNil(t, SumField.GetDefaultPrimitiveFields())
	assert.Nil(t, Unknown.GetPrimitiveFields(function.FuncType(128)))
	assert.Nil(t, Unknown.GetDefaultPrimitiveFields())
}

func TestIsSupportFunc(t *testing.T) {
	assert.True(t, SumField.IsFuncSupported(function.Sum))
	assert.True(t, SumField.IsFuncSupported(function.Min))
	assert.True(t, SumField.IsFuncSupported(function.Max))
	assert.False(t, SumField.IsFuncSupported(function.Histogram))

	assert.True(t, MaxField.IsFuncSupported(function.Max))
	assert.False(t, MaxField.IsFuncSupported(function.Histogram))

	assert.True(t, MinField.IsFuncSupported(function.Min))
	assert.False(t, MinField.IsFuncSupported(function.Histogram))

	assert.True(t, HistogramField.IsFuncSupported(function.Min))
	assert.True(t, HistogramField.IsFuncSupported(function.Sum))
	assert.True(t, HistogramField.IsFuncSupported(function.Max))
	assert.True(t, HistogramField.IsFuncSupported(function.Histogram))

	assert.False(t, Unknown.IsFuncSupported(function.Histogram))
}
