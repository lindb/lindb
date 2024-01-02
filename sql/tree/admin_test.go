package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdmin_FlushDatabase(t *testing.T) {
	stmt, err := GetParser().CreateStatement("flush database test", NewNodeIDAllocator())
	assert.NoError(t, err)
	assert.Equal(t, &FlushDatabase{Database: "test"}, stmt)
}

func TestAdmin_CompactDatabase(t *testing.T) {
	stmt, err := GetParser().CreateStatement("compact database test", NewNodeIDAllocator())
	assert.NoError(t, err)
	assert.Equal(t, &CompactDatabase{Database: "test"}, stmt)
}
