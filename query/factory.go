package query

import (
	"context"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

// executorFactory implements parallel.ExecutorFactory
type executorFactory struct{}

// NewExecutorFactory creates executor factory
func NewExecutorFactory() parallel.ExecutorFactory {
	return &executorFactory{}
}

// NewStorageExecutor creates storage executor
func (*executorFactory) NewStorageExecutor(
	queryFlow flow.StorageQueryFlow,
	database tsdb.Database,
	storageExecuteCtx parallel.StorageExecuteContext,
) parallel.Executor {
	return newStorageExecutor(queryFlow, database, storageExecuteCtx)
}

// NewMetadataStorageExecutor creates the metadata executor in storage side
func (*executorFactory) NewMetadataStorageExecutor(
	database tsdb.Database,
	shardIDs []int32,
	request *stmt.Metadata,
) parallel.MetadataExecutor {
	return newMetadataStorageExecutor(database, shardIDs, request)
}

// NewStorageExecutor creates broker executor
func (*executorFactory) NewBrokerExecutor(
	ctx context.Context,
	databaseName string,
	namespace string,
	sql string,
	replicaStateMachine replica.StatusStateMachine,
	nodeStateMachine broker.NodeStateMachine,
	databaseStateMachine database.DBStateMachine,
	jobManager parallel.JobManager,
) parallel.BrokerExecutor {
	return newBrokerExecutor(ctx, databaseName, namespace, sql,
		replicaStateMachine, nodeStateMachine, databaseStateMachine,
		jobManager)
}

// NewMetadataBrokerExecutor creates the metadata executor in broker side
func (*executorFactory) NewMetadataBrokerExecutor(
	ctx context.Context,
	databaseName string,
	request *stmt.Metadata,
	replicaStateMachine replica.StatusStateMachine,
	nodeStateMachine broker.NodeStateMachine,
	jobManager parallel.JobManager,
) parallel.MetadataExecutor {
	return newMetadataBrokerExecutor(ctx, databaseName, request, nodeStateMachine, replicaStateMachine, jobManager)
}

// NewStorageExecuteContext creates the storage execute context in storage side
func (*executorFactory) NewStorageExecuteContext(namespace string, shardIDs []int32, query *stmt.Query) parallel.StorageExecuteContext {
	return newStorageExecuteContext(namespace, shardIDs, query)
}
