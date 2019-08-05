package indextbl

import (
	"testing"

	"github.com/RoaringBitmap/roaring"

	"github.com/stretchr/testify/assert"
)

func Test_trie_tree(t *testing.T) {
	tree := newTrieTree()
	assert.NotNil(t, tree)

	bitmap1 := roaring.New()
	bitmap1.AddRange(1, 100)
	tree.Add("football", 1, bitmap1)
	tree.Add("football", 11, bitmap1)
	tree.Add("football", 121, bitmap1)
	assert.Equal(t, 1, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	bitmap2 := roaring.New()
	bitmap2.AddRange(323, 400)

	tree.Add("foo", 2, bitmap2)
	assert.Equal(t, 2, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	tree.Add("f", 344, roaring.New())
	tree.Add("fo", 45, roaring.New())
	assert.Equal(t, 4, tree.KeyNum())
	assert.Equal(t, 8, tree.NodeNum())

	tree.Add("feet", 45, roaring.New())
	assert.Equal(t, 5, tree.KeyNum())
	assert.Equal(t, 11, tree.NodeNum())

	tree.Add("bike", 4, roaring.New())
	tree.Add("bike.bke", 5, roaring.New())

	tree.Add("a", 6, roaring.New())
	tree.Add("ab", 7, roaring.New())
	tree.Add("abcd", 8, roaring.New())
	assert.Equal(t, 10, tree.KeyNum())
	assert.Equal(t, 23, tree.NodeNum())

	tree.Add("", 323333, roaring.New())
	assert.Equal(t, 10, tree.KeyNum())
	assert.Equal(t, 23, tree.NodeNum())

	tree.Reset()
	assert.Zero(t, tree.KeyNum())
	assert.Zero(t, tree.NodeNum())
}

func Test_trie_MarshalBinary(t *testing.T) {
	tree := newTrieTree()
	tree.Add("hello", 9, roaring.New())
	tree.Add("world", 12, roaring.New())

	tree.Reset()
	trie := tree.(*trieTree)
	assert.Len(t, trie.nodesBuf1, 0)
	assert.Len(t, trie.nodesBuf2, 0)
	assert.Len(t, trie.root.children, 0)

	tree.Add("eleme", 1, roaring.New())
	tree.Add("eleme", 1, roaring.New())
	tree.Add("eleme", 3, roaring.New())
	tree.Add("eleme", 2, roaring.New())

	tree.Add("eleme.ci", 2, roaring.New())
	tree.Add("eleme.ci.etrace", 3, roaring.New())
	tree.Add("eleme.bdi", 4, roaring.New())
	tree.Add("eleme.other", 5, roaring.New())
	tree.Add("etrace", 6, roaring.New())
	tree.Add("java", 7, roaring.New())
	tree.Add("javascript", 8, roaring.New())
	tree.Add("j", 9, roaring.New())

	bin := tree.MarshalBinary()
	assert.NotNil(t, bin)

	assert.Equal(t, "ejltaervmaaecs.ecbcorditii.hpeettrrace", string(bin.labels)[1:])
	assert.Len(t, bin.values, 9)

	tree.Reset()
}

func Benchmark_MarshalBinary(b *testing.B) {
	tree := newTrieTree()
	rb := roaring.New()

	for i := 0; i < b.N; i++ {
		tree.Add("eleme", 1, rb)
		tree.Add("eleme.ci", 2, rb)
		tree.Add("eleme.ci.etrace", 3, rb)
		tree.Add("eleme.bdi", 4, rb)
		tree.Add("eleme.other", 5, rb)
		tree.Add("etrace", 6, rb)
		tree.Add("java", 7, rb)
		tree.Add("javascript", 8, rb)
		tree.Add("j", 9, rb)

		tree.MarshalBinary()
		tree.Reset()
	}
}
