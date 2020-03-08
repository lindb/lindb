package trie

import (
	"bytes"
	"io"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newHostNameTrie() SuccinctTrie {
	host := []string{
		"nj11",
		"nj-2",
		"nj-3",
		"sh-4",
		"sh-5",
		"sh-6000",
		"bj-777",
		"b",
		"abcdef",
		"abcdefg",
		"bj-9"}
	var keys [][]byte
	var values [][]byte
	for _, key := range host {
		keys = append(keys, []byte(key))
		values = append(values, []byte{1})
	}
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	tree := NewBuilder().Build(keys, values, 1)
	return tree
}

func TestTrie_Get_Iterator(t *testing.T) {
	tree := newHostNameTrie()
	expects := []struct {
		key string
		ok  bool
	}{
		{"nj1", false},
		{"nj11", true},
		{"nj111", false},
		{"nj-2", true},
		{"sh-4", true},
		{"sh-5999", false},
		{"sh-6000", true},
		{"sh-600000000", false},
		{"bj-77", false},
		{"bj-777", true},
		{"", false},
		{"b", true},
		{"a", false},
		{"c", false},
		{"abcde", false},
		{"abcdef", true},
		{"Abcdef", false},
		{"abcdeF", false},
		{"abcdefg", true},
	}

	itr := tree.NewIterator()
	for _, expect := range expects {
		_, ok := tree.Get([]byte(expect.key))
		assert.Equalf(t, expect.ok, ok, "key: %s", expect.key)

		itr.Seek([]byte(expect.key))
		if itr.Valid() && expect.ok {
			assert.Equalf(t, expect.key, string(itr.Key()), "key: %s", expect.key)
		}
	}
}

func TestIterator(t *testing.T) {
	tree := newHostNameTrie()
	itr := tree.NewIterator()

	itr.SeekToFirst()
	assert.Equal(t, []byte("abcdef"), itr.Key())
	itr.Next()
	assert.Equal(t, []byte("abcdefg"), itr.Key())
}

func TestIterator_Seek(t *testing.T) {
	tree := newHostNameTrie()
	itr := tree.NewIterator()

	expects := []struct {
		input   string
		outputs []string
	}{
		{"nj", []string{"nj-2", "nj-3", "nj11"}},
		{"nj1", []string{"nj11"}},
		{"nj-", []string{"nj-2", "nj-3"}},
		{"nj-3", []string{"nj-3"}},
		{"nj-33", nil},
		{"A", nil},
		{"a", []string{"abcdef", "abcdefg"}},
		{"c", nil},
		{"abcdef", []string{"abcdef", "abcdefg"}},
		{"abcdeG", nil},
		{"abcdefg", []string{"abcdefg"}},
	}

	for _, expect := range expects {
		itr.Seek([]byte(expect.input))
		var keys []string
		for itr.Valid() {
			key := itr.Key()
			if !bytes.HasPrefix(key, []byte(expect.input)) {
				break
			}
			keys = append(keys, string(key))
			itr.Next()
		}
		assert.Equalf(t, expect.outputs, keys, "seek:%", expect.input)
	}

	itr.tree.height = 0
	itr.Seek(nil)
}

func TestPrefixIterator(t *testing.T) {
	tree := newHostNameTrie()
	getKeys := func(prefix []byte) []string {
		var keys []string
		itr := tree.NewPrefixIterator(prefix)
		for itr.Valid() {
			keys = append(keys, string(itr.Key()))
			_ = itr.Value()
			itr.Next()
		}
		itr.Next()
		assert.False(t, itr.Valid())
		return keys
	}
	assert.Len(t, getKeys(nil), 11)
	assert.Len(t, getKeys([]byte("b")), 3)
	assert.Len(t, getKeys([]byte("bj")), 2)
	assert.Len(t, getKeys([]byte("n")), 3)
	assert.Len(t, getKeys([]byte("abcde")), 2)
}

func TestBitVector_String(t *testing.T) {
	_ = hasBMI2
	var bv bitVector
	bv.Init([][]uint64{{1, 2}, {3}}, []uint32{2, 2})
	assert.Equal(t, "1011", bv.String())
}

func Test_Select64(t *testing.T) {
	assert.Equal(t, select64(0xffffff, 2), select64Broadword(0xffffff, 2))
}

func Test_popcountBlock(t *testing.T) {
	assert.Zero(t, popcountBlock([]uint64{0xfff, 0xffff}, 10, 0))
}

func Test_ByteSlice_IntSlice_Convert(t *testing.T) {
	assert.Nil(t, u64SliceToBytes(nil))
	assert.Nil(t, bytesToU64Slice(nil))
	assert.Nil(t, u32SliceToBytes(nil))
	assert.Nil(t, bytesToU32Slice(nil))
}

type mockWriter struct {
	round      int
	errorRound int
}

func (mw *mockWriter) Write(p []byte) (n int, err error) {
	defer func() {
		mw.round++
	}()
	if mw.round == mw.errorRound {
		return 0, io.ErrClosedPipe
	}
	return 0, nil
}

func TestTrie_WriteTo(t *testing.T) {
	tree := newHostNameTrie()
	for i := 0; i < 32; i++ {
		assert.Error(t, tree.WriteTo(&mockWriter{errorRound: i}))
	}
}

func TestLabelVector_Unmarshal(t *testing.T) {
	var v labelVector
	_, err := v.Unmarshal(nil)
	assert.Error(t, err)

	_, err = v.Unmarshal([]byte{1, 1, 1, 1})
	assert.Error(t, err)
}

func TestTrie_UnmarshalBinary_WithError(t *testing.T) {
	goodData, _ := newHostNameTrie().MarshalBinary()
	tree := NewTrie()
	treeImpl := tree.(*trie)

	makeCorruptData := func(left []byte) []byte {
		idx := len(goodData) - len(left)
		dst := make([]byte, len(goodData))
		copy(dst, goodData)
		dst[idx] = 0xff
		dst[idx+1] = 0xff
		dst[idx+2] = 0xff
		return dst
	}

	// empty tree
	assert.Error(t, tree.UnmarshalBinary(nil))
	// label vector unmarshal failure
	assert.Error(t, tree.UnmarshalBinary(makeCorruptData(goodData[4:])))
	buf1, _ := treeImpl.labelVec.Unmarshal(goodData[4:])
	// hasChildVec unmarshal failure
	assert.Error(t, tree.UnmarshalBinary(makeCorruptData(buf1)))
	buf1, _ = treeImpl.hasChildVec.Unmarshal(buf1)
	// loudsVec unmarshal failure
	assert.Error(t, tree.UnmarshalBinary(makeCorruptData(buf1)))
	buf1, _ = treeImpl.loudsVec.Unmarshal(buf1)
	// prefixVec unmarshal failure
	assert.Error(t, tree.UnmarshalBinary(makeCorruptData(buf1)))
	buf1, _ = treeImpl.prefixVec.Unmarshal(buf1)
	// suffix unmarshal failure
	assert.Error(t, tree.UnmarshalBinary(makeCorruptData(buf1)))

}
