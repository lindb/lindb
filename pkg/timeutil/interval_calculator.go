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
	"time"
)

var (
	yearCalculator  Calculator = &year{}
	monthCalculator Calculator = &month{}
	dayCalculator   Calculator = &day{}
)

// IntervalCalculator calculates the timestamp for each interval type
type IntervalCalculator interface {
	// GetSegment returns segment name by given timestamp
	GetSegment(timestamp int64) string
	// ParseSegmentTime parses segment base time based on given segment name
	ParseSegmentTime(segmentName string) (int64, error)
	// CalcSegmentTime calculates segment base time based on given segment name
	CalcSegmentTime(timestamp int64) int64
	// CalcFamily calculates family base time based on given timestamp
	CalcFamily(timestamp int64, segmentTime int64) int
	// CalcFamilyStartTime calculates family start time based on segment time and family time
	CalcFamilyStartTime(segmentTime int64, familyTime int) int64
	// CalcFamilyEndTime calculates family end time based on family start time
	CalcFamilyEndTime(familyStartTime int64) int64
	// CalcSlot calculates field store slot index based on given timestamp and base time
	CalcSlot(timestamp, baseTime, interval int64) int
	// CalcTimeWindows calculates the number of time window between start and end time
	CalcTimeWindows(start, end int64) int
}

// Calculator represents calculate timestamp for each interval type
type Calculator interface {
	// GetSegment returns segment name by given timestamp
	GetSegment(timestamp int64) string
	// ParseSegmentTime parses segment base time based on given segment name
	ParseSegmentTime(segmentName string) (int64, error)
	// CalcSegmentTime calculates segment base time based on given segment name
	CalcSegmentTime(timestamp int64) int64
	// CalcFamily calculates family base time based on given timestamp
	CalcFamily(timestamp int64, segmentTime int64) int
	// CalcFamilyStartTime calculates family start time based on segment time and family time
	CalcFamilyStartTime(segmentTime int64, familyTime int) int64
	// CalcFamilyEndTime calculates family end time based on family start time
	CalcFamilyEndTime(familyStartTime int64) int64
	// CalcSlot calculates field store slot index based on given timestamp and base time
	CalcSlot(timestamp, baseTime, interval int64) int
	// CalcTimeWindows calculates the number of time window between start and end time
	CalcTimeWindows(start, end int64) int
}

// day implements Calculator interface for day interval type
type day struct{}

// CalcSlot calculates field store slot index based on given timestamp and base time for day interval type
func (d *day) CalcSlot(timestamp, baseTime, interval int64) int {
	return int(((timestamp - baseTime) % OneHour) / interval)
}

// GetSegment returns segment name by given timestamp for day interval type
func (d *day) GetSegment(timestamp int64) string {
	return FormatTimestamp(timestamp, "20060102")
}

// ParseSegmentTime parses segment base time based on given segment name for day interval type
func (d *day) ParseSegmentTime(segmentName string) (int64, error) {
	return ParseTimestamp(segmentName, "20060102")
}

// CalcSegmentTime calculates segment base time based on given segment name for day interval type
func (d *day) CalcSegmentTime(timestamp int64) int64 {
	t := time.Unix(timestamp/1000, 0)
	t2 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// CalcFamily calculates family base time based on given timestamp for day interval type
func (d *day) CalcFamily(timestamp int64, segmentTime int64) int {
	return int((timestamp - segmentTime) / OneHour)
}

// CalcFamilyStartTime calculates family start time based on segment time and family time for day interval type
func (d *day) CalcFamilyStartTime(segmentTime int64, familyTime int) int64 {
	return segmentTime + int64(familyTime)*OneHour
}

// CalcFamilyEndTime calculates family end time based on family start time for day interval type
func (d *day) CalcFamilyEndTime(familyStartTime int64) int64 {
	return familyStartTime + OneHour - 1
}

// CalcTimeWindows calculates the number of time window between start and end time for day interval type
func (d *day) CalcTimeWindows(start, end int64) int {
	t1 := start / OneHour * OneHour
	t2 := end / OneHour * OneHour
	return int((t2-t1)/OneHour) + 1
}

// month implements Calculator interface for month interval type
type month struct{}

// CalcSlot calculates field store slot index based on given timestamp and base time for month interval type
func (m *month) CalcSlot(timestamp, baseTime, interval int64) int {
	return int(((timestamp - baseTime) % OneDay) / interval)
}

// GetSegment returns segment name by given timestamp for month interval type
func (m *month) GetSegment(timestamp int64) string {
	return FormatTimestamp(timestamp, "200601")
}

// ParseSegmentTime parses segment base time based on given segment name for month interval type
func (m *month) ParseSegmentTime(segmentName string) (int64, error) {
	return ParseTimestamp(segmentName, "200601")
}

// CalcSegmentTime calculates segment base time based on given segment name for month interval type
func (m *month) CalcSegmentTime(timestamp int64) int64 {
	t := time.Unix(timestamp/1000, 0)
	t2 := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// CalcFamily calculates family base time based on given timestamp for month interval type
func (m *month) CalcFamily(timestamp int64, segmentTime int64) int {
	t := time.Unix(timestamp/1000, 0)
	return t.Day()
}

// CalcFamilyStartTime calculates family start time based on segment time and family for month interval type
func (m *month) CalcFamilyStartTime(segmentTime int64, familyTime int) int64 {
	t := time.Unix(segmentTime/1000, 0)
	t2 := time.Date(t.Year(), t.Month(), familyTime, 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// CalcFamilyEndTime calculates family end time based on family start time for month interval type
func (m *month) CalcFamilyEndTime(familyStartTime int64) int64 {
	t := time.Unix(familyStartTime/1000, 0)
	t2 := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, time.Local)
	return t2.UnixNano()/1000000 - 1
}

// CalcTimeWindows calculates the number of time window between start and end time for month interval type
func (m *month) CalcTimeWindows(start, end int64) int {
	t1 := time.Unix(start/1000, 0)
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 := time.Unix(end/1000, 0)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(t2.Sub(t1).Hours()/24) + 1
}

// year implements Calculator interface for year interval type
type year struct{}

// CalcSlot calculates field store slot index based on given timestamp and base time for year interval type
func (y *year) CalcSlot(timestamp, baseTime, interval int64) int {
	return int((timestamp - baseTime) / interval)
}

// GetSegment returns segment name by given timestamp for day interval type
func (y *year) GetSegment(timestamp int64) string {
	return FormatTimestamp(timestamp, "2006")
}

// ParseSegmentTime parses segment base time based on given segment name for year interval type
func (y *year) ParseSegmentTime(segmentName string) (int64, error) {
	return ParseTimestamp(segmentName, "2006")
}

// CalcSegmentTime calculates segment base time based on given segment name for year interval type
func (y *year) CalcSegmentTime(timestamp int64) int64 {
	t := time.Unix(timestamp/1000, 0)
	t2 := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// CalcFamily calculates family base time based on given timestamp for year interval type
func (y *year) CalcFamily(timestamp int64, segmentTime int64) int {
	t := time.Unix(timestamp/1000, 0)
	return int(t.Month())
}

// CalcFamilyStartTime calculates family start time based on segment time and family for year interval type
func (y *year) CalcFamilyStartTime(segmentTime int64, familyTime int) int64 {
	t := time.Unix(segmentTime/1000, 0)
	t2 := time.Date(t.Year(), time.Month(familyTime), 1, 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// CalcFamilyEndTime calculates family end time based on family start time for year interval type
func (y *year) CalcFamilyEndTime(familyStartTime int64) int64 {
	t := time.Unix(familyStartTime/1000, 0)
	t2 := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.Local)
	return t2.UnixNano()/1000000 - 1
}

// CalcTimeWindows calculates the number of time window between start and end time for year interval type
func (y *year) CalcTimeWindows(start, end int64) int {
	t1 := time.Unix(start/1000, 0)
	t1 = time.Date(t1.Year(), t1.Month(), 0, 0, 0, 0, 0, time.Local)
	t2 := time.Unix(end/1000, 0)
	t2 = time.Date(t2.Year(), t2.Month(), 0, 0, 0, 0, 0, time.Local)
	return int(t2.Sub(t1).Hours()/24/30) + 1
}

// CalcTimestamp returns timestamp based on start time, slot and interval.
func CalcTimestamp(startTime int64, slot int, interval Interval) int64 {
	return interval.Int64()*int64(slot) + startTime
}
