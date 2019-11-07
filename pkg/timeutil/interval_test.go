package timeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IntervalType_String(t *testing.T) {
	assert.Equal(t, "day", Day.String())
}

func Test_Interval_ValueOf(t *testing.T) {
	var i Interval

	assert.NotNil(t, i.ValueOf(" "))

	assert.NotNil(t, i.ValueOf("10t"))

	assert.NotNil(t, i.ValueOf("as"))

	assert.Nil(t, i.ValueOf(" 10 s"))
	assert.Equal(t, 10*OneSecond, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 S"))
	assert.Equal(t, 10*OneSecond, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 m"))
	assert.Equal(t, 10*OneMinute, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 h"))
	assert.Equal(t, 10*OneHour, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 H"))
	assert.Equal(t, 10*OneHour, i.Int64())

	assert.Nil(t, i.ValueOf(" 10d"))
	assert.Equal(t, 10*OneDay, i.Int64())

	assert.Nil(t, i.ValueOf(" 10D"))
	assert.Equal(t, 10*OneDay, i.Int64())

	assert.Nil(t, i.ValueOf(" 10M"))
	assert.Equal(t, 10*OneMonth, i.Int64())

	assert.Nil(t, i.ValueOf(" 10y"))
	assert.Equal(t, 10*OneYear, i.Int64())

	assert.Nil(t, i.ValueOf(" 10Y"))
	assert.Equal(t, 10*OneYear, i.Int64())
}

func Test_IntervalCalculator(t *testing.T) {
	var i Interval

	_ = i.ValueOf("30m")
	assert.NotNil(t, i.Calculator())

	_ = i.ValueOf("1m")
	assert.NotNil(t, i.Calculator())

	_ = i.ValueOf("10d")
	assert.NotNil(t, i.Calculator())
}
