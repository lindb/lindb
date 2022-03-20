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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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

	// database config  sm err
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	err = fct.Start()
	assert.Error(t, err)
	// storage state  sm err
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	err = fct.Start()
	assert.Error(t, err)
	// all state machines are ok
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	err = fct.Start()
	assert.NoError(t, err)
}

func TestStateMachineFactory_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fct := NewStateMachineFactory(context.TODO(), nil, nil)
	fct1 := fct.(*stateMachineFactory)
	sm := discovery.NewMockStateMachine(ctrl)
	fct1.stateMachines = append(fct1.stateMachines, sm, sm)

	sm.EXPECT().Close().Return(fmt.Errorf("err"))
	sm.EXPECT().Close().Return(nil)

	fct.Stop()
}

func TestStateMachineFactory_OnDatabaseConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	fct := NewStateMachineFactory(context.TODO(), nil, stateMgr)
	fct1 := fct.(*stateMachineFactory)
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/key",
	})
	fct1.onDatabaseConfigDeletion("/key")
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/key",
		Value: []byte("value"),
	})
	fct1.onDatabaseConfigChanged("/key", []byte("value"))
}

func TestStateMachineFactory_OnNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	fct := NewStateMachineFactory(context.TODO(), nil, stateMgr)
	fct1 := fct.(*stateMachineFactory)
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/key",
	})
	fct1.onNodeFailure("/key")
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/key",
		Value: []byte("value"),
	})
	fct1.onNodeStartup("/key", []byte("value"))
}

func TestStateMachineFactory_OnStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	fct := NewStateMachineFactory(context.TODO(), nil, stateMgr)
	fct1 := fct.(*stateMachineFactory)
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type: discovery.StorageDeletion,
		Key:  "/key",
	})
	fct1.onStorageDeletion("/key")
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:  discovery.StorageStateChanged,
		Key:   "/key",
		Value: []byte("value"),
	})
	fct1.onStorageStateChange("/key", []byte("value"))
}
