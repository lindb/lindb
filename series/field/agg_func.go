package field

import "fmt"

func init() {
	registerFunc(Sum, &sumAgg{aggType: Sum})
	registerFunc(Count, &sumAgg{aggType: Count})
	registerFunc(Min, &minAgg{aggType: Min})
	registerFunc(Max, &maxAgg{aggType: Max})
}

// FuncType represents field's aggregator function type
type FuncType string

var aggFuncMap = make(map[AggType]AggFunc)

// registerFunc register aggregator function for given func type, if have duplicate func type, panic
func registerFunc(funcType AggType, aggFunc AggFunc) {
	if _, ok := aggFuncMap[funcType]; ok {
		panic(fmt.Sprintf("agg func type already registered: %d", funcType))
	}
	aggFuncMap[funcType] = aggFunc
}

//GetAggFunc returns aggregator function by given func type
func GetAggFunc(funcType AggType) AggFunc {
	return aggFuncMap[funcType]
}

// AggFunc represents field's aggregator function for int64 or float64 value
type AggFunc interface {
	// AggregateInt aggregates two int64 values into one
	AggregateInt(a, b int64) int64
	// AggregateInt aggregates two float64 values into one
	AggregateFloat(a, b float64) float64
	// AggType return aggregator type
	AggType() AggType
}

// sumAgg represents sum aggregator
type sumAgg struct {
	aggType AggType
}

func (s *sumAgg) AggType() AggType { return s.aggType }

// AggregateInt returns a+b for int64 value
func (s *sumAgg) AggregateInt(a, b int64) int64 {
	return a + b
}

// AggregateInt returns a+b for float64 value
func (s *sumAgg) AggregateFloat(a, b float64) float64 {
	return a + b
}

// minAgg represents min aggregator
type minAgg struct {
	aggType AggType
}

func (m *minAgg) AggType() AggType { return m.aggType }

// AggregateInt returns the smaller of two int64 values
func (m *minAgg) AggregateInt(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// AggregateInt returns the smaller of two float64 values
func (m *minAgg) AggregateFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// maxAgg represents max aggregator
type maxAgg struct {
	aggType AggType
}

func (m *maxAgg) AggType() AggType { return m.aggType }

// AggregateInt returns the greater of two int64 values
func (m *maxAgg) AggregateInt(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// AggregateFloat returns the greater of two float64 values
func (m *maxAgg) AggregateFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
