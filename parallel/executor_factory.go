package parallel

import (
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./executor_factory.go -destination=./executor_factory_mock.go -package=parallel

// ExecutorFactory represents the executor factory that creates storage/broker executor
type ExecutorFactory interface {
	// NewStorageExecutor creates the storage executor based on params
	NewStorageExecutor(engine tsdb.Engine, shardIDs []int32, query *stmt.Query) Executor
	// NewBrokerExecutor creates the broker executor based on params
	NewBrokerExecutor(database string, sql string,
		replicaStateMachine replica.StatusStateMachine, nodeStateMachine broker.NodeStateMachine,
		jobManager JobManager) Executor
}
