package field

import (
	"testing"

	"github.com/lindb/lindb/aggregation/function"

	"github.com/stretchr/testify/assert"
)

func Test_Sum_getPrimitiveFields(t *testing.T) {
	assert.True(t, newSumSchema().getPrimitiveFields(function.Sum)[uint16(1)] == Sum)
	assert.Equal(t, 1, len(newSumSchema().getPrimitiveFields(function.Sum)))

	assert.True(t, newSumSchema().getDefaultPrimitiveFields()[uint16(1)] == Sum)
	assert.Equal(t, sumAggregator, newSumSchema().GetAggFunc(uint16(1)))
	assert.Equal(t, 1, len(newSumSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newSumSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Min_getPrimitiveFields(t *testing.T) {
	assert.True(t, newMinSchema().getPrimitiveFields(function.Min)[uint16(1)] == Min)
	assert.Equal(t, 1, len(newMinSchema().getPrimitiveFields(function.Min)))

	assert.True(t, newMinSchema().getDefaultPrimitiveFields()[uint16(1)] == Min)
	assert.Equal(t, minAggregator, newMinSchema().GetAggFunc(uint16(1)))
	assert.Equal(t, 1, len(newMinSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newMinSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Max_getPrimitiveFields(t *testing.T) {
	assert.True(t, newMaxSchema().getPrimitiveFields(function.Max)[uint16(1)] == Max)
	assert.Equal(t, 1, len(newMaxSchema().getPrimitiveFields(function.Max)))

	assert.True(t, newMaxSchema().getDefaultPrimitiveFields()[uint16(1)] == Max)
	assert.Equal(t, maxAggregator, newMaxSchema().GetAggFunc(uint16(1)))
	assert.Equal(t, 1, len(newMaxSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newMaxSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Gauge_getPrimitiveFields(t *testing.T) {
	assert.True(t, newGaugeSchema().getPrimitiveFields(function.Replace)[uint16(1)] == Replace)
	assert.Equal(t, 1, len(newGaugeSchema().getPrimitiveFields(function.Replace)))

	assert.True(t, newGaugeSchema().getDefaultPrimitiveFields()[uint16(1)] == Replace)
	assert.Equal(t, replaceAggregator, newGaugeSchema().GetAggFunc(uint16(1)))
	assert.Equal(t, 1, len(newGaugeSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newGaugeSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Summary_getPrimitiveFields(t *testing.T) {
	assert.True(t, newSummarySchema().getDefaultPrimitiveFields()[uint16(2)] == Sum)
	assert.Equal(t, 1, len(newSummarySchema().getDefaultPrimitiveFields()))

	assert.Equal(t, sumAggregator, newSummarySchema().GetAggFunc(uint16(1)))
	assert.Equal(t, sumAggregator, newSummarySchema().GetAggFunc(uint16(2)))
	assert.Equal(t, maxAggregator, newSummarySchema().GetAggFunc(uint16(3)))
	assert.Equal(t, minAggregator, newSummarySchema().GetAggFunc(uint16(4)))
	assert.Equal(t, replaceAggregator, newSummarySchema().GetAggFunc(uint16(5)))

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
