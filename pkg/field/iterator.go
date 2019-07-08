package field

type Iterator interface {
	Next() bool
	ValueType() ValueType
	AggType() AggType
	PrimitiveFieldID() uint8
	IntValue() int64
	FloatValue() float64
}
