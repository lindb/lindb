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

package parallel

import (
	"context"
	"errors"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source=./context.go -destination=./context_mock.go -package=parallel

// TaskType represents the distribution query task type
type TaskType int

const (
	RootTask TaskType = iota + 1
	IntermediateTask
)

// ExecuteContext represents the execute context
type ExecuteContext interface {
	// Emit emits the time series event, and merges the events
	Emit(event *series.TimeSeriesEvent)
	// Complete completes the task with err if task execute fail
	Complete(err error)
}

// StorageExecuteContext represents the storage execute context
type StorageExecuteContext interface {
	// QueryStats returns the storage query stats
	QueryStats() *models.StorageStats
}

// BrokerExecuteContext represents the broker execute context
type BrokerExecuteContext interface {
	ExecuteContext

	// ResultCh returns the result chan
	ResultCh() chan *series.TimeSeriesEvent
	// ResultSet returns the final result set
	ResultSet() (*models.ResultSet, error)
}

type brokerExecuteContext struct {
	resultCh   chan *series.TimeSeriesEvent
	err        error
	query      *stmt.Query
	expression aggregation.Expression
	resultSet  *models.ResultSet

	stats     *models.QueryStats
	startTime int64
}

func NewBrokerExecuteContext(startTime int64, query *stmt.Query) BrokerExecuteContext {
	ctx := &brokerExecuteContext{
		startTime: startTime,
		resultCh:  make(chan *series.TimeSeriesEvent),
		resultSet: models.NewResultSet(),
		query:     query,
	}
	if query != nil {
		ctx.expression = aggregation.NewExpression(query.TimeRange, query.Interval.Int64(), query.SelectItems)
	}
	return ctx
}

func (c *brokerExecuteContext) Emit(event *series.TimeSeriesEvent) {
	//TODO merge stats for cross idc query?
	c.stats = event.Stats
	if event.Err != nil {
		c.err = event.Err
		return
	}
	start := timeutil.NowNano()
	groupByKeys := c.query.GroupBy
	groupByKeysLength := len(groupByKeys)
	for _, ts := range event.SeriesList {
		var tags map[string]string
		if groupByKeysLength > 0 {
			tagValues := tag.SplitTagValues(ts.Tags())
			if groupByKeysLength != len(tagValues) {
				// if tag values not match group by tag keys, ignore this time series
				continue
			}
			// build group by tags for final result
			tags = make(map[string]string)
			for idx, tagKey := range groupByKeys {
				tags[tagKey] = tagValues[idx]
			}
		}
		timeSeries := models.NewSeries(tags)
		c.resultSet.AddSeries(timeSeries)
		c.expression.Eval(ts)
		rs := c.expression.ResultSet()
		for fieldName, values := range rs {
			if values == nil {
				continue
			}
			points := models.NewPoints()
			it := values.Iterator()
			for it.HasNext() {
				slot, val := it.Next()
				points.AddPoint(timeutil.CalcTimestamp(c.query.TimeRange.Start, slot, c.query.Interval), val)
			}
			timeSeries.AddField(fieldName, points)
		}
		c.expression.Reset()
	}
	if c.stats != nil {
		c.stats.ExpressCost = timeutil.NowNano() - start
	}
}

func (c *brokerExecuteContext) Complete(err error) {
	if err != nil {
		c.err = err
		close(c.resultCh)
	}
}

func (c *brokerExecuteContext) ResultCh() chan *series.TimeSeriesEvent {
	return c.resultCh
}

func (c *brokerExecuteContext) ResultSet() (*models.ResultSet, error) {
	if c.err == nil {
		c.resultSet.MetricName = c.query.MetricName
		c.resultSet.StartTime = c.query.TimeRange.Start
		c.resultSet.EndTime = c.query.TimeRange.End
		c.resultSet.Interval = c.query.Interval.Int64()
	}
	if c.stats != nil {
		c.stats.Cost = timeutil.NowNano() - c.startTime
	}
	c.resultSet.Stats = c.stats

	return c.resultSet, c.err
}

type JobContext interface {
	Plan() *models.PhysicalPlan
	Query() *stmt.Query
	Emit(event *series.TimeSeriesEvent)
	Complete()
	ResultSet() chan *series.TimeSeriesEvent
	Context() context.Context
	Completed() bool
}

type jobContext struct {
	resultSet chan *series.TimeSeriesEvent
	plan      *models.PhysicalPlan
	query     *stmt.Query
	ctx       context.Context
	cancel    context.CancelFunc

	completed atomic.Bool
}

func NewJobContext(ctx context.Context, resultSet chan *series.TimeSeriesEvent, plan *models.PhysicalPlan, query *stmt.Query) JobContext {
	c, cancel := context.WithCancel(ctx)
	return &jobContext{
		resultSet: resultSet,
		plan:      plan,
		query:     query,
		ctx:       c,
		cancel:    cancel,
	}
}

func (c *jobContext) Plan() *models.PhysicalPlan {
	return c.plan
}

func (c *jobContext) Query() *stmt.Query {
	return c.query
}
func (c *jobContext) ResultSet() chan *series.TimeSeriesEvent {
	return c.resultSet
}

func (c *jobContext) Complete() {
	if c.completed.CAS(false, true) {
		//TODO send result
		close(c.resultSet)
	}
}
func (c *jobContext) Completed() bool {
	return c.completed.Load()
}

func (c *jobContext) Emit(event *series.TimeSeriesEvent) {
	c.resultSet <- event
}

func (c *jobContext) Context() context.Context {
	return c.ctx
}

// TaskContext represents the task context for distribution query and computing
type TaskContext interface {
	// TaskID returns the task id under current node
	TaskID() string
	// Type returns the task type
	TaskType() TaskType
	// ParentNode returns the parent node's indicator for sending task result
	ParentNode() string
	// ParentTaskID returns the parent node's task id for tracking task
	ParentTaskID() string
	// ReceiveResult marks receive result, decreases the num. of task tracking
	ReceiveResult(resp *pb.TaskResponse)
	// Completed returns if the task is completes
	Completed() bool
	// Error returns task's error
	Error() error
}

// taskContext represents the task context for tacking task execution state
type taskContext struct {
	taskID       string
	taskType     TaskType
	parentTaskID string
	parentNode   string
	merger       ResultMerger

	err           error
	expectResults *atomic.Int32
	completed     atomic.Bool
}

// newTaskContext creates the task context based on params
func newTaskContext(taskID string, taskType TaskType, parentTaskID string, parentNode string,
	expectResults int32, merger ResultMerger) TaskContext {
	return &taskContext{
		taskID:        taskID,
		taskType:      taskType,
		parentTaskID:  parentTaskID,
		parentNode:    parentNode,
		merger:        merger,
		expectResults: atomic.NewInt32(expectResults),
	}
}

func (c *taskContext) TaskType() TaskType {
	return c.taskType
}

// ParentNode returns the parent node's indicator for sending task result
func (c *taskContext) ParentNode() string {
	return c.parentNode
}

// ParentTaskID returns the parent node's task id for tracking task
func (c *taskContext) ParentTaskID() string {
	return c.parentTaskID
}

// TaskID returns the task id under current node
func (c *taskContext) TaskID() string {
	return c.taskID
}

// ReceiveResult marks receive result, decreases the num. of task tracking,
// if no pending task marks this task completed
func (c *taskContext) ReceiveResult(resp *pb.TaskResponse) {
	defer func() {
		// check if task completed,
		// if yes, closes the merger
		if c.Completed() && c.completed.CAS(false, true) {
			c.merger.close()
		}
	}()
	if len(resp.ErrMsg) > 0 {
		c.expectResults.Store(0)
		c.err = errors.New(resp.ErrMsg)
		return
	}
	// task is completed need return it
	if c.Completed() {
		return
	}
	// merge the response
	c.merger.merge(resp)
	// if task is completed, reduces expect result count
	if resp.Completed {
		c.expectResults.Dec()
	}
}

// Error returns task's error
func (c *taskContext) Error() error {
	return c.err
}

// Completed returns if the task is completes
func (c *taskContext) Completed() bool {
	return c.expectResults.Load() == 0
}
