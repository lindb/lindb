package util

import (
	"strings"
)

/**
 * Aggregation function utils
 */

const DownSampling = "d_"

var SimpleFunction = []string{SUM.String(), COUNT.String(), MIN.String(), MAX.String()}

// ValueOf get FunctionType by function name
func ValueOf(functionName string) FunctionType {
	if strings.HasPrefix(functionName, DownSampling) {
		functionName = functionName[len(DownSampling):]
	}
	return GetFunctionType(functionName)
}

// IsDownSamplingOrAggregator judge function name is down sampling or aggregator
func IsDownSamplingOrAggregator(function string) bool {
	if strings.HasPrefix(function, DownSampling) {
		return true
	}
	for i := range SimpleFunction {
		if SimpleFunction[i] == function {
			return true
		}
	}
	return false
}
