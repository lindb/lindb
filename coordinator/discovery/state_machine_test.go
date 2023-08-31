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

package discovery

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestStateMachineType_String(t *testing.T) {
	assert.Equal(t, DatabaseConfigStateMachine.String(), "DatabaseConfigStateMachine")
	assert.Equal(t, ShardAssignmentStateMachine.String(), "ShardAssignmentStateMachine")
	assert.Equal(t, LiveNodeStateMachine.String(), "LiveNodeStateMachine")
	assert.Equal(t, StorageStatusStateMachine.String(), "StorageStatusStateMachine")
	assert.Equal(t, StorageConfigStateMachine.String(), "StorageConfigStateMachine")
	assert.Equal(t, StorageNodeStateMachine.String(), "StorageNodeStateMachine")
	assert.Equal(t, (StateMachineType(0)).String(), "Unknown")
	assert.Equal(t, BrokerConfigStateMachine.String(), "BrokerConfigStateMachine")
	assert.Equal(t, BrokerNodeStateMachine.String(), "BrokerNodeStateMachine")
	assert.Equal(t, DatabaseLimitsStateMachine.String(), "DatabaseLimitsStateMachine")
}

func TestNewMockStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	discoveryFct := NewMockFactory(ctrl)
	discovery := NewMockDiscovery(ctrl)
	discoveryFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery).AnyTimes()

	// case 1: new err
	discovery.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	sm, err := NewStateMachine(context.TODO(),
		DatabaseConfigStateMachine, discoveryFct,
		"/test", true, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, sm)
	// case 2: new ok
	discovery.EXPECT().Discovery(gomock.Any()).Return(nil)
	sm, err = NewStateMachine(context.TODO(),
		DatabaseConfigStateMachine, discoveryFct,
		"/test", true, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, sm)
}

func TestStateMachine_OnCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sm := newStataMachine(t, ctrl)
	sm1 := sm.(*stateMachine)
	// case 1: create fn is nil
	sm.OnCreate("/test", nil)
	// case 2: test normal case
	flag := false
	sm1.onCreateFn = func(key string, _ []byte) {
		flag = true
		assert.Equal(t, "/test", key)
	}
	sm.OnCreate("/test", nil)
	assert.True(t, flag)
	// case 3: state machine is not running
	sm1.running.Store(false)
	sm1.onCreateFn = func(_ string, _ []byte) {
		flag = false
	}
	sm.OnCreate("/test", nil)
	assert.True(t, flag)
}

func TestStateMachine_OnDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sm := newStataMachine(t, ctrl)
	sm1 := sm.(*stateMachine)
	// case 1: delete fn is nil
	sm.OnDelete("/test")
	// case 2: test normal case
	flag := false
	sm1.onDeleteFn = func(key string) {
		flag = true
		assert.Equal(t, "/test", key)
	}
	sm.OnDelete("/test")
	assert.True(t, flag)
	// case 3: state machine is not running
	sm1.running.Store(false)
	sm1.onDeleteFn = func(_ string) {
		flag = false
	}
	sm.OnDelete("/test")
	assert.True(t, flag)
}

func TestStateMachine_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	discovery := NewMockDiscovery(ctrl)
	discovery.EXPECT().Close().AnyTimes()

	sm := newStataMachine(t, ctrl)
	sm1 := sm.(*stateMachine)
	sm1.discovery = discovery
	err := sm.Close()
	assert.NoError(t, err)
	// state machine is not running, do nothing
	err = sm.Close()
	assert.NoError(t, err)
}

func TestStateMachine_ExploreData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, fn func(key, value []byte)) error {
			fn([]byte("key"), encoding.JSONMarshal(&models.Master{ElectTime: 11}))
			return nil
		})
	rs, err := ExploreData(context.TODO(), repo, models.StateMachineInfo{
		Path: "/test/a",
		CreateState: func() interface{} {
			return &models.Master{}
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, rs)

	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, fn func(key, value []byte)) error {
			fn([]byte("key"), []byte("value"))
			return nil
		})
	rs, err = ExploreData(context.TODO(), repo, models.StateMachineInfo{
		Path: "/test/a",
		CreateState: func() interface{} {
			return &models.Master{}
		},
	})
	assert.NoError(t, err)
	assert.Empty(t, rs)

	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	rs, err = ExploreData(context.TODO(), repo, models.StateMachineInfo{})
	assert.Error(t, err)
	assert.Nil(t, rs)
}

func newStataMachine(t *testing.T, ctrl *gomock.Controller) StateMachine {
	discoveryFct := NewMockFactory(ctrl)
	discovery := NewMockDiscovery(ctrl)
	discoveryFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery).AnyTimes()

	discovery.EXPECT().Discovery(gomock.Any()).Return(nil)
	sm, err := NewStateMachine(context.TODO(),
		DatabaseConfigStateMachine, discoveryFct,
		"/test", true, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, sm)
	return sm
}
