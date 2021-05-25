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

package ltoml

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Duration(t *testing.T) {
	assert.Equal(t, Duration(time.Minute).Duration(), time.Minute)

	marshalF := func(duration time.Duration) string {
		txt, _ := Duration(duration).MarshalText()
		return string(txt)
	}
	unmarshalF := func(txt string) time.Duration {
		var d Duration
		_ = d.UnmarshalText([]byte(txt))
		return d.Duration()
	}
	assert.Equal(t, "1m0s", marshalF(time.Minute))
	assert.Equal(t, "10s", marshalF(time.Second*10))

	assert.Equal(t, time.Second, unmarshalF("1s"))
	assert.Equal(t, time.Minute, unmarshalF("1m"))
	assert.Equal(t, time.Hour, unmarshalF("3600s"))

	assert.Zero(t, unmarshalF(""))
	assert.Zero(t, unmarshalF("1fs"))
}
