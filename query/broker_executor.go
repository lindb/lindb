package query

import (
	"context"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/series"
)

// brokerExecutor represents the broker query executor,
// 1) chooses the storage nodes that the data is relatively complete
// 2) chooses broker nodes for root and intermediate computing from all available broker nodes
// 3) storage node as leaf computing node does filtering and atomic compute
// 4) intermediate computing nodes are optional, only need if has group by query, does order by for grouping
// 4) root computing node does function and expression computing ???? //TODO  need?
// 5) finally returns result set to user  ???? //TODO  need?
//
// NOTICE: there are some scenarios:
// 1) some assignment shards not in query replica shards,
//    maybe some expectant results are lost in data in offline shard, WHY can query not completely data,
//    because of for the system availability.
type brokerExecutor struct {
	database string
	sql      string

	replicaStateMachine replica.StatusStateMachine
	nodeStateMachine    broker.NodeStateMachine

	resultSet chan *series.TimeSeriesEvent

	jobManager parallel.JobManager

	ctx context.Context
	err error
}

// newBrokerExecutor creates the execution which executes the job of parallel query
func newBrokerExecutor(ctx context.Context, database string, sql string,
	replicaStateMachine replica.StatusStateMachine, nodeStateMachine broker.NodeStateMachine,
	jobManager parallel.JobManager) parallel.Executor {
	exec := &brokerExecutor{
		sql:                 sql,
		database:            database,
		replicaStateMachine: replicaStateMachine,
		nodeStateMachine:    nodeStateMachine,
		jobManager:          jobManager,
		ctx:                 ctx,
	}
	return exec
}

// Execute executes search logic in broker level,
// 1) get metadata based on params
// 2) build execute plan
// 3) run distribution query job
func (e *brokerExecutor) Execute() <-chan *series.TimeSeriesEvent {
	//FIXME need using storage's replica state ???
	storageNodes := e.replicaStateMachine.GetQueryableReplicas(e.database)
	if len(storageNodes) == 0 {
		e.err = errNoAvailableStorageNode
		return nil
	}

	brokerNodes := e.nodeStateMachine.GetActiveNodes()
	plan := newBrokerPlan(e.sql, storageNodes, e.nodeStateMachine.GetCurrentNode(), brokerNodes)
	if err := plan.Plan(); err != nil {
		e.err = err
		return nil
	}
	brokerPlan := plan.(*brokerPlan)
	brokerPlan.physicalPlan.Database = e.database
	e.resultSet = make(chan *series.TimeSeriesEvent)
	if err := e.jobManager.SubmitJob(parallel.NewJobContext(e.resultSet, brokerPlan.physicalPlan, brokerPlan.query)); err != nil {
		e.err = err
		close(e.resultSet)
		return nil
	}
	return e.resultSet
}

// Error returns the execution error
func (e *brokerExecutor) Error() error {
	return e.err
}
