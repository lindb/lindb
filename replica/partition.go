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
	"sync"
	"sync/atomic"

	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./partition.go -destination=./partition_mock.go -package=replica

var (
	// for testing
	newLocalReplicatorFn  = NewLocalReplicator
	newRemoteReplicatorFn = NewRemoteReplicator
)

const (
	replicatorTypeLocal  = "local"
	replicatorTypeRemote = "remote"
)

// Partition represents a partition of writeTask ahead log.
type Partition interface {
	io.Closer
	// BuildReplicaForLeader builds replica relation when handle writeTask connection.
	BuildReplicaForLeader(leader models.NodeID, replicas []models.NodeID) error
	// BuildReplicaForFollower builds replica relation when handle replica connection.
	BuildReplicaForFollower(leader models.NodeID, replica models.NodeID) error
	// ReplicaLog writes msg that leader sends replica msg.
	// return appended index, if success.
	ReplicaLog(replicaIdx int64, msg []byte) (int64, error)
	// WriteLog writes msg that leader handle client writeTask request.
	WriteLog(msg []byte) error
	// ReplicaAckIndex returns the index which replica appended index.
	ReplicaAckIndex() int64
	// ResetReplicaIndex resets replica index.
	ResetReplicaIndex(idx int64)
	// IsExpire returns partition if it is expired.
	IsExpire() bool
	// Path returns the path of partition.
	Path() string
	// Stop stops replicator channel.
	Stop()
	// getReplicaState returns each family's log replica state.
	getReplicaState() models.FamilyLogReplicaState
	// StartReplica iterates over all replicators and copies data.
	StartReplica()
	// replicaLoop starts replica loop
	replicaLoop()
	// replica tries to consume message
	replica(nodeID models.NodeID, replicator Replicator)
	// recovery rebuilds replication relation based on local partition.
	recovery(leader models.NodeID) error
}

// partition implements Partition interface.
type partition struct {
	ctx           context.Context
	cancel        context.CancelFunc
	currentNodeID models.NodeID
	db            string
	log           queue.FanOutQueue
	shardID       models.ShardID
	shard         tsdb.Shard
	family        tsdb.DataFamily
	running       *atomic.Bool

	replicators map[models.NodeID]Replicator
	cliFct      rpc.ClientStreamFactory
	stateMgr    storage.StateManager

	mutex sync.Mutex

	statistics           *metrics.StorageWriteAheadLogStatistics
	replicatorStatistics map[models.NodeID]*metrics.StorageReplicatorRunnerStatistics

	logger logger.Logger
}

// NewPartition creates a writeTask ahead log partition(db+shard+family time+leader).
func NewPartition(
	ctx context.Context,
	shard tsdb.Shard,
	family tsdb.DataFamily,
	currentNodeID models.NodeID,
	log queue.FanOutQueue,
	cliFct rpc.ClientStreamFactory,
	stateMgr storage.StateManager,
) Partition {
	c, cancel := context.WithCancel(ctx)
	p := &partition{
		ctx:                  c,
		cancel:               cancel,
		log:                  log,
		db:                   shard.Database().Name(),
		shardID:              shard.ShardID(),
		shard:                shard,
		family:               family,
		running:              &atomic.Bool{},
		currentNodeID:        currentNodeID,
		cliFct:               cliFct,
		stateMgr:             stateMgr,
		replicators:          make(map[models.NodeID]Replicator),
		statistics:           metrics.NewStorageWriteAheadLogStatistics(shard.Database().Name(), shard.ShardID().String()),
		replicatorStatistics: make(map[models.NodeID]*metrics.StorageReplicatorRunnerStatistics),
		logger:               logger.GetLogger("Replica", "Partition"),
	}
	return p
}

// ReplicaLog writes msg that leader sends replica msg.
// return appended index, if success.
func (p *partition) ReplicaLog(replicaIdx int64, msg []byte) (int64, error) {
	if p.closed.Load() {
		return 0, constants.ErrPartitionClosed
	}
	appendIdx := p.log.Queue().AppendedSeq() + 1
	if replicaIdx != appendIdx {
		return appendIdx, nil
	}
	p.statistics.ReceiveReplicaSize.Add(float64(len(msg)))
	if err := p.log.Queue().Put(msg); err != nil {
		p.statistics.ReplicaWALFailures.Incr()
		return -1, err
	}
	p.statistics.ReplicaWAL.Incr()
	return appendIdx, nil
}

// ReplicaAckIndex returns the index which replica appended index.
func (p *partition) ReplicaAckIndex() int64 {
	return p.log.Queue().AppendedSeq()
}

// ResetReplicaIndex resets replica index.
func (p *partition) ResetReplicaIndex(idx int64) {
	p.log.SetAppendedSeq(idx - 1)
}

// Path returns the path of partition.
func (p *partition) Path() string {
	return p.log.Path()
}

// IsExpire returns partition if it is expired.
func (p *partition) IsExpire() bool {
	p.log.Sync()       // sync acknowledged sequence of each ConsumerGroup
	p.log.Queue().GC() // try gc old data in queue

	opt := p.shard.Database().GetOption()
	ahead, _ := opt.GetAcceptWritableRange()
	timeRange := p.family.TimeRange()
	now := timeutil.Now()
	// add 15 minute buffer
	if timeRange.End+ahead+15*timeutil.OneMinute > now {
		return false
	}
	// partition is expired, check if all write ahead logs have been replicated
	hasData := false
	ns := p.log.ConsumerGroupNames()
	for _, name := range ns {
		consumerGroup, _ := p.log.GetOrCreateConsumerGroup(name)
		if !consumerGroup.IsEmpty() {
			hasData = true
			continue
		}
		// no data consume, can stop this replicator
		p.stopReplicator(name)
	}
	// no data means all data can be deleted
	return !hasData
}

// stopReplicator stops the replicator when no data can consume.
func (p *partition) stopReplicator(node string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.log.StopConsumerGroup(node)

	nodeID := models.ParseNodeID(node)
	// shutdown replicator if exist
	replicator, ok := p.replicators[nodeID]
	if ok {
		replicator.Close()
		// copy on write
		replicators := make(map[models.NodeID]Replicator, len(p.replicators)-1)
		for id := range p.replicators {
			if id != nodeID {
				replicators[id] = p.replicators[id]
			}
		}
		p.replicators = replicators
	}
}

// WriteLog writes msg that leader sends replica msg.
func (p *partition) WriteLog(msg []byte) error {
	if p.closed.Load() {
		return constants.ErrPartitionClosed
	}
	if len(msg) == 0 {
		return nil
	}
	p.statistics.ReceiveWriteSize.Add(float64(len(msg)))
	if err := p.log.Queue().Put(msg); err != nil {
		p.statistics.WriteWALFailures.Incr()
		return err
	}
	p.statistics.WriteWAL.Incr()
	return nil
}

// BuildReplicaForLeader builds replica relation when handle writeTask connection.
// local replicator: replica node == current node.
// remote replicator: replica node != current node.
func (p *partition) BuildReplicaForLeader(
	leader models.NodeID, replicas []models.NodeID,
) error {
	if leader != p.currentNodeID {
		return fmt.Errorf("leader not equals current node")
	}

	for _, replicaNodeID := range replicas {
		if err := p.buildReplica(leader, replicaNodeID); err != nil {
			p.logger.Error(
				"leader failed building replication channel to follower",
				logger.String("leader", leader.String()),
				logger.String("follower", replicaNodeID.String()),
				logger.Error(err),
			)
			return err
		}
	}
	return nil
}

// BuildReplicaForFollower builds replica relation when handle replica connection.
func (p *partition) BuildReplicaForFollower(leader, replica models.NodeID) error {
	if replica != p.currentNodeID {
		return fmt.Errorf("replica not equals current node")
	}
	err := p.buildReplica(leader, replica)
	if err != nil {
		p.logger.Error("follower failed building replication channel from leader",
			logger.Int("leader", leader.Int()),
			logger.Int("follower", replica.Int()),
		)
	}
	return err
}

// StartReplica iterates over all replicators and copies data.
func (p *partition) StartReplica() {
	if p.running.CompareAndSwap(false, true) {
		go p.replicaLoop()
	}
}

// replicaLoop starts replica loop
func (p *partition) replicaLoop() {
	for p.running.Load() {
		for nodeID, replicator := range p.replicators {
			p.replica(nodeID, replicator)
		}
	}
}

// replica tries to consume message
func (p *partition) replica(nodeID models.NodeID, replicator Replicator) {
	var replicatorStatistics = p.replicatorStatistics[nodeID]

	defer func() {
		if recovered := recover(); recovered != nil {
			replicatorStatistics.ReplicaPanics.Incr()
			p.logger.Error("panic when replica data",
				logger.Any("err", recovered),
				logger.Stack(),
			)
		}
	}()

	if replicator.IsReady() && replicator.Connect() {
		seq := replicator.Consume()
		if seq >= 0 {
			var replicatorType string
			switch replicator.(type) {
			case *localReplicator:
				replicatorType = replicatorTypeLocal
			case *remoteReplicator:
				replicatorType = replicatorTypeRemote
			}
			p.logger.Debug("replica write ahead log",
				logger.String("type", replicatorType),
				logger.String("replicator", replicator.String()),
				logger.Int64("index", seq))
			data, err := replicator.GetMessage(seq)
			if err != nil {
				replicator.IgnoreMessage(seq)
				replicatorStatistics.ConsumeMessageFailures.Incr()
				p.logger.Warn("cannot get replica message data, ignore replica message",
					logger.String("replicator", replicator.String()),
					logger.Int64("index", seq), logger.Error(err))
			} else {
				replicatorStatistics.ConsumeMessage.Incr()
				replicator.Replica(seq, data)
				replicatorStatistics.ReplicaBytes.Add(float64(len(data)))
			}
		}
	} else {
		p.logger.Warn("replica is not ready", logger.String("replicator", replicator.String()))
	}
}

// Close shutdowns all replica workers.
func (p *partition) Close() error {
	if p.closed.CAS(false, true) {
		// close log
		p.log.Close()
	}
	return nil
}

// Stop stops replicator channel.
func (p *partition) Stop() {
	p.running.Store(false)
	p.stop()
}

// stop stops replicator channel.
func (p *partition) stop() {
	// 1. cancel context of partition(will stop replicator)
	p.cancel()

	// 2. stop the peer of replicator
	var waiter sync.WaitGroup
	waiter.Add(len(p.replicators))
	for k := range p.replicators {
		r := p.replicators[k]
		go func() {
			r.Close()
			waiter.Done()
		}()
	}
	waiter.Wait()
}

// getReplicaState returns each family's log replica state.
func (p *partition) getReplicaState() models.FamilyLogReplicaState {
	replicators := p.log.ConsumerGroupNames()
	var stateOfReplicators []models.ReplicaPeerState
	for _, name := range replicators {
		fanout, err := p.log.GetOrCreateConsumerGroup(name)
		if err != nil {
			p.logger.Error("get fan out error when get replica state, ignore it")
			continue
		}
		peerState := models.ReplicaPeerState{
			Replicator: name,
			Consume:    fanout.ConsumedSeq(),
			ACK:        fanout.AcknowledgedSeq(),
			Pending:    fanout.Pending(),
		}
		nodeID := models.ParseNodeID(name)
		if replicator, ok := p.replicators[nodeID]; ok {
			var replicatorType string
			switch replicator.(type) {
			case *localReplicator:
				replicatorType = replicatorTypeLocal
			case *remoteReplicator:
				replicatorType = replicatorTypeRemote
			}
			peerState.ReplicatorType = replicatorType
			peerState.State = replicator.State().state
			peerState.StateErrMsg = replicator.State().errMsg
		}

		stateOfReplicators = append(stateOfReplicators, peerState)
	}
	return models.FamilyLogReplicaState{
		ShardID:     p.shardID,
		FamilyTime:  timeutil.FormatTimestamp(p.family.FamilyTime(), timeutil.DataTimeFormat2),
		Append:      p.log.Queue().AppendedSeq(),
		Replicators: stateOfReplicators,
	}
}

// buildReplica builds replica replication based on leader/follower node.
func (p *partition) buildReplica(leader, replica models.NodeID) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.replicators[replica]; ok {
		// exist
		return nil
	}
	walConsumer, err := p.log.GetOrCreateConsumerGroup(fmt.Sprintf("%d", replica))
	if err != nil {
		return err
	}
	var (
		replicator  Replicator
		replicaType string
	)
	channel := ReplicatorChannel{
		State: &models.ReplicaState{
			Database:   p.shard.Database().Name(),
			ShardID:    p.shardID,
			Leader:     leader,
			Follower:   replica,
			FamilyTime: p.family.TimeRange().Start,
		},
		ConsumerGroup: walConsumer,
	}
	if replica == p.currentNodeID {
		// local replicator
		replicator = newLocalReplicatorFn(&channel, p.shard, p.family)
		replicaType = replicatorTypeLocal
	} else {
		// build remote replicator
		replicator = newRemoteReplicatorFn(p.ctx, &channel, p.stateMgr, p.cliFct)
		replicaType = replicatorTypeRemote
	}

	var (
		state                = replicator.ReplicaState()
		replicators          = make(map[models.NodeID]Replicator, len(p.replicators))
		replicatorStatistics = make(map[models.NodeID]*metrics.StorageReplicatorRunnerStatistics, len(p.replicatorStatistics))
	)

	for nodeID, replicator0 := range p.replicators {
		replicators[nodeID] = replicator0
	}
	for nodeID, statistics := range p.replicatorStatistics {
		replicatorStatistics[nodeID] = statistics
	}

	// copy on write
	replicatorStatistics[replica] = metrics.NewStorageReplicatorRunnerStatistics(replicaType, state.Database, state.ShardID.String())
	//	the order is to first use replicators and then replicatorStatistics,
	//	so replicatorStatistics must be assigned first in concurrent scenarios,
	//	perhaps in the future, this should be refactored to encapsulate two variables into a single structure.
	p.replicatorStatistics = replicatorStatistics
	replicators[replica] = replicator
	p.replicators = replicators

	return nil
}

// recovery rebuilds replication relation based on local partition.
func (p *partition) recovery(leader models.NodeID) error {
	replicatorNames := p.log.ConsumerGroupNames()
	for _, replica := range replicatorNames {
		if err := p.buildReplica(leader, models.ParseNodeID(replica)); err != nil {
			return err
		}
	}
	return nil
}
