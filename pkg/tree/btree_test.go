package tree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	tree := New(BytesCompare)

	for i := 0; i < 10000; i++ {
		tree.Set([]byte(fmt.Sprintf("%s%d", "key-", i)), i)
		tree.Put([]byte(fmt.Sprintf("%s%d", "key-", i)), func(oldV interface{}, exists bool) (newV interface{}, write bool) {
			return oldV, exists
		})
	}
	for i := 0; i < 10000; i++ {
		v, _ := tree.Get([]byte(fmt.Sprintf("%s%d", "key-", i)))
		assert.Equal(t, i, v.(int))
	}
	k, v := tree.First()
	assert.Equal(t, []byte("key-0"), k)
	assert.Equal(t, 0, v)

	k, v = tree.Last()
	assert.Equal(t, []byte("key-9999"), k)
	assert.Equal(t, 9999, v)

	assert.Equal(t, 10000, tree.Len())

	it, ok := tree.Seek([]byte("key-1"))
	if ok {
		var count = 0
		for {
			_, _, err := it.Next()
			//_, _, _ = it.Prev()
			if nil == err {
				count++
			} else {
				break
			}
		}
		it.Close()
		assert.Equal(t, 9999, count)
	}

	it, err := tree.SeekFirst()
	if nil == err {
		var count = 0
		for {
			_, _, err := it.Next()
			if nil == err {
				count++
			} else {
				break
			}
		}
		it.Close()
		assert.Equal(t, 10000, count)
	}

	it, err = tree.SeekLast()
	if nil == err {
		var count = 0
		for {
			_, _, err := it.Prev()
			if nil == err {
				count++
			} else {
				break
			}
		}
		it.Close()
		assert.Equal(t, 10000, count)
	}

	ok = tree.Delete([]byte("key-0"))
	assert.Equal(t, true, ok)

	tree.Close()
}
