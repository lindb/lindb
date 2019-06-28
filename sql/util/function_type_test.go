package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FunctionType(t *testing.T) {
	assert.Equal(t, "sum", SUM.String())
	assert.Equal(t, "histogram", HISTOGRAM.String())
	assert.Equal(t, "sum", FunctionType(1).String())
	assert.Equal(t, "avg", FunctionType(5).String())
	assert.Equal(t, "min", GetFunctionType("min").String())
}
