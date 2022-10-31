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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
)

var (
	node = models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	}
)

func TestClientConnFactory(t *testing.T) {
	fct := GetBrokerClientConnFactory()

	conn1, err := fct.GetClientConn(&models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	})
	assert.NoError(t, err)

	conn11, err := fct.GetClientConn(&models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	})
	assert.NoError(t, err)

	conn2, err := fct.GetClientConn(&models.StatelessNode{
		HostIP:   "1.1.1.1",
		GRPCPort: 456,
	})
	assert.NoError(t, err)

	assert.Same(t, conn1, conn11)
	assert.NotSame(t, conn1, conn2)

	grpcDialFn = func(_ string, _ ...grpc.DialOption) (*grpc.ClientConn, error) {
		return nil, fmt.Errorf("err")
	}
	defer func() {
		grpcDialFn = grpc.Dial
	}()
	// connect failure
	conn3, err := fct.GetClientConn(&models.StatelessNode{
		HostIP:   "1.1.1.1",
		GRPCPort: 789,
	})
	assert.Error(t, err)
	assert.Nil(t, conn3)

	// test close
	err = fct.CloseClientConn(&models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	})
	assert.NoError(t, err)
}

func TestClientStreamFactory_CreateTaskClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTaskServiceClientFn = protoCommonV1.NewTaskServiceClient
		ctrl.Finish()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connFct := NewMockClientConnFactory(ctrl)

	factory := NewClientStreamFactory(ctx, &models.StatelessNode{HostIP: "127.0.0.2", GRPCPort: 9000}, connFct)
	target := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 9000}

	// case 1: get conn failure
	connFct.EXPECT().GetClientConn(gomock.Any()).Return(nil, fmt.Errorf("err"))
	cli, err := factory.CreateTaskClient(&target)
	assert.Nil(t, cli)
	assert.Error(t, err)

	// case 2: create client failure
	connFct.EXPECT().GetClientConn(gomock.Any()).Return(&grpc.ClientConn{}, nil)
	taskServiceClient := protoCommonV1.NewMockTaskServiceClient(ctrl)
	newTaskServiceClientFn = func(cc *grpc.ClientConn) protoCommonV1.TaskServiceClient {
		return taskServiceClient
	}
	taskServiceClient.EXPECT().Handle(gomock.Any()).Return(nil, fmt.Errorf("err"))
	cli, err = factory.CreateTaskClient(&target)
	assert.Nil(t, cli)
	assert.Error(t, err)
}

func TestClientStreamFactory_CreateReplicaServiceClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTaskServiceClientFn = protoCommonV1.NewTaskServiceClient
		ctrl.Finish()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connFct := NewMockClientConnFactory(ctrl)

	factory := NewClientStreamFactory(ctx, &models.StatelessNode{HostIP: "127.0.0.2", GRPCPort: 9000}, connFct)
	target := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 9000}
	// case 1: get conn failure
	connFct.EXPECT().GetClientConn(gomock.Any()).Return(nil, fmt.Errorf("err"))
	cli, err := factory.CreateReplicaServiceClient(&target)
	assert.Nil(t, cli)
	assert.Error(t, err)
	// case 2: create client successfully
	connFct.EXPECT().GetClientConn(gomock.Any()).Return(&grpc.ClientConn{}, nil)
	cli, err = factory.CreateReplicaServiceClient(&target)
	assert.NotNil(t, cli)
	assert.NoError(t, err)
}

func TestClientStreamFactory_CreateWriteServiceClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTaskServiceClientFn = protoCommonV1.NewTaskServiceClient
		ctrl.Finish()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connFct := NewMockClientConnFactory(ctrl)

	factory := NewClientStreamFactory(ctx, &models.StatelessNode{HostIP: "127.0.0.2", GRPCPort: 9000}, connFct)
	target := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 9000}
	// case 1: get conn failure
	connFct.EXPECT().GetClientConn(gomock.Any()).Return(nil, fmt.Errorf("err"))
	cli, err := factory.CreateWriteServiceClient(&target)
	assert.Nil(t, cli)
	assert.Error(t, err)
	// case 2: create client successfully
	connFct.EXPECT().GetClientConn(gomock.Any()).Return(&grpc.ClientConn{}, nil)
	cli, err = factory.CreateWriteServiceClient(&target)
	assert.NotNil(t, cli)
	assert.NoError(t, err)
}

func TestClientStreamFactory_GetValueFromContext(t *testing.T) {
	val, err := GetStringFromContext(context.TODO(), "test_key")
	assert.Error(t, err)
	assert.Empty(t, val)

	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs("key", "value1", "key", "value2"))
	val, err = GetStringFromContext(ctx, "key")
	assert.Error(t, err)
	assert.Empty(t, val)

	ctx = metadata.NewIncomingContext(context.TODO(), metadata.Pairs(constants.RPCMetaKeyLogicNode, "1.1.1.1:9999"))
	val, err = GetStringFromContext(ctx, constants.RPCMetaKeyLogicNode)
	assert.NoError(t, err)
	assert.Equal(t, "1.1.1.1:9999", val)

	node, err := GetLogicNodeFromContext(context.TODO())
	assert.Error(t, err)
	assert.Nil(t, node)

	node, err = GetLogicNodeFromContext(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, node)
}
