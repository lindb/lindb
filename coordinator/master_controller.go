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

package coordinator

import (
	"context"
	"fmt"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/elect"
	masterpkg "github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./master_controller.go -destination=./master_controller_mock.go -package=coordinator

var log = logger.GetLogger("coordinator", "MasterController")

// MasterCfg represents the config for masterController creating
type MasterCfg struct {
	// basic
	Ctx  context.Context
	TTL  int64 // masterController elect keepalive ttl
	Node models.Node
	Repo state.Repository

	// factory
	DiscoveryFactory discovery.Factory
	RepoFactory      state.RepositoryFactory
}

// MasterController represents all metadata/state controller, only has one active master in broker cluster.
// MasterController will control all storage cluster metadata, update state, then notify each broker node.
type MasterController interface {
	// Start starts master do election master, if success build master context,
	// starts state machine do cluster coordinate such metadata, cluster state etc.
	Start()
	// IsMaster returns current node if is master
	IsMaster() bool
	// GetMaster returns the current master info
	GetMaster() *models.Master
	// Stop stops master if current node is master, cleanup master context and stops state machine
	Stop()
	// FlushDatabase submits the coordinator task for flushing memory database by cluster and database name
	FlushDatabase(cluster string, databaseName string) error
	// GetStateManager returns master's state manager.
	GetStateManager() masterpkg.StateManager
}

// masterController implements MasterController interface
type masterController struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg      *MasterCfg
	stateMgr masterpkg.StateManager

	// create by runtime
	stateMachineFct *masterpkg.StateMachineFactory
	elect           elect.Election

	mutex sync.Mutex
}

// NewMasterController create MasterController for current node
func NewMasterController(cfg *MasterCfg) MasterController {
	ctx, cancel := context.WithCancel(cfg.Ctx)
	m := &masterController{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}
	// create master election
	m.elect = elect.NewElection(ctx, cfg.Repo, cfg.Node, cfg.TTL, m)
	return m
}

// OnFailOver invoked after master electing, current node become a new master
func (m *masterController) OnFailOver() error {
	log.Info("starting master fail over")
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var err error
	stateMgr := masterpkg.NewStateManager(m.ctx, m.cfg.Repo, m.cfg.RepoFactory)
	stateMachineFct := masterpkg.NewStateMachineFactory(m.ctx, m.cfg.DiscoveryFactory, stateMgr)
	// first need set state machine factory in state manager
	stateMgr.SetStateMachineFactory(stateMachineFct)

	defer func() {
		if err != nil {
			stateMachineFct.Stop()
			stateMgr.Close()
			m.stateMgr = nil
			m.stateMachineFct = nil
		} else {
			m.stateMachineFct = stateMachineFct
			m.stateMgr = stateMgr
		}
	}()
	// start master state machine
	err = stateMachineFct.Start()
	if err != nil {
		return fmt.Errorf("start master state machine error:%s", err)
	}

	return nil
}

// OnResignation invoked current node is master, before re-electing
func (m *masterController) OnResignation() {
	log.Info("starting master resign")
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.stateMachineFct != nil {
		m.stateMachineFct.Stop()
		m.stateMachineFct = nil
	}

	if m.stateMgr != nil {
		m.stateMgr.Close()
	}
}

// IsMaster returns current node if is master
func (m *masterController) IsMaster() bool {
	return m.elect.IsMaster()
}

// GetMaster returns the current master info
func (m *masterController) GetMaster() *models.Master {
	return m.elect.GetMaster()
}

func (m *masterController) GetStateManager() masterpkg.StateManager {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.stateMgr
}

// Start starts master do election master, if success build master context,
// starts state machine do cluster coordinate such metadata, cluster state etc.
func (m *masterController) Start() {
	m.elect.Initialize()
	m.elect.Elect()
}

// Stop stops master if current node is master, cleanup master context and stops state machine
func (m *masterController) Stop() {
	// close master elect
	m.elect.Close()

	m.cancel()
}

// FlushDatabase submits the coordinator task for flushing memory database by cluster and database name
func (m *masterController) FlushDatabase(cluster string, databaseName string) error {
	if m.IsMaster() {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		storage := m.stateMgr.GetStorageCluster(cluster)
		if storage == nil {
			return constants.ErrNoStorageCluster
		}
		return storage.FlushDatabase(databaseName)
	}
	return nil
}
