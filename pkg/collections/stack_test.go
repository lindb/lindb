package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack()
	assert.True(t, s.Empty())

	s.Push(1)
	assert.False(t, s.Empty())
	assert.Equal(t, 1, s.len)

	assert.Equal(t, 1, s.Peek().(int))

	assert.Equal(t, 1, s.Pop().(int))
	assert.True(t, s.Empty())

	s.Push(1)
	s.Push(2)

	assert.Equal(t, 2, s.len)
	assert.Equal(t, 2, s.Pop().(int))
	assert.Equal(t, 1, s.len)

	s.Pop()

	assert.True(t, s.Empty())
	assert.Nil(t, s.Peek())
	assert.Nil(t, s.Pop())
}
