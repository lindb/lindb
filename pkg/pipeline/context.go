package pipeline

import "sync/atomic"

// TaskContext for sharing something
type TaskContext interface {
	// RetainTask inc the ref counter of the TaskContext
	RetainTask()

	// CompleteTask dec the ref counter of the TaskContext, when the counter equals 0 means that all task are finished
	CompleteTask()
}

// BaseTaskContext implements TaskContext
type BaseTaskContext struct {
	taskCounter int32
}

// RetainTask inc the taskCounter
func (ctx *BaseTaskContext) RetainTask() {
	atomic.AddInt32(&ctx.taskCounter, 1)
}

// CompleteTask dec the taskCounter
func (ctx *BaseTaskContext) CompleteTask() {
	atomic.AddInt32(&ctx.taskCounter, -1)
}
