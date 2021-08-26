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
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./state_manager.go -destination=./state_manager_mock.go -package=broker

type StateManager interface {
	// OnDatabaseCfgChange triggers when database create/modify.
	OnDatabaseCfgChange(key string, data []byte)
	// OnDatabaseCfgDelete triggers when database delete.
	OnDatabaseCfgDelete(key string)
	// OnNodeStartup triggers when node online.
	OnNodeStartup(key string, data []byte)
	// OnNodeFailure trigger when node offline.
	OnNodeFailure(key string)
	// OnStorageStateChange triggers when node online.
	OnStorageStateChange(key string, data []byte)
	// OnStorageDelete trigger when node offline.
	OnStorageDelete(key string)

	// read api as below:

	// GetCurrentNode returns the current node.
	GetCurrentNode() models.StatelessNode
	// GetLiveNodes returns all live broker nodes.
	GetLiveNodes() []models.StatelessNode
	// GetDatabaseCfg returns the database config by name.
	GetDatabaseCfg(databaseName string) (models.Database, bool)
	// GetQueryableReplicas returns the queryable replicas，
	// and chooses the leader replica if the shard has multi-replica.
	// returns storage node => shard id list
	GetQueryableReplicas(databaseName string) map[string][]models.ShardID
}

type stateManager struct {
	currentNode models.StatelessNode
	// state cache
	storages  map[string]*models.StorageState
	databases map[string]models.Database
	nodes     map[string]models.StatelessNode // broker live nodes

	// connection manager
	connectionManager rpc.ConnectionManager
	taskClientFactory rpc.TaskClientFactory
	cm                replication.ChannelManager

	mutex sync.RWMutex

	logger *logger.Logger
}

func NewStateManager(
	currentNode models.StatelessNode,
	connectionManager rpc.ConnectionManager,
	taskClientFactory rpc.TaskClientFactory,
	cm replication.ChannelManager,
) StateManager {
	return &stateManager{
		currentNode:       currentNode,
		connectionManager: connectionManager,
		taskClientFactory: taskClientFactory,
		cm:                cm,
		storages:          make(map[string]*models.StorageState),
		databases:         make(map[string]models.Database),
		nodes:             make(map[string]models.StatelessNode),
		logger:            logger.GetLogger("broker", "stateManager"),
	}
}

// OnDatabaseCfgChange triggers when database create/modify.
func (m *stateManager) OnDatabaseCfgChange(key string, data []byte) {
	m.logger.Info("database config modified",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := models.Database{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		m.logger.Error("database config modified but unmarshal error", logger.Error(err))
		return
	}

	if len(cfg.Name) == 0 {
		m.logger.Error("database name cannot be empty")
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.databases[cfg.Name] = cfg
}

// OnDatabaseCfgDelete triggers when database delete.
func (m *stateManager) OnDatabaseCfgDelete(key string) {
	m.logger.Info("database config deleted",
		logger.String("key", key))

	_, databaseName := filepath.Split(key)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.databases, databaseName)
}

func (m *stateManager) OnNodeStartup(key string, data []byte) {
	m.logger.Info("new node online",
		logger.String("key", key),
		logger.String("data", string(data)))

	node := &models.StatelessNode{}
	if err := encoding.JSONUnmarshal(data, node); err != nil {
		m.logger.Error("new node online but unmarshal error", logger.Error(err))
		return
	}

	_, fileName := filepath.Split(key)
	nodeID := fileName

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.connectionManager.CreateConnection(node)

	m.nodes[nodeID] = *node
}

func (m *stateManager) OnNodeFailure(key string) {
	_, fileName := filepath.Split(key)
	nodeID := fileName

	m.logger.Info("node online => offline",
		logger.String("nodeID", nodeID),
		logger.String("key", key))

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.connectionManager.CloseConnection(nodeID)

	delete(m.nodes, nodeID)
}

func (m *stateManager) OnStorageStateChange(key string, data []byte) {
	m.logger.Info("storage state is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	newState := &models.StorageState{}
	if err := encoding.JSONUnmarshal(data, newState); err != nil {
		m.logger.Error("storage state is changed but unmarshal error", logger.Error(err))
		return
	}
	if len(newState.Name) == 0 {
		m.logger.Error("storage name is empty")
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logger.Info("storage state is changed", logger.String("storage", newState.Name))

	oldState, ok := m.storages[newState.Name]
	if ok {
		liveNodesSet := make(map[string]struct{})
		for idx := range newState.LiveNodes {
			node := newState.LiveNodes[idx]
			liveNodesSet[node.Indicator()] = struct{}{}
			// try create connection for live node
			m.connectionManager.CreateConnection(&node)
		}

		// close old deal node connection
		for _, node := range oldState.LiveNodes {
			target := node.Indicator()
			if _, exist := liveNodesSet[target]; !exist {
				m.connectionManager.CloseConnection(target)
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

	//TODO need modify
	m.buildShardAssign(newState)
}

func (m *stateManager) OnStorageDelete(key string) {
	_, name := filepath.Split(key)

	m.logger.Info("storage is deleted",
		logger.String("storage", name),
		logger.String("key", key))

	m.mutex.Lock()
	defer m.mutex.Unlock()

	state, ok := m.storages[name]
	if ok {
		// close all connection [current broker node=>storage live nodes]
		for _, node := range state.LiveNodes {
			m.connectionManager.CloseConnection(node.Indicator())
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
	return
}

// GetDatabaseCfg returns the database config by name.
func (m *stateManager) GetDatabaseCfg(databaseName string) (models.Database, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	database, ok := m.databases[databaseName]
	return database, ok
}

// GetQueryableReplicas returns the queryable replicas
// returns storage node => shard id list
func (m *stateManager) GetQueryableReplicas(databaseName string) map[string][]models.ShardID {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 1. check database if exist
	database, ok := m.databases[databaseName]
	if !ok {
		return nil
	}

	// 2. check shards if exist
	storageState, ok := m.storages[database.Storage]
	if !ok {
		m.logger.Warn("database not run on any storage",
			logger.String("storage", database.Storage),
			logger.String("database", databaseName))
		return nil
	}
	// check if has live nodes
	liveNodes := storageState.LiveNodes
	if len(liveNodes) == 0 {
		m.logger.Warn("there is no live node for this storage",
			logger.String("storage", database.Storage),
			logger.String("database", databaseName))
		return nil
	}
	shards, ok := storageState.ShardStates[databaseName]
	if !ok {
		m.logger.Warn("database's shard state be lost",
			logger.String("storage", database.Storage),
			logger.String("database", databaseName))
		return nil
	}

	if len(shards) == 0 {
		m.logger.Warn("there is no shard for this database",
			logger.String("storage", database.Storage),
			logger.String("database", databaseName))
		return nil
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
	return result
}

// buildShardAssign builds the wal replica channel and related replicators for the shard assignment
func (m *stateManager) buildShardAssign(storageState *models.StorageState) {
	liveNodes := storageState.LiveNodes
	for db, shards := range storageState.ShardStates {
		numOfShard := len(shards)
		for shardID, shardState := range shards {
			m.createReplicaChannel(db, numOfShard, shardID, shardState, liveNodes)
		}
	}
}

// createReplicaChannel creates wal replica channel for spec database's shard
func (m *stateManager) createReplicaChannel(db string,
	numOfShard int, shardID models.ShardID,
	shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode,
) {
	ch, err := m.cm.CreateChannel(db, int32(numOfShard), shardID)
	if err != nil {
		m.logger.Error("create replica channel", logger.Error(err))
		return
	}
	m.logger.Info("create replica channel successfully", logger.String("db", db), logger.Any("shardID", shardID))

	m.startReplicator(ch, db, shardID, shardState, liveNodes)
}

// startReplicator starts wal replicator for spec database's shard
func (m *stateManager) startReplicator(ch replication.Channel,
	db string, shardID models.ShardID, shardState models.ShardState,
	liveNodes map[models.NodeID]models.StatefulNode,
) {

	for _, replicaID := range shardState.Replica.Replicas {
		target, ok := liveNodes[replicaID]
		if ok {
			_, err := ch.GetOrCreateReplicator(&target)
			if err != nil {
				m.logger.Error("start replicator", logger.Error(err))
				continue
			}
			m.logger.Info("create replicator successfully", logger.String("db", db),
				logger.Any("shardID", shardID), logger.String("target", target.Indicator()))
		}
	}
}
