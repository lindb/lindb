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
	"context"
	"fmt"
	"io"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source=./state_machine.go -destination=./state_machine_mock.go -package=discovery

// NewStateMachineFn represents new state machine function.
var NewStateMachineFn = NewStateMachine

// StateMachineType represents state machine type.
type StateMachineType int

const (
	DatabaseConfigStateMachine StateMachineType = iota + 1
	ShardAssignmentStateMachine
	LiveNodeStateMachine
	StorageStatusStateMachine
	StorageConfigStateMachine
	StorageNodeStateMachine
	BrokerConfigStateMachine
	BrokerNodeStateMachine
)

// String returns state machine type desc.
func (st StateMachineType) String() string {
	switch st {
	case DatabaseConfigStateMachine:
		return "DatabaseConfigStateMachine"
	case ShardAssignmentStateMachine:
		return "ShardAssignmentStateMachine"
	case LiveNodeStateMachine:
		return "LiveNodeStateMachine"
	case StorageStatusStateMachine:
		return "StorageStatusStateMachine"
	case StorageConfigStateMachine:
		return "StorageConfigStateMachine"
	case StorageNodeStateMachine:
		return "StorageNodeStateMachine"
	case BrokerConfigStateMachine:
		return "BrokerConfigStateMachine"
	case BrokerNodeStateMachine:
		return "BrokerNodeStateMachine"
	default:
		return "Unknown"
	}
}

// StateMachineEventHandle represents handle state machine event.
type StateMachineEventHandle interface {
	// EmitEvent emits discovery event when state changed.
	EmitEvent(event *Event)
	// Close cleans the resource.
	Close()
}

// StateMachineFactory represents maintain all state machines for each role.
type StateMachineFactory interface {
	// Start starts all state machines, do init logic.
	Start() error
	// Stop stops all state machines, clean all resources.
	Stop()
}

// Listener represents discovery resource event callback interface,
// includes create/delete/cleanup operation.
type Listener interface {
	// OnCreate is resource creation callback.
	OnCreate(key string, resource []byte)
	// OnDelete is resource deletion callback.
	OnDelete(key string)
}

// StateMachine represents state changed event state machine.
// Like node online/offline, database create events etc.
type StateMachine interface {
	Listener
	io.Closer
}

// stateMachine implements StateMachine interface.
type stateMachine struct {
	ctx    context.Context
	cancel context.CancelFunc

	stateMachineType StateMachineType
	discovery        Discovery

	onCreateFn func(key string, resource []byte)
	onDeleteFn func(key string)

	running *atomic.Bool

	logger *logger.Logger
}

// NewStateMachine creates a state machine instance.
func NewStateMachine(ctx context.Context,
	stateMachineType StateMachineType,
	discoveryFactory Factory,
	watchPath string,
	needInitialize bool,
	onCreateFn func(key string, resource []byte),
	onDeleteFn func(key string),
) (StateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	stateMachine := &stateMachine{
		ctx:              c,
		cancel:           cancel,
		stateMachineType: stateMachineType,
		onCreateFn:       onCreateFn,
		onDeleteFn:       onDeleteFn,
		running:          atomic.NewBool(true),
		logger:           logger.GetLogger("Coordinator", "StateMachine"),
	}

	// new state discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(watchPath, stateMachine)
	if err := stateMachine.discovery.Discovery(needInitialize); err != nil {
		return nil, fmt.Errorf("discovery state error:%s", err)
	}

	stateMachine.logger.Info("state machine start successfully",
		logger.String("type", stateMachineType.String()))
	return stateMachine, nil
}

// OnCreate watches state changed, such as node online event.
func (sm *stateMachine) OnCreate(key string, resource []byte) {
	if !sm.running.Load() {
		sm.logger.Warn("state machine is stopped",
			logger.String("type", sm.stateMachineType.String()))
		return
	}
	sm.logger.Info("discovery new state",
		logger.String("type", sm.stateMachineType.String()),
		logger.String("key", key))

	if sm.onCreateFn != nil {
		sm.onCreateFn(key, resource)
	}
}

// OnDelete watches state deleted, such as node offline event.
func (sm *stateMachine) OnDelete(key string) {
	if !sm.running.Load() {
		sm.logger.Warn("state machine is stopped",
			logger.String("type", sm.stateMachineType.String()))
		return
	}
	sm.logger.Info("discovery state removed",
		logger.String("type", sm.stateMachineType.String()),
		logger.String("key", key))
	if sm.onDeleteFn != nil {
		sm.onDeleteFn(key)
	}
}

// Close closes state machine, stops watch change event.
func (sm *stateMachine) Close() error {
	if sm.running.CAS(true, false) {
		defer func() {
			sm.cancel()
		}()

		sm.discovery.Close()

		sm.logger.Info("state machine stop successfully",
			logger.String("type", sm.stateMachineType.String()))
	}
	return nil
}
