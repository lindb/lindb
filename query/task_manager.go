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
	"fmt"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/logger"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./task_manager.go -destination=./task_manager_mock.go -package=query

// TaskManager represents the task manager for current node.
// FIXME: need remove when target offline
type TaskManager interface {
	rpc.TaskReceiver

	// AddTask adds task context by request id.
	AddTask(requestID string, taskCtx context.TaskContext)
	// RemoveTask removes task context by request id.
	RemoveTask(requestID string)
}

// taskManager implements the task manager interface, tracks all task of the current node.
type taskManager struct {
	workerPool concurrent.Pool                // workers for
	tasks      map[string]context.TaskContext // request id => task context

	statistics *metrics.QueryStatistics
	mutex      sync.RWMutex

	logger *logger.Logger
}

// NewTaskManager creates the task manager.
func NewTaskManager(workerPool concurrent.Pool, registry *linmetric.Registry) TaskManager {
	mgr := &taskManager{
		workerPool: workerPool,
		tasks:      make(map[string]context.TaskContext),
		statistics: metrics.NewQueryStatistics(registry),
		logger:     logger.GetLogger("Query", "TaskManager"),
	}
	mgr.statistics.AliveTask.SetGetValueFn(func(val *atomic.Float64) {
		mgr.mutex.Lock()
		defer mgr.mutex.Unlock()

		val.Store(float64(len(mgr.tasks)))
	})
	return mgr
}

// AddTask adds task context by request id.
func (mgr *taskManager) AddTask(requestID string, taskCtx context.TaskContext) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.statistics.CreatedTasks.Incr()
	mgr.tasks[requestID] = taskCtx
}

// RemoveTask removes task context by request id.
func (mgr *taskManager) RemoveTask(requestID string) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.statistics.ExpireTasks.Incr()
	delete(mgr.tasks, requestID)
}

// Receive receives task response from rpc handler asynchronous.
func (mgr *taskManager) Receive(resp *protoCommonV1.TaskResponse, fromNode string) error {
	taskCtx := mgr.get(resp.RequestID)
	if taskCtx == nil {
		mgr.statistics.OmitResponse.Incr()
		return fmt.Errorf("request may be evicted")
	}
	mgr.statistics.EmitResponse.Incr()
	mgr.workerPool.Submit(taskCtx.Context(), concurrent.NewTask(func() {
		// for root task and intermediate task, handle task response
		taskCtx.HandleResponse(resp, fromNode)
	}, nil))
	return nil
}

// get returns the task context by request id.
func (mgr *taskManager) get(requestID string) context.TaskContext {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	if taskCtx, ok := mgr.tasks[requestID]; ok {
		return taskCtx
	}
	return nil
}
