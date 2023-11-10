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
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/internal/linmetric"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
)

func TestTransportManager_SendResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)

	transportMgr := NewTransportManager(nil, taskServerFactory, linmetric.RootRegistry)

	// empty stream
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(nil)
	assert.Error(t, transportMgr.SendResponse("1", &protoCommonV1.TaskResponse{}))

	// send stream error
	stream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(stream).Times(2)
	stream.EXPECT().Send(gomock.Any()).Return(io.ErrClosedPipe)
	assert.Error(t, transportMgr.SendResponse("1", &protoCommonV1.TaskResponse{}))

	// send ok
	stream.EXPECT().Send(gomock.Any()).Return(nil)
	assert.Nil(t, transportMgr.SendResponse("1", &protoCommonV1.TaskResponse{}))
}

func TestTransportManager_SendRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskClientFactory := rpc.NewMockTaskClientFactory(ctrl)

	transportMgr := NewTransportManager(taskClientFactory, nil, linmetric.RootRegistry)

	// empty stream
	taskClientFactory.EXPECT().GetTaskClient(gomock.Any()).Return(nil)
	assert.Error(t, transportMgr.SendRequest("1", &protoCommonV1.TaskRequest{}))

	// send stream error
	client := protoCommonV1.NewMockTaskService_HandleClient(ctrl)
	taskClientFactory.EXPECT().GetTaskClient(gomock.Any()).Return(client).Times(2)
	client.EXPECT().Send(gomock.Any()).Return(io.ErrClosedPipe)
	assert.Error(t, transportMgr.SendRequest("1", &protoCommonV1.TaskRequest{}))

	// send ok
	client.EXPECT().Send(gomock.Any()).Return(nil)
	assert.Nil(t, transportMgr.SendRequest("1", &protoCommonV1.TaskRequest{}))
}
