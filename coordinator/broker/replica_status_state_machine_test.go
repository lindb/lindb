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

package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
)

func TestStatusStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	discovery1.EXPECT().Discovery(gomock.Any()).Return(fmt.Errorf("err"))
	_, err := NewReplicaStatusStateMachine(context.TODO(), factory)
	assert.Error(t, err)

	discovery1.EXPECT().Discovery(gomock.Any()).Return(nil)
	sm, err := NewReplicaStatusStateMachine(context.TODO(), factory)
	assert.NoError(t, err)
	assert.NotNil(t, sm)

	sm.OnCreate("/data/err1", []byte{1, 1, 3})

	replicaStatus := []models.ReplicaState{{
		Database: "11",
	}}
	brokerReplicaState := models.BrokerReplicaState{Replicas: replicaStatus}

	data, _ := json.Marshal(&brokerReplicaState)
	sm.OnCreate("/data/1.1.1.1:9000", data)
	assert.Equal(t, brokerReplicaState, sm.GetReplicas("1.1.1.1:9000"))

	sm.OnDelete("/data/1.1.1.1:9000")
	assert.Equal(t, 0, len(sm.GetReplicas("1.1.1.1:9000").Replicas))

	// broker 1:
	replicaStatus = []models.ReplicaState{
		{
			Database:     "test_db_2",
			Target:       models.Node{IP: "1.1.1.2", Port: 2090},
			Pending:      50,
			ReplicaIndex: 50,
			ShardID:      1,
		},
		{
			Database:     "test_db",
			Target:       models.Node{IP: "1.1.1.2", Port: 2090},
			Pending:      50,
			ReplicaIndex: 50,
			ShardID:      1,
		},
		{
			Database:     "test_db",
			Target:       models.Node{IP: "1.1.1.3", Port: 2090},
			Pending:      10,
			ReplicaIndex: 90,
			ShardID:      1,
		},
	}
	data, _ = json.Marshal(models.BrokerReplicaState{Replicas: replicaStatus})
	sm.OnCreate("/broker/2.1.1.1:2080", data)

	// broker 2:
	replicaStatus = []models.ReplicaState{
		{
			Database:     "test_db_2",
			Target:       models.Node{IP: "1.1.1.2", Port: 2090},
			Pending:      50,
			ReplicaIndex: 50,
			ShardID:      2,
		},
		{
			Database:     "test_db",
			Target:       models.Node{IP: "1.1.1.2", Port: 2090},
			Pending:      50,
			ReplicaIndex: 50,
			ShardID:      2,
		},
		{
			Database:     "test_db",
			Target:       models.Node{IP: "1.1.1.3", Port: 2090},
			Pending:      10,
			ReplicaIndex: 90,
			ShardID:      2,
		},
	}
	data, _ = json.Marshal(models.BrokerReplicaState{Replicas: replicaStatus})
	sm.OnCreate("/broker/2.1.1.2:2080", data)

	r := sm.GetQueryableReplicas("test_db")
	assert.Equal(t, 1, len(r))
	shards := r["1.1.1.3:2090"]
	sort.Slice(shards, func(i, j int) bool {
		return shards[i] < shards[j]
	})
	assert.Equal(t, []int32{1, 2}, shards)

	r = sm.GetQueryableReplicas("test_db_2")
	assert.Equal(t, 1, len(r))
	shards = r["1.1.1.2:2090"]
	sort.Slice(shards, func(i, j int) bool {
		return shards[i] < shards[j]
	})
	assert.Equal(t, []int32{1, 2}, shards)

	r = sm.GetQueryableReplicas("test_db_not_exist")
	assert.Nil(t, r)

	discovery1.EXPECT().Close()
	err = sm.Close()
	assert.NoError(t, err)

	err = sm.Close()
	assert.NoError(t, err)

	// after close, get empty data
	assert.Nil(t, sm.GetQueryableReplicas("test_db_2"))
	assert.Equal(t, models.BrokerReplicaState{}, sm.GetReplicas("1.1.1.1:9000"))
}
