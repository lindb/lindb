package series

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_Append(t *testing.T) {
	block := NewBlock(1, 1)
	assert.False(t, block.Append(0, 10.0))
	assert.False(t, block.Append(1, 10.0))
	assert.True(t, block.Append(2, 10.0))
}
