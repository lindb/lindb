package parallel

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

func TestTaskReceiver_Receive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobManager := NewMockJobManager(ctrl)
	taskManager := NewMockTaskManager(ctrl)
	jobManager.EXPECT().GetTaskManager().Return(taskManager).AnyTimes()

	receiver := NewTaskReceiver(jobManager)
	taskManager.EXPECT().Get("taskID").Return(nil)
	err := receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.Nil(t, err)

	merger := NewMockResultMerger(ctrl)
	taskCtx := newTaskContext("taskID", RootTask, "parentTaskID", "parentNode", 1, merger)
	c := taskCtx.(*taskContext)
	c.err = fmt.Errorf("err")
	merger.EXPECT().Merge(gomock.Any())
	taskManager.EXPECT().Complete("taskID")
	taskManager.EXPECT().Get("taskID").Return(taskCtx)
	ch := make(chan *series.TimeSeriesEvent)
	jobCtx := NewJobContext(ch, nil, nil)
	jobManager.EXPECT().GetJob(gomock.Any()).Return(jobCtx)
	a := int32(0)
	go func() {
		for r := range ch {
			if r.Err != nil {
				atomic.AddInt32(&a, 1)
			}
		}
	}()

	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID", Completed: true})
	assert.Nil(t, err)
	time.Sleep(300)
	assert.Equal(t, int32(1), atomic.LoadInt32(&a))
}
