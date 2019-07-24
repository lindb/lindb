package timeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseTimestamp(t *testing.T) {
	date := "20191212 10:10:10"
	_, err := ParseTimestamp(date)
	assert.Nil(t, err)

	date = "2019-12-12 10:10:10"
	_, err = ParseTimestamp(date)
	assert.Nil(t, err)

	date = "2019/12/12 10:10:10"
	_, err = ParseTimestamp(date)
	assert.Nil(t, err)
}
