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

	"github.com/stretchr/testify/assert"
)

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
