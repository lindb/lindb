package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
)

func Test_Wrong_Type_Set(t *testing.T) {
	bs := newBlockStore(30)
	b1 := bs.allocIntBlock()
	b1.setFloatValue(0, 10.0)
	assert.Equal(t, 0.0, b1.getFloatValue(0))

	b2 := bs.allocFloatBlock()
	b2.setIntValue(0, 20)
	assert.Equal(t, int64(0), b2.getIntValue(0))
}

func Test_allocBlock(t *testing.T) {
	bs := newBlockStore(30)
	assert.NotNil(t, bs.allocBlock(field.Integer))
	assert.NotNil(t, bs.allocBlock(field.Float))
	assert.Nil(t, bs.allocBlock(0))
}
