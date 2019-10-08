package function

// FuncType is the definition of function type
type FuncType int

const (
	Sum FuncType = iota + 1
	Min
	Max
	Count
	Avg
	Histogram
	Stddev

	Unknown
)

// String return the function's name
func (t FuncType) String() string {
	switch t {
	case Sum:
		return "sum"
	case Min:
		return "min"
	case Max:
		return "max"
	case Count:
		return "count"
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
