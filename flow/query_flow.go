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

package flow

import (
	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./query_flow.go -destination=./query_flow_mock.go -package=flow

// StorageQueryFlow represents the storage query engine execute flow.
type StorageQueryFlow interface {
	// Prepare prepares the query flow, builds the flow execute context based on group aggregator specs.
	Prepare(
		interval timeutil.Interval,
		intervalRatio int,
		timeRange timeutil.TimeRange,
		aggregatorSpecs aggregation.AggregatorSpecs,
	)
	// Filtering does the filtering task.
	Filtering(task concurrent.Task)
	// Grouping does the grouping task.
	Grouping(task concurrent.Task)
	// Load does the load task.
	Load(task concurrent.Task)
	// Reduce reduces the down sampling aggregator's result.
	Reduce(tags string, it series.GroupedIterator)
	// ReduceTagValues reduces the group by tag values.
	ReduceTagValues(tagKeyIndex int, tagValues map[uint32]string)
	// Complete completes the query flow with error.
	Complete(err error)
}

// QueryTask represents query task for data search flow.
type QueryTask interface {
	// BeforeRun invokes before task run.
	BeforeRun()
	// Run executes task query logic.
	Run() error
	// AfterRun invokes after task run.
	AfterRun()
}
