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

	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
)

// brokerExecutor implements parallel.BrokerExecutor
type brokerExecutor struct {
	database string
	sql      string
	query    *stmt.Query

	replicaStateMachine  broker.ReplicaStatusStateMachine
	nodeStateMachine     discovery.ActiveNodeStateMachine
	databaseStateMachine broker.DatabaseStateMachine

	jobManager parallel.JobManager

	ctx context.Context

	executeCtx parallel.BrokerExecuteContext
}

// newBrokerExecutor creates the execution which executes the job of parallel query
func newBrokerExecutor(ctx context.Context, database string, sql string,
	replicaStateMachine broker.ReplicaStatusStateMachine, nodeStateMachine discovery.ActiveNodeStateMachine,
	databaseStateMachine broker.DatabaseStateMachine,
	jobManager parallel.JobManager) parallel.BrokerExecutor {
	exec := &brokerExecutor{
		sql:                  sql,
		database:             database,
		replicaStateMachine:  replicaStateMachine,
		nodeStateMachine:     nodeStateMachine,
		databaseStateMachine: databaseStateMachine,
		jobManager:           jobManager,
		ctx:                  ctx,
	}
	return exec
}

// Execute executes search logic in broker level,
// 1) get metadata based on params
// 2) build execute plan
// 3) run distribution query job
func (e *brokerExecutor) Execute() {
	startTime := timeutil.NowNano()

	databaseCfg, ok := e.databaseStateMachine.GetDatabaseCfg(e.database)
	if !ok {
		e.executeCtx = parallel.NewBrokerExecuteContext(startTime, nil)
		e.executeCtx.Complete(errDatabaseNotExist)
		return
	}

	//FIXME need using storage's replica state ???
	storageNodes := e.replicaStateMachine.GetQueryableReplicas(e.database)
	brokerNodes := e.nodeStateMachine.GetActiveNodes()
	plan := newBrokerPlan(e.sql, databaseCfg, storageNodes, e.nodeStateMachine.GetCurrentNode(), brokerNodes)

	var err error
	if len(storageNodes) == 0 {
		err = errNoAvailableStorageNode
	} else {
		err = plan.Plan()
	}

	// maybe plan doesn't execute(query statement is nil), because storage nodes is empty
	brokerPlan := plan.(*brokerPlan)
	e.executeCtx = parallel.NewBrokerExecuteContext(startTime, brokerPlan.query)

	if err != nil {
		e.executeCtx.Complete(err)
		return
	}

	brokerPlan.physicalPlan.Database = e.database
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
