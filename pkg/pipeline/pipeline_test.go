package pipeline

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stageMessage struct {
	message string
	ctx     TaskContext
}

func (message *stageMessage) SetContext(ctx TaskContext) {
	message.ctx = ctx
}

func (message *stageMessage) GetContext() TaskContext {
	return message.ctx
}

type StageTask struct {
	name    string
	size    int32
	Counter int32
	router  Router
	mutex   sync.Mutex
}

func NewStageTask(name string, size int32) Task {
	return &StageTask{
		name: name,
		size: size,
	}
}

func (task *StageTask) Name() string {
	return task.name
}

func (task *StageTask) Size() int32 {
	return task.size
}

// Process generate two messages from one message
func (task *StageTask) Process(ctx TaskContext, message Message) {
	msg, ok := message.(*stageMessage)
	if ok {
		atomic.AddInt32(&task.Counter, 1)
		router := task.GetRouter()
		if router != nil {
			router.Tell(ctx, &stageMessage{
				message: msg.message + "xxx1",
			})
			router.Tell(ctx, &stageMessage{
				message: msg.message + "xxx2",
			})
		}
	}
}

func (task *StageTask) SetRouter(router Router) {
	task.mutex.Lock()
	defer task.mutex.Unlock()
	task.router = router
}

func (task *StageTask) GetRouter() Router {
	task.mutex.Lock()
	defer task.mutex.Unlock()
	return task.router
}

func (task *StageTask) Shutdown() {
	fmt.Println("shutdown")
}

type ConfigTest struct {
	taskSize int
}

func (config *ConfigTest) GetTaskSize() int {
	return config.taskSize
}

func (config *ConfigTest) NewTask() Task {
	return NewStageTask("test", 10)
}

func TestPipeline_Tell_AddStage(t *testing.T) {
	pipeline := new(Pipeline)
	stage1 := NewStage(&ConfigTest{
		taskSize: 2,
	})
	pipeline.AddStage(stage1)

	stage2 := NewStage(&ConfigTest{
		taskSize: 2,
	})
	pipeline.AddStage(stage2)

	ctx := new(BaseTaskContext)

	firstCounter := 4
	sendPipelineMessage(pipeline, ctx, firstCounter)
	time.Sleep(time.Second)
	assert.Equal(t, int32(8), getStageCounter(stage2))

	stage3 := NewStage(&ConfigTest{
		taskSize: 2,
	})
	pipeline.AddStage(stage3)

	sendPipelineMessage(pipeline, ctx, firstCounter)
	time.Sleep(time.Second)
	assert.Equal(t, int32(16), getStageCounter(stage3))

	assert.Equal(t, int32(0), atomic.LoadInt32(&ctx.taskCounter))

	pipeline.Shutdown()
}

func sendPipelineMessage(pipeline *Pipeline, ctx TaskContext, size int) {
	for i := 0; i < size; i++ {
		pipeline.Tell(ctx, &stageMessage{message: "pipeline"})
	}
}

func getStageCounter(stage *Stage) int32 {
	realCount := int32(0)
	for i := 0; i < len(stage.runs); i++ {
		stageTask, ok := (stage.runs[i].Task).(*StageTask)
		if ok {
			realCount += atomic.LoadInt32(&stageTask.Counter)
		}
	}
	return realCount
}
