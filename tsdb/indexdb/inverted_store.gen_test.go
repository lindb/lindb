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

// Code generated by tmpl; DO NOT EDIT.
// https://github.com/benbjohnson/tmpl
//
// Source: int_map_test.tmpl

package indexdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"
)

// hack test
func _assertInvertedStoreData(t *testing.T, keys []uint32, m *InvertedStore) {
	for _, key := range keys {
		found, highIdx := m.keys.ContainsAndRankForHigh(key)
		assert.True(t, found)
		lowIdx := m.keys.RankForLow(key, highIdx-1)
		assert.True(t, found)
		assert.NotNil(t, m.values[highIdx-1][lowIdx-1])
	}
}

func TestInvertedStore_Put(t *testing.T) {
	m := NewInvertedStore()
	m.Put(1, roaring.New())
	m.Put(8, roaring.New())
	m.Put(3, roaring.New())
	m.Put(5, roaring.New())
	m.Put(6, roaring.New())
	m.Put(7, roaring.New())
	m.Put(4, roaring.New())
	m.Put(2, roaring.New())
	// test insert new high
	m.Put(2000000, roaring.New())
	m.Put(2000001, roaring.New())
	// test insert new high
	m.Put(200000, roaring.New())

	_assertInvertedStoreData(t, []uint32{1, 2, 3, 4, 5, 6, 7, 8, 200000, 2000000, 2000001}, m)
	assert.Equal(t, 11, m.Size())
	assert.Len(t, m.Values(), 3)

	err := m.WalkEntry(func(key uint32, value *roaring.Bitmap) error {
		return fmt.Errorf("err")
	})
	assert.Error(t, err)

	keys := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 200000, 2000000, 2000001}
	idx := 0
	err = m.WalkEntry(func(key uint32, value *roaring.Bitmap) error {
		assert.Equal(t, keys[idx], key)
		idx++
		return nil
	})
	assert.NoError(t, err)
}

func TestInvertedStore_Get(t *testing.T) {
	m := NewInvertedStore()
	_, ok := m.Get(uint32(10))
	assert.False(t, ok)
	m.Put(1, roaring.New())
	m.Put(8, roaring.New())
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

func TestInvertedStore_Keys(t *testing.T) {
	m := NewInvertedStore()
	m.Put(1, roaring.New())
	m.Put(8, roaring.New())
	assert.Equal(t, roaring.BitmapOf(1, 8), m.Keys())
}

func TestInvertedStore_tryOptimize(t *testing.T) {
	m := NewInvertedStore()
	for i := 0; i < 100; i++ {
		m.Put(uint32(i), roaring.New())
	}
}
