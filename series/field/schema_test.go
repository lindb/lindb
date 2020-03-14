package field

import (
	"testing"

	"github.com/lindb/lindb/aggregation/function"

	"github.com/stretchr/testify/assert"
)

func Test_Sum_getPrimitiveFields(t *testing.T) {
	assert.True(t, newSumSchema().getPrimitiveFields(function.Sum)[0].AggType == Sum)
	assert.True(t, newSumSchema().getPrimitiveFields(function.Sum)[0].FieldID == PrimitiveID(1))
	assert.Equal(t, 1, len(newSumSchema().getPrimitiveFields(function.Sum)))
	assert.Equal(t, []PrimitiveID{1}, newSumSchema().GetAllPrimitiveFields())

	assert.True(t, newSumSchema().getDefaultPrimitiveFields()[0].AggType == Sum)
	assert.True(t, newSumSchema().getDefaultPrimitiveFields()[0].FieldID == PrimitiveID(1))
	assert.Equal(t, sumAggregator, newSumSchema().GetAggFunc(PrimitiveID(1)))
	assert.Equal(t, 1, len(newSumSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newSumSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Increase_getPrimitiveFields(t *testing.T) {
	assert.True(t, newIncreaseSchema().getPrimitiveFields(function.Sum)[0].AggType == Sum)
	assert.True(t, newIncreaseSchema().getPrimitiveFields(function.Sum)[0].FieldID == PrimitiveID(1))
	assert.Equal(t, 1, len(newIncreaseSchema().getPrimitiveFields(function.Sum)))
	assert.Equal(t, []PrimitiveID{1}, newIncreaseSchema().GetAllPrimitiveFields())

	assert.True(t, newIncreaseSchema().getDefaultPrimitiveFields()[0].AggType == Sum)
	assert.True(t, newIncreaseSchema().getDefaultPrimitiveFields()[0].FieldID == PrimitiveID(1))
	assert.Equal(t, sumAggregator, newIncreaseSchema().GetAggFunc(PrimitiveID(1)))
	assert.Equal(t, 1, len(newIncreaseSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newIncreaseSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Min_getPrimitiveFields(t *testing.T) {
	assert.True(t, newMinSchema().getPrimitiveFields(function.Min)[0].AggType == Min)
	assert.True(t, newMinSchema().getPrimitiveFields(function.Min)[0].FieldID == PrimitiveID(1))
	assert.Equal(t, 1, len(newMinSchema().getPrimitiveFields(function.Min)))
	assert.Equal(t, []PrimitiveID{1}, newMinSchema().GetAllPrimitiveFields())

	assert.True(t, newMinSchema().getDefaultPrimitiveFields()[0].AggType == Min)
	assert.True(t, newMinSchema().getDefaultPrimitiveFields()[0].FieldID == PrimitiveID(1))
	assert.Equal(t, minAggregator, newMinSchema().GetAggFunc(PrimitiveID(1)))
	assert.Equal(t, 1, len(newMinSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newMinSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Max_getPrimitiveFields(t *testing.T) {
	assert.True(t, newMaxSchema().getPrimitiveFields(function.Max)[0].AggType == Max)
	assert.True(t, newMaxSchema().getPrimitiveFields(function.Max)[0].FieldID == PrimitiveID(1))
	assert.Equal(t, 1, len(newMaxSchema().getPrimitiveFields(function.Max)))
	assert.Equal(t, []PrimitiveID{1}, newMaxSchema().GetAllPrimitiveFields())

	assert.True(t, newMaxSchema().getDefaultPrimitiveFields()[0].AggType == Max)
	assert.True(t, newMaxSchema().getDefaultPrimitiveFields()[0].FieldID == PrimitiveID(1))
	assert.Equal(t, maxAggregator, newMaxSchema().GetAggFunc(PrimitiveID(1)))
	assert.Equal(t, 1, len(newMaxSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newMaxSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Gauge_getPrimitiveFields(t *testing.T) {
	assert.True(t, newGaugeSchema().getPrimitiveFields(function.Replace)[0].AggType == Replace)
	assert.True(t, newGaugeSchema().getPrimitiveFields(function.Replace)[0].FieldID == PrimitiveID(1))
	assert.Equal(t, 1, len(newGaugeSchema().getPrimitiveFields(function.Replace)))
	assert.Equal(t, []PrimitiveID{1}, newGaugeSchema().GetAllPrimitiveFields())

	assert.True(t, newGaugeSchema().getDefaultPrimitiveFields()[0].AggType == Replace)
	assert.True(t, newGaugeSchema().getDefaultPrimitiveFields()[0].FieldID == PrimitiveID(1))
	assert.Equal(t, replaceAggregator, newGaugeSchema().GetAggFunc(PrimitiveID(1)))
	assert.Equal(t, 1, len(newGaugeSchema().getDefaultPrimitiveFields()))

	assert.Nil(t, newGaugeSchema().getPrimitiveFields(function.FuncType(128)))
}

func Test_Summary_getPrimitiveFields(t *testing.T) {
	assert.True(t, newSummarySchema().getDefaultPrimitiveFields()[0].AggType == Sum)
	assert.True(t, newSummarySchema().getDefaultPrimitiveFields()[0].FieldID == PrimitiveID(2))
	assert.Equal(t, 1, len(newSummarySchema().getDefaultPrimitiveFields()))
	assert.Equal(t, []PrimitiveID{1, 2, 3, 4, 50, 75, 90, 95, 99, 39, 49}, newSummarySchema().GetAllPrimitiveFields())

	assert.Equal(t, sumAggregator, newSummarySchema().GetAggFunc(PrimitiveID(1)))
	assert.Equal(t, sumAggregator, newSummarySchema().GetAggFunc(PrimitiveID(2)))
	assert.Equal(t, minAggregator, newSummarySchema().GetAggFunc(PrimitiveID(3)))
	assert.Equal(t, maxAggregator, newSummarySchema().GetAggFunc(PrimitiveID(4)))
	assert.Equal(t, replaceAggregator, newSummarySchema().GetAggFunc(PrimitiveID(5)))

	assert.Equal(t, 1, len(newSummarySchema().getPrimitiveFields(function.Sum)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Sum)[0].AggType == Sum)
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Sum)[0].FieldID == PrimitiveID(1))
	assert.Equal(t, 1, len(newSummarySchema().getPrimitiveFields(function.Count)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Count)[0].AggType == Sum)
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Count)[0].FieldID == PrimitiveID(2))
	assert.Equal(t, 1, len(newSummarySchema().getPrimitiveFields(function.Max)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Max)[0].AggType == Max)
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Max)[0].FieldID == PrimitiveID(4))
	assert.Equal(t, 1, len(newSummarySchema().getPrimitiveFields(function.Min)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Min)[0].AggType == Min)
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Min)[0].FieldID == PrimitiveID(3))
	assert.Equal(t, 2, len(newSummarySchema().getPrimitiveFields(function.Avg)))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Avg)[0].AggType == Sum)
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Avg)[0].FieldID == PrimitiveID(1))
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Avg)[1].AggType == Sum)
	assert.True(t, newSummarySchema().getPrimitiveFields(function.Avg)[1].FieldID == PrimitiveID(2))

	assert.Nil(t, newSummarySchema().getPrimitiveFields(function.FuncType(128)))
}
