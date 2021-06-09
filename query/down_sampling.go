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

package query

import (
	"github.com/lindb/lindb/pkg/timeutil"
)

// downSamplingTimeRange returns down sampling time range and interval ratio
func downSamplingTimeRange(queryInterval,
	storageInterval timeutil.Interval,
	queryTimeRange timeutil.TimeRange,
) (
	timeRange timeutil.TimeRange,
	intervalRatio int,
	interval timeutil.Interval,
) {
	// 1. calc interval, default use storage interval's interval if user not input
	interval = storageInterval
	intervalRatio = 1
	if queryInterval > 0 {
		intervalRatio = timeutil.CalIntervalRatio(queryInterval.Int64(), interval.Int64())
		interval = queryInterval
	}
	// 2. truncate time range
	timeRange = timeutil.TimeRange{
		Start: timeutil.Truncate(queryTimeRange.Start, interval.Int64()),
		End:   timeutil.Truncate(queryTimeRange.End, interval.Int64()),
	}
	return
}
