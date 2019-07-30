package parallel

import "sync/atomic"

// TaskContext represents the task context for distribution query and computing
type TaskContext interface {
	// ParentNode returns the parent node's indicator for sending task result
	ParentNode() string
	// ParentTaskID returns the parent node's task id for tracking task
	ParentTaskID() string
	// TaskID returns the task id under current node
	TaskID() string
	// ReceiveResult marks receive result, decreases the num. of task tracking
	ReceiveResult()
	// Completed returns if the task is completes
	Completed() bool
}

// taskContext represents the task context for tacking task execution state
type taskContext struct {
	taskID       string
	parentTaskID string
	parentNode   string

	expectResults int32
	completed     bool
}

// newTaskContext creates the task context based on params
func newTaskContext(taskID string, parentTaskID string, parentNode string, expectResults int32) TaskContext {
	return &taskContext{
		taskID:        taskID,
		parentTaskID:  parentTaskID,
		parentNode:    parentNode,
		expectResults: expectResults,
		completed:     false,
	}
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
func (c *taskContext) ReceiveResult() {
	pendingTask := atomic.AddInt32(&c.expectResults, -1)
	if pendingTask == 0 {
		c.completed = true
	}
}

// Completed returns if the task is completes
func (c *taskContext) Completed() bool {
	return c.completed
}
