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

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
)

// StateMachinePaths represents the paths which broker state machine need watch.
var StateMachinePaths = make(map[string]models.StateMachineInfo)

func init() {
	StateMachinePaths["LiveNode"] = models.StateMachineInfo{
		Path: constants.LiveNodesPath,
		CreateState: func() interface{} {
			return &models.StatelessNode{}
		},
	}
	StateMachinePaths["DatabaseConfig"] = models.StateMachineInfo{
		Path: constants.DatabaseConfigPath,
		CreateState: func() interface{} {
			return &models.Database{}
		},
	}
	StateMachinePaths["StorageState"] = models.StateMachineInfo{
		Path: constants.StorageStatePath,
		CreateState: func() interface{} {
			return &models.StorageState{}
		},
	}
}

// stateMachineFactory implements discovery.StateMachineFactory.
type stateMachineFactory struct {
	ctx              context.Context
	discoveryFactory discovery.Factory
	stateMgr         StateManager

	stateMachines []discovery.StateMachine

	logger *logger.Logger
}

// NewStateMachineFactory creates a state machine factory instance.
func NewStateMachineFactory(
	ctx context.Context,
	discoveryFactory discovery.Factory,
	stateMgr StateManager,
) discovery.StateMachineFactory {
	return &stateMachineFactory{
		ctx:              ctx,
		discoveryFactory: discoveryFactory,
		stateMgr:         stateMgr,
		logger:           logger.GetLogger("broker", "StateMachineFactory"),
	}
}

// Start starts all broker's related state machines.
func (f *stateMachineFactory) Start() (err error) {
	f.logger.Debug("starting LiveNodeStateMachine")
	sm, err := f.createBrokerLiveNodeStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)

	f.logger.Debug("starting DatabaseConfigStateMachine")
	sm, err = f.createDatabaseCfgStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)

	f.logger.Debug("starting StorageStatusStateMachine")
	sm, err = f.createStorageStatusStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)

	f.logger.Info("started BrokerStateMachines")
	return nil
}

// Stop stops all broker's related state machines.
func (f *stateMachineFactory) Stop() {
	f.logger.Info("stopping broker state machines...")
	for _, sm := range f.stateMachines {
		if err := sm.Close(); err != nil {
			f.logger.Error("close state machine error", logger.Error(err))
		}
	}
}

// createBrokerLiveNodeStateMachine creates broker live node state machine.
func (f *stateMachineFactory) createBrokerLiveNodeStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachineFn(
		f.ctx,
		discovery.LiveNodeStateMachine,
		f.discoveryFactory,
		constants.LiveNodesPath,
		true,
		f.onNodeStartup,
		f.onNodeFailure,
	)
}

// createDatabaseCfgStateMachine creates database config state machine.
func (f *stateMachineFactory) createDatabaseCfgStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachineFn(
		f.ctx,
		discovery.DatabaseConfigStateMachine,
		f.discoveryFactory,
		constants.DatabaseConfigPath,
		true,
		f.onDatabaseConfigChanged,
		f.onDatabaseConfigDeletion,
	)
}

// createStorageStatusStateMachine creates storage status state machine.
func (f *stateMachineFactory) createStorageStatusStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachineFn(
		f.ctx,
		discovery.StorageStatusStateMachine,
		f.discoveryFactory,
		constants.StorageStatePath,
		true,
		f.onStorageStateChange,
		f.onStorageDeletion,
	)
}

// onDatabaseConfigChanged triggers when database config modified(create/update)
func (f *stateMachineFactory) onDatabaseConfigChanged(key string, data []byte) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type:  discovery.DatabaseConfigChanged,
		Key:   key,
		Value: data,
	})
}

// onDatabaseConfigDeletion triggers when database is deletion.
func (f *stateMachineFactory) onDatabaseConfigDeletion(key string) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type: discovery.DatabaseConfigDeletion,
		Key:  key,
	})
}

// onNodeStartup triggers when node online.
func (f *stateMachineFactory) onNodeStartup(key string, data []byte) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   key,
		Value: data,
	})
}

// onNodeFailure triggers when node offline.
func (f *stateMachineFactory) onNodeFailure(key string) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  key,
	})
}

// onStorageStateChange triggers when storage state changed.
func (f *stateMachineFactory) onStorageStateChange(key string, data []byte) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type:  discovery.StorageStateChanged,
		Key:   key,
		Value: data,
	})
}

// onStorageDeletion triggers when storage is deletion.
func (f *stateMachineFactory) onStorageDeletion(key string) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type: discovery.StorageDeletion,
		Key:  key,
	})
}
