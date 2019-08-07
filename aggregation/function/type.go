package function

// FuncType is the definition of function type
type FuncType int

const (
	Sum FuncType = iota + 1
	Min
	Max
	Avg
	Histogram
	Stddev

	Unknown
)

// FuncTypeString return the function's name
func FuncTypeString(funcType FuncType) string {
	switch funcType {
	case Sum:
		return "sum"
	case Min:
		return "min"
	case Max:
		return "max"
	case Avg:
		return "avg"
	case Histogram:
		return "histogram"
	case Stddev:
		return "stddev"
	default:
		return "unknown"
	}
}
