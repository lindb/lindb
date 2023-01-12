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
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./state_manager.go -destination=./state_manager_mock.go -package=broker

// StateManager represents broker state manager, maintains broker node/database/storage states in memory.
type StateManager interface {
	flow.NodeChoose
	discovery.StateMachineEventHandle

	// GetCurrentNode returns the current node.
	GetCurrentNode() models.StatelessNode
	// GetLiveNodes returns all live broker nodes.
	GetLiveNodes() []models.StatelessNode
	// GetDatabaseCfg returns the database config by name.
	GetDatabaseCfg(databaseName string) (models.Database, bool)
	// GetDatabases returns current database config list.
	GetDatabases() []models.Database
	// GetQueryableReplicas returns the queryable replicas，
	// and chooses the leader replica if the shard has multi-replica.
	// returns storage node => shard id list
	GetQueryableReplicas(databaseName string) (map[string][]models.ShardID, error)
	// GetStorage returns storage state by name.
	GetStorage(name string) (*models.StorageState, bool)
	// GetStorageList returns all storage state list.
	GetStorageList() (rs []*models.StorageState)

	WatchShardStateChangeEvent(fn func(databaseCfg models.Database,
		shards map[models.ShardID]models.ShardState,
		liveNodes map[models.NodeID]models.StatefulNode,
	))
}

// stateManager implements StateManager.
type stateManager struct {
	ctx    context.Context
	cancel context.CancelFunc

	// state cache
	currentNode models.StatelessNode
	storages    map[string]*models.StorageState // storage state
	databases   map[string]models.Database      // database config
	nodes       map[string]models.StatelessNode // live nodes of broker cluster

	callbacks []func(databaseCfg models.Database,
		shards map[models.ShardID]models.ShardState,
		liveNodes map[models.NodeID]models.StatefulNode,
	)
	// connection manager
	connectionManager rpc.ConnectionManager
	//FIXME: remove it???
	taskClientFactory rpc.TaskClientFactory

	events chan *discovery.Event
	mutex  sync.RWMutex

	statistics *metrics.StateManagerStatistics
	logger     *logger.Logger
}

// NewStateManager creates a broker state manager instance.
func NewStateManager(
	ctx context.Context,
	currentNode models.StatelessNode,
	connectionManager rpc.ConnectionManager,
	taskClientFactory rpc.TaskClientFactory,
) StateManager {
	c, cancel := context.WithCancel(ctx)
	mgr := &stateManager{
		ctx:               c,
		cancel:            cancel,
		currentNode:       currentNode,
		connectionManager: connectionManager,
		taskClientFactory: taskClientFactory,
		storages:          make(map[string]*models.StorageState),
		databases:         make(map[string]models.Database),
		nodes:             make(map[string]models.StatelessNode),
		events:            make(chan *discovery.Event, 10),
		statistics:        metrics.NewStateManagerStatistics(strings.ToLower(constants.BrokerRole)),
		logger:            logger.GetLogger("Broker", "StateManager"),
	}

	// start consume discovery event task
	go mgr.consumeEvent()

	return mgr
}

// Choose chooses the compute nodes then builds physical plan.
// if need node num > 1, need pick live broker nodes as compute node,
// else pick storage replica node as leaf node.
func (m *stateManager) Choose(database string, numOfNodes int) ([]*models.PhysicalPlan, error) {
	// FIXME: need using storage's replica state ???
	replicas, err := m.GetQueryableReplicas(database)
	if err != nil {
		return nil, err
	}
	nodesLen := len(replicas)
	if nodesLen == 0 {
		return nil, constants.ErrReplicaNotFound
	}
	if numOfNodes > 1 && nodesLen > 1 {
		// build compute target nodes.
		return []*models.PhysicalPlan{flow.BuildPhysicalPlan(database, m.GetLiveNodes(), numOfNodes)}, nil
	}
	// build leaf storage nodes.
	physicalPlan := &models.PhysicalPlan{
		Database: database,
	}
	for storageNode, shardIDs := range replicas {
		physicalPlan.AddTarget(&models.Target{
			Indicator: storageNode,
			ShardIDs:  shardIDs,
		})
	}
	return []*models.PhysicalPlan{physicalPlan}, nil
}

func (m *stateManager) WatchShardStateChangeEvent(fn func(databaseCfg models.Database,
	shards map[models.ShardID]models.ShardState,
	liveNodes map[models.NodeID]models.StatefulNode,
)) {
	if fn != nil {
		m.mutex.Lock()
		m.callbacks = append(m.callbacks, fn)
		m.mutex.Unlock()
	}
}

// EmitEvent emits discovery event when state changed.
func (m *stateManager) EmitEvent(event *discovery.Event) {
	m.events <- event
}

// Close cleans the resource(stop the task).
func (m *stateManager) Close() {
	m.cancel()
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

// processEvent processes each event, if panic will ignore the event handle, maybe lost the state in broker.
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

	var err error
	switch event.Type {
	case discovery.DatabaseConfigChanged:
		err = m.onDatabaseCfgChange(event.Key, event.Value)
	case discovery.DatabaseConfigDeletion:
		m.onDatabaseCfgDelete(event.Key)
	case discovery.NodeStartup:
		err = m.onNodeStartup(event.Key, event.Value)
	case discovery.NodeFailure:
		m.onNodeFailure(event.Key)
	case discovery.StorageStateChanged:
		err = m.onStorageStateChange(event.Key, event.Value)
	case discovery.StorageStateDeletion:
		m.onStorageDelete(event.Key)
	}
	if err != nil {
		m.statistics.HandleEventFailure.WithTagValues(eventType).Incr()
	} else {
		m.statistics.HandleEvents.WithTagValues(eventType).Incr()
	}
}

// onDatabaseCfgChange triggers when database create/modify.
func (m *stateManager) onDatabaseCfgChange(key string, data []byte) error {
	m.logger.Info("database config is modified",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := models.Database{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		m.logger.Error("database config modified but unmarshal error", logger.Error(err))
		return err
	}

	if cfg.Name == "" {
		m.logger.Error("database name cannot be empty")
		return constants.ErrNameEmpty
	}

	m.databases[cfg.Name] = cfg
	return nil
}

// onDatabaseCfgDelete triggers when database is deletion.
func (m *stateManager) onDatabaseCfgDelete(key string) {
	m.logger.Info("database config deleted",
		logger.String("key", key))

	_, databaseName := filepath.Split(key)

	delete(m.databases, databaseName)
}

// onNodeStartup triggers when broker node online.
func (m *stateManager) onNodeStartup(key string, data []byte) error {
	m.logger.Info("new broker node online",
		logger.String("key", key),
		logger.String("data", string(data)))

	node := &models.StatelessNode{}
	if err := encoding.JSONUnmarshal(data, node); err != nil {
		m.logger.Error("new broker node online but unmarshal error", logger.Error(err))
		return err
	}

	_, fileName := filepath.Split(key)
	nodeID := fileName

	m.connectionManager.CreateConnection(node)

	m.nodes[nodeID] = *node

	return nil
}

// onNodeFailure triggers when broker node offline.
func (m *stateManager) onNodeFailure(key string) {
	_, fileName := filepath.Split(key)
	nodeID := fileName

	m.logger.Info("broker node online => offline",
		logger.String("nodeID", nodeID),
		logger.String("key", key))

	node, ok := m.nodes[nodeID]
	if ok {
		m.connectionManager.CloseConnection(&node)
	}

	delete(m.nodes, nodeID)
}

// onStorageStateChange triggers when storage cluster state changed.
func (m *stateManager) onStorageStateChange(key string, data []byte) error {
	m.logger.Info("storage state is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	newState := &models.StorageState{}
	if err := encoding.JSONUnmarshal(data, newState); err != nil {
		m.logger.Error("storage state is changed but unmarshal error", logger.Error(err))
		return err
	}
	if newState.Name == "" {
		m.logger.Error("storage name is empty")
		return constants.ErrNameEmpty
	}

	if oldState, ok := m.storages[newState.Name]; ok {
		liveNodesSet := make(map[string]struct{})
		for idx := range newState.LiveNodes {
			node := newState.LiveNodes[idx]
			liveNodesSet[node.Indicator()] = struct{}{}
			// try to create connection for live node
			m.connectionManager.CreateConnection(&node)
		}

		// close old deal node connection
		for _, node := range oldState.LiveNodes {
			target := node.Indicator()
			if _, exist := liveNodesSet[target]; !exist {
				m.connectionManager.CloseConnection(&node)
			}
		}
	} else {
		// create connection current broker node connect to storage live node
		for idx := range newState.LiveNodes {
			node := newState.LiveNodes[idx]
			m.connectionManager.CreateConnection(&node)
		}
	}
	// set state into cache
	m.storages[newState.Name] = newState

	m.logger.Info("storage state is changed successful, start notify shard state change",
		logger.String("storage", newState.Name))

	m.notifyShardStateChange(newState)
	return nil
}

// onStorageDelete triggers when storage cluster is deletion.
func (m *stateManager) onStorageDelete(key string) {
	_, name := filepath.Split(key)

	m.logger.Info("storage is deleted",
		logger.String("storage", name),
		logger.String("key", key))

	if state, ok := m.storages[name]; ok {
		// close all connection [current broker node=>storage live nodes]
		for _, node := range state.LiveNodes {
			m.connectionManager.CloseConnection(&node)
		}

		delete(m.storages, name)
	}
}

// GetCurrentNode returns the current broker node.
func (m *stateManager) GetCurrentNode() models.StatelessNode {
	return m.currentNode
}

// GetLiveNodes returns all live broker nodes.
func (m *stateManager) GetLiveNodes() (rs []models.StatelessNode) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, node := range m.nodes {
		rs = append(rs, node)
	}

	// return nodes in order(by ip)
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].HostIP < rs[j].HostIP
	})
	return
}

// GetDatabaseCfg returns the database config by name.
func (m *stateManager) GetDatabaseCfg(databaseName string) (models.Database, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	database, ok := m.databases[databaseName]
	return database, ok
}

// GetDatabases returns current database config list.
func (m *stateManager) GetDatabases() (rs []models.Database) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for idx := range m.databases {
		rs = append(rs, m.databases[idx])
	}
	return
}

// GetStorage returns storage state by name.
func (m *stateManager) GetStorage(name string) (*models.StorageState, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	storage, ok := m.storages[name]
	return storage, ok
}

// GetStorageList returns all storage state list.
func (m *stateManager) GetStorageList() (rs []*models.StorageState) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, s := range m.storages {
		rs = append(rs, s)
	}
	// return storage in order(by node)
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Name < rs[j].Name
	})
	return
}

// GetQueryableReplicas returns the queryable replicas, else return detail error msg.::x
// returns storage node => shard id list
func (m *stateManager) GetQueryableReplicas(databaseName string) (map[string][]models.ShardID, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 1. check database if exist
	database, ok := m.databases[databaseName]
	if !ok {
		return nil, constants.ErrDatabaseNotFound
	}

	// 2. check shards if exist
	storageState, ok := m.storages[database.Storage]
	if !ok {
		m.logger.Warn("database not run on any storage",
			logger.String("storage", database.Storage),
			logger.String("database", databaseName))
		return nil, constants.ErrNoStorageCluster
	}
	// check if it has live nodes
	liveNodes := storageState.LiveNodes
	if len(liveNodes) == 0 {
		m.logger.Warn("there is no live node for this storage",
			logger.String("storage", database.Storage),
			logger.String("database", databaseName))
		return nil, constants.ErrNoLiveNode
	}
	shards := storageState.ShardStates[databaseName]
	if len(shards) == 0 {
		m.logger.Warn("there is no shard for this database",
			logger.String("storage", database.Storage),
			logger.String("database", databaseName))
		return nil, constants.ErrShardNotFound
	}

	result := make(map[string][]models.ShardID)
	for shardID, shardState := range shards {
		if shardState.State == models.OnlineShard {
			node := liveNodes[shardState.Leader]
			nodeID := node.Indicator()
			result[nodeID] = append(result[nodeID], shardID)
		} else {
			m.logger.Warn("shard is not online ignore it, maybe query data will be lost",
				logger.String("storage", database.Storage),
				logger.String("database", databaseName),
				logger.Any("shard", shardState.ID))
		}
	}
	return result, nil
}

// buildShardAssign builds the data write channel and related shard state.
func (m *stateManager) notifyShardStateChange(storageState *models.StorageState) {
	liveNodes := storageState.LiveNodes
	for db, shards := range storageState.ShardStates {
		databaseCfg := m.databases[db]
		for _, fn := range m.callbacks {
			fn(databaseCfg, shards, liveNodes)
		}
	}
}
