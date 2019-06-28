package util

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_GetStringValue(t *testing.T) {
	assert.Equal(t, "sum", GetStringValue("sum"))
	assert.Equal(t, "sum", GetStringValue("'sum'"))
	assert.Equal(t, "'sum", GetStringValue("'sum"))
	assert.Equal(t, "sum", GetStringValue("\"sum\""))
	assert.Equal(t, "", GetStringValue(""))
}
