package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// OneSecond is the number of millisecond for a second
	OneSecond int64 = 1000
	// OneMinute is the number of millisecond for a minute
	OneMinute = 60 * OneSecond
	// OneHour is the number of millisecond for an hour
	OneHour = 60 * OneMinute
	// OneDay is the number of millisecond for a day
	OneDay = 24 * OneHour
	// OneWeek is the number of millisecond for a week
	OneWeek = 7 * OneDay
	// OneMonth is the number of millisecond for a month
	OneMonth = 31 * OneDay
	// OneYear is the number of millisecond for a year
	OneYear = 365 * OneDay
	//TODO ????
	dataTimeFormat1 = "20060102 15:04:05"
	dataTimeFormat2 = "2006-01-02 15:04:05"
	dataTimeFormat3 = "2006/01/02 15:04:05"
)

// FormatTimestamp returns timestamp format based on layout
func FormatTimestamp(timestamp int64, layout string) string {
	t := time.Unix(timestamp/1000, 0)
	return t.Format(layout)
}

// ParseTimestamp parses timestamp str value based on layout using local zone
func ParseTimestamp(timestampStr string, layout ...string) (int64, error) {
	var format string
	if len(layout) > 0 {
		format = layout[0]
	} else {
		switch {
		case strings.Index(timestampStr, "-") > 0:
			format = dataTimeFormat2
		case strings.Index(timestampStr, "/") > 0:
			format = dataTimeFormat3
		default:
			format = dataTimeFormat1
		}
	}
	tm, err := time.ParseInLocation(format, timestampStr, time.Local)
	if err != nil {
		return 0, err
	}
	return tm.UnixNano() / 1000000, nil
}

// Now returns t as a Unix time, the number of millisecond elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func Now() int64 {
	return time.Now().UnixNano() / 1000000
}

// Truncate truncates timestamp based on interval
func Truncate(timestamp, interval int64) int64 {
	return timestamp / interval * interval
}

// CalPointCount calculates point counts between start time and end time by interval
func CalPointCount(startTime, endTime, interval int64) int {
	diff := endTime - startTime
	pointCount := diff / interval
	if diff%interval > 0 {
		pointCount++
	}
	if pointCount == 0 {
		pointCount = 1
	}
	return int(pointCount)
}

// CalIntervalRatio calculates the interval ratio for query,
// if query interval < storage interval return 1.
func CalIntervalRatio(queryInterval, storageInterval int64) int {
	if queryInterval < storageInterval {
		return 1
	}
	return int(queryInterval / storageInterval)
}

// ParseInterval parses the interval str, return number of interval(millisecond),
// if parse fail, return 0 and err
func ParseInterval(intervalStr string) (int64, error) {
	var unit, interval int64
	var unitStr string
	switch {
	case strings.HasSuffix(intervalStr, "s"):
		unitStr = "s"
		unit = OneSecond
	case strings.HasSuffix(intervalStr, "S"):
		unitStr = "S"
		unit = OneSecond
	case strings.HasSuffix(intervalStr, "m"):
		unitStr = "m"
		unit = OneMinute
	case strings.HasSuffix(intervalStr, "h"):
		unitStr = "h"
		unit = OneHour
	case strings.HasSuffix(intervalStr, "H"):
		unitStr = "H"
		unit = OneHour
	case strings.HasSuffix(intervalStr, "d"):
		unitStr = "d"
		unit = OneDay
	case strings.HasSuffix(intervalStr, "D"):
		unitStr = "D"
		unit = OneDay
	case strings.HasSuffix(intervalStr, "M"):
		unitStr = "M"
		unit = OneMonth
	case strings.HasSuffix(intervalStr, "y"):
		unitStr = "y"
		unit = OneYear
	case strings.HasSuffix(intervalStr, "Y"):
		unitStr = "Y"
		unit = OneYear
	default:
		return 0, fmt.Errorf("unknown interval")
	}
	intervalStr = strings.Replace(intervalStr, unitStr, "", 1)
	intervalStr = strings.Trim(intervalStr, " ")
	interval, err := strconv.ParseInt(intervalStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return interval * unit, nil
}
