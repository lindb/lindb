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

package model

import (
	"bytes"
	"container/heap"
	"encoding/binary"
	"io"
	"math"
	"regexp"
	"sort"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/pkg/trie"
)

// for testing
var (
	getTrieFn = trie.GetTrie
)

// KVs represents key/value pairs.
type KVs struct {
	Keys [][]byte
	IDs  []uint32
}

func (m *KVs) Len() int { return len(m.Keys) }

func (m *KVs) Less(i, j int) bool { return bytes.Compare(m.Keys[i], m.Keys[j]) < 0 }

func (m *KVs) Swap(i, j int) {
	m.Keys[i], m.Keys[j] = m.Keys[j], m.Keys[i]
	m.IDs[i], m.IDs[j] = m.IDs[j], m.IDs[i]
}

type trieEntry struct {
	tree trie.SuccinctTrie
	buf  []byte
}

type tries []*trieEntry

// TrieBucket represents a bucket include multiple succinct trie.
type TrieBucket struct {
	kvs       tries
	blockSize int
}

// NewTrieBucket creates a trie bucket with default block size.
func NewTrieBucket() *TrieBucket {
	return NewTrieBucketWithBlockSize(math.MaxUint16)
}

// NewTrieBucketWithBlockSize create a trie bucket with block size.
func NewTrieBucketWithBlockSize(blockSize int) *TrieBucket {
	return &TrieBucket{
		blockSize: blockSize,
	}
}

// Unmarshal unmarshals kv tries.
func (b *TrieBucket) Unmarshal(block []byte) error {
	for len(block) > 0 {
		size := binary.LittleEndian.Uint32(block[:4])
		tree := getTrieFn()
		end := 4 + size
		err := tree.UnmarshalBinary(block[4:end])
		if err != nil {
			return err
		}
		b.kvs = append(b.kvs, &trieEntry{tree: tree, buf: block[:end]})
		block = block[end:]
	}
	return nil
}

// Write writes kv tries based on block size.
func (b *TrieBucket) Write(w io.Writer) error {
	sort.Slice(b.kvs, func(i, j int) bool {
		return b.kvs[i].tree.Size() > b.kvs[j].tree.Size()
	})
	var pendingKVs tries
	for idx := range b.kvs {
		tree := b.kvs[idx]
		if tree.tree.Size() >= b.blockSize {
			if _, err := w.Write(tree.buf); err != nil {
				return err
			}
		} else {
			pendingKVs = append(pendingKVs, tree)
		}
	}
	pendingSize := len(pendingKVs)
	switch pendingSize {
	case 0:
		// if no pending kvs return
		return nil
	case 1:
		if _, err := w.Write(pendingKVs[0].buf); err != nil {
			return err
		}
	default:
		var keys [][]byte
		var ids []uint32

		for _, kv := range pendingKVs {
			itr := kv.tree.NewPrefixIterator(nil)
			for itr.Valid() {
				key := itr.Key()
				// NOTE: need copy key, because trie iterator reuse key when iterate
				k := make([]byte, len(key))
				copy(k, key)
				keys = append(keys, k)
				ids = append(ids, itr.Value())
				itr.Next()
			}
		}

		builder := NewTrieBucketBuilder(b.blockSize, w)
		return builder.Write(keys, ids)
	}

	return nil
}

// Release releases bucket resource.
func (b *TrieBucket) Release() {
	for idx := range b.kvs {
		trie.PutTrie(b.kvs[idx].tree)
	}
}

// GetValue returns value by key.
func (b *TrieBucket) GetValue(key []byte) (id uint32, ok bool) {
	for _, kvs := range b.kvs {
		id, ok = kvs.tree.Get(key)
		if ok {
			return
		}
	}
	return
}

// GetValues returns all values.
func (b *TrieBucket) GetValues() (ids []uint32) {
	for _, kvs := range b.kvs {
		ids = append(ids, kvs.tree.Values()...)
	}
	return
}

// CollectKVs collects key/value pairs by values.
func (b *TrieBucket) CollectKVs(values *roaring.Bitmap, result map[uint32]string) {
	// TODO: need refactor?(implements find key by value)
	for _, kv := range b.kvs {
		itr := kv.tree.NewPrefixIterator(nil)
		for itr.Valid() {
			val := itr.Value()
			if values.Contains(val) {
				result[val] = string(itr.Key())
				values.Remove(val)
			}
			if values.IsEmpty() {
				return
			}
			itr.Next()
		}
	}
}

func (b *TrieBucket) Suggest(prefix string, limit int) (rs []string) {
	prefixBytes := strutil.String2ByteSlice(prefix)
	var its []*trie.PrefixIterator
	for _, kv := range b.kvs {
		its = append(its, kv.tree.NewPrefixIterator(prefixBytes))
	}
	it := NewMergedIterator(its)
	for it.HasNext() {
		key := it.Key()
		rs = append(rs, string(key))
		if len(rs) >= limit {
			return
		}
	}
	return
}

// FindValuesByRegexp returns values by regexp expression.
func (b *TrieBucket) FindValuesByRegexp(rp *regexp.Regexp, ids []uint32) []uint32 {
	literalPrefix, _ := rp.LiteralPrefix()
	literalPrefixByte := strutil.String2ByteSlice(literalPrefix)
	for _, kv := range b.kvs {
		itr := kv.tree.NewPrefixIterator(literalPrefixByte)
		for itr.Valid() {
			if rp.Match(itr.Key()) {
				ids = append(ids, itr.Value())
			}
			itr.Next()
		}
	}
	return ids
}

// FindValuesByLike returns values by like expression.
func (b *TrieBucket) FindValuesByLike(prefix, subKey []byte, check func(a, b []byte) bool, ids []uint32) []uint32 {
	for _, kv := range b.kvs {
		itr := kv.tree.NewPrefixIterator(prefix)
		for itr.Valid() {
			if check(itr.Key(), subKey) {
				ids = append(ids, itr.Value())
			}
			itr.Next()
		}
	}
	return ids
}

// mergedIterator iterates over some iterator in key order
type mergedIterator struct {
	its []*trie.PrefixIterator
	pq  priorityQueue

	curKey []byte
}

// NewMergedIterator create merged iterator for multi iterators
func NewMergedIterator(its []*trie.PrefixIterator) *mergedIterator {
	it := &mergedIterator{
		its: its,
	}
	it.initQueue()
	return it
}

// initQueue initializes the priority queue
func (m *mergedIterator) initQueue() {
	i := 0
	for _, it := range m.its {
		if it.Valid() {
			m.pq = append(m.pq, &item{
				it:    it,
				key:   it.Key(),
				index: i,
			})
			it.Next()
			i++
		}
	}
	if len(m.pq) > 0 {
		heap.Init(&m.pq)
	}
}

// HasNext returns if the iteration has more element.
// It returns false if the iterator is exhausted.
func (m *mergedIterator) HasNext() bool {
	result := len(m.pq) > 0
	if result {
		// pop item and get value
		val := heap.Pop(&m.pq)
		item := val.(*item)
		m.curKey = item.key

		// if it has value, push back queue and adjust priority
		it := item.it
		if it.Valid() {
			item.key = it.Key()
			m.pq.Push(item)
			m.pq.update(item)

			it.Next()
		}
	}
	return result
}

// Key returns the key of the current key
func (m *mergedIterator) Key() []byte {
	return m.curKey
}

// item represents an item under priority queue, using key as priority.
type item struct {
	it *trie.PrefixIterator

	key []byte

	index int
}

// priorityQueue implements heap.Interface and holds Items.
type priorityQueue []*item

// Len returns the number of elements in priority queue
func (pq priorityQueue) Len() int { return len(pq) }

// Less compares key of item
func (pq priorityQueue) Less(i, j int) bool { return bytes.Compare(pq[i].key, pq[j].key) < 0 }

// Swap swaps the elements with indexes i and j.
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = j
	pq[j].index = i
}

// Push pushes a item into priority queue
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns element of length -1
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority by the key of item
func (pq *priorityQueue) update(item *item) {
	heap.Fix(pq, item.index)
}
