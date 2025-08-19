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

	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/planner/printer"
)

type SQLTask struct {
	id          model.TaskID
	currentTime int64
	fragment    *plan.PlanFragment
	partitions  []int
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

	lock sync.RWMutex
}

func NewTaskManager(ctx context.Context) TaskManager {
	mgr := &taskManager{
		ctx:   ctx,
		tasks: make(map[model.TaskID]*SQLTask),
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
		currentTime: req.RequestContext.CurrentTime,
		id:          req.TaskID,
		fragment:    fragment,
		partitions:  req.Partitions,
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
			output := buffer.NewPartitionOutputBuffer(task.id, task.fragment)
			mgr.taskPool.Submit(context.TODO(), concurrent.NewTask(func() {
				fmt.Println(task)
				planPrinter := printer.NewPlanPrinter(printer.NewTextRender(0))
				fmt.Println("******************")
				fmt.Println(planPrinter.PrintLogicPlan(task.fragment.Root))
				fmt.Println("******************")

				fct := NewTaskExecutionFactory()
				exec := fct.Create(task) // TODO:

				outputCh := make(chan *types.Page)
				defer func() {
					close(outputCh)
				}()

				go func() {
					for page := range outputCh {
						// TODO: can merge page?
						output.AddPage(page)
						fmt.Println("send page done...")
					}
					fmt.Println("task done 2.....")
					output.Complete()
					fmt.Println("task done.....")
				}()

				if err := exec.Execute(outputCh); err != nil {
					output.AddPage(&types.Page{Error: err.Error()})
				}
				fmt.Printf("task exec result\n")
			}, func(err error) {
				fmt.Printf("task exec fail %v\n", err)
				output.AddPage(&types.Page{Error: err.Error()})
			}))
		case <-mgr.ctx.Done():
			return
		}
	}
}
