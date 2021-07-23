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
	"github.com/lindb/lindb/pkg/logger"
)

type stateMachineFactory struct {
	ctx              context.Context
	discoveryFactory discovery.Factory
	stateMgr         StateManager

	stateMachines []discovery.StateMachine

	logger *logger.Logger
}

func NewStateMachineFactory(
	ctx context.Context,
	discoveryFactory discovery.Factory,
	stateMgr StateManager,
) discovery.StateMachineFactory {
	return &stateMachineFactory{
		ctx:              ctx,
		discoveryFactory: discoveryFactory,
		stateMgr:         stateMgr,
		logger:           logger.GetLogger("coordinator", "BrokerStateMachines"),
	}
}

// Start starts related state machines for broker.
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

// Stop stops the broker's state machines.
func (f *stateMachineFactory) Stop() {
	for _, sm := range f.stateMachines {
		if err := sm.Close(); err != nil {
			f.logger.Error("close state machine error", logger.Error(err))
		}
	}
}

// createBrokerLiveNodeStateMachine creates broker live node state machine.
func (f *stateMachineFactory) createBrokerLiveNodeStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.LiveNodeStateMachine,
		f.discoveryFactory,
		constants.LiveNodesPath,
		true,
		f.stateMgr.OnNodeStartup,
		f.stateMgr.OnNodeFailure,
	)
}

// createDatabaseCfgStateMachine creates database config state machine.
func (f *stateMachineFactory) createDatabaseCfgStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.DatabaseConfigStateMachine,
		f.discoveryFactory,
		constants.DatabaseConfigPath,
		true,
		f.stateMgr.OnDatabaseCfgChange,
		f.stateMgr.OnDatabaseCfgDelete,
	)
}

// createStorageStatusStateMachine creates storage status state machine.
func (f *stateMachineFactory) createStorageStatusStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.StorageStatusStateMachine,
		f.discoveryFactory,
		constants.StorageStatePath,
		true,
		f.stateMgr.OnStorageStateChange,
		f.stateMgr.OnStorageDelete,
	)
}
