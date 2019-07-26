package function

// Type is the definition of function type
type Type int

const (
	Sum Type = iota + 1
	Min
	Max
	Avg
	Histogram
	Stddev

	Unknown
)
