package types

import "fmt"

func GetAccurateType(lhs, rhs DataType) DataType {
	fmt.Printf("get accurate type lhs %v, rhs %v\n", lhs, rhs)
	switch lhs {
	case DataTypeInt:
		switch rhs {
		case DataTypeInt:
			return DataTypeInt
		case DataTypeFloat:
			return DataTypeFloat
		case DataTypeTimeSeries:
			return DataTypeTimeSeries
		}
	case DataTypeFloat:
		switch rhs {
		case DataTypeFloat:
			return DataTypeFloat
		case DataTypeTimeSeries:
			return DataTypeTimeSeries
		}
	case DataTypeTimeSeries:
		return DataTypeTimeSeries
	case DataTypeSum, DataTypeFirst, DataTypeLast, DataTypeMin, DataTypeMax, DataTypeHistogram:
		return DataTypeTimeSeries
	}
	// TODO: add unknown type and string
	return DataTypeFloat
}
