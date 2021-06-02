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
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestStorageStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskClientFactory := rpc.NewMockTaskClientFactory(ctrl)

	streamFactory := rpc.NewMockClientStreamFactory(ctrl)
	clientStream := pb.NewMockTaskService_HandleClient(ctrl)
	clientStream.EXPECT().CloseSend().Return(nil).AnyTimes()
	streamFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(clientStream, nil).AnyTimes()

	repo := state.NewMockRepository(ctrl)
	factory := discovery.NewMockFactory(ctrl)
	factory.EXPECT().GetRepo().Return(repo).AnyTimes()
	discovery1 := discovery.NewMockDiscovery(ctrl)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	_, err := NewStorageStateMachine(context.TODO(), factory, taskClientFactory)
	assert.NotNil(t, err)

	storageState := models.NewStorageState()
	storageState.Name = "test"
	data, _ := json.Marshal(storageState)
	data2, _ := json.Marshal(models.NewStorageState())

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Value: data},
		{Value: []byte{1, 1, 2}},
		{Value: data2},
	}, nil)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err = NewStorageStateMachine(context.TODO(), factory, taskClientFactory)
	assert.NotNil(t, err)

	// normal case
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Value: data},
		{Value: []byte{1, 1, 3}},
	}, nil)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(nil)
	stateMachine, err := NewStorageStateMachine(context.TODO(), factory, taskClientFactory)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(stateMachine.List()))
	assert.Equal(t, *storageState, *(stateMachine.List()[0]))

	storageState2 := models.NewStorageState()
	storageState2.Name = "test2"
	data3, _ := json.Marshal(storageState2)

	stateMachine.OnCreate("/data/test2", data3)
	assert.Equal(t, 2, len(stateMachine.List()))

	stateMachine.OnDelete("/data/test")
	assert.Equal(t, 1, len(stateMachine.List()))
	assert.Equal(t, *storageState2, *(stateMachine.List()[0]))

	discovery1.EXPECT().Close()
	_ = stateMachine.Close()
	assert.Equal(t, 0, len(stateMachine.List()))
}
