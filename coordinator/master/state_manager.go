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
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	statepkg "github.com/lindb/lindb/pkg/state"
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
	// GetDatabases returns the current databases.
	GetDatabases() []models.Database
	// GetStorages returns the current storage cluster list.
	GetStorages() []config.StorageCluster
	// GetShardAssignments returns the current shard assignment list.
	GetShardAssignments() []models.ShardAssignment
	// GetStorageStates returns current storage state list.
	GetStorageStates() []*models.StorageState
}

// stateManager implements StateManager.
type stateManager struct {
	ctx    context.Context
	cancel context.CancelFunc

	repoFactory     statepkg.RepositoryFactory
	stateMachineFct *StateMachineFactory

	masterRepo statepkg.Repository
	elector    ReplicaLeaderElector

	newStorageClusterFn func(ctx context.Context, cfg *config.StorageCluster,
		stateMgr StateManager,
		repoFactory statepkg.RepositoryFactory) (cluster StorageCluster, err error)

	storages         map[string]StorageCluster
	databases        map[string]*models.Database
	shardAssignments map[string]*models.ShardAssignment

	events chan *discovery.Event

	running *atomic.Bool
	mutex   sync.RWMutex

	statistics            *metrics.StateManagerStatistics
	shardLeaderStatistics *metrics.ShardLeaderStatistics
	logger                *logger.Logger
}

// NewStateManager creates a StateManager instance.
func NewStateManager(
	ctx context.Context,
	masterRepo statepkg.Repository,
	repoFactory statepkg.RepositoryFactory,
) StateManager {
	c, cancel := context.WithCancel(ctx)
	mgr := &stateManager{
		ctx:                   c,
		cancel:                cancel,
		masterRepo:            masterRepo,
		repoFactory:           repoFactory,
		storages:              make(map[string]StorageCluster),
		databases:             make(map[string]*models.Database),
		shardAssignments:      make(map[string]*models.ShardAssignment),
		elector:               newReplicaLeaderElector(),
		events:                make(chan *discovery.Event, 10),
		running:               atomic.NewBool(true),
		newStorageClusterFn:   newStorageCluster,
		statistics:            metrics.NewStateManagerStatistics(strings.ToLower(constants.MasterRole)),
		shardLeaderStatistics: metrics.NewShardLeaderStatistics(),
		logger:                logger.GetLogger("Master", "StateManager"),
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

// processEvent processes each event, if panic will ignore the event handle, maybe lost the state in storage/.
func (m *stateManager) processEvent(event *discovery.Event) {
	eventType := event.Type.String()
	defer func() {
		if err := recover(); err != nil {
			m.statistics.Panics.WithTagValues(eventType).Incr()
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
	var err error
	switch event.Type {
	case discovery.StorageConfigChanged:
		err = m.onStorageConfigChange(event.Key, event.Value)
	case discovery.StorageConfigDeletion:
		m.onStorageConfigDelete(event.Key)
	case discovery.DatabaseConfigChanged:
		err = m.onDatabaseCfgChange(event.Key, event.Value)
	case discovery.DatabaseConfigDeletion:
		err = m.onDatabaseCfgDelete(event.Key)
	case discovery.ShardAssignmentChanged:
		err = m.onShardAssignmentChange(event.Key, event.Value)
	case discovery.NodeStartup:
		err = m.onStorageNodeStartup(event.Attributes[storageNameKey], event.Key, event.Value)
	case discovery.NodeFailure:
		err = m.onStorageNodeFailure(event.Attributes[storageNameKey], event.Key)
	}
	if err != nil {
		m.statistics.HandleEventFailure.WithTagValues(eventType).Incr()
	} else {
		m.statistics.HandleEvents.WithTagValues(eventType).Incr()
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
func (m *stateManager) onDatabaseCfgChange(key string, data []byte) error {
	m.logger.Info("do shard assignment, because database config is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := &models.Database{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		m.logger.Error("do shard assignment, because database config is changed, but unmarshal error",
			logger.Error(err))
		return err
	}

	m.shardAssignment(cfg)
	return nil
}

// onDatabaseCfgDelete triggers when database config is deletion.
func (m *stateManager) onDatabaseCfgDelete(key string) error {
	m.logger.Info("database config deleted",
		logger.String("key", key))
	name := strings.TrimPrefix(key, constants.GetDatabaseConfigPath(""))
	databaseCfg, ok := m.databases[name]
	if !ok {
		return constants.ErrDatabaseNotFound
	}
	delete(m.databases, name)
	delete(m.shardAssignments, name)

	storage := m.storages[databaseCfg.Storage]
	// remove database state from storage cluster
	storage.GetState().DropDatabase(name)

	// finally, sync storage state
	if err := m.syncState(storage.GetState()); err != nil {
		return err
	}
	if err := storage.DropDatabaseAssignment(name); err != nil {
		m.logger.Error("drop database assignment failure",
			logger.String("storage", databaseCfg.Storage),
			logger.String("database", name),
			logger.Error(err))
		return err
	}
	return nil
}

// onShardAssignmentChange triggers when shard assignment modify.
func (m *stateManager) onShardAssignmentChange(key string, data []byte) error {
	m.logger.Info("database's shard assignment is changed",
		logger.String("key", key),
		logger.String("data", string(data)))
	shardAssignment := &models.ShardAssignment{}
	if err := encoding.JSONUnmarshal(data, shardAssignment); err != nil {
		m.logger.Error("database's shard assignment is changed, but unmarshal error",
			logger.Error(err))
		return err
	}
	m.shardAssignments[shardAssignment.Name] = shardAssignment

	databaseCfg := m.databases[shardAssignment.Name]

	storage := m.storages[databaseCfg.Storage]

	m.initializeShardState(storage, shardAssignment)
	return m.syncState(storage.GetState())
}

// onStorageConfigChange triggers when storage config create/modify.
func (m *stateManager) onStorageConfigChange(key string, data []byte) error {
	m.logger.Info("storage config is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := &config.StorageCluster{}
	if err := encoding.JSONUnmarshal(data, cfg); err != nil {
		m.logger.Error("storage config modified but unmarshal error", logger.Error(err))
		return err
	}

	if err := m.register(cfg); err != nil {
		m.logger.Error("register new storage cluster", logger.Error(err))
		return err
	}
	return nil
}

// onStorageConfigDelete triggers when storage config is deletion.
func (m *stateManager) onStorageConfigDelete(key string) {
	m.logger.Info("storage config deleted",
		logger.String("key", key))

	name := strings.TrimPrefix(key, constants.StorageConfigPath)

	m.unRegister(name)
}

// onStorageNodeStartup triggers when storage node online
func (m *stateManager) onStorageNodeStartup(storageName, key string, data []byte) error {
	m.logger.Info("new storage node online in storage cluster",
		logger.String("storage", storageName),
		logger.String("key", key),
		logger.String("data", string(data)))

	node := models.StatefulNode{}
	if err := json.Unmarshal(data, &node); err != nil {
		m.logger.Error("new storage node online in storage cluster but unmarshal error", logger.Error(err))
		return err
	}

	cluster := m.storages[storageName]
	s := cluster.GetState()

	s.NodeOnline(node)

	m.onNodeStartup(s, node)

	return m.syncState(s)
}

// onStorageNodeFailure triggers when storage node offline.
func (m *stateManager) onStorageNodeFailure(storageName, key string) error {
	m.logger.Info("a storage node offline in storage cluster",
		logger.String("storage", storageName),
		logger.String("key", key))

	_, nodeIDStr := filepath.Split(key)
	id, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		m.logger.Error("parse offline node id err", logger.Error(err))
		return nil
	}

	cluster := m.storages[storageName]
	s := cluster.GetState()
	// 1. set node offline
	nodeID := models.NodeID(id)
	s.NodeOffline(nodeID)
	// 2. do node offline state change
	m.onNodeFailure(s, nodeID)

	return m.syncState(s)
}

// register starts storage state machine which watch storage state change.
func (m *stateManager) register(cfg *config.StorageCluster) error {
	name := cfg.Config.Namespace
	if name == "" {
		return constants.ErrNameEmpty
	}

	// check storage if it's exist, just config modify
	_, exist := m.storages[name]
	if exist {
		// shutdown old storageCluster state machine if exist
		m.unRegister(name)
	}

	// TODO add config
	cfg.Config.DialTimeout = ltoml.Duration(5 * time.Second)
	cfg.Config.Timeout = ltoml.Duration(5 * time.Second)
	cluster, err := m.newStorageClusterFn(m.ctx, cfg, m, m.repoFactory)
	if err != nil {
		return err
	}
	m.storages[name] = cluster
	if exist {
		// if storage is existed, need to load shard assigment for this storage
		for k := range m.shardAssignments {
			shardAssignment := m.shardAssignments[k]
			databaseCfg := m.databases[shardAssignment.Name]
			if databaseCfg.Storage == name {
				// if shard assigment belong current cluster, need to do initialize shard state
				m.initializeShardState(cluster, shardAssignment)
			}
		}
	}
	// start storage cluster state machine.
	go func() {
		// need start storage cluster state machine in background,
		// because maybe load too many storage nodes when state machine init, emits too many event into event chan,
		// if chan is full, will be blocked, then trigger data race.
		if err := cluster.Start(); err != nil {
			// need lock
			m.mutex.Lock()
			defer m.mutex.Unlock()

			m.unRegister(name)
			m.logger.Warn("start storage cluster failure", logger.String("storage", name), logger.Error(err))
		}
	}()
	return nil
}

// deleteCluster deletes the storageCluster if exist.
func (m *stateManager) unRegister(name string) {
	if cluster, ok := m.storages[name]; ok {
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
func (m *stateManager) shardAssignment(databaseCfg *models.Database) {
	if databaseCfg.Name == "" {
		m.logger.Error("database name cannot be empty")
		return
	}

	m.databases[databaseCfg.Name] = databaseCfg

	// get shard assignment from repo, maybe mem state is not sync.
	shardAssign, err := m.GetShardAssign(databaseCfg.Name)
	if err != nil && err != statepkg.ErrNotExist {
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
		_, err := m.createShardAssignment(cluster, databaseCfg, -1, -1)
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
		if err := m.modifyShardAssignment(cluster, databaseCfg, shardAssign); err != nil {
			m.logger.Error("modify shard assignment error",
				logger.String("storage", databaseCfg.Storage),
				logger.Any("databaseCfg", databaseCfg),
				logger.Error(err))
			return
		}
	default:
		// TODO remove it ???
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
		if shardStates, ok := state.ShardStates[db]; ok {
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
}

func (m *stateManager) onNodeFailure(state *models.StorageState, nodeID models.NodeID) {
	// 1. find all leaders on failure node, need do leader elect
	leadersOnOfflineNode := state.LeadersOnNode(nodeID)
	m.logger.Debug("leader node is offline need elect new leader for shard",
		logger.Any("shards", leadersOnOfflineNode))

	liveNodes := state.LiveNodes
	for db, shards := range leadersOnOfflineNode {
		shardAssignment := state.ShardAssignments[db]
		shardStates := state.ShardStates[db]
		for _, shardID := range shards {
			leader, err := m.elector.ElectLeader(shardAssignment, liveNodes, shardID)
			shardState := shardStates[shardID]
			m.shardLeaderStatistics.LeaderElections.Incr()
			if err != nil {
				shardState.State = models.OfflineShard
				shardState.Leader = models.NoLeader
				m.shardLeaderStatistics.LeaderElectFailures.Incr()
				m.logger.Warn("elect shard leader err",
					logger.String("db", shardAssignment.Name),
					logger.Any("shard", shardID), logger.Error(err))
			} else {
				shardState.State = models.OnlineShard
				shardState.Leader = leader
				m.logger.Info("elect new leader for shard",
					logger.String("db", shardAssignment.Name),
					logger.Any("shard", shardID),
					logger.Any("leader", leader))
			}
			shardStates[shardID] = shardState
		}
	}
}

// syncState syncs storage state into state repo.
func (m *stateManager) syncState(state *models.StorageState) error {
	// TODO add timeout
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	data := encoding.JSONMarshal(state)
	if err := m.masterRepo.Put(ctx, constants.GetStorageStatePath(state.Name), data); err != nil {
		m.logger.Error("sync storage state error", logger.String("storage", state.Name), logger.Error(err))
		return err
	}
	m.logger.Info("sync storage state successfully", logger.String("storage", state.Name))
	return nil
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
	// TODO need calc resource and pick related node for store data

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
	if len(shardAssign.Shards) > cfg.NumOfShard { // reduce shardAssign's shards
		// TODO implement the reduce shards, is needed?
		panic("not implemented")
	} else if len(shardAssign.Shards) < cfg.NumOfShard { // add shardAssign's shards
		liveNodes, err := cluster.GetLiveNodes()
		if err != nil {
			return err
		}
		if len(liveNodes) == 0 {
			return constants.ErrNoLiveNode
		}
		// TODO need calc resource and pick related node for store data

		var nodeIDs []models.NodeID
		for idx := range liveNodes {
			node := liveNodes[idx]
			nodeIDs = append(nodeIDs, node.ID)
			nodes[node.ID] = &node
		}

		// generate shard assignment based on node ids and config
		// TODO check start shard id
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

// GetShardAssign returns shard assignment by database name, return not exist err if it's not exist.
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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	cluster = m.storages[name]
	return
}

// GetDatabases returns the current databases.
func (m *stateManager) GetDatabases() (rs []models.Database) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, db := range m.databases {
		rs = append(rs, *db)
	}
	return
}

// GetStorages returns the current storage cluster list.
func (m *stateManager) GetStorages() (rs []config.StorageCluster) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, storage := range m.storages {
		rs = append(rs, *storage.GetConfig())
	}
	return
}

// GetShardAssignments returns the current shard assignment list.
func (m *stateManager) GetShardAssignments() (rs []models.ShardAssignment) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, shardAssignment := range m.shardAssignments {
		rs = append(rs, *shardAssignment)
	}
	return
}

// GetStorageStates returns current storage state list.
func (m *stateManager) GetStorageStates() (rs []*models.StorageState) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, storage := range m.storages {
		rs = append(rs, storage.GetState())
	}
	return
}

// initializeShardState initializes the shard state based on shard assignment for storage cluster.
func (m *stateManager) initializeShardState(storage StorageCluster, shardAssignment *models.ShardAssignment) {
	storageState := storage.GetState()
	liveNodes := storageState.LiveNodes
	shardStates := make(map[models.ShardID]models.ShardState)
	for shardID, replicas := range shardAssignment.Shards {
		leader, err := m.elector.ElectLeader(shardAssignment, liveNodes, shardID)
		shardState := models.ShardState{ID: shardID, Replica: *replicas}
		m.shardLeaderStatistics.LeaderElections.Incr()
		if err != nil {
			shardState.State = models.OfflineShard
			shardState.Leader = models.NoLeader
			m.shardLeaderStatistics.LeaderElectFailures.Incr()
			m.logger.Warn("elect shard leader err",
				logger.String("db", shardAssignment.Name),
				logger.Any("shard", shardID), logger.Error(err))
		} else {
			shardState.State = models.OnlineShard
			shardState.Leader = leader
		}
		shardStates[shardID] = shardState
	}
	// TODO set shard assignments
	storageState.ShardAssignments[shardAssignment.Name] = shardAssignment
	storageState.ShardStates[shardAssignment.Name] = shardStates
}
