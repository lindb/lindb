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
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

type StorageClusterState struct {
	state             *models.StorageState
	connectionManager *connectionManager
	logger            *logger.Logger
}

func newStorageClusterState(taskClientFactory rpc.TaskClientFactory, logger *logger.Logger) *StorageClusterState {
	return &StorageClusterState{
		connectionManager: &connectionManager{
			RoleFrom:          "broker",
			RoleTo:            "storage",
			connections:       make(map[string]struct{}),
			taskClientFactory: taskClientFactory,
		},
		logger: logger,
	}
}

func (s *StorageClusterState) SetState(state *models.StorageState) {
	s.logger.Debug("set new storage cluster state", logger.String(state.Name, state.String()))
	var activeNodes []string
	for _, node := range state.GetActiveNodes() {
		activeNodes = append(activeNodes, node.Node.Indicator())
	}
	s.connectionManager.closeInactiveNodeConnections(activeNodes)

	for _, node := range state.ActiveNodes {
		s.logger.Info("storage node is online",
			logger.String("node", node.Node.Indicator()),
			logger.Int64("nodeOnlineTime", node.OnlineTime),
		)
		s.connectionManager.createConnection(node.Node)
	}

	s.state = state
	s.logger.Debug("set new storage cluster successfully")
}

func (s *StorageClusterState) close() {
	s.connectionManager.closeAll()
	s.logger.Debug("close storage cluster state successfully")
}
