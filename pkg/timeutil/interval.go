package timeutil

import (
	"errors"
	"strconv"
	"strings"
)

// IntervalType defines interval type
type IntervalType string

// String implements stringer
func (t IntervalType) String() string {
	return string(t)
}

// Interval types.
const (
	Day     IntervalType = "day"
	Month   IntervalType = "month"
	Year    IntervalType = "year"
	Unknown IntervalType = "unknown"
)

var ErrUnknownInterval = errors.New("unknown interval")

// Interval is the interval value in millisecond
type Interval int64

// ValueOf parses the interval str, return number of interval(millisecond),
func (i *Interval) ValueOf(intervalStr string) error {
	intervalBytes := []byte(strings.Replace(intervalStr, " ", "", -1))
	if len(intervalBytes) <= 1 {
		return ErrUnknownInterval
	}
	unixSuffix := string(intervalBytes[len(intervalBytes)-1])
	valuePrefix := string(intervalBytes[:len(intervalBytes)-1])

	var unit int64
	switch unixSuffix {
	case "s", "S":
		unit = OneSecond
	case "m":
		unit = OneMinute
	case "h", "H":
		unit = OneHour
	case "d", "D":
		unit = OneDay
	case "M":
		unit = OneMonth
	case "y", "Y":
		unit = OneYear
	default:
		return ErrUnknownInterval
	}
	value, err := strconv.ParseInt(valuePrefix, 10, 64)
	if err != nil {
		return ErrUnknownInterval
	}
	*i = Interval(value * unit)
	return nil
}

func (i Interval) Int64() int64 {
	return int64(i)
}

func (i Interval) Type() IntervalType {
	switch {
	case i.Int64() >= OneHour:
		return Year
	case i.Int64() >= 5*OneMinute:
		return Month
	default:
		return Day
	}
}

func (i Interval) Calculator() IntervalCalculator {
	switch i.Type() {
	case Year:
		return yearCalculator
	case Month:
		return monthCalculator
	default:
		return dayCalculator
	}
}
