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
	"sync"
	"testing"

	"github.com/lindb/common/pkg/logger"
	commontimeutil "github.com/lindb/common/pkg/timeutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

func TestPartition_BuildReplicaRelation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newLocalReplicatorFn = NewLocalReplicator
		newRemoteReplicatorFn = NewRemoteReplicator
		ctrl.Finish()
	}()
	database := tsdb.NewMockDatabase(ctrl)
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()
	database.EXPECT().Name().Return("test").AnyTimes()
	r := NewMockReplicator(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	r.EXPECT().String().Return("TestPartition_BuildReplicaRelation").AnyTimes()
	r.EXPECT().ReplicaState().Return(&models.ReplicaState{}).AnyTimes()
	r.EXPECT().Pending().Return(int64(10)).AnyTimes()
	newLocalReplicatorFn = func(_ *ReplicatorChannel, _ tsdb.Shard, _ tsdb.DataFamily) Replicator {
		return r
	}
	newRemoteReplicatorFn = func(_ context.Context, _ *ReplicatorChannel,
		_ storage.StateManager, _ rpc.ClientStreamFactory,
	) Replicator {
		return r
	}

	log := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	log.EXPECT().Queue().Return(q).AnyTimes()
	log.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(nil, nil).MaxTimes(3)
	family.EXPECT().TimeRange().Return(timeutil.TimeRange{}).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, log, nil, nil)
	p1 := p.(*partition)
	err := p.BuildReplicaForLeader(2, []models.NodeID{1, 2, 3})
	assert.Error(t, err)
	assert.Len(t, p1.replicators, 0)

	r.EXPECT().IsReady().Return(true).AnyTimes()
	r.EXPECT().Connect().Return(true).AnyTimes()
	r.EXPECT().Consume().Return(int64(-1)).AnyTimes()
	err = p.BuildReplicaForLeader(1, []models.NodeID{1, 2, 3})
	assert.NoError(t, err)
	assert.Len(t, p1.replicators, 3)

	// ignore re-build
	err = p.BuildReplicaForLeader(1, []models.NodeID{1, 2, 3})
	assert.NoError(t, err)
	assert.Len(t, p1.replicators, 3)

	q.EXPECT().AppendedSeq().Return(int64(10))
	assert.Equal(t, int64(10), p.ReplicaAckIndex())
	log.EXPECT().SetAppendedSeq(int64(99))
	p.ResetReplicaIndex(100)
	log.EXPECT().Path().Return("path")
	assert.Equal(t, "path", p.Path())

	// create consume group failure
	p = NewPartition(context.TODO(), shard, family, 1, log, nil, nil)
	log.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = p.BuildReplicaForLeader(1, []models.NodeID{1, 2, 3})
	assert.Error(t, err)
}

func TestPartition_BuildReplicaForFollower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newLocalReplicatorFn = NewLocalReplicator
		newRemoteReplicatorFn = NewRemoteReplicator
		ctrl.Finish()
	}()
	r := NewMockReplicator(ctrl)
	database := tsdb.NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	shard.EXPECT().Database().Return(database).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()
	r.EXPECT().String().Return("TestPartition_BuildReplicaForFollower").AnyTimes()
	r.EXPECT().ReplicaState().Return(&models.ReplicaState{}).AnyTimes()
	r.EXPECT().Pending().Return(int64(10)).AnyTimes()
	newLocalReplicatorFn = func(_ *ReplicatorChannel, _ tsdb.Shard, _ tsdb.DataFamily) Replicator {
		return r
	}
	newRemoteReplicatorFn = func(_ context.Context, _ *ReplicatorChannel,
		_ storage.StateManager, _ rpc.ClientStreamFactory,
	) Replicator {
		return r
	}

	log := queue.NewMockFanOutQueue(ctrl)
	log.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(nil, nil)
	family.EXPECT().TimeRange().Return(timeutil.TimeRange{}).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, log, nil, nil)
	err := p.BuildReplicaForFollower(2, 2)
	assert.Error(t, err)

	r.EXPECT().IsReady().Return(true).AnyTimes()
	r.EXPECT().Connect().Return(true).AnyTimes()
	r.EXPECT().Consume().Return(int64(-1)).AnyTimes()
	err = p.BuildReplicaForFollower(2, 1)
	assert.NoError(t, err)

	// create fan ot failure
	log.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(nil, fmt.Errorf("err"))
	p = NewPartition(context.TODO(), shard, family, 1, log, nil, nil)
	err = p.BuildReplicaForFollower(2, 1)
	assert.Error(t, err)
}

func TestPartition_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newLocalReplicatorFn = NewLocalReplicator
		newRemoteReplicatorFn = NewRemoteReplicator
		ctrl.Finish()
	}()
	r := NewMockReplicator(ctrl)
	r.EXPECT().ReplicaState().Return(&models.ReplicaState{}).AnyTimes()
	database := tsdb.NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	shard.EXPECT().Database().Return(database).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()
	l := queue.NewMockFanOutQueue(ctrl)
	l.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(nil, nil).AnyTimes()
	r.EXPECT().String().Return("TestPartition_Close").AnyTimes()
	r.EXPECT().Pending().Return(int64(10)).AnyTimes()
	newLocalReplicatorFn = func(_ *ReplicatorChannel, _ tsdb.Shard, _ tsdb.DataFamily) Replicator {
		return r
	}
	newRemoteReplicatorFn = func(_ context.Context, _ *ReplicatorChannel,
		_ storage.StateManager, _ rpc.ClientStreamFactory,
	) Replicator {
		return r
	}

	l.EXPECT().Close().MaxTimes(2)
	family.EXPECT().TimeRange().Return(timeutil.TimeRange{}).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, l, nil, nil)
	err := p.Close()
	assert.NoError(t, err)
	r.EXPECT().IsReady().Return(true).AnyTimes()
	r.EXPECT().Connect().Return(true).AnyTimes()
	r.EXPECT().Consume().Return(int64(-1)).AnyTimes()
	err = p.BuildReplicaForLeader(1, []models.NodeID{1, 2, 3})
	assert.NoError(t, err)
	err = p.Close()
	assert.NoError(t, err)
}

func TestPartition_WriteLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	l := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	l.EXPECT().Queue().Return(q).AnyTimes()
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, l, nil, nil)
	q.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("err"))
	err := p.WriteLog([]byte{1})
	assert.Error(t, err)
	// msg is empty
	err = p.WriteLog(nil)
	assert.NoError(t, err)
	q.EXPECT().Put(gomock.Any()).Return(nil)
	err = p.WriteLog([]byte{1})
	assert.NoError(t, err)
}

func TestPartition_ReplicaLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	l := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	l.EXPECT().Queue().Return(q).AnyTimes()
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, l, nil, nil)
	// case 1: replica idx err
	q.EXPECT().AppendedSeq().Return(int64(8))
	idx, err := p.ReplicaLog(10, []byte{1})
	assert.NoError(t, err)
	assert.Equal(t, idx, int64(9))

	// case 2: put err
	q.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("err"))
	q.EXPECT().AppendedSeq().Return(int64(9))
	idx, err = p.ReplicaLog(10, []byte{1})
	assert.Error(t, err)
	assert.Equal(t, idx, int64(-1))

	// case 3: put ok
	q.EXPECT().Put(gomock.Any()).Return(nil)
	q.EXPECT().AppendedSeq().Return(int64(9))
	idx, err = p.ReplicaLog(10, []byte{1})
	assert.NoError(t, err)
	assert.Equal(t, idx, int64(10))
}

func TestPartition_getReplicaState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	l := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	l.EXPECT().Queue().Return(q).AnyTimes()
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, l, nil, nil)
	p1 := p.(*partition)

	var wait sync.WaitGroup

	r := NewMockReplicator(ctrl)
	r.EXPECT().String().Return("str").AnyTimes()
	r.EXPECT().IsReady().Return(true).Times(2)
	r.EXPECT().Connect().Return(true).Times(2)
	r.EXPECT().Consume().Return(int64(-1)).Times(2)
	r.EXPECT().Close().AnyTimes()
	r.EXPECT().State().Return(&state{state: models.ReplicatorReadyState}).AnyTimes()
	r.EXPECT().IsReady().DoAndReturn(func() bool {
		wait.Done()
		p.Stop()
		return false
	}).AnyTimes()

	p1.mutex.Lock()
	p1.replicators[models.NodeID(1)] = r
	p1.replicators[models.NodeID(2)] = r
	p1.mutex.Unlock()

	l.EXPECT().ConsumerGroupNames().Return([]string{"1", "2"})
	fan := queue.NewMockConsumerGroup(ctrl)
	l.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(nil, fmt.Errorf("err"))
	l.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(fan, nil)
	fan.EXPECT().ConsumedSeq().Return(int64(1))
	fan.EXPECT().AcknowledgedSeq().Return(int64(1))
	fan.EXPECT().Pending().Return(int64(1))
	q.EXPECT().AppendedSeq().Return(int64(1))
	state := p.getReplicaState()
	assert.NotNil(t, state)

	wait.Add(2)
	p.StartReplica()
	wait.Wait()
}

func TestPartition_IsExpire(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)

	log := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	q.EXPECT().GC().AnyTimes()
	log.EXPECT().Sync().AnyTimes()
	log.EXPECT().Queue().Return(q).AnyTimes()
	log.EXPECT().ConsumerGroupNames().Return([]string{"1"}).AnyTimes()
	r := NewMockReplicator(ctrl)
	p := &partition{
		shard:  shard,
		family: family,
		log:    log,
		replicators: map[models.NodeID]Replicator{
			models.NodeID(1): r,
		},
	}
	cg := queue.NewMockConsumerGroup(ctrl)
	log.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(cg, nil).AnyTimes()

	t.Run("partition not expire", func(t *testing.T) {
		db.EXPECT().GetOption().Return(&option.DatabaseOption{Ahead: "1h"})
		family.EXPECT().TimeRange().Return(timeutil.TimeRange{End: commontimeutil.Now()})
		assert.False(t, p.IsExpire())
	})
	t.Run("partition not expire, when ahead=0", func(t *testing.T) {
		db.EXPECT().GetOption().Return(&option.DatabaseOption{})
		family.EXPECT().TimeRange().Return(timeutil.TimeRange{End: commontimeutil.Now()})
		assert.False(t, p.IsExpire())
	})

	t.Run("partition is expire, but has data need replica", func(t *testing.T) {
		db.EXPECT().GetOption().Return(&option.DatabaseOption{Ahead: "1h"})
		family.EXPECT().TimeRange().Return(timeutil.TimeRange{End: commontimeutil.Now() - commontimeutil.OneHour - 16*commontimeutil.OneMinute})
		cg.EXPECT().IsEmpty().Return(false)
		assert.False(t, p.IsExpire())
	})
	t.Run("partition is expire, no data replica, can stop replicator", func(t *testing.T) {
		db.EXPECT().GetOption().Return(&option.DatabaseOption{Ahead: "1h"})
		p.mutex.Lock()
		assert.Len(t, p.replicators, 1)
		p.mutex.Unlock()
		family.EXPECT().TimeRange().Return(timeutil.TimeRange{End: commontimeutil.Now() - commontimeutil.OneHour - 16*commontimeutil.OneMinute})
		cg.EXPECT().IsEmpty().Return(true)
		log.EXPECT().StopConsumerGroup(gomock.Any())
		r.EXPECT().Close()
		assert.True(t, p.IsExpire())
		p.mutex.Lock()
		assert.Empty(t, p.replicators)
		p.mutex.Unlock()
	})
}

func TestPartition_recovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newLocalReplicatorFn = NewLocalReplicator
		ctrl.Finish()
	}()

	shard := tsdb.NewMockShard(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	db.EXPECT().Name().Return("test").AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	log := queue.NewMockFanOutQueue(ctrl)
	log.EXPECT().ConsumerGroupNames().Return([]string{"1"}).AnyTimes()

	p := &partition{
		shard:         shard,
		family:        family,
		currentNodeID: 1,
		replicators:   make(map[models.NodeID]Replicator),
		log:           log,
	}

	t.Run("recovery failure", func(t *testing.T) {
		log.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(nil, fmt.Errorf("err"))
		err := p.recovery(1)
		assert.Error(t, err)
	})
	t.Run("recovery successfully", func(t *testing.T) {
		q := queue.NewMockConsumerGroup(ctrl)
		log.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(q, nil)
		family.EXPECT().TimeRange().Return(timeutil.TimeRange{Start: commontimeutil.Now()})
		r := NewMockReplicator(ctrl)
		r.EXPECT().ReplicaState().Return(&models.ReplicaState{}).AnyTimes()
		newLocalReplicatorFn = func(channel *ReplicatorChannel, shard tsdb.Shard, family tsdb.DataFamily) Replicator {
			return r
		}
		err := p.recovery(1)
		assert.NoError(t, err)
	})
}

func TestReplicatorPeer_replica_panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := &partition{
		replicatorStatistics: map[models.NodeID]*metrics.StorageReplicatorRunnerStatistics{
			1: metrics.NewStorageReplicatorRunnerStatistics("local", "test", "1"),
		},
		logger: logger.GetLogger("Replica", "Partition"),
	}
	replicator := NewMockReplicator(ctrl)
	replicator.EXPECT().IsReady().DoAndReturn(func() bool {
		panic("err")
	})
	p.replica(1, replicator)
}

func TestNewReplicator_replica3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	replicator := NewMockReplicator(ctrl)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	database := tsdb.NewMockDatabase(ctrl)
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	log := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	log.EXPECT().Queue().Return(q).AnyTimes()
	log.EXPECT().GetOrCreateConsumerGroup(gomock.Any()).Return(nil, nil).MaxTimes(3)
	p := NewPartition(ctx, shard, family, 1, log, nil, nil)
	p.(*partition).replicatorStatistics = map[models.NodeID]*metrics.StorageReplicatorRunnerStatistics{
		1: metrics.NewStorageReplicatorRunnerStatistics("local", "test", "1"),
	}
	p.(*partition).replicators = map[models.NodeID]Replicator{
		1: replicator,
	}

	replicator.EXPECT().String().Return("str").AnyTimes()
	replicator.EXPECT().ReplicaState().Return(&models.ReplicaState{}).AnyTimes()
	replicator.EXPECT().Pending().Return(int64(19)).AnyTimes()
	replicator.EXPECT().Close().AnyTimes()
	replicator.EXPECT().IgnoreMessage(gomock.Any()).AnyTimes()
	replicator.EXPECT().Pause().AnyTimes()

	// loop 1: no data
	replicator.EXPECT().IsReady().Return(true)
	replicator.EXPECT().Connect().Return(true)
	replicator.EXPECT().Consume().Return(int64(-1)) // no data
	// loop 2: get message err
	replicator.EXPECT().IsReady().Return(true)
	replicator.EXPECT().Connect().Return(true)
	replicator.EXPECT().Consume().Return(int64(1))                          // has data
	replicator.EXPECT().GetMessage(int64(1)).Return(nil, fmt.Errorf("err")) // get message err
	// loop 3: do replica
	replicator.EXPECT().IsReady().Return(true)
	replicator.EXPECT().Connect().Return(true)
	replicator.EXPECT().Consume().Return(int64(1))            // has data
	replicator.EXPECT().GetMessage(int64(1)).Return(nil, nil) // get message
	replicator.EXPECT().Replica(gomock.Any(), gomock.Any())   // replica

	var wait sync.WaitGroup
	replicator.EXPECT().IsReady().DoAndReturn(func() bool {
		wait.Done()
		p.Stop()
		return false
	}).AnyTimes()

	wait.Add(1)
	p.StartReplica()
	wait.Wait()
}

func TestPartition_replicaLoop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	log := queue.NewMockFanOutQueue(ctrl)
	log.EXPECT().StopConsumerGroup(gomock.Any()).AnyTimes()
	r1 := NewMockReplicator(ctrl)
	r2 := NewMockReplicator(ctrl)
	p := &partition{
		ctx:    ctx,
		cancel: cancel,
		replicators: map[models.NodeID]Replicator{
			1: r1,
			2: r2,
		},
		running: atomic.NewBool(false),
		log:     log,
		logger:  logger.GetLogger("Replica", "Partition"),
	}
	var wait sync.WaitGroup
	for _, r := range []*MockReplicator{r1, r2} {
		r.EXPECT().String().Return("str").AnyTimes()
		r.EXPECT().IsReady().Return(true)
		r.EXPECT().Connect().Return(true)
		r.EXPECT().Consume().Return(int64(-1))
		r.EXPECT().IsReady().DoAndReturn(func() bool {
			wait.Done()
			p.running.Store(false)
			return false
		}).AnyTimes()
	}
	p.running.Store(true)
	wait.Add(len(p.replicators))
	go p.replicaLoop()
	wait.Wait()
}

func TestPartition_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	log := queue.NewMockFanOutQueue(ctrl)
	log.EXPECT().StopConsumerGroup(gomock.Any()).AnyTimes()
	r1 := NewMockReplicator(ctrl)
	r2 := NewMockReplicator(ctrl)
	p := &partition{
		ctx:    ctx,
		cancel: cancel,
		replicators: map[models.NodeID]Replicator{
			1: r1,
			2: r2,
		},
		running: &atomic.Bool{},
		log:     log,
	}

	r1.EXPECT().Close()
	r2.EXPECT().Close()

	p.stopReplicator("1")
	assert.Len(t, p.replicators, 1)
	_, ok := p.replicators[1]
	assert.False(t, ok)

	p.Stop()

	assert.False(t, p.running.Load())
}

func TestPartition_getReplicaState2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cliFct := rpc.NewMockClientStreamFactory(ctrl)
	stateMgr := storage.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchNodeStateChangeEvent(gomock.Any(), gomock.Any()).AnyTimes()
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, true).AnyTimes()
	cg := queue.NewMockConsumerGroup(ctrl)
	fq := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	fq.EXPECT().Queue().Return(q).AnyTimes()
	cg.EXPECT().Queue().Return(fq).AnyTimes()
	rc := &ReplicatorChannel{
		State: &models.ReplicaState{
			Database: "test",
			ShardID:  0,
			Leader:   1,
			Follower: 2,
		},
		ConsumerGroup: cg,
	}

	log, err := queue.NewFanOutQueue("test", 1024)
	assert.NoError(t, err)
	_, err = log.GetOrCreateConsumerGroup("1")
	assert.NoError(t, err)
	_, err = log.GetOrCreateConsumerGroup("2")
	assert.NoError(t, err)

	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()

	local := &localReplicator{}
	remote := NewRemoteReplicator(context.Background(), rc, stateMgr, cliFct)

	p := &partition{
		replicators: map[models.NodeID]Replicator{
			1: local,
			2: remote,
		},
		log:    log,
		family: family,
	}

	p.getReplicaState()
}

func TestPartition_remote_replica(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cliFct := rpc.NewMockClientStreamFactory(ctrl)
	replicaCli := protoReplicaV1.NewMockReplicaServiceClient(ctrl)
	cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil).AnyTimes()
	replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil).AnyTimes()
	replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
		AckIndex: 10,
	}, nil).AnyTimes()
	replicaCli.EXPECT().Reset(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	stateMgr := storage.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchNodeStateChangeEvent(gomock.Any(), gomock.Any()).AnyTimes()
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, true).AnyTimes()
	cg := queue.NewMockConsumerGroup(ctrl)
	fq := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	fq.EXPECT().Queue().Return(q).AnyTimes()
	cg.EXPECT().Queue().Return(fq).AnyTimes()
	cg.EXPECT().ConsumedSeq().Return(int64(10)).AnyTimes()
	cg.EXPECT().Consume().Return(int64(1))
	q.EXPECT().Get(gomock.Any()).Return([]byte{}, nil)
	rc := &ReplicatorChannel{
		State: &models.ReplicaState{
			Database: "test",
			ShardID:  0,
			Leader:   1,
			Follower: 2,
		},
		ConsumerGroup: cg,
	}
	log, err := queue.NewFanOutQueue("test", 1024)
	assert.NoError(t, err)
	_, err = log.GetOrCreateConsumerGroup("2")
	assert.NoError(t, err)
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()

	remote := NewRemoteReplicator(context.Background(), rc, stateMgr, cliFct)
	p := &partition{
		replicators: map[models.NodeID]Replicator{
			2: remote,
		},
		log:    log,
		family: family,
		replicatorStatistics: map[models.NodeID]*metrics.StorageReplicatorRunnerStatistics{
			2: metrics.NewStorageReplicatorRunnerStatistics("remote", "test", "1"),
		},
		logger: logger.GetLogger("Replica", "Partition"),
	}

	p.replica(2, remote)
}

func TestPartition_local_replica(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cliFct := rpc.NewMockClientStreamFactory(ctrl)
	replicaCli := protoReplicaV1.NewMockReplicaServiceClient(ctrl)
	cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil).AnyTimes()
	replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil).AnyTimes()
	replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
		AckIndex: 10,
	}, nil).AnyTimes()
	replicaCli.EXPECT().Reset(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	stateMgr := storage.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchNodeStateChangeEvent(gomock.Any(), gomock.Any()).AnyTimes()
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, true).AnyTimes()
	cg := queue.NewMockConsumerGroup(ctrl)
	fq := queue.NewMockFanOutQueue(ctrl)
	q := queue.NewMockQueue(ctrl)
	fq.EXPECT().Queue().Return(q).AnyTimes()
	cg.EXPECT().Queue().Return(fq).AnyTimes()
	cg.EXPECT().ConsumedSeq().Return(int64(10)).AnyTimes()
	cg.EXPECT().Consume().Return(int64(1)).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(commontimeutil.Now()).AnyTimes()
	family.EXPECT().TimeRange().Return(timeutil.TimeRange{}).AnyTimes()
	family.EXPECT().AckSequence(gomock.Any(), gomock.Any()).AnyTimes()
	family.EXPECT().Retain().AnyTimes()
	family.EXPECT().ValidateSequence(gomock.Any(), gomock.Any()).Return(true)
	family.EXPECT().CommitSequence(gomock.Any(), gomock.Any()).AnyTimes()
	fq.EXPECT().Queue().Return(q).AnyTimes()
	cg.EXPECT().Queue().Return(fq).AnyTimes()
	cg.EXPECT().ConsumedSeq().Return(int64(10)).AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	database := tsdb.NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("database name").AnyTimes()
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	log, err := queue.NewFanOutQueue("test", 1024)
	assert.NoError(t, err)

	p := NewPartition(context.TODO(), shard, family, 1, log, nil, stateMgr)
	err = p.BuildReplicaForFollower(1, 1)
	assert.NoError(t, err)
	partition := p.(*partition)
	err = partition.WriteLog([]byte("test msg"))
	assert.NoError(t, err)
	partition.replica(1, partition.replicators[1])
}

func TestPartition_WriteLog_After_Close(t *testing.T) {
	dirPath := t.TempDir()
	log, err := queue.NewFanOutQueue(dirPath, 1024*1014)
	ctx, cancel := context.WithCancel(context.TODO())
	assert.NoError(t, err)
	p := &partition{
		ctx:        ctx,
		cancel:     cancel,
		statistics: metrics.NewStorageWriteAheadLogStatistics("test", "0"),
		closed:     atomic.NewBool(false),
		log:        log,
	}
	err = p.WriteLog([]byte("test"))
	assert.NoError(t, err)
	p.Close()
	err = p.WriteLog([]byte("test"))
	assert.Equal(t, constants.ErrPartitionClosed, err)
}
