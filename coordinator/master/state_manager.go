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
	"encoding/json"
	"io"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./state_manager.go -destination=./state_manager_mock.go -package=master

type StateManager interface {
	io.Closer

	SetStateMachineFactory(stateMachineFct *StateMachineFactory)
	GetStateMachineFactory() *StateMachineFactory

	// OnDatabaseCfgChange triggers when database create/modify.
	OnDatabaseCfgChange(key string, data []byte)
	// OnDatabaseCfgDelete triggers when database delete.
	OnDatabaseCfgDelete(key string)

	OnShardAssignmentChange(key string, data []byte)
	OnShardAssignmentDelete(key string)

	OnStorageConfigChange(key string, data []byte)
	OnStorageConfigDelete(key string)

	OnStorageNodeStartup(storageName string, key string, data []byte)
	OnStorageNodeFailure(storageName string, key string)

	Register(cfg config.StorageCluster) error
	UnRegister(name string)
	ShardAssignment(databaseCfg models.Database)
	NodeStartup(name string, node models.StatefulNode)
	NodeFailure(name string, nodeID models.NodeID)

	// GetStorageCluster returns cluster controller for maintain the metadata of storage cluster.
	GetStorageCluster(name string) StorageCluster
}

type stateManager struct {
	repoFactory       state.RepositoryFactory
	controllerFactory task.ControllerFactory
	stateMachineFct   *StateMachineFactory

	ctx        context.Context
	masterRepo state.Repository
	elector    ReplicaLeaderElector

	storages  map[string]StorageCluster
	databases map[string]models.Database

	mutex sync.RWMutex

	logger *logger.Logger
}

func NewStateManager(
	ctx context.Context,
	masterRepo state.Repository,
	repoFactory state.RepositoryFactory,
	controllerFactory task.ControllerFactory,
) StateManager {
	return &stateManager{
		ctx:               ctx,
		masterRepo:        masterRepo,
		repoFactory:       repoFactory,
		controllerFactory: controllerFactory,
		storages:          make(map[string]StorageCluster),
		databases:         make(map[string]models.Database),
		elector:           newReplicaLeaderElector(),
		logger:            logger.GetLogger("master", "stateManager"),
	}
}

func (m *stateManager) SetStateMachineFactory(stateMachineFct *StateMachineFactory) {
	m.stateMachineFct = stateMachineFct
}

func (m *stateManager) GetStateMachineFactory() *StateMachineFactory {
	return m.stateMachineFct
}

// OnDatabaseCfgChange triggers when database create/modify.
func (m *stateManager) OnDatabaseCfgChange(key string, data []byte) {
	m.logger.Info("do shard assignment, because database config is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := models.Database{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		m.logger.Error("do shard assignment, because database config is changed, but unmarshal error",
			logger.Error(err))
		return
	}

	m.databases[cfg.Name] = cfg

	m.ShardAssignment(cfg)
}

// OnDatabaseCfgDelete triggers when database delete.
func (m *stateManager) OnDatabaseCfgDelete(key string) {
	panic("need impl")
}

func (m *stateManager) OnShardAssignmentChange(key string, data []byte) {
	m.logger.Info("database's shard assignment is changed",
		logger.String("key", key),
		logger.String("data", string(data)))
	shardAssignment := &models.ShardAssignment{}
	if err := encoding.JSONUnmarshal(data, shardAssignment); err != nil {
		m.logger.Error("database's shard assignment is changed, but unmarshal error",
			logger.Error(err))
		return
	}
	databaseCfg := m.databases[shardAssignment.Name]

	storage := m.storages[databaseCfg.Storage]
	s := storage.GetState()

	liveNodes := s.LiveNodes
	shardStates := make(map[models.ShardID]models.ShardState)
	for shardID, replicas := range shardAssignment.Shards {
		leader, err := m.elector.ElectLeader(shardAssignment, liveNodes, shardID)
		shardState := models.ShardState{Replica: *replicas}
		if err != nil {
			shardState.State = models.OfflineShard
			shardState.Leader = models.NoLeader
			//TODO add log and metric
		} else {
			shardState.State = models.OnlineShard
			shardState.Leader = leader
		}
		shardStates[shardID] = shardState
	}
	//TODO set shard assignments
	s.ShardAssignments[shardAssignment.Name] = shardAssignment
	s.ShardStates[shardAssignment.Name] = shardStates

	m.syncState(s)
}

func (m *stateManager) OnShardAssignmentDelete(key string) {
	panic("need impl")
}

func (m *stateManager) OnStorageConfigChange(key string, data []byte) {
	m.logger.Info("storage config is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := config.StorageCluster{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		m.logger.Error("storage config modified but unmarshal error", logger.Error(err))
		return
	}

	if err := m.Register(cfg); err != nil {
		m.logger.Error("register new storage cluster", logger.Error(err))
		return
	}
}

func (m *stateManager) OnStorageConfigDelete(key string) {
	m.logger.Info("storage config deleted",
		logger.String("key", key))

	_, name := filepath.Split(key)

	m.UnRegister(name)
}

func (m *stateManager) OnStorageNodeStartup(storageName string, key string, data []byte) {
	m.logger.Info("new storage node online in storage cluster",
		logger.String("storage", storageName),
		logger.String("key", key),
		logger.String("data", string(data)))

	node := models.StatefulNode{}
	if err := json.Unmarshal(data, &node); err != nil {
		m.logger.Error("new storage node online in storage cluster but unmarshal error", logger.Error(err))
		return
	}

	m.NodeStartup(storageName, node)
}

func (m *stateManager) OnStorageNodeFailure(storageName string, key string) {
	m.logger.Info("a storage node offline in storage cluster",
		logger.String("storage", storageName),
		logger.String("key", key))

	_, nodeIDStr := filepath.Split(key)
	id, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		m.logger.Error("parse offline node id err", logger.Error(err))
		return
	}

	m.NodeFailure(storageName, models.NodeID(id))
}

func (m *stateManager) Register(cfg config.StorageCluster) error {
	if len(cfg.Name) == 0 {
		return constants.ErrNameEmpty
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// shutdown old storageCluster state machine if exist
	m.UnRegister(cfg.Name)

	cluster, err := newStorageCluster(m.ctx, cfg, m, m.repoFactory, m.controllerFactory)
	if err != nil {
		return err
	}
	m.storages[cfg.Name] = cluster
	// start storage cluster state machine
	if err := cluster.Start(); err != nil {
		m.UnRegister(cfg.Name)
		m.logger.Info("start storage cluster failure", logger.String("storage", cfg.Name), logger.Error(err))
		return err
	}
	return nil
}

// deleteCluster deletes the storageCluster if exist
func (m *stateManager) UnRegister(name string) {
	cluster, ok := m.storages[name]
	if ok {
		// need cleanup storageCluster resource
		cluster.Close()
		delete(m.storages, name)
	}
}

func (m *stateManager) Close() error {
	for name := range m.storages {
		m.UnRegister(name)
	}
	return nil
}

func (m *stateManager) ShardAssignment(databaseCfg models.Database) {
	if len(databaseCfg.Name) == 0 {
		m.logger.Error("database name cannot be empty")
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// get shard assignment from repo, maybe mem state is not sync.
	shardAssign, err := m.GetShardAssign(databaseCfg.Name)
	if err != nil && err != state.ErrNotExist {
		m.logger.Error("get shard assign error", logger.Error(err))
		return
	}

	cluster := m.storages[databaseCfg.Storage]
	s := cluster.GetState()

	// build shard assignment for creation database, generate related coordinator task
	if shardAssign == nil {
		shardAssign, err := m.createShardAssignment(cluster, &databaseCfg, -1, -1)
		if err != nil {
			m.logger.Error("create shard assignment error",
				logger.String("storage", databaseCfg.Storage),
				logger.Any("databaseCfg", databaseCfg),
				logger.Error(err))
			return
		}

		// add new shard assignments
		s.ShardAssignments[databaseCfg.Name] = shardAssign //TODO remove it???
	} else if len(shardAssign.Shards) != databaseCfg.NumOfShard {
		if err := m.modifyShardAssignment(cluster, &databaseCfg, shardAssign); err != nil {
			m.logger.Error("modify shard assignment error",
				logger.String("storage", databaseCfg.Storage),
				logger.Any("databaseCfg", databaseCfg),
				logger.Error(err))
			return
		}
	}
}

func (m *stateManager) NodeStartup(name string, node models.StatefulNode) {
	cluster := m.storages[name]
	s := cluster.GetState()

	s.NodeOnline(node)
	m.onNodeStartup(s, node)
}

func (m *stateManager) NodeFailure(name string, nodeID models.NodeID) {
	cluster := m.storages[name]
	s := cluster.GetState()

	m.onNodeFailure(s, nodeID)
}

func (m *stateManager) onNodeStartup(state *models.StorageState, node models.StatefulNode) {
	// 1. do when a new node come up is send it the entire list of shards that it is supposed to host.
	replicasOnOnlineNode := state.ReplicasOnNode(node.ID)
	for db, shards := range replicasOnOnlineNode {
		shardStates := state.ShardStates[db]
		for _, shardID := range shards {
			shardState := shardStates[shardID]
			if shardState.State == models.OfflineShard {
				shardState.State = models.OnlineShard
				shardState.Leader = node.ID
			}
			shardStates[shardID] = shardState
		}
	}
}

func (m *stateManager) onNodeFailure(state *models.StorageState, nodeID models.NodeID) {
	// 1. find all leaders on failure node, need do leader elect
	leadersOnOfflineNode := state.LeadersOnNode(nodeID)

	liveNodes := state.LiveNodes
	for db, shards := range leadersOnOfflineNode {
		shardAssignment := state.ShardAssignments[db]
		shardStates := state.ShardStates[db]
		for _, shardID := range shards {
			leader, err := m.elector.ElectLeader(shardAssignment, liveNodes, shardID)
			shardState := shardStates[shardID]
			if err != nil {
				shardState.State = models.OfflineShard
				shardState.Leader = models.NoLeader
				//TODO add log and metric
			} else {
				shardState.State = models.OnlineShard
				shardState.Leader = leader
			}
			shardStates[shardID] = shardState
		}
	}
}

func (m *stateManager) syncState(state *models.StorageState) {
	//TODO add timeout
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	data := encoding.JSONMarshal(state)
	if err := m.masterRepo.Put(ctx, constants.GetStorageStatePath(state.Name), data); err != nil {
		//TODO add log????
		m.logger.Error("sync storage state error", logger.String("storage", state.Name), logger.Error(err))
		return
	}
	m.logger.Info("sync storage state successfully", logger.String("storage", state.Name))
}

// createShardAssignment creates shard assignment for spec storageCluster
// 1) generate shard assignment
// 2) save shard assignment into related storage storageCluster
// 3) submit create shard coordinator task(storage node will execute it when receive task event)
func (m *stateManager) createShardAssignment(
	cluster StorageCluster, cfg *models.Database,
	startShardID models.ShardID, fixedStartIndex int,
) (*models.ShardAssignment, error) {
	liveNodes := cluster.GetState().LiveNodes
	if len(liveNodes) == 0 {
		return nil, constants.ErrNoLiveNode
	}
	databaseName := cfg.Name
	//TODO need calc resource and pick related node for store data

	var nodeIDs []models.NodeID
	nodes := make(map[models.NodeID]*models.StatefulNode)
	for idx := range liveNodes {
		node := liveNodes[idx]
		nodeIDs = append(nodeIDs, node.ID)
		nodes[node.ID] = &node
	}

	// generate shard assignment based on node ids and config
	shardAssign, err := ShardAssignment(nodeIDs, cfg, fixedStartIndex, startShardID)
	if err != nil {
		return nil, err
	}

	m.logger.Info("create shard assign",
		logger.String("database", databaseName),
		logger.Any("shardAssign", shardAssign))

	data := encoding.JSONMarshal(shardAssign)
	if err := m.masterRepo.Put(m.ctx, constants.GetDatabaseAssignPath(databaseName), data); err != nil {
		return nil, err
	}
	// save shard assignment into related storage storageCluster
	if err := cluster.CreateShards(databaseName, shardAssign, cfg.Option); err != nil {
		return nil, err
	}

	return shardAssign, nil
}

func (m *stateManager) modifyShardAssignment(
	cluster StorageCluster, cfg *models.Database,
	shardAssign *models.ShardAssignment,
) error {
	if len(shardAssign.Shards) > cfg.NumOfShard { //reduce shardAssign's shards
		//TODO implement the reduce shards, is needed?
		panic("not implemented")
	} else if len(shardAssign.Shards) < cfg.NumOfShard { // add shardAssign's shards
		activeNodes := cluster.GetState().LiveNodes
		if len(activeNodes) == 0 {
			return constants.ErrNoLiveNode
		}
		//TODO need calc resource and pick related node for store data

		var nodeIDs []models.NodeID
		for _, node := range activeNodes {
			nodeIDs = append(nodeIDs, node.ID)
		}

		// generate shard assignment based on node ids and config
		//TODO check start shard id
		err := ModifyShardAssignment(nodeIDs, cfg, shardAssign, -1, models.ShardID(len(shardAssign.Shards)))
		if err != nil {
			return err
		}
	}
	databaseName := cfg.Name
	m.logger.Info("modify shard assign",
		logger.String("database", databaseName),
		logger.Any("shardAssign", shardAssign))

	data := encoding.JSONMarshal(shardAssign)
	if err := m.masterRepo.Put(m.ctx, constants.GetDatabaseAssignPath(databaseName), data); err != nil {
		return err
	}

	// save shard assignment into related storage storageCluster
	if err := cluster.CreateShards(databaseName, shardAssign, cfg.Option); err != nil {
		return err
	}
	return nil
}

// GetShardAssign returns shard assignment by database name, return not exist err if it not exist
func (m *stateManager) GetShardAssign(databaseName string) (*models.ShardAssignment, error) {
	data, err := m.masterRepo.Get(m.ctx, constants.GetDatabaseAssignPath(databaseName))
	if err != nil {
		return nil, err
	}
	shardAssign := &models.ShardAssignment{}
	if err := encoding.JSONUnmarshal(data, shardAssign); err != nil {
		return nil, err
	}
	return shardAssign, nil
}

// GetStorageCluster returns cluster controller for maintain the metadata of storage cluster
func (m *stateManager) GetStorageCluster(name string) (cluster StorageCluster) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	cluster = m.storages[name]
	return
}
