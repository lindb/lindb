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

package root

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/pkg/logger"
)

func TestStateMachineFactory_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	stateMgr.EXPECT().SetStateMachineFactory(gomock.Any()).AnyTimes()
	discoveryFct := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	fct := NewStateMachineFactory(context.TODO(), discoveryFct, stateMgr)

	// storage broker config sm err
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	err := fct.Start()
	assert.Error(t, err)

	// database config sm err
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	err = fct.Start()
	assert.Error(t, err)

	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	err = fct.Start()
	assert.NoError(t, err)
}

func TestStateMachineFactory_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fct := &stateMachineFactory{
		logger: logger.GetLogger("Test", "Root"),
	}
	sm := discovery.NewMockStateMachine(ctrl)
	fct.stateMachines = append(fct.stateMachines, sm, sm)

	sm.EXPECT().Close().Return(fmt.Errorf("err"))
	sm.EXPECT().Close().Return(nil)

	fct.Stop()
}

func TestStateMachineFactory_BrokerCfg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	stateMgr.EXPECT().SetStateMachineFactory(gomock.Any())
	discoveryFct := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	fct := NewStateMachineFactory(context.TODO(), discoveryFct, stateMgr)

	fct1 := fct.(*stateMachineFactory)

	sm, err := fct1.createBrokerConfigStateMachine()

	assert.NoError(t, err)
	assert.NotNil(t, sm)
	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:  discovery.BrokerConfigChanged,
		Key:   "/test",
		Value: []byte("value"),
	})
	sm.OnCreate("/test", []byte("value"))

	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type: discovery.BrokerConfigDeletion,
		Key:  "/test",
	})
	sm.OnDelete("/test")
}

func TestStateMachineFactory_BrokerNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	stateMgr.EXPECT().SetStateMachineFactory(gomock.Any())
	discoveryFct := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	fct := NewStateMachineFactory(context.TODO(), discoveryFct, stateMgr)
	fct1 := fct.(*stateMachineFactory)

	sm, err := fct1.createBrokerNodeStateMachine("test", discoveryFct)
	assert.NoError(t, err)
	assert.NotNil(t, sm)

	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:       discovery.NodeStartup,
		Key:        "/test",
		Value:      []byte("value"),
		Attributes: map[string]string{brokerNameKey: "test"},
	})
	sm.OnCreate("/test", []byte("value"))

	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:       discovery.NodeFailure,
		Key:        "/test",
		Attributes: map[string]string{brokerNameKey: "test"},
	})
	sm.OnDelete("/test")
}

func TestStateMachineFactory_DatabaseCfg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := NewMockStateManager(ctrl)
	stateMgr.EXPECT().SetStateMachineFactory(gomock.Any())
	discoveryFct := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	fct := NewStateMachineFactory(context.TODO(), discoveryFct, stateMgr)
	fct1 := fct.(*stateMachineFactory)

	sm, err := fct1.createDatabaseConfigStateMachine()
	assert.NoError(t, err)
	assert.NotNil(t, sm)

	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/test",
		Value: []byte("value"),
	})
	sm.OnCreate("/test", []byte("value"))

	stateMgr.EXPECT().EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/test",
	})
	sm.OnDelete("/test")
}
