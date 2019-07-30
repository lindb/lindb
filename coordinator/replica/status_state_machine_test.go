package replica

import (
	"context"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/constants"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/pathutil"
	"github.com/eleme/lindb/pkg/state"
)

func TestStatusStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)
	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().WatchPrefix(gomock.Any(), constants.ReplicaStatePath).Return(eventCh)

	sm, err := NewStatusStateMachine(context.TODO(), repo)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, sm)

	replicaStatus := []models.ReplicaState{{
		Cluster:  "test",
		Database: "11",
	}}
	data, _ := json.Marshal(models.BrokerReplicaState{Replicas: replicaStatus})

	// wrong event
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: pathutil.GetReplicaStatePath("1.1.1.1:2080"), Value: nil},
		},
	})
	assert.Equal(t, 0, len(sm.GetReplicas("1.1.1.1:2080").Replicas))
	// modify event
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: pathutil.GetReplicaStatePath("1.1.1.1:2080"), Value: data},
		},
	})
	assert.Equal(t, replicaStatus, sm.GetReplicas("1.1.1.1:2080").Replicas)
	assert.Equal(t, 0, len(sm.GetReplicas("1.1.1.2:2080").Replicas))
	// delete event
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeDelete,
		KeyValues: []state.EventKeyValue{
			{Key: pathutil.GetReplicaStatePath("1.1.1.1:2080")},
		},
	})
	assert.Equal(t, 0, len(sm.GetReplicas("1.1.1.1:2080").Replicas))
}

func TestStatusStateMachine_GetQueryableReplicas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)
	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().WatchPrefix(gomock.Any(), constants.ReplicaStatePath).Return(eventCh)

	sm, err := NewStatusStateMachine(context.TODO(), repo)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(sm.GetQueryableReplicas("test")))

	// broker 1:
	replicaStatus := []models.ReplicaState{
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
	data, _ := json.Marshal(models.BrokerReplicaState{Replicas: replicaStatus})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: pathutil.GetReplicaStatePath("2.1.1.1:2080"), Value: data},
		},
	})

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
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: pathutil.GetReplicaStatePath("2.1.1.2:2080"), Value: data},
		},
	})

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
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(10 * time.Millisecond)
}
