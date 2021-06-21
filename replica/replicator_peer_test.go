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
	"time"

	"github.com/golang/mock/gomock"
)

func TestReplicatorPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	replicator := NewMockReplicator(ctrl)
	replicator.EXPECT().IsReady().Return(false).AnyTimes()
	replicator.EXPECT().String().Return("str").AnyTimes()
	peer := NewReplicatorPeer(replicator)
	peer.Startup()
	peer.Startup()
	time.Sleep(10 * time.Millisecond)
	peer.Shutdown()
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

	// loop 1: no data
	replicator.EXPECT().IsReady().Return(true)
	replicator.EXPECT().Consume().Return(int64(0)) //no data
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
	replicator.EXPECT().IsReady().Return(false).AnyTimes()
	peer := NewReplicatorPeer(replicator)
	peer.Startup()
	peer.Shutdown()
	time.Sleep(100 * time.Millisecond)
}
