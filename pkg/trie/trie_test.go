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

package trie

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type kvPair struct {
	keys   [][]byte
	values []uint32
}

func (p kvPair) Len() int {
	return len(p.keys)
}

func (p kvPair) Less(i, j int) bool {
	return bytes.Compare(p.keys[i], p.keys[j]) < 0
}

func (p kvPair) Swap(i, j int) {
	p.keys[i], p.keys[j] = p.keys[j], p.keys[i]
	p.values[i], p.values[j] = p.values[j], p.values[i]
}

func (p kvPair) Sort() {
	sort.Sort(p)
}

func newTestIPs(batchSize int) (ips [][]byte, ranks []uint32) {
	var count int
	for x := 10; x > 0; x-- {
		for y := 1; y < batchSize; y++ {
			for z := batchSize - 1; z > 0; z-- {
				ips = append(ips, []byte(fmt.Sprintf("%d.%d.%d.%d", x, y, y, z)))
				count++
				ranks = append(ranks, uint32(count))
			}
		}
	}
	kvPair{keys: ips, values: ranks}.Sort()
	return
}

func TestTriePool(t *testing.T) {
	trie := GetTrie()
	assert.NotNil(t, trie)
	PutTrie(trie)
	PutTrie(nil)
}

func TestBuilder_Reset(t *testing.T) {
	builder := NewBuilder()
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < 5; i++ {
		var keys [][]byte
		var values []uint32
		count := uint32(r.Uint64() % 10000)
		for i := uint32(0); i < count; i++ {
			var scratch [8]byte
			binary.LittleEndian.PutUint64(scratch[:], r.Uint64())
			keys = append(keys, scratch[:])
			values = append(values, i)
		}
		kvPair{keys: keys, values: values}.Sort()
		builder.Build(keys, values)
		w := bytes.NewBuffer([]byte{})
		err := builder.Write(w)
		assert.NoError(t, err)
		data := w.Bytes()
		tree := NewTrie()
		assert.NoError(t, tree.UnmarshalBinary(data))
		c := 0
		for idx, key := range keys {
			val, ok := tree.Get(key)
			assert.True(t, ok)
			assert.Equal(t, values[idx], val)
			c++
		}
		assert.Equal(t, len(keys), c)
		builder.Reset()
	}
}

func TestIterator_SeekToLast(t *testing.T) {
	ips, ranks := newTestIPs(2)
	builder := NewBuilder()
	builder.Build(ips, ranks)
	tree := builder.Trie()
	itr := tree.NewIterator()

	itr.SeekToLast()
	var count = len(ips)

	for itr.Valid() {
		assert.Equal(t, ips[count-1], itr.Key())
		assert.Equal(t, ranks[count-1], itr.Value())

		itr.Prev()
		count--
	}
	assert.Zero(t, count)
	itr.Prev()
	assert.False(t, itr.Valid())
}

func TestTrie_Get(t *testing.T) {
	ips, ranks := newTestIPs(2)
	builder := NewBuilder()
	builder.Build(ips, ranks)
	tree := builder.Trie()
	for idx, ip := range ips {
		value, ok := tree.Get(ip)
		assert.True(t, ok)
		assert.Equal(t, ranks[idx], value)
	}
	for _, ip := range ips {
		_, ok := tree.Get(ip[:len(ip)-3])
		assert.False(t, ok)
	}
}

func TestTrie_UnmarshalBinary(t *testing.T) {
	ips, ranks := newTestIPs(2)
	builder := NewBuilder()
	builder.Build(ips, ranks)
	w := bytes.NewBuffer([]byte{})
	err := builder.Write(w)
	data := w.Bytes()
	assert.NoError(t, err)
	assert.Equal(t, len(data), builder.MarshalSize())

	// unmarshal
	tree2 := NewTrie()
	assert.NoError(t, tree2.UnmarshalBinary(data))
	count := 0
	for idx, ip := range ips {
		val, ok := tree2.Get(ip)
		assert.True(t, ok)
		assert.Equal(t, ranks[idx], val)
		count++
	}
	assert.Equal(t, len(ips), count)
}

func Test_Trie_ASCII(t *testing.T) {
	var keys [][]byte
	var values []uint32
	for i := 0; i < 255; i++ {
		keys = append(keys, []byte{12, 13, 14, 15, uint8(i)})
		values = append(values, uint32(1))
	}
	builder := NewBuilder()
	builder.Build(keys, values)
	tree := builder.Trie()
	itr := tree.NewIterator()

	itr.Seek([]byte{10, 12, 13, 14})
	itr.Seek([]byte{11, 12, 13, 14})
	itr.Seek([]byte{12})
	itr.Seek([]byte{12, 13})
	itr.Seek([]byte{12, 13, 14})
	itr.Seek([]byte{12, 13, 14, 14})
	itr.Seek([]byte{12, 13, 14, 15, 1})
	itr.Seek([]byte{12, 13, 14, 15, 1, 1})
	itr.Seek([]byte{12, 13, 14, 16})
	itr.Seek([]byte{13, 12, 13, 14})
}

func Test_Trie_words(t *testing.T) {
	var keys [][]byte
	var values []uint32
	keysString := []string{
		"a", "ab", "b", "abc", "abcdefgh", "abcdefghijklmnopqrstuvwxyz", "abcdefghijkl", "zzzzzz", "ice",
	}
	for idx, key := range keysString {
		keys = append(keys, []byte(key))
		values = append(values, uint32(idx))
	}
	kvPair{keys: keys, values: values}.Sort()
	builder := NewBuilder()
	builder.Build(keys, values)
	tree := builder.Trie()
	examples := []struct {
		input string
		ok    bool
	}{
		{"a", true},
		{"ab", true},
		{"b", true},
		{"bb", false},
		{"abc", true},
		{"abcd", false},
		{"abcdefghijklmnopqrstuvwxyz", true},
		{"abcdefghijkl", true},
		{"abcdefghijklm", false},
		{"zzzzzz", true},
		{"zzzzz", false},
		{"zzzzzzz", false},
		{"i", false},
		{"ice", true},
		{"ic", false},
		{"ices", false},
	}

	for _, example := range examples {
		_, ok := tree.Get([]byte(example.input))
		assert.Equalf(t, example.ok, ok, example.input)
	}
}

func Test_Trie_idList(t *testing.T) {
	var ids [][]byte
	var indexes []uint32
	for i := 0; i < 10000; i++ {
		// windows rand wrond data
		// rand.Seed(time.Now().UnixNano())
		id := fmt.Sprintf("%d", rand.Int63n(math.MaxInt64))
		ids = append(ids, []byte(id))
		indexes = append(indexes, uint32(i))
	}
	kvPair{keys: ids, values: indexes}.Sort()
	builder := NewBuilder()
	builder.Build(ids, indexes)
	tree := builder.Trie()
	if len(indexes) == 0 || len(ids) == 0 {
		panic("length is zero")
	}
	count := 0
	for idx := range ids {
		value, ok := tree.Get(ids[idx])
		assert.True(t, ok)
		assert.Equal(t, indexes[idx], value)
		count++
	}
	assert.Equal(t, len(ids), count)
	w := bytes.NewBuffer([]byte{})
	err := builder.Write(w)
	data := w.Bytes()
	assert.NoError(t, err)
	tree2 := NewTrie()
	assert.NoError(t, tree2.UnmarshalBinary(data))

	itr := tree2.NewIterator()
	itr.SeekToFirst()
	// prev is empty
	itr.Prev()
	assert.False(t, itr.Valid())
	var idx = 0
	for itr.Valid() {
		assert.Equal(t, indexes[idx], itr.Value())
		assert.Equal(t, ids[idx], itr.Key())
		itr.Next()
		idx++
	}
	itr.Next()
	assert.False(t, itr.Valid())
}

func assertTestData(t *testing.T, path string) {
	var keys [][]byte
	var values []uint32
	f, err := os.Open(path)
	assert.Nil(t, err)
	r, err := gzip.NewReader(f)
	assert.Nil(t, err)

	data, err := io.ReadAll(r)
	assert.Nil(t, err)
	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		keys = append(keys, []byte(line))
		values = append(values, uint32(i))
	}
	kvPair{keys: keys, values: values}.Sort()
	builder := NewBuilder()
	builder.Build(keys, values)
	tree := builder.Trie()
	if len(keys) == 0 || len(values) == 0 {
		panic("length is zero")
	}
	for idx := range keys {
		value, ok := tree.Get(keys[idx])
		assert.True(t, ok)
		assert.Equal(t, values[idx], value)
	}

	w := bytes.NewBuffer([]byte{})
	err = builder.Write(w)
	data = w.Bytes()
	assert.NoError(t, err)
	tree2 := NewTrie()
	assert.NoError(t, tree2.UnmarshalBinary(data))

	itr := tree2.NewIterator()
	itr.SeekToFirst()
	var idx = 0
	for itr.Valid() {
		assert.Equal(t, values[idx], itr.Value())
		assert.Equal(t, keys[idx], itr.Key())
		itr.Next()
		idx++
	}
	assert.Equal(t, len(keys), idx)
}

func Test_Trie_TestData_Words(t *testing.T) {
	assertTestData(t, "testdata/words.txt.gz")
}

func Test_Trie_TestData_UUID(t *testing.T) {
	assertTestData(t, "testdata/uuid.txt.gz")
}

func Test_Trie_TestData_Hsk_words(t *testing.T) {
	assertTestData(t, "testdata/hsk_words.txt.gz")
}

func TestBuildPrefixKeys(t *testing.T) {
	keys := [][]byte{
		{1},
		{1, 1},
		{1, 1, 1},
		{1, 1, 1, 1},
		{2},
		{2, 2},
		{2, 2, 2},
	}
	vals := genSeqVals(len(keys))
	checker := newFullSuRFChecker(keys, vals)
	buildAndCheckSuRF(t, keys, vals, checker)
}

func TestBuildCompressPath(t *testing.T) {
	keys := [][]byte{
		{1, 1, 1},
		{1, 1, 1, 2, 2},
		{1, 1, 1, 2, 2, 2},
		{1, 1, 1, 2, 2, 3},
		{2, 1, 3},
		{2, 2, 3},
		{2, 3, 1, 1, 1, 1, 1, 1, 1},
		{2, 3, 1, 1, 1, 2, 2, 2, 2},
	}
	vals := genSeqVals(len(keys))
	checker := newFullSuRFChecker(keys, vals)
	buildAndCheckSuRF(t, keys, vals, checker)
}

func TestBuildSuffixKeys(t *testing.T) {
	keys := [][]byte{
		bytes.Repeat([]byte{1}, 30),
		bytes.Repeat([]byte{2}, 30),
		bytes.Repeat([]byte{3}, 30),
		bytes.Repeat([]byte{4}, 30),
	}
	vals := genSeqVals(len(keys))
	checker := newFullSuRFChecker(keys, vals)
	buildAndCheckSuRF(t, keys, vals, checker)
}

func TestRandomKeysSparse(t *testing.T) {
	keys := genRandomKeys(200, 60, 0) // 2000000, 60, 0
	vals := genSeqVals(len(keys))
	checker := newFullSuRFChecker(keys, vals)
	buildAndCheckSuRF(t, keys, vals, checker)
}

func TestRandomKeysPrefixGrowth(t *testing.T) {
	keys := genRandomKeys(100, 10, 20) // 100, 10, 200
	vals := genSeqVals(len(keys))
	checker := newFullSuRFChecker(keys, vals)
	buildAndCheckSuRF(t, keys, vals, checker)
}

func TestSeekKeys(t *testing.T) {
	keys := genRandomKeys(50, 10, 30) // 50, 10, 300
	insert, seek, vals := splitKeys(keys)
	checker := func(t *testing.T, surf SuccinctTrie) {
		it := surf.NewIterator()
		for i, k := range seek {
			it.Seek(k)
			require.True(t, it.Valid())
			require.True(t, it.Value() <= vals[i])
		}
	}

	buildAndCheckSuRF(t, insert, vals, checker)
}

func genSeqVals(n int) []uint32 {
	vals := make([]uint32, n)
	for i := 0; i < n; i++ {
		vals[i] = uint32(i)
	}
	return vals
}

func buildAndCheckSuRF(t *testing.T, keys [][]byte, vals []uint32, checker func(t *testing.T, surf SuccinctTrie)) {
	b := NewBuilder()
	b.Build(keys, vals)
	surf := b.Trie()
	checker(t, surf)
}

func newFullSuRFChecker(keys [][]byte, vals []uint32) func(t *testing.T, surf SuccinctTrie) {
	return func(t *testing.T, surf SuccinctTrie) {
		for i, k := range keys {
			val, ok := surf.Get(k)
			require.True(t, ok)
			require.EqualValues(t, vals[i], val)
		}

		var i int
		it := surf.NewIterator()
		for it.SeekToFirst(); it.Valid(); it.Next() {
			require.Truef(t, bytes.HasPrefix(keys[i], it.Key()), "%v %v %d", keys[i], it.Key(), i)
			require.EqualValues(t, vals[i], it.Value())
			i++
		}
		require.Equal(t, len(keys), i)

		i = len(keys) - 1
		for it.SeekToLast(); it.Valid(); it.Prev() {
			require.True(t, bytes.HasPrefix(keys[i], it.Key()))
			require.EqualValues(t, vals[i], it.Value())
			i--
		}
		require.Equal(t, -1, i)

		for i, k := range keys {
			it.Seek(k)
			require.EqualValues(t, k, it.Key())
			require.EqualValues(t, vals[i], it.Value())
		}
	}
}

// max key length is `initLen * (round + 1)`
// max result size is (initSize + initSize * (round + 1)) * (round + 1) / 2
// you can use small round (0 is allowed) to generate a sparse key set,
// or use a large round to generate a key set which has many common prefixes.
func genRandomKeys(initSize, initLen, round int) [][]byte {
	keys := make([][]byte, initSize)
	random := rand.New(rand.NewSource(rand.Int63n(math.MaxInt64)))

	for i := range keys {
		keys[i] = make([]byte, rand.Intn(initLen)+1)
		rand.Read(keys[i])
	}

	for r := 1; r <= round; r++ {
		for i := 0; i < initSize*r; i++ {
			k := make([]byte, len(keys[i])+rand.Intn(initLen)+1)
			copy(k, keys[i])
			random.Read(k[len(keys[i]):])
			keys = append(keys, k)
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})

	var prev []byte
	result := keys[:0]
	for _, k := range keys {
		if bytes.Equal(prev, k) {
			continue
		}
		prev = k
		result = append(result, k)
	}
	for i := len(result); i < len(keys); i++ {
		keys[i] = nil
	}

	return result
}

func splitKeys(keys [][]byte) (a, b [][]byte, aIdx []uint32) {
	a = keys[:0]
	b = make([][]byte, 0, len(keys)/2)
	aIdx = make([]uint32, 0, len(keys)/2)
	for i := 0; i < len(keys) & ^1; i += 2 {
		b = append(b, keys[i])
		a = append(a, keys[i+1])
		aIdx = append(aIdx, uint32(i+1))
	}
	return
}
