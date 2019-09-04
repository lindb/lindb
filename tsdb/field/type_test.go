package field

import (
	"testing"

	"github.com/lindb/lindb/aggregation/function"

	"github.com/stretchr/testify/assert"
)

func Test_GetPrimitiveFields(t *testing.T) {
	assert.NotNil(t, GetPrimitiveFields(SumField, function.Sum))

	assert.Nil(t, GetPrimitiveFields(Type(128), function.FuncType(128)))

	GetPrimitiveFieldsValue()
}
