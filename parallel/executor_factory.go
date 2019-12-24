package parallel

import (
	"context"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
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
		shardIDs []int32,
		query *stmt.Query,
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
		replicaStateMachine replica.StatusStateMachine,
		nodeStateMachine broker.NodeStateMachine,
		jobManager JobManager,
	) BrokerExecutor

	// NewMetadataBrokerExecutor creates the metadata executor in broker side
	NewMetadataBrokerExecutor(
		ctx context.Context,
		databaseName string,
		request *stmt.Metadata,
		replicaStateMachine replica.StatusStateMachine,
		nodeStateMachine broker.NodeStateMachine,
		jobManager JobManager,
	) MetadataExecutor
}
