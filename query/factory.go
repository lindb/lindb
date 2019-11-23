package query

import (
	"context"

	"github.com/lindb/lindb/coordinator/broker"
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
	shardIDs []int32,
	query *stmt.Query,
) parallel.Executor {
	return newStorageExecutor(queryFlow, database, shardIDs, query)
}

// NewStorageExecutor creates broker executor
func (*executorFactory) NewBrokerExecutor(
	ctx context.Context,
	databaseName string,
	sql string,
	replicaStateMachine replica.StatusStateMachine,
	nodeStateMachine broker.NodeStateMachine,
	jobManager parallel.JobManager,
) parallel.BrokerExecutor {
	return newBrokerExecutor(ctx, databaseName, sql, replicaStateMachine, nodeStateMachine, jobManager)
}
