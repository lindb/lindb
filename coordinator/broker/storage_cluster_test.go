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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
)

func TestStorageClusterState_SetState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskClientFactory := rpc.NewMockTaskClientFactory(ctrl)
	taskClientFactory.EXPECT().CloseTaskClient(gomock.Any()).AnyTimes()

	state := newStorageClusterState(taskClientFactory, storageFSMLogger)

	storageState := models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}})

	taskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(nil)
	state.SetState(storageState)
	assert.Equal(t, 1, len(state.connectionManager.connections))
	assert.NotNil(t, state.connectionManager.connections["1.1.1.1:9000"])

	storageState = models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.2", Port: 9000}})
	taskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(fmt.Errorf("err"))
	state.SetState(storageState)
	assert.Equal(t, 0, len(state.connectionManager.connections))

	storageState = models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.2", Port: 9000}})
	taskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(nil)
	state.SetState(storageState)
	assert.Equal(t, 1, len(state.connectionManager.connections))

	state.close()
	assert.Equal(t, 0, len(state.connectionManager.connections))
}
