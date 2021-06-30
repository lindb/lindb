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

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/inif"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./storage_state_machine.go -destination=./storage_state_machine_mock.go -package=broker

// StorageStateMachine represents storage cluster state state machine.
// Each broker node will start this state machine which watch storage cluster state change event.
type StorageStateMachine interface {
	inif.Listener
	io.Closer

	// List lists currently all alive storage cluster's state
	List() []*models.StorageState
}

// storageStateMachine implements StorageStateMachine interface.
type storageStateMachine struct {
	discovery         discovery.Discovery
	ctx               context.Context
	cancel            context.CancelFunc
	taskClientFactory rpc.TaskClientFactory

	storageClusters map[string]*StorageClusterState

	mutex   sync.RWMutex
	running *atomic.Bool

	logger *logger.Logger
}

// NewStorageStateMachine creates state machine, init data if exist, then starts watch change event
func NewStorageStateMachine(ctx context.Context,
	discoveryFactory discovery.Factory, taskClientFactory rpc.TaskClientFactory) (StorageStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	stateMachine := &storageStateMachine{
		taskClientFactory: taskClientFactory,
		ctx:               c,
		cancel:            cancel,
		storageClusters:   make(map[string]*StorageClusterState),
		running:           atomic.NewBool(false),
		logger:            logger.GetLogger("coordinator", "StorageStateMachine"),
	}

	// new storage config discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.StorageClusterNodeStatePath, stateMachine)
	if err := stateMachine.discovery.Discovery(true); err != nil {
		return nil, fmt.Errorf("discovery storage cluster state error:%s", err)
	}

	stateMachine.running.Store(true)
	stateMachine.logger.Info("storage state machine is started")
	return stateMachine, nil
}

// List lists currently all alive storage cluster's state
func (s *storageStateMachine) List() (rs []*models.StorageState) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.running.Load() {
		return
	}

	for _, storageState := range s.storageClusters {
		rs = append(rs, storageState.state)
	}

	return
}

// OnCreate modifies storage cluster's state, such as trigger by storage create event.
func (s *storageStateMachine) OnCreate(key string, resource []byte) {
	s.logger.Info("discovery new storage cluster create",
		logger.String("key", key),
		logger.String("data", string(resource)))

	storageState := models.NewStorageState()
	if err := json.Unmarshal(resource, storageState); err != nil {
		s.logger.Error("discovery new storage state but unmarshal error",
			logger.String("data", string(resource)), logger.Error(err))
		return
	}
	if len(storageState.Name) == 0 {
		s.logger.Error("cluster name is empty")
		return
	}
	s.logger.Info("storage cluster state change", logger.String("cluster", storageState.Name))
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//TODO need check if same state, maybe state is same, such as system start
	state, ok := s.storageClusters[storageState.Name]
	if !ok {
		state = newStorageClusterState(s.taskClientFactory)
		s.storageClusters[storageState.Name] = state
	}
	state.SetState(storageState)
}

// OnDelete deletes storage cluster's state when cluster offline
func (s *storageStateMachine) OnDelete(key string) {
	_, name := filepath.Split(key)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	state, ok := s.storageClusters[name]
	if ok {
		state.close()
		delete(s.storageClusters, name)
	}
}

// Close closes state machine, stops watch change event
func (s *storageStateMachine) Close() error {
	if s.running.CAS(true, false) {
		s.mutex.Lock()
		defer func() {
			s.cancel()
			s.mutex.Unlock()
		}()

		s.discovery.Close()

		// close all storage cluster's states
		for _, state := range s.storageClusters {
			state.close()
		}

		s.logger.Info("storage cluster state machine is stopped.")
	}
	return nil
}
