package parallel

import (
	"context"
	"errors"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source=./context.go -destination=./context_mock.go -package=parallel

var execLogger = logger.GetLogger("parallel", "execute")

// TaskType represents the distribution query task type
type TaskType int

const (
	RootTask TaskType = iota + 1
	IntermediateTask
)

// ExecuteContext represents the execute context
type ExecuteContext interface {
	// RetainTask adds the task count
	RetainTask(tasks int32)
	// Emit emits the time series event, and merges the events
	Emit(event *series.TimeSeriesEvent)
	// Complete completes the task with err if task execute fail
	Complete(err error)
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
}

func NewBrokerExecuteContext(query *stmt.Query) BrokerExecuteContext {
	ctx := &brokerExecuteContext{
		resultCh:  make(chan *series.TimeSeriesEvent),
		resultSet: models.NewResultSet(),
		query:     query,
	}
	if query != nil {
		ctx.expression = aggregation.NewExpression(query.TimeRange, query.Interval, query.SelectItems)
	}
	return ctx
}

func (c *brokerExecuteContext) RetainTask(tasks int32) {
}

func (c *brokerExecuteContext) Emit(event *series.TimeSeriesEvent) {
	if event.Err != nil {
		c.err = event.Err
		return
	}

	for _, ts := range event.SeriesList {
		timeSeries := models.NewSeries(ts.Tags())
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
				points.AddPoint(int64(slot)*c.query.Interval+c.query.TimeRange.Start, val)
			}
			timeSeries.AddField(fieldName, points)
		}
		c.expression.Reset()
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
	c.resultSet.MetricName = c.query.MetricName
	c.resultSet.StartTime = c.query.TimeRange.Start
	c.resultSet.EndTime = c.query.TimeRange.End
	c.resultSet.Interval = c.query.Interval
	return c.resultSet, c.err
}

// storageExecuteContext represents the storage query executor context
type storageExecuteContext struct {
	ctx         context.Context
	taskCounter atomic.Int32 // pending task ref counter
	stream      pb.TaskService_HandleServer
	req         *pb.TaskRequest

	timeSeriesList []*pb.TimeSeries

	completed atomic.Bool

	err error
}

func newStorageExecutorContext(ctx context.Context,
	req *pb.TaskRequest,
	stream pb.TaskService_HandleServer,
) ExecuteContext {
	return &storageExecuteContext{
		ctx:    ctx,
		req:    req,
		stream: stream,
	}
}

func (c *storageExecuteContext) RetainTask(tasks int32) {
	c.taskCounter.Add(tasks)
}

func (c *storageExecuteContext) Emit(event *series.TimeSeriesEvent) {
	if c.completed.Load() {
		return
	}
	if event.Err != nil {
		c.err = event.Err
		return
	}

	for _, ts := range event.SeriesList {
		fields := make(map[string][]byte)
		for ts.HasNext() {
			fieldIt := ts.Next()
			data, err := series.MarshalIterator(fieldIt)
			if err != nil || len(data) == 0 {
				continue
			}

			fields[fieldIt.FieldName()] = data
		}
		if len(fields) > 0 {
			c.timeSeriesList = append(c.timeSeriesList, &pb.TimeSeries{
				Tags:   ts.Tags(),
				Fields: fields,
			})
		}
	}
}

func (c *storageExecuteContext) Complete(err error) {
	newVal := c.taskCounter.Dec()
	if err != nil {
		c.err = err
	}
	// if all tasks completed, close result channel
	if newVal == 0 {
		c.completed.Store(true)
		errMsg := ""
		var data []byte
		if c.err != nil {
			errMsg = c.err.Error()
		} else {
			seriesList := pb.TimeSeriesList{
				TimeSeriesList: c.timeSeriesList,
			}
			// no error
			data, _ = seriesList.Marshal()
		}

		// send result to upstream
		if err := c.stream.Send(&pb.TaskResponse{
			JobID:     c.req.JobID,
			TaskID:    c.req.ParentTaskID,
			Completed: true,
			Payload:   data,
			ErrMsg:    errMsg,
		}); err != nil {
			execLogger.Error("send storage execute result", logger.Error(err))
		}
	}
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

	// check if task completed,
	// if yes, closes the merger
	if c.Completed() {
		c.merger.close()
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
