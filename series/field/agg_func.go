package field

import "math"

var (
	sumAggregator   = sumAgg{aggType: Sum}
	countAggregator = sumAgg{aggType: Count}
	minAggregator   = minAgg{aggType: Min}
	maxAggregator   = maxAgg{aggType: Max}
)

// AggFunc returns aggregator function by given func type
func (t AggType) AggFunc() AggFunc {
	switch t {
	case Sum:
		return sumAggregator
	case Count:
		return countAggregator
	case Min:
		return minAggregator
	case Max:
		return maxAggregator
	default:
		return nil
	}
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

func (s sumAgg) AggType() AggType                    { return s.aggType }
func (s sumAgg) AggregateInt(a, b int64) int64       { return a + b }
func (s sumAgg) AggregateFloat(a, b float64) float64 { return a + b }

// minAgg represents min aggregator
type minAgg struct {
	aggType AggType
}

func (m minAgg) AggType() AggType { return m.aggType }
func (m minAgg) AggregateInt(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func (m minAgg) AggregateFloat(a, b float64) float64 { return math.Min(a, b) }

// maxAgg represents max aggregator
type maxAgg struct {
	aggType AggType
}

func (m maxAgg) AggType() AggType { return m.aggType }
func (m maxAgg) AggregateInt(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func (m maxAgg) AggregateFloat(a, b float64) float64 { return math.Max(a, b) }
