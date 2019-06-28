package util

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_ValueOf(t *testing.T) {
	assert.Equal(t, "sum", ValueOf("sum").String())
	assert.Equal(t, "avg", ValueOf("d_avg").String())
}

func Test_IsDownSamplingOrAggregator(t *testing.T) {
	assert.True(t, true, IsDownSamplingOrAggregator("d_avg") == true)
	assert.True(t, true, IsDownSamplingOrAggregator("avg") == false)
	assert.True(t, true, IsDownSamplingOrAggregator("sum") == true)
}
