package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortUint32(t *testing.T) {
	var uints = []uint32{3, 4, 1, 9, 5, 6, 2, 8, 7, 0}
	sortArray := SortUint32(uints)
	assert.Equal(t, sortArray, []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}
