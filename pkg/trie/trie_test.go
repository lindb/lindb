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

package trie_test

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

	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/trie"
)

type kvPair struct {
	keys   [][]byte
	values [][]byte
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

func newTestIPs(batchSize int) (ips, ranks [][]byte) {
	var count int
	for x := 10; x > 0; x-- {
		for y := 1; y < batchSize; y++ {
			for z := batchSize - 1; z > 0; z-- {
				ips = append(ips, []byte(fmt.Sprintf("%d.%d.%d.%d", x, y, y, z)))
				count++
				var rank = make([]byte, 4)
				binary.LittleEndian.PutUint32(rank, uint32(count))
				ranks = append(ranks, rank)
			}
		}
	}
	kvPair{keys: ips, values: ranks}.Sort()
	return
}

func TestBuilder_Reset(t *testing.T) {
	ips, ranks := newTestIPs(2)
	builder := trie.NewBuilder()
	assert.True(t, sort.IsSorted(&kvPair{keys: ips, values: ranks}))

	for i := 0; i < 3; i++ {
		tree := builder.Build(ips, ranks, 3)
		itr := tree.NewIterator()
		itr.SeekToFirst()
		var count int
		for itr.Valid() {
			assert.Equal(t, ips[count], itr.Key())
			assert.Equal(t, ranks[count][:3], itr.Value())

			itr.Next()
			count++
		}
		assert.Equal(t, len(ips), count)

		itr.Next()
		assert.False(t, itr.Valid())
		builder.Reset()
	}
}

func TestIterator_SeekToLast(t *testing.T) {
	ips, ranks := newTestIPs(2)
	builder := trie.NewBuilder()
	tree := builder.Build(ips, ranks, 3)
	itr := tree.NewIterator()

	itr.SeekToLast()
	var count = len(ips)

	for itr.Valid() {
		assert.Equal(t, ips[count-1], itr.Key())
		assert.Equal(t, ranks[count-1][:3], itr.Value())

		itr.Prev()
		count--
	}
	assert.Zero(t, count)
	itr.Prev()
	assert.False(t, itr.Valid())
}

func TestTrie_Get(t *testing.T) {
	ips, ranks := newTestIPs(2)
	builder := trie.NewBuilder()
	tree := builder.Build(ips, ranks, 3)
	for _, ip := range ips {
		_, ok := tree.Get(ip)
		assert.True(t, ok)
	}
	for _, ip := range ips {
		_, ok := tree.Get(ip[:len(ip)-3])
		assert.False(t, ok)
	}
}

func TestTrie_UnmarshalBinary(t *testing.T) {
	ips, ranks := newTestIPs(2)
	builder := trie.NewBuilder()
	tree := builder.Build(ips, ranks, 3)
	data, err := tree.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, tree.MarshalSize(), int64(len(data)))

	// unmarshal
	tree2 := trie.NewTrie()
	assert.NoError(t, tree2.UnmarshalBinary(data))
	for _, ip := range ips {
		_, ok := tree2.Get(ip)
		assert.True(t, ok)
	}
	data2, _ := tree2.MarshalBinary()
	assert.Equal(t, len(data), len(data2))
}

func Test_Trie_ASCII(t *testing.T) {
	var keys [][]byte
	var values [][]byte
	for i := 0; i < 255; i++ {
		keys = append(keys, []byte{12, 13, 14, 15, uint8(i)})
		values = append(values, []byte{1})
	}
	builder := trie.NewBuilder()
	tree := builder.Build(keys, values, 1)
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
	var values [][]byte
	keysString := []string{
		"a", "ab", "b", "abc", "abcdefgh", "abcdefghijklmnopqrstuvwxyz", "abcdefghijkl", "zzzzzz", "ice",
	}
	for idx, key := range keysString {
		keys = append(keys, []byte(key))
		values = append(values, []byte{uint8(idx)})
	}
	kvPair{keys: keys, values: values}.Sort()
	builder := trie.NewBuilder()
	tree := builder.Build(keys, values, 1)
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
	var scratch [4]byte
	var ids [][]byte
	var indexes [][]byte
	for i := 0; i < 10000; i++ {
		// windows rand wrond data
		// rand.Seed(time.Now().UnixNano())
		id := fmt.Sprintf("%d", rand.Int63n(math.MaxInt64))
		ids = append(ids, []byte(id))
		binary.LittleEndian.PutUint32(scratch[:], uint32(i))
		indexes = append(indexes, append([]byte{}, scratch[:]...))
	}
	kvPair{keys: ids, values: indexes}.Sort()
	builder := trie.NewBuilder()
	tree := builder.Build(ids, indexes, 4)
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
	fmt.Println(count)
	data, err := tree.MarshalBinary()
	assert.Nil(t, err)
	tree2 := trie.NewTrie()
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
	fmt.Println(idx)
	itr.Next()
	assert.False(t, itr.Valid())
}

func assertTestData(t *testing.T, path string) {
	var scratch [4]byte
	var keys [][]byte
	var values [][]byte
	f, err := os.Open(path)
	assert.Nil(t, err)
	r, err := gzip.NewReader(f)
	assert.Nil(t, err)

	data, err := io.ReadAll(r)
	assert.Nil(t, err)
	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		keys = append(keys, []byte(line))
		binary.LittleEndian.PutUint32(scratch[:], uint32(i))
		values = append(values, append([]byte{}, scratch[:]...))
	}
	kvPair{keys: keys, values: values}.Sort()
	builder := trie.NewBuilder()
	tree := builder.Build(keys, values, 4)

	if len(keys) == 0 || len(values) == 0 {
		panic("length is zero")
	}
	for idx := range keys {
		value, ok := tree.Get(keys[idx])
		assert.True(t, ok)
		assert.Equal(t, values[idx], value)
	}

	data, err = tree.MarshalBinary()
	assert.Nil(t, err)
	tree2 := trie.NewTrie()
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
