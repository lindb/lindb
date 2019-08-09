package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplicaState(t *testing.T) {
	replica := ReplicaState{
		Cluster:      "cluster",
		Database:     "db",
		ShardID:      int32(1),
		WriteIndex:   100,
		ReplicaIndex: 50,
	}

	assert.Equal(t, "cluster/db/1", replica.ShardIndicator())
	assert.Equal(t, int64(50), replica.Pending())
}
