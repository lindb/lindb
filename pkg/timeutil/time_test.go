package timeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const date = "20191212 10:11:10"

func Test_ParseTimestamp(t *testing.T) {
	_, err := ParseTimestamp(date)
	assert.Nil(t, err)

	_, err = ParseTimestamp(date)
	assert.Nil(t, err)

	_, err = ParseTimestamp(date)
	assert.Nil(t, err)
}

func TestCalPointCount(t *testing.T) {
	time, _ := ParseTimestamp(date)
	assert.Equal(t, 1, CalPointCount(time, time, 10*OneSecond))
	assert.Equal(t, 10, CalPointCount(time, time+47*OneSecond, 5*OneSecond))
	assert.Equal(t, 100, CalPointCount(time, time+1000*OneSecond, 10*OneSecond))
}

func TestCalIntervalRatio(t *testing.T) {
	assert.Equal(t, 1, CalIntervalRatio(10, 100))
	assert.Equal(t, 5, CalIntervalRatio(55, 10))
	assert.Equal(t, 10, CalIntervalRatio(1000, 100))
}
