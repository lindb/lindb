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

package master

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
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
)

func TestStateManager_Close(t *testing.T) {
	mgr := NewStateManager(context.TODO(), nil, nil)
	fct := &StateMachineFactory{}
	mgr.SetStateMachineFactory(fct)
	assert.Equal(t, fct, mgr.GetStateMachineFactory())

	mgr.Close()
}

func TestStateManager_DropDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sc := NewMockStorageCluster(ctrl)
	repo := state.NewMockRepository(ctrl)
	mgr := NewStateManager(context.TODO(), nil, nil)
	mgr1 := mgr.(*stateManager)
	mgr1.mutex.Lock()
	shardAssignment := models.NewShardAssignment("test-db")
	mgr1.shardAssignments["test-db"] = shardAssignment
	mgr1.shardAssignments["test-db2"] = shardAssignment
	mgr1.shardAssignments["test-db3"] = shardAssignment
	mgr1.databases["test-db"] = &models.Database{}
	mgr1.databases["test-db2"] = &models.Database{}
	mgr1.databases["test-db3"] = &models.Database{}
	mgr1.storage = sc
	mgr1.masterRepo = repo
	mgr1.mutex.Unlock()

	// case 1: database not exist
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/database/config/test",
	})
	// case 2: drop database config
	sc.EXPECT().GetState().Return(models.NewStorageState()).MaxTimes(2)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	sc.EXPECT().DropDatabaseAssignment(gomock.Any()).Return(fmt.Errorf("err"))
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/database/config/test-db",
	})

	// case 3: sync state failure
	sc.EXPECT().GetState().Return(models.NewStorageState()).MaxTimes(2)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("errj"))
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/database/config/test-db2",
	})

	// case 4: successfully
	sc.EXPECT().GetState().Return(models.NewStorageState()).MaxTimes(2)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	sc.EXPECT().DropDatabaseAssignment(gomock.Any()).Return(nil)
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/database/config/test-db3",
	})

	time.Sleep(100 * time.Millisecond)
	mgr.Close()
}

func TestStateManager_NotRunning(t *testing.T) {
	mgr := NewStateManager(context.TODO(), nil, nil)
	mgr1 := mgr.(*stateManager)
	mgr1.running.Store(false)
	// case 1: not running
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  "/shard/assign/test",
	})
	time.Sleep(100 * time.Millisecond)
	mgr.Close()
}

func TestStateManager_DatabaseCfg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	repo := state.NewMockRepository(ctrl)
	mgr := NewStateManager(context.TODO(), repo, nil)
	mgr1 := mgr.(*stateManager)
	storage1 := NewMockStorageCluster(ctrl)
	mgr1.mutex.Lock()
	mgr1.storage = storage1
	mgr1.mutex.Unlock()

	// case 1: unmarshal cfg err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/database/test",
		Value: []byte("value"),
	})
	// case 2: database name is empty
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/database/test",
		Value: encoding.JSONMarshal(&models.Database{}),
	})
	db := &models.Database{
		Name:          "test",
		NumOfShard:    3,
		ReplicaFactor: 2,
		Option:        &option.DatabaseOption{},
	}
	data := encoding.JSONMarshal(db)
	// case 3: get shard assign err
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/database/test",
		Value: data,
	})
	// case 4: modify shard assign err
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("value"), nil)
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/database/test",
		Value: data,
	})
	// case 4: live node error
	mgr1.mutex.Lock()
	storage1.EXPECT().GetLiveNodes().Return(nil, fmt.Errorf("err"))
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("{}"), nil)
	assert.NoError(t, mgr1.onDatabaseCfgChange("/database/test", data))
	mgr1.mutex.Unlock()
	// case 5: create shard assign err
	mgr1.mutex.Lock()
	storage1.EXPECT().GetLiveNodes().Return(nil, nil)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
	assert.NoError(t, mgr1.onDatabaseCfgChange("/database/test", data))
	mgr1.mutex.Unlock()
	// case 6: trigger modify event
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(encoding.JSONMarshal(&models.ShardAssignment{
		Shards: map[models.ShardID]*models.Replica{1: nil, 2: nil, 3: nil},
	}), nil)
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   "/database/test",
		Value: data,
	})

	time.Sleep(100 * time.Millisecond)
	assert.Len(t, mgr.GetDatabases(), 1)
	mgr.Close()
}

func TestStateManager_ShardAssignment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	repo := state.NewMockRepository(ctrl)
	storage := NewMockStorageCluster(ctrl)
	mgr := NewStateManager(context.TODO(), repo, nil)
	mgr1 := mgr.(*stateManager)
	elector := NewMockReplicaLeaderElector(ctrl)
	mgr1.mutex.Lock()
	mgr1.elector = elector
	mgr1.databases["test"] = &models.Database{}
	mgr1.storage = storage
	mgr1.mutex.Unlock()
	// case 1: unmarshal err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.ShardAssignmentChanged,
		Key:   "/shard/assign/test",
		Value: []byte("valuek"),
	})
	// case 2: put state err
	data := encoding.JSONMarshal(&models.ShardAssignment{
		Name:   "test",
		Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{2, 3}}, 2: {Replicas: []models.NodeID{2, 3}}},
	})
	storage.EXPECT().GetState().Return(models.NewStorageState()).AnyTimes()
	elector.EXPECT().ElectLeader(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.NodeID(2), nil)
	elector.EXPECT().ElectLeader(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.NodeID(0), fmt.Errorf("err"))
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.ShardAssignmentChanged,
		Key:   "/shard/assign/test",
		Value: data,
	})
	// case 2: put state err
	elector.EXPECT().ElectLeader(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.NodeID(2), nil)
	elector.EXPECT().ElectLeader(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.NodeID(0), fmt.Errorf("err"))
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.ShardAssignmentChanged,
		Key:   "/shard/assign/test",
		Value: data,
	})
	time.Sleep(100 * time.Millisecond)
	assert.Len(t, mgr.GetShardAssignments(), 1)
	mgr.Close()
}

func TestStateManager_createShardAssign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	repo := state.NewMockRepository(ctrl)
	storage := NewMockStorageCluster(ctrl)
	mgr := NewStateManager(context.TODO(), repo, nil)
	mgr1 := mgr.(*stateManager)
	// case 1: get live nodes err
	storage.EXPECT().GetLiveNodes().Return(nil, fmt.Errorf("err"))
	shardAssign, err := mgr1.createShardAssignment(storage, &models.Database{Name: "test"}, -1, -1)
	assert.Error(t, err)
	assert.Nil(t, shardAssign)
	// case 2: no live nodes
	storage.EXPECT().GetLiveNodes().Return(nil, nil)
	shardAssign, err = mgr1.createShardAssignment(storage, &models.Database{Name: "test"}, 0, 0)
	assert.Error(t, err)
	assert.Nil(t, shardAssign)
	// case 3: assign shard err
	storage.EXPECT().GetLiveNodes().Return([]models.StatefulNode{{ID: 1}, {ID: 2}, {ID: 3}}, nil).AnyTimes()
	shardAssign, err = mgr1.createShardAssignment(storage, &models.Database{Name: "test"}, -1, -1)
	assert.Error(t, err)
	assert.Nil(t, shardAssign)
	// case 4: save shard assign err
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	shardAssign, err = mgr1.createShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3, ReplicaFactor: 2},
		-1, -1)
	assert.Error(t, err)
	assert.Nil(t, shardAssign)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	storage.EXPECT().SaveDatabaseAssignment(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	shardAssign, err = mgr1.createShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3, ReplicaFactor: 2},
		-1, -1)
	assert.Error(t, err)
	assert.Nil(t, shardAssign)
	// case 5:ok
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	storage.EXPECT().SaveDatabaseAssignment(gomock.Any(), gomock.Any()).Return(nil)
	shardAssign, err = mgr1.createShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3, ReplicaFactor: 2},
		-1, -1)
	assert.NoError(t, err)
	assert.NotNil(t, shardAssign)
}

func TestStateManager_modifyShardAssign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	repo := state.NewMockRepository(ctrl)
	storage := NewMockStorageCluster(ctrl)
	mgr := NewStateManager(context.TODO(), repo, nil)
	mgr1 := mgr.(*stateManager)
	// case 1: no impl
	assert.Panics(t, func() {
		_ = mgr1.modifyShardAssignment(storage,
			&models.Database{Name: "test"},
			&models.ShardAssignment{Shards: map[models.ShardID]*models.Replica{1: {}, 2: {}}})
	})
	// case 2: get live nodes err
	storage.EXPECT().GetLiveNodes().Return(nil, fmt.Errorf("err"))
	err := mgr1.modifyShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3},
		&models.ShardAssignment{Shards: map[models.ShardID]*models.Replica{1: {}, 2: {}}})
	assert.Error(t, err)
	// case 3: no live nodes
	storage.EXPECT().GetLiveNodes().Return(nil, nil)
	err = mgr1.modifyShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3},
		&models.ShardAssignment{Shards: map[models.ShardID]*models.Replica{1: {}, 2: {}}})
	assert.Error(t, err)
	// case 4: modify err
	storage.EXPECT().GetLiveNodes().Return([]models.StatefulNode{{ID: 1}, {ID: 2}, {ID: 3}}, nil).AnyTimes()
	err = mgr1.modifyShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3},
		&models.ShardAssignment{Shards: map[models.ShardID]*models.Replica{1: {}, 2: {}}})
	assert.Error(t, err)

	// case 5: save err
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = mgr1.modifyShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3, ReplicaFactor: 2},
		&models.ShardAssignment{Shards: map[models.ShardID]*models.Replica{1: {}, 2: {}}})
	assert.Error(t, err)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	storage.EXPECT().SaveDatabaseAssignment(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = mgr1.modifyShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3, ReplicaFactor: 2},
		&models.ShardAssignment{Shards: map[models.ShardID]*models.Replica{1: {}, 2: {}}})
	assert.Error(t, err)
	// case 6: ok
	storage.EXPECT().SaveDatabaseAssignment(gomock.Any(), gomock.Any()).Return(nil)
	err = mgr1.modifyShardAssignment(storage,
		&models.Database{Name: "test", NumOfShard: 3, ReplicaFactor: 2},
		&models.ShardAssignment{Shards: map[models.ShardID]*models.Replica{1: {}, 2: {}}})
	assert.NoError(t, err)
}

func TestStateManager_StorageNodeStartup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	repo := state.NewMockRepository(ctrl)
	storage := NewMockStorageCluster(ctrl)
	mgr := NewStateManager(context.TODO(), repo, nil)
	mgr1 := mgr.(*stateManager)
	mgr1.mutex.Lock()
	mgr1.storage = storage
	mgr1.mutex.Unlock()
	// case 1: unmarshal err
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/test/1",
		Value: []byte("dd"),
	})
	// case 2: sync err
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	storage.EXPECT().GetState().Return(models.NewStorageState())
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/test/1",
		Value: []byte(`{"id":1}`),
	})
	// case 3: change shard state,but shard state not found
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	storage.EXPECT().GetState().Return(&models.StorageState{
		LiveNodes: map[models.NodeID]models.StatefulNode{},
		ShardAssignments: map[string]*models.ShardAssignment{"test": {
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3, 4}}},
		}},
	})
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/test/1",
		Value: []byte(`{"id":1}`),
	})
	// case 4: change shard state ok
	storage.EXPECT().GetState().Return(&models.StorageState{
		LiveNodes:   map[models.NodeID]models.StatefulNode{},
		ShardStates: map[string]map[models.ShardID]models.ShardState{"test": {1: {}}},
		ShardAssignments: map[string]*models.ShardAssignment{"test": {
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3, 4}}},
		}},
	})
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   "/test/1",
		Value: []byte(`{"id":1}`),
	})
	time.Sleep(100 * time.Millisecond)

	storage.EXPECT().GetState().Return(&models.StorageState{})
	assert.NotNil(t, mgr.GetStorageState())
	mgr.Close()
}

func TestStateManager_StorageNodeFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	repo := state.NewMockRepository(ctrl)
	storage := NewMockStorageCluster(ctrl)
	mgr := NewStateManager(context.TODO(), repo, nil)
	mgr1 := mgr.(*stateManager)
	mgr1.mutex.Lock()
	mgr1.storage = storage
	mgr1.mutex.Unlock()
	// case 1: unmarshal node id err
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/test/test_1",
	})
	// case 2: sync err
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	storage.EXPECT().GetState().Return(models.NewStorageState())
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/test/1",
	})
	// case 3: change shard state,but elect leader err
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	storage.EXPECT().GetState().Return(&models.StorageState{
		LiveNodes:   map[models.NodeID]models.StatefulNode{},
		ShardStates: map[string]map[models.ShardID]models.ShardState{"test": {1: {Leader: 1}}},
		ShardAssignments: map[string]*models.ShardAssignment{"test": {
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3, 4}}},
		}},
	})
	mgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  "/test/1",
	})
	// case 4: change shard state ok, leader elect success
	shardStates := map[string]map[models.ShardID]models.ShardState{"test": {1: {Leader: 1}}}
	liveNodes := map[models.NodeID]models.StatefulNode{1: {ID: 1}, 2: {ID: 2}}
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	storage.EXPECT().GetState().Return(&models.StorageState{
		LiveNodes:   liveNodes,
		ShardStates: shardStates,
		ShardAssignments: map[string]*models.ShardAssignment{"test": {
			Shards: map[models.ShardID]*models.Replica{1: {Replicas: []models.NodeID{1, 2, 3, 4}}},
		}},
	})
	mgr.EmitEvent(&discovery.Event{
		Type:       discovery.NodeFailure,
		Key:        "/test/1",
		Attributes: map[string]string{},
	})

	time.Sleep(300 * time.Millisecond)
	// get new shard state
	mgr1.mutex.Lock()
	assert.Equal(t, shardStates["test"][1].Leader, models.NodeID(2))
	assert.Len(t, liveNodes, 1)
	assert.Equal(t, liveNodes[models.NodeID(2)].ID, models.NodeID(2))
	mgr1.mutex.Unlock()
	mgr.Close()
}

func TestStateManager_onDatabaseLimits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := state.NewMockRepository(ctrl)
	storage := NewMockStorageCluster(ctrl)
	mgr := NewStateManager(context.TODO(), repo, nil)
	mgr1 := mgr.(*stateManager)
	mgr1.mutex.Lock()
	mgr1.storage = storage
	mgr1.databases["db"] = &models.Database{}
	mgr1.databases["db1"] = &models.Database{}
	mgr1.mutex.Unlock()

	// case 1: database not exist
	mgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseLimitsChanged,
		Key:   "/database/limit/db2",
		Value: []byte("dd"),
	})
	// case 3: set limits failure
	mgr1.mutex.Lock()
	storage.EXPECT().SetDatabaseLimits(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	mgr1.mutex.Unlock()
	assert.Error(t, mgr1.onDatabaseLimitsChange("/database/limit/db", []byte("dd")))
	// case 4: set limits successfully
	mgr1.mutex.Lock()
	storage.EXPECT().SetDatabaseLimits(gomock.Any(), gomock.Any()).Return(nil)
	mgr1.mutex.Unlock()
	assert.NoError(t, mgr1.onDatabaseLimitsChange("/database/limit/db", []byte("dd")))
	time.Sleep(100 * time.Millisecond)
}
