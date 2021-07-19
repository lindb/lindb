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

package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
)

func TestStorageStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskClientFactory := rpc.NewMockTaskClientFactory(ctrl)

	streamFactory := rpc.NewMockClientStreamFactory(ctrl)
	clientStream := protoCommonV1.NewMockTaskService_HandleClient(ctrl)
	clientStream.EXPECT().CloseSend().Return(nil).AnyTimes()
	streamFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(clientStream, nil).AnyTimes()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)

	// case 1: discovery err
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	discovery1.EXPECT().Discovery(true).Return(fmt.Errorf("err"))
	stateMachine, err := NewStorageStateMachine(context.TODO(), factory, taskClientFactory)
	assert.Error(t, err)
	assert.Nil(t, stateMachine)

	// normal case
	discovery1.EXPECT().Discovery(true).Return(nil)
	stateMachine, err = NewStorageStateMachine(context.TODO(), factory, taskClientFactory)
	assert.NoError(t, err)

	storageState2 := models.NewStorageState()
	storageState2.Name = "test2"
	data3, _ := json.Marshal(storageState2)

	stateMachine.OnCreate("/data/test2", data3)
	assert.Equal(t, 1, len(stateMachine.List()))
	storageState1 := models.NewStorageState()
	storageState1.Name = "test"
	data, _ := json.Marshal(storageState1)
	stateMachine.OnCreate("/data/test", data)
	assert.Equal(t, 2, len(stateMachine.List()))

	// cfg data err
	stateMachine.OnCreate("/data/test3", []byte{1, 2, 2})
	assert.Equal(t, 2, len(stateMachine.List()))

	// name empty
	storageState3 := models.NewStorageState()
	data3, _ = json.Marshal(storageState3)
	stateMachine.OnCreate("/data/test5", data3)
	assert.Equal(t, 2, len(stateMachine.List()))

	stateMachine.OnDelete("/data/test")
	assert.Equal(t, 1, len(stateMachine.List()))
	assert.Equal(t, *storageState2, *(stateMachine.List()[0]))

	discovery1.EXPECT().Close()
	_ = stateMachine.Close()
	_ = stateMachine.Close()
	assert.Equal(t, 0, len(stateMachine.List()))
}
