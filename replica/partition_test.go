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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
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
	r := NewMockReplicator(ctrl)
	r.EXPECT().String().Return("test").AnyTimes()
	newLocalReplicatorFn = func(_ tsdb.Shard) Replicator {
		return r
	}
	newRemoteReplicatorFn = func(_ *ReplicatorChannel,
		_ rpc.ClientStreamFactory) Replicator {
		return r
	}

	p := NewPartition(1, nil, 1, nil, nil)
	err := p.BuildReplicaForLeader(2, []models.NodeID{1, 2, 3})
	assert.Error(t, err)

	r.EXPECT().IsReady().Return(false).AnyTimes()
	err = p.BuildReplicaForLeader(1, []models.NodeID{1, 2, 3})
	assert.NoError(t, err)
	// ignore re-build
	err = p.BuildReplicaForLeader(1, []models.NodeID{1, 2, 3})
	assert.NoError(t, err)

	p1 := p.(*partition)
	assert.Len(t, p1.peers, 3)
}

func TestPartition_BuildReplicaForFollower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newLocalReplicatorFn = NewLocalReplicator
		newRemoteReplicatorFn = NewRemoteReplicator
		ctrl.Finish()
	}()
	r := NewMockReplicator(ctrl)
	r.EXPECT().String().Return("test").AnyTimes()
	newLocalReplicatorFn = func(_ tsdb.Shard) Replicator {
		return r
	}
	newRemoteReplicatorFn = func(_ *ReplicatorChannel,
		_ rpc.ClientStreamFactory) Replicator {
		return r
	}

	p := NewPartition(1, nil, 1, nil, nil)
	err := p.BuildReplicaForFollower(2, 2)
	assert.Error(t, err)

	r.EXPECT().IsReady().Return(false).AnyTimes()
	err = p.BuildReplicaForFollower(2, 1)
	assert.NoError(t, err)
}

func TestPartition_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newLocalReplicatorFn = NewLocalReplicator
		newRemoteReplicatorFn = NewRemoteReplicator
		ctrl.Finish()
	}()
	r := NewMockReplicator(ctrl)
	l := queue.NewMockFanOutQueue(ctrl)
	r.EXPECT().String().Return("test").AnyTimes()
	newLocalReplicatorFn = func(_ tsdb.Shard) Replicator {
		return r
	}
	newRemoteReplicatorFn = func(_ *ReplicatorChannel,
		_ rpc.ClientStreamFactory) Replicator {
		return r
	}

	l.EXPECT().Close().MaxTimes(2)
	p := NewPartition(1, nil, 1, l, nil)
	err := p.Close()
	assert.NoError(t, err)
	r.EXPECT().IsReady().Return(false).AnyTimes()
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
	p := NewPartition(1, nil, 1, l, nil)
	l.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("err"))
	err := p.WriteLog([]byte{1})
	assert.Error(t, err)
	// msg is empty
	err = p.WriteLog(nil)
	assert.NoError(t, err)
}

func TestPartition_ReplicaLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	l := queue.NewMockFanOutQueue(ctrl)
	p := NewPartition(1, nil, 1, l, nil)
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
