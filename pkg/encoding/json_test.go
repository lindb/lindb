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

package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUser struct {
	Name string
}

func TestJSONCodec(t *testing.T) {
	user := mockUser{Name: "LinDB"}
	data := JSONMarshal(&user)
	user1 := mockUser{}
	err := JSONUnmarshal(data, &user1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user, user)
	err = JSONUnmarshal([]byte{1, 1, 1}, &user1)
	assert.NotNil(t, err)
}

func Test_JSONMarshal(t *testing.T) {
	assert.Len(t, JSONMarshal(make(chan struct{}, 1)), 0)
	assert.True(t, len(JSONMarshal(nil)) > 0)
}
