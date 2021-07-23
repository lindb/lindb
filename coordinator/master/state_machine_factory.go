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

package master

import (
	"context"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/pkg/logger"
)

type StateMachineFactory struct {
	ctx              context.Context
	discoveryFactory discovery.Factory
	stateMgr         StateManager

	stateMachines []discovery.StateMachine

	logger *logger.Logger
}

func NewStateMachineFactory(ctx context.Context,
	discoveryFactory discovery.Factory,
	stateMgr StateManager,
) *StateMachineFactory {
	return &StateMachineFactory{
		ctx:              ctx,
		discoveryFactory: discoveryFactory,
		stateMgr:         stateMgr,
		logger:           logger.GetLogger("coordinator", "MasterStateMachines"),
	}
}

// Start starts related state machines for broker.
func (f *StateMachineFactory) Start() (err error) {
	f.logger.Debug("starting StorageConfigStateMachine")
	sm, err := f.createStorageConfigStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)

	f.logger.Debug("starting ReplicaLeaderStateMachine")
	sm, err = f.createShardAssignStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)
	f.logger.Debug("starting ShardAssignmentStateMachine")
	sm, err = f.createReplicaLeaderStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)

	f.logger.Info("started MasterStateMachines")
	return nil
}

// Stop stops the broker's state machines.
func (f *StateMachineFactory) Stop() {
	for _, sm := range f.stateMachines {
		if err := sm.Close(); err != nil {
			f.logger.Error("close state machine error", logger.Error(err))
		}
	}
}

func (f *StateMachineFactory) createStorageConfigStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.StorageConfigStateMachine,
		f.discoveryFactory,
		constants.StorageConfigPath,
		true,
		f.stateMgr.OnStorageConfigChange,
		f.stateMgr.OnStorageConfigDelete,
	)
}

func (f *StateMachineFactory) createShardAssignStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.ShardAssignmentStateMachine,
		f.discoveryFactory,
		constants.DatabaseConfigPath,
		true,
		f.stateMgr.OnDatabaseCfgChange,
		f.stateMgr.OnDatabaseCfgDelete,
	)
}

func (f *StateMachineFactory) createReplicaLeaderStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.ReplicaLeaderStateMachine,
		f.discoveryFactory,
		constants.ShardAssigmentPath,
		true,
		f.stateMgr.OnShardAssignmentChange,
		f.stateMgr.OnShardAssignmentDelete,
	)
}

func (f *StateMachineFactory) createStorageNodeStateMachine(storageName string,
	discoveryFactory discovery.Factory,
) (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.StorageNodeStateMachine,
		discoveryFactory,
		constants.LiveNodesPath,
		true,
		func(key string, resource []byte) {
			f.stateMgr.OnStorageNodeStartup(storageName, key, resource)
		},
		func(key string) {
			f.stateMgr.OnStorageNodeFailure(storageName, key)
		},
	)
}
