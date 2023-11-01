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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
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

	cases := []struct {
		name    string
		prepare func(r *remoteReplicator)
		ready   bool
	}{
		{
			name: "replicator is ready",
			prepare: func(r *remoteReplicator) {
				r.state.Store(&state{state: models.ReplicatorReadyState})
			},
			ready: true,
		},
		{
			name: "create replica cli failure",
			prepare: func(_ *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			ready: false,
		},
		{
			name: "get replica stream ack err",
			prepare: func(_ *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			ready: false,
		},
		{
			name: "replica idx == current node",
			prepare: func(_ *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				cg.EXPECT().ConsumedSeq().Return(int64(10))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
			},
			ready: true,
		},
		{
			name: "remote replica ack index < current smallest ack, but reset remote replica index err",
			prepare: func(_ *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				q.EXPECT().AppendedSeq().Return(int64(10))
				cg.EXPECT().ConsumedSeq().Return(int64(12))
				cg.EXPECT().AcknowledgedSeq().Return(int64(13))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				replicaCli.EXPECT().Reset(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			ready: false,
		},
		{
			name: " remote replica ack index < current smallest ack, reset success",
			prepare: func(_ *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				q.EXPECT().AppendedSeq().Return(int64(10))
				cg.EXPECT().ConsumedSeq().Return(int64(7))
				cg.EXPECT().AcknowledgedSeq().Return(int64(8))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 5,
				}, nil)
				replicaCli.EXPECT().Reset(gomock.Any(), gomock.Any()).Return(nil, nil)
				cg.EXPECT().SetConsumedSeq(int64(8))
			},
			ready: true,
		},
		{
			name: "remote replica ack index > current append index, maybe leader lost data",
			prepare: func(_ *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				q.EXPECT().AppendedSeq().Return(int64(5))
				cg.EXPECT().ConsumedSeq().Return(int64(12))
				cg.EXPECT().AcknowledgedSeq().Return(int64(9))
				cg.EXPECT().Ack(int64(10))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				fq.EXPECT().SetAppendedSeq(int64(10))
				cg.EXPECT().SetConsumedSeq(int64(10))
				cg.EXPECT().ConsumedSeq().Return(int64(10))
			},
			ready: true,
		},
		{
			name: "remote replica ack index > current append index, maybe leader lost data, reset replica index failure",
			prepare: func(_ *remoteReplicator) {
				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				q.EXPECT().AppendedSeq().Return(int64(5))
				cg.EXPECT().ConsumedSeq().Return(int64(12))
				cg.EXPECT().AcknowledgedSeq().Return(int64(9))
				cg.EXPECT().Ack(int64(10))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				fq.EXPECT().SetAppendedSeq(int64(10))
				cg.EXPECT().SetConsumedSeq(int64(10))
				cg.EXPECT().ConsumedSeq().Return(int64(12))
			},
			ready: false,
		},
		{
			name: "reconnect after failure",
			prepare: func(r *remoteReplicator) {
				r.state.Store(&state{state: models.ReplicatorFailureState})
				stream := protoReplicaV1.NewMockReplicaService_ReplicaClient(ctrl)
				stream.EXPECT().CloseSend().Return(fmt.Errorf("err"))
				r.replicaStream = stream

				cliFct.EXPECT().CreateReplicaServiceClient(gomock.Any()).Return(replicaCli, nil)
				q.EXPECT().AppendedSeq().Return(int64(5))
				cg.EXPECT().ConsumedSeq().Return(int64(12))
				cg.EXPECT().AcknowledgedSeq().Return(int64(9))
				cg.EXPECT().Ack(int64(10))
				replicaCli.EXPECT().GetReplicaAckIndex(gomock.Any(), gomock.Any()).Return(&protoReplicaV1.GetReplicaAckIndexResponse{
					AckIndex: 10,
				}, nil)
				fq.EXPECT().SetAppendedSeq(int64(10))
				cg.EXPECT().SetConsumedSeq(int64(10))
				cg.EXPECT().ConsumedSeq().Return(int64(10))
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
	cg := queue.NewMockConsumerGroup(ctrl)
	q := queue.NewMockQueue(ctrl)
	fq := queue.NewMockFanOutQueue(ctrl)
	fq.EXPECT().Queue().Return(q).AnyTimes()
	cg.EXPECT().Queue().Return(fq).AnyTimes()
	cg.EXPECT().ConsumedSeq().Return(int64(10)).AnyTimes()
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
		ConsumerGroup: cg,
	}

	r := NewRemoteReplicator(context.TODO(), rc, stateMgr, cliFct)
	// case 1: node ready
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, true)
	assert.True(t, r.IsReady())
	// case 2: node online->offline
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, false)
	r1 := r.(*remoteReplicator)
	r1.rwMutex.Lock()
	r1.state.Store(&state{state: models.ReplicatorInitState})
	r1.rwMutex.Unlock()
	s := r1.State()
	assert.Equal(t, state{state: models.ReplicatorInitState}, *s)
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
	q := queue.NewMockConsumerGroup(ctrl)
	rc := &ReplicatorChannel{
		State: &models.ReplicaState{
			Database: "test",
			ShardID:  0,
			Leader:   1,
			Follower: 2,
		},
		ConsumerGroup: q,
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
	// invalid ack sequence
	cli.EXPECT().Send(gomock.Any()).Return(nil)
	cli.EXPECT().Recv().Return(&protoReplicaV1.ReplicaResponse{
		AckIndex:     1,
		ReplicaIndex: 2,
	}, nil)
	r.Replica(1, []byte{})
}

func TestRemoteReplicator_Connect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	stateMgr := storage.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchNodeStateChangeEvent(gomock.Any(), gomock.Any()).AnyTimes()
	stateMgr.EXPECT().GetLiveNode(gomock.Any()).Return(models.StatefulNode{}, true).AnyTimes()
	cliFct := rpc.NewMockClientStreamFactory(ctrl)
	replicaCli := protoReplicaV1.NewMockReplicaServiceClient(ctrl)
	q := queue.NewMockConsumerGroup(ctrl)
	fq := queue.NewMockFanOutQueue(ctrl)
	q.EXPECT().Queue().Return(fq).AnyTimes()
	rc := &ReplicatorChannel{
		State: &models.ReplicaState{
			Database: "test",
			ShardID:  0,
			Leader:   1,
			Follower: 2,
		},
		ConsumerGroup: q,
	}

	cases := []struct {
		name    string
		prepare func(r *remoteReplicator)
		ready   bool
	}{
		{
			name: "stream exist",
			prepare: func(r *remoteReplicator) {
				r.replicaStream = protoReplicaV1.NewMockReplicaService_ReplicaClient(ctrl)
			},
			ready: true,
		},
		{
			name: "create stream failure",
			prepare: func(r *remoteReplicator) {
				r.replicaCli = replicaCli
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			ready: false,
		},
		{
			name: "create stream successfully",
			prepare: func(r *remoteReplicator) {
				r.replicaCli = replicaCli
				replicaCli.EXPECT().Replica(gomock.Any()).Return(nil, nil)
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
			ready := r.Connect()
			assert.Equal(t, tt.ready, ready)
		})
	}
}

func TestRemoteReplicator_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("no stream", func(_ *testing.T) {
		r := &remoteReplicator{}
		r.Close()
	})
	t.Run("close stream failure", func(_ *testing.T) {
		stream := protoReplicaV1.NewMockReplicaService_ReplicaClient(ctrl)
		r := &remoteReplicator{
			replicaStream: stream,
			statistics:    metrics.NewStorageRemoteReplicatorStatistics("test", "0"),
			logger:        logger.GetLogger("Test", "RemoteReplicator"),
			replicator: replicator{channel: &ReplicatorChannel{
				State: &models.ReplicaState{
					Database: "test",
					ShardID:  0,
					Leader:   1,
					Follower: 2,
				},
			}},
		}
		stream.EXPECT().CloseSend().Return(fmt.Errorf("err"))
		r.Close()
	})
}
