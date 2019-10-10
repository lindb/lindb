package query

import (
	"context"
	"net/http"
	"time"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/models"
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
	if exec.Error() != nil {
		api.Error(w, exec.Error())
		return
	}
	stmt := exec.Statement()
	resultSet := models.NewResultSet()
	if results != nil {
		startTime := stmt.TimeRange.Start
		expression := aggregation.NewExpression(stmt.TimeRange, stmt.Interval, stmt.SelectItems)
		for result := range results {
			if result.Err != nil {
				err = result.Err
				break
			}
			for _, ts := range result.SeriesList {
				series := models.NewSeries(ts.Tags())
				resultSet.AddSeries(series)
				expression.Eval(ts)
				rs := expression.ResultSet()
				for fieldName, values := range rs {
					points := models.NewPoints()
					it := values.Iterator()
					for it.HasNext() {
						slot, val := it.Next()
						points.AddPoint(int64(slot)*stmt.Interval+startTime, val)
					}
					series.AddField(fieldName, points)
				}
				expression.Reset()
			}
		}
	}
	if err != nil {
		api.Error(w, err)
		return
	}
	resultSet.MetricName = stmt.MetricName
	api.OK(w, resultSet)
}
