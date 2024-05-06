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
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/timeutil"
	"github.com/stretchr/testify/assert"
)

type retention struct {
	Retention Interval `toml:"retention" json:"retention,omitempty"`
}

func TestInterval_MarshalText(t *testing.T) {
	cases := []struct {
		name   string
		in     Interval
		assert []byte
	}{
		{
			name:   "10s",
			in:     Interval(10 * 1000),
			assert: []byte("10s"),
		},
		{
			name:   "5day",
			in:     Interval(5 * 24 * 60 * 60 * 1000),
			assert: []byte("5d"),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.in.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.assert, val)
		})
	}
}

func TestInterval_UnmarshalText(t *testing.T) {
	cases := []struct {
		name   string
		in     []byte
		assert Interval
	}{
		{
			"10s",
			[]byte("10s"),
			Interval(10 * 1000),
		},
		{
			"5day",
			[]byte("5d"),
			Interval(5 * 24 * 60 * 60 * 1000),
		},
		{
			"empty",
			[]byte(""),
			Interval(0),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var rs Interval
			err := rs.UnmarshalText(tt.in)
			assert.NoError(t, err)
			assert.Equal(t, tt.assert, rs)
		})
	}
}

func TestInterval_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name    string
		in      []byte
		wantErr bool
		assert  retention
	}{
		{
			"json_10s",
			[]byte(`{"retention":"10s"}`),
			false,
			retention{Retention: Interval(10 * 1000)},
		},
		{
			"json_5min",
			[]byte(`{"retention":"5m"}`),
			false,
			retention{Retention: Interval(5 * 60 * 1000)},
		},
		{
			"json_5day",
			[]byte(`{"retention":"5d"}`),
			false,
			retention{Retention: Interval(5 * 24 * 60 * 60 * 1000)},
		},
		{
			"json_5day",
			[]byte(`{"retention":"5d"}`),
			false,
			retention{Retention: Interval(5 * 24 * 60 * 60 * 1000)},
		},
		{
			"invalid interval",
			[]byte(`{"retention":12}`),
			true,
			retention{Retention: Interval(0)},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var rs retention
			err := encoding.JSONUnmarshal(tt.in, &rs)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.assert, rs)
		})
	}

	defer func() {
		unmarshalFn = jsoniter.Unmarshal
	}()

	unmarshalFn = func(data []byte, v interface{}) error {
		return fmt.Errorf("err")
	}
	interval := Interval(10)
	err := (&interval).UnmarshalJSON([]byte("test"))
	assert.Error(t, err)
}

func TestInterval_JSONMarshal(t *testing.T) {
	cases := []struct {
		name   string
		in     retention
		assert []byte
	}{
		{
			name:   "json_10s",
			in:     retention{Retention: Interval(10 * 1000)},
			assert: []byte(`{"retention":"10s"}`),
		},
		{
			name:   "json_5min",
			in:     retention{Retention: Interval(5 * 60 * 1000)},
			assert: []byte(`{"retention":"5m"}`),
		},
		{
			name:   "json_5day",
			in:     retention{Retention: Interval(5 * 24 * 60 * 60 * 1000)},
			assert: []byte(`{"retention":"5d"}`),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			val := encoding.JSONMarshal(&tt.in)
			assert.Equal(t, tt.assert, val)
		})
	}
}

func TestInterval_String(t *testing.T) {
	cases := []struct {
		name   string
		in     Interval
		assert string
	}{
		{
			name:   "10s",
			in:     Interval(10 * 1000),
			assert: "10s",
		},
		{
			name:   "70s",
			in:     Interval(70 * 1000),
			assert: "70s",
		},
		{
			name:   "5min",
			in:     Interval(5 * 60 * 1000),
			assert: "5m",
		},
		{
			name:   "65min",
			in:     Interval(65 * 60 * 1000),
			assert: "65m",
		},
		{
			name:   "5hour",
			in:     Interval(5 * 60 * 60 * 1000),
			assert: "5h",
		},
		{
			name:   "25hour",
			in:     Interval(25 * 60 * 60 * 1000),
			assert: "25h",
		},
		{
			name:   "5day",
			in:     Interval(5 * 24 * 60 * 60 * 1000),
			assert: "5d",
		},
		{
			name:   "35day",
			in:     Interval(35 * 24 * 60 * 60 * 1000),
			assert: "35d",
		},
		{
			name:   "5month",
			in:     Interval(5 * 30 * 24 * 60 * 60 * 1000),
			assert: "5M",
		},
		{
			name:   "15month",
			in:     Interval(15 * 30 * 24 * 60 * 60 * 1000),
			assert: "15M",
		},
		{
			name:   "5year",
			in:     Interval(5 * 365 * 24 * 60 * 60 * 1000),
			assert: "5y",
		},
		{
			name:   "455day",
			in:     Interval(455 * 24 * 60 * 60 * 1000),
			assert: "455d",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			strVal := tt.in.String()
			assert.Equal(t, tt.assert, strVal)
		})
	}
}

func Test_IntervalType_String(t *testing.T) {
	assert.Equal(t, "day", Day.String())
}

func Test_Interval_ValueOf(t *testing.T) {
	var i Interval

	assert.NotNil(t, i.ValueOf(" "))

	assert.NotNil(t, i.ValueOf("10t"))

	assert.NotNil(t, i.ValueOf("as"))

	assert.Nil(t, i.ValueOf(" 10 s"))
	assert.Equal(t, 10*timeutil.OneSecond, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 S"))
	assert.Equal(t, 10*timeutil.OneSecond, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 m"))
	assert.Equal(t, 10*timeutil.OneMinute, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 h"))
	assert.Equal(t, 10*timeutil.OneHour, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 H"))
	assert.Equal(t, 10*timeutil.OneHour, i.Int64())

	assert.Nil(t, i.ValueOf(" 10d"))
	assert.Equal(t, 10*timeutil.OneDay, i.Int64())

	assert.Nil(t, i.ValueOf(" 10D"))
	assert.Equal(t, 10*timeutil.OneDay, i.Int64())

	assert.Nil(t, i.ValueOf(" 10M"))
	assert.Equal(t, 10*timeutil.OneMonth, i.Int64())

	assert.Nil(t, i.ValueOf(" 10y"))
	assert.Equal(t, 10*timeutil.OneYear, i.Int64())

	assert.Nil(t, i.ValueOf(" 10Y"))
	assert.Equal(t, 10*timeutil.OneYear, i.Int64())
}

func Test_IntervalCalculator(t *testing.T) {
	var i Interval

	_ = i.ValueOf("30m")
	assert.NotNil(t, i.Calculator())

	_ = i.ValueOf("1m")
	assert.NotNil(t, i.Calculator())

	_ = i.ValueOf("10d")
	assert.NotNil(t, i.Calculator())
}

func Test_CalcQueryInterval(t *testing.T) {
	now := timeutil.Now()
	cases := []struct {
		name           string
		timeRange      TimeRange
		queryInterval  Interval
		targetInterval Interval
	}{
		{
			name:           "use input interval",
			timeRange:      TimeRange{},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(timeutil.OneSecond),
		},
		{
			name:           "<3hour",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + 2*timeutil.OneHour},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(10 * timeutil.OneSecond),
		},
		{
			name:           "<6hour",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + 4*timeutil.OneHour},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(30 * timeutil.OneSecond),
		},
		{
			name:           "<12hour",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + 11*timeutil.OneHour},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(timeutil.OneMinute),
		},
		{
			name:           "<1day",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + 23*timeutil.OneHour},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(2 * timeutil.OneMinute),
		},
		{
			name:           "<2day",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + 47*timeutil.OneHour},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(5 * timeutil.OneMinute),
		},
		{
			name:           "<7day",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + 7*timeutil.OneDay - 1},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(10 * timeutil.OneMinute),
		},
		{
			name:           "<1month",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + timeutil.OneMonth - 1},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(timeutil.OneHour),
		},
		{
			name:           "<2month",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + 2*timeutil.OneMonth - 1},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(4 * timeutil.OneHour),
		},
		{
			name:           "<3month",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + 3*timeutil.OneMonth - 1},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(12 * timeutil.OneHour),
		},
		{
			name:           ">3month",
			timeRange:      TimeRange{Start: timeutil.Now(), End: now + timeutil.OneYear},
			queryInterval:  Interval(timeutil.OneSecond),
			targetInterval: Interval(timeutil.OneDay),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			interval := CalcQueryInterval(tt.timeRange, tt.queryInterval)
			assert.Equal(t, tt.targetInterval, interval)
		})
	}
}

func TestInterval_CalcQuerySlotRange(t *testing.T) {
	t1, _ := timeutil.ParseTimestamp("20190101 00:00:00", "20060102 15:04:05")
	t2, _ := timeutil.ParseTimestamp("20190101 03:10:00", "20060102 15:04:05")
	slotRange := Interval(timeutil.OneMinute).CalcSlotRange(t1, TimeRange{
		Start: t1,
		End:   t2,
	})
	assert.Equal(t, SlotRange{
		Start: 0,
		End:   59,
	}, slotRange)
}
