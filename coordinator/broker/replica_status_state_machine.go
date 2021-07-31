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
	"path/filepath"
	"sort"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/inif"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source=./replica_status_state_machine.go -destination=./replica_status_state_machine_mock.go -package=broker

// ReplicaStatusStateMachine represents the status of database's replicas
// Each broker node need start this state machine,
type ReplicaStatusStateMachine interface {
	inif.Listener
	io.Closer

	// GetQueryableReplicas returns the queryable replicas，
	// and chooses the fastest replica if the shard has multi-replica.
	// returns storage node => shard id list
	GetQueryableReplicas(database string) map[string][]int32
	// GetReplicas returns the replica state list under this broker by broker's indicator
	GetReplicas(broker string) models.BrokerReplicaState
}

// replicaStatusStateMachine implements status state machine,
// watches replica state path for listening modify event which broker uploaded
type replicaStatusStateMachine struct {
	discovery discovery.Discovery

	ctx    context.Context
	cancel context.CancelFunc

	running *atomic.Bool
	mutex   sync.RWMutex
	// brokers: broker node => replica list under this broker
	brokers map[string]models.BrokerReplicaState

	logger *logger.Logger
}

// NewReplicaStatusStateMachine creates a replica's status state machine
func NewReplicaStatusStateMachine(ctx context.Context, factory discovery.Factory) (ReplicaStatusStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	sm := &replicaStatusStateMachine{
		running: atomic.NewBool(false),
		brokers: make(map[string]models.BrokerReplicaState),
		logger:  logger.GetLogger("coordinator", "ReplicaStatusStateMachine"),
		ctx:     c,
		cancel:  cancel,
	}

	// new replica status discovery
	sm.discovery = factory.CreateDiscovery(constants.ReplicaStatePath, sm)
	if err := sm.discovery.Discovery(true); err != nil {
		return nil, fmt.Errorf("discovery database status error:%s", err)
	}

	sm.running.Store(true)
	sm.logger.Info("replica status state machine is started")

	return sm, nil
}

// GetQueryableReplicas returns the queryable replicas
// returns storage node => shard id list
func (sm *replicaStatusStateMachine) GetQueryableReplicas(database string) map[string][]int32 {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if !sm.running.Load() {
		return nil
	}

	// 1. find shards by given database's name
	shards := make(map[string][]models.ReplicaState)

	for _, brokerReplicaState := range sm.brokers {
		for _, replica := range brokerReplicaState.Replicas {
			if replica.Database != database {
				continue
			}
			shardID := replica.ShardIndicator()
			shards[shardID] = append(shards[shardID], replica)
		}
	}

	if len(shards) == 0 {
		return nil
	}

	result := make(map[string][]int32)
	for _, replicas := range shards {
		replicaList := replicas
		if len(replicaList) > 1 {
			// has multi-replica, chooses the fastest
			// sort replicas based pending msg
			sort.Slice(replicaList, func(i, j int) bool {
				return replicaList[i].Pending < replicaList[j].Pending
			})
		}
		nodeID := replicaList[0].Target.Indicator()
		result[nodeID] = append(result[nodeID], replicaList[0].ShardID)
	}

	return result
}

// GetReplicas returns the replica state list under this broker by broker's indicator
func (sm *replicaStatusStateMachine) GetReplicas(broker string) models.BrokerReplicaState {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if !sm.running.Load() {
		return models.BrokerReplicaState{}
	}

	return sm.brokers[broker]
}

// Close closes state machine, stops watch change event
func (sm *replicaStatusStateMachine) Close() error {
	if sm.running.CAS(true, false) {
		defer sm.cancel()

		sm.discovery.Close()
		sm.logger.Info("replica status state machine is stopped.")
	}
	return nil
}

// OnCreate updates the broker's replica status when broker upload replica state
func (sm *replicaStatusStateMachine) OnCreate(key string, resource []byte) {
	sm.logger.Debug("discovery new broker online",
		logger.String("key", key),
		logger.String("data", string(resource)))

	brokerReplicaState := models.BrokerReplicaState{}
	if err := encoding.JSONUnmarshal(resource, &brokerReplicaState); err != nil {
		sm.logger.Error("discovery replica status but unmarshal error", logger.Error(err))
		return
	}
	_, broker := filepath.Split(key)

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.brokers[broker] = brokerReplicaState
}

// OnDelete deletes the broker's replica status when broker offline.
func (sm *replicaStatusStateMachine) OnDelete(key string) {
	sm.logger.Info("discovery broker offline remove",
		logger.String("key", key))

	_, broker := filepath.Split(key)
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.brokers, broker)
}
