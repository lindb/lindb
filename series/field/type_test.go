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
	assert.NotNil(t, GetPrimitiveFields(SumField, function.Sum))
	assert.Nil(t, GetPrimitiveFields(Type(128), function.FuncType(128)))
	GetPrimitiveFieldsValue()
}
