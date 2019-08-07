package function

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncTypeString(t *testing.T) {
	assert.Equal(t, "sum", FuncTypeString(Sum))
	assert.Equal(t, "min", FuncTypeString(Min))
	assert.Equal(t, "max", FuncTypeString(Max))
	assert.Equal(t, "avg", FuncTypeString(Avg))
	assert.Equal(t, "histogram", FuncTypeString(Histogram))
	assert.Equal(t, "stddev", FuncTypeString(Stddev))
	assert.Equal(t, "unknown", FuncTypeString(Unknown))
}
