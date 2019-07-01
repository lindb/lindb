package aggregation

import "github.com/eleme/lindb/pkg/field"

type PrimitiveAggregator interface {
	ValueType() field.ValueType
	Aggregate(idx int, value interface{})
	Values() interface{}
}

type aggregator interface {
	aggregateInt(idx int, value int64)
	aggregateFloat(idx int, value float64)
	values() interface{}
}

type primitiveAggregator struct {
	agg aggregator

	valueType  field.ValueType
	pointCount int
}

func NewPrimitiveAggregator(valueType field.ValueType,
	aggFunc field.AggFunc,
	pointCount int) PrimitiveAggregator {
	var agg aggregator
	switch valueType {
	case field.Integer:
		agg = newIntegerAggregator(pointCount, aggFunc)
	case field.Float:
		agg = newFloatAggregator(pointCount, aggFunc)
	}

	return &primitiveAggregator{
		agg:        agg,
		valueType:  valueType,
		pointCount: pointCount,
	}
}

func (agg *primitiveAggregator) ValueType() field.ValueType {
	return agg.valueType
}

func (agg *primitiveAggregator) Aggregate(idx int, value interface{}) {
	if idx < 0 || idx >= agg.pointCount {
		return
	}
	switch val := value.(type) {
	case int64:
		agg.agg.aggregateInt(idx, val)
	case int:
		agg.agg.aggregateInt(idx, int64(val))
	case float64:
		agg.agg.aggregateFloat(idx, val)
	}

}
func (agg *primitiveAggregator) Values() interface{} {
	return agg.agg.values()
}

type noopAggregator struct {
}

func (n *noopAggregator) aggregateInt(idx int, value int64) {

}

func (n *noopAggregator) aggregateFloat(idx int, value float64) {

}
func (n *noopAggregator) values() interface{} {
	return nil
}

type integerAggregator struct {
	noopAggregator
	aggFunc field.AggFunc
	vals    []int64
}

func newIntegerAggregator(pointCount int, aggFunc field.AggFunc) aggregator {
	return &integerAggregator{
		vals:    make([]int64, pointCount),
		aggFunc: aggFunc,
	}
}

func (agg *integerAggregator) aggregateInt(idx int, value int64) {
	agg.vals[idx] = agg.aggFunc.AggregateInt(agg.vals[idx], value)
}
func (agg *integerAggregator) values() interface{} {
	return agg.vals
}

type floatAggregator struct {
	noopAggregator
	aggFunc field.AggFunc
	vals    []float64
}

func newFloatAggregator(pointCount int, aggFunc field.AggFunc) aggregator {
	return &floatAggregator{
		vals:    make([]float64, pointCount),
		aggFunc: aggFunc,
	}
}

func (agg *floatAggregator) aggregateFloat(idx int, value float64) {
	agg.vals[idx] = agg.aggFunc.AggregateFloat(agg.vals[idx], value)
}
func (agg *floatAggregator) values() interface{} {
	return agg.vals
}
