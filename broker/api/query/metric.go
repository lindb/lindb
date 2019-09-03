package query

import (
	"context"
	"net/http"
	"time"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/logger"
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
	var err error
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
	//TODO add timeout cfg
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()
	exec := m.executorFactory.NewBrokerExecutor(ctx, db, sql, m.replicaStateMachine, m.nodeStateMachine, m.jobManager)
	results := exec.Execute()
	resultsSet := ResultSet{Fields: make(map[string]DataPoint)}
	if results != nil {
		for result := range results {
			if result.Err != nil {
				err = result.Err
				break
			}
			ts := result.Series
			for ts.HasNext() {
				field := ts.Next()

				for field.HasNext() {
					pIt := field.Next()
					segmentStartTime := field.SegmentStartTime()

					for pIt.HasNext() {
						var dataPoint DataPoint
						ok := false
						dataPoint, ok = resultsSet.Fields[field.FieldName()]
						if !ok {
							dataPoint = DataPoint{Points: make(map[int64]float64)}
							resultsSet.Fields[field.FieldName()] = dataPoint
						}
						slot, val := pIt.Next()
						//TODO need fix it
						dataPoint.Points[int64(slot*10000)+segmentStartTime] = val
					}
				}
			}
		}
	}
	if exec.Error() != nil {
		err = exec.Error()
	}
	if err != nil {
		api.Error(w, err)
		return
	}
	logger.GetLogger("re", "re").Info("result", logger.Any("R", resultsSet))
	api.OK(w, resultsSet)
}

type ResultSet struct {
	Fields map[string]DataPoint
}

type DataPoint struct {
	Points map[int64]float64
}
