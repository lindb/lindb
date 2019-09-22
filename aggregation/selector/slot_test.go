package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexSlotSelector_IndexOf(t *testing.T) {
	selector := NewIndexSlotSelector(10, 120, 1)
	assert.Equal(t, 110, selector.PointCount())
	start, end := selector.Range()
	assert.Equal(t, 10, start)
	assert.Equal(t, 120, end)
	idx, completed := selector.IndexOf(5)
	assert.Equal(t, -1, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(10)
	assert.Equal(t, 0, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(15)
	assert.Equal(t, 5, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(105)
	assert.Equal(t, 95, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(121)
	assert.Equal(t, -1, idx)
	assert.True(t, completed)

	selector = NewIndexSlotSelector(10, 130, 3)
	idx, completed = selector.IndexOf(12)
	assert.Equal(t, 0, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(13)
	assert.Equal(t, 1, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(130)
	assert.Equal(t, 40, idx)
	assert.False(t, completed)
	idx, completed = selector.IndexOf(131)
	assert.Equal(t, -1, idx)
	assert.True(t, completed)
}
