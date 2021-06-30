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

package coordinator

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/replication"
)

func TestStateMachineFactory_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	discoveryFactory.EXPECT().GetRepo().Return(repo).AnyTimes()
	cm := replication.NewMockChannelManager(ctrl)
	cm.EXPECT().SyncReplicatorState().AnyTimes()

	factory := NewStateMachineFactory(&StateMachineCfg{
		Ctx:              context.TODO(),
		Repo:             repo,
		CurrentNode:      models.Node{IP: "1.1.1.1", Port: 9000},
		DiscoveryFactory: discoveryFactory,
		ChannelManager:   cm,
	})

	// test node state machine
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	nodeSM, err := factory.CreateActiveNodeStateMachine()
	assert.Error(t, err)
	assert.Nil(t, nodeSM)

	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	nodeSM, err = factory.CreateActiveNodeStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, nodeSM)

	// test storage state machine
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	storageStateSM, err := factory.CreateStorageStateMachine()
	assert.Error(t, err)
	assert.Nil(t, storageStateSM)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	storageStateSM, err = factory.CreateStorageStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, storageStateSM)

	// test replica status state machine
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	replicaStatusSM, err := factory.CreateReplicaStatusStateMachine()
	assert.Error(t, err)
	assert.Nil(t, replicaStatusSM)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	replicaStatusSM, err = factory.CreateReplicaStatusStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, replicaStatusSM)

	// test replicator state machine
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	replicatorSM, err := factory.CreateReplicatorStateMachine()
	assert.Error(t, err)
	assert.Nil(t, replicatorSM)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	replicatorSM, err = factory.CreateReplicatorStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, replicatorSM)

	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	dbSM, err := factory.CreateDatabaseStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, dbSM)
}
