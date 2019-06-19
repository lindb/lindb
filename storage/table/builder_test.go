package table

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testKVFile = "test_kv.test"
)

func Test_BuildStore(t *testing.T) {
	var builder, err = NewStoreBuilder(testKVFile)
	defer os.Remove(testKVFile)
	defer builder.Close()

	assert.Nil(t, err)

	added := builder.Add(1, []byte("test"))
	assert.True(t, added)

	added = builder.Add(1, []byte("test"))
	assert.False(t, added)

	NewReader(testKVFile)
}
