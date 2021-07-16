// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package parallel

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	pb "github.com/lindb/lindb/proto/gen/v1/common"
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
	merger.EXPECT().close()
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
