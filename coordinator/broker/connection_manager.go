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

var connectionManagerLogger = logger.GetLogger("coordinator", "ConnectionManager")

// connectionManger manages the rpc connections
// not thread-safe
type connectionManager struct {
	RoleFrom          string
	RoleTo            string
	connections       map[string]struct{}
	taskClientFactory rpc.TaskClientFactory
}

func (manager *connectionManager) createConnection(target models.Node) {
	if err := manager.taskClientFactory.CreateTaskClient(target); err == nil {
		connectionManagerLogger.Info("established connection successfully",
			logger.String("target", target.Indicator()),
			logger.String("from", manager.RoleFrom),
			logger.String("to", manager.RoleTo),
		)
		manager.connections[target.Indicator()] = struct{}{}
	} else {
		connectionManagerLogger.Error("failed to establish connection",
			logger.String("target", target.Indicator()),
			logger.String("from", manager.RoleFrom),
			logger.String("to", manager.RoleTo),
			logger.Error(err),
		)
		delete(manager.connections, target.Indicator())
	}
}

func (manager *connectionManager) closeConnection(target string) {
	closed, err := manager.taskClientFactory.CloseTaskClient(target)
	delete(manager.connections, target)

	if closed {
		if err == nil {
			connectionManagerLogger.Info("closed connection successfully",
				logger.String("target", target),
				logger.String("from", manager.RoleFrom),
				logger.String("to", manager.RoleTo),
			)
		} else {
			connectionManagerLogger.Error("failed to close connection",
				logger.String("target", target),
				logger.String("from", manager.RoleFrom),
				logger.String("to", manager.RoleTo),
				logger.Error(err),
			)
		}
	} else {
		connectionManagerLogger.Debug("unable to close a non-existent connection",
			logger.String("target", target),
			logger.String("from", manager.RoleFrom),
			logger.String("to", manager.RoleTo),
		)
	}
}

func (manager *connectionManager) closeAll() {
	for target := range manager.connections {
		manager.closeConnection(target)
	}
}

func (manager *connectionManager) closeInactiveNodeConnections(activeNodes []string) {
	activeNodesSet := make(map[string]struct{})
	for _, node := range activeNodes {
		activeNodesSet[node] = struct{}{}
	}
	for target := range manager.connections {
		if _, exist := activeNodesSet[target]; !exist {
			manager.closeConnection(target)
		}
	}
}
