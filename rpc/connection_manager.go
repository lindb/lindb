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

package rpc

import (
	"io"
	"sync"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./connection_manager.go -destination=./connection_manager_mock.go -package=rpc

type ConnectionManager interface {
	io.Closer

	CreateConnection(target models.Node)
	CloseConnection(target string)
}

type connectionManager struct {
	connections   map[string]struct{}
	taskClientFct TaskClientFactory

	mutex sync.Mutex

	logger *logger.Logger
}

func NewConnectionManager(taskClientFct TaskClientFactory) ConnectionManager {
	return &connectionManager{
		taskClientFct: taskClientFct,
		connections:   make(map[string]struct{}),
		logger:        logger.GetLogger("RPC", "ConnectionManager"),
	}
}

func (m *connectionManager) CreateConnection(target models.Node) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.taskClientFct.CreateTaskClient(target); err == nil {
		m.logger.Info("established connection successfully",
			logger.String("target", target.Indicator()),
		)
		m.connections[target.Indicator()] = struct{}{}
	} else {
		m.logger.Error("failed to establish connection",
			logger.String("target", target.Indicator()),
			logger.Error(err),
		)
		delete(m.connections, target.Indicator())
	}
}

func (m *connectionManager) CloseConnection(target string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.closeConnection(target)
}

func (m *connectionManager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for target := range m.connections {
		m.closeConnection(target)
	}
	return nil
}

func (m *connectionManager) closeConnection(target string) {
	closed, err := m.taskClientFct.CloseTaskClient(target)
	delete(m.connections, target)

	if closed {
		if err == nil {
			m.logger.Info("closed connection successfully",
				logger.String("target", target),
			)
		} else {
			m.logger.Error("failed to close connection",
				logger.String("target", target),
				logger.Error(err),
			)
		}
	} else {
		m.logger.Debug("unable to close a non-existent connection",
			logger.String("target", target),
		)
	}
}
