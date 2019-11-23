package timeutil

import (
	"fmt"
	"strconv"
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
	assert.Equal(t, 1, CalIntervalRatio(10, 0))
	assert.Equal(t, 5, CalIntervalRatio(55, 10))
	assert.Equal(t, 10, CalIntervalRatio(1000, 100))
}

func Test_Now(t *testing.T) {
	assert.Len(t, strconv.FormatUint(uint64(Now()), 10), 13)
}

func Test_FormatTimestamp(t *testing.T) {
	fmt.Println(FormatTimestamp(Now()*1000, dataTimeFormat2))
}

func TestTruncate(t *testing.T) {
	now, _ := ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	t1, _ := ParseTimestamp("20190702 19:10:40", "20060102 15:04:05")
	assert.Equal(t, t1, Truncate(now, 10*OneSecond))
	t1, _ = ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
	assert.Equal(t, t1, Truncate(now, 10*OneMinute))
}
