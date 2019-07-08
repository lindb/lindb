package interval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/timeutil"
)

func TestInit(t *testing.T) {
	calc, err := GetCalculator(Day)
	assert.Nil(t, err)
	assert.NotNil(t, calc)
	calc, err = GetCalculator(Month)
	assert.Nil(t, err)
	assert.NotNil(t, calc)
	calc, err = GetCalculator(Year)
	assert.Nil(t, err)
	assert.NotNil(t, calc)

	calc, err = GetCalculator(44)
	assert.NotNil(t, err)
	assert.Nil(t, calc)
}

func TestString(t *testing.T) {
	assert.Equal(t, "day", Day.String())
	assert.Equal(t, "month", Month.String())
	assert.Equal(t, "year", Year.String())

	t1, err := ParseType("day")
	assert.Nil(t, err)
	assert.Equal(t, Day, t1)
	t1, err = ParseType("month")
	assert.Nil(t, err)
	assert.Equal(t, Month, t1)
	t1, err = ParseType("year")
	assert.Nil(t, err)
	assert.Equal(t, Year, t1)

	t1, err = ParseType("year111")
	assert.NotNil(t, err)
	assert.Equal(t, Unknown, t1)
}
func TestCalSlot(t *testing.T) {
	now, _ := timeutil.ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	t1, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
	calc, _ := GetCalculator(Day)
	assert.Equal(t, 64, calc.CalSlot(now, t1, 10000))
	assert.Equal(t, 10, calc.CalSlot(now, t1, 60000))

	now, _ = timeutil.ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	t1, _ = timeutil.ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc, _ = GetCalculator(Month)
	assert.Equal(t, 19, calc.CalSlot(now, t1, timeutil.OneHour))
	assert.Equal(t, 19*12+2, calc.CalSlot(now, t1, 60000*5))

	now, _ = timeutil.ParseTimestamp("20190710 19:10:48", "20060102 15:04:05")
	t1, _ = timeutil.ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc, _ = GetCalculator(Year)
	assert.Equal(t, 9, calc.CalSlot(now, t1, timeutil.OneDay))
}
func TestGetSegment(t *testing.T) {
	t2, _ := timeutil.ParseTimestamp("02/07/2019", "02/01/2006")
	calc, _ := GetCalculator(Day)
	assert.Equal(t, "20190702", calc.GetSegment(t2))
	calc, _ = GetCalculator(Month)
	assert.Equal(t, "201907", calc.GetSegment(t2))
	calc, _ = GetCalculator(Year)
	assert.Equal(t, "2019", calc.GetSegment(t2))
}

func TestCalSegment(t *testing.T) {
	t2, _ := timeutil.ParseTimestamp("20190702", "20060102")
	calc, _ := GetCalculator(Day)
	t1, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t2, t1)

	t2, _ = timeutil.ParseTimestamp("201907", "200601")
	calc, _ = GetCalculator(Month)
	t1, _ = calc.ParseSegmentTime("201907")
	assert.Equal(t, t2, t1)

	t2, _ = timeutil.ParseTimestamp("2019", "2006")
	calc, _ = GetCalculator(Year)
	t1, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t2, t1)
}

func TestCalSegmentTime(t *testing.T) {
	now, _ := timeutil.ParseTimestamp("20190702 12:30:30", "20060102 15:04:05")

	t1, _ := timeutil.ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc, _ := GetCalculator(Day)
	assert.Equal(t, t1, calc.CalSegmentTime(now))

	t1, _ = timeutil.ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc, _ = GetCalculator(Month)
	assert.Equal(t, t1, calc.CalSegmentTime(now))

	t1, _ = timeutil.ParseTimestamp("20190101 00:00:00", "20060102 15:04:05")
	calc, _ = GetCalculator(Year)
	assert.Equal(t, t1, calc.CalSegmentTime(now))
}

func TestCalFamily(t *testing.T) {
	now, _ := timeutil.ParseTimestamp("20190702 12:30:30", "20060102 15:04:05")

	t1, _ := timeutil.ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc, _ := GetCalculator(Day)
	t2, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 12, calc.CalFamily(now, t2))

	t1, _ = timeutil.ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc, _ = GetCalculator(Month)
	t2, _ = calc.ParseSegmentTime("201907")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 2, calc.CalFamily(now, t2))

	t1, _ = timeutil.ParseTimestamp("20190101 00:00:00", "20060102 15:04:05")
	calc, _ = GetCalculator(Year)
	t2, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 7, calc.CalFamily(now, t2))
}

func TestCalFamilyStartTime(t *testing.T) {
	t1, _ := timeutil.ParseTimestamp("20190702 12:00:00", "20060102 15:04:05")
	calc, _ := GetCalculator(Day)
	t2, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t1, calc.CalFamilyStartTime(t2, 12))

	t1, _ = timeutil.ParseTimestamp("20190710 00:00:00", "20060102 15:04:05")
	calc, _ = GetCalculator(Month)
	t2, _ = calc.ParseSegmentTime("201907")
	assert.Equal(t, t1, calc.CalFamilyStartTime(t2, 10))

	t1, _ = timeutil.ParseTimestamp("20191001 00:00:00", "20060102 15:04:05")
	calc, _ = GetCalculator(Year)
	t2, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t1, calc.CalFamilyStartTime(t2, 10))
}
