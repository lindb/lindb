package interval

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
)

func TestRegister(t *testing.T) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			errStr := r.(string)
			err = errors.New(errStr)
		}
		assert.NotNil(t, err)
	}()

	register(Day, &day{})
}

func TestInit(t *testing.T) {
	calc := GetCalculator(Day)
	assert.NotNil(t, calc)
	calc = GetCalculator(Month)
	assert.NotNil(t, calc)
	calc = GetCalculator(Year)
	assert.NotNil(t, calc)

	calc = GetCalculator(Unknown)
	_, ok := calc.(*day)
	assert.True(t, ok)
	assert.NotNil(t, calc)
}

func TestCalcSlot(t *testing.T) {
	now, _ := timeutil.ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	t1, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
	calc := GetCalculator(Day)
	assert.Equal(t, 64, calc.CalcSlot(now, t1, 10000))
	assert.Equal(t, 10, calc.CalcSlot(now, t1, 60000))

	now, _ = timeutil.ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	t1, _ = timeutil.ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc = GetCalculator(Month)
	assert.Equal(t, 19, calc.CalcSlot(now, t1, timeutil.OneHour))
	assert.Equal(t, 19*12+2, calc.CalcSlot(now, t1, 60000*5))

	now, _ = timeutil.ParseTimestamp("20190710 19:10:48", "20060102 15:04:05")
	t1, _ = timeutil.ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc = GetCalculator(Year)
	assert.Equal(t, 9, calc.CalcSlot(now, t1, timeutil.OneDay))
}
func TestGetSegment(t *testing.T) {
	t2, _ := timeutil.ParseTimestamp("02/07/2019", "02/01/2006")
	calc := GetCalculator(Day)
	assert.Equal(t, "20190702", calc.GetSegment(t2))
	calc = GetCalculator(Month)
	assert.Equal(t, "201907", calc.GetSegment(t2))
	calc = GetCalculator(Year)
	assert.Equal(t, "2019", calc.GetSegment(t2))
}

func TestCalSegment(t *testing.T) {
	t2, _ := timeutil.ParseTimestamp("20190702", "20060102")
	calc := GetCalculator(Day)
	t1, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t2, t1)

	t2, _ = timeutil.ParseTimestamp("201907", "200601")
	calc = GetCalculator(Month)
	t1, _ = calc.ParseSegmentTime("201907")
	assert.Equal(t, t2, t1)

	t2, _ = timeutil.ParseTimestamp("2019", "2006")
	calc = GetCalculator(Year)
	t1, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t2, t1)
}

func TestCalcSegmentTime(t *testing.T) {
	now, _ := timeutil.ParseTimestamp("20190702 12:30:30", "20060102 15:04:05")

	t1, _ := timeutil.ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc := GetCalculator(Day)
	assert.Equal(t, t1, calc.CalcSegmentTime(now))

	t1, _ = timeutil.ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc = GetCalculator(Month)
	assert.Equal(t, t1, calc.CalcSegmentTime(now))

	t1, _ = timeutil.ParseTimestamp("20190101 00:00:00", "20060102 15:04:05")
	calc = GetCalculator(Year)
	assert.Equal(t, t1, calc.CalcSegmentTime(now))
}

func TestCalcFamily(t *testing.T) {
	now, _ := timeutil.ParseTimestamp("20190702 12:30:30", "20060102 15:04:05")

	t1, _ := timeutil.ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc := GetCalculator(Day)
	t2, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 12, calc.CalcFamily(now, t2))

	t1, _ = timeutil.ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc = GetCalculator(Month)
	t2, _ = calc.ParseSegmentTime("201907")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 2, calc.CalcFamily(now, t2))

	t1, _ = timeutil.ParseTimestamp("20190101 00:00:00", "20060102 15:04:05")
	calc = GetCalculator(Year)
	t2, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 7, calc.CalcFamily(now, t2))
}

func TestCalcFamilyTimeRange(t *testing.T) {
	t1, _ := timeutil.ParseTimestamp("20190702 12:00:00", "20060102 15:04:05")
	calc := GetCalculator(Day)
	t2, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t1, calc.CalcFamilyStartTime(t2, 12))
	t3, _ := timeutil.ParseTimestamp("20190702 13:00:00", "20060102 15:04:05")
	assert.Equal(t, t3-1, calc.CalcFamilyEndTime(t1))

	t1, _ = timeutil.ParseTimestamp("20191231 00:00:00", "20060102 15:04:05")
	calc = GetCalculator(Month)
	t2, _ = calc.ParseSegmentTime("201912")
	assert.Equal(t, t1, calc.CalcFamilyStartTime(t2, 31))
	t3, _ = timeutil.ParseTimestamp("20200101 00:00:00", "20060102 15:04:05")
	assert.Equal(t, t3-1, calc.CalcFamilyEndTime(t1))

	t1, _ = timeutil.ParseTimestamp("20191201 00:00:00", "20060102 15:04:05")
	calc = GetCalculator(Year)
	t2, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t1, calc.CalcFamilyStartTime(t2, 12))
	t3, _ = timeutil.ParseTimestamp("20200101 00:00:00", "20060102 15:04:05")
	assert.Equal(t, t3-1, calc.CalcFamilyEndTime(t1))
}

func TestCalcIntervalType(t *testing.T) {
	assert.Equal(t, Year, CalcIntervalType(timeutil.OneHour))
	assert.Equal(t, Month, CalcIntervalType(5*timeutil.OneMinute))
	assert.Equal(t, Day, CalcIntervalType(10*timeutil.OneSecond))
}
