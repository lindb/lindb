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
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/strutil"
)

func BenchmarkIndexKVStore(b *testing.B) {
	name := "./index_kv" + b.TempDir()
	defer func() {
		_ = os.RemoveAll(name)
	}()
	store, _ := kv.GetStoreManager().CreateStore(name, kv.StoreOption{Levels: 2})
	family, _ := store.CreateFamily("index", familyOption)
	indexStore := NewIndexKVStore(family, 1000, 10*time.Minute)

	seq := uint32(0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	buckets := imap.NewIntMap[map[string]uint32]()
	var scratch [8]byte
	for bucketID := uint32(0); bucketID < 100; bucketID++ {
		kvs := make(map[string]uint32)
		buckets.Put(bucketID, kvs)
		for i := 0; i < 1000; i++ {
			binary.LittleEndian.PutUint64(scratch[:], r.Uint64())
			id, _ := indexStore.GetOrCreateValue(bucketID, scratch[:], func() uint32 {
				seq++
				return seq
			})
			kvs[string(scratch[:])] = id
		}
	}

	// flush
	_ = indexStore.Flush()

	for i := 0; i < b.N; i++ {
		_ = buckets.WalkEntry(func(key uint32, value map[string]uint32) error {
			for k := range value {
				_, _ = indexStore.GetOrCreateValue(key, strutil.String2ByteSlice(k), func() uint32 {
					panic("err")
				})
			}
			return nil
		})
	}
}
