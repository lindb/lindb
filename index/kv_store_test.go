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

package index

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/stmt"
)

var familyOption = kv.FamilyOption{Merger: string(v1.IndexKVMerger)}

func TestIndexKVStore(t *testing.T) {
	name := "./index_kv"
	defer func() {
		_ = os.RemoveAll(name)
	}()
	indexStore := createIndexKVStore(t, name)

	seq := uint32(0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	buckets := imap.NewIntMap[map[string]uint32]()
	var scratch [8]byte
	for bucketID := uint32(0); bucketID < 100; bucketID++ {
		kvs := make(map[string]uint32)
		buckets.Put(bucketID, kvs)
		for i := 0; i < 1000; i++ {
			binary.LittleEndian.PutUint64(scratch[:], r.Uint64())
			id, isNew, err0 := indexStore.GetOrCreateValue(bucketID, scratch[:], func() uint32 {
				seq++
				return seq
			})
			assert.True(t, isNew)
			assert.NoError(t, err0)
			kvs[string(scratch[:])] = id
		}
	}

	test := func() {
		_ = buckets.WalkEntry(func(key uint32, value map[string]uint32) error {
			for k, v := range value {
				id, isNew, err0 := indexStore.GetOrCreateValue(key, strutil.String2ByteSlice(k), func() uint32 {
					panic("err")
				})
				assert.False(t, isNew)
				assert.NoError(t, err0)
				assert.Equal(t, id, v)

				id, ok, err0 := indexStore.GetValue(key, strutil.String2ByteSlice(k))
				assert.NoError(t, err0)
				assert.True(t, ok)
				assert.Equal(t, id, v)
			}
			return nil
		})
		values, err := indexStore.GetValues(10)
		assert.NoError(t, err)
		assert.Len(t, values, 1000)
		result := make(map[uint32]string)
		assert.NoError(t, indexStore.CollectKVs(100, roaring.New(), result))
		assert.Len(t, result, 0)
		assert.NoError(t, indexStore.CollectKVs(100, roaring.BitmapOf(1), result))
		assert.Len(t, result, 0)
		assert.NoError(t, indexStore.CollectKVs(0, roaring.BitmapOf(1), result))
		assert.Len(t, result, 1)
		rs, err := indexStore.Suggest(10, string([]byte{0}), 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, rs)
	}
	test()

	// flush
	indexStore.PrepareFlush()
	assert.NoError(t, indexStore.Flush())

	test()
}

func TestIndexKVStore_Compact(t *testing.T) {
	name := "./index_compact"
	defer func() {
		_ = os.RemoveAll(name)
	}()
	store, err := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	assert.NotNil(t, store)
	family, err := store.CreateFamily("index", familyOption)

	indexStore := NewIndexKVStore(family, 10, time.Minute)
	assert.NoError(t, err)

	write := func(i int) {
		id, isNew, err0 := indexStore.GetOrCreateValue(1, []byte(fmt.Sprintf("key-%d", i)), func() uint32 {
			return uint32(i)
		})
		assert.True(t, isNew)
		assert.NoError(t, err0)
		assert.Equal(t, uint32(i), id)
		indexStore.PrepareFlush()
		assert.NoError(t, indexStore.Flush())
	}
	for i := 0; i < 5; i++ {
		write(i)
	}
	id, ok, err := indexStore.GetValue(1, []byte("key-1"))
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, uint32(1), id)

	family.Compact()
	time.Sleep(200 * time.Millisecond)

	id, ok, err = indexStore.GetValue(1, []byte("key-1"))
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, uint32(1), id)

	assert.NoError(t, kv.GetStoreManager().CloseStore(name))
}

func TestIndexKVStore_Find(t *testing.T) {
	name := "./index_find"
	defer func() {
		_ = os.RemoveAll(name)
	}()
	indexStore := createIndexKVStore(t, name)
	keysString := []string{
		"a", "ab", "b", "abc", "abcdefgh", "abcdefghijklmnopqrstuvwxyz", "abcdefghijkl", "zzzzzz", "ice",
	}
	for idx, key := range keysString {
		v, isNew, err := indexStore.GetOrCreateValue(100, []byte(key), func() uint32 { return uint32(idx) })
		assert.NoError(t, err)
		assert.True(t, isNew)
		assert.Equal(t, uint32(idx), v)
	}
	cases := []struct {
		expr    stmt.TagFilter
		name    string
		size    int
		wantErr bool
	}{
		{
			name: "find data(equals expr)",
			expr: &stmt.EqualsExpr{
				Value: "ab",
			},
			size: 1,
		},
		{
			name: "no data(equals expr)",
			expr: &stmt.EqualsExpr{
				Value: "ab0",
			},
		},
		{
			name: "find data(in expr)",
			expr: &stmt.InExpr{
				Values: []string{"ab0", "a", "abc", "ab"},
			},
			size: 3,
		},
		{
			name: "find no data(in expr)",
			expr: &stmt.InExpr{
				Values: []string{"ab0"},
			},
		},
		{
			name: "find data(regex expr)",
			expr: &stmt.RegexExpr{
				Regexp: "^abc",
			},
			size: 4,
		},
		{
			name: "find no data(regex expr)",
			expr: &stmt.RegexExpr{
				Regexp: "hh",
			},
		},
		{
			name: "find data(like prefix expr)",
			expr: &stmt.LikeExpr{
				Value: "abc*",
			},
			size: 4,
		},
		{
			name: "find no data(like expr)",
			expr: &stmt.LikeExpr{
				Value: "hh",
			},
		},
		{
			name: "find no data(like empty expr)",
			expr: &stmt.LikeExpr{},
		},
		{
			name: "find data(like suffix expr)",
			expr: &stmt.LikeExpr{
				Value: "*abc",
			},
			size: 1,
		},
		{
			name: "find data(like prefix/suffix expr)",
			expr: &stmt.LikeExpr{
				Value: "*abc*",
			},
			size: 4,
		},
	}
	test := func() {
		for i := range cases {
			tt := cases[i]
			t.Run(tt.name, func(t *testing.T) {
				ids, err := indexStore.FindValuesByExpr(100, tt.expr)
				if (err != nil) != tt.wantErr {
					t.Fatal(tt.name)
				}
				assert.Len(t, ids, tt.size)
			})
		}
	}

	// from mem store
	test()

	indexStore.PrepareFlush()
	assert.NoError(t, indexStore.Flush())

	// from kv store
	test()
}

func TestIndexKVStore_Read_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		regexpCompile = regexp.Compile
	}()
	kvFamily := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	kvFamily.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	snapshot.EXPECT().Close().AnyTimes()
	store := NewIndexKVStore(kvFamily, 10, time.Second)

	t.Run("load error when get values", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		ids, err := store.GetValues(100)
		assert.Error(t, err)
		assert.Empty(t, ids)
	})

	t.Run("load error when do in expr", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		ids, err := store.FindValuesByExpr(10, &stmt.InExpr{Values: []string{"abc"}})
		assert.Error(t, err)
		assert.Empty(t, ids)
	})

	t.Run("load error when do like expr", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		ids, err := store.FindValuesByExpr(10, &stmt.LikeExpr{Value: "abc*"})
		assert.Error(t, err)
		assert.Empty(t, ids)
	})

	t.Run("load error when do regex expr", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		ids, err := store.FindValuesByExpr(10, &stmt.RegexExpr{Regexp: "^abc"})
		assert.Error(t, err)
		assert.Empty(t, ids)
	})

	t.Run("load error when collect", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		err := store.CollectKVs(10, roaring.New(), nil)
		assert.Error(t, err)
	})

	t.Run("compile regex error", func(t *testing.T) {
		regexpCompile = func(expr string) (*regexp.Regexp, error) {
			return nil, fmt.Errorf("err")
		}
		ids, err := store.FindValuesByExpr(10, &stmt.RegexExpr{})
		assert.Error(t, err)
		assert.Empty(t, ids)
	})
	t.Run("suggest error", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		rs, err := store.Suggest(10, "test", 10)
		assert.Error(t, err)
		assert.Empty(t, rs)
	})
}

func TestIndexKVStore_Flush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	kvFamily := kv.NewMockFamily(ctrl)
	kvFamily.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	kvFlusher := kv.NewMockFlusher(ctrl)
	kvFamily.EXPECT().NewFlusher().Return(kvFlusher).AnyTimes()
	kvFlusher.EXPECT().Release().AnyTimes()
	store := NewIndexKVStore(kvFamily, 10, time.Minute)
	store1 := store.(*indexKVStore)
	flusher := v1.NewMockIndexKVFlusher(ctrl)

	cases := []struct {
		prepare func()
		name    string
		wantErr bool
	}{
		{
			name: "no data need flush",
			prepare: func() {
				store1.immutable = nil
			},
			wantErr: false,
		},
		{
			name: "create index kv flusher error",
			prepare: func() {
				newIndexKVFlusher = func(_ int, _ kv.Flusher) (v1.IndexKVFlusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "empty kvs",
			prepare: func() {
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "close flusher error",
			prepare: func() {
				flusher.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "write kvs error",
			prepare: func() {
				store1.immutable = imap.NewIntMap[map[string]uint32]()
				store1.immutable.Put(100, map[string]uint32{"test": 1})
				flusher.EXPECT().PrepareBucket(gomock.Any())
				flusher.EXPECT().WriteKVs(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
	}

	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newIndexKVFlusher = v1.NewIndexKVFlusher
			}()
			newIndexKVFlusher = func(_ int, _ kv.Flusher) (v1.IndexKVFlusher, error) {
				return flusher, nil
			}
			store1.immutable = imap.NewIntMap[map[string]uint32]()
			store1.immutable.Put(100, map[string]uint32{})
			tt.prepare()
			err := store.Flush()
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

func createIndexKVStore(t *testing.T, name string) IndexKVStore {
	store, err := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	assert.NotNil(t, store)
	family, err := store.CreateFamily("index", familyOption)

	indexStore := NewIndexKVStore(family, 10, time.Minute)
	assert.NoError(t, err)
	return indexStore
}
