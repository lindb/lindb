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

package discovery

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/inif"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./active_node_state_machine.go -destination=./active_node_state_machine_mock.go -package=discovery

// ActiveNodeStateMachine represents active nodes state machine,
// listens node online/offline change event.
type ActiveNodeStateMachine interface {
	inif.Listener
	io.Closer

	// GetCurrentNode returns the current node.
	GetCurrentNode() models.Node
	// GetActiveNodes returns all active nodes.
	GetActiveNodes() []models.ActiveNode
}

// activeNodeStateMachine implements node state machine interface,
// watches active node path.
type activeNodeStateMachine struct {
	currentNode models.Node
	discovery   Discovery

	ctx    context.Context
	cancel context.CancelFunc

	mutex   sync.RWMutex
	running *atomic.Bool

	nodes             map[string]models.ActiveNode
	connectionManager *ConnectionManager

	logger *logger.Logger
}

// NewActiveNodeStateMachine creates a active node state machine,
// and starts discovery for watching node state change event.
func NewActiveNodeStateMachine(ctx context.Context, currentNode models.Node,
	discoveryFactory Factory, taskClientFactory rpc.TaskClientFactory,
) (ActiveNodeStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	stateMachine := &activeNodeStateMachine{
		ctx:         c,
		cancel:      cancel,
		currentNode: currentNode,
		connectionManager: &ConnectionManager{
			RoleFrom:          "broker",
			RoleTo:            "broker",
			Connections:       make(map[string]struct{}),
			TaskClientFactory: taskClientFactory,
		},
		nodes:   make(map[string]models.ActiveNode),
		running: atomic.NewBool(false),
		logger:  logger.GetLogger("coordinator", "ActiveNodeStateMachine"),
	}
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.ActiveNodesPath, stateMachine)
	if err := stateMachine.discovery.Discovery(true); err != nil {
		return nil, fmt.Errorf("discovery broker active error:%s", err)
	}

	stateMachine.running.Store(true)
	stateMachine.logger.Info("active node state machine is started.")

	return stateMachine, nil
}

// GetCurrentNode returns the current broker node.
func (s *activeNodeStateMachine) GetCurrentNode() models.Node {
	return s.currentNode
}

// GetActiveNodes returns all active broker nodes
func (s *activeNodeStateMachine) GetActiveNodes() (rs []models.ActiveNode) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.running.Load() {
		s.logger.Warn("get active nodes when state machine is not running")
		return
	}

	for _, node := range s.nodes {
		rs = append(rs, node)
	}
	return
}

// OnCreate adds node into active node list when node online.
func (s *activeNodeStateMachine) OnCreate(key string, resource []byte) {
	s.logger.Info("discovery new node online in cluster",
		logger.String("key", key),
		logger.String("data", string(resource)))

	node := models.ActiveNode{}
	if err := encoding.JSONUnmarshal(resource, &node); err != nil {
		s.logger.Error("discovery node online but unmarshal error", logger.Error(err))
		return
	}

	_, fileName := filepath.Split(key)
	nodeID := fileName

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.connectionManager.CreateConnection(node.Node)

	s.nodes[nodeID] = node
}

// OnDelete removes node into active node list when node offline.
func (s *activeNodeStateMachine) OnDelete(key string) {
	_, fileName := filepath.Split(key)
	nodeID := fileName

	s.logger.Info("discovery a node offline from cluster",
		logger.String("nodeID", nodeID),
		logger.String("key", key))

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connectionManager.CloseConnection(nodeID)
	delete(s.nodes, nodeID)
}

// Close closes state machine, then releases resource.
func (s *activeNodeStateMachine) Close() error {
	if s.running.CAS(true, false) {
		s.mutex.Lock()
		defer func() {
			s.mutex.Unlock()
			s.cancel()
		}()
		s.connectionManager.CloseAll()
		s.discovery.Close()

		s.logger.Info("active node state machine is stopped.")
	}

	return nil
}
