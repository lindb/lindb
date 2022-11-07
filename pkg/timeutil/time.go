// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package timeutil

import (
	"strings"
	"time"
)

var (
	parseTimeFunc = time.ParseInLocation
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
	OneMonth = 30 * OneDay
	// OneYear is the number of millisecond for a year
	OneYear = 365 * OneDay

	DataTimeFormat1 = "20060102 15:04:05"
	DataTimeFormat2 = "2006-01-02 15:04:05"
	DataTimeFormat3 = "2006/01/02 15:04:05"
	DataTimeFormat4 = "20060102150405"
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
			format = DataTimeFormat2
		case strings.Index(timestampStr, "/") > 0:
			format = DataTimeFormat3
		case strings.Index(timestampStr, " ") > 0:
			format = DataTimeFormat1
		default:
			format = DataTimeFormat4
		}
	}
	tm, err := parseTimeFunc(format, timestampStr, time.Local)
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

// NowNano returns t as a unix time, the number of nanoseconds elapsed
func NowNano() int64 {
	return time.Now().UnixNano()
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
	if storageInterval == 0 || queryInterval < storageInterval {
		return 1
	}
	return int(queryInterval / storageInterval)
}
