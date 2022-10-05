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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestReplicator_String(t *testing.T) {
	time, _ := timeutil.ParseTimestamp("2019-12-12 10:11:10")
	r := replicator{channel: &ReplicatorChannel{
		State: &models.ReplicaState{
			Database:   "test",
			ShardID:    1,
			Leader:     1,
			Follower:   2,
			FamilyTime: time,
		},
	}}

	assert.Equal(t, "[database:test,shard:1,family:2019-12-12 10:11:10,from(leader):1,to(follower):2]", r.String())
}

func TestReplicator_Base(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cg := queue.NewMockConsumerGroup(ctrl)
	fanoutQ := queue.NewMockFanOutQueue(ctrl)
	cg.EXPECT().Queue().Return(fanoutQ).AnyTimes()
	q := queue.NewMockQueue(ctrl)
	fanoutQ.EXPECT().Queue().Return(q).AnyTimes()
	state := &models.ReplicaState{
		Database:   "test",
		ShardID:    1,
		Leader:     1,
		Follower:   2,
		FamilyTime: 1,
	}
	r := replicator{
		channel: &ReplicatorChannel{
			State:         state,
			ConsumerGroup: cg,
		},
	}
	assert.Equal(t, state, r.ReplicaState())
	r.Replica(0, []byte{1, 2, 3})
	assert.True(t, r.IsReady())
	assert.True(t, r.Connect())
	cg.EXPECT().Consume().Return(int64(10))
	assert.Equal(t, int64(10), r.Consume())
	q.EXPECT().Get(int64(10)).Return([]byte{1, 2, 3}, nil)
	rs, err := r.GetMessage(10)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, rs)
	cg.EXPECT().ConsumedSeq().Return(int64(10))
	assert.Equal(t, int64(11), r.ReplicaIndex())
	cg.EXPECT().AcknowledgedSeq().Return(int64(10))
	assert.Equal(t, int64(10), r.AckIndex())

	q.EXPECT().AppendedSeq().Return(int64(10))
	assert.Equal(t, int64(11), r.AppendIndex())
	fanoutQ.EXPECT().SetAppendedSeq(int64(10))
	r.ResetAppendIndex(int64(11))

	cg.EXPECT().Ack(int64(10))
	r.SetAckIndex(int64(10))

	cg.EXPECT().SetConsumedSeq(int64(9))
	r.ResetReplicaIndex(int64(10))
	cg.EXPECT().Pending().Return(int64(10))
	assert.Equal(t, int64(10), r.Pending())

	r.Close()
}
