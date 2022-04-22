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

package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetStringValue(t *testing.T) {
	assert.Equal(t, "sum", GetStringValue("sum"))
	assert.Equal(t, "sum", GetStringValue("'sum'"))
	assert.Equal(t, "'sum", GetStringValue("'sum"))
	assert.Equal(t, "sum", GetStringValue("\"sum\""))
	assert.Equal(t, "", GetStringValue(""))
}

func Test_ByteSlice2String(t *testing.T) {
	s := []byte("abc")
	assert.Equal(t, "abc", ByteSlice2String(s))
}

func Test_String2ByteSlice(t *testing.T) {
	s := "abc"
	assert.Equal(t, []byte("abc"), String2ByteSlice(s))
}

func Test_DeDupString(t *testing.T) {
	assert.Len(t, DeDupStringSlice(nil), 0)
	assert.Len(t, DeDupStringSlice([]string{"a", "v"}), 2)
	assert.Len(t, DeDupStringSlice([]string{"a", "a", "b", "v"}), 3)
}

func Test_RandomString(t *testing.T) {
	t.Log(RandStringBytes(20))
}
