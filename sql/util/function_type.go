package util

type FunctionType int32

const (
	Sum       = "sum"
	Count     = "count"
	Min       = "min"
	Max       = "max"
	Avg       = "avg"
	Mean      = "mean"
	Histogram = "histogram"
)

const (
	SUM FunctionType = iota + 1
	COUNT
	MIN
	MAX
	AVG
	MEAN
	HISTOGRAM
)

// String override FunctionType to string method,default `sum`
func (f FunctionType) String() string {
	switch f {
	case SUM:
		return Sum
	case COUNT:
		return Count
	case MIN:
		return Min
	case MAX:
		return Max
	case AVG:
		return Avg
	case MEAN:
		return Mean
	case HISTOGRAM:
		return Histogram
	default:
		return Sum
	}
}

// GetFunctionType get FunctionType by Function name,default `SUM`
func GetFunctionType(name string) FunctionType {
	switch name {
	case Sum:
		return FunctionType(1)
	case Count:
		return FunctionType(2)
	case Min:
		return FunctionType(3)
	case Max:
		return FunctionType(4)
	case Avg:
		return FunctionType(5)
	case Mean:
		return FunctionType(6)
	case Histogram:
		return FunctionType(7)
	default:
		return FunctionType(1)
	}
}
