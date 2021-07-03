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
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./node_state_machine.go -destination=./node_state_machine_mock.go -package=broker

// NodeStateMachine represents broker nodes state machine,
// listens node online/offline change event
type NodeStateMachine interface {
	discovery.Listener
	// GetCurrentNode returns the current broker node
	GetCurrentNode() models.Node
	// GetActiveNodes returns all active broker nodes
	GetActiveNodes() []models.ActiveNode
	// Close closes state machine, then releases resource
	Close() error
}

// nodeStateMachine implements node state machine interface,
// watches active node path.
type nodeStateMachine struct {
	currentNode models.Node
	discovery   discovery.Discovery

	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.RWMutex
	// brokers: broker node => replica list under this broker
	nodes             map[string]models.ActiveNode
	connectionManager *connectionManager
	logger            *logger.Logger
}

// NewNodeStateMachine creates a node state machine, and starts discovery for watching node state change event
func NewNodeStateMachine(
	ctx context.Context,
	currentNode models.Node,
	discoveryFactory discovery.Factory,
	taskClientFactory rpc.TaskClientFactory,
) (NodeStateMachine, error) {
	c, cancel := context.WithCancel(ctx)

	stateMachine := &nodeStateMachine{
		ctx:         c,
		cancel:      cancel,
		currentNode: currentNode,
		connectionManager: &connectionManager{
			RoleFrom:          "broker",
			RoleTo:            "broker",
			connections:       make(map[string]struct{}),
			taskClientFactory: taskClientFactory,
		},
		nodes:  make(map[string]models.ActiveNode),
		logger: logger.GetLogger("coordinator", "BrokerNodeStateMachine"),
	}
	// new replica status discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.ActiveNodesPath+"/data", stateMachine)
	if err := stateMachine.discovery.Discovery(); err != nil {
		return nil, fmt.Errorf("discovery broker node error:%s", err)
	}
	return stateMachine, nil
}

// GetCurrentNode returns the current broker node
func (s *nodeStateMachine) GetCurrentNode() models.Node {
	return s.currentNode
}

// GetActiveNodes returns all active broker nodes
func (s *nodeStateMachine) GetActiveNodes() []models.ActiveNode {
	var result []models.ActiveNode
	s.mutex.RLock()
	for _, node := range s.nodes {
		result = append(result, node)
	}
	s.mutex.RUnlock()
	return result
}

// OnCreate adds node into active node list when node online
func (s *nodeStateMachine) OnCreate(key string, resource []byte) {
	node := models.ActiveNode{}
	if err := json.Unmarshal(resource, &node); err != nil {
		s.logger.Error("discovery node online but unmarshal error",
			logger.String("data", string(resource)), logger.Error(err))
		return
	}
	_, fileName := filepath.Split(key)
	nodeID := fileName
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Info("peer broker is online",
		logger.String("node", node.Node.Indicator()),
		logger.Int64("nodeOnlineTime", node.OnlineTime),
	)
	s.connectionManager.createConnection(node.Node)

	s.nodes[nodeID] = node
}

// OnDelete removes node into active node list when node offline
func (s *nodeStateMachine) OnDelete(key string) {
	_, fileName := filepath.Split(key)
	nodeID := fileName
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.nodes, nodeID)

	s.logger.Info("peer broker is offline",
		logger.String("node", nodeID),
	)
	s.connectionManager.closeConnection(nodeID)
}

// Close closes state machine, then releases resource
func (s *nodeStateMachine) Close() error {
	s.discovery.Close()
	s.mutex.Lock()
	s.nodes = make(map[string]models.ActiveNode)
	s.mutex.Unlock()
	s.cancel()

	s.connectionManager.closeAll()
	return nil
}
