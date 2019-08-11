package replica

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

	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err := NewStatusStateMachine(context.TODO(), factory)
	assert.NotNil(t, err)

	discovery1.EXPECT().Discovery().Return(nil)
	sm, err := NewStatusStateMachine(context.TODO(), factory)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, sm)

	sm.OnCreate("/data/err1", []byte{1, 1, 3})

	replicaStatus := []models.ReplicaState{{
		Cluster:  "test",
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
			Cluster:      "test",
			Database:     "test_db_2",
			TO:           models.Node{IP: "1.1.1.2", Port: 2090},
			WriteIndex:   100,
			ReplicaIndex: 50,
			ShardID:      1,
		},
		{
			Cluster:      "test",
			Database:     "test_db",
			TO:           models.Node{IP: "1.1.1.2", Port: 2090},
			WriteIndex:   100,
			ReplicaIndex: 50,
			ShardID:      1,
		},
		{
			Cluster:      "test",
			Database:     "test_db",
			TO:           models.Node{IP: "1.1.1.3", Port: 2090},
			WriteIndex:   100,
			ReplicaIndex: 90,
			ShardID:      1,
		},
	}
	data, _ = json.Marshal(models.BrokerReplicaState{Replicas: replicaStatus})
	sm.OnCreate("/broker/2.1.1.1:2080", data)

	// broker 2:
	replicaStatus = []models.ReplicaState{
		{
			Cluster:      "test",
			Database:     "test_db_2",
			TO:           models.Node{IP: "1.1.1.2", Port: 2090},
			WriteIndex:   100,
			ReplicaIndex: 50,
			ShardID:      2,
		},
		{
			Cluster:      "test",
			Database:     "test_db",
			TO:           models.Node{IP: "1.1.1.2", Port: 2090},
			WriteIndex:   100,
			ReplicaIndex: 50,
			ShardID:      2,
		},
		{
			Cluster:      "test",
			Database:     "test_db",
			TO:           models.Node{IP: "1.1.1.3", Port: 2090},
			WriteIndex:   100,
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
}
