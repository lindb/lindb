package types

import "fmt"

func GetAccurateType(lhs, rhs DataType) DataType {
	fmt.Printf("get accurate type lhs %v, rhs %v\n", lhs, rhs)
	switch lhs {
	case DTInt:
		switch rhs {
		case DTInt:
			return DTInt
		case DTFloat:
			return DTFloat
		case DTString:
			return DTString
		case DTTimeSeries:
			return DTTimeSeries
		case DTTimestamp:
			return DTTimestamp
		}
	case DTFloat:
		switch rhs {
		case DTFloat:
			return DTFloat
		case DTString:
			return DTString
		case DTTimeSeries:
			return DTTimeSeries
		case DTTimestamp:
			return DTTimestamp
		}
	case DTString:
		return DTString
	case DTTimeSeries:
		return DTTimeSeries
	case DTTimestamp:
		return DTTimestamp
	}
	// TODO: add unknown type and string
	return DTUnknown
}
