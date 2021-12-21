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

package unique

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDStore_New_err(t *testing.T) {
	p := t.TempDir()
	_, _ = NewIDStore(p)

	// cannot create duplicate
	store, err := NewIDStore(p)
	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestIDStore_CRUD(t *testing.T) {
	p := t.TempDir()
	store, err := NewIDStore(p)
	assert.NoError(t, err)
	assert.NotNil(t, store)

	k := []byte("k")
	v := []byte("v")
	err = store.Put(k, v)
	assert.NoError(t, err)

	cases := []struct {
		name    string
		key     []byte
		value   []byte
		exist   bool
		wantErr bool
		prepare func()
	}{
		{
			name: "key not exist",
			key:  []byte("key_not_exit"),
		},
		{
			name:  "get value by key",
			key:   []byte("key"),
			value: []byte("value"),
			exist: true,
			prepare: func() {
				err := store.Put([]byte("key"), []byte("value"))
				assert.NoError(t, err)
			},
		},
		{
			name: "delete key",
			key:  []byte("del_key"),
			prepare: func() {
				key := []byte("del_key")
				err := store.Put(key, []byte("del_val"))
				assert.NoError(t, err)
				_, exist, _ := store.Get(key)
				assert.True(t, exist)
				err = store.Delete(key)
				assert.NoError(t, err)
			},
		},
		{
			name:  "merge key",
			key:   []byte("merge_key"),
			value: []byte("value0value1"),
			exist: true,
			prepare: func() {
				key := []byte("merge_key")
				for i := 0; i < 2; i++ {
					err := store.Merge(key, []byte(fmt.Sprintf("value%d", i)))
					assert.NoError(t, err)
				}
			},
		},
		{
			name:  "over write merge key",
			key:   []byte("over_write_merge"),
			value: []byte("final"),
			exist: true,
			prepare: func() {
				key := []byte("over_write_merge")
				for i := 0; i < 2; i++ {
					err := store.Merge(key, []byte(fmt.Sprintf("value%d", i)))
					assert.NoError(t, err)
				}
				v, _, _ := store.Get(key)
				assert.Equal(t, []byte("value0value1"), v)

				// overwrite
				err = store.Put(key, []byte("final"))
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			v, exist, err := store.Get(tt.key)

			assert.Equal(t, tt.value, v)
			assert.Equal(t, tt.exist, exist)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	err = store.Close()
	assert.NoError(t, err)
	// test re-open
	store, err = NewIDStore(p)
	assert.NoError(t, err)
	assert.NotNil(t, store)

	val, exist, err := store.Get(k)
	assert.Equal(t, v, val)
	assert.True(t, exist)
	assert.NoError(t, err)
}

func TestIdStore_IterIDKeys(t *testing.T) {
	p := t.TempDir()
	store, err := NewIDStore(p)
	assert.NoError(t, err)
	mock(t, store)

	cases := []struct {
		name   string
		prefix string
		limit  int
		length int
	}{
		{
			name:   "test limit",
			prefix: "ns",
			limit:  1,
			length: 1,
		},
		{
			name:   "test not limit",
			prefix: "ns",
			limit:  100,
			length: 10,
		},
		{
			name:   "test prefix",
			prefix: "ns-9",
			limit:  1,
			length: 1,
		},
		{
			name:   "not found",
			prefix: "ns-99",
			limit:  0,
			length: 0,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			keys, err := store.IterKeys([]byte(tt.prefix), tt.limit)
			assert.NoError(t, err)
			assert.Len(t, keys, tt.length)
		})
	}
}

func mock(t *testing.T, store IDStore) {
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("ns-%d", i)
		err := store.Put([]byte(key), []byte(key))
		assert.NoError(t, err)
	}
	err := store.Put([]byte("word"), []byte("word"))
	assert.NoError(t, err)
}
