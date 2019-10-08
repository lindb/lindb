package field

import (
	"testing"

	"github.com/lindb/lindb/aggregation/function"

	"github.com/stretchr/testify/assert"
)

func Test_Sum_getPrimitiveFields(t *testing.T) {
	assert.True(t, newSumSchema().getPrimitiveFields(function.Sum)[uint16(1)] == Sum)
	assert.Equal(t, 1, len(newSumSchema().getPrimitiveFields(function.Sum)))

	assert.Nil(t, newSumSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Summary_getPrimitiveFields(t *testing.T) {
	assert.Equal(t, 1, len(newSummarySchema().getPrimitiveFields(function.Sum)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Sum)[uint16(1)] == Sum)
	assert.Equal(t, 1, len(newSummarySchema().getPrimitiveFields(function.Count)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Count)[uint16(2)] == Sum)
	assert.Equal(t, 1, len(newSummarySchema().getPrimitiveFields(function.Max)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Max)[uint16(3)] == Max)
	assert.Equal(t, 1, len(newSummarySchema().getPrimitiveFields(function.Min)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Min)[uint16(4)] == Min)
	assert.Equal(t, 2, len(newSummarySchema().getPrimitiveFields(function.Avg)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Avg)[uint16(1)] == Sum)
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Avg)[uint16(2)] == Sum)

	assert.Nil(t, newSummarySchema().getPrimitiveFields(function.FuncType(128)))
}
