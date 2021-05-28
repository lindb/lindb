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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalcSlot(t *testing.T) {
	now, _ := ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	t1, _ := ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
	calc := dayCalculator
	assert.Equal(t, 64, calc.CalcSlot(now, t1, 10000))
	assert.Equal(t, 10, calc.CalcSlot(now, t1, 60000))

	now, _ = ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	t1, _ = ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc = monthCalculator
	assert.Equal(t, 19, calc.CalcSlot(now, t1, OneHour))
	assert.Equal(t, 19*12+2, calc.CalcSlot(now, t1, 60000*5))

	now, _ = ParseTimestamp("20190710 19:10:48", "20060102 15:04:05")
	t1, _ = ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc = yearCalculator
	assert.Equal(t, 9, calc.CalcSlot(now, t1, OneDay))
}
func TestGetSegment(t *testing.T) {
	t2, _ := ParseTimestamp("02/07/2019", "02/01/2006")
	calc := dayCalculator
	assert.Equal(t, "20190702", calc.GetSegment(t2))
	calc = monthCalculator
	assert.Equal(t, "201907", calc.GetSegment(t2))
	calc = yearCalculator
	assert.Equal(t, "2019", calc.GetSegment(t2))
}

func TestCalSegment(t *testing.T) {
	t2, _ := ParseTimestamp("20190702", "20060102")
	calc := dayCalculator
	t1, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t2, t1)

	t2, _ = ParseTimestamp("201907", "200601")
	calc = monthCalculator
	t1, _ = calc.ParseSegmentTime("201907")
	assert.Equal(t, t2, t1)

	t2, _ = ParseTimestamp("2019", "2006")
	calc = yearCalculator
	t1, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t2, t1)
}

func TestCalcSegmentTime(t *testing.T) {
	now, _ := ParseTimestamp("20190702 12:30:30", "20060102 15:04:05")

	t1, _ := ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc := dayCalculator
	assert.Equal(t, t1, calc.CalcSegmentTime(now))

	t1, _ = ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc = monthCalculator
	assert.Equal(t, t1, calc.CalcSegmentTime(now))

	t1, _ = ParseTimestamp("20190101 00:00:00", "20060102 15:04:05")
	calc = yearCalculator
	assert.Equal(t, t1, calc.CalcSegmentTime(now))
}

func TestCalcFamily(t *testing.T) {
	now, _ := ParseTimestamp("20190702 12:30:30", "20060102 15:04:05")

	t1, _ := ParseTimestamp("20190702 00:00:00", "20060102 15:04:05")
	calc := dayCalculator
	t2, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 12, calc.CalcFamily(now, t2))

	t1, _ = ParseTimestamp("20190701 00:00:00", "20060102 15:04:05")
	calc = monthCalculator
	t2, _ = calc.ParseSegmentTime("201907")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 2, calc.CalcFamily(now, t2))

	t1, _ = ParseTimestamp("20190101 00:00:00", "20060102 15:04:05")
	calc = yearCalculator
	t2, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t1, t2)
	assert.Equal(t, 7, calc.CalcFamily(now, t2))
}

func TestCalcFamilyTimeRange(t *testing.T) {
	t1, _ := ParseTimestamp("20190702 12:00:00", "20060102 15:04:05")
	calc := dayCalculator
	t2, _ := calc.ParseSegmentTime("20190702")
	assert.Equal(t, t1, calc.CalcFamilyStartTime(t2, 12))
	t3, _ := ParseTimestamp("20190702 13:00:00", "20060102 15:04:05")
	assert.Equal(t, t3-1, calc.CalcFamilyEndTime(t1))

	t1, _ = ParseTimestamp("20191231 00:00:00", "20060102 15:04:05")
	calc = monthCalculator
	t2, _ = calc.ParseSegmentTime("201912")
	assert.Equal(t, t1, calc.CalcFamilyStartTime(t2, 31))
	t3, _ = ParseTimestamp("20200101 00:00:00", "20060102 15:04:05")
	assert.Equal(t, t3-1, calc.CalcFamilyEndTime(t1))

	t1, _ = ParseTimestamp("20191201 00:00:00", "20060102 15:04:05")
	calc = yearCalculator
	t2, _ = calc.ParseSegmentTime("2019")
	assert.Equal(t, t1, calc.CalcFamilyStartTime(t2, 12))
	t3, _ = ParseTimestamp("20200101 00:00:00", "20060102 15:04:05")
	assert.Equal(t, t3-1, calc.CalcFamilyEndTime(t1))
}

func TestCalTimeWindow(t *testing.T) {
	calc := dayCalculator
	assert.Equal(t, 2, calc.CalcTimeWindows(3600000, 3600000*2))

	calc = monthCalculator
	assert.Equal(t, 2, calc.CalcTimeWindows(86400000, 86400000*2))

	calc = yearCalculator
	assert.Equal(t, 2, calc.CalcTimeWindows(2592000000, 2592000000*2))
}

func TestCalcTimestamp(t *testing.T) {
	var i Interval

	_ = i.ValueOf("1m")
	now := Now()
	n := time.Now()
	timestamp := CalcTimestamp(now, 10, i)
	assert.True(t, timestamp >= n.Add(10*time.Minute).Unix()*1000)
}
