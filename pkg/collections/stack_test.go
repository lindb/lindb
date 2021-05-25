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

package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack()
	assert.True(t, s.Empty())

	s.Push(1)
	assert.False(t, s.Empty())
	assert.Equal(t, 1, s.Size())

	assert.Equal(t, 1, s.Peek().(int))

	assert.Equal(t, 1, s.Pop().(int))
	assert.True(t, s.Empty())

	s.Push(1)
	s.Push(2)

	assert.Equal(t, 2, s.Size())
	assert.Equal(t, 2, s.Pop().(int))
	assert.Equal(t, 1, s.Size())

	s.Pop()

	assert.Equal(t, 0, s.Size())
	assert.True(t, s.Empty())
	assert.Nil(t, s.Peek())
	assert.Nil(t, s.Pop())
}
