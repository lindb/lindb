package aggregation

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/sql/stmt"
)

func TestBinary_eval(t *testing.T) {
	assert.Equal(t, float64(10), eval(stmt.ADD, 4, 6))
	assert.Equal(t, float64(-2), eval(stmt.SUB, 4, 6))
	assert.Equal(t, float64(24), eval(stmt.MUL, 4, 6))
	assert.Equal(t, 0.5, eval(stmt.DIV, 4, 8))
	assert.Equal(t, float64(0), eval(stmt.DIV, 4, 0))

	// wrong binary operator
	assert.Equal(t, float64(0), eval(stmt.OR, 4, 8))
}

func TestBinaryEval(t *testing.T) {
	assert.Nil(t, binaryEval(stmt.DIV, collections.NewFloatArray(10), collections.NewFloatArray(10)))

	fa := collections.NewFloatArray(10)
	fa.SetValue(0, 1.1)
	fa.SetValue(5, 5.5)
	fa.SetValue(8, 9.9)
	result := binaryEval(stmt.ADD, collections.NewFloatArray(10), fa)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, fa, result)

	result = binaryEval(stmt.SUB, collections.NewFloatArray(10), fa)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, -1.1, result.GetValue(0))
	assert.Equal(t, -5.5, result.GetValue(5))
	assert.Equal(t, -9.9, result.GetValue(8))
	result = binaryEval(stmt.MUL, collections.NewFloatArray(10), fa)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 0.0, result.GetValue(0))
	assert.Equal(t, 0.0, result.GetValue(5))
	assert.Equal(t, 0.0, result.GetValue(8))
	result = binaryEval(stmt.MUL, fa, collections.NewFloatArray(10))
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 0.0, result.GetValue(0))
	assert.Equal(t, 0.0, result.GetValue(5))
	assert.Equal(t, 0.0, result.GetValue(8))
	result = binaryEval(stmt.SUB, fa, collections.NewFloatArray(10))
	assert.Equal(t, fa, result)
	result = binaryEval(stmt.ADD, fa, collections.NewFloatArray(10))
	assert.Equal(t, fa, result)

	fa = collections.NewFloatArray(10)
	fa.SetValue(0, 1.1)
	fa.SetValue(5, 5.5)
	fa2 := collections.NewFloatArray(10)
	fa2.SetValue(0, 1.1)
	fa2.SetValue(8, 9.9)
	result = binaryEval(stmt.ADD, fa, fa2)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 2.2, result.GetValue(0))
	assert.Equal(t, 5.5, result.GetValue(5))
	assert.Equal(t, 9.9, result.GetValue(8))
	result = binaryEval(stmt.SUB, fa, fa2)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 0.0, result.GetValue(0))
	assert.Equal(t, 5.5, result.GetValue(5))
	assert.Equal(t, -9.9, result.GetValue(8))
	result = binaryEval(stmt.MUL, fa, fa2)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 1.21, math.Floor(result.GetValue(0)*100)/100)
	assert.Equal(t, 0.0, result.GetValue(5))
	assert.Equal(t, 0.0, result.GetValue(8))
	result = binaryEval(stmt.DIV, fa, fa2)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 1.0, result.GetValue(0))
	assert.Equal(t, 0.0, result.GetValue(5))
	assert.Equal(t, 0.0, result.GetValue(8))
}
