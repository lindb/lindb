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

package brokerquery

import (
	"context"
	"errors"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=brokerquery

// TaskType represents the distribution query task type
type TaskType int

const (
	RootTask TaskType = iota + 1
	IntermediateTask
)

var (
	ErrTimeout = errors.New("exceed timeout")
)

// Executor represents a query executor both storage/broker side.
// When returning query results the following is the order in which processing takes place:
// 1) filtering
// 2) Scanning
// 3) Grouping if need
// 4) Down sampling
// 5) Aggregation
// 6) Functions
// 7) Expressions
// Execute execute query
// 1) plan query language
// 2) aggregator data from time series(memory/file/network)

// BrokerExecutor represents the broker query executor,
// 1) chooses the storage nodes that the data is relatively complete
// 2) chooses broker nodes for root and intermediate computing from all available broker nodes
// 3) storage node as leaf computing node does filtering and atomic compute
// 4) intermediate computing nodes are optional, only need if has group by query, does order by for grouping
// 4) root computing node does function and expression computing
// 5) finally returns result set to user
//
// NOTICE: there are some scenarios:
// 1) some assignment shards not in query replica shards,
//    maybe some expectant results are lost in data in offline shard, WHY can query not completely data,
//    because of for the system availability.

type MetricQuery interface {
	WaitResponse() (*models.ResultSet, error)
}

// MetadataExecutor represents the metadata query executor, includes:
// 1. suggest metric name
// 2. suggest tag keys by spec metric name
// 3. suggest tag values by spec metric name and tag key
// 4. suggest fields by spec metric name

type MetaDataQuery interface {
	WaitResponse() ([]string, error)
}

// Factory is the handler for executing querying tasks
type Factory interface {
	NewMetricQuery(
		ctx context.Context,
		databaseName string,
		sql string,
	) MetricQuery

	NewMetadataQuery(
		ctx context.Context,
		databaseName string,
		stmt *stmt.Metadata,
	) MetaDataQuery
}
