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

package metric

import (
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

func execute(ctx *ExecutionContext, pool concurrent.Pool, handle func()) {
	pool.Submit(ctx.GetTaskContext(), concurrent.NewTask(func() {
		handle()
		ctx.CompleteTask()
	}, func(err error) {
		ctx.CompleteTask()
	}))
}

func calcTimeRangeAndInterval(
	timeRange timeutil.TimeRange, interval timeutil.Interval,
	cfg *models.DatabaseConfig,
) (timeutil.TimeRange, timeutil.Interval) {
	option := cfg.Option
	targetInterval := interval
	if targetInterval <= 0 {
		// if query interval not set, first set it using the smallest interval in storage option.
		targetInterval = option.Intervals[0].Interval
	}
	// re-calc query interval based on query time range
	targetInterval = timeutil.CalcQueryInterval(timeRange, targetInterval)
	// TODO: need add test
	storageInterval := option.FindMatchSmallestInterval(targetInterval)
	targetInterval = max(storageInterval, targetInterval)
	intervalVal := storageInterval.Int64()
	return timeutil.TimeRange{
		Start: timeutil.Truncate(timeRange.Start, intervalVal),
		End:   timeutil.Truncate(timeRange.End, intervalVal),
	}, targetInterval
}
