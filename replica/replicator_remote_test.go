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
	"time"

	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRemoteReplicator_IsReady(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	cliFct := rpc.NewMockClientStreamFactory(ctrl)
	stateMgr := storage.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchNodeStateChangeEvent(gomock.Any(), gomock.Any()).AnyTimes()
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, true).AnyTimes()
	replicaCli := protoReplicaV1.NewMockReplicaServiceClient(ctrl)
	q := queue.NewMockFanOut(ctrl)
	fq := queue.NewMockFanOutQueue(ctrl)
	q.EXPECT().Queue().Return(fq).AnyTimes()
	rc := &ReplicatorChannel{
		State: &models.ReplicaState{
			Database: "test",
			ShardID:  0,
			Leader:   1,
			Follower: 2,
		},
		Queue: q,
	}

	cases := []struct {
		name    string
		prepare func(r *remoteReplicator)
		ready   bool
	}{
		{
			name: "replicator is ready",
			prepare: func(r *remoteReplicator) {
				r.state = ReplicatorReadyState
			},
			ready: true,
		},
		{
			name: "create replica cli failure",
			prepare: func(r *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			ready: false,
		},
		{
			name: "get replica stream err",
			prepare: func(r *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			ready: false,
		},
		{
			name: "get replica stream ack err",
			prepare: func(r *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil)
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			ready: false,
		},
		{
			name: "replica idx == current node",
			prepare: func(r *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil)
				q.EXPECT().HeadSeq().Return(int64(11))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
			},
			ready: true,
		},
		{
			name: "remote replica ack index < current smallest ack, but reset remote replica index err",
			prepare: func(r *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil)
				fq.EXPECT().HeadSeq().Return(int64(10))
				q.EXPECT().HeadSeq().Return(int64(12))
				q.EXPECT().TailSeq().Return(int64(13))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				replicaCli.EXPECT().Reset(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			ready: false,
		},
		{
			name: " remote replica ack index < current smallest ack, reset success",
			prepare: func(r *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil)
				fq.EXPECT().HeadSeq().Return(int64(10))
				q.EXPECT().HeadSeq().Return(int64(12))
				q.EXPECT().TailSeq().Return(int64(13))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				replicaCli.EXPECT().Reset(gomock.Any(), gomock.Any()).Return(nil, nil)
				q.EXPECT().SetHeadSeq(int64(11))
			},
			ready: true,
		},
		{
			name: "remote replica ack index > current append index, maybe leader lost data",
			prepare: func(r *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil)
				fq.EXPECT().HeadSeq().Return(int64(5))
				q.EXPECT().HeadSeq().Return(int64(12))
				q.EXPECT().TailSeq().Return(int64(9))
				q.EXPECT().Ack(int64(10))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				fq.EXPECT().SetAppendSeq(int64(11))
				q.EXPECT().SetHeadSeq(int64(10)).Return(fmt.Errorf("err"))
				q.EXPECT().HeadSeq().Return(int64(11))
			},
			ready: true,
		},
		{
			name: "remote replica ack index > current append index, maybe leader lost data, reset replica index failure",
			prepare: func(r *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil)
				fq.EXPECT().HeadSeq().Return(int64(5))
				q.EXPECT().HeadSeq().Return(int64(12))
				q.EXPECT().TailSeq().Return(int64(9))
				q.EXPECT().Ack(int64(10))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				fq.EXPECT().SetAppendSeq(int64(11))
				q.EXPECT().SetHeadSeq(int64(10)).Return(fmt.Errorf("err"))
				q.EXPECT().HeadSeq().Return(int64(12))
			},
			ready: false,
		},
		{
			name: "reconnect after failure",
			prepare: func(r *remoteReplicator) {
				r.state = ReplicatorFailureState
				stream := protoReplicaV1.NewMockReplicaService_ReplicaClient(ctrl)
				stream.EXPECT().CloseSend().Return(fmt.Errorf("err"))
				r.replicaStream = stream

				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil)
				fq.EXPECT().HeadSeq().Return(int64(5))
				q.EXPECT().HeadSeq().Return(int64(12))
				q.EXPECT().TailSeq().Return(int64(9))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				fq.EXPECT().SetAppendSeq(int64(11))
				q.EXPECT().SetHeadSeq(int64(10)).Return(fmt.Errorf("err"))
				q.EXPECT().Ack(int64(10))
				q.EXPECT().HeadSeq().Return(int64(11))
			},
			ready: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := NewRemoteReplicator(context.TODO(), rc, stateMgr, cliFct)
			r1 := r.(*remoteReplicator)
			if tt.prepare != nil {
				tt.prepare(r1)
			}
			ready := r.IsReady()
			assert.Equal(t, tt.ready, ready)
		})
	}
}

func TestRemoteReplicator_NodeStateChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	cliFct := rpc.NewMockClientStreamFactory(ctrl)
	stateMgr := storage.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchNodeStateChangeEvent(gomock.Any(), gomock.Any()).AnyTimes()
	q := queue.NewMockFanOut(ctrl)
	fq := queue.NewMockFanOutQueue(ctrl)
	q.EXPECT().Queue().Return(fq).AnyTimes()
	q.EXPECT().HeadSeq().Return(int64(11)).AnyTimes()
	replicaCli := protoReplicaV1.NewMockReplicaServiceClient(ctrl)
	cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil).AnyTimes()
	replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil).AnyTimes()
	replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
		AckIndex: 10,
	}, nil).AnyTimes()
	replicaCli.EXPECT().Reset(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	rc := &ReplicatorChannel{
		State: &models.ReplicaState{
			Database: "test",
			ShardID:  0,
			Leader:   1,
			Follower: 2,
		},
		Queue: q,
	}

	r := NewRemoteReplicator(context.TODO(), rc, stateMgr, cliFct)
	// case 1: node ready
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, true)
	assert.True(t, r.IsReady())
	// case 2: node online->offline
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, false)
	r1 := r.(*remoteReplicator)
	r1.rwMutex.Lock()
	r1.state = ReplicatorInitState
	r1.rwMutex.Unlock()
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		// mock node online async after 100ms
		time.Sleep(100 * time.Millisecond)
		r1.handleNodeStateChangeEvent(models.NodeOnline)
		wait.Done()
	}()
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, true)
	assert.True(t, r.IsReady()) // wait node online
	wait.Wait()
}

func TestRemoteReplicator_Replica(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	cliFct := rpc.NewMockClientStreamFactory(ctrl)
	stateMgr := storage.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchNodeStateChangeEvent(gomock.Any(), gomock.Any()).AnyTimes()
	q := queue.NewMockFanOut(ctrl)
	rc := &ReplicatorChannel{
		State: &models.ReplicaState{
			Database: "test",
			ShardID:  0,
			Leader:   1,
			Follower: 2,
		},
		Queue: q,
	}

	r := NewRemoteReplicator(context.TODO(), rc, stateMgr, cliFct)
	r1 := r.(*remoteReplicator)
	cli := protoReplicaV1.NewMockReplicaService_ReplicaClient(ctrl)
	r1.replicaStream = cli

	cli.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	r.Replica(1, []byte{})

	cli.EXPECT().Send(gomock.Any()).Return(nil)
	cli.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	r.Replica(1, []byte{})

	cli.EXPECT().Send(gomock.Any()).Return(nil)
	cli.EXPECT().Recv().Return(&protoReplicaV1.ReplicaResponse{
		AckIndex:     1,
		ReplicaIndex: 1,
	}, nil)
	q.EXPECT().Ack(int64(1))
	r.Replica(1, []byte{})
}
