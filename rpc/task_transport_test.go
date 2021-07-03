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

package rpc

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc/proto/common"
)

const testGRPCPort = 9999

func TestTaskServerFactory(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	fct := NewTaskServerFactory()

	stream := fct.GetStream((&node).Indicator())
	assert.Nil(t, stream)

	mockServerStream := common.NewMockTaskService_HandleServer(ctl)

	epoch := fct.Register((&node).Indicator(), mockServerStream)
	stream = fct.GetStream((&node).Indicator())
	assert.NotNil(t, stream)

	nodes := fct.Nodes()
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, node, nodes[0])

	ok := fct.Deregister(10, (&node).Indicator())
	assert.False(t, ok)
	ok = fct.Deregister(epoch, (&node).Indicator())
	assert.True(t, ok)
	// parse node error
	fct.Register("node_err", mockServerStream)
	nodes = fct.Nodes()
	assert.Equal(t, 0, len(nodes))
}

func TestTaskClientFactory(t *testing.T) {
	ctl := gomock.NewController(t)
	defer func() {
		ctl.Finish()
	}()

	mockClientConnFct := NewMockClientConnFactory(ctl)

	mockTaskClient := common.NewMockTaskService_HandleClient(ctl)
	mockTaskClient.EXPECT().Recv().Return(nil, nil).AnyTimes()
	mockTaskClient.EXPECT().CloseSend().Return(fmt.Errorf("err")).AnyTimes()
	taskService := common.NewMockTaskServiceClient(ctl)

	fct := NewTaskClientFactory(models.Node{IP: "127.0.0.1", Port: 123})
	receiver := NewMockTaskReceiver(ctl)
	receiver.EXPECT().Receive(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	fct.SetTaskReceiver(receiver)
	fct1 := fct.(*taskClientFactory)
	fct1.connFct = mockClientConnFct
	fct1.newTaskServiceClientFunc = func(cc *grpc.ClientConn) common.TaskServiceClient {
		return taskService
	}

	target := models.Node{IP: "127.0.0.1", Port: testGRPCPort}
	conn, _ := grpc.Dial(target.Indicator(), grpc.WithInsecure())
	mockClientConnFct.EXPECT().GetClientConn(target).Return(conn, nil).AnyTimes()
	taskService.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(mockTaskClient, nil).AnyTimes()
	err := fct.CreateTaskClient(target)
	assert.NoError(t, err)
	tc := fct1.taskStreams[(&target).Indicator()]
	tc.running.Store(false)
	fct1.mutex.Lock()
	tc.cli = mockTaskClient
	fct1.mutex.Unlock()

	// not create new one if exist
	target = models.Node{IP: "127.0.0.1", Port: testGRPCPort}
	err = fct.CreateTaskClient(target)
	assert.NoError(t, err)

	cli := fct.GetTaskClient((&target).Indicator())
	assert.NotNil(t, cli)

	cli = fct.GetTaskClient((&models.Node{IP: "", Port: testGRPCPort}).Indicator())
	assert.Nil(t, cli)

	closed, err := fct.CloseTaskClient((&target).Indicator())
	assert.NotNil(t, err)
	assert.False(t, closed)

	closed, err = fct.CloseTaskClient((&models.Node{IP: "127.0.0.1", Port: 1000}).Indicator())
	assert.Nil(t, err)
	assert.False(t, closed)
}

func TestTaskClientFactory_handler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	receiver := NewMockTaskReceiver(ctrl)
	fct := NewTaskClientFactory(models.Node{IP: "127.0.0.1", Port: 123})
	fct.SetTaskReceiver(receiver)

	target := models.Node{IP: "127.0.0.1", Port: 321}
	conn, _ := grpc.Dial(target.Indicator(), grpc.WithInsecure())
	mockClientConnFct := NewMockClientConnFactory(ctrl)
	mockTaskClient := common.NewMockTaskService_HandleClient(ctrl)
	mockTaskClient.EXPECT().CloseSend().Return(fmt.Errorf("err")).AnyTimes()
	taskService := common.NewMockTaskServiceClient(ctrl)

	factory := fct.(*taskClientFactory)
	factory.newTaskServiceClientFunc = func(cc *grpc.ClientConn) common.TaskServiceClient {
		return taskService
	}
	factory.connFct = mockClientConnFct
	taskClient := &taskClient{
		targetID: "test",
		target:   target,
	}
	taskClient.running.Store(true)
	gomock.InOrder(
		mockClientConnFct.EXPECT().GetClientConn(target).Return(nil, fmt.Errorf("err")),
		mockClientConnFct.EXPECT().GetClientConn(target).Return(conn, nil),
		taskService.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")),
		mockClientConnFct.EXPECT().GetClientConn(target).Return(conn, nil),
		taskService.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(mockTaskClient, nil),
		mockTaskClient.EXPECT().Recv().Return(nil, fmt.Errorf("err")),
		mockClientConnFct.EXPECT().GetClientConn(target).Return(conn, nil),
		taskService.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(mockTaskClient, nil),
		mockTaskClient.EXPECT().Recv().Return(nil, nil),
		receiver.EXPECT().Receive(gomock.Any()).Return(nil),
		mockTaskClient.EXPECT().Recv().Return(nil, nil),
		receiver.EXPECT().Receive(gomock.Any()).DoAndReturn(func(req *common.TaskResponse) error {
			taskClient.running.Store(false)
			return fmt.Errorf("err")
		}),
	)
	factory.handleTaskResponse(taskClient)
}
