package parallel

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

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
	merger.EXPECT().merge(gomock.Any())
	merger.EXPECT().close()
	taskManager.EXPECT().Complete("taskID")
	taskManager.EXPECT().Get("taskID").Return(taskCtx)
	ch := make(chan *series.TimeSeriesEvent)
	jobCtx := NewJobContext(context.TODO(), ch, nil, nil)
	jobManager.EXPECT().GetJob(gomock.Any()).Return(jobCtx)
	a := atomic.NewInt32(0)

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for r := range ch {
			if r.Err != nil {
				a.Inc()
			}
			wait.Done()
		}
	}()

	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID", Completed: true})
	assert.Nil(t, err)
	wait.Wait()
	assert.Equal(t, int32(1), a.Load())
}

func TestTaskReceiver_Receive_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobManager := NewMockJobManager(ctrl)
	taskManager := NewMockTaskManager(ctrl)
	jobManager.EXPECT().GetTaskManager().Return(taskManager).AnyTimes()
	receiver := NewTaskReceiver(jobManager)

	merger := NewMockResultMerger(ctrl)
	taskCtx := newTaskContext("taskID", RootTask, "parentTaskID", "parentNode", 1, merger)
	taskManager.EXPECT().Complete("taskID").MaxTimes(2)
	taskManager.EXPECT().Get("taskID").Return(taskCtx).MaxTimes(2)
	ch := make(chan *series.TimeSeriesEvent)
	jobCtx := NewJobContext(context.TODO(), ch, nil, nil)
	jobManager.EXPECT().GetJob(gomock.Any()).Return(jobCtx).MaxTimes(2)
	a := atomic.NewInt32(0)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for r := range ch {
			if r.Err != nil {
				a.Inc()
			}
			wait.Done()
		}
	}()

	err := receiver.Receive(&pb.TaskResponse{TaskID: "taskID", Completed: true, ErrMsg: "error"})
	assert.Nil(t, err)
	// ignore response
	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID", Completed: true})
	assert.Nil(t, err)
	wait.Wait()
	assert.Equal(t, int32(1), a.Load())
}
