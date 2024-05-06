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
	"time"

	"github.com/lindb/common/pkg/encoding"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

func TestStateManager_Close(t *testing.T) {
	mgr := NewStateManager(context.TODO(), nil, &models.StatefulNode{}, nil)
	mgr.Close()
}

func TestStateManager_Handle_Event_Panic(t *testing.T) {
	mgr := NewStateManager(context.TODO(), nil, &models.StatefulNode{ID: 1}, nil)
	// case 1: panic
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.ShardAssignmentChanged,
		Key:  "/shard/assign/test",
		Value: encoding.JSONMarshal(&models.ShardAssignment{
			Name:   "test",
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3}}},
		}),
	})
	time.Sleep(100 * time.Millisecond)
	mgr.Close()
}

func TestStateManager_Node(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		getConnFct = rpc.GetStorageClientConnFactory
		ctrl.Finish()
	}()

	conFct := rpc.NewMockClientConnFactory(ctrl)
	getConnFct = func() rpc.ClientConnFactory {
		return conFct
	}

	c := 0
	mgr := NewStateManager(context.TODO(), nil, &models.StatefulNode{ID: 1}, nil)
	// test register nil event handler
	mgr.WatchNodeStateChangeEvent(models.NodeID(1), nil)
	mgr.WatchNodeStateChangeEvent(models.NodeID(1), func(_ models.NodeStateType) {
		c++
	})
	// case 1: unmarshal node info err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/lives/1.1.1.1:9000",
		Value: []byte("221"),
	})
	// case 2: cache node
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeStartup,
		Key:  "/lives/1",
		Value: encoding.JSONMarshal(&models.StatefulNode{ID: 1, StatelessNode: models.StatelessNode{
			HostIP: "1.1.1.1",
		}}),
	})
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeStartup,
		Key:  "/lives/4",
		Value: encoding.JSONMarshal(&models.StatefulNode{ID: 4, StatelessNode: models.StatelessNode{
			HostIP: "1.1.1.4",
		}}),
	})
	time.Sleep(100 * time.Millisecond) // wait
	node, ok := mgr.GetLiveNode(models.NodeID(1))
	assert.True(t, ok)
	assert.Equal(t, models.StatefulNode{ID: 1, StatelessNode: models.StatelessNode{
		HostIP: "1.1.1.1",
	}}, node)
	assert.Len(t, mgr.GetLiveNodes(), 2)

	// case 4: remove not exist node
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/lives/2",
	})
	// case 5: remove node
	conFct.EXPECT().CloseClientConn(gomock.Any()).Return(fmt.Errorf("err"))
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/lives/1",
	})
	conFct.EXPECT().CloseClientConn(gomock.Any()).Return(nil)
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/lives/4",
	})
	// case 6: remove node, node id err
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/lives/wrong_id",
	})
	time.Sleep(100 * time.Millisecond) // wait

	node, ok = mgr.GetLiveNode(models.NodeID(1))
	assert.False(t, ok)
	assert.Equal(t, models.StatefulNode{}, node)

	assert.True(t, c > 0)

	mgr.Close()
}

func TestStateManager_OnShardAssignment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	engine := tsdb.NewMockEngine(ctrl)
	mgr := NewStateManager(context.TODO(), repo, &models.StatefulNode{ID: 1}, engine)
	// case 1: create shard storage engine err
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("{}"), nil)
	engine.EXPECT().CreateShards(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.ShardAssignmentChanged,
		Key:  "/shard/assign/test",
		Value: encoding.JSONMarshal(&models.ShardAssignment{
			Name:   "test",
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3}}},
		}),
	})
	// case 2: get db config err
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.ShardAssignmentChanged,
		Key:  "/shard/assign/test",
		Value: encoding.JSONMarshal(&models.ShardAssignment{
			Name:   "test",
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3}}},
		}),
	})
	// case 2: unmarshal db config err
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("abc"), nil)
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.ShardAssignmentChanged,
		Key:  "/shard/assign/test",
		Value: encoding.JSONMarshal(&models.ShardAssignment{
			Name:   "test",
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3}}},
		}),
	})
	// case 1: create shard storage engine successfully
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("{}"), nil)
	engine.EXPECT().CreateShards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.ShardAssignmentChanged,
		Key:  "/shard/assign/test",
		Value: encoding.JSONMarshal(&models.ShardAssignment{Name: "test",
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3}}},
		}),
	})
	// case 1: unmarshal shard assign err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.ShardAssignmentChanged,
		Key:   "/shard/assign/test",
		Value: []byte("xx"),
	})
	// case 2: shard assignment is nil
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.ShardAssignmentChanged,
		Key:   "/shard/assign/test",
		Value: encoding.JSONMarshal(&models.ShardAssignment{}),
	})
	// case 3: other replica
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.ShardAssignmentChanged,
		Key:  "/shard/assign/test",
		Value: encoding.JSONMarshal(&models.ShardAssignment{
			Name:   "test",
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{2, 3}}},
		}),
	})
	time.Sleep(100 * time.Millisecond)
	assert.Len(t, mgr.GetShardAssignments(), 1)
	mgr.Close()
}

func TestStateManager_onDatabaseLimits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	engine := tsdb.NewMockEngine(ctrl)
	mgr := NewStateManager(context.TODO(), nil, nil, engine)

	// case 1: decode limit failure
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseLimitsChanged,
		Key:   "/database/limit/db2",
		Value: []byte("dd"),
	})
	engine.EXPECT().SetDatabaseLimits("db2", gomock.Any())
	// case 1: set limits
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseLimitsChanged,
		Key:   "/database/limit/db2",
		Value: []byte(models.NewDefaultLimits().TOML()),
	})
	time.Sleep(100 * time.Millisecond)
}
