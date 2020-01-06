package function

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
)

func TestFuncCall_Unknown(t *testing.T) {
	result := FuncCall(Unknown, collections.NewFloatArray(10))
	assert.Nil(t, result)
}

func TestFuncCall_Sum(t *testing.T) {
	result := FuncCall(Sum, nil)
	assert.Nil(t, result)
	result = FuncCall(Sum)
	assert.Nil(t, result)

	array1 := collections.NewFloatArray(10)
	array2 := collections.NewFloatArray(20)
	result = FuncCall(Sum, array1, array2)
	assert.Equal(t, array1, result)
}

func TestFuncCall_Avg(t *testing.T) {
	result := FuncCall(Avg, nil)
	assert.Nil(t, result)
	result = FuncCall(Avg)
	assert.Nil(t, result)

	array1 := collections.NewFloatArray(10)
	array1.SetValue(1, 10.0)
	array2 := collections.NewFloatArray(20)
	array2.SetValue(1, 5.0)
	result = FuncCall(Avg, array1)
	assert.Nil(t, result)
	result = FuncCall(Avg, array1, array2)
	assert.Equal(t, 2.0, result.GetValue(1))
}
