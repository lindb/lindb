package query

import (
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

type executorFactory struct {
}

func NewExecutorFactory() parallel.ExecutorFactory {
	return &executorFactory{}
}

func (*executorFactory) NewStorageExecutor(engine tsdb.Engine, shardIDs []int32, query *stmt.Query) parallel.Executor {
	return newStorageExecutor(engine, shardIDs, query)
}

func (*executorFactory) NewBrokerExecutor(database string, sql string,
	replicaStateMachine replica.StatusStateMachine, nodeStateMachine broker.NodeStateMachine,
	jobManager parallel.JobManager) parallel.Executor {
	return newBrokerExecutor(database, sql, replicaStateMachine, nodeStateMachine, jobManager)
}
