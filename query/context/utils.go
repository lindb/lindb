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

package context

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
)

// calcTimeRangeAndInterval calculates the query time range and interval based on input params and database config.
func calcTimeRangeAndInterval(statement *stmt.Query, cfg models.Database) {
	option := cfg.Option
	interval := statement.Interval
	if interval <= 0 {
		// if query interval not set, first set it using the smallest interval in storage option.
		interval = option.Intervals[0].Interval
	}
	// re-calc query interval based on query time range
	interval = timeutil.CalcQueryInterval(statement.TimeRange, interval)
	storageInterval := option.FindMatchSmallestInterval(interval)
	intervalVal := storageInterval.Int64()
	statement.TimeRange.Start = timeutil.Truncate(statement.TimeRange.Start, intervalVal)
	statement.TimeRange.End = timeutil.Truncate(statement.TimeRange.End, intervalVal)
	if statement.AutoGroupByTime {
		// fill group by interval if not set
		statement.Interval = timeutil.Interval(statement.TimeRange.End-statement.TimeRange.Start) + storageInterval
	}
	// if auto calc interval < user input, need to use use input
	if interval < statement.Interval {
		interval = statement.Interval
	}
	intervalRatio := timeutil.CalIntervalRatio(interval.Int64(), storageInterval.Int64())
	// truncate query interval
	interval = timeutil.Interval(storageInterval.Int64() * int64(intervalRatio))

	statement.StorageInterval = storageInterval
	statement.Interval = interval
	statement.IntervalRatio = intervalRatio
}
