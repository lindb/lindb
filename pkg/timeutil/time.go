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

package timeutil

// Truncate truncates timestamp based on interval
func Truncate(timestamp, interval int64) int64 {
	return timestamp / interval * interval
}

// CalPointCount calculates point counts between start time and end time by interval
func CalPointCount(startTime, endTime, interval int64) int {
	diff := endTime - startTime
	pointCount := diff / interval
	if diff%interval > 0 {
		pointCount++
	}
	if pointCount == 0 {
		pointCount = 1
	}
	return int(pointCount)
}

// CalIntervalRatio calculates the interval ratio for query,
// if query interval < storage interval return 1.
func CalIntervalRatio(queryInterval, storageInterval int64) int {
	if storageInterval == 0 || queryInterval < storageInterval {
		return 1
	}
	return int(queryInterval / storageInterval)
}
