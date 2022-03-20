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
	"encoding/json"
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

func Test_Duration_JSON(t *testing.T) {
	type Example struct {
		Cost Duration `json:"cost"`
		A    int      `json:"a"`
	}
	example := Example{Cost: Duration(time.Nanosecond * 1102), A: 23}
	data, err := json.Marshal(example)
	assert.NoError(t, err)
	_, err = example.Cost.MarshalJSON()
	assert.NoError(t, err)
	_, err = example.Cost.MarshalText()
	assert.NoError(t, err)

	var newExample Example
	assert.NoError(t, json.Unmarshal(data, &newExample))
	assert.Equal(t, Duration(time.Nanosecond*1102), newExample.Cost)

	txt := "{\"cost\": 322}"
	assert.NoError(t, json.Unmarshal([]byte(txt), &newExample))
	assert.Equal(t, Duration(time.Nanosecond*322), newExample.Cost)

	assert.Error(t, json.Unmarshal([]byte(`{"cost": "xxxx"}`), &newExample))
	assert.Error(t, json.Unmarshal([]byte(`{"cost": "xxxx"}`), &newExample))

	assert.Error(t, json.Unmarshal([]byte(`{"cost": null}`), &newExample))
	assert.Error(t, json.Unmarshal([]byte{1, 0}, &newExample))

	assert.Nil(t, json.Unmarshal([]byte(`{"cost": "22.265928ms"}`), &newExample))
}

func Test_Size(t *testing.T) {
	type Example struct {
		Size Size `json:"size"`
	}
	example1 := Example{Size: 10240}
	assert.Equal(t, "10 KiB", example1.Size.String())

	txt, err := example1.Size.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "10 KiB", string(txt))

	txt, err = example1.Size.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"10 KiB"`, string(txt))

	var s2 Example
	assert.Error(t, json.Unmarshal([]byte(``), &s2))
	assert.Error(t, json.Unmarshal([]byte(`{"size": null`), &s2))
	assert.Error(t, json.Unmarshal([]byte(`{"size": true`), &s2))
	assert.Error(t, json.Unmarshal([]byte(`{"size": "10 MiB"`), &s2))
	assert.NoError(t, json.Unmarshal([]byte(`{"size": "10 MiB"}`), &s2))
	assert.Equal(t, Size(0xa00000), s2.Size)
	assert.Error(t, json.Unmarshal([]byte(`{"size": "10 iB"}`), &s2))
	assert.NoError(t, json.Unmarshal([]byte(`{"size": 1000}`), &s2))
	assert.NoError(t, json.Unmarshal([]byte(`{"size": "969 B"}`), &s2))
	assert.NoError(t, json.Unmarshal([]byte(`{"size": 10}`), &s2))
	assert.Error(t, json.Unmarshal([]byte("{\"size\": \"\"\"\"}"), &s2))
}
