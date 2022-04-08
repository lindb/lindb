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

package option

import (
	"sort"
	"testing"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseOption_Validate(t *testing.T) {
	cases := []struct {
		name    string
		in      DatabaseOption
		wantErr bool
	}{
		{
			"empty intervals",
			DatabaseOption{},
			true,
		},
		{
			"ahead invalid",
			DatabaseOption{Intervals: Intervals{{}}, Ahead: "aa"},
			true,
		},
		{
			"behind invalid",
			DatabaseOption{Intervals: Intervals{{}}, Behind: "aa"},
			true,
		},
		{
			"interval cannot be negative",
			DatabaseOption{Intervals: Intervals{{}}, Behind: "0h"},
			true,
		},
		{
			"validation pass",
			DatabaseOption{Intervals: Intervals{{}}, Behind: "1h", Ahead: "1h"},
			false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if err := tt.in.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseOption_GetAcceptWritableRange(t *testing.T) {
	cases := []struct {
		name    string
		in      DatabaseOption
		prepare func(in *DatabaseOption)
		assert  func(ahead, behind int64)
	}{
		{
			"default option",
			DatabaseOption{},
			func(in *DatabaseOption) {
				in.Default()
			},
			func(ahead, behind int64) {
				assert.Equal(t, constants.MetricMaxAheadDuration, ahead)
				assert.Equal(t, constants.MetricMaxBehindDuration, behind)
			},
		},
		{

			"get accept writable range",
			DatabaseOption{Ahead: "10s", Behind: "20s"},
			func(in *DatabaseOption) {
				in.Default()
			},
			func(ahead, behind int64) {
				assert.Equal(t, int64(10000), ahead)
				assert.Equal(t, int64(20000), behind)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(&tt.in)
			}
			ahead, behind := tt.in.GetAcceptWritableRange()
			if tt.assert != nil {
				tt.assert(ahead, behind)
			}
		})
	}
}

func TestInterval_String(t *testing.T) {
	assert.Equal(t, "10s->1M",
		Interval{
			Interval:  timeutil.Interval(10 * timeutil.OneSecond),
			Retention: timeutil.Interval(timeutil.OneMonth),
		}.String(),
	)
}

func TestIntervals_Sort(t *testing.T) {
	intervals := Intervals{
		{timeutil.Interval(timeutil.OneMinute), timeutil.Interval(timeutil.OneMonth)},
		{timeutil.Interval(timeutil.OneHour), timeutil.Interval(timeutil.OneMonth)},
		{timeutil.Interval(timeutil.OneSecond), timeutil.Interval(timeutil.OneMonth)},
	}
	sort.Sort(intervals)
	assert.Equal(t, Intervals{
		{timeutil.Interval(timeutil.OneSecond), timeutil.Interval(timeutil.OneMonth)},
		{timeutil.Interval(timeutil.OneMinute), timeutil.Interval(timeutil.OneMonth)},
		{timeutil.Interval(timeutil.OneHour), timeutil.Interval(timeutil.OneMonth)},
	}, intervals)

	assert.Equal(t, "[1s->1M,1m->1M,1h->1M]", intervals.String())
}

func TestDatabaseOption_FindMatchSmallestInterval(t *testing.T) {
	opt := DatabaseOption{Intervals: Intervals{
		{timeutil.Interval(timeutil.OneSecond), timeutil.Interval(timeutil.OneMonth)},
		{timeutil.Interval(timeutil.OneMinute), timeutil.Interval(timeutil.OneMonth)},
		{timeutil.Interval(timeutil.OneHour), timeutil.Interval(timeutil.OneMonth)},
	}}
	interval := opt.FindMatchSmallestInterval(timeutil.Interval(timeutil.OneMinute * 3))
	assert.Equal(t, timeutil.Interval(timeutil.OneMinute), interval)
}
