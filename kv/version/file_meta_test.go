package version

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
)

func TestFileMeta(t *testing.T) {
	f := NewFileMeta(10, 2, 40, 1024)
	assert.Equal(t, table.FileNumber(10), f.GetFileNumber())
	assert.Equal(t, uint32(2), f.GetMinKey())
	assert.Equal(t, uint32(40), f.GetMaxKey())
	assert.Equal(t, int32(1024), f.GetFileSize())
}
