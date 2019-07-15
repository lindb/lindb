package pipeline

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRuntimeRouter_Tell(t *testing.T) {
	runs := make([]*Runtime, 2, 2)
	for i := 0; i < len(runs); i++ {
		runs[i], _ = NewTaskRuntime(NewPrintTask("PrintTask", 10))
	}
	router := NewRuntimeRouter(runs)
	ctx := new(BaseTaskContext)
	for i := 0; i < 100; i++ {
		router.Tell(ctx, &printMessage{
			message: "message",
		})
	}
	time.Sleep(time.Second)
	printTask0, _ := (runs[0].Task).(*PrintTask)
	printTask1, _ := (runs[1].Task).(*PrintTask)
	assert.Equal(t, int32(50), atomic.LoadInt32(&printTask0.Counter))
	assert.Equal(t, int32(50), atomic.LoadInt32(&printTask1.Counter))
}
