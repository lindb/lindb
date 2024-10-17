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

package storage

import (
	"context"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
)

// StateMachinePaths represents the paths which storage state machine need watch.
var StateMachinePaths = make(map[string]models.StateMachineInfo)

func init() {
	StateMachinePaths[constants.LiveNode] = models.StateMachineInfo{
		Path:    constants.StorageLiveNodesPath,
		Comment: "Storage alive nodes",
		CreateState: func() interface{} {
			return &models.StatefulNode{}
		},
	}
	StateMachinePaths[constants.ShardAssignment] = models.StateMachineInfo{
		Path:    constants.ShardAssignmentPath,
		Comment: "Database shard assignment.",
		CreateState: func() interface{} {
			return &models.ShardAssignment{}
		},
	}
}

// StateMachineFactory represents storage state machine maintainer.
type StateMachineFactory struct {
	ctx              context.Context
	discoveryFactory discovery.Factory
	stateMgr         StateManager
	logger           logger.Logger
	stateMachines    []discovery.StateMachine
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
		logger:           logger.GetLogger("Storage", "StateMachineFactory"),
	}
}

// Start starts all storage's related state machines.
func (f *StateMachineFactory) Start() (err error) {
	f.logger.Debug("starting LiveNodeStateMachine")
	sm, err := f.createStorageLiveNodeStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)

	f.logger.Debug("starting ShardAssignStateMachine")
	sm, err = f.createShardAssignStateMachine()
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

	f.logger.Info("started StorageStateMachines")
	return nil
}

// Stop stops all storage's related state machines.
func (f *StateMachineFactory) Stop() {
	for _, sm := range f.stateMachines {
		if err := sm.Close(); err != nil {
			f.logger.Error("close state machine error", logger.Error(err))
		}
	}
}

// createShardAssignStateMachine creates shard assignment state machine.
func (f *StateMachineFactory) createShardAssignStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.ShardAssignmentStateMachine,
		f.discoveryFactory,
		constants.ShardAssignmentPath,
		true,
		f.onShardAssignmentChange,
		nil,
	)
}

// createStorageLiveNodeStateMachine creates storage live node state machine.
func (f *StateMachineFactory) createStorageLiveNodeStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.LiveNodeStateMachine,
		f.discoveryFactory,
		constants.StorageLiveNodesPath,
		true,
		f.onNodeStartup,
		f.onNodeFailure,
	)
}

// onNodeStartup triggers when storage node online.
func (f *StateMachineFactory) onNodeStartup(key string, data []byte) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type:  discovery.NodeStartup,
		Key:   key,
		Value: data,
	})
}

// onNodeFailure triggers when storage node offline.
func (f *StateMachineFactory) onNodeFailure(key string) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type: discovery.NodeFailure,
		Key:  key,
	})
}

// onShardAssignmentChange triggers when shard assignment changed after database config modified.
func (f *StateMachineFactory) onShardAssignmentChange(key string, data []byte) {
	f.stateMgr.EmitEvent(&discovery.Event{
		Type:  discovery.ShardAssignmentChanged,
		Key:   key,
		Value: data,
	})
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
