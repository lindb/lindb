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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

type mockTaskDispatcher struct {
}

func (d *mockTaskDispatcher) Dispatch(ctx context.Context, stream pb.TaskService_HandleServer, req *pb.TaskRequest) {
	panic("err")
}

var cfg = config.Query{
	MaxWorkers:  10,
	IdleTimeout: ltoml.Duration(time.Second * 5),
	Timeout:     ltoml.Duration(time.Second * 10),
}

func TestTaskHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dispatcher := NewMockTaskDispatcher(ctrl)
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	taskServerFactory.EXPECT().Register(gomock.Any(), gomock.Any())
	taskServerFactory.EXPECT().Deregister(gomock.Any(), gomock.Any()).Return(true)
	handler := NewTaskHandler(cfg, taskServerFactory, dispatcher)

	server := pb.NewMockTaskService_HandleServer(ctrl)
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs())
	server.EXPECT().Context().Return(ctx)
	err := handler.Handle(server)
	assert.NotNil(t, err)

	ctx = rpc.CreateIncomingContextWithNode(context.TODO(), models.Node{IP: "1.1.1.1", Port: 9000})
	server.EXPECT().Context().Return(ctx)
	server.EXPECT().Recv().Return(nil, nil)
	server.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	dispatcher.EXPECT().Dispatch(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	_ = handler.Handle(server)
}

func TestTaskHandler_dispatch(t *testing.T) {
	handler := NewTaskHandler(cfg, nil, &mockTaskDispatcher{})
	// test dispatch panic
	handler.dispatch(nil, nil)
}
