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
	"errors"
	"fmt"
	"sync"

	coCtx "github.com/lindb/lindb/coordinator/context"
	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/elect"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

//go:generate mockgen -source=./master.go -destination=./master_mock.go -package=coordinator

var log = logger.GetLogger("coordinator", "Master")

var (
	errNoCluster = errors.New("cluster not exist")
)

// MasterCfg represents the config for master creating
type MasterCfg struct {
	// basic
	Ctx  context.Context
	TTL  int64 // master elect keepalive ttl
	Node models.Node
	Repo state.Repository

	// factory
	DiscoveryFactory  discovery.Factory
	ControllerFactory task.ControllerFactory
	ClusterFactory    storage.ClusterFactory
	RepoFactory       state.RepositoryFactory

	// service
	StorageStateService service.StorageStateService
	ShardAssignService  service.ShardAssignService

	// broker state machine
	BrokerSM *BrokerStateMachines
}

// Master represents all metadata/state controller, only has one active master in broker cluster.
// Master will control all storage cluster metadata, update state, then notify each broker node.
type Master interface {
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
}

// master implements master interface
type master struct {
	cfg *MasterCfg

	// create by runtime
	masterCtx *coCtx.MasterContext
	ctx       context.Context
	cancel    context.CancelFunc
	elect     elect.Election

	mutex sync.Mutex
}

// NewMaster create master for current node
func NewMaster(cfg *MasterCfg) Master {
	ctx, cancel := context.WithCancel(cfg.Ctx)
	m := &master{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}
	// create master election
	m.elect = elect.NewElection(ctx, cfg.Repo, cfg.Node, cfg.TTL, m)
	return m
}

// OnFailOver invoked after master electing, current node become a new master
func (m *master) OnFailOver() error {
	log.Info("starting master fail over")
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var err error
	stateMachine := &coCtx.StateMachine{}
	newCtx := coCtx.NewMasterContext(stateMachine)
	defer func() {
		if err != nil {
			newCtx.Close()
			m.masterCtx = nil
		} else {
			m.masterCtx = newCtx
		}
	}()

	stateMachine.StorageCluster, err = storage.NewClusterStateMachine(m.ctx, m.cfg.Repo,
		m.cfg.ControllerFactory, m.cfg.DiscoveryFactory, m.cfg.ClusterFactory, m.cfg.RepoFactory,
		m.cfg.StorageStateService, m.cfg.ShardAssignService)
	if err != nil {
		return fmt.Errorf("start storage cluster state machine errer:%s", err)
	}

	stateMachine.DatabaseAdmin, err = database.NewAdminStateMachine(m.ctx, m.cfg.DiscoveryFactory, stateMachine.StorageCluster)
	if err != nil {
		return fmt.Errorf("start database admin state machine error:%s", err)
	}

	return nil
}

// OnResignation invoked current node is master, before re-electing
func (m *master) OnResignation() {
	log.Info("starting master resign")
	if m.masterCtx != nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		m.masterCtx.Close()
		m.masterCtx = nil
	}
}

// IsMaster returns current node if is master
func (m *master) IsMaster() bool {
	return m.elect.IsMaster()
}

// GetMaster returns the current master info
func (m *master) GetMaster() *models.Master {
	return m.elect.GetMaster()
}

// Start starts master do election master, if success build master context,
// starts state machine do cluster coordinate such metadata, cluster state etc.
func (m *master) Start() {
	m.elect.Initialize()
	m.elect.Elect()
}

// Stop stops master if current node is master, cleanup master context and stops state machine
func (m *master) Stop() {
	// close master elect
	m.elect.Close()

	m.cancel()
}

// FlushDatabase submits the coordinator task for flushing memory database by cluster and database name
func (m *master) FlushDatabase(cluster string, databaseName string) error {
	if m.IsMaster() {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		cluster := m.masterCtx.StateMachine.StorageCluster.GetCluster(cluster)
		if cluster == nil {
			return errNoCluster
		}
		return cluster.FlushDatabase(databaseName)
	}
	return nil
}
