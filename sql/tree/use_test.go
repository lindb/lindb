package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseStatement(t *testing.T) {
	stmt, err := GetParser().CreateStatement("use test", NewNodeIDAllocator())
	assert.NoError(t, err)
	assert.Equal(t, &Use{
		Database: &Identifier{
			Value: "test",
		},
	}, stmt)
}
