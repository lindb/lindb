package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AggregatorUnit(t *testing.T) {
	fun0 := ValueOf("max")
	fun1 := ValueOf("avg")
	unit := NewAggregatorUnit("a", &fun0, &fun1)
	assert.Equal(t, "a", unit.GetField())
	assert.Equal(t, "max", unit.GetDownSampling().String())
	assert.Equal(t, "avg", unit.GetAggregator().String())

	fun0 = ValueOf("d_max")
	fun1 = ValueOf("d_avg")
	unit = NewAggregatorUnit("a", &fun0, &fun1)
	assert.Equal(t, "a", unit.GetField())
	assert.Equal(t, "max", unit.GetDownSampling().String())
	assert.Equal(t, "avg", unit.GetAggregator().String())

	fun1 = ValueOf("d_avg")
	unit = NewAggregatorUnit("a", nil, &fun1)
	assert.Equal(t, "a", unit.GetField())
	assert.True(t, true, unit.GetDownSampling() == nil)
	assert.Equal(t, "avg", unit.GetAggregator().String())

	unit = NewAggregatorUnit("a", nil, nil)
	assert.Equal(t, "a", unit.GetField())
	assert.True(t, true, unit.GetDownSampling() == nil)
	assert.True(t, true, unit.GetAggregator() == nil)
}
