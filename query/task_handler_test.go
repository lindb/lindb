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

package query

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
)

type mockTaskProcessor struct {
}

func (d *mockTaskProcessor) Process(_ context.Context, _ protoCommonV1.TaskService_HandleServer, _ *protoCommonV1.TaskRequest) {
	panic("err")
}

var cfg = config.Query{
	QueryConcurrency: 10,
	IdleTimeout:      ltoml.Duration(time.Second * 5),
	Timeout:          ltoml.Duration(time.Second * 10),
}

func TestTaskHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	processor := NewMockTaskProcessor(ctrl)
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	taskServerFactory.EXPECT().Register(gomock.Any(), gomock.Any())
	taskServerFactory.EXPECT().Deregister(gomock.Any(), gomock.Any()).Return(true)
	handler := NewTaskHandler(cfg, taskServerFactory, processor,
		concurrent.NewPool("", 10, time.Second, linmetric.NewScope("22")))

	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs())
	server.EXPECT().Context().Return(ctx)
	err := handler.Handle(server)
	assert.NotNil(t, err)

	ctx = rpc.CreateIncomingContextWithNode(context.TODO(), &models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000})
	server.EXPECT().Context().Return(ctx)
	server.EXPECT().Recv().Return(nil, nil)
	server.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	processor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	_ = handler.Handle(server)
}

func TestTaskHandler_dispatch(t *testing.T) {
	handler := NewTaskHandler(cfg, nil, &mockTaskProcessor{},
		concurrent.NewPool("", 10, time.Second, linmetric.NewScope("22")))
	// test process panic
	handler.process(nil, nil)
}
