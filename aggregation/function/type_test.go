package function

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncTypeString(t *testing.T) {
	assert.Equal(t, "sum", Sum.String())
	assert.Equal(t, "min", Min.String())
	assert.Equal(t, "max", Max.String())
	assert.Equal(t, "count", Count.String())
	assert.Equal(t, "avg", Avg.String())
	assert.Equal(t, "replace", Replace.String())
	assert.Equal(t, "histogram", Histogram.String())
	assert.Equal(t, "stddev", Stddev.String())
	assert.Equal(t, "unknown", Unknown.String())
}
