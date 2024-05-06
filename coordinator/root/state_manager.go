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

package root

import (
	"context"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	statepkg "github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./state_manager.go -destination=./state_manager_mock.go -package=root

// StateManager represents root state manager, state coordinator.
type StateManager interface {
	discovery.StateMachineEventHandle
	flow.NodeChoose
	// SetStateMachineFactory sets state machine factory.
	SetStateMachineFactory(stateMachineFct *stateMachineFactory)
	// GetStateMachineFactory returns state machine factory.
	GetStateMachineFactory() *stateMachineFactory
	// GetBrokerStates returns current broker state list.
	GetBrokerStates() []models.BrokerState
	// GetBrokerState returns current broker state by name.
	GetBrokerState(name string) (models.BrokerState, bool)
	// GetDatabase returns the logic database config by name.
	GetDatabase(name string) (models.LogicDatabase, bool)
	// GetDatabases returns all logic databases' config.
	GetDatabases() []models.LogicDatabase
	// GetLiveNodes returns all live broker nodes.
	GetLiveNodes() []models.StatelessNode
}

// stateManager implements StateManager.
type stateManager struct {
	ctx                context.Context
	cancel             context.CancelFunc
	repoFactory        statepkg.RepositoryFactory
	stateMachineFct    *stateMachineFactory
	brokers            map[string]BrokerCluster
	databases          map[string]*models.LogicDatabase
	nodes              map[string]models.StatelessNode // live nodes of root cluster
	events             chan *discovery.Event
	running            *atomic.Bool
	newBrokerClusterFn func(cfg *config.BrokerCluster,
		stateMgr StateManager,
		repoFactory statepkg.RepositoryFactory) (cluster BrokerCluster, err error)
	// connection manager
	connectionManager rpc.ConnectionManager

	mutex sync.RWMutex

	statistics *metrics.StateManagerStatistics
	logger     logger.Logger
}

// NewStateManager creates a StateManager instance.
func NewStateManager(
	ctx context.Context,
	repoFactory statepkg.RepositoryFactory,
	connectionManager rpc.ConnectionManager,
) StateManager {
	c, cancel := context.WithCancel(ctx)
	mgr := &stateManager{
		ctx:                c,
		cancel:             cancel,
		repoFactory:        repoFactory,
		brokers:            make(map[string]BrokerCluster),
		databases:          make(map[string]*models.LogicDatabase),
		events:             make(chan *discovery.Event, 10),
		nodes:              make(map[string]models.StatelessNode),
		running:            atomic.NewBool(true),
		connectionManager:  connectionManager,
		statistics:         metrics.NewStateManagerStatistics(linmetric.RootRegistry),
		newBrokerClusterFn: newBrokerCluster,
		logger:             logger.GetLogger("Root", "StateManager"),
	}

	// start consume event then do coordinator
	go mgr.consumeEvent()

	return mgr
}

// Choose chooses the compute nodes then builds physical plan.
func (s *stateManager) Choose(database string, numOfNodes int) ([]*models.PhysicalPlan, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	databaseCfg, ok := s.databases[database]
	if !ok {
		return nil, nil
	}

	var rs []*models.PhysicalPlan

	getDatabase := func(db, defalutValue string) string {
		if db != "" {
			return db
		}
		return defalutValue
	}

	routers := databaseCfg.Routers
	for _, router := range routers {
		broker, ok := s.brokers[router.Broker]
		if !ok {
			s.logger.Warn("broker cluster is offline, will ingore this cluster", logger.String("broker", router.Broker))
			continue
		}
		rs = append(rs, flow.BuildPhysicalPlan(getDatabase(router.Database, database), broker.GetState().GetLiveNodes(), numOfNodes))
	}
	return rs, nil
}

// EmitEvent emits discovery event when state changed.
func (s *stateManager) EmitEvent(event *discovery.Event) {
	s.events <- event
}

// consumeEvent consumes the discovery event, then handles the event by each event type.
func (s *stateManager) consumeEvent() {
	for {
		select {
		case event := <-s.events:
			s.processEvent(event)
		case <-s.ctx.Done():
			s.logger.Info("consume discovery event task is stopped")
			return
		}
	}
}

// processEvent processes each event, if panic will ignore the event handle, maybe lost the state in storage.
func (s *stateManager) processEvent(event *discovery.Event) {
	eventType := event.Type.String()
	defer func() {
		if err := recover(); err != nil {
			s.statistics.Panics.WithTagValues(eventType, constants.RootRole).Incr()
			s.logger.Error("panic when process discovery event, lost the state",
				logger.Any("err", err), logger.Stack())
		}
	}()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running.Load() {
		s.logger.Warn("root state manager is closed")
		return
	}
	var err error
	switch event.Type {
	case discovery.BrokerConfigChanged:
		err = s.onBrokerConfigChange(event.Key, event.Value)
	case discovery.BrokerConfigDeletion:
		s.onBrokerConfigDelete(event.Key)
	case discovery.DatabaseConfigChanged:
		err = s.onDatabaseCfgChange(event.Key, event.Value)
	case discovery.DatabaseConfigDeletion:
		s.onDatabaseCfgDelete(event.Key)
	case discovery.NodeStartup:
		if event.Attributes == nil {
			err = s.onNodeStartup(event.Key, event.Value)
		} else {
			err = s.onBrokerNodeStartup(event.Attributes[brokerNameKey], event.Key, event.Value)
		}
	case discovery.NodeFailure:
		if event.Attributes == nil {
			s.onNodeFailure(event.Key)
		} else {
			s.onBrokerNodeFailure(event.Attributes[brokerNameKey], event.Key)
		}
	}
	if err != nil {
		s.statistics.HandleEventFailure.WithTagValues(eventType, constants.RootRole).Incr()
	} else {
		s.statistics.HandleEvents.WithTagValues(eventType, constants.RootRole).Incr()
	}
}

// onBrokerConfigDelete triggers when storage config is deletion.
func (s *stateManager) onBrokerConfigDelete(key string) {
	s.logger.Info("broker config deleted",
		logger.String("key", key))

	name := strings.TrimPrefix(key, constants.BrokerConfigPath)

	s.unRegister(name)
}

// onBrokerConfigChange triggers when broker config create/modify.
func (s *stateManager) onBrokerConfigChange(key string, data []byte) error {
	s.logger.Info("broker config is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := &config.BrokerCluster{}
	if err := encoding.JSONUnmarshal(data, cfg); err != nil {
		s.logger.Error("broker config modified but unmarshal error", logger.Error(err))
		return err
	}

	if err := s.register(cfg); err != nil {
		s.logger.Error("register new broker cluster", logger.Error(err))
		return err
	}
	return nil
}

// onDatabaseCfgChange triggers when database create/modify.
func (s *stateManager) onDatabaseCfgChange(key string, data []byte) error {
	s.logger.Info("database config is changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := &models.LogicDatabase{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		s.logger.Error("database config is changed, but unmarshal error",
			logger.Error(err))
		return err
	}
	s.databases[cfg.Name] = cfg
	return nil
}

// onDatabaseCfgDelete triggers when database config is deletion.
func (s *stateManager) onDatabaseCfgDelete(key string) {
	s.logger.Info("database config deleted",
		logger.String("key", key))
	name := strings.TrimPrefix(key, constants.GetDatabaseConfigPath(""))
	delete(s.databases, name)
}

// onNodeStartup triggers when root node online.
func (s *stateManager) onNodeStartup(key string, data []byte) error {
	s.logger.Info("new root node online",
		logger.String("key", key),
		logger.String("data", string(data)))

	node := &models.StatelessNode{}
	if err := encoding.JSONUnmarshal(data, node); err != nil {
		s.logger.Error("new root node online but unmarshal error", logger.Error(err))
		return err
	}

	_, fileName := filepath.Split(key)
	nodeID := fileName

	s.nodes[nodeID] = *node

	return nil
}

// onNodeFailure triggers when root node offline.
func (s *stateManager) onNodeFailure(key string) {
	_, fileName := filepath.Split(key)
	nodeID := fileName

	s.logger.Info("root node online => offline",
		logger.String("nodeID", nodeID),
		logger.String("key", key))

	delete(s.nodes, nodeID)
}

// onBrokerNodeStartup triggers when broker node online
func (s *stateManager) onBrokerNodeStartup(brokerName, key string, data []byte) error {
	s.logger.Info("new broker node online in broker cluster",
		logger.String("broker", brokerName),
		logger.String("key", key),
		logger.String("data", string(data)))

	node := models.StatelessNode{}
	if err := encoding.JSONUnmarshal(data, &node); err != nil {
		s.logger.Error("new broker node online in storage cluster but unmarshal error", logger.Error(err))
		return err
	}
	_, nodeID := filepath.Split(key)

	s.connectionManager.CreateConnection(&node)

	cluster := s.brokers[brokerName]
	state := cluster.GetState()
	state.NodeOnline(nodeID, node)
	return nil
}

// onBrokerNodeFailure triggers when broker node offline.
func (s *stateManager) onBrokerNodeFailure(brokerName, key string) {
	s.logger.Info("a broker node offline in broker cluster",
		logger.String("broker", brokerName),
		logger.String("key", key))

	_, nodeID := filepath.Split(key)

	cluster := s.brokers[brokerName]
	state := cluster.GetState()

	node, ok := state.LiveNodes[nodeID]
	if ok {
		s.connectionManager.CloseConnection(&node)
	}
	state.NodeOffline(nodeID)
}

// register starts storage state machine which watch storage state change.
func (s *stateManager) register(cfg *config.BrokerCluster) error {
	if cfg.Config == nil || cfg.Config.Namespace == "" {
		return constants.ErrNameEmpty
	}
	name := cfg.Config.Namespace
	// check broker if it's exist, just config modify
	_, exist := s.brokers[name]
	if exist {
		// shutdown old storageCluster state machine if exist
		s.unRegister(name)
	}

	cluster, err := s.newBrokerClusterFn(cfg, s, s.repoFactory)
	if err != nil {
		return err
	}
	s.brokers[name] = cluster
	// start broker cluster state machine.
	go func() {
		// need start broker cluster state machine in background,
		// because maybe load too many broker nodes when state machine init, emits too many event into event chan,
		// if chan is full, will be blocked, then trigger data race.
		if err := cluster.Start(); err != nil {
			// need lock
			s.mutex.Lock()
			defer s.mutex.Unlock()

			s.unRegister(name)
			s.logger.Warn("start broker cluster failure", logger.String("broker", name), logger.Error(err))
		}
	}()
	return nil
}

// deleteCluster deletes the brokerCluster if exist.
func (s *stateManager) unRegister(name string) {
	if cluster, ok := s.brokers[name]; ok {
		// need cleanup broker cluster resource
		cluster.Close()

		delete(s.brokers, name)

		s.logger.Info("cleanup broker cluster resource finished", logger.String("broker", name))
	}
}

// GetBrokerStates returns current broker state list.
func (s *stateManager) GetBrokerStates() (rs []models.BrokerState) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, broker := range s.brokers {
		rs = append(rs, *broker.GetState())
	}
	// return nodes in order(by ip)
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Name < rs[j].Name
	})
	return
}

// GetBrokerState returns current broker state by name.
func (s *stateManager) GetBrokerState(name string) (models.BrokerState, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	broker, ok := s.brokers[name]
	if !ok {
		return models.BrokerState{}, false
	}
	return *broker.GetState(), true
}

// GetDatabase returns the logic database config by name.
func (s *stateManager) GetDatabase(name string) (models.LogicDatabase, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	database, ok := s.databases[name]
	if !ok {
		return models.LogicDatabase{}, false
	}
	return *database, true
}

// GetDatabases returns all logic databases' config.
func (s *stateManager) GetDatabases() (rs []models.LogicDatabase) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, db := range s.databases {
		rs = append(rs, *db)
	}

	// return nodes in order(by ip)
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Name < rs[j].Name
	})
	return
}

// SetStateMachineFactory sets state machine factory.
func (s *stateManager) SetStateMachineFactory(stateMachineFct *stateMachineFactory) {
	s.stateMachineFct = stateMachineFct
}

// GetStateMachineFactory returns state machine factory.
func (s *stateManager) GetStateMachineFactory() *stateMachineFactory {
	return s.stateMachineFct
}

// GetLiveNodes returns all live broker nodes.
func (s *stateManager) GetLiveNodes() (rs []models.StatelessNode) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, node := range s.nodes {
		rs = append(rs, node)
	}

	// return nodes in order(by ip)
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Indicator() < rs[j].Indicator()
	})
	return
}

// Close implements StateManager
func (s *stateManager) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running.CompareAndSwap(true, false) {
		s.logger.Info("starting close state manager")
		for name := range s.brokers {
			s.unRegister(name)
		}
	}
	s.cancel()
}
