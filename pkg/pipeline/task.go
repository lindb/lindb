package pipeline

import "sync/atomic"

// Task is used to process message
type Task interface {
	// Name returns the name of the Task
	Name() string

	// Size returns the max channel size for messages
	Size() int32

	// Process process a message and may generate new messages
	Process(ctx TaskContext, message Message)

	// SetRouter set the Router for the messages produced by Process method
	SetRouter(router Router)
}

// Runtime contains a channel 、a goroutine and an Task,
// a goroutine for getting messages from a channel and using an Task to process messages
type Runtime struct {
	mailbox chan Message
	Task    Task
	closed  int32
}

// NewTaskRuntime returns a Runtime wrapping an Task and starts a goroutine to loop for processing messages
func NewTaskRuntime(task Task) (*Runtime, error) {
	var runtime = &Runtime{
		mailbox: make(chan Message, task.Size()),
		Task:    task,
	}
	go runtime.process()
	return runtime, nil
}

// Tell sends a message to an Runtime
func (runtime *Runtime) Tell(message Message) {
	runtime.mailbox <- message
}

// process loops for getting messages from a mailbox and uses an Task to process messages
func (runtime *Runtime) process() {
	defer func() {
		atomic.StoreInt32(&runtime.closed, 1)
	}()
	for message := range runtime.mailbox {
		// close current goroutine when the message is a ShutdownMessage
		_, ok := message.(*ShutdownMessage)
		if ok {
			return
		}

		ctx := message.GetContext()
		// use current Task to process message
		runtime.Task.Process(ctx, message)

		// check the TaskContext is fully finished
		if ctx != nil {
			ctx.CompleteTask()
		}
	}
}
