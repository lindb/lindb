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

package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexSlotSelector_IndexOf(t *testing.T) {
	selector := NewIndexSlotSelector(10, 120, 1)
	assert.Equal(t, 110, selector.PointCount())
	start, end := selector.Range()
	assert.Equal(t, 10, start)
	assert.Equal(t, 120, end)
	idx, completed := selector.IndexOf(5)
	assert.Equal(t, -1, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(10)
	assert.Equal(t, 0, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(15)
	assert.Equal(t, 5, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(105)
	assert.Equal(t, 95, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(121)
	assert.Equal(t, -1, idx)
	assert.True(t, completed)

	selector = NewIndexSlotSelector(10, 130, 3)
	idx, completed = selector.IndexOf(12)
	assert.Equal(t, 0, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(13)
	assert.Equal(t, 1, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(130)
	assert.Equal(t, 40, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(131)
	assert.Equal(t, -1, idx)
	assert.True(t, completed)
}
