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

package metrics

import "github.com/lindb/lindb/internal/linmetric"

// ConcurrentStatistics represents concurrent pool statistics.
type ConcurrentStatistics struct {
	WorkersAlive       *linmetric.BoundGauge   // current workers count in use
	WorkersCreated     *linmetric.BoundCounter // workers created count since start
	WorkersKilled      *linmetric.BoundCounter // workers killed since start
	TasksConsumed      *linmetric.BoundCounter // tasks consumed count
	TasksRejected      *linmetric.BoundCounter // tasks rejected count
	TasksPanic         *linmetric.BoundCounter // tasks execute panic count
	TasksWaitingTime   *linmetric.BoundCounter // tasks waiting total time
	TasksExecutingTime *linmetric.BoundCounter // tasks executing total time with waiting period
}

// NewConcurrentStatistics creates concurrent statistics.
func NewConcurrentStatistics(poolName string, registry *linmetric.Registry) *ConcurrentStatistics {
	scope := registry.NewScope("lindb.concurrent.pool", "pool_name", poolName)
	return &ConcurrentStatistics{
		WorkersAlive:       scope.NewGauge("workers_alive"),
		WorkersCreated:     scope.NewCounter("workers_created"),
		WorkersKilled:      scope.NewCounter("workers_killed"),
		TasksConsumed:      scope.NewCounter("tasks_consumed"),
		TasksRejected:      scope.NewCounter("tasks_rejected"),
		TasksPanic:         scope.NewCounter("tasks_panic"),
		TasksWaitingTime:   scope.NewCounter("tasks_waiting_duration_sum"),
		TasksExecutingTime: scope.NewCounter("tasks_executing_duration_sum"),
	}
}

// LimitStatistics represents rate limit statistics.
type LimitStatistics struct {
	Throttles *linmetric.BoundCounter // counter reaches the max-concurrency
	Timeouts  *linmetric.BoundCounter // counter pending and then timeout
}

// NewLimitStatistics creates a rate limit statistics.
func NewLimitStatistics(limitType string, registry *linmetric.Registry) *LimitStatistics {
	scope := registry.NewScope("lindb.concurrent.limit", "type", limitType)
	return &LimitStatistics{
		Throttles: scope.NewCounter("throttle_requests"),
		Timeouts:  scope.NewCounter("timeout_requests"),
	}
}
