package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFamilySlotRange(t *testing.T) {
	sr := newFamilySlotRange(10, 20)
	sr.setSlot(15)
	start, end := sr.getSlotRange()
	assert.Equal(t, uint16(10), start)
	assert.Equal(t, uint16(20), end)
	sr.setSlot(5)
	start, end = sr.getSlotRange()
	assert.Equal(t, uint16(5), start)
	assert.Equal(t, uint16(20), end)
	sr.setSlot(27)
	start, end = sr.getSlotRange()
	assert.Equal(t, uint16(5), start)
	assert.Equal(t, uint16(27), end)
}
