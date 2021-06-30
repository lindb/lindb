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

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/replication"
)

func TestReplicatorStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replication.NewMockChannelManager(ctrl)
	cm.EXPECT().SyncReplicatorState().AnyTimes()
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	shardAssign := models.NewShardAssignment("test")
	shardAssign.Nodes[1] = &models.Node{IP: "1.1.1.1", Port: 9000}
	shardAssign.AddReplica(1, 1)

	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	sm, err := NewReplicatorStateMachine(context.TODO(), cm, discoveryFactory)
	assert.Error(t, err)
	assert.Nil(t, sm)

	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	sm, err = NewReplicatorStateMachine(context.TODO(), cm, discoveryFactory)
	assert.NoError(t, err)
	assert.NotNil(t, sm)

	data := encoding.JSONMarshal(shardAssign)
	cm.EXPECT().CreateChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	sm.OnCreate("/test/path", data)

	// test on create event
	sm.OnCreate("/test/path", []byte{1, 2, 3})
	ch := replication.NewMockChannel(ctrl)
	cm.EXPECT().CreateChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(ch, nil)
	ch.EXPECT().GetOrCreateReplicator(gomock.Any()).Return(nil, fmt.Errorf("err"))
	sm.OnCreate("/test/path", data)

	cm.EXPECT().CreateChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(ch, nil)
	ch.EXPECT().GetOrCreateReplicator(gomock.Any()).Return(nil, nil)
	sm.OnCreate("/test/path", data)

	s := sm.(*replicatorStateMachine)
	assert.Equal(t, 1, len(s.shardAssigns))
	assert.NotNil(t, s.shardAssigns["test"])

	sm.OnDelete("/shard/test")
	assert.Equal(t, 0, len(s.shardAssigns))

	discovery1.EXPECT().Close()
	err = sm.Close()
	assert.NoError(t, err)

	err = sm.Close()
	assert.NoError(t, err)
}
