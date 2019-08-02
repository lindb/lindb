package table

import (
	"os"
	"testing"

	"github.com/lindb/lindb/pkg/fileutil"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	var builder, err = NewStoreBuilder(testKVPath, 10)
	defer os.RemoveAll(testKVPath)

	_ = builder.Add(1, []byte("test"))
	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	err = builder.Close()
	assert.Nil(t, err)

	cache := NewCache(testKVPath)

	var reader, err2 = cache.GetReader("", 10)
	if err2 != nil {
		t.Error(err2)
	}
	defer reader.Close()

	assert.Equal(t, []byte("test"), reader.Get(1))
	assert.Equal(t, []byte("test10"), reader.Get(10))
	cache.Close()
}

func TestStoreIterator(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	var builder, err = NewStoreBuilder(testKVPath, 10)
	defer os.RemoveAll(testKVPath)

	_ = builder.Add(1, []byte("test"))
	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	err = builder.Close()
	assert.Nil(t, err)

	cache := NewCache(testKVPath)
	var reader, err2 = cache.GetReader("", 10)
	if err2 != nil {
		t.Error(err2)
	}
	defer reader.Close()
	it := reader.Iterator()
	assert.True(t, it.Next())
	assert.Equal(t, uint32(1), it.Key())
	assert.Equal(t, []byte("test"), it.Value())

	assert.True(t, it.Next())
	assert.Equal(t, uint32(10), it.Key())
	assert.Equal(t, []byte("test10"), it.Value())

	assert.False(t, it.Next())
}
