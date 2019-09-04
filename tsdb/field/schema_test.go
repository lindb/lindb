package field

import (
	"testing"

	"github.com/lindb/lindb/aggregation/function"

	"github.com/stretchr/testify/assert"
)

func Test_getPrimitiveFields(t *testing.T) {
	assert.NotNil(t, newSumSchema().getPrimitiveFields(function.Sum))

	assert.Nil(t, newSumSchema().getPrimitiveFields(function.FuncType(128)))
}
