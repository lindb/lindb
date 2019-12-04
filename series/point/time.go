package point

import (
	"math"
)

const (
	minNanoTime   = int64(math.MinInt64)
	maxNanoTime   = int64(math.MaxInt64)
	minMicroTime  = minNanoTime / 1000
	maxMicroTime  = maxNanoTime / 1000
	minMilliTime  = minMicroTime / 1000
	maxMilliTime  = maxMicroTime / 1000
	minSecondTime = minMilliTime / 1000
	maxSecondTime = maxMilliTime / 1000
)

// MilliSecondOf calculates the given time, and converts it to milliseconds.
func MilliSecondOf(timestamp int64) int64 {
	switch {
	case minSecondTime <= timestamp && timestamp <= maxSecondTime:
		return timestamp * 1000
	case minMilliTime <= timestamp && timestamp <= maxMilliTime:
		return timestamp
	case minMicroTime <= timestamp && timestamp <= maxMicroTime:
		return timestamp / 1000
	default:
		return timestamp / 1000 / 1000
	}
}
