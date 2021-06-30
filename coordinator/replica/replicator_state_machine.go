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

package replica

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/inif"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/replication"
)

//go:generate mockgen -source=./replicator_state_machine.go -destination=./replicator_state_machine_mock.go -package=replica

// ReplicatorStateMachine represents the replicator state machine in broker
type ReplicatorStateMachine interface {
	inif.Listener
	io.Closer
}

// replicatorStateMachine implements the state machine interface,
// watches shard assignment change event, then builds replicators
type replicatorStateMachine struct {
	discovery discovery.Discovery
	cm        replication.ChannelManager

	mutex   sync.RWMutex
	running *atomic.Bool
	// shardAssigns: db's name => shard assignment
	shardAssigns map[string]*models.ShardAssignment

	ctx    context.Context
	cancel context.CancelFunc

	logger *logger.Logger
}

// NewReplicatorStateMachine creates the state machine
func NewReplicatorStateMachine(ctx context.Context,
	cm replication.ChannelManager,
	discoveryFactory discovery.Factory,
) (ReplicatorStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	stateMachine := &replicatorStateMachine{
		ctx:          c,
		cancel:       cancel,
		cm:           cm,
		shardAssigns: make(map[string]*models.ShardAssignment),
		running:      atomic.NewBool(false),
		logger:       logger.GetLogger("coordinator", "ReplicatorStateMachine"),
	}
	// new database's shard assign discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.DatabaseAssignPath, stateMachine)
	if err := stateMachine.discovery.Discovery(true); err != nil {
		return nil, fmt.Errorf("discovery replicator error:%s", err)
	}

	stateMachine.running.Store(true)
	// final sync new replicator state
	stateMachine.cm.SyncReplicatorState()
	stateMachine.logger.Info("replicator state machine is started")

	return stateMachine, nil
}

// OnCreate triggers on shard assignment creation, builds related replicators
func (sm *replicatorStateMachine) OnCreate(key string, resource []byte) {
	sm.logger.Info("discovery new database shard assignment create in cluster",
		logger.String("key", key),
		logger.String("data", string(resource)))

	shardAssign := &models.ShardAssignment{}
	if err := encoding.JSONUnmarshal(resource, shardAssign); err != nil {
		sm.logger.Error("unmarshal shard assign", logger.Error(err))
		return
	}
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.buildShardAssign(shardAssign)

	if sm.running.Load() {
		// final sync new replicator state
		sm.cm.SyncReplicatorState()
	}
}

// OnDelete trigger on database deletion, destroy related replicators for deletion database
func (sm *replicatorStateMachine) OnDelete(key string) {
	sm.logger.Info("discovery a database shard assignment delete from cluster",
		logger.String("key", key))

	_, dbName := filepath.Split(key)

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	//FIXME: need remove replicator and channel when database delete?
	delete(sm.shardAssigns, dbName)
}

// Close closes the state machine
func (sm *replicatorStateMachine) Close() error {
	if sm.running.CAS(true, false) {
		sm.mutex.Lock()
		defer func() {
			sm.mutex.Unlock()
			sm.cancel()
		}()
		sm.discovery.Close()

		sm.logger.Info("replicator state machine is stopped")
	}
	return nil
}

// buildShardAssign builds the wal replica channel and related replicators for the shard assignment
func (sm *replicatorStateMachine) buildShardAssign(shardAssign *models.ShardAssignment) {
	sm.shardAssigns[shardAssign.Name] = shardAssign
	shards := shardAssign.Shards

	numOfShard := len(shards)
	for shardID := range shards {
		sm.createReplicaChannel(numOfShard, shardID, shardAssign)
	}
}

// createReplicaChannel creates wal replica channel for spec database's shard
func (sm *replicatorStateMachine) createReplicaChannel(numOfShard, shardID int, shardAssign *models.ShardAssignment) {
	db := shardAssign.Name
	ch, err := sm.cm.CreateChannel(db, int32(numOfShard), int32(shardID))
	if err != nil {
		sm.logger.Error("create replica channel", logger.Error(err))
		return
	}
	sm.logger.Info("create replica channel successfully", logger.String("db", db), logger.Any("shardID", shardID))

	sm.startReplicator(ch, shardID, shardAssign)
}

// startReplicator starts wal replicator for spec database's shard
func (sm *replicatorStateMachine) startReplicator(ch replication.Channel, shardID int, shardAssign *models.ShardAssignment) {
	replica := shardAssign.Shards[shardID]
	db := shardAssign.Name

	for _, replicaID := range replica.Replicas {
		target := shardAssign.Nodes[replicaID]
		if target != nil {
			_, err := ch.GetOrCreateReplicator(*target)
			if err != nil {
				sm.logger.Error("start replicator", logger.Error(err))
				continue
			}
			sm.logger.Info("create replicator successfully", logger.String("db", db),
				logger.Any("shardID", shardID), logger.String("target", target.Indicator()))
		}
	}
}
