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
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
)

func TestReplicatorPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	mockReplicator := NewMockReplicator(ctrl)
	mockReplicator.EXPECT().IsReady().Return(false).AnyTimes()
	mockReplicator.EXPECT().String().Return("str").AnyTimes()
	mockReplicator.EXPECT().State().Return(&models.ReplicaState{}).AnyTimes()
	peer := NewReplicatorPeer(mockReplicator)
	peer.Startup()
	peer.Startup()
	time.Sleep(10 * time.Millisecond)
	peer.Shutdown()
	peer.Shutdown()
	time.Sleep(10 * time.Millisecond)

	ch := make(chan struct{})
	peer = &replicatorPeer{
		runner: &replicatorRunner{
			replicator: &remoteReplicator{
				replicator: replicator{
					channel: &ReplicatorChannel{
						State: &models.ReplicaState{},
					},
				},
			},
			running: atomic.NewBool(true),
			closed:  ch,
		},
		running: atomic.NewBool(true),
	}
	go func() {
		ch <- struct{}{}
	}()
	peer.Shutdown()
	time.Sleep(10 * time.Millisecond)
}

func TestNewReplicator_runner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	replicator := NewMockReplicator(ctrl)
	replicator.EXPECT().String().Return("str").AnyTimes()
	replicator.EXPECT().State().Return(&models.ReplicaState{}).AnyTimes()
	replicator.EXPECT().Pending().Return(int64(19)).AnyTimes()
	peer := NewReplicatorPeer(replicator)
	var wait sync.WaitGroup

	// loop 1: no data
	replicator.EXPECT().IsReady().Return(true)
	replicator.EXPECT().Consume().Return(int64(-1)) // no data
	// loop 2: get message err
	replicator.EXPECT().IsReady().Return(true)
	replicator.EXPECT().Consume().Return(int64(1))                          // has data
	replicator.EXPECT().GetMessage(int64(1)).Return(nil, fmt.Errorf("err")) // get message err
	// loop 3: do replica
	replicator.EXPECT().IsReady().Return(true)
	replicator.EXPECT().Consume().Return(int64(1))            // has data
	replicator.EXPECT().GetMessage(int64(1)).Return(nil, nil) // get message
	replicator.EXPECT().Replica(gomock.Any(), gomock.Any())   // replica
	// other loop
	replicator.EXPECT().IsReady().DoAndReturn(func() bool {
		wait.Done()
		peer.Shutdown()
		return false
	}).AnyTimes()
	peer.Startup()
	wait.Add(1)
	wait.Wait()
}

func TestReplicatorPeer_newReplicatorRunner(t *testing.T) {
	rc := &ReplicatorChannel{
		State: &models.ReplicaState{
			Database: "test",
			ShardID:  models.ShardID(1),
		},
	}
	assert.NotNil(t, newReplicatorRunner(&localReplicator{
		replicator: replicator{
			channel: rc,
		},
	}))
	assert.NotNil(t, newReplicatorRunner(&remoteReplicator{
		replicator: replicator{
			channel: rc,
		},
	}))
}

func TestReplicatorPeer_replica_panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	replicator := NewMockReplicator(ctrl)
	r := &replicatorRunner{
		replicator: replicator,
	}
	replicator.EXPECT().IsReady().DoAndReturn(func() bool {
		panic("err")
	})
	assert.Panics(t, r.replica)
}
