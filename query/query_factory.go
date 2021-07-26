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
	"context"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

type queryFactory struct {
	replicaStateMachine  broker.ReplicaStatusStateMachine
	nodeStateMachine     discovery.ActiveNodeStateMachine
	databaseStateMachine broker.DatabaseStateMachine
	taskManager          TaskManager
}

func NewQueryFactory(
	replicaStateMachine broker.ReplicaStatusStateMachine,
	nodeStateMachine discovery.ActiveNodeStateMachine,
	databaseStateMachine broker.DatabaseStateMachine,
	taskManager TaskManager,
) Factory {
	return &queryFactory{
		replicaStateMachine:  replicaStateMachine,
		nodeStateMachine:     nodeStateMachine,
		databaseStateMachine: databaseStateMachine,
		taskManager:          taskManager,
	}
}

func (qh *queryFactory) NewMetricQuery(
	ctx context.Context,
	databaseName string,
	sql string,
) MetricQuery {
	return newMetricQuery(ctx, databaseName, sql, qh)
}

func (qh *queryFactory) NewMetadataQuery(
	ctx context.Context,
	database string,
	stmt *stmt.Metadata,
) MetaDataQuery {
	return newMetadataQuery(ctx, database, stmt, qh)
}

func (qh *queryFactory) newStorageMetricQuery(
	queryFlow flow.StorageQueryFlow,
	database tsdb.Database,
	storageExecuteCtx StorageExecuteContext,
) storageMetricQuery {
	return newStorageMetricQuery(queryFlow, database, storageExecuteCtx)
}

func (qh *queryFactory) newStorageMetadataQuery(
	database tsdb.Database,
	shardIDs []int32,
	request *stmt.Metadata,
) storageMetadataQuery {
	return newStorageMetadataQuery(database, shardIDs, request)
}

// newStorageQueryContext creates the storage execute context in storage side
func (qh *queryFactory) newStorageQueryContext(
	shardIDs []int32,
	query *stmt.Query,
) StorageExecuteContext {
	return newStorageExecuteContext(shardIDs, query)
}
