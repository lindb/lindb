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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	masterpkg "github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
)

func TestMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(eventCh).AnyTimes()
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil).AnyTimes()
	discovery1.EXPECT().Close().AnyTimes()
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	node1 := models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 8000}
	master1 := NewMasterController(&MasterCfg{
		Ctx:              context.TODO(),
		Repo:             repo,
		Node:             &node1,
		TTL:              1,
		DiscoveryFactory: discoveryFactory,
	})
	master1.Start()
	data := encoding.JSONMarshal(&models.Master{Node: &node1})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.Equal(t, &node1, master1.GetMaster().Node)
	assert.True(t, master1.IsMaster())
	assert.NotNil(t, master1.GetStateManager())

	// re-elect
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeDelete,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.False(t, master1.IsMaster())

	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.True(t, master1.IsMaster())

	master1.Stop()
}

func TestMaster_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(eventCh).AnyTimes()
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discovery1.EXPECT().Close().AnyTimes()
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	node1 := models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 8000}
	master1 := NewMasterController(&MasterCfg{
		Ctx:              context.TODO(),
		Repo:             repo,
		Node:             &node1,
		TTL:              1,
		DiscoveryFactory: discoveryFactory,
	})
	master1.Start()

	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	data := encoding.JSONMarshal(&models.Master{Node: &node1})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.False(t, master1.IsMaster())
	assert.Nil(t, master1.GetMaster())

	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.False(t, master1.IsMaster())
	assert.Nil(t, master1.GetMaster())

	master1.Stop()
}

func TestMaster_FlushDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(eventCh).AnyTimes()
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil).AnyTimes()
	discovery1.EXPECT().Close().AnyTimes()
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	node1 := models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 8000}
	master1 := NewMasterController(&MasterCfg{
		Ctx:              context.TODO(),
		Repo:             repo,
		Node:             &node1,
		TTL:              1,
		DiscoveryFactory: discoveryFactory,
	})
	err := master1.FlushDatabase("test", "test")
	assert.NoError(t, err)

	master1.Start()
	data := encoding.JSONMarshal(&models.Master{Node: &node1})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.True(t, master1.IsMaster())
	err = master1.FlushDatabase("test", "test")
	assert.Error(t, err)

	m1 := master1.(*masterController)
	m1.mutex.Lock()
	statMgr := masterpkg.NewMockStateManager(ctrl)
	m1.stateMgr = statMgr
	m1.mutex.Unlock()

	cluster1 := masterpkg.NewMockStorageCluster(ctrl)
	statMgr.EXPECT().GetStorageCluster(gomock.Any()).Return(cluster1)
	cluster1.EXPECT().FlushDatabase(gomock.Any()).Return(nil)
	err = master1.FlushDatabase("test", "test")
	assert.NoError(t, err)
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(10 * time.Millisecond)
}
