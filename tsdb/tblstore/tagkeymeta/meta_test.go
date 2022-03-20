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

package tagkeymeta

import (
	"fmt"
	"sync"
	"testing"

	"github.com/lindb/lindb/kv"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func Test_newTagKeyMeta_error_cases(t *testing.T) {
	// case1: block length too short
	_, err := newTagKeyMeta([]byte{})
	assert.Error(t, err)

	// case2: data corruption
	_, err = newTagKeyMeta([]byte{
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
		4, 4, 4, 4,
		5})
	assert.Error(t, err)
}

var testOnce sync.Once
var testData []byte

func buildTestTrieData() []byte {
	testOnce.Do(func() {
		kvFlusher := kv.NewNopFlusher()
		flusher, _ := NewFlusher(kvFlusher)
		flusher.EnsureSize(10 * 10 * 10 * 10)

		count := uint32(0)
		for a := 1; a <= 10; a++ {
			for b := 1; b <= 10; b++ {
				for c := 1; c <= 10; c++ {
					for d := 1; d <= 10; d++ {
						flusher.FlushTagValue([]byte(fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)), count)
						count++
					}
				}
			}
		}
		_ = flusher.FlushTagKeyID(1, 10*10*10*10)

		testData = kvFlusher.Bytes()
	})
	return testData
}

func TestTagKeyMeta_TagValueIDSeq(t *testing.T) {
	tagKeyMeta, err := newTagKeyMeta(buildTestTrieData())
	assert.NoError(t, err)
	assert.NotNil(t, tagKeyMeta)

	assert.Equal(t, uint32(10*10*10*10), tagKeyMeta.TagValueIDSeq())
}

func TestTagKeyMetas_GetTagValueIDs(t *testing.T) {
	// normal cases
	meta, _ := newTagKeyMeta(buildTestTrieData())
	var tagKeyMetas = TagKeyMetas{meta}
	bitmap, err := tagKeyMetas.GetTagValueIDs()
	assert.NoError(t, err)
	assert.Equal(t, uint64(10*10*10*10), bitmap.GetCardinality())

	// bitmap corrupted
	tagKeyMetaImpl := meta.(*tagKeyMeta)
	tagKeyMetaImpl.bitmapData = nil
	_, err = tagKeyMetas.GetTagValueIDs()
	assert.Error(t, err)
}

func TestTagKeyMeta_FindTagValueID(t *testing.T) {
	meta, _ := newTagKeyMeta(buildTestTrieData())
	tagValueIDs := meta.FindTagValueID("1.1.1.1")
	assert.Len(t, tagValueIDs, 1)
	assert.Equal(t, uint32(0), tagValueIDs[0])

	tagValueIDs = meta.FindTagValueID("1.1.1.")
	assert.Len(t, tagValueIDs, 0)

	tagValueIDs = meta.FindTagValueID("9.9.9.9")
	assert.Equal(t, uint32(8888), tagValueIDs[0])

	tagValueIDs = meta.FindTagValueID("9.9.9.9.9")
	assert.Len(t, tagValueIDs, 0)
}

func TestTagKeyMeta_FindTagValueIDs(t *testing.T) {
	meta, _ := newTagKeyMeta(buildTestTrieData())

	expected := []uint32{
		1111, 2222, 3333,
		4444, 5555, 6666, 7777}
	assert.Equal(t, expected, meta.FindTagValueIDs([]string{
		"2.2.2.2", "3.3.3.3", "4.4.4.4",
		"5.5.5.5", "6.6.6.6", "7.7.7.7", "8.8.8.8"}))
}

func TestTagKeyMeta_FindTagValueIDsByLike(t *testing.T) {
	meta, _ := newTagKeyMeta(buildTestTrieData())
	// case1: exactly search
	tagValueIDs := meta.FindTagValueIDsByLike("1.1.1.1")
	assert.Len(t, tagValueIDs, 1)

	// case2: prefix search
	// 1.1.1.1
	// 1.1.1.10
	// 1.1.1.2
	// 1.1.1.3
	// 1.1.1.4
	// 1.1.1.5
	// 1.1.1.6
	// 1.1.1.7
	// 1.1.1.8
	// 1.1.1.9
	assert.Equal(t, []uint32{0, 9, 1, 2, 3, 4, 5, 6, 7, 8},
		meta.FindTagValueIDsByLike("1.1.1.*"))

	// case3: suffix search
	// 1.1.1.1
	// 10.1.1.1
	// 2.1.1.1
	// 3.1.1.1
	// 4.1.1.1
	// 5.1.1.1
	// 6.1.1.1
	// 7.1.1.1
	// 8.1.1.1
	// 9.1.1.1
	assert.Equal(t, []uint32{0, 0x2328, 0x3e8, 0x7d0, 0xbb8, 0xfa0, 0x1388, 0x1770, 0x1b58, 0x1f40},
		meta.FindTagValueIDsByLike("*.1.1.1"))

	// case4: middle search
	assert.Len(t, meta.FindTagValueIDsByLike("*.1.1.*"), 100)

	// case5: nil search
	assert.Len(t, meta.FindTagValueIDsByLike(""), 0)
}

func TestTagKeyMeta_FindTagValueIDsByRegex(t *testing.T) {
	meta, _ := newTagKeyMeta(buildTestTrieData())

	// case1: bad pattern
	assert.Len(t, meta.FindTagValueIDsByRegex("1["), 0)

	// case2: prefix regex
	assert.Len(t, meta.FindTagValueIDsByRegex("1\\.1\\.1\\.[1-3]"), 4)

	// case3: regex all
	assert.Len(t, meta.FindTagValueIDsByRegex(".*"), 10000)
}

func TestTagKeyMeta_CollectTagValues(t *testing.T) {
	meta, _ := newTagKeyMeta(buildTestTrieData())
	// case1: normal
	bitmap := roaring.BitmapOf(2, 1, 3, 4, 5, 111111, 222222)
	tagValues := make(map[uint32]string)
	assert.NoError(t, meta.CollectTagValues(bitmap, tagValues))
	assert.Len(t, tagValues, 5)

	// case2: bitmap corrupted
	tagKeyMetaImpl := meta.(*tagKeyMeta)
	tagKeyMetaImpl.bitmapData = nil
	assert.Error(t, tagKeyMetaImpl.CollectTagValues(bitmap, tagValues))

	// case3: bitmap not found
	meta, _ = newTagKeyMeta(buildTestTrieData())
	bitmap = roaring.BitmapOf(10001, 10002, 10003)
	tagValues = make(map[uint32]string)
	assert.Len(t, tagValues, 0)
	assert.NoError(t, meta.CollectTagValues(bitmap, tagValues))

	// case4: offset corrupted
	meta, _ = newTagKeyMeta(buildTestTrieData())
	tagKeyMetaImpl = meta.(*tagKeyMeta)
	tagKeyMetaImpl.offsetsData = nil
	assert.Error(t, tagKeyMetaImpl.CollectTagValues(roaring.BitmapOf(1, 2), tagValues))
}

func TestTagKeyMeta_CollectTagValues_WithSwap(t *testing.T) {
	kvFlusher := kv.NewNopFlusher()
	flusher, _ := NewFlusher(kvFlusher)
	flusher.FlushTagValue([]byte("x"), 1)
	flusher.FlushTagValue([]byte("t"), 2)
	flusher.FlushTagValue([]byte("sfd"), 3)
	flusher.FlushTagValue([]byte("b"), 4)
	flusher.FlushTagValue([]byte("bc"), 5)
	_ = flusher.FlushTagKeyID(1, 1)
	data := kvFlusher.Bytes()
	meta, _ := newTagKeyMeta(data)
	tagValues := make(map[uint32]string)
	assert.NoError(t, meta.CollectTagValues(roaring.BitmapOf(1, 2, 3, 4, 5), tagValues))
	assert.Len(t, tagValues, 5)

	assert.Len(t, meta.FindTagValueID("bcd"), 0)
}

func TestTagKeyMeta_Error(t *testing.T) {
	kvFlusher := kv.NewNopFlusher()
	flusher, _ := NewFlusher(kvFlusher)
	flusher.FlushTagValue([]byte("x"), 1)
	flusher.FlushTagValue([]byte("t"), 2)
	_ = flusher.FlushTagKeyID(1, 1)
	data := kvFlusher.Bytes()
	meta, _ := newTagKeyMeta(data)
	metaImpl := meta.(*tagKeyMeta)

	// destroy the meta trie data
	metaImpl.trieBlock = append([]byte{1, 2, 3, 4}, metaImpl.trieBlock...)

	// FindTagValueIDsByRegex error
	assert.Len(t, meta.FindTagValueIDsByRegex("x"), 0)
	// FindTagValueIDsByLike error
	assert.Len(t, meta.FindTagValueIDsByLike("x*"), 0)
	assert.Len(t, meta.FindTagValueIDsByLike("*x*"), 0)
	assert.Len(t, meta.FindTagValueIDsByLike("*x"), 0)
	// FindTagValueIDs error
	assert.Len(t, meta.FindTagValueID("x"), 0)
	// FindTagValueIDs error
	assert.Len(t, meta.FindTagValueIDs([]string{"x"}), 0)
	// CollectTagValues error
	assert.Error(t, meta.CollectTagValues(roaring.BitmapOf(1, 2), map[uint32]string{}))
}
