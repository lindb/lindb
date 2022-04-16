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
	"path/filepath"
	"strconv"
	"sync"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./state_manager.go -destination=./state_manager_mock.go -package=storage

// for test
var getConnFct = rpc.GetStorageClientConnFactory

// StateManager represents storage state manager, maintains storage node in memory.
type StateManager interface {
	discovery.StateMachineEventHandle

	// GetLiveNode returns storage live node by node id, return false if not exist.
	GetLiveNode(nodeID models.NodeID) (models.StatefulNode, bool)
	// WatchNodeStateChangeEvent registers node state change event handle.
	WatchNodeStateChangeEvent(nodeID models.NodeID, fn func(state models.NodeStateType))
	// GetLiveNodes returns the current live nodes.
	GetLiveNodes() []models.StatefulNode
	// GetDatabaseAssignments returns the current database assignments.
	GetDatabaseAssignments() []*models.DatabaseAssignment
}

// stateManager implements StateManager.
type stateManager struct {
	ctx    context.Context
	cancel context.CancelFunc

	engine              tsdb.Engine
	current             *models.StatefulNode
	nodes               map[models.NodeID]models.StatefulNode // storage live nodes
	watches             map[models.NodeID][]func(state models.NodeStateType)
	databaseAssignments map[string]*models.DatabaseAssignment

	events chan *discovery.Event

	mutex sync.RWMutex

	logger *logger.Logger

	statistics struct {
		nodeStartUps *linmetric.BoundCounter
		nodeFailures *linmetric.BoundCounter
		shardAssigns *linmetric.BoundCounter
		panics       *linmetric.BoundCounter
	}
}

// NewStateManager creates a StateManager instance.
func NewStateManager(
	ctx context.Context,
	current *models.StatefulNode,
	engine tsdb.Engine,
) StateManager {
	c, cancel := context.WithCancel(ctx)
	mgr := &stateManager{
		ctx:                 c,
		cancel:              cancel,
		current:             current,
		engine:              engine,
		nodes:               make(map[models.NodeID]models.StatefulNode),
		databaseAssignments: make(map[string]*models.DatabaseAssignment),
		events:              make(chan *discovery.Event, 10),
		watches:             make(map[models.NodeID][]func(state models.NodeStateType)),
		logger:              logger.GetLogger("storage", "StateManager"),
	}
	scope := linmetric.StorageRegistry.NewScope("lindb.storage.state_manager")
	eventVec := scope.NewCounterVec("emit_events", "type")
	mgr.statistics.nodeStartUps = eventVec.WithTagValues("node_joins")
	mgr.statistics.nodeFailures = eventVec.WithTagValues("node_leaves")
	mgr.statistics.shardAssigns = eventVec.WithTagValues("shard_assigns")
	mgr.statistics.panics = scope.NewCounter("panics")

	// start consume discovery event task
	go mgr.consumeEvent()

	return mgr
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

// processEvent processes each event, if panic will ignore the event handle, maybe lost the state in storage/.
func (m *stateManager) processEvent(event *discovery.Event) {
	defer func() {
		if err := recover(); err != nil {
			m.statistics.panics.Incr()
			m.logger.Error("panic when process discovery event, lost the state",
				logger.Any("err", err), logger.Stack())
		}
	}()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch event.Type {
	case discovery.NodeStartup:
		m.statistics.nodeStartUps.Incr()
		m.onNodeStartup(event.Key, event.Value)
	case discovery.NodeFailure:
		m.statistics.nodeFailures.Incr()
		m.onNodeFailure(event.Key)
	case discovery.ShardAssignmentChanged:
		m.statistics.shardAssigns.Incr()
		m.onShardAssignmentChange(event.Key, event.Value)
	}
}

// onShardAssignmentChange triggers when shard assignment changed after database config modified.
func (m *stateManager) onShardAssignmentChange(key string, data []byte) {
	m.logger.Info("shard assignment is changed",
		logger.String("key", key),
		logger.String("data", string(data)))
	param := models.DatabaseAssignment{}
	if err := encoding.JSONUnmarshal(data, &param); err != nil {
		return
	}
	if param.ShardAssignment == nil {
		return
	}

	m.databaseAssignments[param.ShardAssignment.Name] = &param

	var shardIDs []models.ShardID
	for shardID, replica := range param.ShardAssignment.Shards {
		if replica.Contain(m.current.ID) {
			shardIDs = append(shardIDs, shardID)
		}
	}
	if len(shardIDs) == 0 {
		return
	}
	if err := m.engine.CreateShards(
		param.ShardAssignment.Name,
		param.Option,
		shardIDs...,
	); err != nil {
		m.logger.Error("create shard storage engine err",
			logger.String("db", param.ShardAssignment.Name),
			logger.Any("shards", shardIDs),
			logger.Error(err))
		return
	}
}

// onNodeStartup triggers when storage node online.
func (m *stateManager) onNodeStartup(key string, data []byte) {
	m.logger.Info("new node online",
		logger.String("key", key),
		logger.String("data", string(data)))

	node := &models.StatefulNode{}
	if err := encoding.JSONUnmarshal(data, node); err != nil {
		m.logger.Error("new node online but unmarshal error", logger.Error(err))
		return
	}

	m.nodes[node.ID] = *node

	// notify node online
	watches := m.watches[node.ID]
	for _, handle := range watches {
		handle(models.NodeOnline)
	}
}

// onNodeFailure triggers when storage node offline.
func (m *stateManager) onNodeFailure(key string) {
	_, fileName := filepath.Split(key)

	m.logger.Info("node online => offline",
		logger.String("nodeID", fileName),
		logger.String("key", key))

	id, err := strconv.ParseInt(fileName, 10, 64)
	if err != nil {
		m.logger.Error("parse offline node id err", logger.Error(err))
		return
	}

	nodeID := models.NodeID(id)
	node, ok := m.nodes[nodeID]
	if !ok {
		// node not exist in alive node list
		return
	}
	delete(m.nodes, nodeID)

	// notify node offline
	watches := m.watches[nodeID]
	for _, handle := range watches {
		handle(models.NodeOffline)
	}
	// try close offline node connection in pool
	if err := getConnFct().CloseClientConn(&node); err != nil {
		m.logger.Error("close connection for offline node err", logger.Error(err))
	}
}

// GetLiveNode returns storage live node by node id, return false if not exist.
func (m *stateManager) GetLiveNode(nodeID models.NodeID) (models.StatefulNode, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	node, ok := m.nodes[nodeID]
	return node, ok
}

// WatchNodeStateChangeEvent registers node state change event handle.
func (m *stateManager) WatchNodeStateChangeEvent(nodeID models.NodeID, fn func(state models.NodeStateType)) {
	if fn == nil {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	watches := m.watches[nodeID]
	watches = append(watches, fn)
	m.watches[nodeID] = watches
}

// GetLiveNodes returns the current live nodes.
func (m *stateManager) GetLiveNodes() (rs []models.StatefulNode) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for idx := range m.nodes {
		rs = append(rs, m.nodes[idx])
	}
	return
}

// GetDatabaseAssignments returns the current database assignments.
func (m *stateManager) GetDatabaseAssignments() (rs []*models.DatabaseAssignment) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, databaseAssignment := range m.databaseAssignments {
		rs = append(rs, databaseAssignment)
	}
	return
}
