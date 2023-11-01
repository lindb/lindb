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

package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
)

func TestStateMachineFactory_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	discoveryFct := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	fct := NewStateMachineFactory(context.TODO(), discoveryFct, nil)

	// live node sm err
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	err := fct.Start()
	assert.Error(t, err)

	// shard assignment sm err
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	err = fct.Start()
	assert.Error(t, err)
	// database limits sm err
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil).MaxTimes(2)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	err = fct.Start()
	assert.Error(t, err)
	// all state machines are ok
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil).MaxTimes(3)
	err = fct.Start()
	assert.NoError(t, err)
}

func TestStateMachineFactory_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fct := NewStateMachineFactory(context.TODO(), nil, nil)
	sm := discovery.NewMockStateMachine(ctrl)
	fct.stateMachines = append(fct.stateMachines, sm, sm)

	sm.EXPECT().Close().Return(fmt.Errorf("err"))
	sm.EXPECT().Close().Return(nil)

	fct.Stop()
}

func TestStateMachineFactory_OnNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	fct := NewStateMachineFactory(context.TODO(), nil, stateMgr)
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/key",
	})
	fct.onNodeFailure("/key")
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/key",
		Value: []byte("value"),
	})
	fct.onNodeStartup("/key", []byte("value"))
}

func TestStateMachineFactory_OnShardAssign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	fct := NewStateMachineFactory(context.TODO(), nil, stateMgr)
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:  discovery.ShardAssignmentChanged,
		Key:   "/key",
		Value: []byte("value"),
	})
	fct.onShardAssignmentChange("/key", []byte("value"))
}

func TestStateMachineFactory_CreateState(t *testing.T) {
	assert.NotNil(t, StateMachinePaths[constants.LiveNode].CreateState())
	assert.NotNil(t, StateMachinePaths[constants.ShardAssignment].CreateState())
}

func TestStateMachineFactory_DatabaseLimits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	discoveryFct := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	fct := NewStateMachineFactory(context.TODO(), discoveryFct, stateMgr)

	sm, err := fct.createDatabaseLimitsStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, sm)

	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseLimitsChanged,
		Key:   "/test",
		Value: []byte("value"),
	})
	sm.OnCreate("/test", []byte("value"))
	sm.OnDelete("/test")
}
