package point

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MilliSecondOf(t *testing.T) {
	assert.Equal(t, int64(0), MilliSecondOf(0))
	assert.Equal(t, int64(1576808994000), MilliSecondOf(1576808994))
	assert.Equal(t, int64(1576808994000), MilliSecondOf(1576808994000))
	assert.Equal(t, int64(1576808994000), MilliSecondOf(1576808994000000))
	assert.Equal(t, int64(1576808994000), MilliSecondOf(1576808994000000000))
}
