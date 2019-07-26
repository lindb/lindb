package table

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/eleme/lindb/pkg/fileutil"

	"github.com/stretchr/testify/assert"
)

const (
	testKVPath = "test_builder"
)

func Test_magicNumber(t *testing.T) {
	code := []byte("eleme-ci")
	assert.Len(t, code, 8)
	assert.Equal(t, magicNumberOffsetFile, binary.BigEndian.Uint64(code))
}

func Test_BuildStore(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
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
