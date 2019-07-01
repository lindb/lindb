package interval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/util"
)

func TestInit(t *testing.T) {
	assert.NotNil(t, GetCalculator(Day))
	assert.NotNil(t, GetCalculator(Month))
	assert.NotNil(t, GetCalculator(Year))
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

func TestGetSegment(t *testing.T) {
	t2, _ := util.ParseTimestamp("02/07/2019", "02/01/2006")
	assert.Equal(t, "20190702", GetCalculator(Day).GetSegment(t2))
	assert.Equal(t, "201907", GetCalculator(Month).GetSegment(t2))
	assert.Equal(t, "2019", GetCalculator(Year).GetSegment(t2))
}

func TestCalSegment(t *testing.T) {
	t2, _ := util.ParseTimestamp("20190702", "20060102")
	t1, _ := GetCalculator(Day).ParseSegmentTime("20190702")
	assert.Equal(t, t2, t1)

	t2, _ = util.ParseTimestamp("201907", "200601")
	t1, _ = GetCalculator(Month).ParseSegmentTime("201907")
	assert.Equal(t, t2, t1)

	t2, _ = util.ParseTimestamp("2019", "2006")
	t1, _ = GetCalculator(Year).ParseSegmentTime("2019")
	assert.Equal(t, t2, t1)
}
