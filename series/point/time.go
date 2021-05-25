// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
