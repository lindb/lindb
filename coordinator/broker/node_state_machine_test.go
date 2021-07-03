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
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
)

var currentNode = models.Node{IP: "1.1.1.2", Port: 2080}

func TestNodeStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	taskClientFactory := rpc.NewMockTaskClientFactory(ctrl)

	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err := NewNodeStateMachine(context.TODO(), currentNode, factory, taskClientFactory)
	assert.Error(t, err)

	// normal case
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(nil)
	stateMachine, err := NewNodeStateMachine(context.TODO(), currentNode, factory, taskClientFactory)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))
	assert.Equal(t, currentNode, stateMachine.GetCurrentNode())
}

func TestNodeStateMachine_Listener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	taskClientFactory := rpc.NewMockTaskClientFactory(ctrl)

	// normal case
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(nil)
	stateMachine, err := NewNodeStateMachine(context.TODO(), currentNode, factory, taskClientFactory)
	assert.NoError(t, err)

	taskClientFactory.EXPECT().CreateTaskClient(gomock.Any())
	activeNode := models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}}
	data, _ := json.Marshal(&activeNode)
	stateMachine.OnCreate("/data/test", data)

	assert.Equal(t, 1, len(stateMachine.GetActiveNodes()))
	assert.Equal(t, activeNode, stateMachine.GetActiveNodes()[0])
	assert.Equal(t, currentNode, stateMachine.GetCurrentNode())

	taskClientFactory.EXPECT().CreateTaskClient(gomock.Any())
	stateMachine.OnCreate("/data/test2", []byte{1, 1})
	assert.Equal(t, 1, len(stateMachine.GetActiveNodes()))

	taskClientFactory.EXPECT().CloseTaskClient(gomock.Any()).Return(true, nil)
	stateMachine.OnDelete("/data/test")
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))

	// add
	stateMachine.OnCreate("/data/test", data)
	assert.Equal(t, 1, len(stateMachine.GetActiveNodes()))

	taskClientFactory.EXPECT().CloseTaskClient(gomock.Any()).Return(true, io.ErrClosedPipe)
	discovery1.EXPECT().Close()
	_ = stateMachine.Close()
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))
}
