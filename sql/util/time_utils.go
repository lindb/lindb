package util

import (
	"strings"
	"time"
)

const (
	OneSeconds      int64 = 1000
	OneMinute             = 60 * OneSeconds
	OneHour               = 60 * OneMinute
	OneDay                = 24 * OneHour
	OneWeek               = 7 * OneDay
	OneMonth              = 31 * OneDay
	OneYear               = 365 * OneDay
	dataTimeFormat1       = "20060102 15:04:05"
	dataTimeFormat2       = "2006-01-02 15:04:05"
	dataTimeFormat3       = "2006/01/02 15:04:05"
)

// ParseTimestamp format date to timestamp
func ParseTimestamp(date string) int64 {
	loc, _ := time.LoadLocation("PRC")
	var format string
	switch {
	case strings.Index(date, "-") > 0:
		format = dataTimeFormat2
	case strings.Index(date, "/") > 0:
		format = dataTimeFormat3
	default:
		format = dataTimeFormat1
	}
	t, _ := time.ParseInLocation(format, date, loc)
	timestamp := t.UnixNano() / 1000000
	return timestamp
}

// NowTimestamp get now time to timestamp
func NowTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}
