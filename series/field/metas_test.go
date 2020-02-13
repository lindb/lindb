package field

import (
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Metas(t *testing.T) {
	var metas = Metas{}
	ids := make(map[uint16]struct{})
	for i := uint16(0); i < 230; i++ {
		ids[i] = struct{}{}
	}

	for i := range ids {
		metas = metas.Insert(Meta{ID: ID(i), Type: SumField, Name: strconv.Itoa(int(i))})
	}
	sort.Sort(metas)

	// GetFromName
	m, ok := metas.GetFromName("203")
	assert.True(t, ok)
	assert.Equal(t, ID(203), m.ID)
	clone := metas.Clone()
	m, ok = clone.GetFromName("204")
	assert.True(t, ok)
	assert.Equal(t, ID(204), m.ID)
	_, ok = metas.GetFromName("250")
	assert.False(t, ok)

	// GetFromID
	m, ok = metas.GetFromID(ID(204))
	assert.True(t, ok)
	assert.Equal(t, ID(204), m.ID)
	_, ok = metas.GetFromID(ID(230))
	assert.False(t, ok)

	// Intersects
	ml, ok := metas.Intersects([]ID{1, 203, 250})
	assert.False(t, ok)
	assert.Len(t, ml, 2)
	ml, ok = metas.Intersects([]ID{1, 203, 204})
	assert.True(t, ok)
	assert.Len(t, ml, 3)
}
