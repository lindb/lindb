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

package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryOPString(t *testing.T) {
	assert.Equal(t, "and", BinaryOPString(AND))
	assert.Equal(t, "or", BinaryOPString(OR))

	assert.Equal(t, "+", BinaryOPString(ADD))
	assert.Equal(t, "-", BinaryOPString(SUB))
	assert.Equal(t, "*", BinaryOPString(MUL))
	assert.Equal(t, "/", BinaryOPString(DIV))

	assert.Equal(t, "=", BinaryOPString(EQUAL))
	assert.Equal(t, "!=", BinaryOPString(NOTEQUAL))
	assert.Equal(t, ">", BinaryOPString(GREATER))
	assert.Equal(t, ">=", BinaryOPString(GREATEREQUAL))
	assert.Equal(t, "<", BinaryOPString(LESS))
	assert.Equal(t, "<=", BinaryOPString(LESSEQUAL))
	assert.Equal(t, "like", BinaryOPString(LIKE))

	assert.Equal(t, "unknown", BinaryOPString(UNKNOWN))
}
