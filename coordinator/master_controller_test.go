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
	"github.com/lindb/lindb/coordinator/elect"
	masterpkg "github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
)

func TestNewMasterController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dis := discovery.NewMockDiscovery(ctrl)
	disFct := discovery.NewMockFactory(ctrl)
	disFct.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(dis).AnyTimes()
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "create master controller failure",
			prepare: func() {
				dis.EXPECT().Discovery(true).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create master controller failure",
			prepare: func() {
				dis.EXPECT().Discovery(true).Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cfg := &MasterCfg{
				Ctx:              context.TODO(),
				DiscoveryFactory: disFct,
			}
			if tt.prepare != nil {
				tt.prepare()
			}

			mc, err := NewMasterController(cfg)
			if ((err != nil) != tt.wantErr && mc == nil) || (!tt.wantErr && mc == nil) {
				t.Errorf("NewMasterController() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMasterController_OnFailOver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newStateMgrFn = masterpkg.NewStateManager
		ctrl.Finish()
	}()

	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discovery1.EXPECT().Close().AnyTimes()
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	stateMgr := masterpkg.NewMockStateManager(ctrl)
	stateMgr.EXPECT().Close().AnyTimes()
	stateMgr.EXPECT().SetStateMachineFactory(gomock.Any()).AnyTimes()
	newStateMgrFn = func(ctx context.Context, masterRepo state.Repository,
		repoFactory state.RepositoryFactory) masterpkg.StateManager {
		return stateMgr
	}
	registry := discovery.NewMockRegistry(ctrl)

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "start state machine failure",
			prepare: func() {
				discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "register master done failure",
			prepare: func() {
				discovery1.EXPECT().Discovery(gomock.Any()).Return(nil).MaxTimes(3)
				registry.EXPECT().Register(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "elect master successfully",
			prepare: func() {
				discovery1.EXPECT().Discovery(gomock.Any()).Return(nil).MaxTimes(3)
				registry.EXPECT().Register(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mc := &masterController{
				ctx: context.TODO(),
				cfg: &MasterCfg{
					DiscoveryFactory: discoveryFactory,
				},
				registry:   registry,
				statistics: metrics.NewMasterStatistics(),
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			err := mc.OnFailOver()

			if (err != nil) != tt.wantErr {
				t.Errorf("OnFailOver() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				assert.NotNil(t, mc.GetStateManager())
			}
		})
	}
}

func TestMasterController_OnResignation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := masterpkg.NewMockStateManager(ctrl)
	stateMgr.EXPECT().Close().AnyTimes()
	registry := discovery.NewMockRegistry(ctrl)

	mc := &masterController{
		stateMgr: stateMgr,
		cfg: &MasterCfg{
			Node: &models.StatelessNode{},
		},
		registry:        registry,
		stateMachineFct: masterpkg.NewStateMachineFactory(context.TODO(), nil, stateMgr),
		statistics:      metrics.NewMasterStatistics(),
	}
	// resign failure
	registry.EXPECT().Deregister(gomock.Any()).Return(fmt.Errorf("err"))
	mc.OnResignation()
	// resign successfully
	registry.EXPECT().Deregister(gomock.Any()).Return(nil)
	mc.OnResignation()
}

func TestMasterController_Start_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registry := discovery.NewMockRegistry(ctrl)
	ctx, cancel := context.WithCancel(context.TODO())
	masterElect := elect.NewMockElection(ctrl)
	mc := &masterController{
		ctx:      ctx,
		cancel:   cancel,
		elect:    masterElect,
		registry: registry,
		cfg: &MasterCfg{
			Node: &models.StatelessNode{},
		},
	}
	gomock.InOrder(
		masterElect.EXPECT().Initialize(),
		masterElect.EXPECT().Elect(),
	)
	mc.Start()
	master := &models.Master{}
	masterElect.EXPECT().IsMaster().Return(true)
	masterElect.EXPECT().GetMaster().Return(master)
	assert.True(t, mc.IsMaster())
	assert.Equal(t, master, mc.GetMaster())
	masterElect.EXPECT().Close()
	registry.EXPECT().Close().Return(fmt.Errorf("err"))
	mc.Stop()
}

func TestMasterController_Elect_Listener(t *testing.T) {
	mc := &masterController{}
	// err
	mc.OnCreate("", []byte("err"))
	master1 := &models.Master{}
	data := encoding.JSONMarshal(master1)
	c := 0
	mc.WatchMasterElected(func(master *models.Master) {
		assert.Equal(t, master1, master)
		c++
	})
	mc.OnCreate("", data)
	assert.Equal(t, 1, c)

	mc.OnDelete("")
}

func TestMasterController_FlushDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	masterElect := elect.NewMockElection(ctrl)
	stateMgr := masterpkg.NewMockStateManager(ctrl)
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "isn't master",
			prepare: func() {
				masterElect.EXPECT().IsMaster().Return(false)
			},
			wantErr: false,
		},
		{
			name: "storage not found",
			prepare: func() {
				masterElect.EXPECT().IsMaster().Return(true)
				stateMgr.EXPECT().GetStorageCluster("test").Return(nil)
			},
			wantErr: true,
		},
		{
			name: "flush database successfully",
			prepare: func() {
				masterElect.EXPECT().IsMaster().Return(true)
				storage := masterpkg.NewMockStorageCluster(ctrl)
				stateMgr.EXPECT().GetStorageCluster("test").Return(storage)
				storage.EXPECT().FlushDatabase("db").Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mc := &masterController{
				elect:    masterElect,
				stateMgr: stateMgr,
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			err := mc.FlushDatabase("test", "db")
			if (err != nil) != tt.wantErr {
				t.Errorf("FlushDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
