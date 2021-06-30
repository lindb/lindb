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
	"io"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/inif"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"

	"go.uber.org/atomic"
)

//go:generate mockgen -source=./database_state_machine.go -destination=./database_state_machine_mock.go -package=broker

// DatabaseStateMachine represents alive database config state machine,
// listens database create/delete change event.
type DatabaseStateMachine interface {
	inif.Listener
	io.Closer

	// GetDatabaseCfg returns the database config by name.
	GetDatabaseCfg(databaseName string) (models.Database, bool)
}

// databaseStateMachine implements DatabaseStateMachine
type databaseStateMachine struct {
	discovery discovery.Discovery

	databases map[string]models.Database
	running   *atomic.Bool

	mutex  sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	logger *logger.Logger
}

// NewDatabaseStateMachine creates database config state machine instance.
func NewDatabaseStateMachine(
	ctx context.Context,
	discoveryFactory discovery.Factory,
) (DatabaseStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	// new admin state machine instance
	stateMachine := &databaseStateMachine{
		ctx:       c,
		cancel:    cancel,
		running:   atomic.NewBool(false),
		databases: make(map[string]models.Database),
		logger:    logger.GetLogger("coordinator", "DatabaseStateMachine"),
	}

	// new database config discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.DatabaseConfigPath, stateMachine)
	if err := stateMachine.discovery.Discovery(true); err != nil {
		return nil, fmt.Errorf("discovery database config error:%s", err)
	}

	stateMachine.running.Store(true)
	stateMachine.logger.Info("database state machine is started")

	return stateMachine, nil
}

// OnCreate adds database config into list when database creation
func (sm *databaseStateMachine) OnCreate(key string, resource []byte) {
	sm.logger.Info("discovery new database create in cluster",
		logger.String("key", key),
		logger.String("data", string(resource)))

	cfg := models.Database{}
	if err := json.Unmarshal(resource, &cfg); err != nil {
		sm.logger.Error("discovery database create but unmarshal error", logger.Error(err))
		return
	}

	if len(cfg.Name) == 0 {
		sm.logger.Error("database name cannot be empty")
		return
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.databases[cfg.Name] = cfg
}

// OnDelete removes database config from list when database deletion.
func (sm *databaseStateMachine) OnDelete(key string) {
	sm.logger.Info("discovery a database delete from cluster",
		logger.String("key", key))

	_, databaseName := filepath.Split(key)

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.databases, databaseName)
}

// GetDatabaseCfg returns the database config by name.
func (sm *databaseStateMachine) GetDatabaseCfg(databaseName string) (models.Database, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if !sm.running.Load() {
		sm.logger.Warn("get database cfg when state machine is not running")
		return models.Database{}, false
	}

	database, ok := sm.databases[databaseName]
	return database, ok
}

// Close closes database config state machine, stops watch change event.
func (sm *databaseStateMachine) Close() error {
	if sm.running.CAS(true, false) {
		sm.mutex.Lock()
		defer func() {
			sm.mutex.Unlock()
			sm.cancel()
		}()

		sm.discovery.Close()
		sm.logger.Info("database config state machine is stopped.")
	}
	return nil
}
