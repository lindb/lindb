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
	"fmt"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

//go:generate mockgen -source=./task_manager.go -destination=./task_manager_mock.go -package=parallel

// TaskManager represents the task manager for current node
type TaskManager interface {
	// AllocTaskID allocates the task id for new task, before task submits
	AllocTaskID() string
	// Submit submits the task, saving task context for task tracking
	Submit(taskCtx TaskContext)
	// Complete completes the task by task id
	Complete(taskID string)
	// Get returns the task context by task id
	Get(taskID string) TaskContext

	// SendRequest sends the task request to target node based on node's indicator
	SendRequest(targetNodeID string, req *pb.TaskRequest) error
	// SendResponse sends the task response to parent node
	SendResponse(targetNodeID string, resp *pb.TaskResponse) error
}

// taskManager implements the task manager interface, tracks all task of the current node
type taskManager struct {
	currentNodeID     string
	seq               *atomic.Int64
	taskClientFactory rpc.TaskClientFactory
	taskServerFactory rpc.TaskServerFactory

	tasks sync.Map
}

// NewTaskManager creates the task manager
func NewTaskManager(currentNode models.Node,
	taskClientFactory rpc.TaskClientFactory, taskServerFactory rpc.TaskServerFactory) TaskManager {
	return &taskManager{
		currentNodeID:     (&currentNode).Indicator(),
		taskClientFactory: taskClientFactory,
		taskServerFactory: taskServerFactory,
		seq:               atomic.NewInt64(0),
	}
}

// AllocTaskID allocates the task id for new task, before task submits
func (t *taskManager) AllocTaskID() string {
	seq := t.seq.Inc()
	return fmt.Sprintf("%s-%d", t.currentNodeID, seq)
}

// Submit submits the task, saving task context for task tracking
func (t *taskManager) Submit(taskCtx TaskContext) {
	//TODO check duplicate
	t.tasks.Store(taskCtx.TaskID(), taskCtx)
}

// Complete completes the task by task id
func (t *taskManager) Complete(taskID string) {
	t.tasks.Delete(taskID)
}

// Get returns the task context by task id
func (t *taskManager) Get(taskID string) TaskContext {
	task, ok := t.tasks.Load(taskID)
	if !ok {
		return nil
	}
	taskCtx, ok := task.(TaskContext)
	if !ok {
		return nil
	}
	return taskCtx
}

// SendRequest sends the task request to target node based on node's indicator,
// if fail, returns err
func (t *taskManager) SendRequest(targetNodeID string, req *pb.TaskRequest) error {
	// todo: query from other broker
	client := t.taskClientFactory.GetTaskClient(targetNodeID)
	if client == nil {
		return fmt.Errorf("SendRequest: %w, targetNodeID: %s", errNoSendStream, targetNodeID)
	}
	if err := client.Send(req); err != nil {
		return fmt.Errorf("%w, targetNodeID: %s", errTaskSend, targetNodeID)
	}
	return nil
}

// SendResponse sends the task response to parent node,
// if fail, returns err
func (t *taskManager) SendResponse(parentNodeID string, resp *pb.TaskResponse) error {
	stream := t.taskServerFactory.GetStream(parentNodeID)
	if stream == nil {
		return fmt.Errorf("SendResponse: %w, parentNodeID: %s", errNoSendStream, parentNodeID)
	}
	if err := stream.Send(resp); err != nil {
		return fmt.Errorf("SendResponse: %w, parentNodeID: %s", errResponseSend, parentNodeID)
	}
	return nil
}
