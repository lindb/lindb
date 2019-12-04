package field

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Metas(t *testing.T) {
	var metas = Metas{}
	ids := make(map[uint16]struct{})
	for i := uint16(0); i < 1000; i++ {
		ids[i] = struct{}{}
	}

	for i := range ids {
		metas = metas.Insert(Meta{ID: i, Type: SumField, Name: strconv.Itoa(int(i))})
	}

	// GetFromName
	m, ok := metas.GetFromName("304")
	assert.True(t, ok)
	assert.Equal(t, uint16(304), m.ID)
	clone := metas.Clone()
	m, ok = clone.GetFromName("304")
	assert.True(t, ok)
	assert.Equal(t, uint16(304), m.ID)
	_, ok = metas.GetFromName("1001")
	assert.False(t, ok)

	// GetFromID
	m, ok = metas.GetFromID(304)
	assert.True(t, ok)
	assert.Equal(t, uint16(304), m.ID)
	_, ok = metas.GetFromID(1001)
	assert.False(t, ok)

	// Intersects
	ml, ok := metas.Intersects([]uint16{1, 303, 1001})
	assert.False(t, ok)
	assert.Len(t, ml, 2)
	ml, ok = metas.Intersects([]uint16{1, 303, 304})
	assert.True(t, ok)
	assert.Len(t, ml, 3)
}
