package query

import (
	"context"
	"net/http"
	"time"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/parallel"
)

// MetricAPI represents the metric query api
type MetricAPI struct {
	replicaStateMachine replica.StatusStateMachine
	nodeStateMachine    broker.NodeStateMachine
	executorFactory     parallel.ExecutorFactory
	jobManager          parallel.JobManager
}

// NewMetricAPI creates the metric query api
func NewMetricAPI(replicaStateMachine replica.StatusStateMachine, nodeStateMachine broker.NodeStateMachine,
	executorFactory parallel.ExecutorFactory, jobManager parallel.JobManager) *MetricAPI {
	return &MetricAPI{
		replicaStateMachine: replicaStateMachine,
		nodeStateMachine:    nodeStateMachine,
		executorFactory:     executorFactory,
		jobManager:          jobManager,
	}
}

// Search searches the metric data based on database and sql.
func (m *MetricAPI) Search(w http.ResponseWriter, r *http.Request) {
	db, err := api.GetParamsFromRequest("db", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	namespace, _ := api.GetParamsFromRequest("ns", r, constants.DefaultNamespace, false)
	sql, err := api.GetParamsFromRequest("sql", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	//FIXME add timeout cfg
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	exec := m.executorFactory.NewBrokerExecutor(ctx, db, namespace, sql, m.replicaStateMachine, m.nodeStateMachine, m.jobManager)
	exec.Execute()

	brokerExecutor := exec.(parallel.BrokerExecutor)
	exeCtx := brokerExecutor.ExecuteContext()

	//FIXME timeout logic use select
	resultCh := exeCtx.ResultCh()
	for result := range resultCh {
		exeCtx.Emit(result)
	}

	resultSet, err := exeCtx.ResultSet()
	if err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, resultSet)
}
