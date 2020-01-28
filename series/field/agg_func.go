package field

import "math"

var (
	sumAggregator     = sumAgg{aggType: Sum}
	countAggregator   = sumAgg{aggType: Count}
	minAggregator     = minAgg{aggType: Min}
	maxAggregator     = maxAgg{aggType: Max}
	replaceAggregator = replaceAgg{aggType: Replace}
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
	case Replace:
		return replaceAggregator
	default:
		return nil
	}
}

// AggFunc represents field's aggregator function for int64 or float64 value
type AggFunc interface {
	// Aggregate aggregates two float64 values into one
	Aggregate(a, b float64) float64
	// AggType return aggregator type
	AggType() AggType
}

// sumAgg represents sum aggregator
type sumAgg struct {
	aggType AggType
}

func (s sumAgg) AggType() AggType               { return s.aggType }
func (s sumAgg) Aggregate(a, b float64) float64 { return a + b }

// minAgg represents min aggregator
type minAgg struct {
	aggType AggType
}

func (m minAgg) AggType() AggType               { return m.aggType }
func (m minAgg) Aggregate(a, b float64) float64 { return math.Min(a, b) }

// maxAgg represents max aggregator
type maxAgg struct {
	aggType AggType
}

func (m maxAgg) AggType() AggType               { return m.aggType }
func (m maxAgg) Aggregate(a, b float64) float64 { return math.Max(a, b) }

// replaceAgg represents replace aggregator
type replaceAgg struct {
	aggType AggType
}

func (m replaceAgg) AggType() AggType               { return m.aggType }
func (m replaceAgg) Aggregate(a, b float64) float64 { return b }
