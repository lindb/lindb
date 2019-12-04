package point

import (
	"testing"

	"github.com/lindb/lindb/series/field"

	"github.com/stretchr/testify/assert"
)

func Test_FieldIterator(t *testing.T) {
	var itr = new(FieldIterator)
	type testField struct {
		fieldName  string
		fieldValue interface{}
		fieldType  field.Type
	}
	except := func(data []byte, fields []testField) {
		itr.Reset(data)
		for _, f := range fields {
			assert.True(t, itr.Next())
			assert.Equal(t, f.fieldName, string(itr.Name()))
			assert.Equal(t, f.fieldType, itr.Type())
			v, _ := itr.Float64Value()
			assert.Equal(t, f.fieldValue, v)
		}
		assert.False(t, itr.Next())
	}
	// normal cases
	except(
		[]byte("n=0.5,timerCount_SUM=32,timerSum_SUM=120,timerMax_MAX=180,timerMin_MIN=10"),
		[]testField{
			{fieldName: "n", fieldValue: float64(0.5), fieldType: field.Unknown},
			{fieldName: "timerCount", fieldValue: float64(32), fieldType: field.SumField},
			{fieldName: "timerSum", fieldValue: float64(120), fieldType: field.SumField},
			{fieldName: "timerMax", fieldValue: float64(180), fieldType: field.MaxField},
			{fieldName: "timerMin", fieldValue: float64(10), fieldType: field.MinField},
		},
	)
	// escape cases
	except(
		[]byte("x\\,_SMY=20.222,h\\ g_HGM=32"),
		[]testField{
			{fieldName: "x,", fieldValue: float64(20.222), fieldType: field.SummaryField},
			{fieldName: "h g", fieldValue: float64(32), fieldType: field.HistogramField},
		},
	)
	// parse failure cases
	except(
		[]byte("a_SUM=20c,b_SUM=3 2"),
		[]testField{
			{fieldName: "a", fieldValue: float64(0), fieldType: field.SumField},
			{fieldName: "b", fieldValue: float64(0), fieldType: field.SumField},
		},
	)
}

func Test_FieldIterator_Int64Value(t *testing.T) {
	var itr = new(FieldIterator)
	itr.valueBuf = []byte("1232.32")
	v, err := itr.Int64Value()
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), v)

	itr.valueBuf = []byte("1232")
	v, err = itr.Int64Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1232), v)
}
