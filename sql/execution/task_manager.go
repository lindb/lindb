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

package execution

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/planner/printer"
)

type SQLTask struct {
	ID          model.TaskID
	CurrentTime int64
	Fragment    *plan.PlanFragment
	Partitions  []int

	Database     string
	OutputStream string // default output stream name

	// NodeSource is the gRPC address of the node executing this task.
	// It is embedded in error messages as "[role@addr]" so the caller
	// can identify which remote node failed.
	NodeSource string
}

type TaskManager interface {
	SubmitTask(req *model.TaskRequest, fragment *plan.PlanFragment)
	GetTask(taskID model.TaskID) *SQLTask
}

type taskManager struct {
	ctx      context.Context
	tasks    map[model.TaskID]*SQLTask
	taskCh   chan *SQLTask
	taskPool concurrent.Pool

	// nodeSource is the node role and address of this node, embedded in remote-task
	// error messages as "[role@addr]".
	nodeSource string

	lock sync.RWMutex
}

// NewTaskManager creates a task manager for the given node.
// role should be one of constants.StorageRole / constants.BrokerRole / constants.RootRole.
// addr is the node's gRPC address (ip:port).
func NewTaskManager(ctx context.Context, role, addr string) TaskManager {
	mgr := &taskManager{
		ctx:        ctx,
		tasks:      make(map[model.TaskID]*SQLTask),
		nodeSource: role + "@" + addr,
		// TODO: add config
		taskCh: make(chan *SQLTask, 100),
		taskPool: concurrent.NewPool("task-exec",
			10, time.Minute, metrics.NewConcurrentStatistics("task-exec", linmetric.BrokerRegistry)), // TODO: fix it
	}

	go mgr.dispatchTask()

	return mgr
}

func (mgr *taskManager) GetTask(taskID model.TaskID) *SQLTask {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()

	return mgr.tasks[taskID]
}

func (mgr *taskManager) SubmitTask(req *model.TaskRequest, fragment *plan.PlanFragment) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	task := &SQLTask{
		CurrentTime: req.RequestContext.CurrentTime,
		ID:          req.TaskID,
		Fragment:    fragment,
		Partitions:  req.Partitions,
		NodeSource:  mgr.nodeSource,
	}

	mgr.tasks[req.TaskID] = task
	mgr.taskCh <- task
}

func (mgr *taskManager) CompleteTask(taskID model.TaskID) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	delete(mgr.tasks, taskID)
}

func (mgr *taskManager) dispatchTask() {
	// FIXME:
	for {
		select {
		case task := <-mgr.taskCh:
			output := buffer.NewPartitionOutputBuffer(task.ID, task.Fragment, task.NodeSource)
			mgr.taskPool.Submit(context.TODO(), concurrent.NewTask(func() {
				planPrinter := printer.NewPlanPrinter(printer.NewTextRender(0))
				fmt.Println("******************")
				fmt.Println(planPrinter.PrintLogicPlan(task.Fragment.Root))
				fmt.Println("******************")

				fct := NewTaskExecutionFactory()
				exec := fct.Create(context.Background(), task) // TODO:

				outputCh := make(chan arrow.RecordBatch)
				defer func() {
					close(outputCh)
				}()

				go func() {
					for record := range outputCh {
						// TODO: can merge page?
						output.AddRecord(record)
					}
					output.Complete()
				}()

				if err := exec.Execute(outputCh); err != nil {
					output.Fail(err.Error())
				}
			}, func(err error) {
				// Pool-level panic: propagate back to broker.
				output.Fail(err.Error())
			}))
		case <-mgr.ctx.Done():
			return
		}
	}
}
