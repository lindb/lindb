package query

import (
	"fmt"
	"net/http"

	"github.com/lindb/lindb/broker/api"
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
	sql, err := api.GetParamsFromRequest("sql", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	exec := m.executorFactory.NewBrokerExecutor(db, sql, m.replicaStateMachine, m.nodeStateMachine, m.jobManager)
	results := exec.Execute()
	if results != nil {
		for result := range results {
			//TODO need handle result
			fmt.Println(result)
		}
	}
	if err := exec.Error(); err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, "ok")
}
