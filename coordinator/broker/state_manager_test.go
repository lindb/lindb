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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
)

func TestStateManager_Close(t *testing.T) {
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, nil, nil, nil)
	mgr.Close()
}

func TestStateManager_Handle_Event_Panic(t *testing.T) {
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, nil, nil, nil)
	// case 1: panic
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/1.1.1.1:9000",
	})
	time.Sleep(100 * time.Millisecond)
	mgr.Close()
}

func TestStateManager_DatabaseConfig(t *testing.T) {
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, nil, nil, nil)
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
	mgr := NewStateManager(context.TODO(), models.StatelessNode{HostIP: "3.3.3.3"}, cm, nil, nil)
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
		Value: []byte(`{"HostIp":"1.1.1.1"}`),
	})
	time.Sleep(time.Second) // wait
	nodes := mgr.GetLiveNodes()
	assert.Equal(t, []models.StatelessNode{{HostIP: "1.1.1.1"}}, nodes)

	// case 4: remove not exist node
	cm.EXPECT().CloseConnection("2.2.2.2:9000")
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/lives/2.2.2.2:9000",
	})
	// case 5: remove node
	cm.EXPECT().CloseConnection("1.1.1.1:9000")
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
	cm := replica.NewMockChannelManager(ctrl)
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, connectionMgr, nil, cm)
	// case 1: unmarshal storage state err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.StorageStateChanged,
		Key:   "/lin/storage",
		Value: []byte("221"),
	})
	// case 2: storage name is empty
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.StorageStateChanged,
		Key:   "/lin/storage",
		Value: []byte("{}"),
	})
	// case 3: new storage state
	connectionMgr.EXPECT().CreateConnection(gomock.Any()).MaxTimes(2)
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  "/lin/storage/test",
		Value: encoding.JSONMarshal(&models.StorageState{
			Name: "test",
			LiveNodes: map[models.NodeID]models.StatefulNode{1: {
				StatelessNode: models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
			}, 2: {
				StatelessNode: models.StatelessNode{HostIP: "2.2.2.2", GRPCPort: 9000},
			}},
		}),
	})
	// case 4: old storage state, new node online, old node offline
	connectionMgr.EXPECT().CreateConnection(gomock.Any()).MaxTimes(2)
	connectionMgr.EXPECT().CloseConnection("2.2.2.2:9000")
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  "/lin/storage/test",
		Value: encoding.JSONMarshal(&models.StorageState{
			Name: "test",
			LiveNodes: map[models.NodeID]models.StatefulNode{1: {
				StatelessNode: models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
			}, 3: {
				StatelessNode: models.StatelessNode{HostIP: "3.3.3.3", GRPCPort: 9000},
			}},
		}),
	})
	// case 5: remove storage
	connectionMgr.EXPECT().CloseConnection("1.1.1.1:9000")
	connectionMgr.EXPECT().CloseConnection("3.3.3.3:9000")
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageDeletion,
		Key:  "/lin/storage/test",
	})
	// case 6: remove not exist storage
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageDeletion,
		Key:  "/lin/storage/test",
	})
	time.Sleep(time.Second)
}

func TestStateManager_ShardState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	connectionMgr := rpc.NewMockConnectionManager(ctrl)
	cm := replica.NewMockChannelManager(ctrl)
	channel := replica.NewMockChannel(ctrl)
	mgr := NewStateManager(context.TODO(), models.StatelessNode{}, connectionMgr, nil, cm)

	connectionMgr.EXPECT().CreateConnection(gomock.Any()).MaxTimes(2)

	cm.EXPECT().CreateChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	cm.EXPECT().CreateChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(channel, nil)
	channel.EXPECT().SyncShardState(gomock.Any(), gomock.Any())
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  "/lin/storage/test",
		Value: encoding.JSONMarshal(&models.StorageState{
			Name: "test",
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
	cm.EXPECT().CreateChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")).AnyTimes()
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  "/lin/storage/test_1",
		Value: encoding.JSONMarshal(&models.StorageState{
			Name: "test_1",
			ShardStates: map[string]map[models.ShardID]models.ShardState{
				"test_1": {1: models.ShardState{ID: 1, State: models.OnlineShard}, 2: models.ShardState{ID: 2}},
			},
		}),
	})
	connectionMgr.EXPECT().CreateConnection(gomock.Any()).MaxTimes(2)
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageStateChanged,
		Key:  "/lin/storage/test_2",
		Value: encoding.JSONMarshal(&models.StorageState{
			Name: "test_2",
			ShardStates: map[string]map[models.ShardID]models.ShardState{
				"test_2": {},
			},
			LiveNodes: map[models.NodeID]models.StatefulNode{1: {
				StatelessNode: models.StatelessNode{HostIP: "3.1.1.1", GRPCPort: 9000},
			}, 2: {
				StatelessNode: models.StatelessNode{HostIP: "3.2.2.2", GRPCPort: 9000},
			}},
		}),
	})
	time.Sleep(time.Second)

	mgr1 := mgr.(*stateManager)
	mgr1.mutex.Lock()
	mgr1.databases = map[string]models.Database{
		"test_1": {Storage: "test_1"},
		"test_2": {Storage: "test_2"},
		"test":   {Storage: "test_not_exist"},
		"db":     {Storage: "test"}}
	mgr1.mutex.Unlock()

	// db not exist
	replicas, err := mgr.GetQueryableReplicas("test_db")
	assert.Equal(t, err, constants.ErrDatabaseNotFound)
	assert.Empty(t, replicas)

	// storage not exist
	replicas, err = mgr.GetQueryableReplicas("test")
	assert.Equal(t, err, constants.ErrNoStorageCluster)
	assert.Empty(t, replicas)
	// no live nodes
	replicas, err = mgr.GetQueryableReplicas("test_1")
	assert.Equal(t, err, constants.ErrNoLiveNode)
	assert.Empty(t, replicas)
	// no shard
	replicas, err = mgr.GetQueryableReplicas("test_2")
	assert.Equal(t, err, constants.ErrShardNotFound)
	assert.Empty(t, replicas)

	replicas, err = mgr.GetQueryableReplicas("db")
	assert.NoError(t, err)
	assert.Len(t, replicas, 1)
}
