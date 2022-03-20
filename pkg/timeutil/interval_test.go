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

	"github.com/lindb/lindb/pkg/encoding"

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
			"unmarshal_err",
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
			name:   "5min",
			in:     Interval(5 * 60 * 1000),
			assert: "5m",
		},
		{
			name:   "5hour",
			in:     Interval(5 * 60 * 60 * 1000),
			assert: "5h",
		},
		{
			name:   "5day",
			in:     Interval(5 * 24 * 60 * 60 * 1000),
			assert: "5d",
		},
		{
			name:   "5month",
			in:     Interval(5 * 31 * 24 * 60 * 60 * 1000),
			assert: "5M",
		},
		{
			name:   "5year",
			in:     Interval(5 * 365 * 24 * 60 * 60 * 1000),
			assert: "5y",
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
	assert.Equal(t, 10*OneSecond, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 S"))
	assert.Equal(t, 10*OneSecond, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 m"))
	assert.Equal(t, 10*OneMinute, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 h"))
	assert.Equal(t, 10*OneHour, i.Int64())

	assert.Nil(t, i.ValueOf(" 10 H"))
	assert.Equal(t, 10*OneHour, i.Int64())

	assert.Nil(t, i.ValueOf(" 10d"))
	assert.Equal(t, 10*OneDay, i.Int64())

	assert.Nil(t, i.ValueOf(" 10D"))
	assert.Equal(t, 10*OneDay, i.Int64())

	assert.Nil(t, i.ValueOf(" 10M"))
	assert.Equal(t, 10*OneMonth, i.Int64())

	assert.Nil(t, i.ValueOf(" 10y"))
	assert.Equal(t, 10*OneYear, i.Int64())

	assert.Nil(t, i.ValueOf(" 10Y"))
	assert.Equal(t, 10*OneYear, i.Int64())
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
