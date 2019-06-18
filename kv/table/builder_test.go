package table

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/util"
)

const (
	testKVPath = "test_builder"
)

func Test_BuildStore(t *testing.T) {
	_ = util.MkDirIfNotExist(testKVPath)
	var builder, err = NewStoreBuilder(testKVPath, 10)
	defer os.RemoveAll(testKVPath)
	defer builder.Close()

	assert.Nil(t, err)

	err = builder.Add(1, []byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), builder.Count())

	// reject for duplicate key
	err = builder.Add(1, []byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), builder.Count())

	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	assert.Equal(t, uint32(1), builder.MinKey())
	assert.Equal(t, uint32(10), builder.MaxKey())
	assert.Equal(t, int64(10), builder.FileNumber())
	assert.True(t, builder.Size() > 0)
}
