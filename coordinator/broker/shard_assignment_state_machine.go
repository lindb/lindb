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
	"fmt"
	"io"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/inif"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./shard_assignment_state_machine.go -destination=./shard_assignment_state_machine_mock.go -package=broker

// ShardAssignmentStateMachine is database config controller,
// creates shard assignment based on config and active nodes related storage cluster.
// runtime watches database change event, maintain shard assignment and create related coordinator task.
type ShardAssignmentStateMachine interface {
	inif.Listener
	io.Closer
}

// shardAssignmentStateMachine implement ShardAssignmentStateMachine interface.
// all metadata change will store related storage cluster.
type shardAssignmentStateMachine struct {
	storageCluster storage.ClusterStateMachine
	discovery      discovery.Discovery

	mutex   sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
	running *atomic.Bool

	logger *logger.Logger
}

// NewShardAssignmentStateMachine creates shard assignment state machine instance
func NewShardAssignmentStateMachine(ctx context.Context, discoveryFactory discovery.Factory,
	storageCluster storage.ClusterStateMachine) (ShardAssignmentStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	// new shard assignment state machine instance
	stateMachine := &shardAssignmentStateMachine{
		storageCluster: storageCluster,
		ctx:            c,
		running:        atomic.NewBool(false),
		cancel:         cancel,
		logger:         logger.GetLogger("coordinator", "ShardAssignmentStateMachine"),
	}
	// new database config discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.DatabaseConfigPath, stateMachine)
	if err := stateMachine.discovery.Discovery(false); err != nil {
		return nil, fmt.Errorf("discovery database config error:%s", err)
	}
	stateMachine.running.Store(true)
	stateMachine.logger.Info("database shard assignment state machine is started")
	return stateMachine, nil
}

// OnCreate creates shard assignment when receive database create event
func (sm *shardAssignmentStateMachine) OnCreate(key string, resource []byte) {
	sm.logger.Info("discovery new database need shard assignment in cluster",
		logger.String("key", key),
		logger.String("data", string(resource)))

	cfg := models.Database{}
	if err := encoding.JSONUnmarshal(resource, &cfg); err != nil {
		sm.logger.Error("discovery database create but unmarshal error", logger.Error(err))
		return
	}

	if len(cfg.Name) == 0 {
		sm.logger.Error("database name cannot be empty")
		return
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	cluster := sm.storageCluster.GetCluster(cfg.Cluster)
	if cluster == nil {
		sm.logger.Error("storage cluster not exist",
			logger.String("cluster", cfg.Cluster))
		return
	}
	shardAssign, err := cluster.GetShardAssign(cfg.Name)
	if err != nil && err != state.ErrNotExist {
		sm.logger.Error("get shard assign error", logger.Error(err))
		return
	}
	// build shard assignment for creation database, generate related coordinator task
	if shardAssign == nil {
		if err := sm.createShardAssignment(cfg.Name, cluster, &cfg, -1, -1); err != nil {
			sm.logger.Error("create shard assignment error", logger.Error(err))
		}
	} else if len(shardAssign.Shards) != cfg.NumOfShard {
		if err := sm.modifyShardAssignment(cfg.Name, shardAssign, cluster, &cfg); err != nil {
			sm.logger.Error("modify shard assignment error", logger.Error(err))
		}
	}
}

func (sm *shardAssignmentStateMachine) OnDelete(key string) {
	//TODO impl delete database???
}

// Close closes shard assignment state machine, stops watch change event.
func (sm *shardAssignmentStateMachine) Close() error {
	if sm.running.CAS(true, false) {
		defer sm.cancel()
		sm.discovery.Close()

		sm.logger.Info("shard assignment state machine is stopped.")
	}
	return nil
}

// createShardAssignment creates shard assignment for spec cluster
// 1) generate shard assignment
// 2) save shard assignment into related storage cluster
// 3) submit create shard coordinator task(storage node will execute it when receive task event)
func (sm *shardAssignmentStateMachine) createShardAssignment(databaseName string,
	cluster storage.Cluster, cfg *models.Database, fixedStartIndex, startShardID int) error {
	activeNodes := cluster.GetActiveNodes()
	if len(activeNodes) == 0 {
		return fmt.Errorf("active node not found")
	}
	//TODO need calc resource and pick related node for store data
	var nodes = make(map[int]*models.Node)
	for idx, node := range activeNodes {
		nodes[idx] = &node.Node
	}

	var nodeIDs []int
	for idx := range nodes {
		nodeIDs = append(nodeIDs, idx)
	}

	// generate shard assignment based on node ids and config
	shardAssign, err := ShardAssignment(nodeIDs, cfg, fixedStartIndex, startShardID)
	if err != nil {
		return err
	}
	// set nodes and config, storage node will use it when execute create shard task
	shardAssign.Nodes = nodes

	sm.logger.Info("create shard assign",
		logger.String("database", databaseName),
		logger.Any("shardAssign", shardAssign))
	// save shard assignment into related storage cluster
	if err := cluster.SaveShardAssign(databaseName, shardAssign, cfg.Option); err != nil {
		return err
	}
	return nil
}

func (sm *shardAssignmentStateMachine) modifyShardAssignment(databaseName string, shardAssign *models.ShardAssignment,
	cluster storage.Cluster, cfg *models.Database) error {
	if len(shardAssign.Shards) > cfg.NumOfShard { //reduce shardAssign's shards
		//TODO implement the reduce shards, is needed?
		panic("not implemented")
	} else if len(shardAssign.Shards) < cfg.NumOfShard { //add shardAssign's shards
		activeNodes := cluster.GetActiveNodes()
		if len(activeNodes) == 0 {
			return fmt.Errorf("active node not found")
		}
		//TODO need calc resource and pick related node for store data
		var nodes = make(map[int]*models.Node)
		for idx, node := range activeNodes {
			nodes[idx] = &node.Node
		}

		var nodeIDs []int
		for idx := range nodes {
			nodeIDs = append(nodeIDs, idx)
		}

		// generate shard assignment based on node ids and config
		err := ModifyShardAssignment(nodeIDs, cfg, shardAssign, -1, len(shardAssign.Shards))
		if err != nil {
			return err
		}
	}
	sm.logger.Info("modify shard assign",
		logger.String("database", databaseName),
		logger.Any("shardAssign", shardAssign))
	// save shard assignment into related storage cluster
	if err := cluster.SaveShardAssign(databaseName, shardAssign, cfg.Option); err != nil {
		return err
	}
	return nil
}
