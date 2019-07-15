package pipeline

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type printMessage struct {
	message string
}

func (message *printMessage) SetContext(ctx TaskContext) {

}

func (message *printMessage) GetContext() TaskContext {
	return nil
}

type PrintTask struct {
	name    string
	size    int32
	Counter int32
}

func NewPrintTask(name string, size int32) Task {
	return &PrintTask{
		name: name,
		size: size,
	}
}

func (task *PrintTask) Name() string {
	return task.name
}

func (task *PrintTask) Size() int32 {
	return task.size
}

func (task *PrintTask) Process(ctx TaskContext, message Message) {
	_, ok := message.(*printMessage)
	if ok {
		atomic.AddInt32(&task.Counter, 1)
	}
}

func (task *PrintTask) SetRouter(router Router) {

}

func (task *PrintTask) Shutdown() {
	fmt.Println("shutdown")
}

func TestNewTaskRuntime(t *testing.T) {
	runtime, err := NewTaskRuntime(NewPrintTask("PrintTask", 10))
	assert.Nil(t, err)
	assert.NotNil(t, runtime)
}

func TestRuntime_Tell(t *testing.T) {
	size := int32(10)
	count := int32(10)
	runtime, _ := NewTaskRuntime(NewPrintTask("PrintTask", size))
	for i := int32(0); i < count; i++ {
		runtime.Tell(&printMessage{message: "message"})
	}
	time.Sleep(2 * time.Second)
	printTask, ok := (runtime.Task).(*PrintTask)
	assert.Equal(t, true, ok)
	assert.Equal(t, size, printTask.size)
	assert.Equal(t, count, atomic.LoadInt32(&printTask.Counter))
	printTask.Shutdown()
}
