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

package root

import (
	"context"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/pkg/logger"
)

const brokerNameKey = "brokerName"

// stateMachineFactory represents root state matchine maintainer.
type stateMachineFactory struct {
	ctx              context.Context
	discoveryFactory discovery.Factory
	stateMgr         StateManager
	stateMachines    []discovery.StateMachine

	logger *logger.Logger
}

// NewStateMachineFactory creates a StateMachineFactory instance.
func NewStateMachineFactory(
	ctx context.Context,
	discoveryFactory discovery.Factory,
	stateMgr StateManager,
) discovery.StateMachineFactory {
	return &stateMachineFactory{
		ctx:              ctx,
		discoveryFactory: discoveryFactory,
		stateMgr:         stateMgr,
		logger:           logger.GetLogger("Root", "StateMachineFactory"),
	}
}

// Start starts all root's related state machines.
func (f *stateMachineFactory) Start() error {
	f.logger.Debug("starting BrokerConfigStateMachine")
	sm, err := f.createBrokerConfigStateMachine()
	if err != nil {
		return err
	}
	f.stateMachines = append(f.stateMachines, sm)
	return nil
}

// Stop stops all root's related state machines.
func (f *stateMachineFactory) Stop() {
	f.logger.Info("stopping root state machines...")
	for _, sm := range f.stateMachines {
		if err := sm.Close(); err != nil {
			f.logger.Error("close state machine error", logger.Error(err))
		}
	}
}

// createBrokerConfigStateMachine creates broker config state machine.
func (f *stateMachineFactory) createBrokerConfigStateMachine() (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx,
		discovery.BrokerConfigStateMachine,
		f.discoveryFactory,
		constants.BrokerConfigPath,
		true,
		func(key string, data []byte) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:  discovery.BrokerConfigChanged,
				Key:   key,
				Value: data,
			})
		},
		func(key string) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type: discovery.BrokerConfigDeletion,
				Key:  key,
			})
		},
	)
}

// createBrokerNodeStateMachine creates broker node state machine.
func (f *stateMachineFactory) createBrokerNodeStateMachine(
	brokerName string,
	discoveryFactory discovery.Factory,
) (discovery.StateMachine, error) {
	return discovery.NewStateMachine(
		f.ctx, discovery.BrokerNodeStateMachine,
		discoveryFactory,
		constants.LiveNodesPath,
		true,
		func(key string, data []byte) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:       discovery.NodeStartup,
				Key:        key,
				Value:      data,
				Attributes: map[string]string{brokerNameKey: brokerName},
			})
		},
		func(key string) {
			f.stateMgr.EmitEvent(&discovery.Event{
				Type:       discovery.NodeFailure,
				Key:        key,
				Attributes: map[string]string{brokerNameKey: brokerName},
			})
		})
}
