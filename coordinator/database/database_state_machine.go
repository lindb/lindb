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

package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./database_state_machine.go -destination=./database_state_machine_mock.go -package=database

// DBStateMachine represents alive database config state machine,
// listens database create/delete change event
type DBStateMachine interface {
	discovery.Listener

	// GetDatabaseCfg returns the database config by name
	GetDatabaseCfg(databaseName string) (models.Database, bool)

	// Close closes database config state machine, stops watch change event
	io.Closer
}

// dbStateMachine implements DBStateMachine
type dbStateMachine struct {
	discovery discovery.Discovery

	databases map[string]models.Database
	mutex     sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc

	log *logger.Logger
}

// NewDBStateMachine creates database config state machine instance
func NewDBStateMachine(
	ctx context.Context,
	repo state.Repository,
	discoveryFactory discovery.Factory,
) (DBStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	// new admin state machine instance
	stateMachine := &dbStateMachine{
		ctx:       c,
		cancel:    cancel,
		databases: make(map[string]models.Database),
		log:       logger.GetLogger("coordinator", "DBStateMachine"),
	}
	databaseList, err := repo.List(c, constants.DatabaseConfigPath)
	if err != nil {
		return nil, fmt.Errorf("get database config list error:%s", err)
	}

	// init exist database list
	for _, database := range databaseList {
		stateMachine.addDatabase(database.Value)
	}
	// new database config discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.DatabaseConfigPath, stateMachine)
	if err := stateMachine.discovery.Discovery(); err != nil {
		return nil, fmt.Errorf("discovery database config error:%s", err)
	}
	return stateMachine, nil
}

// OnCreate adds database config into list when database creation
func (sm *dbStateMachine) OnCreate(_ string, resource []byte) {
	sm.addDatabase(resource)
}

// OnDelete removes database config from list when database deletion
func (sm *dbStateMachine) OnDelete(key string) {
	_, databaseName := filepath.Split(key)

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.databases, databaseName)
}

// GetDatabaseCfg returns the database config by name
func (sm *dbStateMachine) GetDatabaseCfg(databaseName string) (models.Database, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	database, ok := sm.databases[databaseName]
	return database, ok
}

// Close closes database config state machine, stops watch change event
func (sm *dbStateMachine) Close() error {
	sm.discovery.Close()
	sm.cancel()
	return nil
}

// addDatabase parses and adds database config into state machine.
func (sm *dbStateMachine) addDatabase(resource []byte) {
	cfg := models.Database{}
	if err := json.Unmarshal(resource, &cfg); err != nil {
		sm.log.Error("discovery database create but unmarshal error",
			logger.String("data", string(resource)), logger.Error(err))
		return
	}

	if len(cfg.Name) == 0 {
		sm.log.Error("database name cannot be empty", logger.String("data", string(resource)))
		return
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.databases[cfg.Name] = cfg
}
