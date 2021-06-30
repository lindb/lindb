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

package parallel

import (
	"context"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./executor_factory.go -destination=./executor_factory_mock.go -package=parallel

// ExecutorFactory represents the executor factory that creates storage/broker executor
type ExecutorFactory interface {
	// NewStorageExecutor creates the storage executor based on params
	NewStorageExecutor(
		queryFlow flow.StorageQueryFlow,
		database tsdb.Database,
		storageExecuteCtx StorageExecuteContext,
	) Executor

	// NewMetadataStorageExecutor creates the metadata executor in storage side
	NewMetadataStorageExecutor(
		database tsdb.Database,
		shardIDs []int32,
		request *stmt.Metadata,
	) MetadataExecutor

	// NewBrokerExecutor creates the broker executor based on params
	NewBrokerExecutor(
		ctx context.Context,
		databaseName string,
		sql string,
		replicaStateMachine broker.ReplicaStatusStateMachine,
		nodeStateMachine discovery.ActiveNodeStateMachine,
		databaseStateMachine broker.DatabaseStateMachine,
		jobManager JobManager,
	) BrokerExecutor

	// NewMetadataBrokerExecutor creates the metadata executor in broker side
	NewMetadataBrokerExecutor(
		ctx context.Context,
		databaseName string,
		request *stmt.Metadata,
		replicaStateMachine broker.ReplicaStatusStateMachine,
		nodeStateMachine discovery.ActiveNodeStateMachine,
		jobManager JobManager,
	) MetadataExecutor

	// NewStorageExecuteContext creates the storage execute context in storage side
	NewStorageExecuteContext(shardIDs []int32, query *stmt.Query) StorageExecuteContext
}
