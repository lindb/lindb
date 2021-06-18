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
	"github.com/lindb/lindb/service"
)

func TestStateMachineFactory_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	discoveryFactory.EXPECT().GetRepo().Return(repo).AnyTimes()
	shardAssignSVR := service.NewMockShardAssignService(ctrl)
	cm := replication.NewMockChannelManager(ctrl)
	cm.EXPECT().SyncReplicatorState().AnyTimes()

	factory := NewStateMachineFactory(&StateMachineCfg{
		Ctx:              context.TODO(),
		Repo:             repo,
		CurrentNode:      models.Node{IP: "1.1.1.1", Port: 9000},
		DiscoveryFactory: discoveryFactory,
		ShardAssignSRV:   shardAssignSVR,
		ChannelManager:   cm,
	})

	// test node state machine
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	nodeSM, err := factory.CreateNodeStateMachine()
	assert.NotNil(t, err)
	assert.Nil(t, nodeSM)

	discovery1.EXPECT().Discovery().Return(nil)
	nodeSM, err = factory.CreateNodeStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, nodeSM)

	// test storage state machine
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(2)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	storageStateSM, err := factory.CreateStorageStateMachine()
	assert.NotNil(t, err)
	assert.Nil(t, storageStateSM)
	discovery1.EXPECT().Discovery().Return(nil)
	storageStateSM, err = factory.CreateStorageStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, storageStateSM)

	// test replica status state machine
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(2)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	replicaStatusSM, err := factory.CreateReplicaStatusStateMachine()
	assert.NotNil(t, err)
	assert.Nil(t, replicaStatusSM)
	discovery1.EXPECT().Discovery().Return(nil)
	replicaStatusSM, err = factory.CreateReplicaStatusStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, replicaStatusSM)

	// test replicator state machine
	shardAssignSVR.EXPECT().List().Return(nil, fmt.Errorf("err"))
	replicatorSM, err := factory.CreateReplicatorStateMachine()
	assert.NotNil(t, err)
	assert.Nil(t, replicatorSM)
	shardAssignSVR.EXPECT().List().Return(nil, nil)
	discovery1.EXPECT().Discovery().Return(nil)
	replicatorSM, err = factory.CreateReplicatorStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, replicatorSM)

	discovery1.EXPECT().Discovery().Return(nil)
	dbSM, err := factory.CreateDatabaseStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, dbSM)
}
