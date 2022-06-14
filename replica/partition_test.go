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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
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
	family.EXPECT().FamilyTime().Return(timeutil.Now()).AnyTimes()
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
		_ storage.StateManager, _ rpc.ClientStreamFactory) Replicator {
		return r
	}

	log := queue.NewMockFanOutQueue(ctrl)
	log.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, nil).MaxTimes(3)
	family.EXPECT().TimeRange().Return(timeutil.TimeRange{}).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, log, nil, nil)
	err := p.BuildReplicaForLeader(2, []models.NodeID{1, 2, 3})
	assert.Error(t, err)

	r.EXPECT().IsReady().Return(true).AnyTimes()
	r.EXPECT().Connect().Return(true).AnyTimes()
	r.EXPECT().Consume().Return(int64(-1)).AnyTimes()
	err = p.BuildReplicaForLeader(1, []models.NodeID{1, 2, 3})
	assert.NoError(t, err)
	// ignore re-build
	err = p.BuildReplicaForLeader(1, []models.NodeID{1, 2, 3})
	assert.NoError(t, err)

	p1 := p.(*partition)
	assert.Len(t, p1.peers, 3)

	log.EXPECT().HeadSeq().Return(int64(10))
	assert.Equal(t, int64(9), p.ReplicaAckIndex())
	log.EXPECT().SetAppendSeq(int64(100))
	p.ResetReplicaIndex(100)
	log.EXPECT().Path().Return("path")
	assert.Equal(t, "path", p.Path())

	// create fanout failure
	p = NewPartition(context.TODO(), shard, family, 1, log, nil, nil)
	log.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, fmt.Errorf("err"))
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
	family.EXPECT().FamilyTime().Return(timeutil.Now()).AnyTimes()
	r.EXPECT().String().Return("TestPartition_BuildReplicaForFollower").AnyTimes()
	r.EXPECT().ReplicaState().Return(&models.ReplicaState{}).AnyTimes()
	r.EXPECT().Pending().Return(int64(10)).AnyTimes()
	newLocalReplicatorFn = func(_ *ReplicatorChannel, _ tsdb.Shard, _ tsdb.DataFamily) Replicator {
		return r
	}
	newRemoteReplicatorFn = func(_ context.Context, _ *ReplicatorChannel,
		_ storage.StateManager, _ rpc.ClientStreamFactory) Replicator {
		return r
	}

	log := queue.NewMockFanOutQueue(ctrl)
	log.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, nil)
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
	log.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, fmt.Errorf("err"))
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
	family.EXPECT().FamilyTime().Return(timeutil.Now()).AnyTimes()
	l := queue.NewMockFanOutQueue(ctrl)
	l.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, nil).AnyTimes()
	r.EXPECT().String().Return("TestPartition_Close").AnyTimes()
	r.EXPECT().Pending().Return(int64(10)).AnyTimes()
	newLocalReplicatorFn = func(_ *ReplicatorChannel, _ tsdb.Shard, _ tsdb.DataFamily) Replicator {
		return r
	}
	newRemoteReplicatorFn = func(_ context.Context, _ *ReplicatorChannel,
		_ storage.StateManager, _ rpc.ClientStreamFactory) Replicator {
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
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(timeutil.Now()).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, l, nil, nil)
	l.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("err"))
	err := p.WriteLog([]byte{1})
	assert.Error(t, err)
	// msg is empty
	err = p.WriteLog(nil)
	assert.NoError(t, err)
	l.EXPECT().Put(gomock.Any()).Return(nil)
	err = p.WriteLog([]byte{1})
	assert.NoError(t, err)
}

func TestPartition_ReplicaLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	l := queue.NewMockFanOutQueue(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(timeutil.Now()).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, l, nil, nil)
	// case 1: replica idx err
	l.EXPECT().HeadSeq().Return(int64(8))
	idx, err := p.ReplicaLog(10, []byte{1})
	assert.NoError(t, err)
	assert.Equal(t, idx, int64(8))

	// case 2: put err
	l.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("err"))
	l.EXPECT().HeadSeq().Return(int64(10))
	idx, err = p.ReplicaLog(10, []byte{1})
	assert.Error(t, err)
	assert.Equal(t, idx, int64(-1))

	// case 3: put ok
	l.EXPECT().Put(gomock.Any()).Return(nil)
	l.EXPECT().HeadSeq().Return(int64(10))
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
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(timeutil.Now()).AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	p := NewPartition(context.TODO(), shard, family, 1, l, nil, nil)
	p1 := p.(*partition)
	peer := NewMockReplicatorPeer(ctrl)
	peer.EXPECT().ReplicatorState().Return("remote", &state{state: models.ReplicatorReadyState}).AnyTimes()
	p1.mutex.Lock()
	p1.peers[models.NodeID(1)] = peer
	p1.peers[models.NodeID(2)] = peer
	p1.mutex.Unlock()
	l.EXPECT().FanOutNames().Return([]string{"1", "2"})
	fan := queue.NewMockFanOut(ctrl)
	l.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, fmt.Errorf("err"))
	l.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(fan, nil)
	fan.EXPECT().HeadSeq().Return(int64(1))
	fan.EXPECT().TailSeq().Return(int64(1))
	fan.EXPECT().Pending().Return(int64(1))
	l.EXPECT().HeadSeq().Return(int64(1))
	state := p.getReplicaState()
	assert.NotNil(t, state)
}

func TestPartition_IsExpire(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	db.EXPECT().GetOption().Return(&option.DatabaseOption{Ahead: "1h"}).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)

	log := queue.NewMockFanOutQueue(ctrl)
	log.EXPECT().FanOutNames().Return([]string{"test"}).AnyTimes()
	p := &partition{
		shard:  shard,
		family: family,
		log:    log,
	}
	q := queue.NewMockFanOut(ctrl)
	log.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(q, nil).AnyTimes()
	q.EXPECT().IsEmpty().Return(false)
	assert.False(t, p.IsExpire())

	q.EXPECT().IsEmpty().Return(true).AnyTimes()

	family.EXPECT().TimeRange().Return(timeutil.TimeRange{End: timeutil.Now()})
	assert.False(t, p.IsExpire())

	family.EXPECT().TimeRange().Return(timeutil.TimeRange{End: timeutil.Now() - timeutil.OneHour - 16*timeutil.OneMinute})
	assert.True(t, p.IsExpire())
}

func TestPartition_recovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newLocalReplicatorFn = NewLocalReplicator
		newReplicatorPeerFn = NewReplicatorPeer
		ctrl.Finish()
	}()

	shard := tsdb.NewMockShard(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	db.EXPECT().Name().Return("test").AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)

	log := queue.NewMockFanOutQueue(ctrl)
	log.EXPECT().FanOutNames().Return([]string{"1"}).AnyTimes()
	p := &partition{
		shard:         shard,
		family:        family,
		currentNodeID: 1,
		peers:         make(map[models.NodeID]ReplicatorPeer),
		log:           log,
	}

	t.Run("recovery failure", func(t *testing.T) {
		log.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, fmt.Errorf("err"))
		err := p.recovery(1)
		assert.Error(t, err)
	})
	t.Run("recovery successfully", func(t *testing.T) {
		q := queue.NewMockFanOut(ctrl)
		log.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(q, nil)
		family.EXPECT().TimeRange().Return(timeutil.TimeRange{Start: timeutil.Now()})
		newLocalReplicatorFn = func(channel *ReplicatorChannel, shard tsdb.Shard, family tsdb.DataFamily) Replicator {
			return nil
		}
		peer := NewMockReplicatorPeer(ctrl)
		peer.EXPECT().Startup()
		newReplicatorPeerFn = func(replicator Replicator) ReplicatorPeer {
			return peer
		}
		err := p.recovery(1)
		assert.NoError(t, err)
	})
}

func TestPartition_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peer1 := NewMockReplicatorPeer(ctrl)
	peer2 := NewMockReplicatorPeer(ctrl)
	ctx, cancel := context.WithCancel(context.TODO())
	p := &partition{
		ctx:    ctx,
		cancel: cancel,
		peers: map[models.NodeID]ReplicatorPeer{
			1: peer1,
			2: peer2,
		},
	}
	peer1.EXPECT().Shutdown()
	peer2.EXPECT().Shutdown()
	p.Stop()
}
