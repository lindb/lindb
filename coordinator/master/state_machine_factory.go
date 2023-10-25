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

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
)

const storageNameKey = "storageName"

// StateMachinePaths represents the paths which master state machine need watch.
var StateMachinePaths = make(map[string]models.StateMachineInfo)

func init() {
	StateMachinePaths[constants.Master] = models.StateMachineInfo{
		Path: constants.MasterPath,
		CreateState: func() interface{} {
			return &models.Master{}
		},
	}
	StateMachinePaths[constants.DatabaseConfig] = models.StateMachineInfo{
		Path: constants.DatabaseConfigPath,
		CreateState: func() interface{} {
			return &models.Database{}
		},
	}
	StateMachinePaths[constants.StorageConfig] = models.StateMachineInfo{
		Path: constants.StorageConfigPath,
		CreateState: func() interface{} {
			return &config.StorageCluster{}
		},
	}
	StateMachinePaths[constants.ShardAssignment] = models.StateMachineInfo{
		Path: constants.ShardAssignmentPath,
		CreateState: func() interface{} {
			return &models.ShardAssignment{}
		},
	}
	StateMachinePaths[constants.StorageState] = models.StateMachineInfo{
		Path: constants.StorageStatePath,
		CreateState: func() interface{} {
			return &models.StorageState{}
		},
	}
}

// StateMachineFactory represents master state machine maintainer.
type StateMachineFactory struct {
	ctx              context.Context
	discoveryFactory discovery.Factory
	stateMgr         StateManager

	stateMachines []discovery.StateMachine

	logger logger.Logger
}

// NewStateMachineFactory creates a StateMachineFactory instance.
func NewStateMachineFactory(ctx context.Context,
	discoveryFactory discovery.Factory,
	stateMgr StateManager,
) *StateMachineFactory {
	return &StateMachineFactory{
		ctx:              ctx,
		discoveryFactory: discoveryFactory,
		stateMgr:         stateMgr,
		logger:           logger.GetLogger("Master", "MasterStateMachines"),
	}
}

// Start starts all master related state machines.
func (f *StateMachineFactory) Start() (err error) {
	f.logger.Debug("starting StorageConfigStateMachine")
	sm, err := f.createStorageConfigStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)

	f.logger.Debug("starting DatabaseConfigStateMachine")
	sm, err = f.createDatabaseConfigStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)
	f.logger.Debug("starting ShardAssignmentStateMachine")
	sm, err = f.createShardAssignmentStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)

	f.logger.Debug("starting DatabaseLimitsStateMachine")
	sm, err = f.createDatabaseLimitsStateMachine()
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

// createStorageConfigStateMachine creates storage config state machine.
func (f *StateMachineFactory) createStorageConfigStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.StorageConfigStateMachine,
		f.discoveryFactory,
		constants.StorageConfigPath,
		true,
		func(key string, data []byte) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:  discovery.StorageConfigChanged,
				Key:   key,
				Value: data,
			})
		},
		func(key string) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type: discovery.StorageConfigDeletion,
				Key:  key,
			})
		},
	)
}

// createDatabaseConfigStateMachine crates database config state machine.
func (f *StateMachineFactory) createDatabaseConfigStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.DatabaseConfigStateMachine,
		f.discoveryFactory,
		constants.DatabaseConfigPath,
		true,
		func(key string, data []byte) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:  discovery.DatabaseConfigChanged,
				Key:   key,
				Value: data,
			})
		},
		func(key string) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type: discovery.DatabaseConfigDeletion,
				Key:  key,
			})
		})
}

// createShardAssignmentStateMachine creates shard assignment state machine.
func (f *StateMachineFactory) createShardAssignmentStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.ShardAssignmentStateMachine,
		f.discoveryFactory,
		constants.ShardAssignmentPath,
		true,
		func(key string, data []byte) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:  discovery.ShardAssignmentChanged,
				Key:   key,
				Value: data,
			})
		},
		func(key string) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type: discovery.ShardAssignmentDeletion,
				Key:  key,
			})
		})
}

// createStorageNodeStateMachine creates storage node state machine.
func (f *StateMachineFactory) createStorageNodeStateMachine(storageName string,
	discoveryFactory discovery.Factory,
) (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.StorageNodeStateMachine,
		discoveryFactory,
		constants.LiveNodesPath,
		true,
		func(key string, data []byte) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:       discovery.NodeStartup,
				Key:        key,
				Value:      data,
				Attributes: map[string]string{storageNameKey: storageName},
			})
		},
		func(key string) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:       discovery.NodeFailure,
				Key:        key,
				Attributes: map[string]string{storageNameKey: storageName},
			})
		},
	)
}

// createDatabaseLimitsStateMachine creates database's limits state machine.
func (f *StateMachineFactory) createDatabaseLimitsStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.DatabaseLimitsStateMachine,
		f.discoveryFactory,
		constants.DatabaseLimitPath,
		true,
		func(key string, data []byte) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:  discovery.DatabaseLimitsChanged,
				Key:   key,
				Value: data,
			})
		},
		nil,
	)
}
