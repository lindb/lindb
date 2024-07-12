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

package imap

import (
	"fmt"
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

// hack test
func _assertIntMapData(t *testing.T, keys []uint32, m *IntMap[string]) {
	for _, key := range keys {
		found, highIdx := m.keys.ContainsAndRankForHigh(key)
		assert.True(t, found)
		lowIdx := m.keys.RankForLow(key, highIdx-1)
		assert.True(t, found)
		assert.NotNil(t, m.values[highIdx-1][lowIdx-1])
	}
}

func TestIntMap_Put(t *testing.T) {
	m := NewIntMap[string]()
	m.Put(1, "value1")
	m.Put(8, "value8")
	m.Put(3, "value3")
	m.Put(5, "value5")
	m.Put(6, "value6")
	m.Put(7, "value7")
	m.Put(4, "value4")
	m.Put(2, "value2")
	// test insert new high
	m.Put(2000000, "value2000000")
	m.Put(2000001, "value2000000")
	// test insert new high
	m.PutIfNotExist(200000, "value200000")

	_assertIntMapData(t, []uint32{1, 2, 3, 4, 5, 6, 7, 8, 200000, 2000000, 2000001}, m)
	assert.Equal(t, 11, m.Size())
	assert.Len(t, m.Values(), 3)

	err := m.WalkEntry(func(key uint32, value string) error {
		return fmt.Errorf("err")
	})
	assert.Error(t, err)

	keys := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 200000, 2000000, 2000001}
	idx := 0
	err = m.WalkEntry(func(key uint32, value string) error {
		assert.Equal(t, keys[idx], key)
		idx++
		return nil
	})
	assert.NoError(t, err)
	assert.False(t, m.IsEmpty())
}

func TestIntMap_Get(t *testing.T) {
	m := NewIntMap[string]()
	_, ok := m.Get(uint32(10))
	assert.False(t, ok)
	m.Put(1, "value1")
	m.Put(8, "value8")
	_, ok = m.Get(1)
	assert.True(t, ok)
	_, ok = m.Get(2)
	assert.False(t, ok)
	_, ok = m.Get(0)
	assert.False(t, ok)
	_, ok = m.Get(9)
	assert.False(t, ok)
	_, ok = m.Get(999999)
	assert.False(t, ok)
}

func TestIntMap_Keys(t *testing.T) {
	m := NewIntMap[string]()
	m.Put(1, "value1")
	m.Put(8, "value8")
	assert.Equal(t, roaring.BitmapOf(1, 8), m.Keys())
}

func TestIntMap_tryOptimize(t *testing.T) {
	m := NewIntMap[string]()
	for i := 0; i < 100; i++ {
		m.Put(uint32(i), "value")
	}
}
