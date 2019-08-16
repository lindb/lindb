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
	"github.com/lindb/lindb/pkg/state"
)

func TestStatusStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	factory := discovery.NewMockFactory(ctrl)
	factory.EXPECT().GetRepo().Return(repo).AnyTimes()
	discovery1 := discovery.NewMockDiscovery(ctrl)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	_, err := NewStatusStateMachine(context.TODO(), factory)
	assert.NotNil(t, err)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err = NewStatusStateMachine(context.TODO(), factory)
	assert.NotNil(t, err)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Key: "key", Value: []byte{1, 2, 3}}}, nil)
	discovery1.EXPECT().Discovery().Return(nil)
	sm, err := NewStatusStateMachine(context.TODO(), factory)
	if err != nil {
		t.Fatal(err)
	}
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
	if err != nil {
		t.Fatal(err)
	}
}
