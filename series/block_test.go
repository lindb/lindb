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

package series

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_Append(t *testing.T) {
	block := NewBlock(1, 1)
	assert.False(t, block.Append(0, 10.0))
	assert.False(t, block.Append(1, 10.0))
	assert.True(t, block.Append(2, 10.0))

	block.Clear()
}
