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

package http

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
)

var tokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGF" +
	"zc3dvcmQiOiJhZG1pbjEyMyJ9.YbNGN0V-U5Y3xOIGNXcgbQkK2VV30UDDEZV19FN62hk"

func Test_CreateToken(t *testing.T) {
	user := config.User{UserName: "admin", Password: "admin123"}
	token, err := CreateToken(user)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, tokenStr, token)
}
