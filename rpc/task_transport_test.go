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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
)

const testGRPCPort = 9996

func TestTaskServerFactory(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	fct := NewTaskServerFactory()

	stream := fct.GetStream((&node).Indicator())
	assert.Nil(t, stream)

	mockServerStream := protoCommonV1.NewMockTaskService_HandleServer(ctl)

	epoch := fct.Register((&node).Indicator(), mockServerStream)
	stream = fct.GetStream((&node).Indicator())
	assert.NotNil(t, stream)

	nodes := fct.Nodes()
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, &node, nodes[0])

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

	mockTaskClient := protoCommonV1.NewMockTaskService_HandleClient(ctl)
	mockTaskClient.EXPECT().Recv().Return(&protoCommonV1.TaskResponse{}, nil).AnyTimes()
	mockTaskClient.EXPECT().CloseSend().Return(fmt.Errorf("err")).AnyTimes()
	taskService := protoCommonV1.NewMockTaskServiceClient(ctl)

	fct := NewTaskClientFactory(context.TODO(),
		&models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 123},
		GetBrokerClientConnFactory())
	receiver := NewMockTaskReceiver(ctl)
	receiver.EXPECT().Receive(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	fct.SetTaskReceiver(receiver)
	fct1 := fct.(*taskClientFactory)
	fct1.connFct = mockClientConnFct
	fct1.newTaskServiceClientFunc = func(cc *grpc.ClientConn) protoCommonV1.TaskServiceClient {
		return taskService
	}

	target := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: testGRPCPort}
	conn, _ := grpc.Dial(target.Indicator(), grpc.WithInsecure())
	mockClientConnFct.EXPECT().GetClientConn(&target).Return(conn, nil).AnyTimes()
	taskService.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(mockTaskClient, nil).AnyTimes()
	err := fct.CreateTaskClient(&target)
	assert.NoError(t, err)
	tc := fct1.taskStreams[(&target).Indicator()]
	tc.running.Store(false)
	fct1.mutex.Lock()
	tc.cli = mockTaskClient
	fct1.mutex.Unlock()

	// not create new one if exist
	target = models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: testGRPCPort}
	err = fct.CreateTaskClient(&target)
	assert.NoError(t, err)

	cli := fct.GetTaskClient((&target).Indicator())
	assert.NotNil(t, cli)

	cli = fct.GetTaskClient((&models.StatelessNode{HostIP: "", GRPCPort: testGRPCPort}).Indicator())
	assert.Nil(t, cli)

	closed, err := fct.CloseTaskClient((&target).Indicator())
	assert.NotNil(t, err)
	assert.True(t, closed)

	closed, err = fct.CloseTaskClient((&models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 1000}).Indicator())
	assert.Nil(t, err)
	assert.False(t, closed)

	// wait goroutine exit
	time.Sleep(100 * time.Millisecond)
}

func TestTaskClientFactory_handler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	receiver := NewMockTaskReceiver(ctrl)
	fct := NewTaskClientFactory(ctx, &models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 123}, GetStorageClientConnFactory())
	fct.SetTaskReceiver(receiver)

	target := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 321}
	conn, _ := grpc.Dial(target.Indicator(), grpc.WithInsecure())
	mockClientConnFct := NewMockClientConnFactory(ctrl)
	mockTaskClient := protoCommonV1.NewMockTaskService_HandleClient(ctrl)
	mockTaskClient.EXPECT().CloseSend().Return(fmt.Errorf("err")).AnyTimes()
	taskService := protoCommonV1.NewMockTaskServiceClient(ctrl)

	factory := fct.(*taskClientFactory)
	factory.newTaskServiceClientFunc = func(cc *grpc.ClientConn) protoCommonV1.TaskServiceClient {
		return taskService
	}
	factory.connFct = mockClientConnFct
	taskClient := &taskClient{
		targetID: "test",
		target:   &target,
	}
	taskClient.running.Store(true)
	gomock.InOrder(
		mockClientConnFct.EXPECT().GetClientConn(&target).Return(nil, fmt.Errorf("err")),
		mockClientConnFct.EXPECT().GetClientConn(&target).Return(conn, nil),
		taskService.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")),
		mockClientConnFct.EXPECT().GetClientConn(&target).Return(conn, nil),
		taskService.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(mockTaskClient, nil),
		mockTaskClient.EXPECT().Recv().Return(nil, fmt.Errorf("err")),
		mockClientConnFct.EXPECT().GetClientConn(&target).Return(conn, nil),
		taskService.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(mockTaskClient, nil),
		mockTaskClient.EXPECT().Recv().Return(nil, nil),
		receiver.EXPECT().Receive(gomock.Any(), gomock.Any()).Return(nil),
		mockTaskClient.EXPECT().Recv().Return(&protoCommonV1.TaskResponse{}, nil),
		receiver.EXPECT().Receive(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ *protoCommonV1.TaskResponse, _ string) error {
				taskClient.running.Store(false)
				return fmt.Errorf("err")
			}),
	)
	factory.handleTaskResponse(taskClient)
}
