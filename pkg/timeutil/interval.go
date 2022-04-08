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
	case val >= OneYear:
		return fmt.Sprintf("%dy", val/OneYear)
	case val >= OneMonth:
		return fmt.Sprintf("%dM", val/OneMonth)
	case val >= OneDay:
		return fmt.Sprintf("%dd", val/OneDay)
	case val >= OneHour:
		return fmt.Sprintf("%dh", val/OneHour)
	case val >= OneMinute:
		return fmt.Sprintf("%dm", val/OneMinute)
	default:
		return fmt.Sprintf("%ds", val/OneSecond)
	}
}

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
	if err := jsoniter.Unmarshal(data, &v); err != nil {
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

// CalcQueryInterval returns query interval based on query time range and interval.
func CalcQueryInterval(queryTimeRange TimeRange, queryInterval Interval) Interval {
	diff := queryTimeRange.End - queryTimeRange.Start
	switch {
	case diff < OneHour:
		return queryInterval
	case diff < 3*OneHour:
		return Interval(10 * OneSecond)
	case diff < 6*OneHour:
		return Interval(30 * OneSecond)
	case diff < 12*OneHour:
		return Interval(OneMinute)
	case diff < OneDay:
		return Interval(2 * OneMinute)
	case diff < 2*OneDay:
		return Interval(5 * OneMinute)
	case diff < 7*OneDay:
		return Interval(10 * OneMinute)
	case diff < OneMonth:
		return Interval(OneHour)
	case diff < 2*OneMonth:
		return Interval(4 * OneHour)
	case diff < 3*OneMonth:
		return Interval(12 * OneHour)
	default:
		return Interval(OneDay)
	}
}
