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
	"bytes"
	"math"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/index/model"
	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/tree"
)

// for testing
var (
	regexpCompile     = regexp.Compile
	newIndexKVFlusher = v1.NewIndexKVFlusher
)

// indexKVStore implements IndexKVStore interface.
type indexKVStore struct {
	family   kv.Family
	snapshot version.Snapshot

	// memory store
	mutable   *imap.IntMap[map[string]uint32]
	immutable *imap.IntMap[map[string]uint32]
	// cache
	bucketCache *expirable.LRU[uint32, *model.TrieBucket]

	lock sync.RWMutex
}

// NewIndexKVStore creates an IndexKVStore instance.
func NewIndexKVStore(family kv.Family, cacheSize int, cacheTTL time.Duration) IndexKVStore {
	return &indexKVStore{
		family:   family,
		snapshot: family.GetSnapshot(),
		mutable:  imap.NewIntMap[map[string]uint32](),
		bucketCache: expirable.NewLRU(cacheSize, func(_ uint32, value *model.TrieBucket) {
			value.Release()
		}, cacheTTL),
	}
}

// GetOrCreateValue returns unique id for key, if key not exist, creates a new unique id.
func (s *indexKVStore) GetOrCreateValue(bucketID uint32,
	key []byte,
	createFn func() (uint32, error),
) (id uint32, isNew bool, err error) {
	id, _, isNew, err = s.getOrCreateValue(bucketID, key, createFn)
	return
}

// GetValue returns value based on bucket and key.
func (s *indexKVStore) GetValue(bucketID uint32, key []byte) (id uint32, ok bool, err error) {
	id, ok, _, err = s.getOrCreateValue(bucketID, key, nil)
	return
}

// GetValues returns all values for bucket.
func (s *indexKVStore) GetValues(bucketID uint32) (ids []uint32, err error) {
	snapshot := s.getSnapshot()

	reader := v1.NewIndexKVReader(snapshot)
	bucket, err := reader.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}
	if bucket != nil {
		defer bucket.Release()
		ids = bucket.GetValues()
	}

	// find from memory
	s.lock.RLock()
	defer s.lock.RUnlock()

	ids = s.getValuesFromMem(s.mutable, bucketID, ids)
	ids = s.getValuesFromMem(s.immutable, bucketID, ids)
	return ids, nil
}

// FindValuesByExpr returns values based on filter expr.
func (s *indexKVStore) FindValuesByExpr(bucketID uint32, expr tree.Expr) (ids []uint32, err error) {
	switch expression := expr.(type) {
	case *tree.EqualsExpr:
		key := strutil.String2ByteSlice(expression.Value)
		return s.findValue(bucketID, key, ids)
	case *tree.InExpr:
		values := expression.Values
		for _, value := range values {
			key := strutil.String2ByteSlice(value)
			ids, err = s.findValue(bucketID, key, ids)
			if err != nil {
				return nil, err
			}
		}
	case *tree.LikeExpr:
		return s.FindValuesByLike(bucketID, expression.Value, ids)
	case *tree.RegexExpr:
		rp, err := regexpCompile(expression.Regexp)
		if err != nil {
			return nil, err
		}
		return s.FindValuesByRegexp(bucketID, rp, ids)
	}
	return ids, nil
}

// CollectKVs collects all keys based on bucket and values.
func (s *indexKVStore) CollectKVs(bucketID uint32, values *roaring.Bitmap, result map[uint32]string) error {
	collect := func(mem *imap.IntMap[map[string]uint32]) {
		if mem == nil {
			return
		}
		if values.IsEmpty() {
			return
		}
		kvs, ok := mem.Get(bucketID)
		if !ok {
			return
		}
		for k, v := range kvs {
			if values.IsEmpty() {
				return
			}
			if values.Contains(v) {
				result[v] = k
				values.Remove(v)
			}
		}
	}

	snapshot := s.getSnapshot()

	reader := v1.NewIndexKVReader(snapshot)
	bucket, err := reader.GetBucket(bucketID)
	if err != nil {
		return err
	}

	s.lock.RLock()
	collect(s.mutable)
	collect(s.immutable)
	s.lock.RUnlock()

	if bucket != nil {
		defer bucket.Release()
		bucket.CollectKVs(values, result)
	}

	return nil
}

// Suggest suggests the kv pairs by prefix.
func (s *indexKVStore) Suggest(bucketID uint32, prefix string, limit int) ([]string, error) {
	sortResult := func(rs []string) []string {
		sort.Slice(rs, func(i, j int) bool {
			return rs[i] < rs[j]
		})
		end := len(rs)
		if end >= limit {
			end = limit
		}
		return rs[:end]
	}
	suggest := func(mem *imap.IntMap[map[string]uint32]) []string {
		if mem == nil {
			return nil
		}
		kvs, ok := mem.Get(bucketID)
		if !ok {
			return nil
		}
		var result []string
		for k := range kvs {
			if strings.HasPrefix(k, prefix) {
				result = append(result, k)
			}
		}
		return sortResult(result)
	}

	snapshot := s.getSnapshot()

	reader := v1.NewIndexKVReader(snapshot)
	bucket, err := reader.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}

	var result []string
	s.lock.RLock()
	result = append(result, suggest(s.mutable)...)
	result = append(result, suggest(s.immutable)...)
	s.lock.RUnlock()

	if bucket != nil {
		defer bucket.Release()
		result = append(result, bucket.Suggest(prefix, limit)...)
	}

	return sortResult(result), nil
}

func (s *indexKVStore) PrepareFlush() {
	// swap mutable/immutable store
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.immutable == nil {
		s.immutable = s.mutable
		s.mutable = imap.NewIntMap[map[string]uint32]()
	}
}

// needFlusf returns if memory data need flush.
func (s *indexKVStore) needFlush() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.immutable != nil && !s.immutable.IsEmpty()
}

func (s *indexKVStore) Flush() (err error) {
	if !s.needFlush() {
		return nil
	}
	kvFlusher := s.family.NewFlusher()
	defer kvFlusher.Release()
	// flush buckets
	flusher, err := newIndexKVFlusher(math.MaxInt16, kvFlusher)
	if err != nil {
		return err
	}
	err = s.immutable.WalkEntry(func(bucket uint32, kvs map[string]uint32) error {
		if len(kvs) == 0 {
			return nil
		}
		flusher.PrepareBucket(bucket)
		var keys [][]byte
		var ids []uint32
		for k, v := range kvs {
			keys = append(keys, strutil.String2ByteSlice(k))
			ids = append(ids, v)
		}
		if err0 := flusher.WriteKVs(keys, ids); err0 != nil {
			return err0
		}
		return flusher.CommitBucket()
	})
	if err != nil {
		return err
	}
	err = flusher.Close()
	if err != nil {
		return err
	}
	// clear immutable store, after flush successfully
	s.lock.Lock()
	defer s.lock.Unlock()

	// close old snapshot
	snapshot := s.snapshot
	snapshot.Close()

	s.snapshot = s.family.GetSnapshot()
	s.immutable = nil
	// purge bucket cache, because new kv write
	s.bucketCache.Purge()
	return nil
}

// getSnapshot returns family snapshot.
func (s *indexKVStore) getSnapshot() version.Snapshot {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.snapshot
}

func (s *indexKVStore) getOrCreateValue(bucketID uint32, key []byte,
	createFn func() (uint32, error),
) (id uint32, ok, isNew bool, err error) {
	// get from memory store
	id, ok = s.GetValueFromMem(bucketID, key)
	if ok {
		return id, true, false, nil
	}

	bucket, ok := s.bucketCache.Get(bucketID)
	if !ok {
		// get from kv store(persist)
		snapshot := s.getSnapshot()
		reader := v1.NewIndexKVReader(snapshot)
		bucket, err = reader.GetBucket(bucketID)
		if err != nil {
			return 0, false, false, err
		}
		if bucket != nil {
			s.bucketCache.Add(bucketID, bucket)
		}
	}
	if bucket != nil {
		// check bucket not nil, maybe not exist under kv store
		id, ok = bucket.GetValue(key)
		if ok {
			return id, true, false, nil
		}
	}

	// create new value
	if createFn == nil {
		return 0, false, false, nil
	}
	id, err = s.createValue(bucketID, key, createFn)
	if err != nil {
		return 0, false, false, err
	}
	return id, true, true, nil
}

// createValue creates new value.
func (s *indexKVStore) createValue(bucketID uint32, key []byte, createFn func() (uint32, error)) (uint32, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	kvs, ok := s.mutable.Get(bucketID)
	if !ok {
		kvs = make(map[string]uint32)
		s.mutable.Put(bucketID, kvs)
	}
	// generate and store value
	id, err := createFn()
	if err != nil {
		return 0, err
	}
	kvs[string(key)] = id
	return id, nil
}

// GetValueFromMem returns value from mem store.
func (s *indexKVStore) GetValueFromMem(bucketID uint32, key []byte) (uint32, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	id, ok := s.getValueFromMem(s.mutable, bucketID, key)
	if ok {
		return id, ok
	}
	return s.getValueFromMem(s.immutable, bucketID, key)
}

// FindValuesByRegexp returns values by regexp expr.
func (s *indexKVStore) FindValuesByRegexp(bucketID uint32, rp *regexp.Regexp, ids []uint32) ([]uint32, error) {
	snapshot := s.getSnapshot()

	reader := v1.NewIndexKVReader(snapshot)
	bucket, err := reader.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}
	if bucket != nil {
		defer bucket.Release()
		ids = bucket.FindValuesByRegexp(rp, ids)
	}
	// find from memory
	s.lock.RLock()
	defer s.lock.RUnlock()

	ids = s.findValuesByRegexp(s.mutable, bucketID, rp, ids)
	ids = s.findValuesByRegexp(s.immutable, bucketID, rp, ids)
	return ids, nil
}

// FindValuesByLike returns values by like expr.
func (s *indexKVStore) FindValuesByLike(bucketID uint32, like string, ids []uint32) ([]uint32, error) {
	hashPrefix := strings.HasPrefix(like, "*")
	hasSuffix := strings.HasSuffix(like, "*")
	likeSlice := strutil.String2ByteSlice(like)
	switch {
	case like == "":
		return nil, nil
	// only ends with *
	case !hashPrefix && hasSuffix:
		prefix := likeSlice[:len(likeSlice)-1]
		return s.findValuesByLike(bucketID, prefix, prefix, bytes.HasPrefix, ids)
	// only starts with *
	case hashPrefix && !hasSuffix:
		suffix := likeSlice[1:]
		return s.findValuesByLike(bucketID, nil, suffix, bytes.HasSuffix, ids)
	// starts with and ends with *
	case hashPrefix && hasSuffix:
		middle := likeSlice[1 : len(likeSlice)-1]
		return s.findValuesByLike(bucketID, nil, middle, bytes.Contains, ids)
	default:
		return s.findValue(bucketID, likeSlice, ids)
	}
}

// findValuesByRegexp returns values by regexp from mem store.
func (s *indexKVStore) findValuesByRegexp(mem *imap.IntMap[map[string]uint32], bucketID uint32, rp *regexp.Regexp, ids []uint32) []uint32 {
	if mem == nil {
		return ids
	}
	kvs, ok := mem.Get(bucketID)
	if !ok {
		return ids
	}
	for k, v := range kvs {
		if rp.Match(strutil.String2ByteSlice(k)) {
			ids = append(ids, v)
		}
	}
	return ids
}

func (s *indexKVStore) findValuesByLike(bucketID uint32,
	prefix, subKey []byte,
	check func(a, b []byte) bool, ids []uint32,
) ([]uint32, error) {
	snapshot := s.getSnapshot()
	reader := v1.NewIndexKVReader(snapshot)
	bucket, err := reader.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}
	if bucket != nil {
		defer bucket.Release()
		ids = bucket.FindValuesByLike(prefix, subKey, check, ids)
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	ids = s.findValuesByLikeFormMem(s.mutable, bucketID, subKey, check, ids)
	ids = s.findValuesByLikeFormMem(s.immutable, bucketID, subKey, check, ids)
	return ids, nil
}

func (s *indexKVStore) findValuesByLikeFormMem(
	mem *imap.IntMap[map[string]uint32],
	bucketID uint32, key []byte,
	check func(a, b []byte) bool, ids []uint32,
) []uint32 {
	if mem == nil {
		return ids
	}
	kvs, ok := mem.Get(bucketID)
	if !ok {
		return ids
	}
	for k, v := range kvs {
		if check(strutil.String2ByteSlice(k), key) {
			ids = append(ids, v)
		}
	}
	return ids
}

func (s *indexKVStore) findValue(bucketID uint32, key []byte, ids []uint32) ([]uint32, error) {
	id, ok, err := s.GetValue(bucketID, key)
	if err != nil {
		return nil, err
	}
	if ok {
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *indexKVStore) getValueFromMem(mem *imap.IntMap[map[string]uint32], bucketID uint32, key []byte) (uint32, bool) {
	if mem == nil {
		return 0, false
	}
	kvs, ok := mem.Get(bucketID)
	if !ok {
		return 0, false
	}
	id, ok := kvs[strutil.ByteSlice2String(key)]
	return id, ok
}

func (s *indexKVStore) getValuesFromMem(mem *imap.IntMap[map[string]uint32], bucketID uint32, result []uint32) []uint32 {
	if mem == nil {
		return result
	}
	kvs, ok := mem.Get(bucketID)
	if !ok {
		return result
	}
	for _, v := range kvs {
		result = append(result, v)
	}
	return result
}
