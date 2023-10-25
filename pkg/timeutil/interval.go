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
	"errors"
	"fmt"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/lindb/common/pkg/timeutil"
)

var (
	unmarshalFn = jsoniter.Unmarshal
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

// String returns the string representation of the interval.
func (i Interval) String() string {
	val := i.Int64()
	switch {
	case val%timeutil.OneYear == 0 && val/timeutil.OneYear > 0:
		return fmt.Sprintf("%dy", val/timeutil.OneYear)
	case val%timeutil.OneMonth == 0 && val/timeutil.OneMonth > 0:
		return fmt.Sprintf("%dM", val/timeutil.OneMonth)
	case val%timeutil.OneDay == 0 && val/timeutil.OneDay > 0:
		return fmt.Sprintf("%dd", val/timeutil.OneDay)
	case val%timeutil.OneHour == 0 && val/timeutil.OneHour > 0:
		return fmt.Sprintf("%dh", val/timeutil.OneHour)
	case val%timeutil.OneMinute == 0 && val/timeutil.OneMinute > 0:
		return fmt.Sprintf("%dm", val/timeutil.OneMinute)
	default:
		return fmt.Sprintf("%ds", val/timeutil.OneSecond)
	}
}

// ValueOf parses the interval str, return number of interval(millisecond),
func (i *Interval) ValueOf(intervalStr string) error {
	intervalBytes := []byte(strings.ReplaceAll(intervalStr, " ", ""))
	if len(intervalBytes) <= 1 {
		return ErrUnknownInterval
	}
	unixSuffix := string(intervalBytes[len(intervalBytes)-1])
	valuePrefix := string(intervalBytes[:len(intervalBytes)-1])

	var unit int64
	switch unixSuffix {
	case "s", "S":
		unit = timeutil.OneSecond
	case "m":
		unit = timeutil.OneMinute
	case "h", "H":
		unit = timeutil.OneHour
	case "d", "D":
		unit = timeutil.OneDay
	case "M":
		unit = timeutil.OneMonth
	case "y", "Y":
		unit = timeutil.OneYear
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

// UnmarshalText parses a TOML value into an interval value.
// See https://github.com/BurntSushi/toml
func (i *Interval) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	return i.ValueOf(string(text))
}

// MarshalText converts an interval to a string for decoding toml
func (i Interval) MarshalText() (text []byte, err error) {
	return []byte(i.String()), nil
}

// UnmarshalJSON parses a JSON value into an interval value.
func (i *Interval) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err := unmarshalFn(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		return i.ValueOf(value)
	default:
		return errors.New("invalid interval")
	}
}

// MarshalJSON converts an interval to a string for decoding json
func (i *Interval) MarshalJSON() (data []byte, err error) {
	return jsoniter.Marshal(i.String())
}

func (i Interval) Int64() int64 {
	return int64(i)
}

func (i Interval) Type() IntervalType {
	switch {
	case i.Int64() >= timeutil.OneHour:
		return Year
	case i.Int64() >= 5*timeutil.OneMinute:
		return Month
	default:
		return Day
	}
}

// Calculator returns the calculator for current interval.
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

// CalcSlotRange returns slot range by family time and time range.
func (i Interval) CalcSlotRange(familyTime int64, timeRange TimeRange) SlotRange {
	calc := i.Calculator()
	storageTimeRange := TimeRange{
		Start: familyTime,
		End:   calc.CalcFamilyEndTime(familyTime),
	}
	rs := timeRange.Intersect(storageTimeRange)
	intervalVal := i.Int64()
	return SlotRange{
		Start: uint16(calc.CalcSlot(rs.Start, familyTime, intervalVal)),
		End:   uint16(calc.CalcSlot(rs.End, familyTime, intervalVal)),
	}
}

// CalcQueryInterval returns query interval based on query time range and interval.
func CalcQueryInterval(queryTimeRange TimeRange, queryInterval Interval) Interval {
	diff := queryTimeRange.End - queryTimeRange.Start
	switch {
	case diff < timeutil.OneHour:
		return queryInterval
	case diff < 3*timeutil.OneHour:
		return Interval(10 * timeutil.OneSecond)
	case diff < 6*timeutil.OneHour:
		return Interval(30 * timeutil.OneSecond)
	case diff < 12*timeutil.OneHour:
		return Interval(timeutil.OneMinute)
	case diff < timeutil.OneDay:
		return Interval(2 * timeutil.OneMinute)
	case diff < 2*timeutil.OneDay:
		return Interval(5 * timeutil.OneMinute)
	case diff < 7*timeutil.OneDay:
		return Interval(10 * timeutil.OneMinute)
	case diff < timeutil.OneMonth:
		return Interval(timeutil.OneHour)
	case diff < 2*timeutil.OneMonth:
		return Interval(4 * timeutil.OneHour)
	case diff < 3*timeutil.OneMonth:
		return Interval(12 * timeutil.OneHour)
	default:
		return Interval(timeutil.OneDay)
	}
}
