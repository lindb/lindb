package indextbl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trie_tree(t *testing.T) {
	tree := newTrieTree()
	assert.NotNil(t, tree)

	tree.Add("football", 1)
	tree.Add("football", 11)
	tree.Add("football", 121)
	assert.Equal(t, 1, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	tree.Add("foo", 2)
	assert.Equal(t, 2, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	tree.Add("f", 345)
	tree.Add("fo", 45)
	assert.Equal(t, 4, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	tree.Add("feet", 45)
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

	assert.Equal(t, "ejltaervmaaecs.ecbcorditii.hpeettrrace", string(bin.labels)[1:])
	assert.Equal(t, []uint32{9, 7, 1, 6, 2, 4, 8, 5, 3}, bin.values)

}

func Benchmark_MarshalBinary(b *testing.B) {
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
