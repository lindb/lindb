package trie_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"testing"

	"github.com/lindb/lindb/pkg/trie"

	"github.com/stretchr/testify/assert"
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

func newTestIPs() ([][]byte, [][]byte) {
	var ips [][]byte
	var ranks [][]byte
	var count int
	for x := 10; x > 0; x-- {
		for y := 1; y < 1<<8; y++ {
			for z := 1<<8 - 1; z > 0; z-- {
				ips = append(ips, []byte(fmt.Sprintf("%d.%d.%d.%d", x, y, y, z)))
				count++
				var rank = make([]byte, 4)
				binary.LittleEndian.PutUint32(rank, uint32(count))
				ranks = append(ranks, rank)
			}
		}
	}
	kvPair{keys: ips, values: ranks}.Sort()
	return ips, ranks
}

func TestBuilder_Reset(t *testing.T) {
	ips, ranks := newTestIPs()
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
	ips, ranks := newTestIPs()
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
	ips, ranks := newTestIPs()
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
	ips, ranks := newTestIPs()
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
