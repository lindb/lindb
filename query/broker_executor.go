package query

import (
	"context"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/sql/stmt"
)

// brokerExecutor implements parallel.BrokerExecutor
type brokerExecutor struct {
	database  string
	namespace string
	sql       string
	query     *stmt.Query

	replicaStateMachine replica.StatusStateMachine
	nodeStateMachine    broker.NodeStateMachine

	jobManager parallel.JobManager

	ctx context.Context

	executeCtx parallel.BrokerExecuteContext
}

// newBrokerExecutor creates the execution which executes the job of parallel query
func newBrokerExecutor(ctx context.Context, database string, namespace string, sql string,
	replicaStateMachine replica.StatusStateMachine, nodeStateMachine broker.NodeStateMachine,
	jobManager parallel.JobManager) parallel.BrokerExecutor {
	exec := &brokerExecutor{
		sql:                 sql,
		database:            database,
		namespace:           namespace,
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
func (e *brokerExecutor) Execute() {
	//FIXME need using storage's replica state ???
	storageNodes := e.replicaStateMachine.GetQueryableReplicas(e.database)
	brokerNodes := e.nodeStateMachine.GetActiveNodes()
	plan := newBrokerPlan(e.sql, storageNodes, e.nodeStateMachine.GetCurrentNode(), brokerNodes)
	var err error
	if len(storageNodes) == 0 {
		err = errNoAvailableStorageNode
	} else {
		err = plan.Plan()
	}

	// maybe plan doesn't execute(query statement is nil), because storage nodes is empty
	brokerPlan := plan.(*brokerPlan)
	e.executeCtx = parallel.NewBrokerExecuteContext(brokerPlan.query)

	if err != nil {
		e.executeCtx.Complete(err)
		return
	}

	brokerPlan.physicalPlan.Database = e.database
	brokerPlan.physicalPlan.Namespace = e.namespace
	e.query = brokerPlan.query

	if err := e.jobManager.SubmitJob(parallel.NewJobContext(e.ctx,
		e.executeCtx.ResultCh(), brokerPlan.physicalPlan, e.query),
	); err != nil {
		e.executeCtx.Complete(err)
		return
	}
}

func (e *brokerExecutor) ExecuteContext() parallel.BrokerExecuteContext {
	return e.executeCtx
}
