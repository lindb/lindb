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
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

var connectionManagerLogger = logger.GetLogger("coordinator", "ConnectionManager")

// ConnectionManager manages the rpc Connections
// not thread-safe
type ConnectionManager struct {
	RoleFrom          string
	RoleTo            string
	Connections       map[string]struct{}
	TaskClientFactory rpc.TaskClientFactory
}

func (manager *ConnectionManager) CreateConnection(target models.Node) {
	if err := manager.TaskClientFactory.CreateTaskClient(target); err == nil {
		connectionManagerLogger.Info("established connection successfully",
			logger.String("target", target.Indicator()),
			logger.String("from", manager.RoleFrom),
			logger.String("to", manager.RoleTo),
		)
		manager.Connections[target.Indicator()] = struct{}{}
	} else {
		connectionManagerLogger.Error("failed to establish connection",
			logger.String("target", target.Indicator()),
			logger.String("from", manager.RoleFrom),
			logger.String("to", manager.RoleTo),
			logger.Error(err),
		)
		delete(manager.Connections, target.Indicator())
	}
}

func (manager *ConnectionManager) CloseConnection(target string) {
	closed, err := manager.TaskClientFactory.CloseTaskClient(target)
	delete(manager.Connections, target)

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

func (manager *ConnectionManager) CloseAll() {
	for target := range manager.Connections {
		manager.CloseConnection(target)
	}
}

func (manager *ConnectionManager) CloseInactiveNodeConnections(activeNodes []string) {
	activeNodesSet := make(map[string]struct{})
	for _, node := range activeNodes {
		activeNodesSet[node] = struct{}{}
	}
	for target := range manager.Connections {
		if _, exist := activeNodesSet[target]; !exist {
			manager.CloseConnection(target)
		}
	}
}
