package metricsmeta

import (
	"fmt"
	"math"
	"sync"
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func Test_trie_tree(t *testing.T) {
	tree := newTrieTree()
	assert.NotNil(t, tree)

	tree.Add("football", nil)
	tree.Add("football", nil)
	tree.Add("football", nil)
	assert.Equal(t, 1, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	tree.Add("foo", nil)
	assert.Equal(t, 2, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	tree.Add("f", 1)
	tree.Add("fo", 2)
	assert.Equal(t, 4, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	tree.Add("feet", 3)
	assert.Equal(t, 5, tree.KeyNum())
	assert.Equal(t, 11, tree.NodeNum())

	tree.Add("bike", 4)
	tree.Add("bike.bke", 5)

	tree.Add("a", 6)
	tree.Add("ab", 7)
	tree.Add("abcd", 8)
	assert.Equal(t, 10, tree.KeyNum())
	assert.Equal(t, 23, tree.NodeNum())

	tree.Add("", 323333)
	assert.Equal(t, 10, tree.KeyNum())
	assert.Equal(t, 23, tree.NodeNum())

	tree.Reset()
	assert.Zero(t, tree.KeyNum())
	assert.Zero(t, tree.NodeNum())
}

func Test_trie_MarshalBinary(t *testing.T) {
	tree := newTrieTree()
	tree.Add("hello", 9)
	tree.Add("world", 12)

	tree.Reset()
	trie := tree.(*trieTree)
	assert.Len(t, trie.nodesBuf1, 0)
	assert.Len(t, trie.nodesBuf2, 0)
	assert.Len(t, trie.root.children, 0)

	tree.Add("eleme", 1)
	tree.Add("eleme", 1)
	tree.Add("eleme", 3)
	tree.Add("eleme", 2)

	tree.Add("eleme.ci", 2)
	tree.Add("eleme.ci.etrace", 3)
	tree.Add("eleme.bdi", 4)
	tree.Add("eleme.other", 5)
	tree.Add("etrace", 6)
	tree.Add("java", 7)
	tree.Add("javascript", 8)
	tree.Add("j", 9)

	bin := tree.MarshalBinary()
	assert.NotNil(t, bin)

	assert.Equal(t, "ejltaervmaaecs.ecbcorditii.hpeettrrace", string(bin.labels)[2:])
	assert.Len(t, bin.values, 9)

	tree.Reset()
}

func Benchmark_trie_MarshalBinary(b *testing.B) {
	tree := newTrieTree()

	for i := 0; i < b.N; i++ {
		tree.Add("eleme", 1)
		tree.Add("eleme.ci", 2)
		tree.Add("eleme.ci.etrace", 3)
		tree.Add("eleme.bdi", 4)
		tree.Add("eleme.other", 5)
		tree.Add("etrace", 6)
		tree.Add("java", 7)
		tree.Add("javascript", 8)
		tree.Add("j", 9)

		tree.MarshalBinary()
		tree.Reset()
	}
}

func buildTestTrieTreeData() *trieTreeData {
	/*
		c5   e      f
		d6   l   t  i
		   e  c  r  r
		   m  d2 a  e
		   e1    c  f
		         e3 o
		            x4

		values : 5,2,1,3,4
		indexes: 0,1,2,3,4
	*/

	tree := newTrieTree()
	tree.Add("eleme", 1)   // index: 3
	tree.Add("etcd", 2)    // index: 2
	tree.Add("etrace", 3)  // index: 4
	tree.Add("firefox", 4) // index: 5
	tree.Add("c", 5)       // index: 0
	tree.Add("cd", 6)      // index: 1
	return tree.MarshalBinary()
}

func Test_trieTree_FindOffsetsByEqual(t *testing.T) {
	data := buildTestTrieTreeData()
	// test FindOffsetsByEqual
	assert.Equal(t, []int{3}, data.FindOffsetsByEqual("eleme"))
	assert.Equal(t, []int{2}, data.FindOffsetsByEqual("etcd"))
	assert.Equal(t, []int{4}, data.FindOffsetsByEqual("etrace"))
	assert.Equal(t, []int{5}, data.FindOffsetsByEqual("firefox"))
	assert.Equal(t, []int{0}, data.FindOffsetsByEqual("c"))
	assert.Equal(t, []int{1}, data.FindOffsetsByEqual("cd"))
	assert.Len(t, data.FindOffsetsByEqual("d"), 0)
	assert.Len(t, data.FindOffsetsByEqual("et"), 0)
	assert.Len(t, data.FindOffsetsByEqual("etcd1"), 0)
	assert.Len(t, data.FindOffsetsByEqual("fire"), 0)
	assert.Len(t, data.FindOffsetsByEqual("etrac"), 0)
}

func Test_trieTree_GetValuesByOffsets(t *testing.T) {
	data := buildTestTrieTreeData()
	value, ok := data.GetValueByOffset(3)
	assert.True(t, ok)
	assert.Equal(t, "eleme", value)

	_, ok = data.GetValueByOffset(8)
	assert.False(t, ok)

	// validation failure
	data.labels = append(data.labels, byte(1))
	_, ok = data.GetValueByOffset(1)
	assert.False(t, ok)
}

func Test_trieTree_walkTreeByValue(t *testing.T) {
	data := buildTestTrieTreeData()

	expects := []struct {
		prefixValue string
		exhausted   bool
		nodeNumber  uint64
	}{
		{"", true, 1},
		{"e", true, 3},
		{"z", false, 23},
		{"ellme", false, 23},
		{"elome", false, 23},
		{"elemee", false, 23},
	}
	for _, testCase := range expects {
		exhausted, nodeNumber := data.walkTreeByValue([]byte(testCase.prefixValue))
		assert.Equal(t, testCase.exhausted, exhausted)
		assert.Equal(t, testCase.nodeNumber, nodeNumber)
	}
}

func Test_trieTree_FindOffsetsByIn(t *testing.T) {
	data := buildTestTrieTreeData()
	// test FindOffsetsByIn
	assert.Len(t, data.FindOffsetsByIn([]string{"d", "c"}), 1)
	assert.Equal(t, []int{0}, data.FindOffsetsByIn([]string{"d", "c"}))
	assert.Equal(t, []int{3, 2}, data.FindOffsetsByIn([]string{"eleme", "etcd"}))
	assert.Equal(t, []int{4, 5}, data.FindOffsetsByIn([]string{"etrace", "etrace1", "firefox"}))
}

func Test_trieTree_FindOffsetsByLike(t *testing.T) {
	data := buildTestTrieTreeData()
	// test FindOffsetsByLike
	assert.Equal(t, []int{0}, data.FindOffsetsByLike("c"))
	assert.Equal(t, []int{1}, data.FindOffsetsByLike("cd"))
	assert.Len(t, data.FindOffsetsByLike("et"), 0)
	assert.Len(t, data.FindOffsetsByLike("fire"), 0)
	assert.Nil(t, data.FindOffsetsByLike(""))
	assert.Len(t, data.FindOffsetsByLike("*"), 6)
	assert.Nil(t, data.FindOffsetsByLike("etrace1"))

	assert.Equal(t, []int{2, 1}, data.FindOffsetsByLike("*cd"))
	assert.Equal(t, []int{4, 2}, data.FindOffsetsByLike("et*"))
	assert.Equal(t, []int{4, 2}, data.FindOffsetsByLike("*t*"))
}

func Test_trieTree_FindOffsetsByRegex(t *testing.T) {
	data := buildTestTrieTreeData()
	// test FindOffsetsByRegex
	assert.Len(t, data.FindOffsetsByRegex("et"), 2)
	assert.Len(t, data.FindOffsetsByRegex("cd"), 1)
	assert.Len(t, data.FindOffsetsByRegex("^c[a-d]?"), 2)
	// bad pattern
	assert.Nil(t, data.FindOffsetsByRegex("[a^-#]("))
}

func Test_trieTree_PrefixSearch(t *testing.T) {
	data := buildTestTrieTreeData()
	// test PrefixSearch
	assert.Len(t, data.PrefixSearch("e", 3), 3)
	assert.Len(t, data.PrefixSearch("e", 1), 1)
	assert.Len(t, data.PrefixSearch("etcd1", 1), 0)
}

func Test_trieTree_Iterator(t *testing.T) {
	data := buildTestTrieTreeData()
	// test iterator with prefix
	itr1 := data.Iterator("e")
	assert.True(t, itr1.HasNext())
	value, offset := itr1.Next()
	assert.Equal(t, "etrace", string(value))
	assert.Equal(t, 4, offset)

	assert.True(t, itr1.HasNext())
	value, offset = itr1.Next()
	assert.Equal(t, "etcd", string(value))
	assert.Equal(t, 2, offset)

	assert.True(t, itr1.HasNext())
	value, offset = itr1.Next()
	assert.Equal(t, "eleme", string(value))
	assert.Equal(t, 3, offset)
	assert.False(t, itr1.HasNext())

	// test iterator with no-prefix
	itr2 := data.Iterator("")
	var count = 0
	for itr2.HasNext() {
		count++
	}
	assert.Equal(t, 6, count)

	// has Error
	itr3 := data.Iterator("not-exist")
	assert.False(t, itr3.HasNext())
	assert.False(t, itr3.HasNext())
}

var (
	once4TestTrieTree sync.Once
	testTrieTree      *trieTreeData
)

func prepareTrieTreeData() *trieTreeData {
	once4TestTrieTree.Do(
		func() {
			tree := newTrieTree()
			for x := 0; x < math.MaxUint8; x++ {
				for y := 0; y < math.MaxUint8; y++ {
					// build ip
					seriesID := uint32(x*math.MaxUint8 + y)
					ip := fmt.Sprintf("192.168.%d.%d", x, y)
					r := roaring.New()
					r.Add(seriesID)
					tree.Add(ip, r)
				}
			}
			testTrieTree = tree.MarshalBinary()
		})
	return testTrieTree
}

func BenchmarkTrieTree_LikeSearch(b *testing.B) {
	data := prepareTrieTreeData()
	for i := 0; i < b.N; i++ {
		data.FindOffsetsByLike("192.168.1.1")
	}
}

func BenchmarkTrieTree_EqualSearch(b *testing.B) {
	data := prepareTrieTreeData()
	for i := 0; i < b.N; i++ {
		data.FindOffsetsByEqual("192.168.1.1")
	}
}

func BenchmarkTrieTree_InSearch(b *testing.B) {
	data := prepareTrieTreeData()
	for i := 0; i < b.N; i++ {
		data.FindOffsetsByIn([]string{"192.168.1.1", "192.168.3.2", "192.168.2.2"})
	}
}

func BenchmarkTrieTree_RegexSearch(b *testing.B) {
	data := prepareTrieTreeData()
	for i := 0; i < b.N; i++ {
		data.FindOffsetsByRegex("192\\.168")
	}
}

func BenchmarkTrieTree_PrefixSearch(b *testing.B) {
	data := prepareTrieTreeData()
	for i := 0; i < b.N; i++ {
		data.PrefixSearch("192.168", 200000)
	}
}

func BenchmarkTrieTree_Iterator(b *testing.B) {
	data := prepareTrieTreeData()

	for i := 0; i < b.N; i++ {
		itr := data.Iterator("192.168")
		for itr.HasNext() {
			itr.Next()
		}
	}
}
