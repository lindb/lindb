package pipeline

import "sync/atomic"

// Router routes a TaskContext and a Message to a Runtime Task
type Router interface {
	Tell(ctx TaskContext, message Message)
}

// RuntimeRouter implements Router and contains a list of Runtime Tasks
type RuntimeRouter struct {
	runs    []*Runtime
	counter int32
}

// NewRuntimeRouter returns a RuntimeRouter
func NewRuntimeRouter(runs []*Runtime) Router {
	return &RuntimeRouter{
		runs: runs,
	}
}

// Tell routes the message by round robin
func (router *RuntimeRouter) Tell(ctx TaskContext, message Message) {
	message.SetContext(ctx)
	ctx.RetainTask()
	atomic.AddInt32(&router.counter, 1)
	if atomic.LoadInt32(&router.counter) < 0 {
		atomic.StoreInt32(&router.counter, 0)
	}
	index := atomic.LoadInt32(&router.counter) % int32(len(router.runs))
	router.runs[index].Tell(message)
}
