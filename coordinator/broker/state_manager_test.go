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
	"testing"
	"time"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
)

func TestStateManager_Close(t *testing.T) {
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, nil, nil)
	mgr.Close()
}

func TestStateManager_Handle_Event_Panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	connectionMgr := rpc.NewMockConnectionManager(ctrl)
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, connectionMgr, nil)
	mgr1 := mgr.(*stateManager)
	mgr1.mutex.Lock()
	mgr1.nodes["1.1.1.1:9000"] = models.StatelessNode{}
	mgr1.mutex.Unlock()
	connectionMgr.EXPECT().CloseConnection(gomock.Any()).Do(func(node models.Node) {
		panic("err")
	})
	// case 1: panic
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/1.1.1.1:9000",
	})
	time.Sleep(100 * time.Millisecond)
	mgr.Close()
}

func TestStateManager_DatabaseConfig(t *testing.T) {
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, nil, nil)
	// case 1: unmarshal database config err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/test",
		Value: []byte("221"),
	})
	// case 2: database id empty
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/test",
		Value: []byte("{}"),
	})
	// case 3: cache database config
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/test",
		Value: []byte(`{"name":"test"}`),
	})
	time.Sleep(time.Second) // wait
	databaseCfg, ok := mgr.GetDatabaseCfg("test")
	assert.True(t, ok)
	assert.Equal(t, models.Database{Name: "test"}, databaseCfg)
	assert.Len(t, mgr.GetDatabases(), 1)

	// case 4: remove not exist database config
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/test_not_exist",
	})
	// case 5: remove database config
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/test",
	})
	time.Sleep(time.Second) // wait
	_, ok = mgr.GetDatabaseCfg("test")
	assert.False(t, ok)

	mgr.Close()
}

func TestStateManager_Node(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := rpc.NewMockConnectionManager(ctrl)
	mgr := NewStateManager(context.TODO(), models.StatelessNode{HostIP: "3.3.3.3"}, cm, nil)
	// case 1: unmarshal node info err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/lives/1.1.1.1:9000",
		Value: []byte("221"),
	})
	// case 2: cache node
	cm.EXPECT().CreateConnection(gomock.Any())
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/lives/1.1.1.1:9000",
		Value: []byte(`{"HostIp":"1.1.1.1","GRPCPort":9000}`),
	})
	time.Sleep(time.Second) // wait
	nodes := mgr.GetLiveNodes()
	assert.Equal(t, []models.StatelessNode{{HostIP: "1.1.1.1", GRPCPort: 9000}}, nodes)

	// case 4: remove not exist node
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/lives/2.2.2.2:9000",
	})
	// case 5: remove node
	cm.EXPECT().CloseConnection(&models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000})
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/lives/1.1.1.1:9000",
	})
	time.Sleep(time.Second) // wait
	nodes = mgr.GetLiveNodes()
	assert.Empty(t, nodes)

	assert.Equal(t, models.StatelessNode{HostIP: "3.3.3.3"}, mgr.GetCurrentNode())

	mgr.Close()
}

func TestStateManager_Storage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	connectionMgr := rpc.NewMockConnectionManager(ctrl)
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, connectionMgr, nil)

	// case 1: unmarshal storage state err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.StorageStateChanged,
		Key:   constants.StorageStatePath,
		Value: []byte("221"),
	})
	// case 2: storage name is empty
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.StorageStateChanged,
		Key:   constants.StorageStatePath,
		Value: []byte("{}"),
	})
	// case 3: new storage state
	connectionMgr.EXPECT().CreateConnection(gomock.Any()).MaxTimes(2)
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  constants.StorageStatePath,
		Value: encoding.JSONMarshal(&models.StorageState{
			LiveNodes: map[models.NodeID]models.StatefulNode{1: {
				StatelessNode: models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
			}, 2: {
				StatelessNode: models.StatelessNode{HostIP: "2.2.2.2", GRPCPort: 9000},
			}},
		}),
	})
	// case 4: old storage state, new node online, old node offline
	connectionMgr.EXPECT().CreateConnection(gomock.Any()).MaxTimes(2)
	connectionMgr.EXPECT().CloseConnection(&models.StatefulNode{
		StatelessNode: models.StatelessNode{HostIP: "2.2.2.2", GRPCPort: 9000},
	})
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  constants.StorageStatePath,
		Value: encoding.JSONMarshal(&models.StorageState{
			LiveNodes: map[models.NodeID]models.StatefulNode{1: {
				StatelessNode: models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
			}, 3: {
				StatelessNode: models.StatelessNode{HostIP: "3.3.3.3", GRPCPort: 9000},
			}},
		}),
	})
	time.Sleep(time.Second)
	state := mgr.GetStorage()
	assert.NotNil(t, state)
}

func TestStateManager_ShardState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	connectionMgr := rpc.NewMockConnectionManager(ctrl)
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, connectionMgr, nil)
	c := 0
	mgr.WatchShardStateChangeEvent(func(_ models.Database,
		_ map[models.ShardID]models.ShardState,
		_ map[models.NodeID]models.StatefulNode,
	) {
		c++
	})
	connectionMgr.EXPECT().CreateConnection(gomock.Any()).MaxTimes(2)

	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  constants.StorageStatePath,
		Value: encoding.JSONMarshal(&models.StorageState{
			ShardStates: map[string]map[models.ShardID]models.ShardState{
				"db": {1: models.ShardState{ID: 1, State: models.OnlineShard}, 2: models.ShardState{ID: 2}},
			},
			LiveNodes: map[models.NodeID]models.StatefulNode{1: {
				StatelessNode: models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
			}, 2: {
				StatelessNode: models.StatelessNode{HostIP: "2.2.2.2", GRPCPort: 9000},
			}},
		}),
	})
	mgr1 := mgr.(*stateManager)
	mgr1.mutex.Lock()
	mgr1.databases = map[string]models.Database{
		"test_1": {},
		"test_2": {},
		"test":   {},
		"db":     {},
	}
	mgr1.mutex.Unlock()
	time.Sleep(200 * time.Millisecond)

	// db not exist
	replicas, err := mgr.GetQueryableReplicas("test_db")
	assert.Equal(t, err, constants.ErrDatabaseNotFound)
	assert.Empty(t, replicas)

	// no shard
	replicas, err = mgr.GetQueryableReplicas("test_2")
	assert.Equal(t, err, constants.ErrShardNotFound)
	assert.Empty(t, replicas)

	replicas, err = mgr.GetQueryableReplicas("db")
	assert.NoError(t, err)
	assert.Len(t, replicas, 1)
	assert.True(t, c > 0)

	connectionMgr.EXPECT().CloseConnection(gomock.Any()).MaxTimes(2)
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  constants.StorageStatePath,
		Value: encoding.JSONMarshal(&models.StorageState{
			ShardStates: map[string]map[models.ShardID]models.ShardState{
				"test_1": {1: models.ShardState{ID: 1, State: models.OnlineShard}, 2: models.ShardState{ID: 2}},
			},
		}),
	})
	time.Sleep(time.Second)

	// no live node
	replicas, err = mgr.GetQueryableReplicas("test_1")
	assert.Equal(t, err, constants.ErrNoLiveNode)
	assert.Empty(t, replicas)
}

func TestStateManager_GetLiveNodes(t *testing.T) {
	s := &stateManager{
		nodes: make(map[string]models.StatelessNode),
	}
	nodes := s.GetLiveNodes()
	assert.Empty(t, nodes)

	s.nodes["test1"] = models.StatelessNode{HostIP: "1.1.1.1"}
	s.nodes["test2"] = models.StatelessNode{HostIP: "1.1.2.1"}
	nodes = s.GetLiveNodes()
	assert.Equal(t, []models.StatelessNode{{HostIP: "1.1.1.1"}, {HostIP: "1.1.2.1"}}, nodes)
}

func TestStateManager_Choose(t *testing.T) {
	mgr := &stateManager{
		nodes: map[string]models.StatelessNode{"test": {}},
		databases: map[string]models.Database{
			"test_1": {},
		},
		storageState: &models.StorageState{
			LiveNodes: map[models.NodeID]models.StatefulNode{
				1: {StatelessNode: models.StatelessNode{HostIP: "1.1.1.1"}},
				2: {StatelessNode: models.StatelessNode{HostIP: "1.1.1.2"}},
			},
			ShardStates: map[string]map[models.ShardID]models.ShardState{
				"test_1": {
					1: {
						State:  models.OnlineShard,
						Leader: models.NodeID(1),
					},
					2: {
						State:  models.OnlineShard,
						Leader: models.NodeID(2),
					},
				},
			},
		},
		logger: logger.GetLogger("Test", "StateManager"),
	}
	plans, err := mgr.Choose("test", 2)
	assert.Error(t, err)
	assert.Nil(t, plans)

	plans, err = mgr.Choose("test_2", 2)
	assert.Error(t, err)
	assert.Nil(t, plans)
	plans, err = mgr.Choose("test_1", 2)
	assert.NoError(t, err)
	assert.Len(t, plans, 1)
	plans, err = mgr.Choose("test_1", 1)
	assert.NoError(t, err)
	assert.Len(t, plans, 1)
}

func TestStateManager_onDatabaseLimits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, nil, nil)

	// case 1: decode limit failure
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseLimitsChanged,
		Key:   "/database/limit/db2",
		Value: []byte("dd"),
	})
	// case 1: set limits
	limit2 := models.NewDefaultLimits()
	limit2.MaxFieldsPerMetric = 10000
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseLimitsChanged,
		Key:   "/database/limit/db2",
		Value: []byte(limit2.TOML()),
	})
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, limit2, models.GetDatabaseLimits("db2"))
	assert.Equal(t, models.NewDefaultLimits(), models.GetDatabaseLimits("test"))
}
