package tblstore

import (
	"testing"

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

func Test_trieTreeQuerier(t *testing.T) {
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

	data := tree.MarshalBinary()
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

	// test FindOffsetsByIn
	assert.Len(t, data.FindOffsetsByIn([]string{"d", "c"}), 1)
	assert.Equal(t, []int{0}, data.FindOffsetsByIn([]string{"d", "c"}))
	assert.Equal(t, []int{3, 2}, data.FindOffsetsByIn([]string{"eleme", "etcd"}))
	assert.Equal(t, []int{4, 5}, data.FindOffsetsByIn([]string{"etrace", "etrace1", "firefox"}))

	// test FindOffsetsByLike
	assert.Equal(t, []int{0, 1}, data.FindOffsetsByLike("c"))
	assert.Equal(t, []int{1}, data.FindOffsetsByLike("cd"))
	assert.Equal(t, []int{2, 4}, data.FindOffsetsByLike("et"))
	assert.Equal(t, []int{5}, data.FindOffsetsByLike("fire"))
	assert.Nil(t, data.FindOffsetsByLike(""))
	assert.Nil(t, data.FindOffsetsByLike("etrace1"))
}
