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

package brokerquery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source=./task_manager.go -destination=./task_manager_mock.go -package=brokerquery

// TaskManager represents the task manager for current node
type TaskManager interface {
	// SubmitMetricTask concurrently send query task to multi intermediates and leafs.
	// If intermediates are empty, the root waits response from leafs.
	// Otherwise, the roots waits response from intermediates.
	// 1. api -> metric-query -> SubmitMetricTask (query without intermediate nodes) -> leaf nodes
	//                                                                            -> leaf nodes -> response
	// 2. api -> metric-query -> SubmitMetricTask (query with intermediate nodes) <-> peer broker <->
	SubmitMetricTask(
		physicalPlan *models.PhysicalPlan,
		stmtQuery *stmt.Query,
	) (eventCh <-chan *series.TimeSeriesEvent, err error)

	// SubmitIntermediateMetricTask creates a intermediate task from leaf nodes
	// leaf response will also be merged in task-context.
	// when all intermediate response arrives to root, event will be returned to the caller
	SubmitIntermediateMetricTask(
		physicalPlan *models.PhysicalPlan,
		stmtQuery *stmt.Query,
		parentTaskID string,
	) (eventCh <-chan *series.TimeSeriesEvent)

	// SubmitMetaDataTask concurrently send query metadata task to multi leafs.
	SubmitMetaDataTask(
		physicalPlan *models.PhysicalPlan,
		suggest *stmt.Metadata,
	) (taskResponse <-chan *protoCommonV1.TaskResponse, err error)

	// SendRequest sends the task request to target node based on node's indicator
	SendRequest(targetNodeID string, req *protoCommonV1.TaskRequest) error
	// SendResponse sends the task response to parent node
	SendResponse(targetNodeID string, resp *protoCommonV1.TaskResponse) error

	// Receive receives task response from rpc handler asynchronous
	Receive(req *protoCommonV1.TaskResponse, targetNode string) error
}

// taskManager implements the task manager interface, tracks all task of the current node
type taskManager struct {
	ctx               context.Context
	currentNodeID     string
	seq               *atomic.Int64
	taskClientFactory rpc.TaskClientFactory
	taskServerFactory rpc.TaskServerFactory

	workerPool concurrent.Pool // workers for
	tasks      sync.Map        // taskID -> taskCtx
	logger     *logger.Logger
	ttl        time.Duration

	createdTaskCounter   *linmetric.BoundDeltaCounter
	aliveTaskGauge       *linmetric.BoundGauge
	emitResponseCounter  *linmetric.BoundDeltaCounter
	omitResponseCounter  *linmetric.BoundDeltaCounter
	sentRequestCounter   *linmetric.BoundDeltaCounter
	sentResponsesCounter *linmetric.BoundDeltaCounter
	sentResponseFailures *linmetric.BoundDeltaCounter
	sentRequestFailures  *linmetric.BoundDeltaCounter
}

// NewTaskManager creates the task manager
func NewTaskManager(
	ctx context.Context,
	currentNode models.Node,
	taskClientFactory rpc.TaskClientFactory,
	taskServerFactory rpc.TaskServerFactory,
	taskPool concurrent.Pool,
	ttl time.Duration,
) TaskManager {
	taskManagerScope := linmetric.NewScope("lindb.broker.query")
	tm := &taskManager{
		ctx:                  ctx,
		currentNodeID:        (&currentNode).Indicator(),
		taskClientFactory:    taskClientFactory,
		taskServerFactory:    taskServerFactory,
		seq:                  atomic.NewInt64(0),
		workerPool:           taskPool,
		logger:               logger.GetLogger("query", "TaskManager"),
		ttl:                  ttl,
		createdTaskCounter:   taskManagerScope.NewDeltaCounter("created_tasks"),
		aliveTaskGauge:       taskManagerScope.NewGauge("alive_tasks"),
		emitResponseCounter:  taskManagerScope.NewDeltaCounter("emitted_responses"),
		omitResponseCounter:  taskManagerScope.NewDeltaCounter("omitted_responses"),
		sentRequestCounter:   taskManagerScope.NewDeltaCounter("sent_requests"),
		sentResponsesCounter: taskManagerScope.NewDeltaCounter("sent_responses"),
		sentResponseFailures: taskManagerScope.NewDeltaCounter("sent_responses_failures"),
		sentRequestFailures:  taskManagerScope.NewDeltaCounter("sent_requests_failures"),
	}
	duration := ttl
	if ttl < time.Minute {
		duration = time.Minute
	}
	go tm.cleaner(duration)
	return tm
}

// cleaner cleans expired tasks in
func (t *taskManager) cleaner(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.tasks.Range(func(key, value interface{}) bool {
				taskCtx := value.(TaskContext)
				if taskCtx.Expired(t.ttl) {
					t.aliveTaskGauge.Decr()
					t.tasks.Delete(key)
				}
				return true
			})
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *taskManager) evictTask(taskID string) {
	_, loaded := t.tasks.LoadAndDelete(taskID)
	if loaded {
		t.aliveTaskGauge.Decr()
	}
}

func (t *taskManager) storeTask(taskID string, taskCtx TaskContext) {
	t.tasks.Store(taskID, taskCtx)
	t.createdTaskCounter.Incr()
	t.aliveTaskGauge.Incr()
}

func (t *taskManager) SubmitMetricTask(
	physicalPlan *models.PhysicalPlan,
	stmtQuery *stmt.Query,
) (eventCh <-chan *series.TimeSeriesEvent, err error) {
	rootTaskID := t.AllocTaskID()
	marshalledPhysicalPlan := encoding.JSONMarshal(physicalPlan)
	marshalledPayload, _ := stmtQuery.MarshalJSON()
	responseCh := make(chan *series.TimeSeriesEvent)

	taskCtx := newMetricTaskContext(
		rootTaskID,
		RootTask,
		"",
		"",
		stmtQuery,
		physicalPlan.Root.NumOfTask,
		responseCh,
	)
	t.storeTask(rootTaskID, taskCtx)

	// return the channel for reader, then send the rpc request
	// in case of too early response arriving without reader
	var (
		wg        sync.WaitGroup
		sendError atomic.Error
	)
	// send task to intermediates firstly, then to the leafs
	if len(physicalPlan.Intermediates) > 0 {
		req := &protoCommonV1.TaskRequest{
			ParentTaskID: rootTaskID,
			Type:         protoCommonV1.TaskType_Intermediate,
			RequestType:  protoCommonV1.RequestType_Data,
			PhysicalPlan: marshalledPhysicalPlan,
			Payload:      marshalledPayload,
		}
		wg.Add(len(physicalPlan.Intermediates))
		for _, intermediate := range physicalPlan.Intermediates {
			intermediate := intermediate
			t.workerPool.Submit(func() {
				defer wg.Done()
				if err := t.SendRequest(intermediate.Indicator, req); err != nil {
					sendError.Store(err)
				}
			})
		}
		wg.Wait()
	}
	// notify error to other peer nodes
	if sendError.Load() == nil {
		req := &protoCommonV1.TaskRequest{
			ParentTaskID: rootTaskID,
			Type:         protoCommonV1.TaskType_Leaf,
			RequestType:  protoCommonV1.RequestType_Data,
			PhysicalPlan: marshalledPhysicalPlan,
			Payload:      marshalledPayload,
		}
		wg.Add(len(physicalPlan.Leafs))
		for _, leaf := range physicalPlan.Leafs {
			leaf := leaf
			t.workerPool.Submit(func() {
				defer wg.Done()
				if err := t.SendRequest(leaf.Indicator, req); err != nil {
					sendError.Store(err)
				}
			})
		}
		wg.Wait()
	}

	if sendError.Load() != nil {
		t.evictTask(rootTaskID)
	}
	return responseCh, sendError.Load()
}

func (t *taskManager) SubmitIntermediateMetricTask(
	physicalPlan *models.PhysicalPlan,
	stmtQuery *stmt.Query,
	parentTaskID string,
) (eventCh <-chan *series.TimeSeriesEvent) {
	responseCh := make(chan *series.TimeSeriesEvent)
	taskCtx := newMetricTaskContext(
		parentTaskID,
		IntermediateTask,
		parentTaskID,
		physicalPlan.Root.Indicator,
		stmtQuery,
		int32(len(physicalPlan.Leafs)),
		responseCh,
	)

	t.storeTask(parentTaskID, taskCtx)
	return responseCh
}

func (t *taskManager) SubmitMetaDataTask(
	physicalPlan *models.PhysicalPlan,
	suggest *stmt.Metadata,
) (taskResponse <-chan *protoCommonV1.TaskResponse, err error) {
	taskID := t.AllocTaskID()

	suggestMarshalData, _ := suggest.MarshalJSON()
	req := &protoCommonV1.TaskRequest{
		RequestType:  protoCommonV1.RequestType_Metadata,
		ParentTaskID: taskID,
		PhysicalPlan: encoding.JSONMarshal(physicalPlan),
		Payload:      suggestMarshalData,
	}

	responseCh := make(chan *protoCommonV1.TaskResponse)
	taskCtx := newMetaDataTaskContext(
		taskID,
		RootTask,
		"",
		"",
		physicalPlan.Root.NumOfTask,
		responseCh)

	t.storeTask(taskID, taskCtx)

	var (
		wg        sync.WaitGroup
		sendError atomic.Error
	)
	wg.Add(len(physicalPlan.Leafs))
	for _, leafNode := range physicalPlan.Leafs {
		leafNode := leafNode
		t.workerPool.Submit(func() {
			defer wg.Done()
			if err := t.SendRequest(leafNode.Indicator, req); err != nil {
				sendError.Store(err)
			}
		})
	}
	wg.Wait()
	if sendError.Load() != nil {
		t.evictTask(taskID)
	}
	return responseCh, sendError.Load()
}

// AllocTaskID allocates the task id for new task, before task submits
func (t *taskManager) AllocTaskID() string {
	seq := t.seq.Inc()
	return fmt.Sprintf("%s-%d", t.currentNodeID, seq)
}

// Get returns the task context by task id
func (t *taskManager) Get(taskID string) TaskContext {
	task, ok := t.tasks.Load(taskID)
	if !ok {
		return nil
	}
	return task.(TaskContext)
}

// SendRequest sends the task request to target node based on node's indicator,
// if fail, returns err
func (t *taskManager) SendRequest(targetNodeID string, req *protoCommonV1.TaskRequest) error {
	client := t.taskClientFactory.GetTaskClient(targetNodeID)
	if client == nil {
		t.sentRequestFailures.Incr()
		return fmt.Errorf("SendRequest: %w, targetNodeID: %s", query.ErrNoSendStream, targetNodeID)
	}
	if err := client.Send(req); err != nil {
		t.sentRequestFailures.Incr()
		return fmt.Errorf("%w, targetNodeID: %s", query.ErrTaskSend, targetNodeID)
	}
	t.sentRequestCounter.Incr()
	return nil
}

// SendResponse sends the task response to parent node,
// if fail, returns err
func (t *taskManager) SendResponse(parentNodeID string, resp *protoCommonV1.TaskResponse) error {
	stream := t.taskServerFactory.GetStream(parentNodeID)
	if stream == nil {
		t.sentResponseFailures.Incr()
		return fmt.Errorf("SendResponse: %w, parentNodeID: %s", query.ErrNoSendStream, parentNodeID)
	}
	if err := stream.Send(resp); err != nil {
		t.sentResponseFailures.Incr()
		return fmt.Errorf("SendResponse: %w, parentNodeID: %s", query.ErrResponseSend, parentNodeID)
	}
	t.sentResponsesCounter.Incr()
	return nil
}

func (t *taskManager) Receive(resp *protoCommonV1.TaskResponse, targetNode string) error {
	taskCtx := t.Get(resp.TaskID)
	if taskCtx == nil {
		t.omitResponseCounter.Incr()
		return fmt.Errorf("TaskID: %s may be evicted", resp.TaskID)
	}
	t.emitResponseCounter.Incr()
	t.workerPool.Submit(func() {
		// for root task and intermediate task
		taskCtx.WriteResponse(resp, targetNode)

		if taskCtx.Done() {
			t.evictTask(resp.TaskID)
		}
	})
	return nil
}
