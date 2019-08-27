package timeutil

import (
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
	assert.Equal(t, 5, CalIntervalRatio(55, 10))
	assert.Equal(t, 10, CalIntervalRatio(1000, 100))
}

func TestParseInterval(t *testing.T) {
	_, err := ParseInterval("10t")
	assert.NotNil(t, err)

	_, err = ParseInterval("as")
	assert.NotNil(t, err)
	interval, err := ParseInterval(" 10  s")
	assert.Equal(t, 10*OneSecond, interval)
	assert.Nil(t, err)
	interval, _ = ParseInterval(" 10  S")
	assert.Equal(t, 10*OneSecond, interval)
	interval, _ = ParseInterval("10m")
	assert.Equal(t, 10*OneMinute, interval)
	interval, _ = ParseInterval("10h")
	assert.Equal(t, 10*OneHour, interval)
	interval, _ = ParseInterval("10H")
	assert.Equal(t, 10*OneHour, interval)
	interval, _ = ParseInterval("10d")
	assert.Equal(t, 10*OneDay, interval)
	interval, _ = ParseInterval("10D")
	assert.Equal(t, 10*OneDay, interval)
	interval, _ = ParseInterval("10M")
	assert.Equal(t, 10*OneMonth, interval)
	interval, _ = ParseInterval("10y")
	assert.Equal(t, 10*OneYear, interval)
	interval, _ = ParseInterval("10Y")
	assert.Equal(t, 10*OneYear, interval)
}

func Test_Now(t *testing.T) {
	assert.Len(t, strconv.FormatUint(uint64(Now()), 10), 13)
}

func Test_FormatTimestamp(t *testing.T) {
	t.Log(FormatTimestamp(Now()*1000, dataTimeFormat1))
}
