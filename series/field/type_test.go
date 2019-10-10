package field

import (
	"testing"

	"github.com/lindb/lindb/aggregation/function"

	"github.com/stretchr/testify/assert"
)

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
