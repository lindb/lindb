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
)

//go:generate mockgen -source=./node_state_machine.go -destination=./node_state_machine_mock.go -package=discovery

// NodeStateMachine represents idToNodeMap' meta state machine,
// listens node adds into cluster change event.
type NodeStateMachine interface {
	io.Closer
	inif.Listener

	// GetNodes returns all idToNodeMap info list.
	GetNodes() []models.Node
}

// nodeStateMachine implements node state machine interface, watches node path.
type nodeStateMachine struct {
	discovery Discovery

	ctx    context.Context
	cancel context.CancelFunc

	mutex       sync.RWMutex
	idToNodeMap map[models.NodeID]models.Node
	nodeMap     map[string]models.Node
	running     *atomic.Bool

	logger *logger.Logger
}

// NewNodeStateMachine creates a node state machine, and starts discovery for watching node state change event.
func NewNodeStateMachine(ctx context.Context,
	discoveryFactory Factory,
) (NodeStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	stateMachine := &nodeStateMachine{
		ctx:         c,
		cancel:      cancel,
		idToNodeMap: make(map[models.NodeID]models.Node),
		nodeMap:     make(map[string]models.Node),
		running:     atomic.NewBool(false),
		logger:      logger.GetLogger("coordinator", "NodeStateMachine"),
	}
	// new node state discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.NodesPath, stateMachine)
	if err := stateMachine.discovery.Discovery(true); err != nil {
		return nil, fmt.Errorf("discovery node meta error:%s", err)
	}
	stateMachine.running.Store(true)
	stateMachine.logger.Info("node state machine is started.")
	return stateMachine, nil
}

// GetNodes returns all idToNodeMap info list.
func (s *nodeStateMachine) GetNodes() (rs []models.Node) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if !s.running.Load() {
		// state machine not running return it.
		s.logger.Warn("get nodes when state machine is not running")
		return
	}
	for _, node := range s.idToNodeMap {
		rs = append(rs, node)
	}
	return
}

// OnCreate adds node into active node list when node online
func (s *nodeStateMachine) OnCreate(key string, resource []byte) {
	s.logger.Info("discovery new node metadata in cluster",
		logger.String("key", key),
		logger.String("data", string(resource)))

	node := models.Node{}
	if err := encoding.JSONUnmarshal(resource, &node); err != nil {
		s.logger.Error("discovery new node metadata but unmarshal error", logger.Error(err))
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.idToNodeMap[node.ID] = node
	s.nodeMap[node.Indicator()] = node
}

// OnDelete removes node into active node list when node offline
func (s *nodeStateMachine) OnDelete(key string) {
	_, fileName := filepath.Split(key)
	nodeID := fileName

	s.logger.Info("delete exist node metadata from cluster", logger.String("key", key))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	node, ok := s.nodeMap[nodeID]
	if ok {
		delete(s.nodeMap, node.Indicator())
		delete(s.idToNodeMap, node.ID)
	}
}

// Close closes state machine, then releases resource
func (s *nodeStateMachine) Close() error {
	if s.running.CAS(true, false) {
		s.mutex.Lock()
		defer func() {
			s.mutex.Unlock()
			s.cancel()
		}()

		s.discovery.Close()

		s.logger.Info("node state machine is stopped.")
	}
	return nil
}
