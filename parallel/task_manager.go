package parallel

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/lindb/lindb/models"
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

	// GetTaskSenderManager returns the task sender manager
	GetTaskSenderManager() TaskSenderManager
}

// taskManager implements the task manager interface, tracks all task of the current node
type taskManager struct {
	currentNodeID string
	seq           int64
	taskSender    TaskSenderManager

	tasks sync.Map
}

// NewTaskManager creates the task manager
func NewTaskManager(currentNode models.Node, taskSender TaskSenderManager) TaskManager {
	return &taskManager{
		currentNodeID: (&currentNode).Indicator(),
		taskSender:    taskSender,
	}
}

// AllocTaskID allocates the task id for new task, before task submits
func (t *taskManager) AllocTaskID() string {
	seq := atomic.AddInt64(&t.seq, 1)
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

// GetTaskSenderManager returns the task sender manager
func (t *taskManager) GetTaskSenderManager() TaskSenderManager {
	return t.taskSender
}
