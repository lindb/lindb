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

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/service"
)

//go:generate mockgen -source=./state_machine_factory.go -destination=./state_machine_factory_mock.go -package=coordinator

// StateMachineCfg represents the state machine config
type StateMachineCfg struct {
	Ctx               context.Context
	Repo              state.Repository
	CurrentNode       models.Node
	DiscoveryFactory  discovery.Factory
	ShardAssignSRV    service.ShardAssignService
	ChannelManager    replication.ChannelManager
	TaskClientFactory rpc.TaskClientFactory // rpc task stream create factory
}

// StateMachineFactory represents the state machine create factory
type StateMachineFactory interface {
	// CreateNodeStateMachine creates the node state machine
	CreateNodeStateMachine() (broker.NodeStateMachine, error)
	// CreateStorageStateMachine creates the storage state machine
	CreateStorageStateMachine() (broker.StorageStateMachine, error)
	// CreateReplicaStatusStateMachine creates the shard replica status state machine
	CreateReplicaStatusStateMachine() (replica.StatusStateMachine, error)
	// CreateReplicatorStateMachine creates the shard replicator state machine
	CreateReplicatorStateMachine() (replica.ReplicatorStateMachine, error)
	// CreateDatabaseStateMachine creates the database state machine
	CreateDatabaseStateMachine() (database.DBStateMachine, error)
}

// stateMachineFactory implements the interface, using state machine config for creating
type stateMachineFactory struct {
	cfg *StateMachineCfg
}

// NewStateMachineFactory creates the factory using config
func NewStateMachineFactory(cfg *StateMachineCfg) StateMachineFactory {
	return &stateMachineFactory{cfg: cfg}
}

// CreateNodeStateMachine creates the node state machine, if fail returns err
func (s *stateMachineFactory) CreateNodeStateMachine() (broker.NodeStateMachine, error) {
	return broker.NewNodeStateMachine(s.cfg.Ctx, s.cfg.CurrentNode, s.cfg.DiscoveryFactory, s.cfg.TaskClientFactory)
}

// CreateStorageStateMachine creates the storage state machine, if fail returns err
func (s *stateMachineFactory) CreateStorageStateMachine() (broker.StorageStateMachine, error) {
	return broker.NewStorageStateMachine(s.cfg.Ctx, s.cfg.DiscoveryFactory, s.cfg.TaskClientFactory)
}

// CreateReplicaStatusStateMachine creates the shard replica status state machine, if fail returns err
func (s *stateMachineFactory) CreateReplicaStatusStateMachine() (replica.StatusStateMachine, error) {
	return replica.NewStatusStateMachine(s.cfg.Ctx, s.cfg.DiscoveryFactory)
}

// CreateReplicatorStateMachine creates the shard replicator state machine
func (s *stateMachineFactory) CreateReplicatorStateMachine() (replica.ReplicatorStateMachine, error) {
	return replica.NewReplicatorStateMachine(s.cfg.Ctx, s.cfg.ChannelManager, s.cfg.ShardAssignSRV, s.cfg.DiscoveryFactory)
}

// CreateDatabaseStateMachine creates the database state machine
func (s *stateMachineFactory) CreateDatabaseStateMachine() (database.DBStateMachine, error) {
	return database.NewDBStateMachine(s.cfg.Ctx, s.cfg.Repo, s.cfg.DiscoveryFactory)
}
