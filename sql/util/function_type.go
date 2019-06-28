package util

type FunctionType int32

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
func (this FunctionType) String() string {
	switch this {
	case SUM:
		return "sum"
	case COUNT:
		return "count"
	case MIN:
		return "min"
	case MAX:
		return "max"
	case AVG:
		return "avg"
	case MEAN:
		return "mean"
	case HISTOGRAM:
		return "histogram"
	default:
		return "sum"
	}
}

// GetFunctionType get FunctionType by Function name,default `SUM`
func GetFunctionType(name string) FunctionType {
	switch name {
	case "sum":
		return FunctionType(1)
	case "count":
		return FunctionType(2)
	case "min":
		return FunctionType(3)
	case "max":
		return FunctionType(4)
	case "avg":
		return FunctionType(5)
	case "mean":
		return FunctionType(6)
	case "histogram":
		return FunctionType(7)
	default:
		return FunctionType(1)
	}
}
