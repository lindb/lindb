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

package function

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
)

func TestFuncCall_Unknown(t *testing.T) {
	result := FuncCall(Unknown, collections.NewFloatArray(10))
	assert.Nil(t, result)
}

func TestFuncCall_Sum(t *testing.T) {
	result := FuncCall(Sum, nil)
	assert.Nil(t, result)
	result = FuncCall(Sum)
	assert.Nil(t, result)

	array1 := collections.NewFloatArray(10)
	array2 := collections.NewFloatArray(20)
	result = FuncCall(Sum, array1, array2)
	assert.Equal(t, array1, result)
}
