package util

// AggregatorUnit does for data aggregate
type AggregatorUnit struct {
	field        string
	downSampling *FunctionType
	aggregator   *FunctionType
}

// NewAggregatorUnit build default aggregator unit
func NewAggregatorUnit(field string, downSampling *FunctionType, aggregator *FunctionType) *AggregatorUnit {
	return &AggregatorUnit{
		field:        field,
		downSampling: downSampling,
		aggregator:   aggregator,
	}
}

// GetField get aggregator unit field message
func (a *AggregatorUnit) GetField() string {
	return a.field
}

// GetDownSampling get aggregator unit down sampling message
func (a *AggregatorUnit) GetDownSampling() *FunctionType {
	return a.downSampling
}

// GetAggregator get aggregator unit aggregator message
func (a *AggregatorUnit) GetAggregator() *FunctionType {
	return a.aggregator
}
