package util

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_ParseTimestamp(t *testing.T) {
	date := "20191212 10:10:10"
	timestamp := ParseTimestamp(date)
	assert.Equal(t, int64(1576116610000), timestamp)

	date = "2019-12-12 10:10:10"
	timestamp = ParseTimestamp(date)
	assert.Equal(t, int64(1576116610000), timestamp)

	date = "2019/12/12 10:10:10"
	timestamp = ParseTimestamp(date)
	assert.Equal(t, int64(1576116610000), timestamp)
}
