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
