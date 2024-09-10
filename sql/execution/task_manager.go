package execution

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/planner/printer"
)

type SQLTask struct {
	id         model.TaskID
	fragment   *plan.PlanFragment
	partitions []int
}

type TaskManager interface {
	SubmitTask(req *model.TaskRequest, fragment *plan.PlanFragment)
	GetTask(taskID model.TaskID) *SQLTask
}

type taskManager struct {
	tasks    map[model.TaskID]*SQLTask
	taskCh   chan *SQLTask
	taskPool concurrent.Pool

	lock sync.RWMutex
}

func NewTaskManager() TaskManager {
	mgr := &taskManager{
		tasks: make(map[model.TaskID]*SQLTask),
		// TODO: add config
		taskCh:   make(chan *SQLTask, 100),
		taskPool: concurrent.NewPool("task-exec", 10, time.Minute, metrics.NewConcurrentStatistics("task-exec", linmetric.BrokerRegistry)), // TODO: fix it
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
		id:         req.TaskID,
		fragment:   fragment,
		partitions: req.Partitions,
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
	for {
		select {
		case task := <-mgr.taskCh:
			output := buffer.NewPartitionOutputBuffer(task.id, task.fragment)
			mgr.taskPool.Submit(context.TODO(), concurrent.NewTask(func() {
				fmt.Println(task)
				printer := printer.NewPlanPrinter(printer.NewTextRender(0))
				fmt.Println("******************")
				fmt.Println(printer.PrintLogicPlan(task.fragment.Root))
				fmt.Println("******************")

				fct := NewTaskExecutionFactory()
				exec := fct.Create(task, output) // TODO:

				exec.Execute()

				fmt.Printf("task exec result\n")
			}, func(err error) {
				fmt.Printf("task exec fail %v\n", err)
				output.AddPage(nil)
			}))
		}
	}
}
