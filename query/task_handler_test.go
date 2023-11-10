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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
)

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
	taskServerFactory.EXPECT().Register(gomock.Any(), gomock.Any()).AnyTimes()
	taskServerFactory.EXPECT().Deregister(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	handler := NewTaskHandler(cfg, taskServerFactory, processor,
		concurrent.NewPool("", 10, time.Second,
			metrics.NewConcurrentStatistics("test", linmetric.BrokerRegistry)))

	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs())
	server.EXPECT().Context().Return(ctx)
	err := handler.Handle(server)
	assert.Error(t, err)

	ctx = metadata.NewIncomingContext(ctx,
		metadata.Pairs(constants.RPCMetaKeyLogicNode,
			(&models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}).Indicator()))
	server.EXPECT().Context().Return(ctx).MaxTimes(2)
	server.EXPECT().Recv().Return(nil, nil)
	server.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	processor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	_ = handler.Handle(server)
}

func TestTaskHandler_dispatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	processor := NewMockTaskProcessor(ctrl)
	stream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	req := &protoCommonV1.TaskRequest{}
	handler := NewTaskHandler(cfg, nil, processor,
		concurrent.NewPool("", 10, time.Second,
			metrics.NewConcurrentStatistics("test", linmetric.BrokerRegistry)))
	// test process panic
	processor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx *flow.TaskContext, stream protoCommonV1.TaskService_HandleServer, req *protoCommonV1.TaskRequest) error {
			panic("err")
		})
	handler.process(context.Background(), stream, req)
	processor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	handler.process(context.Background(), stream, req)
	time.Sleep(300 * time.Millisecond)
}
