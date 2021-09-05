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
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./state_manager.go -destination=./state_manager_mock.go -package=master

// StateManager represents master state manager, state coordinator.
type StateManager interface {
	discovery.StateMachineEventHandle

	// SetStateMachineFactory sets state machine factory.
	SetStateMachineFactory(stateMachineFct *StateMachineFactory)
	// GetStateMachineFactory returns state machine factory.
	GetStateMachineFactory() *StateMachineFactory

	// GetStorageCluster returns cluster controller for maintain the metadata of storage cluster.
	GetStorageCluster(name string) StorageCluster
}

// stateManager implements StateManager.
type stateManager struct {
	ctx    context.Context
	cancel context.CancelFunc

	repoFactory       state.RepositoryFactory
	controllerFactory task.ControllerFactory
	stateMachineFct   *StateMachineFactory

	masterRepo state.Repository
	elector    ReplicaLeaderElector

	newStorageClusterFn func(ctx context.Context, cfg config.StorageCluster,
		stateMgr StateManager,
		repoFactory state.RepositoryFactory,
		controllerFactory task.ControllerFactory) (cluster StorageCluster, err error)

	storages  map[string]StorageCluster
	databases map[string]models.Database

	events chan *discovery.Event

	running *atomic.Bool
	mutex   sync.Mutex

	logger *logger.Logger
}

// NewStateManager creates a StateManager instance.
func NewStateManager(
	ctx context.Context,
	masterRepo state.Repository,
	repoFactory state.RepositoryFactory,
	controllerFactory task.ControllerFactory,
) StateManager {
	c, cancel := context.WithCancel(ctx)
	mgr := &stateManager{
		ctx:                 c,
		cancel:              cancel,
		masterRepo:          masterRepo,
		repoFactory:         repoFactory,
		controllerFactory:   controllerFactory,
		storages:            make(map[string]StorageCluster),
		databases:           make(map[string]models.Database),
		elector:             newReplicaLeaderElector(),
		events:              make(chan *discovery.Event, 10),
		running:             atomic.NewBool(true),
		newStorageClusterFn: newStorageCluster,
		logger:              logger.GetLogger("master", "StateManager"),
	}

	// start consume event then do coordinate
	go mgr.consumeEvent()

	return mgr
}

// EmitEvent emits discovery event when state changed.
func (m *stateManager) EmitEvent(event *discovery.Event) {
	m.events <- event
}

// consumeEvent consumes the discovery event, then handles the event by each event type.
func (m *stateManager) consumeEvent() {
	for {
		select {
		case event := <-m.events:
			m.processEvent(event)
		case <-m.ctx.Done():
			m.logger.Info("consume discovery event task is stopped")
			return
		}
	}
}

// processEvent processes each events, if panic will ignore the event handle, maybe lost the state in storage/.
func (m *stateManager) processEvent(event *discovery.Event) {
	defer func() {
		if err := recover(); err != nil {
			//TODO add metric
			m.logger.Error("panic when process discovery event, lost the state",
				logger.Any("err", err), logger.Stack())
		}
	}()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running.Load() {
		m.logger.Warn("master state manager is closed")
		return
	}
	switch event.Type {
	case discovery.StorageConfigChanged:
		m.onStorageConfigChange(event.Key, event.Value)
	case discovery.StorageDeletion:
		m.onStorageConfigDelete(event.Key)
	case discovery.DatabaseConfigChanged:
		m.onDatabaseCfgChange(event.Key, event.Value)
	case discovery.DatabaseConfigDeletion:
		m.onDatabaseCfgDelete(event.Key)
	case discovery.ShardAssignmentChanged:
		m.onShardAssignmentChange(event.Key, event.Value)
	case discovery.NodeStartup:
		m.onStorageNodeStartup(event.Attributes[storageNameKey], event.Key, event.Value)
	case discovery.NodeFailure:
		m.onStorageNodeFailure(event.Attributes[storageNameKey], event.Key)
	}
}

// SetStateMachineFactory sets state machine factory.
func (m *stateManager) SetStateMachineFactory(stateMachineFct *StateMachineFactory) {
	m.stateMachineFct = stateMachineFct
}

// GetStateMachineFactory returns state machine factory.
func (m *stateManager) GetStateMachineFactory() *StateMachineFactory {
	return m.stateMachineFct
}

// onDatabaseCfgChange triggers when database create/modify.
func (m *stateManager) onDatabaseCfgChange(key string, data []byte) {
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

	m.shardAssignment(cfg)
}

// onDatabaseCfgDelete triggers when database config is deletion.
func (m *stateManager) onDatabaseCfgDelete(_ string) {
	panic("need impl")
}

// onShardAssignmentChange triggers when shard assignment modify.
func (m *stateManager) onShardAssignmentChange(key string, data []byte) {
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
		shardState := models.ShardState{ID: shardID, Replica: *replicas}
		if err != nil {
			shardState.State = models.OfflineShard
			shardState.Leader = models.NoLeader
			//TODO add log and metric
			m.logger.Warn("elect shard leader err",
				logger.String("db", shardAssignment.Name),
				logger.Any("shard", shardID), logger.Error(err))
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

// onStorageConfigChange triggers when storage config create/modify.
func (m *stateManager) onStorageConfigChange(key string, data []byte) {
	m.logger.Info("storage config is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := config.StorageCluster{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		m.logger.Error("storage config modified but unmarshal error", logger.Error(err))
		return
	}

	if err := m.register(cfg); err != nil {
		m.logger.Error("register new storage cluster", logger.Error(err))
		return
	}
}

// onStorageConfigDelete triggers when storage config is deletion.
func (m *stateManager) onStorageConfigDelete(key string) {
	m.logger.Info("storage config deleted",
		logger.String("key", key))

	_, name := filepath.Split(key)

	m.unRegister(name)
}

// onStorageNodeStartup triggers when storage node online
func (m *stateManager) onStorageNodeStartup(storageName string, key string, data []byte) {
	m.logger.Info("new storage node online in storage cluster",
		logger.String("storage", storageName),
		logger.String("key", key),
		logger.String("data", string(data)))

	node := models.StatefulNode{}
	if err := json.Unmarshal(data, &node); err != nil {
		m.logger.Error("new storage node online in storage cluster but unmarshal error", logger.Error(err))
		return
	}

	cluster := m.storages[storageName]
	s := cluster.GetState()

	s.NodeOnline(node)

	m.onNodeStartup(s, node)

	m.syncState(s)
}

// onStorageNodeFailure triggers when storage node offline.
func (m *stateManager) onStorageNodeFailure(storageName string, key string) {
	m.logger.Info("a storage node offline in storage cluster",
		logger.String("storage", storageName),
		logger.String("key", key))

	_, nodeIDStr := filepath.Split(key)
	id, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		m.logger.Error("parse offline node id err", logger.Error(err))
		return
	}

	cluster := m.storages[storageName]
	s := cluster.GetState()

	m.onNodeFailure(s, models.NodeID(id))

	m.syncState(s)
}

// register registers start storage state machine which watch storage state change.
func (m *stateManager) register(cfg config.StorageCluster) error {
	if len(cfg.Name) == 0 {
		return constants.ErrNameEmpty
	}

	// shutdown old storageCluster state machine if exist
	m.unRegister(cfg.Name)

	//TODO add config
	cfg.Config.DialTimeout = ltoml.Duration(5 * time.Second)
	cfg.Config.Timeout = ltoml.Duration(5 * time.Second)
	cluster, err := m.newStorageClusterFn(m.ctx, cfg, m, m.repoFactory, m.controllerFactory)
	if err != nil {
		return err
	}
	m.storages[cfg.Name] = cluster
	// start storage cluster state machine.
	// TODO need start in goroutine???
	if err := cluster.Start(); err != nil {
		m.unRegister(cfg.Name)
		m.logger.Info("start storage cluster failure", logger.String("storage", cfg.Name), logger.Error(err))
		return err
	}
	return nil
}

// deleteCluster deletes the storageCluster if exist.
func (m *stateManager) unRegister(name string) {
	cluster, ok := m.storages[name]
	if ok {
		// need cleanup storage cluster resource
		cluster.Close()

		delete(m.storages, name)

		m.logger.Info("cleanup storage cluster resource finished", logger.String("storage", name))
	}
}

func (m *stateManager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.running.CAS(true, false) {
		m.logger.Info("starting close state manager")
		for name := range m.storages {
			m.unRegister(name)
		}
		m.cancel()
	}
}

// shardAssignment does shard assignment.
func (m *stateManager) shardAssignment(databaseCfg models.Database) {
	if len(databaseCfg.Name) == 0 {
		m.logger.Error("database name cannot be empty")
		return
	}

	// get shard assignment from repo, maybe mem state is not sync.
	shardAssign, err := m.GetShardAssign(databaseCfg.Name)
	if err != nil && err != state.ErrNotExist {
		m.logger.Error("get shard assign error", logger.Error(err))
		return
	}

	cluster := m.storages[databaseCfg.Storage]

	switch {
	case shardAssign == nil:
		// build shard assignment for creation database, generate related coordinator task
		m.logger.Info("create shard assignment starting....",
			logger.String("storage", databaseCfg.Storage),
			logger.Any("database", databaseCfg.Name))
		_, err := m.createShardAssignment(cluster, &databaseCfg, -1, -1)
		if err != nil {
			m.logger.Error("create shard assignment error",
				logger.String("storage", databaseCfg.Storage),
				logger.Any("databaseCfg", databaseCfg),
				logger.Error(err))
			return
		}
	case len(shardAssign.Shards) != databaseCfg.NumOfShard:
		m.logger.Info("modify shard assignment starting....",
			logger.String("storage", databaseCfg.Storage),
			logger.Any("database", databaseCfg.Name))
		if err := m.modifyShardAssignment(cluster, &databaseCfg, shardAssign); err != nil {
			m.logger.Error("modify shard assignment error",
				logger.String("storage", databaseCfg.Storage),
				logger.Any("databaseCfg", databaseCfg),
				logger.Error(err))
			return
		}
	default:
		//TODO remove it ???
		m.logger.Info("no data changed, just trigger shard assignment data modify event",
			logger.String("storage", databaseCfg.Storage),
			logger.Any("database", databaseCfg.Name))
		data := encoding.JSONMarshal(shardAssign)
		if err := m.masterRepo.Put(m.ctx, constants.GetDatabaseAssignPath(shardAssign.Name), data); err != nil {
			m.logger.Error("trigger shard assignment data modify event",
				logger.String("storage", databaseCfg.Storage),
				logger.Any("database", databaseCfg.Name),
				logger.Error(err))
			return
		}
	}
}

func (m *stateManager) onNodeStartup(state *models.StorageState, node models.StatefulNode) {
	// 1. do when a new node come up is send it the entire list of shards that it is supposed to host.
	replicasOnOnlineNode := state.ReplicasOnNode(node.ID)
	for db, shards := range replicasOnOnlineNode {
		shardStates, ok := state.ShardStates[db]
		if !ok {
			continue
		}
		for _, shardID := range shards {
			shardState := shardStates[shardID]
			if shardState.State != models.OnlineShard {
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

				m.logger.Warn("elect shard leader err",
					logger.String("db", shardAssignment.Name),
					logger.Any("shard", shardID), logger.Error(err))
			} else {
				shardState.State = models.OnlineShard
				shardState.Leader = leader
			}
			shardStates[shardID] = shardState
		}
	}
}

// syncState syncs storage state into state repo.
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
	liveNodes, err := cluster.GetLiveNodes()
	if err != nil {
		return nil, err
	}

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
	// save shard assignment into related storage repo.
	if err := cluster.SaveDatabaseAssignment(shardAssign, cfg.Option); err != nil {
		return nil, err
	}

	return shardAssign, nil
}

func (m *stateManager) modifyShardAssignment(
	cluster StorageCluster, cfg *models.Database,
	shardAssign *models.ShardAssignment,
) error {
	nodes := make(map[models.NodeID]*models.StatefulNode)
	if len(shardAssign.Shards) > cfg.NumOfShard { //reduce shardAssign's shards
		//TODO implement the reduce shards, is needed?
		panic("not implemented")
	} else if len(shardAssign.Shards) < cfg.NumOfShard { // add shardAssign's shards
		liveNodes, err := cluster.GetLiveNodes()
		if err != nil {
			return err
		}
		if len(liveNodes) == 0 {
			return constants.ErrNoLiveNode
		}
		//TODO need calc resource and pick related node for store data

		var nodeIDs []models.NodeID
		for idx := range liveNodes {
			node := liveNodes[idx]
			nodeIDs = append(nodeIDs, node.ID)
			nodes[node.ID] = &node
		}

		// generate shard assignment based on node ids and config
		//TODO check start shard id
		err = ModifyShardAssignment(nodeIDs, cfg, shardAssign, -1, models.ShardID(len(shardAssign.Shards)))
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

	// save shard assignment into related storage repo.
	if err := cluster.SaveDatabaseAssignment(shardAssign, cfg.Option); err != nil {
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

// GetStorageCluster returns cluster controller for maintain the metadata of storage cluster.
func (m *stateManager) GetStorageCluster(name string) (cluster StorageCluster) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	cluster = m.storages[name]
	return
}
