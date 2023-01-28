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

import (
	"github.com/lindb/lindb/internal/linmetric"
)

// QueryStatistics represents query statistics.
type QueryStatistics struct {
	CreatedTasks *linmetric.BoundCounter // create query task
	ExpireTasks  *linmetric.BoundCounter // task expire, long-term no response
	AliveTask    *linmetric.BoundGauge   // current executing task(alive)
	EmitResponse *linmetric.BoundCounter // emit response to parent node
	OmitResponse *linmetric.BoundCounter // omit response because task evicted
}

// TransportStatistics represents request/response transport statistics.
type TransportStatistics struct {
	SentRequest          *linmetric.BoundCounter // send request success
	SentRequestFailures  *linmetric.BoundCounter // send request failure
	SentResponses        *linmetric.BoundCounter // send response to parent success
	SentResponseFailures *linmetric.BoundCounter // send response failure
}

// StorageQueryStatistics represents storage query statistics.
type StorageQueryStatistics struct {
	MetricQuery         *linmetric.BoundCounter // execute metric query success(just plan it)
	MetricQueryFailures *linmetric.BoundCounter // execute metric query failure
	MetaQuery           *linmetric.BoundCounter // metadata query success
	MetaQueryFailures   *linmetric.BoundCounter // metadata query failure
	OmitRequest         *linmetric.BoundCounter // omit request(task no belong to current node, wrong stream etc.)
}

// NewTransportStatistics creates a transport statistics.
func NewTransportStatistics(registry *linmetric.Registry) *TransportStatistics {
	scope := registry.NewScope("lindb.task.transport")
	return &TransportStatistics{
		SentRequest:          scope.NewCounter("sent_requests"),
		SentResponses:        scope.NewCounter("sent_responses"),
		SentResponseFailures: scope.NewCounter("sent_responses_failures"),
		SentRequestFailures:  scope.NewCounter("sent_requests_failures"),
	}
}

// NewQueryStatistics creates a query statistics.
func NewQueryStatistics(registry *linmetric.Registry) *QueryStatistics {
	scope := registry.NewScope("lindb.query")
	return &QueryStatistics{
		CreatedTasks: scope.NewCounter("created_tasks"),
		AliveTask:    scope.NewGauge("alive_tasks"),
		ExpireTasks:  scope.NewCounter("expire_tasks"),
		EmitResponse: scope.NewCounter("emitted_responses"),
		OmitResponse: scope.NewCounter("omitted_responses"),
	}
}

// NewStorageQueryStatistics creates a storage query statistics.
func NewStorageQueryStatistics() *StorageQueryStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.storage.query")
	return &StorageQueryStatistics{
		MetricQuery:         scope.NewCounter("metric_queries"),
		MetricQueryFailures: scope.NewCounter("metric_query_failures"),
		MetaQuery:           scope.NewCounter("meta_queries"),
		MetaQueryFailures:   scope.NewCounter("meta_query_failures"),
		OmitRequest:         scope.NewCounter("omitted_requests"),
	}
}
