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
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const date = "20191212 10:11:10"

func Test_ParseTimestamp(t *testing.T) {
	defer func() {
		parseTimeFunc = time.ParseInLocation
	}()
	_, err := ParseTimestamp(date)
	assert.Nil(t, err)

	_, err = ParseTimestamp(date)
	assert.Nil(t, err)

	_, err = ParseTimestamp(date)
	assert.Nil(t, err)
	_, err = ParseTimestamp("2019-12-12 10:11:10")
	assert.Nil(t, err)
	_, err = ParseTimestamp("2019/12/12 10:11:10")
	assert.Nil(t, err)

	parseTimeFunc = func(layout, value string, loc *time.Location) (t time.Time, err error) {
		return time.Now(), fmt.Errorf("err")
	}
	_, err = ParseTimestamp(date)
	assert.Error(t, err)
}

func TestCalPointCount(t *testing.T) {
	time1, _ := ParseTimestamp(date)
	assert.Equal(t, 1, CalPointCount(time1, time1, 10*OneSecond))
	assert.Equal(t, 10, CalPointCount(time1, time1+47*OneSecond, 5*OneSecond))
	assert.Equal(t, 100, CalPointCount(time1, time1+1000*OneSecond, 10*OneSecond))
}

func TestCalIntervalRatio(t *testing.T) {
	assert.Equal(t, 1, CalIntervalRatio(10, 100))
	assert.Equal(t, 1, CalIntervalRatio(10, 0))
	assert.Equal(t, 5, CalIntervalRatio(55, 10))
	assert.Equal(t, 10, CalIntervalRatio(1000, 100))
}

func Test_Now(t *testing.T) {
	assert.Len(t, strconv.FormatUint(uint64(Now()), 10), 13)
	assert.Len(t, strconv.FormatUint(uint64(NowNano()), 10), 19)
}

func Test_FormatTimestamp(t *testing.T) {
	fmt.Println(FormatTimestamp(Now()*1000, dataTimeFormat2))
}

func TestTruncate(t *testing.T) {
	now, _ := ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	t1, _ := ParseTimestamp("20190702 19:10:40", "20060102 15:04:05")
	assert.Equal(t, t1, Truncate(now, 10*OneSecond))
	t1, _ = ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
	assert.Equal(t, t1, Truncate(now, 10*OneMinute))
}
